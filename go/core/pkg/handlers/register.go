package handlers

/*

func DoLogin(c *gin.Context, client *client.Client) {
	var json pkg3.LoginRequest
	err := c.BindJSON(&json)
	if err != nil {
		debug.PrintStack()
		c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
		return
	}
	//logger.Info("Executing login for:", zap.Object("json", &json))

	sess, err := client.Login(json.Username, json.Password)
	if err != nil {
		debug.PrintStack()
		c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
		return
	}
	pkg3.PutClientSessionTokenInCache(sess.Id, sess.Token)
	c.JSON(http.StatusOK, gin.H{"sessionId": sess.Id})
}

func DoRegister(c *gin.Context, client *client.Client) {
	var json pkg3.LoginRequest
	err := c.BindJSON(&json)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	//logger.Info("Executing register for:", zap.Object("json", &json))
	sess, err := client.Register(json.Username, json.Password)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = provisionTinkUser(sess, client)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = createUserDir(json.Username)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	pkg3.PutClientSessionTokenInCache(sess.Id, sess.Token)
	c.JSON(http.StatusOK, gin.H{"sessionId": sess.Id})
}


Each user needs
.data/username/accounts.bean
.data/username/transactions/


func createUserDir(username string) error {
	err := os.MkdirAll(fmt.Sprintf("%s/%s", pkg3.DataDirPath, username), os.ModePerm)
	if err != nil {
		pkg.Sugar.Error(err)
		debug.PrintStack()
		return err
	}
	err = os.MkdirAll(fmt.Sprintf("%s/%s/transactions", pkg3.DataDirPath, username), os.ModePerm)
	if err != nil {
		pkg.Sugar.Error(err)
		debug.PrintStack()
		return err
	}
	_, err = os.Create(fmt.Sprintf("%s/%s/accounts.bean", pkg3.DataDirPath, username))
	if err != nil {
		pkg.Sugar.Error(err)
		debug.PrintStack()
		return err
	}
	return nil
}

// Each user needs to be registered as client in Tink organisation
func provisionTinkUser(session client.Session, oauthKeeperClient *client.Client) error {
	body := pkg3.PhiSessionRequest{
		Token:   session.Token,
		Session: session.Session,
	}

	_json, err := json.Marshal(body)
	if err != nil {
		return err
	}

	oauthKeeperClient.SetBearerToken(session.Token)

	resp, err := oauthKeeperClient.SendRequest("POST", OauthKeeperUri+"/tink-user", "application/json", bytes.NewBuffer(_json))
	if err != nil {
		return err
	}

	var res pkg3.PhiClientIdResponse
	err = json.Unmarshal([]byte(resp), &res)
	if err != nil {
		return err
	}

	traits := session.Identity.Traits.(map[string]interface{})
	traits["tink_id"] = res.TinkId

	oryConf := ory.NewConfiguration()
	oryConf.Servers = ory.ServerConfigurations{{URL: pkg3.OriUrl}}
	oryConf.AddDefaultHeader("Authorization", "Bearer "+pkg3.OryToken)
	oryConf.HTTPClient = &http.Client{}
	oryClient := ory.NewAPIClient(oryConf)

	identity, _, err := oryClient.V0alpha2Api.AdminUpdateIdentity(context.Background(), session.Identity.Id).AdminUpdateIdentityBody(
		ory.AdminUpdateIdentityBody{State: *session.Identity.State, Traits: traits}).Execute()
	if err != nil {
		pkg.Sugar.Error(err)
		return err
	}
	session.Identity = *identity
	return nil
}
*/
