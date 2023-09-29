package utils

import (
	"fmt"
	"io"
	"net/http"
	url2 "net/url"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/utils/Return"
)

type Http struct {
	url      *url2.URL
	request  *http.Request
	response *http.Response
	headers  map[string]string
	Error    Return.Error
}

func NewHttp() Http {
	return Http{
		url:      nil,
		request:  nil,
		response: nil,
		headers:  make(map[string]string),
		Error:    Return.New(),
	}
}

func (h *Http) SetUrl(format string, args ...any) Return.Error {
	for range Only.Once {
		var e error

		h.url, e = url2.Parse(fmt.Sprintf(format, args...))
		if e != nil {
			h.Error.SetError(e)
			break
		}
	}

	return h.Error
}

func (h *Http) Get() ([]byte, Return.Error) {
	var body []byte

	for range Only.Once {
		var e error
		h.request, e = http.NewRequest(http.MethodGet, h.url.String(), nil)
		if e != nil {
			h.Error.SetError(e)
			break
		}

		for key, value := range h.headers {
			h.request.Header.Set(key, value)
		}

		response, e := http.DefaultClient.Do(h.request)
		if e != nil {
			h.Error.SetError(e)
			break
		}

		if response.StatusCode != http.StatusOK {
			h.Error.SetError("Something went wrong: [%d]%s", response.StatusCode, response.Status)
			break
		}

		body, e = io.ReadAll(response.Body)
		if e != nil {
			h.Error.SetError(e)
			break
		}
	}

	return body, h.Error
}

func (h *Http) SetHeader(key string, value string) Return.Error {
	h.headers[key] = value
	return h.Error
}
