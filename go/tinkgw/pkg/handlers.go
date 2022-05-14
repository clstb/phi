package pkg

import (
	"context"
	"errors"
	pb "github.com/clstb/phi/go/proto"
	"github.com/clstb/phi/go/tinkgw/pkg/client/tink"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) PrivisionTinkUser(ctx context.Context, r *pb.ProvisionTinkUserRequest) (*pb.ProvisionTinkUserResponse, error) {
	createdUser, err := s.tinkClient.CreateUser(
		r.Id,
		"DE",
		"de_DE",
	)
	if err == nil {
		s.Logger.Info("OK ---> ", createdUser.UserID)
		return &pb.ProvisionTinkUserResponse{TinkId: createdUser.UserID}, nil
	}
	if errors.Is(err, tink.ErrUserExists) {
		s.Logger.Warn(err)
		user, err := s.getUser(r.Id)
		if err != nil {
			s.Logger.Error(err)
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		createdUser.UserID = user.Id
		return &pb.ProvisionTinkUserResponse{TinkId: user.Id}, status.Error(codes.AlreadyExists, err.Error())
	}
	s.Logger.Error(err)
	return nil, status.Error(codes.Internal, err.Error())
}

func (s *Server) PrivisionMockTinkUser(ctx context.Context, r *pb.ProvisionTinkUserRequest) (*pb.ProvisionTinkUserResponse, error) {
	s.Logger.Info("OK ---> b534d4493183487e8e77ce3eeccaae1b")
	return &pb.ProvisionTinkUserResponse{TinkId: "b534d4493183487e8e77ce3eeccaae1b"}, nil
}

/*
func (s *Server) Link(context *gin.Context) {
	logger := s.Logger.With(
		zap.String("handler", "link"),
	)

	session, ok := context.Request.Context().Value("session").(ory.Session)
	if !ok {
		s.Logger.Error("missing session")
		context.AbortWithError(http.StatusUnauthorized, client.NestedHttpError{HttpCode: http.StatusUnauthorized, Description: "missing session"})
		return
	}

	code, err := s.tinkClient.GetAuthorizeGrantDelegateCode(
		"code",
		"",
		session.Identity.Id,
		session.Identity.Id,
		GetAuthorizeGrantDelegateCodeRoles,
	)
	if err != nil {
		logger.Error("tink: authorize grant delegate", zap.Error(err))
		context.AbortWithError(http.StatusFailedDependency, err)
		return
	}

	link := fmt.Sprintf(
		LinkBankAccountUriFormat,
		s.tinkClientId,
		s.callbackURL,
		"DE",    // req.Market TODO,
		"de_DE", // req.Locale TODO,
		code,
	)
	context.Data(http.StatusOK, "text/plain", []byte(link))
}

func (s *Server) getToken(id string) (tink.Token, error) {
	return GetToken(id, s.tinkClient, s.tinkClientId, s.tinkClientSecret)
}

func (s *Server) Token(context *gin.Context) {
	logger := s.Logger.With(
		zap.String("handler", "token"),
	)

	session, ok := context.Value("session").(ory.Session)
	if !ok {
		s.Logger.Error("missing session")
		context.AbortWithError(http.StatusUnauthorized, client.NestedHttpError{
			HttpCode:    http.StatusUnauthorized,
			Description: "missing session",
		})
		return
	}

	token, err := s.getToken(session.Identity.Id)
	if err != nil {
		logger.Error("getting token", zap.Error(err))
		context.AbortWithError(http.StatusFailedDependency,
			client.NestedHttpError{HttpCode: http.StatusFailedDependency, Description: "getting token"},
		)
		return
	}

	context.JSON(http.StatusOK, "Authorization "+"Bearer "+token.AccessToken)
}

func GetToken(id string, tinkClient *tink.Client, tinkClientId string, tinkClientSecret string) (tink.Token, error) {
	code, err := tinkClient.GetAuthorizeGrantCode(
		"",
		id,
		GetAuthorizeGrantCodeRoles,
	)
	if err != nil {
		return tink.Token{}, fmt.Errorf("tink: authorize grant: %w", err)
	}
	token, err := tinkClient.GetToken(
		code,
		"",
		tinkClientId,
		tinkClientSecret,
		"authorization_code",
		GetAuthorizeGrantCodeRoles,
	)
	if err != nil {
		return tink.Token{}, fmt.Errorf("tink: oauth token: %w", err)
	}
	return token, nil
}

*/
