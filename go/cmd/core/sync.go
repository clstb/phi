package main

import (
	"context"
	"fmt"
	"github.com/clstb/phi/go/internal/phi/commands"
	"github.com/clstb/phi/go/pkg/client"
	"github.com/clstb/phi/go/pkg/client/tink"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func SyncLedger(c *gin.Context) {
	var json SyncLedgerRequest
	err := c.BindJSON(&json)
	if err != nil {
		sugar.Error(err)
		debug.PrintStack()
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
		return
	}
	userClient, err := UserClientCache.Get(context.TODO(), json.SessionId)
	if err != nil {
		debug.PrintStack()
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	if userClient == nil {
		sugar.Error("Not logged in")
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Not logged in"})
		return
	}

	providers, err := userClient.(*client.Client).GetProviders("DE")
	if err != nil {
		debug.PrintStack()
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	accounts, err := userClient.(*client.Client).GetAccounts("")
	if err != nil {
		debug.PrintStack()
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	transactions, err := userClient.(*client.Client).GetTransactions("")
	if err != nil {
		debug.PrintStack()
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
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
