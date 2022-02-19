package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	authpb "github.com/clstb/phi/go/pkg/auth/pb"
	"github.com/clstb/phi/go/pkg/config"
	"github.com/clstb/phi/go/pkg/interceptor"
	"github.com/clstb/phi/go/pkg/ledger"
	"github.com/clstb/phi/go/pkg/tink"
	tinkgwpb "github.com/clstb/phi/go/pkg/tinkgw/pb"
	"github.com/shopspring/decimal"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type model struct {
	ctx       *cli.Context
	state     state
	subModels map[state]tea.Model
	help      help.Model

	config config.Config
	ledger ledger.Ledger

	tinkgwClient tinkgwpb.TinkGWClient
	tinkClient   *tink.Client
}

func newModel(ctx *cli.Context) model {
	return model{
		ctx:   ctx,
		state: AUTH,
		subModels: map[state]tea.Model{
			AUTH:        newAuthModel(ctx.Context),
			CLASSIFY:    newClassifyModel(),
			ADD_ACCOUNT: newAddAccountModel(),
			BANKS:       newBanksModel(),
		},
		help: help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		createAuthClient(m.ctx),
		loadConfig(m.ctx.Path("config")),
		loadLedger(m.ctx.Path("ledger")),
		ticker(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case state:
		m.state = msg
		cmds = append(cmds, m.subModels[msg].Init())
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "?":
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+s":
			cmds = append(cmds, sync(m.ledger, m.tinkClient))
		case "ctrl+e":
			cmds = append(cmds, func() tea.Msg { return CLASSIFY })
		case "ctrl+a":
			cmds = append(cmds, func() tea.Msg { return ADD_ACCOUNT })
		case "ctrl+l":
			cmds = append(cmds, openTinkLink(m.ctx.Context, m.tinkgwClient))
		case "ctrl+b":
			cmds = append(cmds, func() tea.Msg { return BANKS })
		}
		switch m.state {
		case DEFAULT:
		default:
			var cmd tea.Cmd
			m.subModels[m.state], cmd = m.subModels[m.state].Update(msg)
			cmds = append(cmds, cmd)
		}
	case tick:
		cmds = append(
			cmds,
			ticker(),
			checkPhiToken(m.config.PhiToken),
			checkTinkToken(
				m.ctx.Context,
				m.config.TinkToken,
				m.tinkgwClient,
			),
		)
	case error:
		panic(msg)
	default:
		for state, subModel := range m.subModels {
			var cmd tea.Cmd
			m.subModels[state], cmd = subModel.Update(msg)
			cmds = append(cmds, cmd)
		}
		switch msg := msg.(type) {
		case ledger.Ledger:
			m.ledger = msg
			cmds = append(cmds, saveLedger(m.ctx.Path("ledger"), m.ledger))
		case tinkgwpb.TinkGWClient:
			m.tinkgwClient = msg
		case *tink.Client:
			m.tinkClient = msg
		case config.Config:
			m.config = msg
			cmds = append(
				cmds,
				createTinkGWClient(m.ctx, m.config.PhiToken),
				createTinkClient(m.ctx, m.config.TinkToken),
			)
		case config.PhiToken:
			m.config.PhiToken = msg
			cmds = append(
				cmds,
				saveConfig(m.ctx.Path("config"), m.config),
				createTinkGWClient(m.ctx, msg),
			)
		case config.TinkToken:
			m.config.TinkToken = msg
			cmds = append(
				cmds,
				saveConfig(m.ctx.Path("config"), m.config),
				createTinkClient(m.ctx, msg),
			)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.state == DEFAULT {
		return ""
	}

	modelView := m.subModels[m.state].View()
	helpView := m.help.View(keys)

	return modelView + "\n\n" + helpView
}

func loadLedger(path string) tea.Cmd {
	return func() tea.Msg {
		var l ledger.Ledger
		filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			f, err := os.OpenFile(path, os.O_RDONLY, info.Mode())
			if err != nil {
				return err
			}

			l = append(l, ledger.Parse(f)...)
			return f.Close()
		})

		return l
	}
}

