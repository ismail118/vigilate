package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var loginTest = []struct {
	name string
	url string
	method string
	postedData url.Values
	expectedResponseCode int
} {
	{
		name: "login-screen",
		url: "/",
		method: "GET",
		expectedResponseCode: http.StatusOK,
	},
	{
		name: "login-screen-post",
		url: "/",
		method: "POST",
		postedData: url.Values{
			"email": {"me@here.com"},
			"password": {"secret"},
		},
		expectedResponseCode: http.StatusSeeOther,
	},
}

func TestLoginScreen(t *testing.T) {
	for _, e := range loginTest {
		if e.method == "GET" {
			req, _ := http.NewRequest(e.method, e.url, nil)

			ctx := getCtx(req)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(Repo.LoginScreen)

			handler.ServeHTTP(rr, req)

			if rr.Code != e.expectedResponseCode {
				t.Errorf("%s, expected %d got %d", e.name, e.expectedResponseCode, rr.Code)
			}
		} else {
			req, _ := http.NewRequest(e.method, e.url, strings.NewReader(e.postedData.Encode()))

			ctx := getCtx(req)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(Repo.Login)

			handler.ServeHTTP(rr, req)

			if rr.Code != e.expectedResponseCode {
				t.Errorf("%s, expected %d got %d", e.name, e.expectedResponseCode, rr.Code)
			}
		}
		
	}
}

func TestDBRepo_PusherAuth(t *testing.T) {
	e := struct {
		name string
		url string
		method string
		postedData url.Values
		expectedResponseCode int
	} {
		name: "auth",
		url: "/pusher/auth",
		method: "POST",
		postedData: url.Values{
			"socket_id": {"671772991.50991907"},
			"channel_name": {"private-channel-1"},
		},
		expectedResponseCode: http.StatusOK,
	}

	req, _ := http.NewRequest(e.method, e.url, strings.NewReader(e.postedData.Encode()))

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PusherAuth)

	handler.ServeHTTP(rr, req)

	if rr.Code != e.expectedResponseCode {
		t.Errorf("%s, expected %d got %d", e.name, e.expectedResponseCode, rr.Code)
	}

	var puserResp struct {
		Auth string `json:"auth"`
	}

	err := json.NewDecoder(rr.Body).Decode(&puserResp)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("got: ", puserResp.Auth)

	if len(puserResp.Auth) == 0 {
		t.Error("empty json response")
	}
}