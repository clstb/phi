package core

import (
	"context"
	"github.com/clstb/phi/go/pkg/client"
	"github.com/gin-gonic/gin"
	"net/http"
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
		Sugar.Error("Not logged in")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "Not logged in"})
		return
	}

	client.SetBearerToken(token.(string))
	link, err := client.GetLink()
	if err != nil {
		Sugar.Error(err)
		debug.PrintStack()
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"link": link})
}
