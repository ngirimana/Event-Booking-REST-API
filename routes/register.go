package routes

import (
	"net/http"
	"strconv"

	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

// @Summary Register for an event
// @Description Register the authenticated user for the specified event
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 201 {object} map[string]string  "successfully registered for event"
// @Failure 400 {object} map[string]string  "invalid event ID"
// @Failure 404 {object} map[string]string "event not found"
// @Failure 500 {object} map[string]string  "could not register for event"
// @Router /events/{id}/register [post]
func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid event id"})
		return
	}

	event, err := models.GetEvent(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "event not found"})
		return
	}

	err = event.Register(userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not register for event"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "successfully registered for event"})

}

// @Summary Cancel registration for an event
// @Description Cancel the authenticated user's registration for the specified event
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]string "successfully cancelled registration"
// @Failure 400 {object} map[string]string "invalid event ID"
// @Failure 404 {object} map[string]string "event not found"
// @Failure 401 {object} map[string]string 	"you can cancel registration for your own event"
// @Failure 500 {object} map[string]string  "could not cancel registration"
// @Router /events/{id}/register [delete]
func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid event id"})
		return
	}

	event, err := models.GetEvent(eventId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "event not found"})
		return
	}

	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "you can cancel registration for your own event"})
		return
	}

	err = event.CancelRegistration(userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not cancel registration"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "successfully cancelled registration"})

}
