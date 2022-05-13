package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/clstb/phi/go/pkg/client"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"runtime/debug"
)

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
	err = provisionTinkUser(sess)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	err = createUserDir(json.Username)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
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
		Sugar.Error(err)
		debug.PrintStack()
		return err
	}
	err = os.MkdirAll(fmt.Sprintf("%s/%s/transactions", DataDirPath, username), os.ModePerm)
	if err != nil {
		Sugar.Error(err)
		debug.PrintStack()
		return err
	}
	_, err = os.Create(fmt.Sprintf("%s/%s/accounts.bean", DataDirPath, username))
	if err != nil {
		Sugar.Error(err)
		debug.PrintStack()
		return err
	}
	return nil
}

// Each user needs to be registered as client in Tink organisation
func provisionTinkUser(session client.Session) error {
	body := PhiSessionRequest{
		Token:   session.Token,
		Session: session.Session,
	}

	_json, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", TinkGWUri+"/api/tink-client", bytes.NewBuffer(_json))

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	Sugar.Info("response Status:", resp.Status)
	Sugar.Info("response Headers:", resp.Header)
	Sugar.Info("response Body:", resp.Body)
	return nil
}
