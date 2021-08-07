# s3protocol

Package s3protocol provides the [http.RoundTripper](https://golang.org/pkg/net/http/#RoundTripper) interface for [Amazon S3 (Simple Storage Service)](https://docs.aws.amazon.com/s3/index.html).

The typical use case is to register the "s3" protocol with a [http.Transport](https://golang.org/pkg/net/http/#Transport), as in:

```go
s := session.Must(session.NewSession(&aws.Config{
    Credentials: credentials.AnonymousCredentials,
}))
s3 := s3protocol.NewTransport(s)

t := &http.Transport{}
t.RegisterProtocol("s3", s3)
c := &http.Client{Transport: t}

resp, err := c.Get("s3://shogo82148-s3protocol/example.txt")
if err != nil {
    // handle error
}
defer resp.Body.Close()
// read resp.Body
```

Amazon S3 supports object versioning.
To access the noncurrent version of an object, use a uri like `s3://[BUCKET_NAME]/[OBJECT_NAME]?versionId=[VERSION_ID]`.
For example,

```go
resp, err := c.Get("s3://shogo82148-s3protocol/example.txt?versionId=null")
```
