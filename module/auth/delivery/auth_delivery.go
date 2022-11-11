package delivery

import (
	"clean_arch_v2/helper"
	"clean_arch_v2/models"
	"clean_arch_v2/module/auth/repository"
	"clean_arch_v2/module/auth/usecase"

	"github.com/gin-gonic/gin"
)

type AuthDelivery struct {
	Usecase models.AuthUsecase
	Response models.MainResponse
	Helper 		models.Helper
}

func NewAuthDelivery(r *gin.Engine, response models.MainResponse) {
	helpers := helper.NewHelper()
	repo := repository.NewAuthRepository()
	usecase := usecase.NewAuthUsecase(repo, helpers)
	
	
	handler := &AuthDelivery{
		Usecase: usecase,
		Response: response,
		Helper: helpers,
	}

	v1 := r.Group("/v1")
	{
		v1.POST("/signup", handler.Signup)
		v1.POST("/signin", handler.Signin)
		v1.GET("/getUser", handler.GetAllUser)
	}
}

func(a AuthDelivery) Signup(c *gin.Context) {
	var reqParam models.AuthModel
	c.Bind(&reqParam)

	a.Usecase.SignUp(&a.Response, &reqParam)
	c.JSON(a.Response.StatusCode, a.Response)
}

func(a AuthDelivery) Signin(c *gin.Context) {
	var reqParam models.ReqSigninModel
	c.Bind(&reqParam)

	a.Usecase.SignIn(&a.Response, &reqParam)
	c.JSON(a.Response.StatusCode, a.Response)
}

func(a AuthDelivery) GetAllUser(c *gin.Context) {
	a.Helper.TokenDecrypt(c, &a.Response)
	
	if a.Response.Payload.Id != "0" {
		a.Usecase.GetAllUser(&a.Response)	
	}

	a.Usecase.GetAllUser(&a.Response)
	c.JSON(a.Response.StatusCode, a.Response)
	
}