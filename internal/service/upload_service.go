package service

import "github.com/despondency/toggl-task/internal/persister"

type UploadServicer interface {
	HandleUpload(fileName string, fileContent []byte) (string, error)
}

type RemoteUploader struct {
	rawFilePersister persister.RawFilePersister

	resultPersister persister.ResultPersister
}

func (ru *RemoteUploader) HandleUpload(fileName string, fileContent []byte) (string, error) {
	err := ru.rawFilePersister.Persist(fileName, fileContent)
	if err != nil {
		return "", err
	}
	return "", nil
}
