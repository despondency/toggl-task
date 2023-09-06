package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	v1 "github.com/despondency/toggl-task/internal/handler/v1"
	"github.com/despondency/toggl-task/internal/service"
	"github.com/gofiber/fiber/v2/log"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestIntegration_UploadReceipt(t *testing.T) {

	testCases := []struct {
		ctx  context.Context
		name string
	}{
		{
			name: "upload receipt1.png",
			ctx:  context.Background(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uploadReceipt("receipt1.png", []string{})
		})
	}
}

func uploadReceipt(receiptNameToUpload string, tags []string) string {
	url := "http://localhost:8084/v1/receipt"
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open(fmt.Sprintf("../../testdata/%s", receiptNameToUpload))
	defer file.Close()
	part1, errFile1 := writer.CreateFormFile("file", filepath.Base(fmt.Sprintf("../../testdata/%s", receiptNameToUpload)))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		log.Fatal(errFile1)
		return ""
	}
	if len(tags) > 0 {
		wr, err := writer.CreateFormField("json")
		if err != nil {
			log.Fatal(err)
		}
		jsonTags := &service.UploadReceiptBody{Tags: tags}
		b, err := json.Marshal(jsonTags)
		if err != nil {
			log.Fatal(err)
		}
		_, err = wr.Write(b)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := writer.Close()
	if err != nil {
		log.Fatal(err)
		return ""
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Fatal(err)
		return ""
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	uploadModelRes := &v1.UploadResultModel{}
	err = json.Unmarshal(body, uploadModelRes)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	log.Infof("id is %s", string(body))
	return uploadModelRes.Id
}
