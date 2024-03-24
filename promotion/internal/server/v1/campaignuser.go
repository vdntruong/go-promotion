package v1

import (
	"net/http"

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
		h.POST("/", r.CreateCampaignUser)
		h.GET("/", r.GetCampaignUsers)
	}
}

func (r *campaignUserRoutes) CreateCampaignUser(ctx *gin.Context) {
	var req dto.CreateCampaignUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.RespondErrorWithCode(ctx, err, http.StatusBadRequest)
		return
	}

	result, err := r.c.CreateCampaignUser(ctx, req.CampaignExtID, req.UserExtID, req.RegisterDate)
	if err != nil {
		resp.RespondError(ctx, err)
		return
	}

	var createdCampUser dto.CampaignUser
	copier.Copy(&createdCampUser, &result)
	resp.RespondDataWithCode(ctx, createdCampUser, http.StatusCreated)
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
