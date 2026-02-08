package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/Rahmans11/final-phase-3/internal/dto"
	"github.com/Rahmans11/final-phase-3/internal/err"
	"github.com/Rahmans11/final-phase-3/internal/service"
	"github.com/Rahmans11/final-phase-3/pkg"
	"github.com/gin-gonic/gin"
)

type FollowsController struct {
	followsService *service.FollowsService
}

func NewFollowsController(followsService *service.FollowsService) *FollowsController {
	return &FollowsController{
		followsService: followsService,
	}
}

// @Summary      Create Follows
// @Description  Create Follows
// @Security	BearerAuth
// @Tags         Follows
// @Produce      json
// @Param	followed_id body	dto.AddFollowed	true "data to add follows"
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.Response
// @Failure      401  {object}  dto.Response
// @Failure      403  {object}  dto.Response
// @Failure      500  {object}  dto.Response
// @Router       /follows/add-followed [post]
func (f FollowsController) AddFollowed(c *gin.Context) {

	var data dto.AddFollowed
	if e := c.ShouldBindJSON(&data); e != nil {
		log.Println("binding", e.Error())
		c.JSON(http.StatusInternalServerError, dto.Response{
			Msg:     "Internal Server Error",
			Success: false,
			Error:   "internal server error",
			Data:    []any{},
		})
		return
	}

	token, _ := c.Get("token")
	accessToken, _ := token.(pkg.JWTClaims)

	if e := f.followsService.AddFollowed(c.Request.Context(), accessToken.Id, data.FollowedId); e != nil {
		log.Println(data.FollowedId)
		log.Println(e.Error())
		if errors.Is(e, err.ErrNoRowsUpdated) {
			c.JSON(http.StatusNotFound, dto.Response{
				Msg:     e.Error(),
				Success: false,
				Error:   "User not found",
				Data:    []any{},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.Response{
			Msg:     "Internal Server Error",
			Success: false,
			Error:   "internal server error",
			Data:    []any{},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Msg:     "OK",
		Success: true,
	})
}
