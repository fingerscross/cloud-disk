package help

import (
	"cloud-disk/core/define"
	"context"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jordan-wright/email"
	_ "github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"path"
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

// 解析token
func AnakyzeToken(token string) (*define.UserClaim, error) {
	uc := new(define.UserClaim)
	claims, err := jwt.ParseWithClaims(token, uc, func(token *jwt.Token) (interface{}, error) {
		return []byte(define.Jwtkey), nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return uc, errors.New("token is invaild")
	}
	return uc, err
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

func CosUpload(r *http.Request) (string, error) {
	u, _ := url.Parse(define.CosBUcket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
			SecretKey: define.TencentSecretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
		},
	})
	file, fileheader, err := r.FormFile("file")
	//腾讯云中的路径与结果名字
	key := "cloud-disk/" + UUID() + path.Ext(fileheader.Filename) //唯一id+后缀

	_, err = client.Object.Put(
		context.Background(), key, file, nil,
	)
	if err != nil {
		panic(err)
	}
	return define.CosBUcket + "/" + key, nil //以访问链接的方式访问文件
}
