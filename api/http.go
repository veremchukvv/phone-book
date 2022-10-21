package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"integ/entities"
	"net/http"
	"strconv"
)

type RelationService interface {
	SaveContacts(ctx context.Context, userID int, contacts entities.ContactList) error
}

type HTTPHandler struct {
	svc RelationService
	log *logrus.Entry
}

func NewHTTPHandler(svc RelationService, log *logrus.Entry) *HTTPHandler {
	return &HTTPHandler{
		svc: svc,
		log: log,
	}
}

type AddContactRequest struct {
	Contacts entities.ContactList `json:"contacts"`
}

func (h *HTTPHandler) AddContacts(c *gin.Context) {
	userIDstr := c.Param("userID")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var contactRequest AddContactRequest

	if err := c.BindJSON(&contactRequest); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := h.svc.SaveContacts(c.Request.Context(), userID, contactRequest.Contacts); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *HTTPHandler) Friends(c *gin.Context) {
	panic("not implemented")
}

func (h *HTTPHandler) Name(c *gin.Context) {
	panic("not implemented")
}
