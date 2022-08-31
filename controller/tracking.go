package controller

import (
	"context"
	"net/http"

	"github.com/Fajar-Islami/scrapping_test/helper"
	"github.com/Fajar-Islami/scrapping_test/model"
	"github.com/Fajar-Islami/scrapping_test/service"
	"github.com/gin-gonic/gin"
)

const URI = "https://gist.githubusercontent.com/nubors/eecf5b8dc838d4e6cc9de9f7b5db236f/raw/d34e1823906d3ab36ccc2e687fcafedf3eacfac9/jne-awb.html"

type TrackingController interface {
	GetDataTracking(c *gin.Context)
}

type trackingControllerImpl struct {
	context context.Context
	service service.TrackingService
}

func NewTrackingController(ctx context.Context, s service.TrackingService) TrackingController {
	return &trackingControllerImpl{
		context: ctx,
		service: s,
	}
}

func (tc *trackingControllerImpl) GetDataTracking(c *gin.Context) {
	res, err := tc.service.GetDataTracking(URI)
	if err != nil {
		// response := helper.BuildErrorResponse("Created article failed", "Something wrong when Created article failed", helper.EmptyObj{})
		response := helper.BuildErrorResponse("Get Data Tracked failed", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	resp := helper.BuildSuccessResponse[model.DataStruct]("060101", "Delivery tracking detail fetched successfully", res)
	c.JSON(http.StatusOK, resp)

}
