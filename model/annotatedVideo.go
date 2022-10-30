package model

type AnnotatedVideo struct {
	//The url of the video file for this dataset
	Id        string  `json:"id,omitempty"`
	VideoLink string  `json:"video_link" binding:"required"`
	Duration  float64 `json:"duration,omitempty"`
	FrameRate float64 `json:"frame_rate,omitempty"`
	//The path of the metadata schema file for this dataset
	Schema   MetadataSchema     `json:"schema" form:"schema_file" binding:"required"`
	Segments []AnnotatedSegment `json:"segments,omitempty" sql:"-"`
}
