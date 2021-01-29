package csv

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"os"
	"sort"
	"strings"

	"github.com/clstb/phi/cmd"
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

	w := csv.NewWriter(output)
	w.Comma = ';'
	for i, record := range records {
		hash := sha256.New()
		_, err = hash.Write([]byte(strings.Join([]string{
			record[0],
			record[2],
			record[4],
			record[7] + " " + record[8],
		}, "")))
		if err != nil {
			return err
		}
		hashStr := hex.EncodeToString(hash.Sum(nil))

		_, ok := hashes[hashStr]
		if ok {
			continue
		}

		if i >= 1 && records[i][2] != records[i-1][2] {
			output.Write([]byte("\n"))
		}

		w.Write(append([]string{"", ""}, []string{record[0], record[2], record[4], record[7] + " " + record[8]}...))
		w.Flush()
	}

	return nil
}
