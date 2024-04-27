
# Twitch EventSub Handler

This repository aims to represent how I will be handling EventSub notifications from Twitch


## [Goose Documentation](https://github.com/pressly/goose?tab=readme-ov-file#install)

```shell
EXPORT GOOSE_DRIVER="postgres"
EXPORT GOOSE_DBSTRING="postgres://test:password@localhost:5432/test"
EXPORT GOOSE_MIGRATION_DIR="migrations/"
go install github.com/pressly/goose/v3/cmd/goose@latest
goose up
```

```powershell
$env:"GOOSE_DRIVER" = "postgres"
$env:"GOOSE_DBSTRING" = "postgres://test:password@localhost:5432/test"
$env:"GOOSE_MIGRATION_DIR" = "migrations/"
go install github.com/pressly/goose/v3/cmd/goose@latest
goose up
```
