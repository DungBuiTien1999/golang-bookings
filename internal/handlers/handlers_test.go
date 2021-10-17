package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DungBuiTien1999/bookings/internal/driver"
	"github.com/DungBuiTien1999/bookings/internal/models"
)

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"general_quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors_suite", "/majors-suite", "GET", http.StatusOK},
	{"search_availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"non-exist", "/haha/hoho", "GET", http.StatusNotFound},
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"admin dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"new reservations", "/admin/reservations-new", "GET", http.StatusOK},
	{"all reservations", "/admin/reservations-all", "GET", http.StatusOK},
	{"show reservation", "/admin/reservations/new/1/show", "GET", http.StatusOK},
	{"show reservation calender", "/admin/reservations-calendar?y=2021&m=10", "GET", http.StatusOK},
	{"handle process mark with year", "/admin/process-reservation/new/1/do?y=2021&m=10", "GET", http.StatusOK},
	{"handle process mark without year", "/admin/process-reservation/new/1/do", "GET", http.StatusOK},
	{"handle delete reservation with year", "/admin/delete-reservation/new/1/do?y=2021&m=10", "GET", http.StatusOK},
	{"handle delete reservation", "/admin/delete-reservation/new/1/do", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestNewRepo(t *testing.T) {
	var db driver.DB
	testRepo := NewRepo(&app, &db)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo, got %s, wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		ID:     1,
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarter",
		},
	}
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset initial session)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 3
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, "2050-01-01")
	endDate, _ := time.Parse(layout, "2050-01-03")
	reservation := models.Reservation{
		RoomID:    1,
		StartDate: startDate,
		EndDate:   endDate,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarter",
		},
	}

	postData := url.Values{}
	postData.Add("first_name", "dung")
	postData.Add("last_name", "bui")
	postData.Add("email", "dung@gmail.com")
	postData.Add("phone", "023186753")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case where reservation is not in session (reset initial session)
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code for missing reservation session: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code for missing body: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case invalid form
	postData = url.Values{}
	postData.Add("first_name", "d")
	postData.Add("last_name", "bui")
	postData.Add("email", "dung@gmail.com")
	postData.Add("phone", "023186753")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code for invalid form data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case insert reservation failure
	reservation.RoomID = 2
	postData = url.Values{}
	postData.Add("first_name", "dung")
	postData.Add("last_name", "bui")
	postData.Add("email", "dung@gmail.com")
	postData.Add("phone", "023186753")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code for insert reservation: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case insert restriction room failure
	reservation.RoomID = 1000
	postData = url.Values{}
	postData.Add("first_name", "dung")
	postData.Add("last_name", "bui")
	postData.Add("email", "dung@gmail.com")
	postData.Add("phone", "023186753")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code for insert restriction room: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostAvailability handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case invalid form
	req, _ = http.NewRequest("POST", "/search-availability", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler returned wrong response code for missing form data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case failure to parse start date
	reqBody = "start=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case failure to parse end date
	reqBody = "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler returned wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case failure to search available rooms
	reqBody = "start=2000-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2000-01-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler returned wrong response code for search available rooms: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case when have non-existently available room
	reqBody = "start=2050-10-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2000-10-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler returned wrong response code for no available room: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	reservation := models.Reservation{}
	req, _ := http.NewRequest("POST", "/search-availability", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case missing reservation session
	req, _ = http.NewRequest("POST", "/search-availability", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ReservationSummary handler returned wrong response code for missing reservation session: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_PostAvailabilityJSON(t *testing.T) {
	// first case - rooms are not available
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=3")

	// create request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// make handler handleFunc
	handler := http.HandlerFunc(Repo.PostAvailabilityJSON)

	// get response recorder
	rr := httptest.NewRecorder()

	// make request to handler
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.Bytes()), &j)
	if err != nil {
		t.Error("failed to parse json - case have not availible room")
	}
	if j.OK {
		t.Error("expected rooms are not available")
	}

	// second case - rooms are available
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	err = json.Unmarshal([]byte(rr.Body.Bytes()), &j)
	if err != nil {
		t.Error("failed to parse json - case have not availible room")
	}
	if !j.OK {
		t.Error("expected rooms are available")
	}

	// third case - missing form data
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	err = json.Unmarshal([]byte(rr.Body.Bytes()), &j)
	if err != nil {
		t.Error("failed to parse json - case missing form data")
	}
	if j.OK {
		t.Error("expected error because missing form data")
	}

	// four case - failure to connect to database
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=3")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	err = json.Unmarshal([]byte(rr.Body.Bytes()), &j)
	if err != nil {
		t.Error("failed to parse json - case failure to connect to database")
	}
	if j.OK {
		t.Error("expected error for failure connect to database")
	}

}

func TestRepository_ChooseRoom(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	req.RequestURI = "/choose/1"
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseRoom)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case failure to parse room id
	req, _ = http.NewRequest("GET", "/choose-room/invalid", nil)
	req.RequestURI = "/choose/invalid"
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.ChooseRoom)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code for invalid roomID: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case missing reservation session
	req, _ = http.NewRequest("GET", "/choose-room/1", nil)
	req.RequestURI = "/choose/1"
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	handler = http.HandlerFunc(Repo.ChooseRoom)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code for missing reservation session: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_BookRoom(t *testing.T) {
	req, _ := http.NewRequest("GET", "/book-room?id=1&sd=2050-01-01&ed=2050-01-03", nil)
	req.RequestURI = "/book-room?id=1&sd=2050-01-01&ed=2050-01-03"
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	handler := http.HandlerFunc(Repo.BookRoom)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case failure to connect to database
	req, _ = http.NewRequest("GET", "/book-room?id=3&sd=2050-01-01&ed=2050-01-03", nil)
	req.RequestURI = "/book-room?id=3&sd=2050-01-01&ed=2050-01-03"
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	handler = http.HandlerFunc(Repo.BookRoom)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code for failure to connect to db: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid-credentials",
		"me@hehe.com",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"dung@hehe.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-credentials",
		"d",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestLogin(t *testing.T) {
	// range through all tests
	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		// create request
		req, _ := http.NewRequest("POST", "user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			// get the URL from the test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

func TestAdminPostShowReservation(t *testing.T) {
	// case valid without redirect to calendar
	formData := url.Values{}
	formData.Add("year", "")
	formData.Add("month", "")
	formData.Add("first_name", "dung")
	formData.Add("last_name", "bui")
	formData.Add("email", "dung@gmail.com")
	formData.Add("phone", "320-334-8878")

	req, _ := http.NewRequest("POST", "/admin/reservations/new/1", strings.NewReader(formData.Encode()))
	req.RequestURI = "/admin/reservations/new/1"
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.AdminPostShowReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("failed valid without redirect to calendar: expected code %d, but got %d", http.StatusSeeOther, rr.Code)
	}
	actualLoc, _ := rr.Result().Location()
	if actualLoc.String() != "/admin/reservations-new" {
		t.Errorf("failed valid without redirect to calendar: expected location /admin/reservations-new, but got %s", actualLoc.String())
	}

	// case valid with redirect to calendar
	formData = url.Values{}
	formData.Add("year", "2021")
	formData.Add("month", "10")
	formData.Add("first_name", "dung")
	formData.Add("last_name", "bui")
	formData.Add("email", "dung@gmail.com")
	formData.Add("phone", "320-334-8878")

	req, _ = http.NewRequest("POST", "/admin/reservations/cal/1", strings.NewReader(formData.Encode()))
	req.RequestURI = "/admin/reservations/cal/1"
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.AdminPostShowReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("failed valid without redirect to calendar: expected code %d, but got %d", http.StatusSeeOther, rr.Code)
	}
	actualLoc, _ = rr.Result().Location()
	if actualLoc.String() != "/admin/reservations-calendar?y=2021&m=10" {
		t.Errorf("failed valid without redirect to calendar: expected location /admin/reservations-calendar?y=2021&m=10, but got %s", actualLoc.String())
	}
}

func TestAdminPostCalendarReservations(t *testing.T) {
	formData := url.Values{}
	formData.Add("y", "2021")
	formData.Add("m", "10")
	formData.Add("remove_block_1_2021-10-10", "2")
	formData.Add("remove_block_2_2021-10-10", "5")
	formData.Add("add_block_2_2021-10-12", "3")

	blockMap1 := make(map[string]int)
	blockMap1["2021-10-10"] = 2
	blockMap1["2021-10-12"] = 4

	blockMap2 := make(map[string]int)
	blockMap2["2021-10-10"] = 5

	req, _ := http.NewRequest("POST", "/admin/reservations-calendar", strings.NewReader(formData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "block_map_1", blockMap1)
	session.Put(ctx, "block_map_2", blockMap2)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.AdminPostCalendarReservations)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("failed case valid date: expected code %d, but got %d", http.StatusSeeOther, rr.Code)
	}
	actualLoc, _ := rr.Result().Location()
	if actualLoc.String() != "/admin/reservations-calendar?y=2021&m=10" {
		t.Errorf("failed case valid date: expected location /admin/reservations-calendar?y=2021&m=10, but got %s", actualLoc.String())
	}
}

func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println()
	}
	return ctx
}
