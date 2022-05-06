package main

import (
	"fmt"
	"github.com/clstb/phi/go/pkg/client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"runtime/debug"
)

func doRegister(c *gin.Context) {
	var json LoginRequest
	err := c.BindJSON(&json)
	if err != nil {
		logger.Error("Error unmarshalling JSON", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	logger.Info("Executing register for:", zap.Object("json", &json))
	newClient := client.NewClient("https://phi.clstb.codes")
	sess, err := newClient.Register(json.Username, json.Password)
	if err != nil {
		sugar.Error(zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	err = createUserDir(json.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
	}
	putClientInCache(sess.Id, newClient)
	c.JSON(http.StatusOK, gin.H{"sessionId": sess.Id})
}

/*
Each user needs
.data/username/accounts.bean
.data/username/transactions/
*/

var DataDirPath = ".data"

func createUserDir(username string) error {
	err := os.MkdirAll(fmt.Sprintf("%s/%s", DataDirPath, username), os.ModePerm)
	if err != nil {
		sugar.Error(err)
		debug.PrintStack()
		return err
	}
	err = os.MkdirAll(fmt.Sprintf("%s/%s/transactions", DataDirPath, username), os.ModePerm)
	if err != nil {
		sugar.Error(err)
		debug.PrintStack()
		return err
	}
	_, err = os.Create(fmt.Sprintf("%s/%s/accounts.bean", DataDirPath, username))
	if err != nil {
		sugar.Error(err)
		debug.PrintStack()
		return err
	}
	return nil
}
