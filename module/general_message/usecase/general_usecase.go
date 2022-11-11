package usecase

import (
	"clean_arch_v2/helper"
	"clean_arch_v2/models"
	"clean_arch_v2/module/auth/repository"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type generalUsecase struct {
	repo 	models.GeneralMRepository
	help	models.Helper
	auth 	models.AuthRepository
}

func NewGeneralMessageUsecase(a models.GeneralMRepository, help models.Helper) models.GeneralMUsecase {
	auth := repository.NewAuthRepository()
	return &generalUsecase{
		repo: a,
		help: help,
		auth: auth,
	}
}

func(gm *generalUsecase) HandleInit(currentConn *models.WebSocketConnection, response *models.MainResponse) {		
	defer func() {
		if r := recover(); r != nil {
				fmt.Println("ERROR", fmt.Sprintf("%v", r))
		}
	}()

	// get history room and chat
	func() {
		fmt.Println("get history room and chat")
		gm.repo.UpdateUserActiveStatus(currentConn.Id, 1)
		listUser := gm.auth.GetAllUser(response)
		historyChat := make([]models.ChatHistory, 0)
		for _, eachUser := range listUser {
			var history models.ChatHistory

			historyData := gm.repo.GetUserHistory(currentConn.Id, eachUser.Id)
			history.Id = eachUser.Id
			history.Chat = *historyData
			historyChat = append(historyChat, history)
		}

		responseModel := models.WebsocketResponseModel{ 
			Data: historyChat,
			Type: "HISTORY",
		}
		
		err := currentConn.WriteJSON(responseModel)
		fmt.Println(err)
	}()

	// loop to get and send message
	for {
		payload := models.WebsocketPayloadModel{}
		err := currentConn.ReadJSON(&payload)
		fmt.Println(payload)
		if err != nil {
			fmt.Println("error", err)
				if strings.Contains(err.Error(), "websocket: close") {
					// set room status to offline
					func() {
						fmt.Println("set room status to offline")
						gm.repo.UpdateUserActiveStatus(currentConn.Id, 0)		
					}()

					// remove offline room from connection
					func() {
						fmt.Println("remove offline line room from connection")
						var tempConn []models.WebSocketConnection
						for _, connValue := range helper.Conn {

							if connValue.Id != currentConn.Id {
									tempConn = append(tempConn, connValue)
							}
						}
						helper.Conn = tempConn
					}()
					return
				}

				fmt.Println("ERROR", err.Error())	
				continue
		}

		// send message
		func() {
			fmt.Println("send chat to other client")
			resChat := models.GeneralMessageListModel{
				IdFrom: currentConn.Id, IdTo: payload.Id, Message: payload.Message, CreatedAt: time.Now(), SendBy: currentConn.Id,
			}
			errc:= gm.repo.AddChat(&resChat)
			fmt.Println("send chat now", errc)

			for _, eachConn := range helper.Conn {
				if eachConn.Id != payload.Id {
					continue
				}
					
				sendData := models.ChatModel{SendBy: currentConn.Id, Message: payload.Message, CreatedAt: time.Now()}
				err := eachConn.WriteJSON(gin.H{
						"Type": "NEW_MESSAGE",
						"Data": sendData,
				})
				fmt.Println(err)
			}
		}()
	}
}