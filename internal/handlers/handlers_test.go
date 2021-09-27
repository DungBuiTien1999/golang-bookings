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

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"general_quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors_suite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"search_availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"make_reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"reservation_summary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	{"search_availability", "/search-availability", "POST", []postData{
		{key: "start", value: "2021-09-26"},
		{key: "end", value: "2021-09-28"},
	}, http.StatusOK},
	{"search_availability_json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2021-09-26"},
		{key: "end", value: "2021-09-28"},
	}, http.StatusOK},
	{"make_reservation", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "Dung"},
		{key: "last_name", value: "Bui"},
		{key: "email", value: "dungbui@gmail.com"},
		{key: "phone", value: "0328816614"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
