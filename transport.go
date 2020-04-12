package s3protocol

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

//go:generate go run codegen.go

// Transport serving the S3 objects.
type Transport struct {
	S3 s3iface.S3API
}

// NewTransport returns a new Transport.
func NewTransport(c client.ConfigProvider) *Transport {
	return &Transport{
		S3: s3.New(c),
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
	path := req.URL.Path

	in, err := newGetObjectInput(req)
	if err != nil {
		return nil, err
	}
	in.Bucket = &host
	in.Key = &path
	ctx := req.Context()
	out, err := t.S3.GetObjectWithContext(ctx, in)
	if err != nil {
		return nil, err
	}

	header, err := makeHeaderFromGetObjectOutput(out)
	if err != nil {
		return nil, err
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
	path := req.URL.Path

	in, err := newHeadObjectInput(req)
	if err != nil {
		return nil, err
	}
	in.Bucket = &host
	in.Key = &path
	ctx := req.Context()
	out, err := t.S3.HeadObjectWithContext(ctx, in)
	if err != nil {
		return nil, err
	}

	header, err := makeHeaderFromHeadObjectOutput(out)
	if err != nil {
		return nil, err
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
