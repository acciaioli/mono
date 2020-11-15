package display

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func String(s string) {
	fmt.Println(s)
}

func Table(headers []string, data [][]string) {
	if data != nil {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t")
		table.SetNoWhiteSpace(true)

		table.SetHeader(headers)
		table.AppendBulk(data)
		table.Render()
	}
}
