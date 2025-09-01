package service

import (
	"context"
	"github.com/star-find-cloud/star-mall/conf"
	"github.com/star-find-cloud/star-mall/pkg/mail"
	"github.com/star-find-cloud/star-mall/repo"
	"gopkg.in/gomail.v2"
)

type PublicService interface {
	// SendVerificationCode 发送验证码
	SendVerificationCode(ctx context.Context, email string) (string, error)
}

type PublicServiceImp struct {
	repo *repo.PublicRepo
}

func NewPublicService(repo *repo.PublicRepo) *PublicServiceImp {
	return &PublicServiceImp{repo: repo}
}

func (s *PublicServiceImp) SendVerificationCode(ctx context.Context, email string) (string, error) {
	c := conf.GetConfig()
	code := mail.GenerateCode()

	m := gomail.NewMessage()
	m.SetHeader("From", c.Mail.From)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "A verification code e-mail from Star Mall")
	m.SetBody("text/html", "【StarMall】您正在访问寻星商场,验证码是<b>"+code+"</b>，有效期5分钟（为避免被欺诈，如非本人主动操作，请勿继续输入验证码，也不要把验证码告知他人）")

	if err := s.repo.SetVerificationCodeCache(ctx, email, code); err != nil {
		return "", err
	}

	//fmt.Println("code: ", code)
	var d *gomail.Dialer
	d = gomail.NewDialer(c.Mail.Host, c.Mail.Port, c.Mail.User, c.Mail.SMTPCode)
	return code, d.DialAndSend(m)
}
