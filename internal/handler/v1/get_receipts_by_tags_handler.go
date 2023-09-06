package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/despondency/toggl-task/internal/service"
	"github.com/gofiber/fiber/v2"
)

type GetReceiptsByTagsRequest struct {
	Tags []string `json:"tags"`
}

type GetReceiptsByTagsResultHandler struct {
	receiptService service.ReceiptServicer
}

func NewGetReceiptsByTagResultHandler(uploadSvc service.ReceiptServicer) *GetReceiptResultHandler {
	return &GetReceiptResultHandler{receiptService: uploadSvc}
}

func (grrh *GetReceiptResultHandler) GetReceiptsByTagHandler() func(c *fiber.Ctx) error {
	handler := func(c *fiber.Ctx) error {
		body := c.Body()
		req := &GetReceiptsByTagsRequest{}
		err := json.Unmarshal(body, req)
		if err != nil {
			c.Response().AppendBodyString(fmt.Sprintf("err: %v", err))
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		r, err := grrh.receiptService.GetReceiptsByTags(context.Background(), req.Tags)
		if err != nil {
			c.Response().AppendBodyString(fmt.Sprintf("err: %v", err))
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		jsonStr, err := json.Marshal(r)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendString(string(jsonStr))
	}
	return handler
}
