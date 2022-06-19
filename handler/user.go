package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	// tangkap input
	// map input ke registeruserinput
	// struct passing ke repository

	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Akun gagal regis", http.StatusBadRequest, "fail", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Akun gagal regis", http.StatusBadRequest, "fail", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// update token
	token := "tokentkoentok"
	formatter := user.UserFormatter(newUser, token)
	response := helper.APIResponse("Akun berhasil regis", http.StatusOK, "sucess", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UpdateUser(c *gin.Context) {
	//TANGKAP INPUT
}
