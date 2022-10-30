package rest_service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/tonymj76/video-annotator/datastore"
	"github.com/tonymj76/video-annotator/model"
	"github.com/tonymj76/video-annotator/utility"
	"net/http"
	"net/url"
)

// RESTServices app services
type RESTServices struct {
	repository datastore.Repository
}

// NewRESTServices create app services and link it to db
func NewRESTServices(datastore datastore.Repository) *RESTServices {
	return &RESTServices{
		repository: datastore,
	}
}

// CreateVideo creates Video and its metadata
func (rs *RESTServices) CreateVideo(ctx *gin.Context) {
	//s, _ := rs.repository.Create(ctx)
	var schema model.AnnotatedVideo
	if err := utility.ReadJSON(ctx, &schema); err != nil {
		JSON(ctx, "failed", http.StatusBadRequest, err)
		return
	}
	url, err := url.Parse(schema.VideoLink)
	if err != nil {
		JSON(ctx, "invalid video Link", http.StatusBadRequest, err)
		return
	}
	videoDetails, err := utility.NewVideoDetails(url.String())
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

// FetchVideos gets all videos and its metadata
func (rs *RESTServices) FetchVideos(ctx *gin.Context) {
	videoSchema, err := rs.repository.FetchVideos(ctx)
	if err != nil {
		JSON(ctx, "failed to fetch videos", http.StatusInternalServerError, err)
		return
	}
	JSON(ctx, "success", http.StatusOK, videoSchema)
}

// FetchVideo get a video and its metadata
func (rs *RESTServices) FetchVideo(ctx *gin.Context) {
	videoSchema, err := rs.repository.FetchVideoByID(ctx)
	if err != nil {
		JSON(ctx, "failed to fetch videos", http.StatusInternalServerError, err)
		return
	}
	JSON(ctx, "success", http.StatusOK, videoSchema)
}

// FetchAnnotationByID get annotation and its metadata
func (rs *RESTServices) FetchAnnotationByID(ctx *gin.Context) {
	segment, err := rs.repository.FetchAnnotationByID(ctx)
	if err != nil {
		JSON(ctx, "failed to fetch segments", http.StatusInternalServerError, err)
		return
	}
	JSON(ctx, "success", http.StatusOK, segment)
}

//FetchAnnotations get annotations and its metadata
func (rs *RESTServices) FetchAnnotations(ctx *gin.Context) {
	segments, err := rs.repository.FetchAnnotations(ctx)
	if err != nil {
		JSON(ctx, "failed to fetch Annotation", http.StatusInternalServerError, err)
		return
	}
	JSON(ctx, "success", http.StatusOK, segments)
}

//CreateAnnotation create annotation and its metadata
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

// UpdateAnnotation updates annotation and it's metadata
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

// DeleteVideo delete a Video
func (rs *RESTServices) DeleteVideo(ctx *gin.Context) {
	value, err := rs.repository.DeleteVideo(ctx)
	if err != nil {
		JSON(ctx, "failed to delete video", http.StatusBadRequest, err)
		return
	}
	JSON(ctx, "success", http.StatusOK, value)
}

// DeleteAnnotationByID delete annotation by id
func (rs *RESTServices) DeleteAnnotationByID(ctx *gin.Context) {
	value, err := rs.repository.DeleteAnnotationByID(ctx)
	if err != nil {
		JSON(ctx, "failed to delete Annotation", http.StatusBadRequest, err)
		return
	}
	JSON(ctx, "success", http.StatusOK, value)
}

// CreateVideoFromDisk creates video schemas from form data
func (rs *RESTServices) CreateVideoFromDisk(ctx *gin.Context) {
	var annoSchema model.AnnotatedVideo
	videoFile, err := ctx.FormFile("videFile")
	if err != nil {
		JSON(ctx, "failed", http.StatusBadRequest, err)
		return
	}
	annoSchema.VideoLink = videoFile.Filename
	schemaFile, err := ctx.FormFile("schemaFile")
	if err != nil {
		JSON(ctx, "failed", http.StatusBadRequest, err)
		return
	}
	video, err := utility.ReadFileHeader(videoFile)

	if err != nil {
		JSON(ctx, "failed", http.StatusBadRequest, err)
		return
	}
	schema, err := utility.ReadFileHeader(schemaFile)
	if err := json.Unmarshal(schema, &annoSchema.Schema); err != nil {
		JSON(ctx, "failed", http.StatusBadRequest, err)
		return
	}
	vd, err := utility.GetVideoDetailsFromByte(video)
	if err != nil {
		JSON(ctx, "failed", http.StatusBadRequest, err)
		return
	}
	annoSchema.Duration = vd.Duration
	annoSchema.FrameRate = vd.Framerate
	savedSchema, err := rs.repository.Create(ctx, &annoSchema)
	if err != nil {
		JSON(ctx, "failed to save video details", http.StatusBadRequest, err)
		return
	}
	JSON(ctx, "success", http.StatusCreated, savedSchema)
}
