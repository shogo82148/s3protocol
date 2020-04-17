package s3protocol

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
	path := strings.TrimPrefix(req.URL.Path, "/")

	in := newGetObjectInput(req)
	in.Bucket = &host
	in.Key = &path
	ctx := req.Context()
	out, err := t.S3.GetObjectWithContext(ctx, in)
	if err != nil {
		return handleError(err)
	}

	header := makeHeaderFromGetObjectOutput(out)
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

	in := newHeadObjectInput(req)
	in.Bucket = &host
	in.Key = &path
	ctx := req.Context()
	out, err := t.S3.HeadObjectWithContext(ctx, in)
	if err != nil {
		return handleError(err)
	}

	header := makeHeaderFromHeadObjectOutput(out)
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

func handleError(err error) (*http.Response, error) {
	if err, ok := awsRequestFailure(err); ok {
		code := err.StatusCode()
		return &http.Response{
			Status:     fmt.Sprintf("%d %s", code, http.StatusText(code)),
			StatusCode: code,
			Proto:      "HTTP/1.0",
			ProtoMajor: 1,
			ProtoMinor: 0,
			Header:     make(http.Header),
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
		Header:     make(http.Header),
		Body:       http.NoBody,
		Close:      true,
	}, nil
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
