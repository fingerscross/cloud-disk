package help

import (
	"bytes"
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
	"io"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GenerateToken(id int, identity, name string, second int) (string, error) {
	uc := define.UserClaim{
		Id:             id,
		Identity:       identity,
		Name:           name,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Second * time.Duration(second)).Unix()},
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

func CosInitPart(ext string) (string, string, error) {
	u, _ := url.Parse(define.CosBUcket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
			SecretKey: define.TencentSecretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
		},
	})
	key := "cloud-disk/" + UUID() + ext //对应存储桶路径
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), key, nil)
	if err != nil {
		return "", "", err
	}
	UploadID := v.UploadID
	fmt.Println(UploadID)
	return key, v.UploadID, nil
}

// 分片上传
func CosPartUpload(r *http.Request) (string, error) {
	u, _ := url.Parse(define.CosBUcket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: define.TencentSecretID, // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: define.TencentSecretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
		},
	})

	key := r.PostForm.Get("key")
	UploadID := r.PostForm.Get("upload_id")
	part_number, err := strconv.Atoi(r.PostForm.Get("part_number"))

	if err != nil {
		return "", err
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}

	//得转化一下 直接用reader不行
	buffer := bytes.NewBuffer(nil)
	io.Copy(buffer, file)

	// opt 可选
	resp, err := client.Object.UploadPart(
		context.Background(), key, UploadID, part_number, bytes.NewReader(buffer.Bytes()), nil,
	)
	if err != nil {
		return "", err
	}
	PartETag := resp.Header.Get("ETag") //MD5值 ： 0.chunk
	fmt.Println(PartETag)
	return strings.Trim(resp.Header.Get("ETag"), "\""), nil
}

// 分片上传完成
func CosPartUploadComplete(key, uploadId string, cs []cos.Object) error {
	u, _ := url.Parse(define.CosBUcket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts, cs...)
	_, _, err := client.Object.CompleteMultipartUpload(
		context.Background(), key, uploadId, opt,
	)
	return err

}
