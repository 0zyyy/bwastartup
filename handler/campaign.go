package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//tangkap parameter handler
// check parameter, manggil service.
// service menentukan apakah getAll, getByid, atau GetBockerId
// repo akses ke db

type campaignHandler struct {
	campaignService campaign.CampaignService
}

func NewCampaignHandler(campaignService campaign.CampaignService) *campaignHandler {
	return &campaignHandler{campaignService}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))
	campaigns, err := h.campaignService.GetCampaign(userId)
	if err != nil {
		response := helper.APIResponse("Failed to fetch", http.StatusBadRequest, "error", campaign.FormatCampaigns(campaigns))
		c.JSON(http.StatusBadRequest, response)
	}

	// kalo sukses
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail", http.StatusBadRequest, "fail", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignAda, err := h.campaignService.GetCampaignById(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail", http.StatusBadRequest, "error", campaign.FormatCampaign(campaignAda))
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Campaign Detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignAda))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) SaveCampaign(c *gin.Context) {
	var input campaign.NewCampaign
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ErrorResponse(err)
		errMsg := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "fail", errMsg)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	currUser := c.MustGet("currentUser").(user.User)
	input.User = currUser
	newCampaign, err := h.campaignService.AddCampaign(input)
	if err != nil {
		// errors := helper.ErrorResponse(err)
		// errMsg := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "fail", err)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.APIResponse("Successfully create campaign", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

// repo data campaign
// panggil service, input dari user, input yang ada di uri
// mapping input ke input struct (handler)
// handler
func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputId campaign.GetCampaignDetailInput
	var inputData campaign.NewCampaign
	err := c.ShouldBindUri(&inputId)
	if err != nil {
		errors := helper.ErrorResponse(err)
		errMsg := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "fail", errMsg)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.ErrorResponse(err)
		errMsg := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "fail", errMsg)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updatedCampaign, err := h.campaignService.UpdateCampaign(inputId, inputData)
	if err != nil {
		errors := helper.ErrorResponse(err)
		errMsg := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "fail", errMsg)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	response := helper.APIResponse("Successfully update campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UploadCampaignImage(c *gin.Context) {
	var input campaign.GetCampaignImageInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.ErrorResponse(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "fail", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "fail", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userID := currentUser.Id

	path := fmt.Sprintf("img/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "fail", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.campaignService.SaveCampaignImage(input, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "fail", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Successfully upload campaign image", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
