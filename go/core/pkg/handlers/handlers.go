package handlers

/*
var OauthKeeperUri = func() string {
	str := os.Getenv("OAUTHKEEPER_URI")
	if len(str) == 0 {
		panic("OAUTHKEEPER_URI == nill")
	}
	return str
}()

func GetTinkLink(c *gin.Context, client *client.Client) {
	var json pkg3.SessionId
	err := c.BindJSON(&json)
	if err != nil {
		debug.PrintStack()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	token, err := pkg3.UserTokenCache.Get(context.TODO(), json.SessionId)
	if err != nil {
		debug.PrintStack()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if token == nil {
		pkg.Sugar.Error("Not logged in")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "Not logged in"})
		return
	}

	client.SetBearerToken(token.(string))
	link, err := client.GetLink()
	if err != nil {
		pkg.Sugar.Error(err)
		debug.PrintStack()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"link": link})
}

func SyncLedger(c *gin.Context, client *client.Client) {
	var json pkg3.SyncLedgerRequest
	err := c.BindJSON(&json)
	if err != nil {
		pkg.Sugar.Error(err)
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

	ledger := commands2.ParseLedger(fmt.Sprintf("%s/%s", pkg3.DataDirPath, json.Username))
	ledger = commands2.UpdateLedger(ledger, providers, accounts, filteredTransactions)
	c.JSON(http.StatusOK, gin.H{})
}
*/
