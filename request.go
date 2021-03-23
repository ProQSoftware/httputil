package httputil

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func SendRequest(method, url, contentType string, body io.Reader) (resp *http.Response, err error) {
	method = strings.ToUpper(method)
	switch method {
	case "POST":
		resp, err = http.Post(url, contentType, body)
	case "GET":
		if body == nil {
			resp, err = http.Get(url)
			break
		}
		fallthrough
	default:
		var req *http.Request

		if req, err = http.NewRequest(method, url, body); err != nil {
			err = NewError(http.StatusInternalServerError, "text/html", err.Error())
			return
		} else if body != nil {
			req.Header.Add("Content-Type", contentType)
		}

		resp, err = http.DefaultClient.Do(req)
	}

	if err != nil {
		err = NewError(http.StatusServiceUnavailable, "text/html", err.Error())
	} else if resp.StatusCode >= 400 {
		if msg, e := ioutil.ReadAll(resp.Body); e != nil {
			panic(e)
		} else {
			err = NewError(resp.StatusCode, resp.Header.Get("Content-Type"), string(msg))
		}

		_ = resp.Body.Close()

		resp = nil
	}

	return
}
