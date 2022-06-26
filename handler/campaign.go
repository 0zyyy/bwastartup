package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
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
		response := helper.APIResponse("Failed to fetch", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
	}

	// kalo sukses
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaigns)
	c.JSON(http.StatusOK, response)
}
