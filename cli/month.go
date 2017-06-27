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
	headerHighlight := color.New(color.FgGreen).Add(color.Bold).SprintFunc()
	fmt.Printf("%s%s\n", padding, headerHighlight(headerMonth))
	fmt.Println("Su Mo Tu We Th Fr Sa Tt")
}

func highlight(now time.Time, date time.Time) (style func(...interface{}) string) {
	normal := color.New(color.FgWhite).Add(color.BgBlack).Add(color.Bold).SprintFunc()
	marked := color.New(color.FgBlack).Add(color.BgBlue).Add(color.Bold).SprintFunc()
	if date.Year() == now.Year() && date.YearDay() == now.YearDay() {
		style = marked
	} else {
		style = normal
	}
	return style
}

func date(month time.Month, day int) time.Time {
	now := time.Now().Local()
	return time.Date(now.Year(), month, day, 0, 0, 0, 0, now.Location())
}

func printRuns(runs []data.Run) {
	now := time.Now().Local()
	header(now)

	first := date(now.Month(), 1)
	last := date(now.Month()+1, 0)
	totalDays := last.Day() + 1

	for i := 0; i < int(first.Weekday()); i++ {
		print("   ")
	}

	weekTotal := 0.0
	for day := 1; day < totalDays; day++ {
		date := date(now.Month(), day)
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

	remainder := time.Saturday - last.Weekday()
	for i := 1; i <= int(remainder); i++ {
		date := date(now.Month(), last.Day()+i)
		fmt.Printf("   ")
		if date.Weekday() == time.Saturday {
			fmt.Printf("%3.1f\n", weekTotal)
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
