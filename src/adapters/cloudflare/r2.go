package cloudflare

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"api.turistikrota.com/upload/src/domain/cdn"
	"github.com/sirupsen/logrus"
)

type repo struct {
	config Config
}

type Config struct {
	AccountId  string
	AccessKey  string
	SecretKey  string
	Bucket     string
	PublicHost string
}

func NewR2(cnf Config) cdn.Repository {
	return &repo{
		config: cnf,
	}
}

func (r *repo) Upload(file []byte, filename string, path ...string) (string, bool) {
	p := "img"
	if len(path) > 0 {
		p = path[0]
	}
	payload := bytes.NewReader(file)
	req, _err := http.NewRequest("PUT", r.makeUrl(filename, p), payload)
	req.Header.Add("content-type", "image/png")
	req.Header.Add("x-amz-content-sha256", "UNSIGNED-PAYLOAD")
	req.Header.Add("Authorization", r.getAuthKey())
	if _err != nil {
		logrus.Error(_err)
		return "", false
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error(err)
		return "", false
	}
	if res.StatusCode == 201 {
		return r.makeUploadedUrl(filename, p), true
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	logrus.Error(errors.New(string(body)))
	return "", false
}

func (r *repo) getAuthKey() string {
	return fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/20230622/auto/s3/aws4_request, SignedHeaders=content-length;content-type;host;x-amz-content-sha256;x-amz-date, Signature=71f89d84ada33d0cb01c7ef702997abb6c9b6c940bcb91f01ba0343144041431", r.config.AccessKey)
}

func (r *repo) getCurrentTime() string {
	return time.Now().Format("20210901T000000Z")
}

func (r *repo) makeUrl(filename string, path string) string {
	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s/%s", r.config.AccountId, r.config.Bucket, path, filename)
}

func (r *repo) makeUploadedUrl(filename string, path string) string {
	return fmt.Sprintf("%s/%s/%s", r.config.PublicHost, path, filename)
}
