package monitor

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	drift "github.com/cyber-life-tech/o2-train/internal/data-drift"
)

func TestNewDataDrift(t *testing.T) {
	t.Parallel()

	dd := NewDataDrift(os.Stdin, os.Stdout)

	assert.Equal(t, os.Stdin, dd.reader)
	assert.Equal(t, os.Stdout, dd.writer)
	assert.Nil(t, dd.hooks)
}

func TestDataDrift_Tick(t *testing.T) {
	t.Parallel()

	type testCase struct {
		readerData        string
		writer            *bytes.Buffer
		expectedError     error
		expectedErrorText string
	}

	testCases := map[string]testCase{
		"[happy] reader/no-writer": {
			readerData: "{}",
		},
		"[happy] data-drift/with-writer": {
			readerData: `{
  "ref_vector": [1],
  "new_vector": [],
  "metric": "ks",
  "threshold": 0.1
}`,
			writer: &bytes.Buffer{},
		},
		"[faulty] no-reader/no-writer": {
			expectedError: ErrTick,
		},
		"[faulty] corrupted-reader/no-writer": {
			readerData:        "%%%",
			expectedErrorText: "failed to json.Unmarshal data",
		},
		"[faulty] invalid-metric/no-writer": {
			readerData: `{
  "ref_vector": [1],
  "new_vector": [1],
  "metric": "not_a_valid_metric",
  "threshold": 0.1
}`,
			expectedError: drift.ErrInvalidMetric,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var reader io.Reader
			if tc.readerData != "" {
				reader = strings.NewReader(tc.readerData)
			}

			dd := NewDataDrift(reader, tc.writer)
			err := dd.Tick(context.Background())

			if tc.writer != nil {
				assert.NotEmpty(t, tc.writer.String())
			}

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else if tc.expectedErrorText != "" {
				assert.ErrorContains(t, err, tc.expectedErrorText)
			}
		})
	}
}
