package main

import (
	"github.com/clstb/phi/go/pkg/client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func doRegister(c *gin.Context) {
	var json Login
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
	putSessionIdInCache(json.Username, sess.Id)
	putClientInCache(sess.Id, newClient)
	c.JSON(http.StatusOK, gin.H{"sessionId": sess.Id})
}
