package model

type AnnotatedVideo struct {
	Id string `json:"id,omitempty"`
	//The url of the video file for this dataset
	VideoLink string `json:"video_link" binding:"required"`
	//Duration of the video
	Duration float64 `json:"duration,omitempty"`
	// FrameRate of the vidoe
	FrameRate float64 `json:"frame_rate,omitempty"`
	//The path of the metadata schema file for this dataset
	Schema MetadataSchema `json:"schema" form:"schema_file" binding:"required"`
	// list of the annotations
	Annotations []AnnotatedSegment `json:"segments,omitempty" sql:"-"`
}
