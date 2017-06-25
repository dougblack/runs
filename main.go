package main

import (
	"context"
	"flag"
	"github.com/dougblack/runs/cli"
	"github.com/google/subcommands"
	"os"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&cli.AddCommand{}, "")
	subcommands.Register(&cli.LastCommand{}, "")
	subcommands.Register(&cli.MonthCommand{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
