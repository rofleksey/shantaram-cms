package service

import (
	"github.com/mojocn/base64Captcha"
)

type Captcha struct{}

func NewCaptcha() *Captcha {
	return &Captcha{}
}

func (s *Captcha) Generate() (string, string, error) {
	driver := base64Captcha.NewDriverDigit(100, 240, 4, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	id, b64s, _, err := captcha.Generate()

	if err != nil {
		return "", "", err
	}

	return id, b64s, nil
}

func (s *Captcha) Verify(id, answer string) bool {
	return base64Captcha.DefaultMemStore.Verify(id, answer, true)
}
