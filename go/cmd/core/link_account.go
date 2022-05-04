package main

import (
	"context"
	"github.com/clstb/phi/go/pkg/client"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getTinkLink(c *gin.Context) {
	var json Session
	err := c.BindJSON(&json)
	if err != nil {
		sugar.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
		return
	}
	userClient, err := UserClientCache.Get(context.TODO(), json.SessionId)
	if err != nil {
		sugar.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	if userClient == nil {
		sugar.Error("Not logged in")
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Not logged in"})
		return
	}
	link, err := userClient.(*client.Client).GetLink()
	if err != nil {
		sugar.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"link": link})
}

func isAccountLinked(c *gin.Context) {
	name := c.Param("name")
	isLinked := getIsAccountLinked(name)
	c.JSON(http.StatusOK, gin.H{"linked": isLinked})
}

func confirmLinked(c *gin.Context) {
	var json Username
	err := c.BindJSON(&json)
	if err != nil {
		sugar.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
		return
	}
	err = setAccountLinked(json.Username)
	if err != nil {
		sugar.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
