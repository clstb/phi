package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func Review(ctx *cli.Context) error {
	fp := ctx.Path("file")
	f, err := os.OpenFile(fp, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)

	var from, to string
	for i, record := range records {
		if i >= 1 && records[i][3] != records[i-1][3] {
			fmt.Fprintf(w, "\t\t\t\t\t\n")
		}

		if record[0] != "" {
			from = record[0]
		}
		if record[1] != "" {
			to = record[1]
		}

		amount := record[5]
		if strings.HasPrefix(amount, "-") {
			amount = color.HiRedString(amount)
		} else {
			amount = color.HiGreenString(amount)
		}

		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\t%s\t%s\n",
			from,
			to,
			record[2],
			record[3],
			record[4],
			amount,
		)
	}
	return w.Flush()
}
