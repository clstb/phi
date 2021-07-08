package tinkgw

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/clstb/phi/pkg/db/tinkgw"
	"github.com/clstb/phi/pkg/fin"
	"github.com/clstb/phi/pkg/pb"
	"github.com/clstb/phi/pkg/tink"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc/metadata"
)

func (s *Server) Sync(
	ctx context.Context,
	req *pb.SyncReq,
) (*pb.SyncRes, error) {
	subStr, ok := ctx.Value("sub").(string)
	if !ok {
		return nil, fmt.Errorf("context: missing subject")
	}
	sub, err := uuid.FromString(subStr)
	if err != nil {
		return nil, err
	}

	jwt, ok := ctx.Value("jwt").(string)
	if !ok {
		return nil, fmt.Errorf("context: missing jwt")
	}

	token, err := s.getToken(ctx, sub)

	res, err := s.tink.Transactions(token.AccessToken, "")
	if err != nil {
		return nil, fmt.Errorf("failed getting transactions: %w", err)
	}

	data := res.Transactions
	for res.NextPageToken != "" {
		res, err = s.tink.Transactions(
			token.AccessToken,
			res.NextPageToken,
		)
		if err != nil {
			return nil, fmt.Errorf("failed getting transactions: %w", err)
		}
		data = append(data, res.Transactions...)
	}

	accountsPB, err := s.core.GetAccounts(
		metadata.AppendToOutgoingContext(
			ctx,
			"authorization",
			fmt.Sprintf("Bearer %s", jwt),
		),
		&pb.AccountsQuery{
			Name: "^Uncategorized$",
		},
	)
	accounts, err := fin.AccountsFromPB(accountsPB)
	if err != nil {
		return nil, err
	}
	account := accounts.ByName("Uncategorized")

	var transactions fin.Transactions
	for _, transaction := range data {
		date, err := time.Parse("2006-01-02", transaction.Dates.Booked)
		if err != nil {
			return nil, err
		}

		amount := fin.NewAmount(
			transaction.Amount.Value.UnscaledValue,
			transaction.Amount.Value.Scale*-1,
			transaction.Amount.CurrencyCode,
		)

		t := fin.Transaction{}
		t.Date = date
		t.Entity = transaction.Descriptions.Original
		if transaction.Descriptions.DisplayDescription != "" {
			t.Reference = sql.NullString{
				String: transaction.Descriptions.DisplayDescription,
				Valid:  true,
			}
		}
		t.User = sub
		t.From = account.ID
		t.To = account.ID
		t.Units = amount.Abs()
		t.Debit = amount.IsNegative()
		t.TinkID = sql.NullString{
			String: transaction.ID,
			Valid:  true,
		}
		transactions = append(transactions, t)
	}

	_, err = s.core.CreateTransactions(
		metadata.AppendToOutgoingContext(
			ctx,
			"authorization",
			fmt.Sprintf("Bearer %s", jwt),
		),
		transactions.PB(),
	)
	if err != nil {
		return nil, err
	}

	return &pb.SyncRes{}, nil
}

func (s *Server) getToken(
	ctx context.Context,
	user uuid.UUID,
) (db.Token, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return db.Token{}, err
	}
	q := db.New(tx)

	tokens, err := q.GetTokens(ctx, user)
	if err != nil {
		return db.Token{}, fmt.Errorf(
			"err: %w; rollback err: %v",
			err,
			tx.Rollback(ctx),
		)
	}
	if len(tokens) != 0 {
		return tokens[0], nil
	}

	code, err := s.tink.AuthorizeGrant(&tink.AuthorizeGrantReq{
		ExternalUserID: user.String(),
		Scope:          "accounts:read,balances:read,transactions:read,provider-consents:read,credentials:refresh",
	})
	if err != nil {
		return db.Token{}, fmt.Errorf(
			"err: %w; rollback err: %v",
			err,
			tx.Rollback(ctx),
		)
	}

	tokenRes, err := s.tink.OAuthToken(&tink.OAuthTokenReq{
		Code:         code,
		ClientID:     s.clientID,
		ClientSecret: s.clientSecret,
		GrantType:    "authorization_code",
	})
	if err != nil {
		return db.Token{}, fmt.Errorf(
			"err: %w; rollback err: %v",
			err,
			tx.Rollback(ctx),
		)
	}

	expiresAt := time.Now().Add(time.Duration(tokenRes.ExpiresIn) * time.Second)
	token, err := q.CreateToken(ctx, db.CreateTokenParams{
		AccessToken:  tokenRes.AccessToken,
		RefreshToken: tokenRes.RefreshToken,
		TokenType:    tokenRes.TokenType,
		ExpiresAt:    expiresAt,
		Scope:        tokenRes.Scope,
		User:         user,
	})
	if err != nil {
		return db.Token{}, fmt.Errorf(
			"err: %w; rollback err: %v",
			err,
			tx.Rollback(ctx),
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return db.Token{}, fmt.Errorf("commit err: %w", err)
	}

	return token, nil
}
