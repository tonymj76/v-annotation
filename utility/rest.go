package utility

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

var ErrEmptyBody = errors.New("body must not be empty")

//ReadJSON parse request body to json and check for errors
func ReadJSON(c *gin.Context, dst any) error {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	if err := c.ShouldBindJSON(dst); err != nil {
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return ErrEmptyBody

		default:
			return err
		}

	}

	return nil
}

// GoClient to make http request
func GoClient(urlPath string) (*http.Response, error) {
	var client http.Client
	client = http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(req)
	return response, err
}
