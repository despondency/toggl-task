package service

import (
	"context"
	"github.com/despondency/toggl-task/internal/model"
	"github.com/despondency/toggl-task/internal/persister"
	"github.com/despondency/toggl-task/internal/scanner"
)

type UploadReceiptBody struct {
	Tags []string `json:"tags"`
}

type UploadPayload struct {
	*UploadReceiptBody
	*FilePayload
}

type FilePayload struct {
	Receipt  []byte
	FileName string
	MimeType string
}

type ReceiptServicer interface {
	CreateReceipt(ctx context.Context, payload *UploadPayload) (*model.Receipt, error)
	GetReceipt(ctx context.Context, uuid string) (*model.Receipt, error)
	GetReceiptsByTags(ctx context.Context, tags []string) ([]*model.Receipt, error)
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

func (ms *MultiServicer) GetReceipt(ctx context.Context, uuid string) (*model.Receipt, error) {
	return ms.resultPersister.Get(ctx, uuid)
}

func (ms *MultiServicer) GetReceiptsByTags(ctx context.Context, tags []string) ([]*model.Receipt, error) {
	return ms.resultPersister.GetByTags(ctx, tags)
}

func (ms *MultiServicer) CreateReceipt(ctx context.Context, payload *UploadPayload) (*model.Receipt, error) {
	err := ms.rawFilePersister.Persist(payload.FileName, payload.Receipt)
	if err != nil {
		return nil, err
	}
	res, err := ms.scanner.Scan(ctx, payload.Receipt, payload.MimeType)
	if err != nil {
		return nil, err
	}
	res.Tags = payload.Tags
	return ms.resultPersister.Persist(ctx, res)
}
