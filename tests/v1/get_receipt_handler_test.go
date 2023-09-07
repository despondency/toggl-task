package v1

import (
	"context"
	"fmt"
	"github.com/despondency/toggl-task/internal/model"
	"github.com/gofiber/fiber/v2/log"
	"io"
	"net/http"
	"testing"
)

func TestIntegration_GetReceipt(t *testing.T) {

	testCases := []struct {
		ctx  context.Context
		name string
	}{
		{
			name: "get receipt1.png scan results",
			ctx:  context.Background(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			getReceipt(uploadReceipt("receipt1.png", []string{}))
		})
	}
}

func getReceipt(receipt *model.Receipt) {
	url := fmt.Sprintf("http://localhost:8084/v1/receipt?id=%s", receipt.Id.Hex())
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Fatal(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Infof("successfully retrieved receipt %s", body)
}
