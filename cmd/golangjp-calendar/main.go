package main

import (
	"os"

	calendar "github.com/vzvu3k6k/golangjp-calendar"
)

func main() {
	calendar.Run(os.Stdout, os.Args[1:])
}
