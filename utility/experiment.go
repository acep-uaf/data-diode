package utility

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func Isolation() {
	data := [][]string{
		{"Objective", "Vary the packet size in a one packet per second test to observe the impact on bandwidth and latency."},
		{"Expectation", "Low data rate (small packets) will deliver poor bandwidth but good latency, high data rate (large packets) will deliver good bandwidth but poor latency."},
		{"Scenario", "Investigate what happens if more data per second is sent than the diode can forward."},
		{"Question", "Is the excess data queued or dropped?"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Experiment", "Description"})

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}

func Validation() {
	fmt.Println(">> Variable Packet Size")
}
