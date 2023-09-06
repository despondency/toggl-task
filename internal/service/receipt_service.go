package service

import (
	"context"
	"github.com/despondency/toggl-task/internal/persister"
	"github.com/despondency/toggl-task/internal/scanner"
)

type UploadReceiptBody struct {
	Tags []string `json:"tags"`
}

type UploadPayload struct {
	*UploadReceiptBody
	FilePayload
}

type FilePayload struct {
	Receipt  []byte
	FileName string
	MimeType string
}

type ReceiptServicer interface {
	CreateReceipt(ctx context.Context, payload UploadPayload) (string, error)
	GetReceipt(ctx context.Context, uuid string) (*persister.ResultModel, error)
}

type MultiServicer struct {
	rawFilePersister persister.RawFilePersister
	resultPersister  persister.ResultPersister
	scanner          scanner.Scanner
}

func NewMultiServicer(rawFilePersister persister.RawFilePersister, resultPersister persister.ResultPersister,
	scanner scanner.Scanner) ReceiptServicer {
	return &MultiServicer{
		rawFilePersister: rawFilePersister,
		resultPersister:  resultPersister,
		scanner:          scanner,
	}
}

func (ms *MultiServicer) GetReceipt(ctx context.Context, uuid string) (*persister.ResultModel, error) {
	return ms.resultPersister.Get(ctx, uuid)
}

func (ms *MultiServicer) CreateReceipt(ctx context.Context, payload UploadPayload) (string, error) {
	err := ms.rawFilePersister.Persist(payload.FileName, payload.Receipt)
	if err != nil {
		return "", err
	}
	res, err := ms.scanner.Scan(ctx, payload.Receipt, payload.MimeType)
	if err != nil {
		return "", err
	}
	return ms.resultPersister.Persist(ctx, &persister.ResultModel{
		Payload: res.Result,
		Tags:    payload.Tags,
	})
}
