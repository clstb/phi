package main

import (
	"fmt"
	"github.com/clstb/phi/go/internal/phi/commands"
	"github.com/clstb/phi/go/pkg/client"
	"github.com/clstb/phi/go/pkg/client/tink"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func SyncLedger(c *gin.Context, client *client.Client) {
	var json SyncLedgerRequest
	err := c.BindJSON(&json)
	if err != nil {
		sugar.Error(err)
		debug.PrintStack()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	providers, err := client.GetProviders("DE")
	if err != nil {
		debug.PrintStack()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	accounts, err := client.GetAccounts("")
	if err != nil {
		debug.PrintStack()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	transactions, err := client.GetTransactions("")
	if err != nil {
		debug.PrintStack()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var filteredTransactions []tink.Transaction
	for _, transaction := range transactions {
		if transaction.Status != "BOOKED" {
			continue
		}
		filteredTransactions = append(filteredTransactions, transaction)
	}

	ledger := commands.ParseLedger(fmt.Sprintf("%s/%s", DataDirPath, json.Username))
	ledger = commands.UpdateLedger(ledger, providers, accounts, filteredTransactions)
	c.JSON(http.StatusOK, gin.H{})
}
