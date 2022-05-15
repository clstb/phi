package handlers

import (
	"context"
	"github.com/clstb/phi/core/internal/config"
	pb "github.com/clstb/phi/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

/*
func (s *CoreServer) LinkBank(ctx *gin.Context) {

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

	gwServiceClient := proto2.NewTinkGWServiceClient(connection)

	user, err := s.GetUserFromCache(json.SessionId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, err)
		return
	}

	res, err := gwServiceClient.CreateTinkLink(context.TODO(), &proto2.TinkIdMessage{TinkId: user.tinkId})
	if err != nil {
		ctx.AbortWithError(http.StatusFailedDependency, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"link": res.TinkLink})
}*/

func (s *CoreServer) AuthCodeLink(ctx *gin.Context) {
	var json LinkRequest
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

	var res *pb.BytesMessage
	if json.Test {
		res, err = gwServiceClient.GetTestAuthLink(context.Background(), &emptypb.Empty{})
	} else {
		res, err = gwServiceClient.GetTestAuthLink(context.Background(), &emptypb.Empty{})
	}
	if err != nil {
		ctx.AbortWithError(s.mapErrorToHttpCode(err), err)
		return
	}

	link := string(res.Arr)
	s.Logger.Info("Retrieved Link ---> ", link)

	ctx.JSON(http.StatusOK, gin.H{"link": link})
}

func (s *CoreServer) ExchangeToToken(ctx *gin.Context) {

	var json AccessCodeRequest
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

	res, err := gwServiceClient.ExchangeAuthCodeToToken(context.Background(), &pb.StringMessage{Value: json.AccessCode})

	if err != nil {
		ctx.AbortWithError(s.mapErrorToHttpCode(err), err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": res.Value})
}
