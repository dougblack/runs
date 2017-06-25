package cli

import (
	"context"
	"flag"
	"fmt"
	"github.com/dougblack/runs/data"
	"github.com/fatih/color"
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
		if date.Year() == runDate.Year() && date.YearDay() == runDate.YearDay() {
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
	headerHighlight := color.New(color.FgBlue).Add(color.Bold).SprintFunc()
	fmt.Printf("%s%s\n", padding, headerHighlight(headerMonth))
	fmt.Println("Su Mo Tu We Th Fr Sa Tt")
}

func highlight(now time.Time, date time.Time) (style func(...interface{}) string) {
	normal := color.New(color.FgWhite).Add(color.BgBlack).Add(color.Bold).SprintFunc()
	marked := color.New(color.FgBlack).Add(color.BgGreen).Add(color.Bold).SprintFunc()
	if date.Year() == now.Year() && date.YearDay() == now.YearDay() {
		style = marked
	} else {
		style = normal
	}
	return style
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
		style := highlight(now, date)
		if dayTotal > 0 {
			fmt.Printf(style("%2d"), int(dayTotal))
		} else {
			fmt.Printf(style("--"))
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
