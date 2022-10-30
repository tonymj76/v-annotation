package utility

import (
	"errors"
	"fmt"
	"github.com/tonymj76/video-annotator/model"
	"reflect"
)

//ValidateManyData validate input request
func ValidateManyData(vs *model.AnnotatedVideo, sg model.AnnotationSegments) error {
	duration := vs.Duration
	schema := vs.Schema
	for _, segment := range sg.Annotations {
		if err := checkStartAndEndTime(duration, segment.Start, segment.End); err != nil {
			return err
		}

		if err := checkMetaData(segment.Metadata, schema.Fields.Segments); err != nil {
			return err
		}
	}
	return nil
}

//ValidateOneData validate input request
func ValidateOneData(vs *model.AnnotatedVideo, ag model.AnnotatedSegment) error {
	duration := vs.Duration
	schema := vs.Schema
	if err := checkStartAndEndTime(duration, ag.Start, ag.End); err != nil {
		return err
	}

	if err := checkMetaData(ag.Metadata, schema.Fields.Segments); err != nil {
		return err
	}
	return nil
}

// start and end time in seconds should not be less than or greater than video duration
func checkStartAndEndTime(d, s, e float64) error {
	if s < 0 {
		return errors.New("start time can't be negative number")
	}
	if e <= s || e > d {
		return errors.New(fmt.Sprintf("end time (%v) is invalid. duration is %v", e, d))
	}
	return nil
}

// check if the metadata of the schema is observed
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
