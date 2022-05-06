package main

import (
	"github.com/clstb/phi/go/pkg/client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
)

func doLogin(c *gin.Context) {
	var json LoginRequest
	err := c.BindJSON(&json)
	if err != nil {
		sugar.Error(zap.Error(err))
		debug.PrintStack()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	logger.Info("Executing login for:", zap.Object("json", &json))
	newClient := client.NewClient("https://phi.clstb.codes")
	sess, err := newClient.Login(json.Username, json.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	newClient.SetBearerToken(sess.Token)
	putClientInCache(sess.Id, newClient)
	c.JSON(http.StatusOK, gin.H{"sessionId": sess.Id})
}
