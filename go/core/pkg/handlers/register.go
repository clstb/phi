package handlers

import (
	"context"
	"fmt"
	"github.com/clstb/phi/go/core/pkg"
	"github.com/clstb/phi/go/core/pkg/client"
	pb "github.com/clstb/phi/go/proto"
	"github.com/gin-gonic/gin"
	ory "github.com/ory/kratos-client-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"runtime/debug"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TinkGwClient struct {
	cc grpc.ClientConnInterface
}

type CoreServer struct {
	AuthClient *client.AuthClient
	Logger     *zap.SugaredLogger
}

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
	PutClientSessionTokenInCache(sess.Id, sess.Token)
	c.JSON(http.StatusOK, gin.H{"sessionId": sess.Id})
}

func (s *CoreServer) DoRegister(c *gin.Context) {
	var json LoginRequest
	err := c.BindJSON(&json)
	if err != nil {
		s.Logger.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	sess, err := s.AuthClient.Register(json.Username, json.Password)
	if err != nil {
		s.Logger.Error(err)
		c.AbortWithError(s.mapErrorToHttpCode(err), err)
		return
	}
	err = s.provisionTinkUser(sess)
	if err != nil {
		s.Logger.Error(err)
		c.AbortWithError(s.mapGRPCErrorToHttpCode(err), err)
		return
	}

	err = createUserDir(json.Username)
	if err != nil {
		s.Logger.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	PutClientSessionTokenInCache(sess.Id, sess.Token)
	c.JSON(http.StatusOK, gin.H{"sessionId": sess.Id})
}

//Each user needs
//.data/username/accounts.bean
//.data/username/transactions/

func createUserDir(username string) error {
	err := os.MkdirAll(fmt.Sprintf("%s/%s", pkg.DataDirPath, username), os.ModePerm)
	if err != nil {
		debug.PrintStack()
		return err
	}
	err = os.MkdirAll(fmt.Sprintf("%s/%s/transactions", pkg.DataDirPath, username), os.ModePerm)
	if err != nil {
		debug.PrintStack()
		return err
	}
	_, err = os.Create(fmt.Sprintf("%s/%s/accounts.bean", pkg.DataDirPath, username))
	if err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func (s *CoreServer) provisionTinkUser(sess client.Session) error {
	connection, err := grpc.Dial(pkg.TinkGWAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return err
	}
	defer connection.Close()

	gwServiceClient := pb.NewTinkGWServiceClient(connection)

	res, err := gwServiceClient.ProvisionMockTinkUser(context.TODO(), &pb.ProvisionTinkUserRequest{Id: sess.Id})
	if err != nil {
		return err
	}

	traits := sess.Identity.Traits.(map[string]interface{})

	traits["tink_id"] = res.TinkId
	_, _, err = s.AuthClient.OryClient.V0alpha2Api.AdminUpdateIdentity(
		context.Background(),
		sess.Identity.Id,
	).AdminUpdateIdentityBody(ory.AdminUpdateIdentityBody{
		State:  *sess.Identity.State,
		Traits: traits,
	}).Execute()
	if err != nil {
		s.Logger.Error("ory: updating identity", zap.Error(err))
		return err
	}
	return nil
}
