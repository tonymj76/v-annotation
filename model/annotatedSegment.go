package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type AnnotatedSegment struct {
	// The starting position of the segment (in seconds)
	Id int `json:",omitempty"`

	VideoSchemaId int `json:"video_schema_id,omitempty" binding:"required"`

	Start float64 `json:"start" binding:"required"`

	// The ending position of the segment (in seconds)
	End float64 `json:"end" binding:"required"`

	//The metadata for the segment
	Metadata MetaData `json:"metadata"`
}

type AnnotationSegments struct {
	Annotations []AnnotatedSegment `json:"annotations" binding:"required"`
}

type MetaData map[string]any

func (m MetaData) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *MetaData) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &m)
}
