package csv

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/clstb/phi/cmd"
	"github.com/clstb/phi/pkg/config"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/urfave/cli/v2"
)

func Parse(ctx *cli.Context) error {
	core, err := cmd.Core(ctx)
	if err != nil {
		return err
	}

	transactionsPB, err := core.GetTransactions(
		ctx.Context,
		&pb.TransactionsQuery{},
	)
	if err != nil {
		return err
	}

	transactions, err := fin.TransactionsFromPB(transactionsPB)
	if err != nil {
		return err
	}

	hashes := make(map[string]struct{})
	for _, transaction := range transactions {
		hashes[transaction.Hash] = struct{}{}
	}

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
		hash := sha256.New()
		_, err = hash.Write([]byte(strings.Join([]string{
			record[fc.Date],
			record[fc.Entity],
			record[fc.Reference],
			record[fc.Amount] + " " + record[fc.Currency],
		}, "")))
		if err != nil {
			return err
		}
		hashStr := hex.EncodeToString(hash.Sum(nil))

		_, ok := hashes[hashStr]
		if ok {
			continue
		}

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
