package api_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sir-wiggles/arc/pkg/webstore/api"
	"github.com/sir-wiggles/arc/pkg/webstore/mock"
)

func TestPhoneSMSCode(t *testing.T) {
	url := "/phone/sms/code"
	testCases := []struct {
		method     string
		url        string
		body       io.Reader
		statusCode int
		psInvoked  bool
	}{
		{"POST", url, bytes.NewReader([]byte(`{"phone": 5554443333}`)), 200, true},
		{"DELETE", url, nil, 405, false},
		{"GET", url, nil, 405, false},
		{"HEAD", url, nil, 405, false},
		{"PATCH", url, nil, 405, false},
		{"PUT", url, nil, 405, false},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("%s#%s", tc.method, tc.url)
		t.Run(name, func(t *testing.T) {
			h := api.NewHandler()
			mps := mock.PhoneService{}

			mps.SendCodeFn = func(string) error {
				return nil
			}

			mcs := mock.CacheService{}

			h.PhoneService = &mps
			h.CacheService = &mcs

			w := httptest.NewRecorder()
			r, err := http.NewRequest(tc.method, tc.url, tc.body)
			if err != nil {
				t.Error(err)
			}

			h.ServeHTTP(w, r)

			rStatusCode := w.Result().StatusCode
			if rStatusCode != tc.statusCode {
				t.Errorf("Expected %d got %d", tc.statusCode, rStatusCode)
			}

			invoked := mps.SendCodeInvoked
			if invoked != tc.psInvoked {
				t.Errorf("Expected %t got %t", tc.psInvoked, invoked)
			}
		})
	}
}

func TestHandler_PhoneVerifyCode(t *testing.T) {
	url := "/phone/sms/verify"
	testCases := []struct {
		method     string
		url        string
		body       io.Reader
		statusCode int
		psInvoked  bool
	}{
		{"POST", url, nil, 200, true},
		{"DELETE", url, nil, 405, false},
		{"GET", url, nil, 405, false},
		{"HEAD", url, nil, 405, false},
		{"PATCH", url, nil, 405, false},
		{"PUT", url, nil, 405, false},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("%s#%s %v", tc.method, tc.url, tc.body)
		t.Run(name, func(t *testing.T) {
			h := api.NewHandler()
			mps := mock.PhoneService{}

			mps.VerifyCodeFn = func(pn string, id int) error {
				return nil
			}

			h.PhoneService = &mps

			w := httptest.NewRecorder()
			r, err := http.NewRequest(tc.method, tc.url, tc.body)
			if err != nil {
				t.Error(err)
			}

			h.ServeHTTP(w, r)

			rStatusCode := w.Result().StatusCode
			if rStatusCode != tc.statusCode {
				t.Errorf("Expected %d got %d", tc.statusCode, rStatusCode)
			}

			invoked := mps.VerifyCodeInvoked
			if invoked != tc.psInvoked {
				t.Errorf("Expected %t got %t", tc.psInvoked, invoked)
			}
		})
	}
}
