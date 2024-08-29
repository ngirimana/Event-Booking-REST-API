package routes

import (
	"net/http"
	"strconv"

	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

// getEvent godoc
// @Summary Get an event by ID
// @Description Retrieves a single event by its ID.
// @Tags events
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Success 200 {object} models.Event "Successfully retrieved event"
// @Failure 400 {object} map[string]string "Invalid event ID"
// @Failure 404 {object} map[string]string "Event not found"
// @Router /events/{id} [get]
func getEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID."})
		return
	}
	event, err := models.GetEvent(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found."})
		return
	}

	context.JSON(http.StatusOK, event)
}

// getEvents godoc
// @Summary Get all events
// @Description Retrieves all events.
// @Tags events
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Event "Successfully retrieved events"
// @Failure 500 {object} map[string]string "Could not fetch events"
// @Router /events [get]
func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch events. Please try again later."})
		return
	}

	context.JSON(http.StatusOK, events)
}

// createEvent godoc
// @Summary Create a new event
// @Description Creates a new event and saves it to the database.
// @Tags events
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer Token"
// @Param event body models.Event true "Event data"
// @Success 201 {object} map[string]interface{} "Event created successfully"
// @Failure 400 {object} map[string]string "Could not parse request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /events [post]
func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	event.UserID = context.GetInt64("userId")

	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Event created!", "event": event})
}

// updateEvent godoc
// @Summary Update an event by ID
// @Description Updates an existing event by its ID.
// @Tags events
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Param event body models.Event true "Updated event data"
// @Success 200 {object} map[string]interface{} "Event updated successfully"
// @Failure 400 {object} map[string]string "Invalid event ID or request data"
// @Failure 401 {object} map[string]string "You are not authorized to update this event."
// @Failure 404 {object} map[string]string "Event not found"
// @Failure 500 {object} map[string]string "Could not update event"
// @Router /events/{id} [put]
func updateEvent(context *gin.Context) {

	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID."})
		return
	}
	userID := context.GetInt64("userId")
	event, error := models.GetEvent(id)
	if error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found."})
		return
	}

	if event.UserID != userID {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to update this event."})
		return
	}
	var updatedEvent models.Event

	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	updatedEvent.ID = id
	err = updatedEvent.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Event updated!", "event": updatedEvent})

}

// deleteEvent godoc
// @Summary Delete an event by ID
// @Description Deletes an event by its ID.
// @Tags events
// @Accept  json
// @Produce  json
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]string "Event deleted successfully"
// @Failure 400 {object} map[string]string "Invalid event ID"
// @Failure 404 {object} map[string]string "Event not found"
// @Failure 401 {object} map[string]string "You are not authorized to delete this event."
// @Failure 500 {object} map[string]string "Could not delete event"
// @Router /events/{id} [delete]
func deleteEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID."})
		return
	}
	event, error := models.GetEvent(id)
	if error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found."})
		return
	}
	userID := context.GetInt64("userId")
	if event.UserID != userID {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to delete this event."})
		return
	}

	err = event.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete event."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Event deleted!"})
}
