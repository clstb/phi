package core

import "os"

const TinkGWUri = "http://localhost:8080"

const OriUrl = "https://youthful-colden-lom8uboh9g.projects.oryapis.com/api/kratos/admin"

const OryToken = "J>ExF`M`=k65&ur"

const DataDirPath = ".data"

const ServerPort = "8099"

var OauthKeeperUri = func() string {
	str := os.Getenv("OAUTHKEEPER_URI")
	if len(str) == 0 {
		panic("OAUTHKEEPER_URI == nill")
	}
	return str
}()
