package datastore

import (
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tonymj76/video-annotator/model"
)

const (
	annotation  = "annotation"
	videoSchema = "video_schema"
)

type Repository interface {
	Create(*gin.Context, *model.AnnotatedVideo) (*model.AnnotatedVideo, error)
	FetchVideos(*gin.Context) ([]*model.AnnotatedVideo, error)
	FetchVideoByID(*gin.Context) (*model.AnnotatedVideo, error)
	FetchAnnotations(ctx *gin.Context) ([]*model.AnnotatedSegment, error)
	FetchAnnotationByID(*gin.Context) (*model.AnnotatedSegment, error)
	CreateAnnotationByID(*gin.Context, *model.AnnotationSegments) (string, error)
	DeleteVideo(ctx *gin.Context) (string, error)
	DeleteAnnotationByID(ctx *gin.Context) (string, error)
}

func (d *DBStore) Create(ctx *gin.Context, video *model.AnnotatedVideo) (*model.AnnotatedVideo, error) {
	row := d.SQLBuilder.Insert(videoSchema).SetMap(map[string]any{
		"video":      video.VideoLink,
		"schema":     video.Schema,
		"duration":   video.Duration,
		"frame_rate": video.FrameRate,
	}).Suffix(`RETURNING "id"`).QueryRowContext(ctx)

	var Id string
	if err := row.Scan(&Id); err != nil {
		return nil, err
	}
	video.Id = Id
	return video, nil
}

func (d *DBStore) FetchVideos(ctx *gin.Context) ([]*model.AnnotatedVideo, error) {
	var videos []*model.AnnotatedVideo
	rows, err := d.fetchMoreFromDB(ctx,
		videoSchema,
		"",
		"partial",
		"id", "duration", "frame_rate", "video", "schema")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		v, err := scanRows(rows)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}
	return videos, nil
}

func (d *DBStore) FetchVideoByID(ctx *gin.Context) (*model.AnnotatedVideo, error) {
	Id := ctx.Param("id")
	if Id == "0" {
		return nil, errors.New("invalid id")
	}
	row := d.fetchOneFromDB(ctx,
		videoSchema,
		"id",
		Id,
		"id", "duration", "frame_rate", "video", "schema")
	return scanRows(row)
}

func (d *DBStore) FetchAnnotationByID(ctx *gin.Context) (*model.AnnotatedSegment, error) {
	annotationID := ctx.Param("annotationID")
	row := d.fetchOneFromDB(ctx,
		annotation,
		"id",
		annotationID,
		"id", "video_schema_id", "start", "end_time", "metadata")
	return scanAnnotationRows(row)
}

func (d *DBStore) FetchAnnotations(ctx *gin.Context) ([]*model.AnnotatedSegment, error) {
	Id := ctx.Param("id")

	var segments []*model.AnnotatedSegment
	rows, err := d.fetchMoreFromDB(ctx,
		annotation,
		"video_schema_id",
		Id,
		"id", "video_schema_id", "start", "end_time", "metadata")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		sr, err := scanAnnotationRows(rows)
		if err != nil {
			return nil, err
		}
		segments = append(segments, sr)
	}
	return segments, nil
}

func (d *DBStore) CreateAnnotationByID(ctx *gin.Context, segments *model.AnnotationSegments) (string, error) {
	videoSchemaID := ctx.Param("id")
	var ids []string
	for _, segment := range segments.Annotations {
		logrus.Println(segment.End)
		row := d.SQLBuilder.Insert(annotation).
			SetMap(map[string]any{
				"video_schema_id": videoSchemaID,
				"start":           segment.Start,
				"end_time":        segment.End,
				"metadata":        segment.Metadata,
			}).Suffix(`RETURNING "id"`).QueryRowContext(ctx)
		var Id string
		if err := row.Scan(&Id); err != nil {
			return "", err
		}
		ids = append(ids, Id)
	}

	return fmt.Sprintf("%v Annotation(s) created", len(ids)), nil
}

func (d *DBStore) DeleteVideo(ctx *gin.Context) (string, error) {
	Id := ctx.Param("id")
	if Id == "0" {
		return "", errors.New("invalid id")
	}
	_, err := d.SQLBuilder.Delete(
		videoSchema,
	).Where(squirrel.Eq{"id": Id}).ExecContext(ctx)

	return "success", err
}

func (d *DBStore) DeleteAnnotationByID(ctx *gin.Context) (string, error) {
	Id := ctx.Param("id")
	if Id == "0" {
		return "", errors.New("invalid id")
	}
	_, err := d.SQLBuilder.Delete(
		annotation,
	).Where(squirrel.Eq{"id": Id}).ExecContext(ctx)

	return "success", err
}
