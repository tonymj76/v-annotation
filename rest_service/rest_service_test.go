package rest_service

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMockRepository_Create(t *testing.T) {

	testCase := []struct {
		name          string
		videoLink     string
		schema        string
		buildStubs    func(mockStore *MockRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			videoLink: "https://cdn.api.video/vod/viqyDj9tTG4LmwTJuwf7LS9/mp4/source.mp4",
			buildStubs: func(mockStore *MockRepository) {
				mockStore.EXPECT().
					Create(gomock.Any(), gomock.Eq("")).
					Times(1).
					Return(nil, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := NewMockRepository(ctrl)
			tc.buildStubs(mock)

		})

	}
}
