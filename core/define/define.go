package define

import (
	"github.com/golang-jwt/jwt/v4"
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
