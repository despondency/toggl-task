package persister

import (
	"fmt"
	"os"
)

type RawFilePersister interface {
	Persist(fileName string, fileContent []byte) error
}

type Local struct {
	folderToWriteTo string
}

func NewLocal(folderToWriteTo string) RawFilePersister {
	return &Local{
		folderToWriteTo: folderToWriteTo,
	}
}

func (l *Local) Persist(fileName string, fileContent []byte) error {
	return os.WriteFile(fmt.Sprintf("%s/%s", l.folderToWriteTo, fileName), fileContent, 0644)
}
