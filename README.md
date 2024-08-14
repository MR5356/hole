# Hole
An ssh server honeypot

## Usage
### Get Help
```shell
[root@hecs-210000 hole]# ./hole -h
Usage:
  hole [flags]
  hole [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  read        read records from database

Flags:
  -h, --help       help for hole
  -p, --port int   ssh server port (default 2222)
  -v, --version    version for hole

Use "hole [command] --help" for more information about a command.
```
### As a ssh server
```shell
hole -p 22
```
### Read ssh login records from database
```shell
hole read
```
### Run as systemd service
```shell
echo '[Unit]
Description=Hole
After=network.target syslog.target

[Service]
Type=simple
ExecStart=/usr/local/bin/hole -p 2222

[Install]
WantedBy=multi-user.target
' > /etc/systemd/system/hole.service
systemctl daemon-reload
systemctl enable hole --now
```