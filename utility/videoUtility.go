package utility

import (
	"errors"
	"fmt"
	mediainfo "github.com/zhulik/go_mediainfo"
	"io/ioutil"
	"strconv"
)

// VideoDetails Represents key metadata details about a video track
type VideoDetails struct {
	//The framerate of the video (in frames per second)
	Framerate float64

	//The duration of the video (in seconds)
	Duration float64
}

// NewVideoDetails Retrieve the metadata for the specified video file
func NewVideoDetails(videoLink string) (*VideoDetails, error) {
	res, err := GoClient(videoLink)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	mi := mediainfo.NewMediaInfo()
	err = mi.OpenMemory(bytes)

	if err != nil {
		return nil, err
	}
	framerate, err := CalcDurationStringToSeconds(mi.Get("FrameRate"), "s")
	duration, err := CalcDurationStringToSeconds(mi.Get("Duration"), "ms")
	return &VideoDetails{
		Framerate: framerate,
		Duration:  duration,
	}, nil
}

func CalcDurationStringToSeconds(t, unit string) (float64, error) {

	switch unit {
	case "ms":
		msec, err := strconv.ParseFloat(t, 64)
		sec := msec / 1000.0
		return sec, err
	case "s":
		duration, err := strconv.ParseFloat(t, 64)
		return duration, err
	default:
		return 0, errors.New(fmt.Sprintf("unknown unit %s in duration", unit))
	}

}
