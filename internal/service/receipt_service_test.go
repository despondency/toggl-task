package service_test

import (
	"context"
	"github.com/despondency/toggl-task/internal/model"
	"github.com/despondency/toggl-task/internal/service"
	persistermock "github.com/despondency/toggl-task/mocks/persister"
	scannermock "github.com/despondency/toggl-task/mocks/scanner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestUnitMultiServicer_CreateReceipt(t *testing.T) {
	type testCase struct {
		ctx            context.Context
		name           string
		fileName       string
		fileContent    []byte
		mimeType       string
		errExpected    string
		idExpected     primitive.ObjectID
		createInstance func(ctx context.Context, tc *testCase, t *testing.T) service.ReceiptServicer
	}
	testCases := []*testCase{
		{
			name:        "first case",
			ctx:         context.Background(),
			idExpected:  primitive.NewObjectID(),
			fileContent: []byte{},
			createInstance: func(ctx context.Context, tc *testCase, t *testing.T) service.ReceiptServicer {
				rawp := persistermock.NewRawFilePersister(t)
				rawp.EXPECT().Persist("", []byte{}).Return(nil)
				rp := persistermock.NewResultPersister(t)
				rp.EXPECT().Persist(ctx, &model.Receipt{}).Return(&model.Receipt{Id: tc.idExpected}, nil)
				s := scannermock.NewScanner(t)
				s.EXPECT().Scan(ctx, []byte{}, mock.AnythingOfType("string")).Return(&model.Receipt{}, nil)
				return service.NewMultiServicer(rawp, rp, s)
			},
		},
	}

	for _, tc := range testCases {
		instance := tc.createInstance(tc.ctx, tc, t)
		id, err := instance.CreateReceipt(tc.ctx, &service.UploadPayload{
			UploadReceiptBody: &service.UploadReceiptBody{},
			FilePayload: &service.FilePayload{
				Receipt:  tc.fileContent,
				FileName: tc.fileName,
				MimeType: tc.mimeType,
			},
		})
		assert.Equal(t, id.Id, tc.idExpected)
		if tc.errExpected != "" && err != nil {
			assert.EqualErrorf(t, err, tc.errExpected, "", err)
		} else {
			require.NoError(t, err)
		}
	}
}
