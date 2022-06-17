# osc-logs

osc-logs constantly fetch API call logs from Outscale to easily consult and keep them.

# Features

osc-logs fetch API call logs every few seconds and only print logs from program's start date.

By default logs are printed as Line-delimited JSON to standard output. See [format documentation](https://docs.outscale.com/api#tocslog) for more details.

We recommand using [jq](https://stedolan.github.io/jq/) utility for additional filtering and formating.

```
Description:
    osc-logs

Options:
    -w, --write      Write all traces inside a file instead of writing to standard output
    -c, --count      Exit after <count> logs
    -i, --interval   Wait a duration defined by <wait> (in seconds) between two calls to Outscale API 
    -p, --profile    Use a specific profile name ("default" is the default profile )
    -I, --ignore     Ignore one or more specific API calls. Values are separated by commas e.g. "--ignore=ReadApiLogs,ReadVms"
```

# Installation

Download latest binary in [Release page](https://github.com/outscale/osc-logs/releases) or run:
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

Example of storing in a file all logs except ReadApiLogs itself:
```
osc-logs --ignore=ReadApiLogs -w logs.json
```

Once logs are recording to `logs.json`, you can separatly work on them.

You can live-view them using `tail -f logs.json` or use more advanced tools like [jq](https://stedolan.github.io/jq/).

Example of showing only calls date, name and status code in a tsv format:
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