func saveLedger(
	path string,
	l ledger.Ledger,
) tea.Cmd {
	return func() tea.Msg {
		transactionsByMonth := l.Transactions().ByMonth()
		for month, transactions := range transactionsByMonth {
			f, err := os.OpenFile(
				fmt.Sprintf("%s/transactions/%s.bean", path, month),
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				os.ModePerm,
			)
			if err != nil {
				return err
			}
			defer f.Close()

			sort.Slice(transactions, func(i, j int) bool {
				return transactions[i].Date.Before(transactions[j].Date)
			})

			for _, transaction := range transactions {
				fmt.Fprint(f, transaction.String())
			}
		}

		f, err := os.OpenFile(
			fmt.Sprintf("%s/accounts.bean", path),
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			os.ModePerm,
		)
		if err != nil {
			return err
		}
		defer f.Close()

		opens := l.Opens()
		sort.Slice(opens, func(i, j int) bool {
			return opens[i].Account < opens[j].Account
		})
		for _, open := range opens {
			fmt.Fprint(f, open.String())
		}

		return nil
	}
}

func checkPhiToken(token config.PhiToken) tea.Cmd {
	return func() tea.Msg {
		expiresAt := time.Unix(token.ExpiresAt, 0)
		if time.Now().Before(expiresAt) {
			return nil
		}
		return AUTH
	}
}

func checkTinkToken(
	ctx context.Context,
	token config.TinkToken,
	client tinkgwpb.TinkGWClient,
) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return nil
		}

		now := time.Now()
		expiresAt := time.Unix(token.ExpiresAt, 0)

		var refreshToken string
		if now.Before(expiresAt) {
			return nil
		} else {
			if now.Before(expiresAt.Add(24 * time.Hour)) {
				refreshToken = token.RefreshToken
			}
		}

		token, err := client.GetToken(ctx, &tinkgwpb.Token{
			RefreshToken: refreshToken,
		})
		if err != nil {
			return err
		}

		return config.TinkToken{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
			ExpiresAt:    token.ExpiresAt,
			Scope:        token.Scope,
		}
	}
}

func fetchTinkProviders(tinkClient *tink.Client) (map[string]tink.Provider, error) {
	providers := map[string]tink.Provider{}

	res, err := tinkClient.Providers("DE")
	if err != nil {
		return nil, err
	}
	for _, provider := range res.Providers {
		providers[provider.FinancialInstitutionID] = provider
	}

	return providers, nil
}

func fetchTinkAccounts(client *tink.Client) (tink.Accounts, error) {
	fetch := func(client *tink.Client) (tink.Accounts, error) {
		accounts, err := client.Accounts("")
		if err != nil {
			return nil, fmt.Errorf("failed getting accounts: %w", err)
		}
		for accounts.NextPageToken != "" {
			res, err := client.Accounts(accounts.NextPageToken)
			if err != nil {
				return nil, fmt.Errorf("failed getting accounts: %w", err)
			}
			accounts.Accounts = append(accounts.Accounts, res.Accounts...)
			accounts.NextPageToken = res.NextPageToken
		}

		return accounts.Accounts, nil
	}
	return fetch(client)
}

func fetchTinkTransactions(client *tink.Client) (tink.Transactions, error) {
	fetch := func(client *tink.Client) (tink.Transactions, error) {
		transactions, err := client.Transactions("")
		if err != nil {
			return nil, fmt.Errorf("failed getting transactions: %w", err)
		}
		for transactions.NextPageToken != "" {
			res, err := client.Transactions(transactions.NextPageToken)
			if err != nil {
				return nil, fmt.Errorf("failed getting transactions: %w", err)
			}
			transactions.Transactions = append(transactions.Transactions, res.Transactions...)
			transactions.NextPageToken = res.NextPageToken
		}

		return transactions.Transactions, nil
	}
	return fetch(client)
}

func createAuthClient(ctx *cli.Context) tea.Cmd {
	return func() tea.Msg {
		conn, err := grpc.Dial(
			ctx.String("auth-host"),
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		)
		if err != nil {
			return err
		}
		return authpb.NewAuthClient(conn)
	}
}

func createTinkGWClient(ctx *cli.Context, token config.PhiToken) tea.Cmd {
	return func() tea.Msg {
		if token.AccessToken == "" {
			return nil
		}

		expiresAt := time.Unix(token.ExpiresAt, 0)
		if time.Now().After(expiresAt) {
			return nil
		}

		conn, err := grpc.Dial(
			ctx.String("tinkgw-host"),
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
			grpc.WithUnaryInterceptor(interceptor.ClientAuthUnary(token.AccessToken)),
			grpc.WithStreamInterceptor(interceptor.ClientAuthStream(token.AccessToken)),
		)
		if err != nil {
			return err
		}
		return tinkgwpb.NewTinkGWClient(conn)
	}
}

