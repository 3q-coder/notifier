# Notifier
Test websocket server in Golang.

## Run
```
go run cmd/notifier/main.go
```

## Test
To run tests, go to the appropriate directory and run `go test` command.

To get code coverage, run `go test -cover`.

To get information about coverage, run:

`go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html`