package datastore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	"github.com/tonymj76/video-annotator/model"
	"strings"
)

func (d *DBStore) fetchOneFromDB(ctx context.Context, table, whereToSearch, itemToSearch string, art ...string) squirrel.RowScanner {
	switch itemToSearch {
	case "all":
		row := d.SQLBuilder.Select("*").From(table).QueryRowContext(ctx)
		return row
	case "partial":
		row := d.SQLBuilder.Select(strings.Join(art, ",")).From(table).QueryRowContext(ctx)
		return row
	default:
		row := d.SQLBuilder.Select(strings.Join(art, ",")).
			From(table).
			Where(fmt.Sprintf("%s = ?", whereToSearch), itemToSearch).
			QueryRowContext(ctx)
		return row
	}
}

func (d *DBStore) fetchMoreFromDB(ctx context.Context, table, whereToSearch, itemToSearch string, art ...string) (*sql.Rows, error) {
	switch itemToSearch {
	case "all":
		rows, err := d.SQLBuilder.Select("*").From(table).QueryContext(ctx)
		return rows, err
	case "partial":
		rows, err := d.SQLBuilder.Select(strings.Join(art, ",")).From(table).QueryContext(ctx)
		return rows, err
	default:
		rows, err := d.SQLBuilder.Select(strings.Join(art, ",")).
			From(table).
			Where(fmt.Sprintf("%s = ?", whereToSearch), itemToSearch).
			QueryContext(ctx)
		return rows, err
	}
}

func scanRows(row squirrel.RowScanner) (*model.AnnotatedVideo, error) {
	var an model.AnnotatedVideo
	err := row.Scan(
		&an.Id,
		&an.Duration,
		&an.FrameRate,
		&an.VideoLink,
		&an.Schema,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.Errorf("no record found: %v", err.Error())
			return nil, err
		}
		return nil, err
	}
	return &an, nil
}

func scanAnnotationRows(row squirrel.RowScanner) (*model.AnnotatedSegment, error) {
	var as model.AnnotatedSegment
	err := row.Scan(
		&as.Id,
		&as.VideoSchemaId,
		&as.Start,
		&as.End,
		&as.Metadata,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.Errorf("no record found: %v", err.Error())
			return nil, err
		}
		return nil, err
	}
	return &as, nil
}
