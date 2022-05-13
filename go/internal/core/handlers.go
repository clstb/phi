package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/clstb/phi/go/internal/phi/commands"
	"github.com/clstb/phi/go/pkg"
	"github.com/clstb/phi/go/pkg/client"
	"github.com/clstb/phi/go/pkg/client/tink"
	"github.com/gin-gonic/gin"
	ory "github.com/ory/kratos-client-go"
	"net/http"
	"os"
	"runtime/debug"
)

func GetTinkLink(c *gin.Context, client *client.Client) {
	var json SessionId
	err := c.BindJSON(&json)
	if err != nil {
		debug.PrintStack()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	token, err := UserTokenCache.Get(context.TODO(), json.SessionId)
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

func DoLogin(c *gin.Context, client *client.Client) {
	var json LoginRequest
	err := c.BindJSON(&json)
	if err != nil {
		debug.PrintStack()
		c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
		return
	}
	//logger.Info("Executing login for:", zap.Object("json", &json))

	sess, err := client.Login(json.Username, json.Password)
	if err != nil {
		debug.PrintStack()
		c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
		return
	}
	putClientSessionTokenInCache(sess.Id, sess.Token)
	c.JSON(http.StatusOK, gin.H{"sessionId": sess.Id})
}

func DoRegister(c *gin.Context, client *client.Client) {
	var json LoginRequest
	err := c.BindJSON(&json)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	//logger.Info("Executing register for:", zap.Object("json", &json))
	sess, err := client.Register(json.Username, json.Password)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = provisionTinkUser(sess, client)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = createUserDir(json.Username)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	putClientSessionTokenInCache(sess.Id, sess.Token)
	c.JSON(http.StatusOK, gin.H{"sessionId": sess.Id})
}

/*
Each user needs
.data/username/accounts.bean
.data/username/transactions/
*/

func createUserDir(username string) error {
	err := os.MkdirAll(fmt.Sprintf("%s/%s", DataDirPath, username), os.ModePerm)
	if err != nil {
		pkg.Sugar.Error(err)
		debug.PrintStack()
		return err
	}
	err = os.MkdirAll(fmt.Sprintf("%s/%s/transactions", DataDirPath, username), os.ModePerm)
	if err != nil {
		pkg.Sugar.Error(err)
		debug.PrintStack()
		return err
	}
	_, err = os.Create(fmt.Sprintf("%s/%s/accounts.bean", DataDirPath, username))
	if err != nil {
		pkg.Sugar.Error(err)
		debug.PrintStack()
		return err
	}
	return nil
}

// Each user needs to be registered as client in Tink organisation
func provisionTinkUser(session client.Session, oauthKeeperClient *client.Client) error {
	body := PhiSessionRequest{
		Token:   session.Token,
		Session: session.Session,
	}

	_json, err := json.Marshal(body)
	if err != nil {
		return err
	}

	oauthKeeperClient.SetBearerToken(session.Token)

	resp, err := oauthKeeperClient.SendRequest("POST", OauthKeeperUri+"/tink-user", "application/json", bytes.NewBuffer(_json))
	if err != nil {
		return err
	}

	var res PhiClientIdResponse
	err = json.Unmarshal([]byte(resp), &res)
	if err != nil {
		return err
	}

	traits := session.Identity.Traits.(map[string]interface{})
	traits["tink_id"] = res.TinkId

	oryConf := ory.NewConfiguration()
	oryConf.Servers = ory.ServerConfigurations{{URL: OriUrl}}
	oryConf.AddDefaultHeader("Authorization", "Bearer "+OryToken)
	oryConf.HTTPClient = &http.Client{}
	oryClient := ory.NewAPIClient(oryConf)

	identity, _, err := oryClient.V0alpha2Api.AdminUpdateIdentity(context.Background(), session.Identity.Id).AdminUpdateIdentityBody(
		ory.AdminUpdateIdentityBody{State: *session.Identity.State, Traits: traits}).Execute()
	if err != nil {
		pkg.Sugar.Error(err)
		return err
	}
	session.Identity = *identity
	return nil
}

func SyncLedger(c *gin.Context, client *client.Client) {
	var json SyncLedgerRequest
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

	ledger := commands.ParseLedger(fmt.Sprintf("%s/%s", DataDirPath, json.Username))
	ledger = commands.UpdateLedger(ledger, providers, accounts, filteredTransactions)
	c.JSON(http.StatusOK, gin.H{})
}
