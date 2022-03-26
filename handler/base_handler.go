package handler

import (
	"net/http"
	"replace-url-gin/usecase"

	"github.com/gin-gonic/gin"
)

type BaseHandler interface {
	HandleReplaceImage(c *gin.Context)
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
