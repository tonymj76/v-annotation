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
	router.Use(AuthAPIKey())
	router.MaxMultipartMemory = 15 << 20 // 15 MiB
	db, err := datastore.NewDBStore(logrus.New())

	if err != nil {
		return nil, err
	}
	service := rest_service.NewRESTServices(db)

	router.POST("/video", service.CreateVideo)
	router.POST("/video-from-disk", service.CreateVideoFromDisk)
	router.POST("/video/:id/annotation", service.CreateAnnotation)
	router.PATCH("/video/:id/annotation/:annotationID", service.UpdateAnnotation)
	router.GET("/video", service.FetchVideos)
	router.GET("/video/:id", service.FetchVideo)
	router.GET("/annotation/:annotationID", service.FetchAnnotationByID)
	router.GET("/video/:id/annotation", service.FetchAnnotations)
	router.DELETE("/video/:id", service.DeleteVideo)
	router.DELETE("/annotation/:id", service.DeleteAnnotationByID)
	return router, nil
}
