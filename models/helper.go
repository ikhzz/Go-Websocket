package models

import "github.com/gin-gonic/gin"

type Helper interface {
	ErrorLog(s string)
	CreateToken(s string) string
	PasswordHash(s string) (string, error)
	PasswordCompare(pass string, compare string) error
	TokenDecrypt(c *gin.Context, response *MainResponse) *MainResponse
}

type (
	TokenModel struct {
		Token string `json:"token"`
		// id have to be inside token, testing purposes
		Id int `json:"id"`
	}
)