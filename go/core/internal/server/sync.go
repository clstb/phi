package server

import (
	"context"
	proto2 "github.com/clstb/phi/proto"
	"github.com/gin-gonic/gin"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
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

	err = s.doSyncRPC(json.Username, json.AccessToken)

	if err != nil {
		s.Logger.Error(err)
		ctx.AbortWithStatusJSON(http.StatusFailedDependency, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (s *CoreServer) doSyncRPC(username string, token string) error {
	connection, err := grpc.Dial(s.LedgerUri, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(grpczap.StreamClientInterceptor(s.Logger.Desugar())),
	)
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
