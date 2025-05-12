package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	osc "github.com/outscale/osc-sdk-go/v2"
	cli "github.com/teris-io/cli"
)

const (
	defaultFetchInterval       = 20 // seconds
	firstPageSize        int32 = 25
	nextPageSize         int32 = 1000
	oscLogsVersion             = "v0.1.4"
)

func displayLogs(args []string, options map[string]string) int {
	if options["version"] == "true" {
		fmt.Println(oscLogsVersion)
		return 0
	}

	var err error
	var ctx context.Context
	var client osc.APIClient
	var lastRequestId string
	if options["profile"] != "" {
		_, ctx, client, err = GenerateConfigurationAndContext(options["profile"])
	} else {
		_, ctx, client, err = GenerateConfigurationAndContext("default")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot load credentials: %v\n", err)
		return 1
	}

	countValue := -1
	if options["count"] != "" {
		countValue, err = strconv.Atoi(options["count"])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: cannot convert --count option into integer")
			return 1
		}
	}

	duration := time.Duration(defaultFetchInterval) * time.Second
	if options["interval"] != "" {
		intervalValue, err := strconv.Atoi(options["interval"])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: cannot convert --interval option to integer")
			return 1
		}
		if intervalValue < 1 {
			fmt.Fprintln(os.Stderr, "the interval must be greater than 0")
			return 1
		} else {
			duration = time.Duration(intervalValue) * time.Second
		}
	}

	callsToIgnore := []string{"ReadApiLogs"}
	if options["ignore"] != "" {
		callsToIgnore = strings.Split(options["ignore"], ",")
	}

	var output io.Writer
	if options["write"] != "" {
		file, err := os.OpenFile(options["write"], os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: cannot open output file: %v\n", err)
			return 1
		}
		defer file.Close()
		output = file
	} else {
		output = os.Stdout
	}
	jsonOut := json.NewEncoder(output)
	tk := time.NewTicker(duration)
	var logDate *string
	logcount := 0
	pageSize := firstPageSize
	for range tk.C {
		req := osc.ReadApiLogsRequest{
			Filters: &osc.FiltersApiLog{
				QueryDateAfter: logDate,
			},
			ResultsPerPage: &pageSize,
		}
		resp, _, err := client.ApiLogApi.ReadApiLogs(ctx).ReadApiLogsRequest(req).Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %v\n", err)
			continue
		}
		logs := resp.GetLogs()
		if len(logs) == 0 || (len(logs) == 1 && logs[0].GetRequestId() == lastRequestId) {
			fmt.Print(".")
			continue
		}
		for _, log := range slices.Backward(logs) {
			if log.GetRequestId() == lastRequestId {
				continue
			}
			if slices.Contains(callsToIgnore, log.GetQueryCallName()) {
				continue
			}
			_, err := fmt.Fprintln(output)
			if err == nil {
				err = jsonOut.Encode(log)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot write log: %v\n", err)
				return 1
			}
			logcount++
			if logcount == countValue {
				return 0
			}
		}

		// logs sorted by time, most recent first
		lastLog := logs[0]
		lastRequestId = lastLog.GetRequestId()
		logDate = lastLog.QueryDate
		pageSize = nextPageSize
	}
	return 0
}
func AddWriteOption() cli.Option {
	return cli.NewOption("write", "Write all traces inside a file instead of writing to standard output").WithChar('w').WithType(cli.TypeString)
}
func AddCountOption() cli.Option {
	return cli.NewOption("count", "Exit after <count> logs").WithChar('c').WithType(cli.TypeInt)
}
func AddIntervalOption() cli.Option {
	text := fmt.Sprintf("Wait a duration defined by <wait> (in seconds) between two calls to Outscale API (default: %d)", defaultFetchInterval)
	return cli.NewOption("interval", text).WithChar('i').WithType(cli.TypeInt)
}
func AddProfileOption() cli.Option {
	return cli.NewOption("profile", "Use a specific profile name (\"default\" is the default profile )").WithChar('p').WithType(cli.TypeString)
}
func AddIgnoreOption() cli.Option {
	return cli.NewOption("ignore", "Ignore one or more specific API calls. Values are separated by commas e.g. \"--ignore=ReadApiLogs,ReadVms\"").WithChar('I').WithType(cli.TypeString)
}
func AddVersionOption() cli.Option {
	return cli.NewOption("version", "Print version to standard output and exit").WithChar('v').WithType(cli.TypeBool)
}
func GenerateConfigurationAndContext(profileName string) (*osc.Configuration, context.Context, osc.APIClient, error) {
	ctx := context.Background()
	configFile, err := osc.LoadDefaultConfigFile()
	if err != nil {
		return nil, ctx, osc.APIClient{}, err
	}
	config, err := configFile.Configuration(profileName)
	if err != nil {
		return nil, ctx, osc.APIClient{}, err
	}
	config.Debug = false
	ctx, err = configFile.Context(ctx, profileName)
	if err != nil {
		return nil, ctx, osc.APIClient{}, err
	}
	client := *osc.NewAPIClient(config)
	return config, ctx, client, err
}

func main() {
	app := cli.New("osc-logs").
		WithAction(displayLogs).
		WithOption(AddWriteOption()).
		WithOption(AddCountOption()).
		WithOption(AddIntervalOption()).
		WithOption(AddProfileOption()).
		WithOption(AddIgnoreOption()).
		WithOption(AddVersionOption())
	ret := app.Run(os.Args, os.Stdout)
	os.Exit(ret)
}
