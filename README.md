# osc-logs

It is a tool allowing users of the 3DS Outscale cloud to easily consult and keep the logs of calls made on the IaaS.

# Features
From the moment the program starts and one or more logs are available, the program displays it on the standard output (JSON format).
Each log will be displayed on a single line in a compact way.
The program can be stopped with ctrl-c.

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

# License

> Copyright Outscale SAS
> BSD-3-Clause

Check [CONTRIBUTING.md](CONTRIBUTING.md) for more details about license testing.