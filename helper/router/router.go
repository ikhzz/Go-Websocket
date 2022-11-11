package router

import (
	"clean_arch_v2/models"
	auth "clean_arch_v2/module/auth/delivery"
	general_message "clean_arch_v2/module/general_message/delivery"

	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	conn []models.WebSocketConnection 
)

func InitRouter(c *gin.Engine, mw io.Writer) {
	
	prepRes := models.MainResponse{
		StatusCode: http.StatusOK,
		Status: true,
		// default 
		Message: "Success access request",
		Data: make([]int, 0),
	}

	auth.NewAuthDelivery(c, prepRes)
	general_message.NewGeneralMessageDelivery(c, prepRes)
}