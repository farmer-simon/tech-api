package my_jwt

import "github.com/dgrijalva/jwt-go"

// 自定义jwt的声明字段信息+标准字段，参考地址：https://blog.csdn.net/codeSquare/article/details/99288718
type CustomClaims struct {
	Platform string `json:"platform"`
	UserId   int64  `json:"user_id"`
	Params   string `json:"params"`
	jwt.StandardClaims
}
