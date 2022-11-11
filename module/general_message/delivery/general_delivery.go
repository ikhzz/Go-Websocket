package delivery

import (
	"clean_arch_v2/helper"
	"clean_arch_v2/models"
	"clean_arch_v2/module/general_message/repository"
	"clean_arch_v2/module/general_message/usecase"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type GeneralDelivery struct {
	Usecase 	models.GeneralMUsecase
	Response 	models.MainResponse
	Helper 		models.Helper
}

func NewGeneralMessageDelivery(r *gin.Engine, response models.MainResponse) {
	helpers := helper.NewHelper()
	repo := repository.NewGeneralMessageRepository()
	usecase :=  usecase.NewGeneralMessageUsecase(repo, helpers)

	handler := &GeneralDelivery{
		Usecase: usecase,
		Response: response,
		Helper: helpers,
	}

	v1 := r.Group("/v1")
	{
		v1.GET("/initConn", handler.InitConn)
	}
}

func(a GeneralDelivery) InitConn(c *gin.Context) {
	openConn, err := websocket.Upgrade(c.Writer, c.Request, c.Writer.Header(), 1024, 1024)
	if err != nil {
		a.Response.Message = err.Error()
		a.Response.Status = false
		c.JSON(a.Response.StatusCode, a.Response)
		return
	}
	queryToken := c.Request.URL.Query().Get("token")
	c.Request.Header.Add("Authorization", "Bearer "+queryToken)
	a.Helper.TokenDecrypt(c, &a.Response)
	
	id, err := strconv.Atoi(a.Response.Payload.Id)
	fmt.Println("ini id",id)
	if err == nil {
		currentConn := models.WebSocketConnection{Conn: openConn, Id: id}
  	helper.Conn = append(helper.Conn, currentConn)
		fmt.Println("current conn", helper.Conn)
		go a.Usecase.HandleInit(&currentConn, &a.Response)
	}
}
