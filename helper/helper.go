package helper

import (
	"clean_arch_v2/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type helper struct {
	jwtSecret string
}

func NewHelper() models.Helper {
	secret := viper.GetString("jwt.key")

	return &helper{
		jwtSecret: secret,
	}
}


func (h *helper) CreateToken(s string) string {
	setToken := jwt.New(jwt.SigningMethodHS512)
	claims := setToken.Claims.(jwt.MapClaims)

	// depend on data on token use model
	claims["id"] = s
	token, err := setToken.SignedString([]byte(h.jwtSecret))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return token
} 

func (h *helper) ErrorLog(s string) {
	fmt.Println(s)
	// var mw io.writer
	// fmt.Fprintf(mw, `
	// level=info
	// datetime="%s" 
	// ip=%s method=%s 
	// url="%s" 
	// proto=%s 
	// status=%d 
	// latency=%s 
	// user="%s" 
	// device="%s" 
	// req_header:"%s" 
	// req_body="%s" 
	// response=%s`,
	// t.Format(time.RFC1123), c.ClientIP(), c.Request.Method, c.Request.URL.String(),
	// c.Request.Proto, c.Writer.Status(), latency, user,
	// device, reqHeaderStr, reqBodyStr, res,
	// )
}

func(g *helper) PasswordHash(s string) (string, error) {
	passbyte, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		
		return "", err
	}

	return string(passbyte), nil
}

func(g *helper) PasswordCompare(pass string, compare string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(compare))
	fmt.Println(err)
	return
}

func (g *helper) TokenDecrypt(c *gin.Context, response *models.MainResponse) *models.MainResponse {
	
	bearer := c.GetHeader("Authorization")
	strSplit := strings.Split(bearer, " ")
	if len(strSplit) < 2 {
		response.Status = false
		response.StatusCode = http.StatusUnauthorized
		response.Message = "token is required"
		return response
	}
	// if token doesnt't have 3 parts it is unhanled error on token claims
	tokenSection := 0
	for _, t := range strSplit[1] {
		if string(t) == "." {
			tokenSection++
		}
	}
	if bearer != "" && len(strSplit) == 2 && tokenSection == 2{
		token, err := jwt.Parse(strSplit[1], func(token *jwt.Token) (interface{}, error) {
			if jwt.GetSigningMethod("HS512") != token.Method {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(g.jwtSecret), nil
		})
		claims,ok := token.Claims.(jwt.MapClaims)
		if err == nil && ok && token.Valid {
			response.Payload.Id = claims["id"].(string)
			return response
		} else {
			response.Status = false
			response.StatusCode = http.StatusUnauthorized
			response.Message = "token not valid"
			return response	
		}
	} else {
		response.Status = false
		response.StatusCode = http.StatusUnauthorized
		response.Message = "token not valid"
		return response
	}
}