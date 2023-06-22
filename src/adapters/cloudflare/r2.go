package cloudflare

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

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
	req.Header.Add("content-type", "image/jpeg")
	req.Header.Add("x-amz-content-sha256", "UNSIGNED-PAYLOAD")
	req.Header.Add("x-amz-date", "20210901T000000Z")
	req.Header.Add("x-amz-acl", "public-read")
	req.Header.Add("Host", "upload.api.turistikrota.com")
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

func (r *repo) makeUrl(filename string, path string) string {
	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s/%s", r.config.AccountId, r.config.Bucket, path, filename)
}

func (r *repo) makeUploadedUrl(filename string, path string) string {
	return fmt.Sprintf("%s/%s/%s", r.config.PublicHost, path, filename)
}

// 1f5d3866562e9650cc4da9557aacc24f accountid
// 2bc7697469760ccbd5f6e3a147f4d3b8 access key id
// d126eb10512ac98c94d47cab5bbcadc98ab169f81437de53e2e1978246a248ce secret key
