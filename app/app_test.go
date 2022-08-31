package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fajar-Islami/scrapping_test/controller"
	"github.com/Fajar-Islami/scrapping_test/helper"
	"github.com/Fajar-Islami/scrapping_test/infrastructure/container"
	redisClient "github.com/Fajar-Islami/scrapping_test/infrastructure/db/redis"
	"github.com/Fajar-Islami/scrapping_test/model"
	"github.com/Fajar-Islami/scrapping_test/repository"
	"github.com/Fajar-Islami/scrapping_test/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func setupApp() (controller.TrackingController, *redis.Client) {
	var (
		cont                                                  = container.New("../.env")
		ctx                context.Context                    = context.Background()
		redisDb            *redis.Client                      = redisClient.NewRedisClient(*cont.Redis)
		trackingRedisRepo  repository.RedisTrackingRepository = repository.NewRedisTrackingRepo(redisDb)
		trackingService    service.TrackingService            = service.NewTrackingService(ctx, trackingRedisRepo)
		trackingController controller.TrackingController      = controller.NewTrackingController(ctx, trackingService)
	)
	return trackingController, redisDb
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestHealthHandler(t *testing.T) {
	fmt.Println("")
	fmt.Println("=================TestHealthHandler=================")
	mockResponse := `{"status":"ok"}`
	r := setupRouter()
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
	fmt.Println("=================EndTestHealthHandler=================")
	fmt.Println("")
}

func TestTrackingHandler(t *testing.T) {
	fmt.Println("")
	fmt.Println("=================TestTrackingHandler=================")

	r := setupRouter()
	controller, redis := setupApp()
	defer redis.Close()
	uri := "/api/tracking"

	r.GET(uri, controller.GetDataTracking)

	req := httptest.NewRequest(http.MethodGet, uri, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responBody, _ := ioutil.ReadAll(w.Body)
	respHelper := helper.ResponJSON[model.DataStruct]{}
	json.Unmarshal(responBody, &respHelper)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "060101", respHelper.Status.Code)
	assert.Equal(t, "Delivery tracking detail fetched successfully", respHelper.Status.Message)

	fmt.Println("=================EndTestTrackingHandler=================")
	fmt.Println("")
}
