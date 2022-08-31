package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Fajar-Islami/scrapping_test/controller"
	"github.com/Fajar-Islami/scrapping_test/infrastructure/container"
	"github.com/Fajar-Islami/scrapping_test/service"
	"github.com/gin-gonic/gin"
)

func Start(cont *container.Container, server *gin.Engine) {

	var (
		ctx context.Context = context.Background()

		trackingService    service.TrackingService       = service.NewTrackingService(ctx)
		trackingController controller.TrackingController = controller.NewTrackingController(ctx, trackingService)
	)

	trackingRouter := server.Group("api/tracking")
	{
		trackingRouter.GET("/", trackingController.GetDataTracking)
	}

	server.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	port := ":" + fmt.Sprint(cont.Apps.Port)
	server.Run(port)

}
