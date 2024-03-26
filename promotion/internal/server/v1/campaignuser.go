package v1

import (
	"promotion/internal/dto"
	"promotion/internal/server/v1/resp"
	"promotion/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type campaignUserRoutes struct {
	c usecase.CampaignUserUsecase
}

func newCampaignUserRoutes(handler *gin.RouterGroup, c usecase.CampaignUserUsecase) {
	r := &campaignUserRoutes{c}
	h := handler.Group("/campaign-users")
	{
		h.GET("/", r.GetCampaignUsers)
	}
}

func (r *campaignUserRoutes) GetCampaignUsers(ctx *gin.Context) {
	var filter dto.CampaignUserFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		resp.RespondError(ctx, err)
		return
	}

	result, err := r.c.GetCampaignUsers(ctx.Request.Context(), filter)
	if err != nil {
		resp.RespondError(ctx, err)
		return
	}

	var campUsers []*dto.CampaignUser
	copier.Copy(&campUsers, &result)
	resp.RespondData(ctx, campUsers)
}
