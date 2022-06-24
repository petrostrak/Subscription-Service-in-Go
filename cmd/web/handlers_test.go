package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/petrostrak/Subscription-Service-in-Go/data"
)

var pageTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
	handler            http.HandlerFunc
	sessionData        map[string]any
	expectedHTML       string
}{
	{
		name:               "home",
		url:                "/",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
		handler:            testApp.HomePage,
	},
	{
		name:               "login",
		url:                "/login",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
		handler:            testApp.LoginPage,
	},
	{
		name:               "post login",
		url:                "/login",
		method:             "POST",
		expectedStatusCode: http.StatusSeeOther,
		handler:            testApp.PostLoginPage,
	},
	{
		name:               "logout",
		url:                "/logout",
		method:             "GET",
		expectedStatusCode: http.StatusSeeOther,
		handler:            testApp.Logout,
		sessionData: map[string]any{
			"userID": 1,
			"user":   data.User{},
		},
	},
	{
		name:               "plans (logged in)",
		url:                "/plans",
		method:             "GET",
		expectedStatusCode: http.StatusOK,
		handler:            testApp.ChooseSubscription,
		sessionData: map[string]any{
			"userID": 1,
			"user":   data.User{},
		},
		expectedHTML: `<h1 class="mt-5">Plans</h1>`,
	},
	{
		name:               "plans (not logged in)",
		url:                "/plans",
		method:             "GET",
		expectedStatusCode: http.StatusTemporaryRedirect,
		handler:            testApp.ChooseSubscription,
	},
	{
		name:               "subscribe (not logged in)",
		url:                "/subscribe",
		method:             "GET",
		expectedStatusCode: http.StatusTemporaryRedirect,
		handler:            testApp.SubscribeToPlan,
	},
}

func Test_Pages(t *testing.T) {
	pathToTemplates = "./templates"

	for _, e := range pageTests {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", e.url, nil)

		ctx := getCtx(req)
		req = req.WithContext(ctx)

		if len(e.sessionData) > 0 {
			for key, value := range e.sessionData {
				testApp.Session.Put(ctx, key, value)
			}
		}

		e.handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("%s failed: expected %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if len(e.expectedHTML) > 0 {
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("%s failed: expected to find %s, but did not", e.name, e.expectedHTML)
			}
		}
	}

}

func TestConfig_PostLoginPage(t *testing.T) {
	pathToTemplates = "./templates"

	postedData := url.Values{
		"email":    {"admin@example.com"},
		"password": {"abc123abc123abc123abc123"},
	}

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	handler := http.HandlerFunc(testApp.PostLoginPage)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Error("wrong code returned")
	}

	if !testApp.Session.Exists(ctx, "userID") {
		t.Error("did not find userID in session")
	}
}
