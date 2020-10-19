package main

import (
    _ "fmt"
    "log"
    "os"

    "github.com/dhuan/giback/pkg/cmd"

    "github.com/urfave/cli/v2"
)

func main() {
  app := &cli.App{
    Name: "boom",
    Usage: "make an explosive entrance",
    Action: cmd.Main,
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
