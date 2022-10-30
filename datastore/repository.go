package datastore

import (
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tonymj76/video-annotator/model"
)

type Repository interface {
	Create(*gin.Context, *model.AnnotatedVideo) (*model.AnnotatedVideo, error)
	FetchVideos(*gin.Context) ([]*model.AnnotatedVideo, error)
	FetchVideoByID(*gin.Context) (*model.AnnotatedVideo, error)
	FetchSegments(ctx *gin.Context) ([]*model.AnnotatedSegment, error)
	FetchSegmentByID(*gin.Context) (*model.AnnotatedSegment, error)
	CreateSegmentByID(*gin.Context, *model.AnnotationSegments) (string, error)
	DeleteVideo(ctx *gin.Context) (string, error)
	DeleteSegmentByID(ctx *gin.Context) (string, error)
}

func (d *DBStore) Create(ctx *gin.Context, video *model.AnnotatedVideo) (*model.AnnotatedVideo, error) {
	row := d.SQLBuilder.Insert("video_schema").SetMap(map[string]any{
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
		"video_schema",
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
		"video_schema",
		"id",
		Id,
		"id", "duration", "frame_rate", "video", "schema")
	return scanRows(row)
}

func (d *DBStore) FetchSegmentByID(ctx *gin.Context) (*model.AnnotatedSegment, error) {
	segmentID := ctx.Param("segmentID")
	row := d.fetchOneFromDB(ctx,
		"segment",
		"id",
		segmentID,
		"id", "video_schema_id", "start", "end_time", "metadata")
	return scanSegmentRows(row)
}

func (d *DBStore) FetchSegments(ctx *gin.Context) ([]*model.AnnotatedSegment, error) {
	Id := ctx.Param("id")

	var segments []*model.AnnotatedSegment
	rows, err := d.fetchMoreFromDB(ctx,
		"segment",
		"video_schema_id",
		Id,
		"id", "video_schema_id", "start", "end_time", "metadata")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		sr, err := scanSegmentRows(rows)
		if err != nil {
			return nil, err
		}
		segments = append(segments, sr)
	}
	return segments, nil
}

func (d *DBStore) CreateSegmentByID(ctx *gin.Context, segments *model.AnnotationSegments) (string, error) {
	videoSchemaID := ctx.Param("id")
	var ids []string
	for _, segment := range segments.Segments {
		logrus.Println(segment.End)
		row := d.SQLBuilder.Insert("segment").
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

	return fmt.Sprintf("%v segment(s) created", len(ids)), nil
}

func (d *DBStore) DeleteVideo(ctx *gin.Context) (string, error) {
	Id := ctx.Param("id")
	if Id == "0" {
		return "", errors.New("invalid id")
	}
	_, err := d.SQLBuilder.Delete(
		"video_schema",
	).Where(squirrel.Eq{"id": Id}).ExecContext(ctx)

	return "success", err
}

func (d *DBStore) DeleteSegmentByID(ctx *gin.Context) (string, error) {
	Id := ctx.Param("id")
	if Id == "0" {
		return "", errors.New("invalid id")
	}
	_, err := d.SQLBuilder.Delete(
		"segment",
	).Where(squirrel.Eq{"id": Id}).ExecContext(ctx)

	return "success", err
}
