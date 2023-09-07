package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/despondency/toggl-task/internal/service"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
)

type UploadReceiptHandler struct {
	receiptSvc service.ReceiptServicer
}

func NewUploadReceiptHandler(uploadSvc service.ReceiptServicer) *UploadReceiptHandler {
	return &UploadReceiptHandler{receiptSvc: uploadSvc}
}

func (urh *UploadReceiptHandler) GetUploadFileHandler() func(c *fiber.Ctx) error {
	handler := func(c *fiber.Ctx) error {
		f, err := c.FormFile("file")
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
		jsonValue := c.FormValue("json")
		uploadReceiptBody := &service.UploadReceiptBody{}
		if jsonValue != "" {
			err = json.Unmarshal([]byte(jsonValue), uploadReceiptBody)
			if err != nil {
				c.Response().AppendBodyString(fmt.Sprintf("failed to extract body"))
				return c.SendStatus(fiber.StatusBadRequest)
			}
		}
		receiptModel, err := urh.receiptSvc.CreateReceipt(context.Background(), &service.UploadPayload{
			UploadReceiptBody: uploadReceiptBody,
			FilePayload: &service.FilePayload{
				Receipt:  buf.Bytes(),
				FileName: f.Filename,
				MimeType: mimeType,
			},
		})
		if err != nil {
			c.Response().AppendBodyString(fmt.Sprintf("could not persist %v", err))
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		err = c.JSON(receiptModel)
		if err != nil {
			c.Response().AppendBodyString("cannot add result uuid to body")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusAccepted)
	}
	return handler
}
