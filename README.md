# domo-cli

go-cli is a Golang CLI for accessing the Domoticz HTTP API. It is based on the [go-domoto](https://github.com/pawal/go-domoto) library implementing the HTTP API library for Domoticz.

Usage is simple:

```sh
$> domo-cli --url https://domoticz.tld -u username -p password list-devices
1: Kub-Belysning
9: Hall Movement
14: Squeezebox Radio

$> domo-cli --url https://domoticz.tld -u username -p password device-toggle 1
```

Full usage:

```
Usage: ./domo-cli OPTIONS <command>

OPTIONS

-h, --help                 print help and exit
-p, --password <password>  Domoticz User password
-u, --user <user>          Domoticz API username
--url <url>                Domoticz API URL (default: http://localhost:8080)
-v, --verbose              log more information

COMMANDS

group-off                  switch group off <id>
list-scenes                list all scenes and groups
device                     info on device <id>
device-on                  switch device on <id>
scene-run                  execute scene <id>
group-on                   switch group on <id>
list-devices               list all devices
scene-info                 list devices in scene/group
device-toggle              toggle device <id>
device-off                 switch device off <id>
```

**Installation:**
When you have a proper Go installation setup, run this:

```sh
$> go install github.com/pawal/domo-cli
```
