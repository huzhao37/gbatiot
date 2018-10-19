package apis

import (
	"net/http"
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	."yuniot/models/auth"
	"yuniot/core"
	"time"
	myjwt "yuniot/middleware"
	jwtgo "github.com/dgrijalva/jwt-go"
)


type LoginResult struct{
	Token string `json:"token"`
	User
}

func Test(c *gin.Context){
	c.JSON(http.StatusOK,gin.H{
		"message":"hello",
	})
}

func GetDataByTime(c *gin.Context) {
	isPass := c.GetBool("isPass")
	if !isPass {
		return
	}
	claims := c.MustGet("claims").(*myjwt.CustomClaims)
	if claims != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "token有效",
			"data":   claims,
		})
	}
}

// 登录
func Auth(c *gin.Context) {
	account := c.Request.PostFormValue("account")
	pwd:=c.Request.PostFormValue("pwd")

	if account!= ""&&pwd!=""{
		user, err := UserLogin(account,pwd)
		if err==nil&&user.Id>0 {
			generateToken(c,user)
		}else {
			c.JSON(http.StatusOK,gin.H{
				"status":-1,
				"msg":"验证失败"+err.Error(),
			})
			return
		}
	}else{
		c.JSON(http.StatusOK,gin.H{
			"status":-1,
			"msg":"json 解析失败",
		})
		return
	}


}
// 生成令牌
func generateToken(c *gin.Context, user User) {
	j := &myjwt.JWT{
		[]byte("newtrekWang"),
	}

	claims := myjwt.CustomClaims{
		user.Userid,
		user.Userroleid,
		user.Account,
		 user.Mobileno,
		jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),// 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600),// 过期时间 一小时
			Issuer: "newtrekWang",//签名的发行者
		},
	}
	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"status":-1,
			"msg":err.Error(),
		})
		return
	}
	fmt.Println(token)
	data := LoginResult{
		User:user,
		Token:token,
	}
	c.JSON(http.StatusOK,gin.H{
		"status":0,
		"msg":"登录成功！",
		"data":data,
	})
	return
}

func LoginApi(c *gin.Context) {
	isPass := c.GetBool("isPass")
	if !isPass {
		return
	}
	account := c.Request.PostFormValue("account")
	pwd:=c.Request.PostFormValue("pwd")
	user, err := UserLogin(account,pwd)
	if err != nil {
		core.Logger.Fatalln(err)
		c.JSON(http.StatusBadGateway, gin.H{
			"msg": "server exception",
		})
	}
	if user.Id==0{
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "login failed",
		})
	}else{
		c.JSON(http.StatusOK, gin.H{
			"msg": user,
		})
	}

}

func GetUsersApi(c *gin.Context){
	isPass := c.GetBool("isPass")
	if !isPass {
		return
	}
	var users = make([]User, 0)
	users, err :=GetUsers()
	if err != nil {
		core.Logger.Fatalln(err)
	}
	//msg := fmt.Sprintf("get successful %d", len(users))
	c.JSON(http.StatusOK, gin.H{
		"msg": users,
	})
}
