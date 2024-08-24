// deadline: in 2-weeks Infinite
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homebrew-ec-foss/event-loop-backend/handlers"
)

func main() {
	r := gin.Default()
	r.Use(handlers.CorsMiddleware())

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	////////////////////////////////////////////////

	// Pre event endpoints
	// handles creation for event
	r.POST("/create", handlers.HandleCreate)

	////////////////////////////////////////////////

	// Generic functions for admin control

	r.GET("/search", func(ctx *gin.Context) {})

	// fetch participants details based on JWT ID
	// fetched based on a qr code
	r.POST("/qrsearch", handlers.HandleQRFetch)

	r.GET("/participant", handlers.HandleParticipantFetch)
	r.POST("/participant", handlers.HandleParticipantUpdate)

	// Additional team addition besides CSV
	// future prospect
	r.POST("/createteam", func(ctx *gin.Context) {})
	r.POST("/login", handlers.HandleLogin)

	////////////////////////////////////////////////

	// Endpoints accessed during events
	// eg: Crossing checkpoints, etc.

	// TODO: Handle checking by scanner
	r.PUT("/checkin", handlers.HandleCheckin)

	r.PUT("/checkout", handlers.HandleCheckout)

	r.PUT("/checkpoint", handlers.HandleCheckpoint)

	err := r.RunTLS("0.0.0.0:8080", "localhost.crt", "localhost.key")
	if err != nil {
		log.Fatal(err)
	}
}
