package csv

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/clstb/phi/pkg/config"
	"github.com/clstb/phi/pkg/core/db"
	"github.com/clstb/phi/pkg/fin"
)

type Loader struct {
	r *csv.Reader
	c config.FileConfig
}

func New(
	f *os.File,
	c *config.Config,
	opts ...Opt,
) (*Loader, error) {
	l := &Loader{
		r: csv.NewReader(f),
	}

	for _, fileConfig := range c.Files {
		re, err := regexp.Compile(fileConfig.Regex)
		if err != nil {
			return nil, err
		}
		if re.MatchString(filepath.Base(f.Name())) {
			l.c = fileConfig
		}
	}
	fmt.Println(c)

	for _, opt := range opts {
		opt(l)
	}

	return l, nil
}

func (l *Loader) Load() (
	transaction fin.Transaction,
	amount db.Amount,
	err error,
) {
	f, err := l.r.Read()
	if err != nil {
		return
	}

	amountStr := f[l.c.Amount]
	currency := f[l.c.Currency]
	dateStr := f[l.c.Date]
	entity := f[l.c.Entity]
	reference := f[l.c.Reference]

	fmt.Println(amountStr, currency)
	amount, err = db.AmountFromString(
		fmt.Sprintf("%s %s", amountStr, currency),
		db.AmountEU,
	)
	if err != nil {
		return
	}

	date, err := time.Parse(l.c.DateFormat, dateStr)
	if err != nil {
		return
	}

	hash := sha256.New()
	_, err = hash.Write([]byte(strings.Join(f, "")))
	if err != nil {
		return
	}
	hashStr := hex.EncodeToString(hash.Sum(nil))

	transaction = fin.NewTransaction(db.Transaction{
		Date:      date,
		Entity:    entity,
		Reference: reference,
		Hash:      hashStr,
	}, nil)
	return
}

type Opt func(l *Loader)

func WithSeperator(seperator rune) Opt {
	return func(l *Loader) {
		l.r.Comma = seperator
	}
}
