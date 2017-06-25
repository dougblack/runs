package cli

import (
    "context"
    "flag"
    "fmt"
    "github.com/google/subcommands"
    "github.com/dougblack/runs/data"
)

type MonthCommand struct {
    month bool
}

func (*MonthCommand) Name() string { return "month" }
func (*MonthCommand) Synopsis() string { return "Show runs of this month" }
func (*MonthCommand) Usage() string {
    return `runs month:
    Show month run.
`
}

func (m *MonthCommand) SetFlags(f *flag.FlagSet) {
    f.BoolVar(&m.month, "month", false, "show month run")
}

func (m *MonthCommand) Execute(_ context.Context, f *flag.FlagSet, _ ... interface{}) subcommands.ExitStatus {
    if f.NArg() != 0 {
        fmt.Printf("Received too many arguments for `month`: %s\n", f.Args())
        return subcommands.ExitUsageError
    }
    runs := data.LastMonth()
    for _, run := range runs {
        fmt.Printf("%s: %.2f miles\n", run.Date.Format("1/2"), run.Miles)
    }
    return subcommands.ExitSuccess
}
