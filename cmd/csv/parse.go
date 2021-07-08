package csv

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"sort"

	"github.com/clstb/phi/pkg/config"
	"github.com/urfave/cli/v2"
)

func Parse(ctx *cli.Context) error {
	filePath := ctx.Path("file")
	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
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

	sort.Slice(records, func(i, j int) bool {
		if records[i][2] != records[j][2] {
			return records[i][2] < records[j][2]
		}

		return records[i][4] < records[j][4]
	})

	outputPath := ctx.Path("output")
	output, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer output.Close()

	configPath := ctx.Path("config")
	config, err := config.Load(configPath)
	if err != nil {
		return err
	}

	fc, err := config.ForFile(filepath.Base(f.Name()))
	if err != nil {
		return err
	}

	w := csv.NewWriter(output)
	w.Comma = ';'
	for i, record := range records {
		if i >= 1 && records[i][fc.Entity] != records[i-1][fc.Entity] {
			output.Write([]byte("\n"))
		}

		w.Write(append(
			[]string{"", ""},
			[]string{
				record[fc.Date],
				record[fc.Entity],
				record[fc.Reference],
				record[fc.Amount] + " " + record[fc.Currency],
			}...),
		)
		w.Flush()
	}

	return nil
}
