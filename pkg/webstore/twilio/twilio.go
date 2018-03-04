package twilio

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"os"

	"github.com/sir-wiggles/arc/pkg/webstore"
)

type Service struct {
	Cache      webstore.CacheService
	HTTPClient *http.Client
}

func NewService() *Service {
	return &Service{
		HTTPClient: &http.Client{},
	}
}

func generateCode() (string, error) {
	var (
		r1  *big.Int
		r2  *big.Int
		err error
	)

	// TODO: should match the method that zapi uses
	if r1, err = rand.Int(rand.Reader, big.NewInt(129)); err != nil {
		return "", err
	}

	if r2, err = rand.Int(rand.Reader, big.NewInt(129)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%03d%03d", r1, r2), nil
}

//SendCode sends a verification code the phone number provided and stores
//the code in cache with the phone number as the key to be looked up during
//validation.
func (s *Service) SendCode(phoneNumber string) error {
	var (
		code    string
		err     error
		form    map[string][]string
		message string
		req     *http.Request
		res     *http.Response
	)

	if code, err = generateCode(); err != nil {
		return err
	}

	// TODO: get message from zapi
	message = url.QueryEscape(fmt.Sprintf("Zume verification code: %s", code))

	// TODO: handle env vars
	form = map[string][]string{
		"Body": []string{message},
		"From": []string{os.Getenv("TWILIO_FROM_PHONE_NUMBER")},
		"To":   []string{phoneNumber},
	}

	req, err = http.NewRequest("POST", os.Getenv("TWILIO_URL"), nil)
	req.SetBasicAuth(os.Getenv("TWILIO_ACCOUNT"), os.Getenv("TWILIO_AUTH_TOKEN"))
	req.Form = form

	if res, err = s.HTTPClient.Do(req); err != nil {
		return err
	}

	// TODO: handle stupid response from twilio 200 with errors
	if res.StatusCode != http.StatusOK {
		return errors.New("Non 200 from Twilio")
	}

	if err = s.Cache.Set(phoneNumber, code); err != nil {
		return err
	}

	return nil
}

func (s *Service) VerifyCode(code int) error {
	return nil
}
