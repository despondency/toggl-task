package v1

import (
	"bytes"
	"fmt"
	"github.com/despondency/toggl-task/internal/persister"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
)

type UploadFileHandler struct {
	persister persister.RawFilePersister
}

func NewUploadFileHandler(persister persister.RawFilePersister) *UploadFileHandler {
	return &UploadFileHandler{persister: persister}
}

func (ufh *UploadFileHandler) GetUploadFileHandler() func(c *fiber.Ctx) error {
	handler := func(c *fiber.Ctx) error {
		f, err := c.FormFile("persister")
		if err != nil {
			c.Response().AppendBodyString("can't find form element 'persister'")
			return c.SendStatus(fiber.StatusBadRequest)
		}
		openedFile, err := f.Open()
		if err != nil {
			c.Response().AppendBodyString("could not open persister")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, openedFile); err != nil {
			c.Response().AppendBodyString("could not copy persister")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		mimeType := http.DetectContentType(buf.Bytes())
		if mimeType != "image/png" && mimeType != "image/jpeg" && mimeType != "application/pdf" {
			c.Response().AppendBodyString(fmt.Sprintf("unknown mime type %s", mimeType))
			return c.SendStatus(fiber.StatusBadRequest)
		}
		err = ufh.persister.Persist(f.Filename, buf.Bytes())
		if err != nil {
			c.Response().AppendBodyString("could not persist persister")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusAccepted)
	}
	return handler
}
