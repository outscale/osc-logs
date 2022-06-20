# osc-logs

osc-logs constantly fetch API call logs from Outscale to easily consult and keep them.

# Features

osc-logs fetch API call logs every few seconds and only print logs from program's start date.

By default logs are printed as Line-delimited JSON to standard output. See [format documentation](https://docs.outscale.com/api#tocslog) for more details.

We recommend using [jq](https://stedolan.github.io/jq/) utility for additional filtering and formating.

```
osc-logs

Options:
    -w, --write      Write all traces inside a file instead of writing to standard output
    -c, --count      Exit after <count> logs
    -i, --interval   Wait a duration defined by <wait> (in seconds) between two calls to Outscale API (default: 10)
    -p, --profile    Use a specific profile name ("default" is the default profile )
    -I, --ignore     Ignore one or more specific API calls. Values are separated by commas e.g. "--ignore=ReadApiLogs,ReadVms"
    -v, --version    Print version to standard output and exit
```

# Installation

Download latest binary in [Release page](https://github.com/outscale/osc-logs/releases).

Alternatively, you can run (or update) osc-logs with this command:
```
go install github.com/outscale/osc-logs@latest
```

# Configuration

osc-logs reads `~/.osc/config.json` file to get its credentials and region details.

Example of `config.json`:
```
{
    "default": {
        "access_key": "MyAccessKey",
        "secret_key": "MySecretKey",
        "region": "eu-west-2"
    }
}
```

# Usage example

Example of storing all logs in a file  except [ReadApiLogs](https://docs.outscale.com/api#readapilogs) call itself:
```
osc-logs --ignore=ReadApiLogs -w logs.json
```

Once logs are recording to `logs.json`, you can separatly work on them.

You can live-view them using `tail -f logs.json` or use more advanced tools like [jq](https://stedolan.github.io/jq/) to query json documents.

Example of showing only calls date, name and status code with tab-separated values:
```
jq -s -r '(.[] | [.QueryDate, .QueryCallName, .ResponseStatusCode]) | @tsv' logs.json
2022-06-17T12:14:28.378111Z	ReadVolumes	200
2022-06-17T12:14:28.378111Z	ReadVolumes	200
2022-06-17T12:14:30.379899Z	CreateVms	200
```

# License

> Copyright Outscale SAS
> BSD-3-Clause

Check [CONTRIBUTING.md](CONTRIBUTING.md) for more details about license testing.
