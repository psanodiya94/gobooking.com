package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var testData = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"rs", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	{"post-search-availability", "/search-availability", "POST", []postData{
		{"start", "2020-01-01"},
		{"end", "2020-01-02"},
	}, http.StatusOK},
	{"post-search-availability-json", "/search-availability-json", "POST", []postData{
		{"start", "2020-01-01"},
		{"end", "2020-01-02"},
	}, http.StatusOK},
	{"make-reservation", "/make-reservation", "POST", []postData{
		{"first_name", "John"},
		{"last_name", "Smith"},
		{"email", "john.smith@mail.com"},
		{"phone", "555-555-5555"},
	}, http.StatusOK},
}

// TestHandlers tests all the handlers
func TestHandlers(t *testing.T) {
	testRoutes := getRoutes()
	testServer := httptest.NewTLSServer(testRoutes)
	defer testServer.Close()

	for _, test := range testData {
		if test.method == "GET" {
			resp, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("For %s, expected %d but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, param := range test.params {
				values.Add(param.key, param.value)
			}

			resp, err := testServer.Client().PostForm(testServer.URL+test.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("For %s, expected %d but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
