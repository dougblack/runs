package cli

import (
	"context"
	"flag"
	"fmt"
	"github.com/dougblack/runs/data"
	"github.com/fatih/color"
	"github.com/google/subcommands"
	"strconv"
	"time"
)

type AddCommand struct {
	add bool
}

func (*AddCommand) Name() string     { return "add" }
func (*AddCommand) Synopsis() string { return "Add a run to the log" }
func (*AddCommand) Usage() string {
	return `runs add [miles]:
    Add a run to the run log.
`
}

func (a *AddCommand) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&a.add, "add", false, "add a new run")
}

func (a *AddCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() <= 0 {
		fmt.Printf("Received incorrect number of arguments for `add`: %s\n", f.Args())
		return subcommands.ExitUsageError
	}
	miles, err := strconv.ParseFloat(f.Arg(0), 32)
	if err != nil {
		fmt.Printf("Failed to convert %s to 32-bit float\n", f.Arg(0))
		return subcommands.ExitUsageError
	}
	date := time.Now().UTC()
	if f.NArg() == 2 {
		date, err = time.Parse("1/2", f.Arg(1))
		date = time.Date(time.Now().UTC().Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		if err != nil {
			panic(err)
		}
	}
	data.AddRun(miles, date)
	dateColor := color.New(color.FgBlue).Add(color.Bold).SprintFunc()
	fmt.Printf("%s: %.2f miles\n", dateColor(date.Format("1/2")), miles)
	return subcommands.ExitSuccess
}
