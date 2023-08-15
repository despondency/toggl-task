package v1

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"io/ioutil"
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
			getReceipt(uploadReceipt())
		})
	}
}

func getReceipt(id string) {
	url := fmt.Sprintf("http://localhost:8084/v1/receipt?id=%s", id)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Infof("successfully retrieved receipt %s", body)
}
