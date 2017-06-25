package cli

import (
	"context"
	"flag"
	"fmt"
	"github.com/dougblack/runs/data"
	"github.com/google/subcommands"
	"strings"
	"time"
)

type MonthCommand struct {
	month bool
}

func (*MonthCommand) Name() string     { return "month" }
func (*MonthCommand) Synopsis() string { return "Show runs of this month" }
func (*MonthCommand) Usage() string {
	return `runs month:
    Show month run.
`
}

func (m *MonthCommand) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&m.month, "month", false, "show month run")
}

func total(now time.Time, date time.Time, runs []data.Run) float64 {
	total := 0.0
	for _, run := range runs {
		runDate := run.Date.In(now.Location())
		sameYear := (date.Year() == runDate.Year())
		sameMonth := (date.Month() == runDate.Month())
		sameDay := (date.Day() == runDate.Day())
		if sameYear && sameMonth && sameDay {
			total += run.Miles
		}
	}
	return total
}

func header(now time.Time) {
	const headerWidth = 26
	headerMonth := fmt.Sprintf("%s %d", now.Month().String(), now.Year())
	leftPadding := (headerWidth - len(headerMonth)) / 2
	padding := strings.Repeat(" ", leftPadding)
	fmt.Printf("%s%s\n", padding, headerMonth)
	fmt.Println("Su Mo Tu We Th Fr Sa Tt")
}

func printRuns(runs []data.Run) {
	now := time.Now().Local()
	header(now)

	totalDays := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Day()
	first := time.Date(now.Year(), now.Month(), 1, 1, 0, 0, 0, now.Location())
	last := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location())

	for i := 0; i < int(first.Weekday()); i++ {
		print("   ")
	}

	weekTotal := 0.0
	for day := 0; day <= totalDays; day++ {
		date := time.Date(now.Year(), now.Month(), day+1, 0, 0, 0, 0, now.Location())
		dayTotal := total(now, date, runs)
		weekTotal = weekTotal + dayTotal
		if dayTotal > 0 {
			fmt.Printf("%2d", int(dayTotal))
		} else {
			fmt.Printf("--")
		}
		fmt.Print(" ")
		if date.Weekday() == time.Saturday {
			fmt.Printf("%3.1f\n", weekTotal)
			weekTotal = 0.0
		}
	}

	for i := int(last.Weekday()); i <= int(time.Sunday); i++ {
		date := time.Date(now.Year(), now.Month()+1, i, 0, 0, 0, 0, now.Location())
		fmt.Printf("--")
		fmt.Print(" ")
		if date.Weekday() == time.Saturday {
			fmt.Printf(" %3.1f\n", weekTotal)
		}
	}
}

func (m *MonthCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 0 {
		fmt.Printf("Received too many arguments for `month`: %s\n", f.Args())
		return subcommands.ExitUsageError
	}
	runs := data.LastMonth()
	printRuns(runs)
	return subcommands.ExitSuccess
}
