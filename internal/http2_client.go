package internal

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mongche/gupload/config"
	"golang.org/x/net/http2"
)

type Http2Client struct {
	address string
	client  *http.Client
}

// Create a new HTTP/2.0 client that provides to upload a file
func NewHttp2Client(cfg config.Http2ClientConfig) (c Http2Client, err error) {
	if cfg.Address == "" {
		err = errors.New("address must be required")
		return
	}

	c.address = cfg.Address

	if cfg.Cert == "" {
		err = errors.New("cert must be required")
		return
	}

	cert, err := ioutil.ReadFile(cfg.Cert)

	if err != nil {
		err = fmt.Errorf("failed to read root certificate, %v", err)
		return
	}

	certPool := x509.NewCertPool()

	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		err = errors.New("failed to root certificate to pool")
		return
	}

	c.client = &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}

	return
}

// Upload a file to the server
func (c *Http2Client) UploadFile(ctx context.Context, file string) (err error) {
	var (
		f *os.File
	)

	if f, err = os.Open(file); err != nil {
		err = fmt.Errorf("faield to open file, %s", file)
		return
	}
	defer f.Close()

	req, err := http.NewRequest("POST", c.address+"/upload", f)

	if err != nil {
		err = errors.New("failed to create the request")
		return
	}

	resp, err := c.client.Do(req)

	if err != nil {
		err = fmt.Errorf("failed to send a request to upload file, %v", err)
		return
	}

	log.Printf("received response %v", resp)

	return
}

func (c *Http2Client) Close() {
	return
}
