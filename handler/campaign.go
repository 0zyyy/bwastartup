package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"
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
