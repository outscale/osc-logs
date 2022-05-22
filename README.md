# osc-logs

## **Introduction**
It is a tool allowing users of the 3DS Outscale cloud to easily consult and keep the logs of calls made on the IaaS.

## **How to use**
From the moment the program starts and one or more logs are available, the program displays it on the standard output (JSON format).
Each log will be displayed on a single line in a compact way.
The program can be stopped with ctrl-c.

### Options:
-w: Write all logs to a JSON file instead of writing to standard output</br>
-i: Wait a defined duration, an interval of 2 seconds between two calls to Outscale API before displaying new logs.</br>
-c: Define a number of logs to display beyond which the program is closed.</br>
-p: View logs by specific profile name.</br>
The user must specify the name of his profile otherwise if he does not specify this option, the default one is used --p profile-name or --profile profile-name</br>
-h: Displays help presenting a description of the program, its version and the options available with their description
