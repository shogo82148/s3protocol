package s3protocol

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type s3mock struct {
	s3iface.S3API
	getObjectWithContext  func(ctx context.Context, in *s3.GetObjectInput, _ ...request.Option) (*s3.GetObjectOutput, error)
	headObjectWithContext func(ctx context.Context, in *s3.HeadObjectInput, _ ...request.Option) (*s3.HeadObjectOutput, error)
}

func (mock *s3mock) GetObjectWithContext(ctx context.Context, in *s3.GetObjectInput, _ ...request.Option) (*s3.GetObjectOutput, error) {
	return mock.getObjectWithContext(ctx, in)
}

func (mock *s3mock) HeadObjectWithContext(ctx context.Context, in *s3.HeadObjectInput, _ ...request.Option) (*s3.HeadObjectOutput, error) {
	return mock.headObjectWithContext(ctx, in)
}

func TestRoundTrip(t *testing.T) {
	mock := &s3mock{
		getObjectWithContext: func(ctx context.Context, in *s3.GetObjectInput, _ ...request.Option) (*s3.GetObjectOutput, error) {
			// a test case for string value
			if aws.StringValue(in.VersionId) != "foobar" {
				t.Errorf("unexpected version id: want %q, got %q", "footbar", aws.StringValue(in.VersionId))
			}

			// a test case for time.Time
			if want, got := time.Date(2015, time.October, 21, 7, 28, 0, 0, time.UTC), in.IfModifiedSince; got == nil || !got.Equal(want) {
				t.Errorf("unexpected If-Modified-Since: want %s, got %s", want, got)
			}

			// a test case for integer value
			if aws.Int64Value(in.PartNumber) != 1 {
				t.Errorf("unexpected part number: want %d, got %d", 1, aws.Int64Value(in.PartNumber))
			}

			return &s3.GetObjectOutput{
				ContentType:  aws.String("plain/text"),
				DeleteMarker: aws.Bool(false),
				LastModified: aws.Time(time.Date(2015, time.October, 21, 7, 28, 0, 0, time.UTC)),
				PartsCount:   aws.Int64(10),
				Body:         ioutil.NopCloser(strings.NewReader("Hello S3!")),
			}, nil
		},
	}
	tr := &http.Transport{}
	tr.RegisterProtocol("s3", &Transport{S3: mock})
	c := &http.Client{Transport: tr}
	req, err := http.NewRequest(http.MethodGet, "s3://bucket-name/object-key?versionId=foobar&partNumber=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("If-Modified-Since", "Wed, 21 Oct 2015 07:28:00 GMT")
	resp, err := c.Do(req)
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

	// a test case for string value
	if resp.Header.Get("Content-Type") != "plain/text" {
		t.Errorf("want %s, got %s", "plain/text", resp.Header.Get("Content-Type"))
	}

	// a test case for boolean value
	if resp.Header.Get("X-Amz-Delete-Marker") != "false" {
		t.Errorf("want %s, got %s", "false", resp.Header.Get("X-Amz-Delete-Marker"))
	}

	// a test case for time.Time
	if resp.Header.Get("Last-Modified") != "Wed, 21 Oct 2015 07:28:00 GMT" {
		t.Errorf("want %s, got %s", "Wed, 21 Oct 2015 07:28:00 GMT", resp.Header.Get("Last-Modified"))
	}

	// a test case for integer value
	if resp.Header.Get("X-Amz-Mp-Parts-Count") != "10" {
		t.Errorf("want %s, got %s", "10", resp.Header.Get("X-Amz-Mp-Parts-Count"))
	}
}

func TestRoundTrip_HEAD(t *testing.T) {
	mock := &s3mock{
		headObjectWithContext: func(ctx context.Context, in *s3.HeadObjectInput, _ ...request.Option) (*s3.HeadObjectOutput, error) {
			if aws.StringValue(in.VersionId) != "foobar" {
				t.Errorf("unexpected version id: want %q, got %q", "footbar", aws.StringValue(in.VersionId))
			}
			return &s3.HeadObjectOutput{
				ContentType: aws.String("image/png"),
			}, nil
		},
	}
	tr := &http.Transport{}
	tr.RegisterProtocol("s3", &Transport{S3: mock})
	c := &http.Client{Transport: tr}
	resp, err := c.Head("s3://bucket-name/object-key?versionId=foobar")
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
	if string(got) != "" {
		t.Errorf(`want "", got %q`, string(got))
	}
	if resp.Header.Get("Content-Type") != "image/png" {
		t.Errorf("want %s, got %s", "image/png", resp.Header.Get("Content-Type"))
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

func TestRoundTrip_NotFound(t *testing.T) {
	mock := &s3mock{
		getObjectWithContext: func(ctx context.Context, in *s3.GetObjectInput, _ ...request.Option) (*s3.GetObjectOutput, error) {
			aerr := awserr.New("not found", "not found", errors.New("not found"))
			return nil, awserr.NewRequestFailure(aerr, http.StatusNotFound, "request-id")
		},
	}
	tr := &http.Transport{}
	tr.RegisterProtocol("s3", &Transport{S3: mock})
	c := &http.Client{Transport: tr}
	req, err := http.NewRequest(http.MethodGet, "s3://bucket-name/object-key?versionId=foobar&partNumber=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected status: want %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "" {
		t.Errorf("want %q, got %q", "", string(got))
	}
}
