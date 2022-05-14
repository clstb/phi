package handlers

import (
	"context"
	"fmt"
	"github.com/clstb/phi/go/core/pkg"
	"github.com/clstb/phi/go/core/pkg/client"
	pb "github.com/clstb/phi/go/proto"
	"github.com/gin-gonic/gin"
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

func DoLogin(c *gin.Context, logger *zap.SugaredLogger, authClient *client.AuthClient) {
	var json LoginRequest
	err := c.BindJSON(&json)
	if err != nil {
		logger.Error(err)
		c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
		return
	}
	logger.Debug("Executing login for: %s: $s\n", json.Username, json.Password)

	sess, err := authClient.Login(json.Username, json.Password)
	if err != nil {
		logger.Error(err)
		c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
		return
	}
	PutClientSessionTokenInCache(sess.Id, sess.Token)
	c.JSON(http.StatusOK, gin.H{"sessionId": sess.Id})
}

func DoRegister(c *gin.Context, logger *zap.SugaredLogger, authClient *client.AuthClient) {
	var json LoginRequest
	err := c.BindJSON(&json)
	if err != nil {
		logger.Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	sess, err := authClient.Register(json.Username, json.Password)
	if err != nil {
		logger.Error(err)
		c.AbortWithError(http.StatusFailedDependency, err)
		return
	}
	id, err = provisionTinkUser(sess.Id)
	if err != nil {
		logger.Error(err)
		c.AbortWithError(http.StatusFailedDependency, err)
		return
	}

	err = createUserDir(json.Username, logger)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	PutClientSessionTokenInCache(sess.Id, sess.Token)
	c.JSON(http.StatusOK, gin.H{"sessionId": sess.Id})
}

//Each user needs
//.data/username/accounts.bean
//.data/username/transactions/

func createUserDir(username string, sugar *zap.SugaredLogger) error {
	err := os.MkdirAll(fmt.Sprintf("%s/%s", pkg.DataDirPath, username), os.ModePerm)
	if err != nil {
		sugar.Error(err)
		debug.PrintStack()
		return err
	}
	err = os.MkdirAll(fmt.Sprintf("%s/%s/transactions", pkg.DataDirPath, username), os.ModePerm)
	if err != nil {
		sugar.Error(err)
		debug.PrintStack()
		return err
	}
	_, err = os.Create(fmt.Sprintf("%s/%s/accounts.bean", pkg.DataDirPath, username))
	if err != nil {
		sugar.Error(err)
		debug.PrintStack()
		return err
	}
	return nil
}

func provisionTinkUser(id string) (*string, error) {
	connection, err := grpc.Dial(pkg.TinkGWAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}
	defer connection.Close()

	gwServiceClient := pb.NewTinkGWServiceClient(connection)

	res, err := gwServiceClient.ProvisionMockTinkUser(context.TODO(), &pb.ProvisionTinkUserRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return &res.TinkId, nil
}
