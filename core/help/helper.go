package help

import (
	"cloud-disk/core/define"
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jordan-wright/email"
	_ "github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"net/smtp"
	"time"
)

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GenerateToken(id int, identity, name string) (string, error) {
	uc := define.UserClaim{
		Id:       id,
		Identity: identity,
		Name:     name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString([]byte(define.Jwtkey)) //加密后的token
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func MailSendCode(mail, code string) error {
	e := email.NewEmail()
	e.From = "Jordan Wright <gu18915702696@163.com>"
	e.To = []string{mail}
	e.Subject = "验证码发送测试"

	e.HTML = []byte("<h1>" + code + "</h1>")
	err := e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "gu18915702696@163.com", define.MailPassword, "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		return err
	}
	return nil
}

func RandCode() string {
	s := "123456789"
	code := ""
	rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < define.CodeLength; i++ {
		code += string(s[rand.Intn(len(s))])
	}
	return code
}

func UUID() string {
	return uuid.NewV4().String()
}
