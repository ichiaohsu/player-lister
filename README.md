# player-lister

A simple asynchronous http request script to request API for players in target teams. After eliminating players showed up in both club and national team, it will print results to stdout.

## Environment

Please use go >= 1.8. Recommended 1.9.2

## Usage

Clone this repository or `main.go` under `$GOPATH/src/`. In the same directory with `main.go` use command
```bash
go run main.go
```

to execute.