package v1

import (
	"net/http"

	"ekyc/internal/dto"
	"ekyc/internal/server/v1/resp"
	"ekyc/internal/usecase"

	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type userRoutes struct {
	c usecase.UserUsecase
}

func newUserRoutes(handler *gin.RouterGroup, c usecase.UserUsecase) {
	r := &userRoutes{c}
	h := handler.Group("/users")
	{
		h.POST("/sign-up", limit.MaxAllowed(100), r.SignUp)
		h.POST("/sign-in", r.SignIn)
	}
}

// @Summary     User signup
// @Description user register new account
// @ID          signup
// @Tags  	    Auth
// @Accept      json
// @Produce     json
// @Param request body dto.SignUpUser true "query params"
// @Success     200 {object} dto.User
// @Router      /users/sign-up [post]
func (r *userRoutes) SignUp(ctx *gin.Context) {
	var req dto.SignUpUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.RespondErrorWithCode(ctx, err, http.StatusBadRequest)
		return
	}

	result, err := r.c.SignUp(ctx, req.Email, req.Password, req.CampaignExtID)
	if err != nil {
		resp.RespondError(ctx, err)
		return
	}

	var createdUser dto.User
	copier.Copy(&createdUser, &result)
	resp.RespondDataWithCode(ctx, createdUser, http.StatusCreated)
}

func (r *userRoutes) SignIn(ctx *gin.Context) {
	var req dto.SignInUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.RespondErrorWithCode(ctx, err, http.StatusBadRequest)
		return
	}

	result, err := r.c.SignIn(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		resp.RespondError(ctx, err)
		return
	}

	resp.RespondData(ctx, result)
}
