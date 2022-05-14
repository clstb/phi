package server

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

	code, err := s.tinkClient.GetDelegatedAutorizationCode(
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
