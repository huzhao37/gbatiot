package middleware

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"time"
	"net/http"
	"log"
)
// 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token :=  c.Request.Header.Get("token")
		if token == ""{
			c.JSON(http.StatusOK,gin.H{
				"status":-1,
				"msg":"请求未携带token，无权限访问",
			})
			c.Set("isPass", false)
			c.Abort()
			return
		}

		log.Print("get token: ",token)

		j := NewJWT()
		// parseToken
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusOK,gin.H{
					"status":-1,
					"msg":"授权已过期",
				})
				c.Set("isPass", false)
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg": err.Error(),
			})
			c.Set("isPass", false)
			c.Abort()
			return
		}
		c.Set("isPass", true)
		c.Set("claims",claims)
	}
}
// 签名
type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired error = errors.New("Token is expired")
	TokenNotValidYet error = errors.New("Token not active yet")
	TokenMalformed error = errors.New("That's not even a token")
	TokenInvalid error = errors.New("Couldn't handle this token:")
	SignKey string = "newtrekWang"
)
// 载荷
type CustomClaims struct {
	ID int `json:"userid"`
	RoleId int `json:"userroleid"`
	Account string `json:"account"`
	Phone string `json:"mobileno"`
	jwt.StandardClaims
}
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}
func GetSignKey() string {
	return SignKey
}
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}


func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(48 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
