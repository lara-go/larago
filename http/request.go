package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	net_http "net/http"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
)

// Request handles http request.
type Request struct {
	request *net_http.Request
	Route   *Route
}

// NewRequest constructor.
func NewRequest(netRequest *net_http.Request) *Request {
	return &Request{
		request: netRequest,
	}
}

// BaseRequest returns base net/http request.
func (r *Request) BaseRequest() *net_http.Request {
	return r.request
}

// IsAjax checks if request was made via ajax.
func (r *Request) IsAjax() bool {
	return r.Header("HTTP_X_REQUESTED_WITH") == "XMLHttpRequest"
}

// Header returns header value.
func (r *Request) Header(name string) string {
	value := r.request.Header.Get(name)

	if value != "" {
		return strings.TrimSpace(value)
	}

	return value
}

// Method returns method name.
func (r *Request) Method() string {
	return r.request.Method
}

// URL returns requested URL.
func (r *Request) URL() string {
	return r.request.RequestURI
}

// IP tries to return real client IP.
func (r *Request) IP() string {
	// Try to get IP from X-Real-IP header.
	realIP := r.Header("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Try to get IP from X-Forwarded-For header.
	realIP = r.Header("X-Forwarded-For")
	idx := strings.IndexByte(realIP, ',')
	if idx >= 0 {
		realIP = realIP[0:idx]
	}
	realIP = strings.TrimSpace(realIP)
	if realIP != "" {
		return realIP
	}

	// Get IP from base request.
	addr := strings.TrimSpace(r.request.RemoteAddr)
	if len(addr) == 0 {
		return ""
	}

	// If address contains port, split it out.
	if ip, _, err := net.SplitHostPort(addr); err == nil {
		return ip
	}

	return addr
}

// HeaderContains checks if header contains requested substring.
func (r *Request) HeaderContains(header, substring string) bool {
	return strings.Contains(r.Header(header), substring)
}

// WantsJSON checks if client wants JSON answer.
func (r *Request) WantsJSON() bool {
	return r.HeaderContains("accept", "application/json")
}

// WantsHTML checks if client wants HTML answer.
func (r *Request) WantsHTML() bool {
	return r.HeaderContains("accept", "text/html")
}

// WantsPlainText checks if client wants plain text answer.
func (r *Request) WantsPlainText() bool {
	return r.HeaderContains("accept", "text/plain")
}

// Param returns route param.
func (r *Request) Param(name string) string {
	return r.Route.Params.ByName(name)
}

// ReadForm unmarshal form request to the structure.
func (r *Request) ReadForm(target interface{}) error {
	decoder := schema.NewDecoder()

	// Decode form values.
	if err := decoder.Decode(target, r.FormValues()); err != nil {
		return err
	}

	return nil
}

// FormValues returns all form values.
func (r *Request) FormValues() url.Values {
	if err := r.parseForm(); err != nil {
		return nil
	}

	return r.request.Form
}

// Parse form values.
func (r *Request) parseForm() error {
	if r.request.Form == nil {
		if err := r.request.ParseForm(); err != nil {
			return err
		}
	}

	return nil
}

// ReadJSON unmarshal json request to the structure.
func (r *Request) ReadJSON(target interface{}) error {
	rawBody, err := r.readBody()
	if err != nil {
		return err
	}

	// Decode JSON body.
	if err := json.Unmarshal(rawBody, target); err != nil {
		return err
	}

	return nil
}

// Read raw body.
func (r *Request) readBody() ([]byte, error) {
	if r.request.Body == nil {
		return nil, errors.New("Body was empty")
	}

	// Read raw body from request.
	rawBody, err := ioutil.ReadAll(r.request.Body)
	if err != nil {
		return nil, err
	}

	// Return parsed body back to base request.
	r.request.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))

	return rawBody, nil
}
