package httpserver

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fxlibraries/loggers"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/url"
	"sync/atomic"
)

var requestID uint32

type Request struct {
	ID          uint32            `json:ID,omitempty"`
	Method      string            `json:"method,omitempty"`
	Body        io.Reader         `json:"-"`
	BodyBuff    *bytes.Buffer     `json:"body,omitempty"`
	RemoteAddr  string            `json:"remoteAddr,omitempty"`
	QueryParams url.Values        `json:"queryPrams,omitempty"`
	UrlParams   map[string]string `json:"URLParams,omitempty"`
	URL         *url.URL          `json:"URL,omitempty"`
	Header      http.Header       `json:"header,omitempty"`
}

func init() {
	requestID = 0
}

// NewRequest New a Request from http.Request
func NewRequest(r *http.Request) *Request {
	var request Request
	request.ID = atomic.AddUint32(&requestID, 1)
	request.Method = r.Method
	request.Body = r.Body
	request.BodyBuff = new(bytes.Buffer)
	request.BodyBuff.ReadFrom(r.Body)
	request.RemoteAddr = r.Header.Get("X-Forwarded-For")
	request.Header = r.Header
	if request.RemoteAddr == "" {
		request.RemoteAddr = r.RemoteAddr
	}
	request.UrlParams = mux.Vars(r)
	request.QueryParams = r.URL.Query()
	return &request
}

// Parse Parse the JSON-encoded Request.Body store the result in the value pointed to by v
// If cache is true, cache the body in BodyBuff
func (self *Request) Parse(v interface{}) error {
	return json.Unmarshal(self.BodyBuff.Bytes(), v)
}

func (self *Request) ParseCache(v interface{}) error {
	return json.Unmarshal(self.BodyBuff.Bytes(), v)
}

func (self *Request) ParseByXML(v interface{}) error {
	loggers.Debug.Println(string(self.BodyBuff.Bytes()))
	return xml.Unmarshal(self.BodyBuff.Bytes(), v)
}
