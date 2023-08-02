package define

import (
	"github.com/golang-jwt/jwt/v4"
	"os"
)

type UserClaim struct {
	Id       int
	Identity string
	Name     string
	jwt.StandardClaims
}

var Jwtkey = "cloud-disk-key"
var MailPassword = "NDHBNEVUQPDIWUPH"

// 验证码长度
var CodeLength = 6

// 验证码过期时间
var CodeExpire = 300

// 腾讯云对象存储
var TencentSecretKey = os.Getenv("TencentSecterKey")
var TencentSecretID = os.Getenv("TencentSecterID")
