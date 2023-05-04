package handler

import (
	"app/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ChangeProdsCount(c *gin.Context) {
	var upd models.StoreChange

	err := c.ShouldBindJSON(&upd)
	if err != nil {
		h.handlerResponse(c, "change store", http.StatusBadRequest, err.Error())
		return
	}

	massage, err := h.storages.Code().Exam(&upd)
	if err != nil {
		h.handlerResponse(c, "_", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "update stock", http.StatusAccepted, massage)

}
