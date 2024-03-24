package v1

import (
	"net/http"

	"promotion/internal/dto"
	"promotion/internal/model"
	"promotion/internal/server/v1/resp"
	"promotion/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type campaignRoutes struct {
	c usecase.CampaignUsecase
}

func newCampaignRoutes(handler *gin.RouterGroup, c usecase.CampaignUsecase) {
	r := &campaignRoutes{c}
	h := handler.Group("/campaigns")
	{
		h.POST("/", r.CreateCampaign)
		h.GET("/", r.GetCampaigns)
	}
}

func (r *campaignRoutes) CreateCampaign(ctx *gin.Context) {
	var dtoCampaign dto.CreateCampaign
	if err := ctx.ShouldBindJSON(&dtoCampaign); err != nil {
		resp.RespondErrorWithCode(ctx, err, http.StatusBadRequest)
		return
	}

	var creatingCampaign model.Campaign
	copier.Copy(&creatingCampaign, &dtoCampaign)

	campaign, err := r.c.CreateCampaign(ctx, creatingCampaign)
	if err != nil {
		resp.RespondError(ctx, err)
		return
	}

	var createdCampaign dto.Campaign
	copier.Copy(&createdCampaign, &campaign)
	resp.RespondDataWithCode(ctx, createdCampaign, http.StatusCreated)
}

func (r *campaignRoutes) GetCampaigns(ctx *gin.Context) {
	result, err := r.c.GetCampaigns(ctx.Request.Context())
	if err != nil {
		resp.RespondError(ctx, err)
		return
	}

	var campaigns []*dto.Campaign
	copier.Copy(&campaigns, &result)
	resp.RespondData(ctx, campaigns)
}
