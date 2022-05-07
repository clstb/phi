package main

import (
	"github.com/clstb/phi/go/pkg/client"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func doLogin(c *gin.Context, client *client.Client) {
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
