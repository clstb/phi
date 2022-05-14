package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *CoreServer) DoLogin(c *gin.Context) {
	var json LoginRequest
	err := c.BindJSON(&json)
	if err != nil {
		s.Logger.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	s.Logger.Debug("Executing login for: %s: $s\n", json.Username, json.Password)

	sess, err := s.AuthClient.Login(json.Username, json.Password)
	if err != nil {
		s.Logger.Error(err)
		c.AbortWithError(s.mapErrorToHttpCode(err), err)
		return
	}
	s.PutClientSessionTokenInCache(sess.Id, sess.Token)
	c.JSON(http.StatusOK, gin.H{"sessionId": sess.Id})
}
