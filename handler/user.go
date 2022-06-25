package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	// tangkap input
	// map input ke registeruserinput
	// struct passing ke repository

	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMsg := helper.ErrorResponse(err)
		errow := gin.H{"errors": errorMsg}
		response := helper.APIResponse("Akun gagal regis", http.StatusUnprocessableEntity, "fail", errow)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Akun gagal regis", http.StatusUnprocessableEntity, "fail", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// update token
	token, err := h.authService.GenerateToken(newUser.Id)
	if err != nil {
		response := helper.ErrorResponse(err)
		c.JSON(http.StatusUnprocessableEntity, response)
	}
	formatter := user.UserFormatter(newUser, token)
	response := helper.APIResponse("Akun berhasil regis", http.StatusOK, "sucess", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMsg := helper.ErrorResponse(err)
		errow := gin.H{"errors": errorMsg}
		response := helper.APIResponse("Gagal login", http.StatusUnprocessableEntity, "fail", errow)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loginin, err := h.userService.Login(input)
	if err != nil {
		response := helper.APIResponse("Ada kesalahan", http.StatusUnprocessableEntity, "fail", err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	token, err := h.authService.GenerateToken(loginin.Id)
	if err != nil {
		response := helper.ErrorResponse(err)
		c.JSON(http.StatusUnprocessableEntity, response)
	}
	formatter := user.UserFormatter(loginin, token)
	c.JSON(http.StatusOK, formatter)
}

func (h *userHandler) CheckEmail(c *gin.Context) {
	var email user.EmailUserInput
	isAvail, err := h.userService.CheckEmailAvail(email)
	if err != nil {
		errorMsg := helper.ErrorResponse(err)
		errow := gin.H{"errors": errorMsg}
		response := helper.APIResponse("Terjadi kesalahan", http.StatusBadRequest, "fail", errow)
		c.JSON(http.StatusBadRequest, response)
	}

	data := gin.H{
		"is_avail": isAvail,
	}

	metaMsg := "email registered"
	if isAvail {
		metaMsg = "email not registered"
	}
	response := helper.APIResponse(metaMsg, http.StatusOK, "sucess", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Upload Avatar Gaiso", http.StatusBadRequest, "fail", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//INI PAS BAGIAN UPLOAD
	path := "/img/" + file.Filename
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Upload Avatar Gaiso", http.StatusBadRequest, "fail", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// INI PAS BAGIAN UPDATE
	userId := 1
	_, err = h.userService.SaveAvatar(userId, file.Filename)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Upload Avatar Gaiso", http.StatusBadRequest, "fail", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Terjadi kesalahan", http.StatusOK, "success", data)
	c.JSON(http.StatusBadRequest, response)
}
