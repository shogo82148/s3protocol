package s3protocol

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

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
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	path := req.URL.Path

	ctx := req.Context()
	out, err := t.S3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &host,
		Key:    &path,
	})
	if err != nil {
		return nil, err
	}

	header := make(http.Header)
	header.Set("Content-Type", aws.StringValue(out.ContentType))
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
