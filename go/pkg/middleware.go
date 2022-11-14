package pkg

import (
	"fmt"
	"github.com/jobala/middleware_pipeline/pipeline"
	"log"
	"net/http"
	"net/http/httputil"
)

type LoggingMiddleware struct{}

func (s LoggingMiddleware) Intercept(pipeline pipeline.Pipeline, req *http.Request) (*http.Response, error) {
	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("REQUEST:")
	fmt.Println("-------------------------------")
	fmt.Println(string(reqDump))
	fmt.Println("-------------------------------")

	res, err := pipeline.Next(req)

	resDump, err := httputil.DumpResponse(res, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("RESPONSE:")
	fmt.Println("-------------------------------")
	fmt.Println(string(resDump))
	fmt.Println("-------------------------------")
	return res, err
}
