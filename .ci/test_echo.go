package main

import (
	"github.com/vladimirvivien/echo"
)

func main() {
	echo.New().Env("CGO_ENABLED=0 GOOS=linux GOARCH=amd64").Runout("go build -o build/amd64/linux/crash-diagnostics .")
}
