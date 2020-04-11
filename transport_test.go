package s3protocol

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type s3mock struct {
	s3iface.S3API
	getObjectWithContext func(ctx context.Context, in *s3.GetObjectInput, _ ...request.Option) (*s3.GetObjectOutput, error)
}

func (mock *s3mock) GetObjectWithContext(ctx context.Context, in *s3.GetObjectInput, _ ...request.Option) (*s3.GetObjectOutput, error) {
	return mock.getObjectWithContext(ctx, in)
}

func TestRoundTrip(t *testing.T) {
	mock := &s3mock{
		getObjectWithContext: func(ctx context.Context, in *s3.GetObjectInput, _ ...request.Option) (*s3.GetObjectOutput, error) {
			return &s3.GetObjectOutput{
				Body: ioutil.NopCloser(strings.NewReader("Hello S3!")),
			}, nil
		},
	}
	tr := &http.Transport{}
	tr.RegisterProtocol("s3", &Transport{S3: mock})
	c := &http.Client{Transport: tr}
	resp, err := c.Get("s3://bucket-name/object-key?versionId=foobar")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status: want %d, got %d", http.StatusOK, resp.StatusCode)
	}
	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "Hello S3!" {
		t.Errorf("want Hello S3!, got %s", string(got))
	}
}

func TestRoundTrip_StatusMethodNotAllowed(t *testing.T) {
	mock := &s3mock{
		getObjectWithContext: func(ctx context.Context, in *s3.GetObjectInput, _ ...request.Option) (*s3.GetObjectOutput, error) {
			panic("not reach")
		},
	}
	tr := &http.Transport{}
	tr.RegisterProtocol("s3", &Transport{S3: mock})
	c := &http.Client{Transport: tr}
	resp, err := c.Post("s3://bucket-name/object-key?versionId=foobar", "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("unexpected status: want %d, got %d", http.StatusOK, resp.StatusCode)
	}
	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "" {
		t.Errorf(`want "", got %q`, string(got))
	}
}
