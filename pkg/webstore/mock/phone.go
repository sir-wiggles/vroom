package mock

type PhoneService struct {
	SendCodeFn      func(string) error
	SendCodeInvoked bool

	VerifyCodeFn      func(string, int) error
	VerifyCodeInvoked bool
}

func (p *PhoneService) SendCode(phoneNumber string) error {
	p.SendCodeInvoked = true
	return p.SendCodeFn(phoneNumber)
}

func (p *PhoneService) VerifyCode(phoneNumber string, code int) error {
	p.VerifyCodeInvoked = true
	return p.VerifyCodeFn(phoneNumber, code)
}
