package networker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type ContentType int

const (
	JSON ContentType = iota
	XML
	TEXT
	FORM
)

const (
	GET     = "GET"
	HEAD    = "HEAD"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
)

type Request struct {
	url           string
	method        string
	headers       map[string]string
	cookies       map[string]string
	query         map[string]string
	ctype         ContentType
	data          map[string]interface{}
	raw_data      string
	authUsername  string
	authPassword  string
	isUseBaseAuth bool
}

func New(method string) *Request {
	return &Request{
		method:        string(method),
		headers:       make(map[string]string),
		cookies:       make(map[string]string),
		isUseBaseAuth: false,
	}
}

func Get(url string, query map[string]string) *Request {
	return &Request{
		url:           url,
		headers:       make(map[string]string),
		cookies:       make(map[string]string),
		method:        GET,
		query:         query,
		isUseBaseAuth: false,
	}
}

func Head(url string, query map[string]string) *Request {
	return &Request{
		url:           url,
		headers:       make(map[string]string),
		cookies:       make(map[string]string),
		method:        HEAD,
		query:         query,
		isUseBaseAuth: false,
	}
}

func Post(url string, query map[string]string, ctype ContentType, data interface{}) *Request {
	r := &Request{
		url:           url,
		headers:       make(map[string]string),
		cookies:       make(map[string]string),
		method:        POST,
		query:         query,
		ctype:         ctype,
		isUseBaseAuth: false,
	}
	r.Body(data)
	return r
}

func Put(url string, query map[string]string, ctype ContentType, data interface{}) *Request {
	r := &Request{
		url:           url,
		headers:       make(map[string]string),
		cookies:       make(map[string]string),
		method:        PUT,
		query:         query,
		ctype:         ctype,
		isUseBaseAuth: false,
	}
	r.Body(data)
	return r
}

func Delete(url string, query map[string]string) *Request {
	return &Request{
		url:           url,
		headers:       make(map[string]string),
		cookies:       make(map[string]string),
		method:        DELETE,
		query:         query,
		isUseBaseAuth: false,
	}
}

func Patch(url string, query map[string]string, ctype ContentType, data interface{}) *Request {
	r := &Request{
		url:           url,
		headers:       make(map[string]string),
		cookies:       make(map[string]string),
		method:        PATCH,
		query:         query,
		ctype:         ctype,
		isUseBaseAuth: false,
	}
	r.Body(data)
	return r
}

func Options(url string, query map[string]string) *Request {
	return &Request{
		url:           url,
		headers:       make(map[string]string),
		cookies:       make(map[string]string),
		method:        OPTIONS,
		query:         query,
		isUseBaseAuth: false,
	}
}

func (r *Request) Url(url string) *Request {
	r.url = url
	return r
}

func (r *Request) Header(key string, value string) *Request {
	r.headers[key] = value
	return r
}

func (r *Request) Query(params map[string]string) *Request {
	if r.query != nil {
		for k, v := range params {
			r.query[k] = v
		}
	} else {
		r.query = params
	}
	return r
}

func (r *Request) Cookie(key string, value string) *Request {
	r.cookies[key] = value
	return r
}

func (r *Request) BaseAuth(username string, password string) *Request {
	r.authUsername = username
	r.authPassword = password
	r.isUseBaseAuth = true
	return r
}

func (r *Request) Body(params interface{}) *Request {
	switch v := reflect.ValueOf(params); v.Kind() {
	case reflect.String:
		r.addRawDataAsString(v.String())
	case reflect.Struct:
		r.addDataAsStruct(v.Interface())
	case reflect.Map:
		r.addDataAsMap(params.(map[string]interface{}))
	default:
	}
	return r
}

func (r *Request) addDataAsMap(params map[string]interface{}) {
	if r.data != nil {
		for k, v := range params {
			r.data[k] = v
		}
	} else {
		r.data = params
	}
}

func (r *Request) addDataAsStruct(object interface{}) {
	objectValue := reflect.ValueOf(object)
	objectType := reflect.TypeOf(object)

	if r.data == nil {
		r.data = make(map[string]interface{})
	}

	for i := 0; i < objectType.NumField(); i++ {
		field := objectValue.Field(i)
		typeField := objectType.Field(i)

		r.data[strings.ToLower(typeField.Name)] = fmt.Sprintf("%v", field.Interface())
	}
}

func (r *Request) addRawDataAsString(data string) {
	r.raw_data = data
}

func (r *Request) Do() ([]byte, *http.Response, error) {
	var bodyReader io.Reader
	if r.method != GET && r.method != DELETE {
		bodyReader = bytes.NewReader(r.prepareRequestBody())
	}

	req, err := http.NewRequest(r.method, r.url, bodyReader)
	if err != nil {
		return nil, nil, err
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	q := req.URL.Query()
	for k, v := range r.query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	for k, v := range r.cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	if r.isUseBaseAuth {
		req.SetBasicAuth(r.authUsername, r.authPassword)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	return body, resp, nil
}

func (r *Request) prepareRequestBody() []byte {
	var bodyData []byte
	if r.ctype == JSON {
		bodyData, _ = json.Marshal(r.data)
		r.Header("Content-Type", "application/json")
	}
	if r.ctype == XML {
		bodyData = []byte(r.raw_data)
		r.Header("Content-Type", "application/xml")
	}
	if r.ctype == TEXT {
		bodyData = []byte(r.raw_data)
		r.Header("Content-Type", "text/plain")
	}
	if r.ctype == FORM {
		bodyData = []byte(r.getURLValues().Encode())
		r.Header("Content-Type", "application/x-www-form-urlencoded")
	}
	return bodyData
}

func (r *Request) getURLValues() url.Values {
	values := url.Values{}
	for k, v := range r.data {
		values.Add(k, fmt.Sprintf("%v", v))
	}
	return values
}
