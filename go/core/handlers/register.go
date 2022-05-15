package handlers

import (
	"context"
	"github.com/clstb/phi/go/core/config"
	proto2 "github.com/clstb/phi/go/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

func (s *CoreServer) DoRegister(c *gin.Context) {
	var json LoginRequest
	err := c.BindJSON(&json)
	if err != nil {
		s.Logger.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tinkId, err := s.provisionTinkUser()

	if err != nil {
		s.Logger.Error(err)
		c.AbortWithError(s.mapErrorToHttpCode(err), err)
		return
	}

	err = provisionFS(json.Username)
	if err != nil {
		s.Logger.Error(err)
		c.AbortWithError(s.mapErrorToHttpCode(err), err)
		return
	}

	sess, err := s.AuthClient.Register(json.Username, json.Password, *tinkId)
	if err != nil {
		s.Logger.Error(err)
		c.AbortWithError(s.mapErrorToHttpCode(err), err)
		return
	}

	if err != nil {
		s.Logger.Error(err)
		c.AbortWithError(s.mapErrorToHttpCode(err), err)
		return
	}

	traits := sess.Identity.Traits.(map[string]string)
	s.PutUserInCache(sess.Id, UserDetails{
		tinkId:   traits["tink_id"],
		username: traits["username"],
	})
	c.JSON(http.StatusOK, gin.H{"sessionId": sess.Id})
}

func (s *CoreServer) provisionTinkUser() (*string, error) {
	connection, err := grpc.Dial(config.TinkGWAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}
	defer connection.Close()

	gwServiceClient := proto2.NewTinkGWServiceClient(connection)

	//res, err := gwServiceClient.ProvisionMockTinkUser(context.TODO(), &emptypb.Empty{})
	res, err := gwServiceClient.ProvisionTinkUser(context.TODO(), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return &res.TinkId, nil
}

func provisionFS(username string) error {

	connection, err := grpc.Dial(config.LedgerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return err
	}
	defer connection.Close()

	gwServiceClient := proto2.NewBeanAccountServiceClient(connection)

	_, err = gwServiceClient.ProvisionFSStructure(context.TODO(), &proto2.StringMessage{Value: username})
	if err != nil {
		return err
	}
	return nil

}
