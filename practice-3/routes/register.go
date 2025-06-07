package routes

import (
	"net/http"
	"strconv"

	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID."})
		return
	}

	userId := context.GetInt64("userId")

	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve event."})
		return
	}

	err = event.RegisterForEvent(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register for event."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Successfully registered for event."})
}

func unregisterFromEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID."})
		return
	}

	userId := context.GetInt64("userId")

	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve event."})
		return
	}

	err = event.CancelRegistration(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not unregister from event."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Successfully unregistered from event."})
}
