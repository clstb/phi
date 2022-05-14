package handlers

import (
	"context"
	"github.com/clstb/phi/go/core/pkg/config"
	pb "github.com/clstb/phi/go/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func (s *CoreServer) linkBank(ctx *gin.Context) {

	var json SessionId
	err := ctx.BindJSON(&json)
	if err != nil {
		s.Logger.Error(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	connection, err := grpc.Dial(config.TinkGWAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer connection.Close()

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	gwServiceClient := pb.NewTinkGWServiceClient(connection)

	user, err := s.GetUserFromCache(json.SessionId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, err)
		return
	}

	res, err := gwServiceClient.CreateTinkLink(context.TODO(), &pb.TinkIdMessage{TinkId: user.tinkId})
	if err != nil {
		ctx.AbortWithError(http.StatusFailedDependency, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"link": res.TinkLink})
}
