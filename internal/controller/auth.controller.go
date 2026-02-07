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

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// @Summary      Login
// @Description  Login to get credential
// @Tags         auth
// @Produce      json
// @Param        user	body	dto.AuthRequest	true	"user body"
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.Response
// @Failure      500  {object}  dto.Response
// @Router       /auth/ [post]
func (a AuthController) Login(c *gin.Context) {
	var loginData dto.AuthRequest
	if e := c.ShouldBindJSON(&loginData); e != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Msg:     "Internal server error",
			Success: false,
			Error:   "Internal server error",
			Data:    []any{},
		})
		return
	}

	data, e := a.authService.Login(c, loginData)
	if e != nil {
		log.Println(e.Error())
		c.JSON(http.StatusBadRequest, dto.Response{
			Msg:     "Bad request",
			Success: false,
			Error:   "Wrong Email or Password",
			Data:    []any{},
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Msg:     "Login Success",
		Success: true,
		Data:    []any{data},
	})
}

// @Summary      Register
// @Description  Register to be user
// @Tags         auth
// @Produce      json
// @Param        user	body	dto.AuthRequest	true	"user body"
// @Success      201  {object}  dto.Response
// @Failure      400  {object}  dto.Response
// @Failure      500  {object}  dto.Response
// @Router       /auth/register [post]
func (a AuthController) Register(c *gin.Context) {
	var newUser dto.AuthRequest
	if e := c.ShouldBindJSON(&newUser); e != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Msg:     "Internal Server Error",
			Success: false,
			Error:   "internal server error",
			Data:    []any{},
		})
		return
	}

	data, e := a.authService.Register(c.Request.Context(), newUser)
	if e != nil {
		if errors.Is(e, err.ExistingEmail) {
			c.JSON(http.StatusBadRequest, dto.Response{
				Msg:     e.Error(),
				Success: false,
				Error:   "Bad Request",
				Data:    []any{},
			})
			return
		}
		if errors.Is(e, err.InvalidFormatEmail) {
			c.JSON(http.StatusBadRequest, dto.Response{
				Msg:     e.Error(),
				Success: false,
				Error:   "Bad Request",
				Data:    []any{},
			})
			return
		}
		if errors.Is(e, err.InvalidFormatPassword) {
			c.JSON(http.StatusBadRequest, dto.Response{
				Msg:     e.Error(),
				Success: false,
				Error:   "Bad Request",
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
		Msg:     "Register Success",
		Success: true,
		Data:    []any{data},
	})
}

// @Summary      Logout
// @Description  Post Logout
// @Security	BearerAuth
// @Tags         auth
// @Produce      json
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.Response
// @Failure      401  {object}  dto.Response
// @Failure      403  {object}  dto.Response
// @Failure      500  {object}  dto.Response
// @Router       /auth/logout [delete]
func (a AuthController) Logout(c *gin.Context) {

	token, _ := c.Get("token")
	accessToken, _ := token.(pkg.JWTClaims)

	log.Println(accessToken.ExpiresAt)
	log.Println(accessToken.ExpiresAt.Time)
	test := accessToken.ExpiresAt.IsZero()
	log.Println(test)
	err := a.authService.Logout(c.Request.Context(), accessToken.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Msg:     "Internal server error",
			Success: false,
			Error:   "Internal server error",
			Data:    []any{},
		})
		return
	}
	c.JSON(http.StatusOK, dto.Response{
		Msg:     "OK",
		Success: true,
		Data:    []any{},
	})
}
