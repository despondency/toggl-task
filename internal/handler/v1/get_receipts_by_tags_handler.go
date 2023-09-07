package v1

import (
	"context"
	"fmt"
	"github.com/despondency/toggl-task/internal/model"
	"github.com/despondency/toggl-task/internal/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

type GetReceiptsByTagsResponse struct {
	Receipts    []*model.Receipt `json:"receipts"`
	TotalAmount float64          `json:"total_amount"`
}

type GetReceiptsByTagsResultHandler struct {
	receiptService service.ReceiptServicer
}

func NewGetReceiptsByTagResultHandler(uploadSvc service.ReceiptServicer) *GetReceiptResultHandler {
	return &GetReceiptResultHandler{receiptService: uploadSvc}
}

func (grrh *GetReceiptResultHandler) GetReceiptsByTagHandler() func(c *fiber.Ctx) error {
	handler := func(c *fiber.Ctx) error {
		queries := c.Query("tags")
		queries = strings.TrimSpace(queries)
		if len(queries) == 0 {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		tagsSeparated := strings.Split(queries, ",")
		if len(tagsSeparated) == 0 {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		if checkIfTagsEmpty(tagsSeparated) {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		r, err := grrh.receiptService.GetReceiptsByTags(context.Background(), tagsSeparated)
		if err != nil {
			c.Response().AppendBodyString(fmt.Sprintf("err: %v", err))
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		resp := &GetReceiptsByTagsResponse{
			Receipts: r,
		}
		var total float64 = 0
		for _, receipt := range r {
			if receipt.TotalAmount != nil {
				parsed, err := strconv.ParseFloat(*receipt.TotalAmount, 64)
				if err != nil {
					c.Response().AppendBodyString(fmt.Sprintf("failed parsing %s amount", *receipt.TotalAmount))
					return c.SendStatus(fiber.StatusInternalServerError)
				}
				total += parsed
			}
		}
		resp.TotalAmount = total
		err = c.JSON(resp)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusOK)
	}
	return handler
}

func checkIfTagsEmpty(separated []string) bool {
	for _, tag := range separated {
		if len(strings.TrimSpace(tag)) == 0 {
			return true
		}
	}
	return false
}
