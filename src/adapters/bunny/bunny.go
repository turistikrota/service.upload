package bunny

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/turistikrota/service.upload/src/domain/cdn"
)

type repo struct {
	config Config
}

type Config struct {
	CdnHost     string
	UploadHost  string
	StorageZone string
	ApiKey      string
}

func New(cnf Config) cdn.Repository {
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
	req.Header.Add("content-type", "application/octet-stream")
	req.Header.Add("AccessKey", r.config.ApiKey)
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
	return fmt.Sprintf("%s/%s/%s/%s", r.config.UploadHost, r.config.StorageZone, path, filename)
}

func (r *repo) makeUploadedUrl(filename string, path string) string {
	return fmt.Sprintf("%s/%s/%s", r.config.CdnHost, path, filename)
}
