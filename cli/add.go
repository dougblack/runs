package cli

import (
    "context"
    "flag"
    "fmt"
    "github.com/google/subcommands"
    "github.com/dougblack/runs/data"
    "strconv"
    "time"
)

type AddCommand struct {
    add bool
}

func (*AddCommand) Name() string { return "add" }
func (*AddCommand) Synopsis() string { return "Add a run to the log" }
func (*AddCommand) Usage() string {
    return `runs add [miles]:
    Add a run to the run log.
`
}

func (r *AddCommand) SetFlags(f *flag.FlagSet) {
    f.BoolVar(&r.add, "add", false, "add a new run")
}

func (r *AddCommand) Execute(_ context.Context, f *flag.FlagSet, _ ... interface{}) subcommands.ExitStatus {
    if f.NArg() != 1 {
        fmt.Printf("Received too many arguments for `add`: %s\n", f.Args())
        return subcommands.ExitUsageError
    }
    miles, err := strconv.ParseFloat(f.Arg(0), 32)
    if err != nil {
        fmt.Printf("Failed to convert %s to 32-bit float\n", f.Arg(0))
        return subcommands.ExitUsageError
    }
    now := time.Now().UTC()
    data.AddRun(miles, now)
    fmt.Printf("Added %.2f miles on %s\n", miles, now.Format("1/2"))
    return subcommands.ExitSuccess
}

