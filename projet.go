package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	osc "github.com/outscale/osc-sdk-go/v2"
	cli "github.com/teris-io/cli"
)

func cliLog() cli.Command {
	return cli.NewCommand("apilog", "Read Api logs").
		WithAction(func(args []string, options map[string]string) int {
			config := osc.NewConfiguration()
			config.Debug = true

			client := osc.NewAPIClient(config)

			ctx := context.WithValue(context.Background(), osc.ContextAWSv4, osc.AWSv4{
				AccessKey: os.Getenv("OSC_ACCESS_KEY"),
				SecretKey: os.Getenv("OSC_SECRET_KEY"),
			})
			req := osc.ReadApiLogsRequest{}
			resp, httpRes, err := client.ApiLogApi.ReadApiLogs(ctx).ReadApiLogsRequest(req).Execute()
			if err != nil {
				fmt.Printf("Error %v", err)
				if httpRes != nil {
					fmt.Fprintln(os.Stderr, httpRes.Status)
				}
				os.Exit(1)
			}

			for _, log := range resp.GetLogs() {
				a, _ := json.Marshal(log)
				fmt.Println(string(a))
			}

			return 0
		})
}

func LogW() cli.Option {
	return cli.NewOption("write", "Write the API logs in a file").WithChar('w').WithType(cli.TypeString)
}

func LogC() cli.Option {
	return cli.NewOption("count", "exit once <count> logs are written").WithChar('c').WithType(cli.TypeInt)
}

func LogI() cli.Option {
	return cli.NewOption("i", "waits a duration defined by <wait> (in seconds) between two calls to Outscale AP ").WithChar('i').WithType(cli.TypeInt)
}

func LogP() cli.Option {
	return cli.NewOption("profil", "use a specific profile name ").WithChar('p').WithType(cli.TypeString)
}
func LogH() cli.Option {
	return cli.NewOption("help", "displays help presenting a description of the program, its version, the options available with their description ").WithChar('h').WithType(cli.TypeString)
}

func main() {

	app := cli.New("OSC-LOG").WithCommand(cliLog()).WithOption(LogW()).WithOption(LogC()).WithOption(LogI()).WithOption(LogP()).WithOption(LogH())
	ret := app.Run(os.Args, os.Stdout)
	os.Exit(ret)

}
