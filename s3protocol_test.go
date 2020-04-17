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
		Region:      aws.String("us-east-1"),
		Credentials: credentials.AnonymousCredentials,
	}))
	s3 := s3protocol.NewTransport(s)

	t := &http.Transport{}
	t.RegisterProtocol("s3", s3)
	c := &http.Client{Transport: t}

	resp, err := c.Get("s3://shogo82148-jis0208/product-.py")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)

	// Output:
	// import sys
	// import itertools
	//
	// a = [l.rstrip() for l in sys.stdin]
	//
	// for l in itertools.product(a, a):
	//     print(''.join(l))
}
