package main

import (
    "context"
    "flag"
    "os"
    "github.com/google/subcommands"
    "github.com/dougblack/runs/cli"
)

func main() {
  subcommands.Register(subcommands.HelpCommand(), "")
  subcommands.Register(subcommands.FlagsCommand(), "")
  subcommands.Register(subcommands.CommandsCommand(), "")
  subcommands.Register(&cli.AddCommand{}, "")
  subcommands.Register(&cli.LastCommand{}, "")

  flag.Parse()
  ctx := context.Background()
  os.Exit(int(subcommands.Execute(ctx)))
}