type AuthorizationRoundTripper struct {
	Token string
	Next  http.RoundTripper
}

func (rt AuthorizationRoundTripper) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+rt.Token)
	return rt.Next.RoundTrip(req)
}

func createTinkClient(ctx *cli.Context, token config.TinkToken) tea.Cmd {
	return func() tea.Msg {
		return tink.NewClient(&http.Client{
			Transport: &AuthorizationRoundTripper{
				Token: token.AccessToken,
				Next:  http.DefaultTransport,
			},
		})
	}
}

func openTinkLink(
	ctx context.Context,
	client tinkgwpb.TinkGWClient,
) tea.Cmd {
	return func() tea.Msg {
		res, err := client.GetLink(ctx, &tinkgwpb.GetLinkReq{
			Market: "DE",
			Locale: "de_DE",
		})
		if err != nil {
			return err
		}

		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", res.Link).Start()
		case "windows":
			err = exec.Command("rundll32", "url.dll,FileProtocolHandler", res.Link).Start()
		case "darwin":
			err = exec.Command("open", res.Link).Start()
		default:
			err = fmt.Errorf("unsupported platform")
		}
		if err != nil {
			return err
		}

		return nil
	}
}

func sync(
	l ledger.Ledger,
	client *tink.Client,
) tea.Cmd {
	return func() tea.Msg {
		providers, err := fetchTinkProviders(client)
		if err != nil {
			return err
		}

		accounts, err := fetchTinkAccounts(client)
		if err != nil {
			return err
		}

		transactions, err := fetchTinkTransactions(client)
		if err != nil {
			return err
		}

		var filteredTransactions tink.Transactions
		for _, transaction := range transactions {
			if transaction.Status != "BOOKED" {
				continue
			}
			filteredTransactions = append(filteredTransactions, transaction)
		}

		return updateLedger(l, providers, accounts, filteredTransactions)
	}
}

func updateLedger(
	l ledger.Ledger,
	providers map[string]tink.Provider,
	accounts tink.Accounts,
	transactions tink.Transactions,
) ledger.Ledger {
	opensByTinkId := l.Opens().ByTinkId()
	for _, account := range accounts {
		_, ok := opensByTinkId[account.ID]
		if ok {
			continue
		}

		l = append(l, ledger.Open{
			Date: "1970-01-01",
			Account: fmt.Sprintf(
				"Assets:%s:%s",
				providers[account.FinancialInstitutionID].DisplayName,
				account.Name,
			),
			Metadata: []ledger.Metadata{
				{
					Key:   "tink_id",
					Value: strconv.Quote(account.ID),
				},
			},
		})
	}
	opensByTinkId = l.Opens().ByTinkId()

	transactionsByTinkId := l.Transactions().ByTinkId()
	for _, transaction := range transactions {
		_, ok := transactionsByTinkId[transaction.ID]
		if ok {
			continue
		}

		amount := ledger.Amount{
			Decimal: decimal.New(
				transaction.Amount.Value.UnscaledValue,
				transaction.Amount.Value.Scale*-1,
			),
			Currency: transaction.Amount.CurrencyCode,
		}
		var balanceAccount string
		if amount.IsNegative() {
			balanceAccount = "Expenses:Unassigned"
		} else {
			balanceAccount = "Income:Unassigned"
		}

		date, _ := time.Parse("2006-01-02", transaction.Dates.Booked)

		l = append(l, ledger.Transaction{
			Date:      date,
			Type:      "*",
			Payee:     transaction.Reference,
			Narration: transaction.Descriptions.Display,
			Metadata: []ledger.Metadata{
				{
					Key:   "tink_id",
					Value: strconv.Quote(transaction.ID),
				},
			},
			Postings: []ledger.Posting{
				{
					Account: balanceAccount,
					Units:   amount.Neg(),
				},
				{
					Account: opensByTinkId[transaction.AccountID].Account,
					Units:   amount,
				},
			},
		})
	}

	return l
}
