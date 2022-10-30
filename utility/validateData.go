package utility

import (
	"errors"
	"fmt"
	"github.com/tonymj76/video-annotator/model"
	"reflect"
)

func ValidateData(vs *model.AnnotatedVideo, sg model.AnnotationSegments) error {
	duration := vs.Duration
	schema := vs.Schema
	for _, segment := range sg.Segments {
		if err := checkStartAndEndSec(duration, segment.Start, segment.End); err != nil {
			return err
		}

		if err := checkMetaData(segment.Metadata, schema.Fields.Segments); err != nil {
			return err
		}
	}
	return nil
}

func checkStartAndEndSec(d, s, e float64) error {
	if s < 0 {
		return errors.New("start time can't be negative number")
	}
	if e <= s || e > d {
		return errors.New(fmt.Sprintf("end time (%v) is invalid. duration is %v", e, d))
	}
	return nil
}

func checkMetaData(m map[string]any, fields []model.SchemaField) error {
	for _, field := range fields {
		value, ok := m[field.Name]
		if !ok {
			return errors.New(fmt.Sprintf("invalid field in metadata: %v is not found in meta data", field.Name))
		}

		if reflect.TypeOf(value).String() != field.Type {
			typeOfV := reflect.TypeOf(value).String()
			if !(typeOfV == "float64" && field.Type == "int") {
				return errors.New(fmt.Sprintf("invalid type in metadata: got %v want %v at %v", typeOfV, field.Type, field.Name))
			}
		}

	}
	return nil
}
