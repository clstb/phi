package core

import "os"

const TinkGWUri = "http://localhost:8080"

const DataDirPath = ".data"

const ServerPort = "8099"

var OauthKeeperUri = func() string {
	str := os.Getenv("OAUTHKEEPER_URI")
	if len(str) == 0 {
		panic("OAUTHKEEPER_URI == nill")
	}
	return str
}()
