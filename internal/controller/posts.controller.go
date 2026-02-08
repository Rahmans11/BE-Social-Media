package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
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

type PostController struct {
	postService *service.PostsService
}

func NewPostController(postService *service.PostsService) *PostController {
	return &PostController{
		postService: postService,
	}
}

// @Summary      Create Post
// @Description  Create Post
// @Security	BearerAuth
// @Tags         Post
// @Accept multipart/form-data
// @Produce      json
// @Param	caption formData string false "data to create post"
// @Param	image formData file false "data to create post"
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.Response
// @Failure      401  {object}  dto.Response
// @Failure      403  {object}  dto.Response
// @Failure      500  {object}  dto.Response
// @Router       /post/ [post]
func (p PostController) CreatePost(c *gin.Context) {

	var data dto.CreatePosts
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

	if data.Image != nil {
		ext := path.Ext(data.Image.Filename)
		log.Println(data.Image.Filename)
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

		if data.Image.Size > 10*1024*1024 {
			c.JSON(http.StatusBadRequest, dto.Response{
				Msg:     "Max file 10 mb",
				Error:   "Bad Request",
				Success: false,
				Data:    []any{},
			})
			return
		}

		filename := fmt.Sprintf("%d_post_%d%s", time.Now().UnixNano(), accessToken.Id, ext)

		if e := c.SaveUploadedFile(data.Image, filepath.Join("public", "post", filename)); e != nil {
			log.Println(e.Error())
			c.JSON(http.StatusInternalServerError, dto.Response{
				Msg:     "Internal Server Error",
				Success: false,
				Error:   "internal server error",
				Data:    []any{},
			})
			return
		}

	}

	post, e := p.postService.CreatePosts(c.Request.Context(), data, accessToken.Id)
	if e != nil {
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

	response := dto.Posts{
		Id:      post.Id,
		UserId:  post.UserId,
		Caption: post.Caption,
		Image:   post.Image,
	}

	c.JSON(http.StatusCreated, dto.Response{
		Msg:     "Created",
		Success: true,
		Data: []any{
			response,
		},
	})
}
