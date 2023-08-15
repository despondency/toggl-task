package v1

import (
	"bytes"
	"context"
	"encoding/json"
	v1 "github.com/despondency/toggl-task/internal/handler/v1"
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
			uploadReceipt()
		})
	}
}

func uploadReceipt() string {
	url := "http://localhost:8084/v1/receipt"
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("../../testdata/receipt1.png")
	defer file.Close()
	part1, errFile1 := writer.CreateFormFile("file", filepath.Base("../../testdata/receipt1.png"))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		log.Error(errFile1)
		return ""
	}
	err := writer.Close()
	if err != nil {
		log.Error(err)
		return ""
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Error(err)
		return ""
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return ""
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return ""
	}
	uploadModelRes := &v1.UploadResultModel{}
	err = json.Unmarshal(body, uploadModelRes)
	if err != nil {
		log.Error(err)
		return ""
	}
	log.Infof("id is %s", string(body))
	return uploadModelRes.Id
}
