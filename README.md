# osc-logs

It is a tool allowing users of the 3DS Outscale cloud to easily consult and keep the logs of calls made on the IaaS.

# **Features**
From the moment the program starts and one or more logs are available, the program displays it on the standard output (JSON format).
Each log will be displayed on a single line in a compact way.
The program can be stopped with ctrl-c.

```
Description:
    osc-logs

Options:
    -w, --write   Write all traces inside a file instead of writing to standard output
    -c, --count   Exit after <count> logs
```

# **License**
Copyright Outscale SAS

BSD-3-Clause

LICENSE folder contain raw licenses terms following spdx naming.

You can check which license apply to which copyright owner through .reuse/dep5 specification.

You can test reuse compliance by running make test-reuse.
