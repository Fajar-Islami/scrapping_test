package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Fajar-Islami/scrapping_test/controller"
	"github.com/Fajar-Islami/scrapping_test/infrastructure/container"
	redisClient "github.com/Fajar-Islami/scrapping_test/infrastructure/db/redis"
	"github.com/Fajar-Islami/scrapping_test/repository"
	"github.com/Fajar-Islami/scrapping_test/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func Start(cont *container.Container, server *gin.Engine) {

	var (
		ctx                context.Context                    = context.Background()
		redisDb            *redis.Client                      = redisClient.NewRedisClient(*cont.Redis)
		trackingRedisRepo  repository.RedisTrackingRepository = repository.NewRedisTrackingRepo(redisDb)
		trackingService    service.TrackingService            = service.NewTrackingService(ctx, trackingRedisRepo)
		trackingController controller.TrackingController      = controller.NewTrackingController(ctx, trackingService)
	)

	defer redisDb.Close()

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
