package handlers

import (
	"context"
	"github.com/clstb/phi/go/core/config"
	proto2 "github.com/clstb/phi/go/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func (s *CoreServer) SyncLedger(ctx *gin.Context) {

	var json SessionId
	err := ctx.BindJSON(&json)
	if err != nil {
		s.Logger.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user, err := s.GetUserFromCache(json.SessionId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, err)
		return
	}

	err = doSyncRPC(user.username)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusFailedDependency, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func doSyncRPC(username string) error {
	connection, err := grpc.Dial(config.LedgerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer connection.Close()

	gwServiceClient := proto2.NewBeanAccountServiceClient(connection)
	_, err = gwServiceClient.SyncLedger(context.TODO(), &proto2.UserNameMessage{Username: username})
	if err != nil {
		return err
	}
	return nil
}
