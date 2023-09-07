package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/despondency/toggl-task/internal/persister"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestIntegration_GetReceiptsByTags(t *testing.T) {
	testCases := []struct {
		ctx                  context.Context
		receiptNamesToUpload []string
		expectedReceipts     int
		tags                 [][]string
		name                 string
	}{
		{
			name:                 "upload receipt1.png 3 times with different tags and fetch",
			receiptNamesToUpload: []string{"receipt1.png", "receipt1.png", "receipt1.png", "receipt1.png"},
			tags: [][]string{
				{"tag1", "tag2"},
				{"tag1"},
				{"tag2"},
				{"tag4"},
			},
			expectedReceipts: 1,
			ctx:              context.Background(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for i, receiptNameToUpload := range tc.receiptNamesToUpload {
				uploadReceipt(receiptNameToUpload, tc.tags[i])
			}
			assert.Equal(t, 1, len(getReceiptsByTags([]string{"tag4"})))
			assert.Equal(t, 2, len(getReceiptsByTags([]string{"tag1"})))
			assert.Equal(t, 2, len(getReceiptsByTags([]string{"tag2"})))
			assert.Equal(t, 1, len(getReceiptsByTags([]string{"tag1", "tag2"})))
		})
	}
}

func getReceiptsByTags(tags []string) []*persister.ResultModel {
	url := fmt.Sprintf("http://localhost:8084/v1/receipts-by-tags?tags=%s", strings.Join(tags, ","))
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Fatal(err)
		return nil
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	resp := make([]*persister.ResultModel, 0)
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return resp
}
