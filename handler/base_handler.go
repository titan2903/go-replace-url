package handler

import (
	"fmt"
	"net/http"
	"replace-url-gin/usecase"

	"github.com/gin-gonic/gin"
)

type BaseHandler interface {
	HandleReplaceImage(c *gin.Context)
	HandleReplaceImageUrl(c *gin.Context)
}

type baseHandler struct {
	usecase usecase.BaseUsecase
}

func (h handler) BaseHandler() BaseHandler {
	return baseHandler{
		usecase: h.usecase.BaseHandler(),
	}
}

func (h baseHandler) HandleReplaceImage(c *gin.Context) {
	err := h.usecase.ReplaceImage()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Something wrong!")
		return
	}

	c.JSON(http.StatusOK, "Success!")
}

func (h baseHandler) HandleReplaceImageUrl(c *gin.Context) {
	fmt.Println("msuk replace")
	err := h.usecase.ReplaceImageUrl()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Something wrong!")
		return
	}

	c.JSON(http.StatusOK, "Success!")
}
