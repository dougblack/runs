package main

import (
    "context"
    "flag"
    "fmt"
    "os"
    "strconv"

    "github.com/google/subcommands"
)

type addRunCommand struct {
    add bool
}

func (*addRunCommand) Name() string { return "add" }
func (*addRunCommand) Synopsis() string { return "Add a run to the log" }
func (*addRunCommand) Usage() string {
    return `runs add [miles]:
    Add a run to the run log.
`
}

func (r *addRunCommand) SetFlags(f *flag.FlagSet) {
    f.BoolVar(&r.add, "add", false, "add a new run")
}

func (r *addRunCommand) Execute(_ context.Context, f *flag.FlagSet, _ ... interface{}) subcommands.ExitStatus {
    if f.NArg() != 1 {
        fmt.Printf("Received too many arguments for `add`: %s\n", f.Args())
        return subcommands.ExitUsageError
    }
    miles, err := strconv.ParseFloat(f.Arg(0), 32)
    if err != nil {
        fmt.Printf("Failed to convert %s to 32-bit float\n", f.Arg(0))
        return subcommands.ExitUsageError
    }
    fmt.Printf("Adding %.2f\n", miles)
    return subcommands.ExitSuccess
}

func main() {
  subcommands.Register(subcommands.HelpCommand(), "")
  subcommands.Register(subcommands.FlagsCommand(), "")
  subcommands.Register(subcommands.CommandsCommand(), "")
  subcommands.Register(&addRunCommand{}, "")

  flag.Parse()
  ctx := context.Background()
  os.Exit(int(subcommands.Execute(ctx)))
}
