package handlers

import (
	"context"
	"encoding/json"
	"github.com/usmanzaheer1995/bed-and-breakfast/internal/models"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},

	//{"post-search-avail", "/search-availability", "POST", []postData{
	//	{key: "start", value: "2020-01-01"},
	//	{key: "end", value: "2020-01-02"},
	//}, http.StatusOK},
}

var postReservationTests = []struct {
	name               string
	params             []postData
	expectedStatusCode int
}{
	{"happy-path", []postData{
		{key: "start_date", value: "2050-01-01"},
		{key: "end_date", value: "2050-01-02"},
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "john@smith.com"},
		{key: "phone", value: "123123456"},
		{key: "room_id", value: "1"},
		{key: "room_name", value: "General's Quarters"},
	}, http.StatusSeeOther},
	{"invalid-post-data", nil, http.StatusTemporaryRedirect},
	{"invalid-start-date", []postData{
		{key: "start_date", value: "invalid"},
		{key: "end_date", value: "2050-01-02"},
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "john@smith.com"},
		{key: "phone", value: "123123456"},
		{key: "room_id", value: "1"},
		{key: "room_name", value: "General's Quarters"},
	}, http.StatusSeeOther},
	{"invalid-end-date", []postData{
		{key: "start_date", value: "2050-01-01"},
		{key: "end_date", value: "invalid"},
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "john@smith.com"},
		{key: "phone", value: "123123456"},
		{key: "room_id", value: "1"},
		{key: "room_name", value: "General's Quarters"},
	}, http.StatusSeeOther},
	{"invalid-room-id", []postData{
		{key: "start_date", value: "2050-01-01"},
		{key: "end_date", value: "2050-01-02"},
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "john@smith.com"},
		{key: "phone", value: "123123456"},
		{key: "room_id", value: "invalid"},
		{key: "room_name", value: "General's Quarters"},
	}, http.StatusSeeOther},
	{"invalid-data", []postData{
		{key: "start_date", value: "2050-01-01"},
		{key: "end_date", value: "2050-01-02"},
		{key: "first_name", value: "a"},
		{key: "last_name", value: "s"},
		{key: "email", value: "john@smith.com"},
		{key: "phone", value: "123123456"},
		{key: "room_id", value: "1"},
		{key: "room_name", value: "General's Quarters"},
	}, http.StatusOK},
	{"failure-to-insert-reservation", []postData{
		{key: "start_date", value: "2050-01-01"},
		{key: "end_date", value: "2050-01-02"},
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "john@smith.com"},
		{key: "phone", value: "123123456"},
		{key: "room_id", value: "2"},
		{key: "room_name", value: "General's Quarters"},
	}, http.StatusTemporaryRedirect},
	{"failure-to-insert-restriction", []postData{
		{key: "start_date", value: "2050-01-01"},
		{key: "end_date", value: "2050-01-02"},
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "john@smith.com"},
		{key: "phone", value: "123123456"},
		{key: "room_id", value: "1000"},
		{key: "room_name", value: "General's Quarters"},
	}, http.StatusTemporaryRedirect},
}

var postAvailabilityJSONTests = []struct {
	name   string
	params []postData
}{
	{"happy-path", []postData{
		{key: "start_date", value: "2050-01-01"},
		{key: "end_date", value: "2050-01-02"},
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "john@smith.com"},
		{key: "phone", value: "123123456"},
		{key: "room_id", value: "1"},
		{key: "room_name", value: "General's Quarters"},
	}},
	{"invalid-post-data", nil},
	{"happy-path", []postData{
		{key: "start_date", value: "2050-01-01"},
		{key: "end_date", value: "2050-01-02"},
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "john@smith.com"},
		{key: "phone", value: "123123456"},
		{key: "room_id", value: "1000"},
		{key: "room_name", value: "General's Quarters"},
	}},
	{"invalid-post-data", nil},
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
				t.Errorf("for %s expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepositoryReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's quarters",
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
		t.Errorf("Reservation handler returned wrong response code. Got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code. Got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code. Got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {

	for _, e := range postReservationTests {
		var body io.Reader
		postData := url.Values{}

		for _, v := range e.params {
			postData.Add(v.key, v.value)
		}

		if len(postData) == 0 {
			body = nil
		} else {
			body = strings.NewReader(postData.Encode())
		}

		req, _ := http.NewRequest("POST", "/make-reservation", body)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.PostReservation)
		handler.ServeHTTP(rr, req)
		if rr.Code != e.expectedStatusCode {
			t.Errorf("reservation handler return wrong response code for %s. Expected %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	for _, e := range postAvailabilityJSONTests {
		var body io.Reader
		postData := url.Values{}

		for _, v := range e.params {
			postData.Add(v.key, v.value)
		}

		if len(postData) == 0 {
			body = nil
		} else {
			body = strings.NewReader(postData.Encode())
		}

		req, _ := http.NewRequest("POST", "/search-availability-json", body)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AvailabilityJSON)
		handler.ServeHTTP(rr, req)

		var j jsonResponse
		err := json.Unmarshal([]byte(rr.Body.String()), &j)

		if err != nil {
			t.Error("failed to parse json")
		}
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
