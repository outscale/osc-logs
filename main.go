package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	osc "github.com/outscale/osc-sdk-go/v2"
	cli "github.com/teris-io/cli"
)

var (
	defaultFetchInterval = 10
	resultsPerPage int32 = 1000
	oscLogsVersion = "v0.1.1"
)

func displayLogs(args []string, options map[string]string) int {
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT)
	logDate := time.Now().UTC().Format("2006-01-02T15:04:05")
	duration := time.Duration(defaultFetchInterval) * time.Second
	var file *os.File
	lineBreak := []byte("\n")
	logcount := 0
	var countValue int
	var err error
	var ctx context.Context
	var client osc.APIClient
	var callsToIgnore []string
	var lastRequestId string
	if options["profile"] != "" {
		_, ctx, client, err = GenerateConfigurationAndContext(options["profile"])
	} else {
		_, ctx, client, err = GenerateConfigurationAndContext("default")
	}

	if options["write"] != "" {
		file, err = os.OpenFile(options["write"], os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: can not open the file")
			os.Exit(1)
		}
		defer file.Close()
	}
	if options["count"] != "" {
		countValue, err = strconv.Atoi(options["count"])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: cannot convert --count option into integer  ")
			os.Exit(1)
		}
	} else {
		countValue = -1
	}
	if options["interval"] != "" {
		intervalValue, err := strconv.Atoi(options["interval"])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:cannot convert --interval option to integer ")
			os.Exit(1)
		}
		if intervalValue < 1 {
			fmt.Fprintln(os.Stderr, "the interval must be greater than 0")
			os.Exit(1)
		} else {
			duration = time.Duration(intervalValue) * time.Second
		}
	}
	tk := time.NewTicker(duration)
	if options["ignore"] != "" {
		callsToIgnore = strings.Split(options["ignore"], ",")
	}
	if options["version"] == "true" {
		fmt.Println(oscLogsVersion)
		os.Exit(0)
	}
	for range tk.C {
		req := osc.ReadApiLogsRequest{
			Filters: &osc.FiltersApiLog{
				QueryDateAfter: &logDate,
			},
			ResultsPerPage: &resultsPerPage,
		}
		resp, httpRes, err := client.ApiLogApi.ReadApiLogs(ctx).ReadApiLogsRequest(req).Execute()
		if err != nil {
			fmt.Printf("Error %v", err)
			if httpRes != nil {
				fmt.Fprintln(os.Stderr, httpRes.Status)
			}
			return 1
		}
		logs := resp.GetLogs()
		if !resp.HasLogs() || len(logs) == 0 {
			continue
		}
		for _, log := range logs {
			if log.GetRequestId() == lastRequestId {
				continue
			}
			if SearchByCallName(log, callsToIgnore) {
				continue
			}
			jsonLog, marshalError := json.Marshal(log)
			if marshalError != nil {
				fmt.Fprintf(os.Stderr, "Error: can not read log output")
				return 1
			}
			logcount = logcount + 1
			if file == nil {
				fmt.Println(string(jsonLog))
			} else {
				logWriting := []byte(string(jsonLog))
				file.Write(logWriting)
				file.Write(lineBreak)
			}
			if logcount == countValue {
				os.Exit(0)
			}
		}

		lastLog := logs[len(logs)-1]
		lastRequestId = lastLog.GetRequestId()
		logDate = *lastLog.QueryDate
		go func() {
			<-stopSignal
			os.Exit(0)
		}()
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
	configFile, err := osc.LoadDefaultConfigFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while loading default configuration file: %s \n", err.Error())
		os.Exit(1)
	}
	config, err := configFile.Configuration(profileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating configuration: %s \n", err.Error())
		os.Exit(1)
	}
	config.Debug = false
	ctx, err := configFile.Context(context.Background(), profileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while creating context: %s \n", err.Error())
		os.Exit(1)
	}
	client := *osc.NewAPIClient(config)
	return config, ctx, client, err
}
func SearchByCallName(log osc.Log, callsToIgnore []string) bool {
	LogCallName := log.QueryCallName
	for _, CallName := range callsToIgnore {
		if *LogCallName == CallName {
			return true
		}
	}
	return false
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
