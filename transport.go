package s3protocol

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//go:generate go run codegen.go

type s3api struct {
	mu  sync.RWMutex
	svc s3iface.S3API
}

func (c *s3api) get(ctx context.Context, t *Transport, bucket string) (s3iface.S3API, error) {
	c.mu.RLock()
	svc := c.svc
	c.mu.RUnlock()
	if svc != nil {
		return svc, nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.svc != nil {
		return c.svc, nil
	}

	region, err := t.getBucketRegion(ctx, bucket)
	if err != nil {
		return nil, err
	}
	var cfg aws.Config
	cfg.Region = aws.String(region)
	svc = s3.New(t.config, &cfg)
	c.svc = svc

	return svc, nil
}

// Transport serving the S3 objects.
type Transport struct {
	config client.ConfigProvider

	// s3 api client for getting the region
	svc s3iface.S3API

	// regional s3 api clients
	s3 sync.Map
}

// NewTransport returns a new Transport.
func NewTransport(c client.ConfigProvider) *Transport {
	svc := s3.New(c)
	req, _ := svc.HeadBucketRequest(&s3.HeadBucketInput{
		Bucket: aws.String("dummy"),
	})
	if aws.StringValue(req.Config.Region) == "" && aws.StringValue(req.Config.Endpoint) == "" {
		// GetBucketRegion needs Config.Region or Config.Endpoint, but neither found.
		// fall back to the default region.
		var cfg aws.Config
		cfg.Region = aws.String("us-east-1")
		svc = s3.New(c, &cfg)
	}

	return &Transport{
		config: c,
		svc:    svc,
	}
}

// RoundTrip implements http.RoundTripper.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	switch req.Method {
	case http.MethodGet:
		return t.getObject(req)
	case http.MethodHead:
		return t.headObject(req)
	}
	return &http.Response{
		Status:     "405 Method Not Allowed",
		StatusCode: http.StatusMethodNotAllowed,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header:     make(http.Header),
		Body:       http.NoBody,
		Close:      true,
	}, nil
}

func (t *Transport) getObject(req *http.Request) (*http.Response, error) {
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	path := strings.TrimPrefix(req.URL.Path, "/")

	ctx := req.Context()
	svc, err := t.getBucketClient(ctx, host)
	if err != nil {
		return handleError(nil, err)
	}

	in := newGetObjectInput(req)
	in.Bucket = &host
	in.Key = &path
	out, err := svc.GetObjectWithContext(ctx, in)
	header := makeHeaderFromGetObjectOutput(out)
	if err != nil {
		return handleError(header, err)
	}

	return &http.Response{
		Status:        "200 OK",
		StatusCode:    http.StatusOK,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
		Header:        header,
		Body:          out.Body,
		ContentLength: aws.Int64Value(out.ContentLength),
		Close:         true,
	}, nil
}

func (t *Transport) headObject(req *http.Request) (*http.Response, error) {
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	path := strings.TrimPrefix(req.URL.Path, "/")

	ctx := req.Context()
	svc, err := t.getBucketClient(ctx, host)
	if err != nil {
		return handleError(nil, err)
	}

	in := newHeadObjectInput(req)
	in.Bucket = &host
	in.Key = &path
	out, err := svc.HeadObjectWithContext(ctx, in)
	header := makeHeaderFromHeadObjectOutput(out)
	if err != nil {
		return handleError(header, err)
	}

	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header:     header,
		Body:       http.NoBody,
		Close:      true,
	}, nil
}

func handleError(header http.Header, err error) (*http.Response, error) {
	if header == nil {
		header = make(http.Header)
	}
	if err, ok := awsRequestFailure(err); ok {
		code := err.StatusCode()
		return &http.Response{
			Status:     fmt.Sprintf("%d %s", code, http.StatusText(code)),
			StatusCode: code,
			Proto:      "HTTP/1.0",
			ProtoMajor: 1,
			ProtoMinor: 0,
			Header:     header,
			Body:       http.NoBody,
			Close:      true,
		}, nil
	}
	return &http.Response{
		Status:     "500 Internal Server Error",
		StatusCode: http.StatusInternalServerError,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header:     header,
		Body:       http.NoBody,
		Close:      true,
	}, nil
}

func (t *Transport) getBucketClient(ctx context.Context, bucket string) (s3iface.S3API, error) {
	if c, ok := t.s3.Load(bucket); ok {
		return c.(*s3api).get(ctx, t, bucket)
	}

	c, _ := t.s3.LoadOrStore(bucket, new(s3api))
	return c.(*s3api).get(ctx, t, bucket)
}

func (t *Transport) getBucketRegion(ctx context.Context, bucket string) (string, error) {
	region, err := s3manager.GetBucketRegionWithClient(ctx, t.svc, bucket)
	if err != nil {
		return "", err
	}
	return region, nil
}

func awsRequestFailure(err error) (awserr.RequestFailure, bool) {
	for err != nil {
		if err, ok := err.(awserr.RequestFailure); ok {
			return err, true
		}
		if aerr, ok := err.(awserr.Error); ok {
			err = aerr.OrigErr()
		}
	}
	return nil, false
}
