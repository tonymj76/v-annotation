package rest_service

import (
	"github.com/gin-gonic/gin"
	"github.com/tonymj76/video-annotator/datastore"
	"github.com/tonymj76/video-annotator/model"
	"github.com/tonymj76/video-annotator/utility"
	"net/http"
)

type RESTServices struct {
	repository datastore.Repository
}

func NewRESTServices(datastore datastore.Repository) *RESTServices {
	return &RESTServices{
		repository: datastore,
	}
}

func (rs *RESTServices) CreateVideo(ctx *gin.Context) {
	//s, _ := rs.repository.Create(ctx)
	var schema model.AnnotatedVideo
	if err := utility.ReadJSON(ctx, &schema); err != nil {
		JSON(ctx, "failed", http.StatusBadRequest, err)
		return
	}
	videoDetails, err := utility.NewVideoDetails(schema.VideoLink)
	if err != nil {
		JSON(ctx, "failed to get video details", http.StatusBadRequest, err)
		return
	}
	schema.Duration = videoDetails.Duration
	schema.FrameRate = videoDetails.Framerate
	savedSchema, err := rs.repository.Create(ctx, &schema)
	if err != nil {
		JSON(ctx, "failed to save video details", http.StatusBadRequest, err)
		return
	}
	JSON(ctx, "success", http.StatusCreated, savedSchema)
}

func (rs *RESTServices) FetchVideos(ctx *gin.Context) {
	videoSchema, err := rs.repository.FetchVideos(ctx)
	if err != nil {
		JSON(ctx, "failed to fetch videos", http.StatusInternalServerError, err)
		return
	}
	JSON(ctx, "success", http.StatusOK, videoSchema)
}

func (rs *RESTServices) FetchVideo(ctx *gin.Context) {
	videoSchema, err := rs.repository.FetchVideoByID(ctx)
	if err != nil {
		JSON(ctx, "failed to fetch videos", http.StatusInternalServerError, err)
		return
	}
	JSON(ctx, "success", http.StatusOK, videoSchema)
}

func (rs *RESTServices) FetchAnnotationByID(ctx *gin.Context) {
	segment, err := rs.repository.FetchAnnotationByID(ctx)
	if err != nil {
		JSON(ctx, "failed to fetch segments", http.StatusInternalServerError, err)
		return
	}
	JSON(ctx, "success", http.StatusOK, segment)
}

func (rs *RESTServices) FetchAnnotations(ctx *gin.Context) {
	segments, err := rs.repository.FetchAnnotations(ctx)
	if err != nil {
		JSON(ctx, "failed to fetch Annotation", http.StatusInternalServerError, err)
		return
	}
	JSON(ctx, "success", http.StatusOK, segments)
}

func (rs *RESTServices) CreateAnnotation(ctx *gin.Context) {
	videoSchema, err := rs.repository.FetchVideoByID(ctx)
	if err != nil {
		JSON(ctx, "failed to fetch videos", http.StatusInternalServerError, err)
		return
	}
	var segments model.AnnotationSegments
	if err := utility.ReadJSON(ctx, &segments); err != nil {
		JSON(ctx, "failed", http.StatusBadRequest, err)
		return
	}
	if err := utility.ValidateManyData(videoSchema, segments); err != nil {
		JSON(ctx, "failed to validate Annotation", http.StatusBadRequest, err)
		return
	}
	result, err := rs.repository.CreateAnnotationByID(ctx, &segments)
	if err != nil {
		JSON(ctx, "failed to save segment", http.StatusBadRequest, err)
		return
	}
	JSON(ctx, "success", http.StatusCreated, result)
}

func (rs *RESTServices) UpdateAnnotation(ctx *gin.Context) {
	videoSchema, err := rs.repository.FetchVideoByID(ctx)
	if err != nil {
		JSON(ctx, "failed to fetch videos", http.StatusInternalServerError, err)
		return
	}
	var segment model.AnnotatedSegment
	if err := utility.ReadJSON(ctx, &segment); err != nil {
		JSON(ctx, "failed", http.StatusBadRequest, err)
		return
	}
	if err := utility.ValidateOneData(videoSchema, segment); err != nil {
		JSON(ctx, "failed to validate Annotation", http.StatusBadRequest, err)
		return
	}
	result, err := rs.repository.UpdateAnnotationByID(ctx, &segment)
	if err != nil {
		JSON(ctx, "failed to save segment", http.StatusBadRequest, err)
		return
	}
	JSON(ctx, "success", http.StatusCreated, result)
}

func (rs *RESTServices) DeleteVideo(ctx *gin.Context) {
	value, err := rs.repository.DeleteVideo(ctx)
	if err != nil {
		JSON(ctx, "failed to delete video", http.StatusBadRequest, err)
		return
	}
	JSON(ctx, "success", http.StatusOK, value)
}

func (rs *RESTServices) DeleteAnnotationByID(ctx *gin.Context) {
	value, err := rs.repository.DeleteAnnotationByID(ctx)
	if err != nil {
		JSON(ctx, "failed to delete Annotation", http.StatusBadRequest, err)
		return
	}
	JSON(ctx, "success", http.StatusOK, value)
}
