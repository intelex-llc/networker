package networker

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRequest(t *testing.T) {
	ts := startTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != GET {
			t.Errorf("Expected request method %s but got %s", GET, r.Method)
		}
	})
	defer ts.Close()

	Get(ts.URL+"/index", nil).Do()
}

func TestQueryRequest(t *testing.T) {
	ts := startTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != GET {
			t.Errorf("Expected request method %s but got %s", GET, r.Method)
		}
		switch r.URL.Path {
		case "/query1":
			q := r.URL.Query()
			if q["foo"][0] != "bar" {
				t.Error("Expected 'foo:bar' but got", q["foo"][0])
			}
		case "/query2":
			q := r.URL.Query()
			if q["foo2"][0] != "bar2" {
				t.Error("Expected 'foo2:bar2' but got", q["foo2"][0])
			}
		}
	})
	defer ts.Close()

	Get(ts.URL+"/query1", map[string]string{"foo": "bar"}).Do()

	New(GET).
		Url(ts.URL + "/query2").
		Query(map[string]string{"foo2": "bar2"}).
		Do()
}

func TestBaseAuth(t *testing.T) {
	ts := startTestServer(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header["Authorization"]
		if len(auth) == 0 {
			t.Errorf("Expected base auth %q", r.Header["Authorization"])
		}
	})
	defer ts.Close()

	Get(ts.URL+"/auth", nil).BaseAuth("foo", "bar").Do()
}

func TestHeaders(t *testing.T) {
	ts := startTestServer(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/headers":
			if r.Header.Get("foo") != "bar" {
				t.Errorf("Expected header foo = bar got %q", r.Header.Get("foo"))
			}
		case "/cookies":
			if r.Header.Get("Cookie") != "foo=bar" {
				t.Errorf("Expected cookie foo=bar got %q", r.Header.Get("Cookie"))
			}
		}
	})
	defer ts.Close()

	Get(ts.URL+"/headers", nil).Header("foo", "bar").Do()
	Get(ts.URL+"/cookies", nil).Cookie("foo", "bar").Do()
}

func TestPostRequest(t *testing.T) {
	ts := startTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != POST {
			t.Errorf("Expected request method %s but got %s", POST, r.Method)
		}
		switch r.URL.Path {
		case "/json":
			body, _ := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if string(body) != `{"foo":"bar"}` {
				t.Error("Expected JSON body {\"foo\":\"bar\"} but got", string(body))
			}
		case "/xml":
			body, _ := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if string(body) != "<language>Golang</language>" {
				t.Error("Expected XML string '<language>Golang</language>' but got", string(body))
			}
		case "/form":
			body, _ := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if string(body) != "foo=bar" {
				t.Error("Expected Form data 'foo=bar' but got", string(body))
			}
		}
	})
	defer ts.Close()

	Post(ts.URL+"/json", nil, JSON, map[string]interface{}{"foo": "bar"}).Do()
	Post(ts.URL+"/xml", nil, XML, "<language>Golang</language>").Do()
	Post(ts.URL+"/form", nil, FORM, struct{ Foo string }{"bar"}).Do()
}

func startTestServer(callback func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(callback))
}
