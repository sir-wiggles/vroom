package twilio_test

import (
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sir-wiggles/arc/pkg/webstore/mock"
	"github.com/sir-wiggles/arc/pkg/webstore/twilio"
)

type fakeRoundTripper struct {
	response *http.Response
	err      error
}

func newRoundTripper(r *http.Response, err error) *fakeRoundTripper {
	return &fakeRoundTripper{r, err}
}

func (rt *fakeRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt.response, rt.err
}

func TestSendCode(t *testing.T) {

	type CalledWith struct {
		Key   string
		Value string
	}

	var calledWith *CalledWith

	testCases := []struct {
		_name    string
		SetFn    func(string, string) error
		hasError bool
		invoked  bool
	}{
		{
			"Test",
			func(k, v string) error {
				calledWith.Key = k
				calledWith.Value = v
				return nil
			},
			false,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc._name, func(t *testing.T) {

			g := NewGomegaWithT(t)

			calledWith = &CalledWith{}
			phone := twilio.NewService()

			res := &http.Response{
				StatusCode: http.StatusOK,
			}
			phone.HTTPClient.Transport = newRoundTripper(res, nil)

			mc := &mock.CacheService{}
			mc.SetFn = tc.SetFn
			phone.Cache = mc
			err := phone.SendCode("5599361530")

			if tc.hasError {
				g.Expect(err).ToNot(BeNil(), "return error")
			} else {
				g.Expect(err).To(BeNil(), "return error")
			}

			g.Expect(mc.SetInvoked).To(Equal(tc.invoked), "invoked")
		})
	}
}
