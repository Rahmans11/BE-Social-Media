package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Rahmans11/final-phase-3/internal/dto"
	"github.com/Rahmans11/final-phase-3/internal/err"
	"github.com/Rahmans11/final-phase-3/internal/service"
	"github.com/Rahmans11/final-phase-3/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ProfileController struct {
	profileService *service.ProfileService
}

func NewProfileController(profileService *service.ProfileService) *ProfileController {
	return &ProfileController{
		profileService: profileService,
	}
}

// @Summary      Profile
// @Description  Get Profile
// @Security	BearerAuth
// @Tags         Profile
// @Produce      json
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.Response
// @Failure      401  {object}  dto.Response
// @Failure      403  {object}  dto.Response
// @Failure      500  {object}  dto.Response
// @Router       /profile/ [get]
func (p ProfileController) GetProfile(c *gin.Context) {

	token, _ := c.Get("token")
	accessToken, _ := token.(pkg.JWTClaims)

	data, err := p.profileService.GetProfile(c.Request.Context(), accessToken.Id)
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
		Data:    []any{data},
	})
}

// @Summary      Other Profile
// @Description  Get Other Profile
// @Security	BearerAuth
// @Tags         Profile
// @Produce      json
// @Param        id	path integer	false	"id param"
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.Response
// @Failure      401  {object}  dto.Response
// @Failure      403  {object}  dto.Response
// @Failure      500  {object}  dto.Response
// @Router       /profile/user/{id} [get]
func (p ProfileController) GetOtherProfile(c *gin.Context) {

	var param dto.ProfileParam
	if err := c.ShouldBindUri(&param); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, dto.Response{
			Msg: "Internal server error",
		})
		return
	}

	data, err := p.profileService.GetOtherProfile(c.Request.Context(), param.Id)
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
		Data:    []any{data},
	})
}

// @Summary      Edit Profile
// @Description  Patch Profile
// @Security	BearerAuth
// @Tags         Profile
// @Accept multipart/form-data
// @Produce      json
// @Param	first_name formData string false "data to edit profile"
// @Param	last_name formData string false "data to edit profile"
// @Param	phone_number formData string false "data to edit profile"
// @Param	avatar formData file false "data to edit profile"
// @Param	bio formData string false "data to edit profile"
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.Response
// @Failure      401  {object}  dto.Response
// @Failure      403  {object}  dto.Response
// @Failure      500  {object}  dto.Response
// @Router       /profile/ [patch]
func (p ProfileController) EditProfile(c *gin.Context) {

	var data dto.EditProfile
	if e := c.ShouldBindWith(&data, binding.FormMultipart); e != nil {
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

	if data.Avatar != nil {
		ext := path.Ext(data.Avatar.Filename)
		log.Println(data.Avatar.Filename)
		re := regexp.MustCompile("^[.](jpg|png)$")
		if !re.MatchString(ext) {
			c.JSON(http.StatusBadRequest, dto.Response{
				Msg:     err.ErrInvalidExt.Error(),
				Error:   "Bad Request",
				Success: false,
				Data:    []any{},
			})
			return
		}

		if data.Avatar.Size > 2*800*600 {
			c.JSON(http.StatusBadRequest, dto.Response{
				Msg:     "Max file 2 mb",
				Error:   "Bad Request",
				Success: false,
				Data:    []any{},
			})
			return
		}

		filename := fmt.Sprintf("%d_profile_%d%s", time.Now().UnixNano(), accessToken.Id, ext)

		if e := c.SaveUploadedFile(data.Avatar, filepath.Join("public", "profile", filename)); e != nil {
			log.Println(e.Error())
			c.JSON(http.StatusInternalServerError, dto.Response{
				Msg:     "Internal Server Error",
				Success: false,
				Error:   "internal server error",
				Data:    []any{},
			})
			return
		}

		data.Avatar.Filename = fmt.Sprint(filename)

		photo, e := p.profileService.GetPhoto(c.Request.Context(), accessToken.Id)
		if e != nil {
			c.JSON(http.StatusInternalServerError, dto.Response{
				Msg:     "Internal server error",
				Success: false,
				Error:   "Internal server error",
				Data:    []any{},
			})
			return
		}

		filePath := fmt.Sprintf("D:/KodaCourse/test/public/profile/%v", photo)

		if photo != "" {
			err := os.Remove(filePath)
			if err != nil {
				log.Printf("Error deleting file: %v", err)
			}
		}

		fmt.Printf("File %s successfully deleted\n", filePath)
	}

	log.Println(data.Avatar)

	if e := p.profileService.UpdateProfile(c.Request.Context(), data, accessToken.Id); e != nil {
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

	if data.Avatar != nil {
		c.JSON(http.StatusOK, dto.Response{
			Msg:     "OK",
			Success: true,
			Data: []any{
				fmt.Sprint(data.Avatar.Filename),
			},
		})
	}

	c.JSON(http.StatusOK, dto.Response{
		Msg:     "OK",
		Success: true,
	})
}
