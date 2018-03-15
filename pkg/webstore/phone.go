package webstore

type PhoneService interface {
	SendCode(string) error
	VerifyCode(string, int) error
}
