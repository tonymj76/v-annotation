package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

var ErrDataType = errors.New("not a valid metadata data type")

type SchemaField struct {
	// Each field must have a unique name (required)
	Name string `json:"name" binding:"required"`
	// Each field must have a datatype (supported types are boolean, float, int, and string)
	Type string `json:"type" binding:"required"`

	// possible valid values for a field (optional)
	Values []string `json:"values"`
}

// SchemaRoot definition for the root node of a parsed metadata object
type SchemaRoot struct {
	// segments that list the metadata fields that apply to entire segments of a video required
	Segments []SchemaField `json:"segments" binding:"required"`
}

type MetadataSchema struct {
	// schema root
	Fields SchemaRoot `json:"fields" binding:"required"`
}

func (m MetadataSchema) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *MetadataSchema) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &m)
}
