package usecase

import (
	"clean_arch_v2/models"
	"clean_arch_v2/module/general_message/repository"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
)

type authUsecase struct {
	repo	models.AuthRepository
	helper 	models.Helper
	general models.GeneralMRepository
}

func NewAuthUsecase(a models.AuthRepository, help models.Helper) models.AuthUsecase {
	gm := repository.NewGeneralMessageRepository()
	return &authUsecase{
		repo: a,
		helper: help,
		general: gm,
	}
}

func(a authUsecase) SignUp(response *models.MainResponse, reqParam *models.AuthModel) *models.MainResponse {
	_, errParseEmail := mail.ParseAddress(reqParam.Email)

	if errParseEmail != nil {
		a.helper.ErrorLog(errParseEmail.Error())
		response.Status = false
		response.StatusCode = http.StatusBadRequest
		response.Message = "email is not valid"
		return response
	}

	paramError :=  make([]string, 0)
	switch {
		case reqParam.Username == "":
			paramError = append(paramError, "username is required")
			fallthrough
		case reqParam.Password == "":
			paramError = append(paramError, "password is required")
			fallthrough
		case len(reqParam.Username) <= 5:
			paramError = append(paramError, "username require more than 5 char")
			fallthrough
		case len(reqParam.Password) <= 5:
			paramError = append(paramError, "password require more than 5 char")
	}

	if len(paramError) > 0 {
		a.helper.ErrorLog(strings.Join(paramError, ", "))
		response.Status = false
		response.StatusCode = http.StatusBadRequest
		response.Message = strings.Join(paramError, ", ")
		return response
	}
	createHash, errHash := a.helper.PasswordHash(reqParam.Password)
	if errHash != nil {
		a.helper.ErrorLog(errParseEmail.Error())
		response.Status = false
		response.StatusCode = http.StatusBadRequest
		response.Message = "hash password failed"
		return response
	}

	reqParam.Password = createHash
	errSignup := a.repo.SignUp(reqParam)
	if errSignup != nil {
		a.helper.ErrorLog(errSignup.Error())
		response.Status = false
		response.StatusCode = http.StatusBadRequest
		response.Message = "sign up failed email already registered"
		return response
	}
	
	token := a.helper.CreateToken(strconv.Itoa(reqParam.Id))
	tokenReponse := models.TokenModel{Token: token}
	
	response.Data = tokenReponse
	response.Message = "sign up success"
	
	return response
}

func(a authUsecase) SignIn(response *models.MainResponse, reqParam *models.ReqSigninModel) *models.MainResponse {
	_, errParseEmail := mail.ParseAddress(reqParam.Email)

	if errParseEmail != nil {
		a.helper.ErrorLog(errParseEmail.Error())
		response.Status = false
		response.StatusCode = http.StatusBadRequest
		response.Message = "email is not valid"
		return response
	}

	if reqParam.Password == "" && len(reqParam.Password) <= 5 {
		response.Status = false
		response.StatusCode = http.StatusBadRequest
		response.Message = "password is required and more than 5 char"
		return response
	}
	
	res, err := a.repo.SignIn(reqParam)
	if err != nil {
		response.Status = false
		response.StatusCode = http.StatusBadRequest
		response.Message = "failed to sign in"
		return response
	}

	errCompare := a.helper.PasswordCompare(res.Password, reqParam.Password)
	if errCompare != nil {
		response.Status = false
		response.StatusCode = http.StatusBadRequest
		response.Message = "password doesnt match"
		return response
	}

	token := a.helper.CreateToken(strconv.Itoa(res.Id))
	tokenReponse := models.TokenModel{Token: token, Id: res.Id}
	
	response.Data = tokenReponse
	response.Message = "sign in success"

	return response
}

func(gm *authUsecase) GetAllUser(response *models.MainResponse) *models.MainResponse {		
	result := gm.repo.GetAllUser(response)
	for i, v := range result {
		idInt, _ := strconv.Atoi(response.Payload.Id)
		res := gm.general.GetUserHistory(idInt, v.Id)
		if len(*res) > 0{
			result[i].IsHistory = 1
		}
	}
	
	response.Data = result
	response.Message = "success get all user"
	if len(result) == 0 {
		response.Message = "success with no user found"
	}

	return response
}