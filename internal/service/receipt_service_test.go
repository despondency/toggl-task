package service_test

import (
	"context"
	"github.com/despondency/toggl-task/internal/persister"
	"github.com/despondency/toggl-task/internal/scanner"
	"github.com/despondency/toggl-task/internal/service"
	persistermock "github.com/despondency/toggl-task/mocks/persister"
	scannermock "github.com/despondency/toggl-task/mocks/scanner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnitMultiServicer_CreateReceipt(t *testing.T) {
	testCases := []struct {
		ctx            context.Context
		name           string
		fileName       string
		fileContent    []byte
		mimeType       string
		errExpected    string
		idExpected     string
		createInstance func(ctx context.Context, t *testing.T) service.ReceiptServicer
	}{
		{
			name:        "first case",
			ctx:         context.Background(),
			idExpected:  "1234",
			fileContent: []byte{},
			createInstance: func(ctx context.Context, t *testing.T) service.ReceiptServicer {
				rawp := persistermock.NewRawFilePersister(t)
				rawp.EXPECT().Persist("", []byte{}).Return(nil)
				rp := persistermock.NewResultPersister(t)
				rp.EXPECT().Persist(ctx, &persister.ResultModel{UUID: "", Payload: "res"}).Return("1234", nil)
				s := scannermock.NewScanner(t)
				s.EXPECT().Scan(ctx, []byte{}, mock.AnythingOfType("string")).Return(&scanner.ScannedResult{Result: "res"}, nil)
				return service.NewMultiServicer(rawp, rp, s)
			},
		},
	}

	for _, tc := range testCases {
		instance := tc.createInstance(tc.ctx, t)
		id, err := instance.CreateReceipt(tc.ctx, tc.fileName, tc.fileContent, tc.mimeType)
		assert.Equal(t, id, tc.idExpected)
		if tc.errExpected != "" && err != nil {
			assert.EqualErrorf(t, err, tc.errExpected, "", err)
		} else {
			require.NoError(t, err)
		}
	}
}
