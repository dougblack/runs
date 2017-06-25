package cli

import (
    "context"
    "flag"
    "fmt"
    "github.com/google/subcommands"
    "github.com/dougblack/runs/data"
)

type LastCommand struct {
    last bool
}

func (*LastCommand) Name() string { return "last" }
func (*LastCommand) Synopsis() string { return "Show last run" }
func (*LastCommand) Usage() string {
    return `runs last:
    Show last run.
`
}

func (l *LastCommand) SetFlags(f *flag.FlagSet) {
    f.BoolVar(&l.last, "last", false, "show last run")
}

func (l *LastCommand) Execute(_ context.Context, f *flag.FlagSet, _ ... interface{}) subcommands.ExitStatus {
    if f.NArg() != 0 {
        fmt.Printf("Received too many arguments for `last`: %s\n", f.Args())
        return subcommands.ExitUsageError
    }
    miles, date := data.LastRun()
    fmt.Printf("%s: %.2f miles\n", date.Format("1/2"), miles)
    return subcommands.ExitSuccess
}
