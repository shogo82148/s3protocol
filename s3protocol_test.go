package s3protocol_test

import (
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/shogo82148/s3protocol"
)

func ExampleNewTransport() {
	s := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.AnonymousCredentials,
	}))
	s3 := s3protocol.NewTransport(s)

	t := &http.Transport{}
	t.RegisterProtocol("s3", s3)
	c := &http.Client{Transport: t}

	resp, err := c.Get("s3://shogo82148-s3protocol/example.txt")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)

	// Output:
	// Hello Amazon S3!
}
