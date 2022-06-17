# osc-logs

osc-logs download API call logs from Outscale easily consult and keep them.

# Features

By default logs are printed as Line-delimited JSON to standard output.

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

Download latest binary in Release page or run:
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

# License

> Copyright Outscale SAS
> BSD-3-Clause

Check [CONTRIBUTING.md](CONTRIBUTING.md) for more details about license testing.
