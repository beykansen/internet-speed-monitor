package pkg

import (
	"encoding/json"
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

func Callback(url string, result *Result) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	reqBody, _ := json.Marshal(result)
	req.SetBody(reqBody)

	// Perform the request
	err := fasthttp.DoTimeout(req, resp, 3*time.Second)
	if err != nil {
		log.Printf("Got error from callback url %s \n", err.Error())
	}
}
