package tink

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type LoggingClient struct {
	httpClient *http.Client
}

func (c *LoggingClient) Get(url string) (resp *http.Response, err error) {
	fmt.Printf("GET -------> %s\n", url)
	res, err := c.httpClient.Get(url)
	if err != nil {
		fmt.Println("------------------")
		fmt.Printf("Error: %v\n", err)
		fmt.Println("------------------")
	}
	fmt.Printf("Status code ------> %v\n", res.StatusCode)
	var b bytes.Buffer
	_, err = b.ReadFrom(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Response Body: %s\n", b.String())
	return res, err
}

func (c *LoggingClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	var b bytes.Buffer
	_, err = b.ReadFrom(body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("POST -------> %s\n", url)
	fmt.Printf("Body: %v\n", b.String())

	res, err := c.httpClient.Post(url, contentType, body)
	if err != nil {
		fmt.Println("------------------")
		fmt.Printf("Error: %v\n", err)
		fmt.Println("------------------")
		return res, err
	}

	fmt.Printf("Status code ------> %v\n", res.StatusCode)
	var bb bytes.Buffer
	_, err = bb.ReadFrom(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Response Body: %s\n", bb.String())
	return res, err
}

func (c *LoggingClient) PostForm(url string, data url.Values) ([]byte, error) {
	fmt.Printf("POST -------> %s\n", url)
	fmt.Printf("Body: %v\n", data)
	res, err := c.httpClient.PostForm(url, data)
	if err != nil {
		fmt.Println("------------------")
		fmt.Printf("Error: %v\n", err)
		fmt.Println("------------------")
		return nil, err
	}
	fmt.Printf("Status code ------> %v\n", res.StatusCode)

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response Body: %s\n", buf)
	return buf, err
}
