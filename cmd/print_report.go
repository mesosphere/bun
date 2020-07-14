package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mitchellh/go-wordwrap"

	"github.com/logrusorgru/aurora"
	"github.com/olekukonko/tablewriter"

	"github.com/mesosphere/bun/v2/checks"
)

// MsgErr is a standard message used in the check summary when errors
// occurs during the check.

func printReport(c checks.Check, r checks.Results, verbose bool) {
	au := aurora.NewAurora(useColors())
	var status string
	var summary string
	switch r.Status() {
	case checks.SOK:
		if !verbose {
			return
		}
		status = au.Bold(au.Green("[" + r.Status() + "]")).String()
		summary = c.OKSummary
	case checks.SProblem:
		status = au.Bold(au.Red("[" + r.Status() + "]")).String()
		summary = c.ProblemSummary
		if len(r.Undefined()) > 0 {
			summary += "\n" + "Couldn't check all hosts. See details below."
		}
	case checks.SUndefined:
		if !verbose {
			return
		}
		status = au.Bold(au.Yellow("[" + r.Status() + "]")).String()
		if len(r.OKs()) == 0 {
			summary = "Couldn't check any hosts because of the error(s)."
		} else {
			summary = "Couldn't check some hosts because of the error(s)." +
				" Please find the details below."
		}
	default:
		panic("Unknown status: " + r.Status())
	}
	var data tableData = make([][]string, 0, len(r)+5)
	data.appendBulk([][]string{
		{au.Bold("Check").String(), c.Name},
		{au.Bold("Status").String(), status},
		{au.Bold("Description").String(), c.Description},
	})
	if r.Status() == checks.SProblem {
		data.append([]string{au.Bold("Cure").String(), c.Cure})
	}
	data.append([]string{au.Bold("Summary").String(), summary})
	data.appendBulk(resultsData(r.Problems()))
	data.appendBulk(resultsData(r.Undefined()))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetAutoWrapText(false)
	table.AppendBulk(data)
	table.Render()
	fmt.Println()
}

func useColors() bool {
	return !outputRedirectedToFile() && !noColor
}

func resultsData(results checks.Results) [][]string {
	au := aurora.NewAurora(useColors())
	data := make([][]string, 0, len(results))
	for _, result := range results {
		var leftColumn string
		switch result.Status {
		case checks.SProblem:
			leftColumn = au.Bold(au.Red("[P]")).String()
		case checks.SOK:
			leftColumn = au.Bold(au.Green("[OK]")).String()
		case checks.SUndefined:
			leftColumn = au.Bold(au.Yellow("[U]")).String()
		default:
			panic("Unknown status: " + result.Status)
		}
		if result.Host.IP != "" {
			leftColumn += fmt.Sprintf(" %v %v", result.Host.Type, result.Host.IP)
		}
		data = append(data, []string{leftColumn, fmt.Sprintf("%v", result.Value)})
	}
	return data
}

type tableData [][]string

func (d *tableData) append(row []string) {
	*d = append(*d, wrap(row))
}

func (d *tableData) appendBulk(rows [][]string) {
	for _, row := range rows {
		d.append(row)
	}
}
func wrap(row []string) []string {
	return []string{
		wordwrap.WrapString(row[0], 50),
		wordwrap.WrapString(row[1], 60),
	}
}

func printSummary(results []checks.Results) {
	au := aurora.NewAurora(useColors())
	total := len(results)
	nP := 0
	nOK := 0
	nU := 0
	for _, r := range results {
		switch r.Status() {
		case checks.SProblem:
			nP++
		case checks.SOK:
			nOK++
		case checks.SUndefined:
			nU++
		default:
			panic("Unknown status " + r.Status())
		}
	}
	var data tableData = make([][]string, 0, 4)
	data.appendBulk([][]string{
		{au.Bold("Problem").String(), strconv.Itoa(nP)},
		{au.Bold("Undefined").String(), strconv.Itoa(nU)},
		{au.Bold("OK").String(), strconv.Itoa(nOK)},
	})
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Summary", ""})
	table.SetFooter([]string{"Total", strconv.Itoa(total)})
	table.AppendBulk(data)
	table.Render()
}
