package controllers

import (
	"bookings_view/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingsViewController struct {
	bookingsViewService service.BookingsViewService
}

func NewBookingsViewController(bookingsViewService service.BookingsViewService) *BookingsViewController {
	return &BookingsViewController{bookingsViewService: bookingsViewService}
}

func (c *BookingsViewController) GetBookingByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "booking id is required"})
		return
	}

	booking, err := c.bookingsViewService.GetBookingByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, booking)
}

func (c *BookingsViewController) GetBookingsByEventID(ctx *gin.Context) {
	eventID := ctx.Param("event_id")
	if eventID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "eventID is required"})
		return
	}

	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	bookings, err := c.bookingsViewService.GetBookingsByEventID(eventID, limit, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, bookings)
}

func (c *BookingsViewController) GetBookingsByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	bookings, err := c.bookingsViewService.GetBookingsByUserID(userID, limit, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, bookings)
}
