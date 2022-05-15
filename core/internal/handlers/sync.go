package handlers

import (
	"context"
	"github.com/clstb/phi/core/internal/config"
	proto2 "github.com/clstb/phi/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func (s *CoreServer) SyncLedger(ctx *gin.Context) {

	var json SyncRequest
	err := ctx.BindJSON(&json)
	if err != nil {
		s.Logger.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	//user, err := s.GetUserFromCache(json.)
	//if err != nil {
	//	ctx.AbortWithStatusJSON(http.StatusNotFound, err)
	//	return
	//}

	err = doSyncRPC(json.Username, json.AccessToken)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusFailedDependency, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func doSyncRPC(username string, token string) error {
	connection, err := grpc.Dial(config.LedgerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer connection.Close()

	gwServiceClient := proto2.NewBeanAccountServiceClient(connection)
	_, err = gwServiceClient.SyncLedger(context.TODO(), &proto2.SyncMessage{Username: username, Token: token})
	if err != nil {
		return err
	}
	return nil
}
