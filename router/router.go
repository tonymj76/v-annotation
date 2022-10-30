package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tonymj76/video-annotator/datastore"
	"github.com/tonymj76/video-annotator/rest_service"
)

func Router() (*gin.Engine, error) {
	router := gin.Default()
	//if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
	//	return nil, err
	//}
	db, err := datastore.NewDBStore(logrus.New())

	if err != nil {
		return nil, err
	}
	service := rest_service.NewRESTServices(db)

	router.POST("/video", service.CreateVideo)
	router.POST("/video/:id/segment", service.CreateSegment)
	router.GET("/video", service.FetchVideos)
	router.GET("/video/:id", service.FetchVideo)
	router.GET("/segment/:segmentID", service.FetchSegmentByID)
	router.GET("/video/:id/segment", service.FetchSegments)
	router.DELETE("/video/:id", service.DeleteVideo)
	router.DELETE("/segment/:id", service.DeleteSegmentByID)

	return router, nil
}
