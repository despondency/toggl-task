package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/despondency/toggl-task/internal/service"
	"github.com/gofiber/fiber/v2"
)

type GetReceiptResultHandler struct {
	receiptService service.ReceiptServicer
}

func NewGetReceiptResultHandler(uploadSvc service.ReceiptServicer) *GetReceiptResultHandler {
	return &GetReceiptResultHandler{receiptService: uploadSvc}
}

func (grrh *GetReceiptResultHandler) GetReceiptHandler() func(c *fiber.Ctx) error {
	handler := func(c *fiber.Ctx) error {
		queryValue := c.Query("id")
		r, err := grrh.receiptService.GetReceipt(context.Background(), queryValue)
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
