package routes

import (
	"github.com/amarjeetdev/ev-charging-app/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.POST("/auth/register", controllers.Register)
	api.POST("/auth/login", controllers.Login)

	api.GET("/stations", controllers.GetAllStations)
	api.GET("/stations/:id", controllers.GetStationDetails)
	api.POST("/stations", controllers.CreateStation)

	api.POST("/bookings", controllers.CreateBooking)
	api.GET("/bookings/user", controllers.GetUserBookings)
}
