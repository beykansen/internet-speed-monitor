package pkg

import (
	"encoding/json"
	"time"

	"github.com/valyala/fasthttp"
)

func Callback(url string, result *Result) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	reqBody, _ := json.Marshal(result)
	req.SetBody(reqBody)

	return fasthttp.DoTimeout(req, resp, 3*time.Second)
}
