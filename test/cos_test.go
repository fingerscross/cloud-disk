package test

import (
	"bytes"
	"cloud-disk/core/define"
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestGetId(t *testing.T) {
	_, err := fmt.Println(define.TencentSecretID)
	fmt.Println(define.TencentSecretKey)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFileUploadByFilepath(t *testing.T) {
	u, _ := url.Parse("https://go-cloud-disk-1310218422.cos.ap-nanjing.myqcloud.com")
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

	//腾讯云中的路径与结果名字
	key := "cloud-disk/笔试题2.docx"

	_, _, err := client.Object.Upload(
		context.Background(), key, "./file/笔试题-【后端实习生版】.docx", nil,
	)
	if err != nil {
		panic(err)
	}

}

func TestFileUploadByReader(t *testing.T) {
	u, _ := url.Parse("https://go-cloud-disk-1310218422.cos.ap-nanjing.myqcloud.com")
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

	//腾讯云中的路径与结果名字
	key := "cloud-disk/笔试题3.docx"
	file, err := os.ReadFile("./file/笔试题-【后端实习生版】.docx")
	if err != nil {
		return
	}
	_, err = client.Object.Put(
		context.Background(), key, bytes.NewReader(file), nil,
	)
	if err != nil {
		panic(err)
	}
}
