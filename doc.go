/*
Package s3protocol provides the http.RoundTripper interface for Amazon S3 (Simple Storage Service).

The typical use case is to register the "s3" protocol with a http.ransport, as in:

	s := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
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

Amazon S3 supports object versioning.
To access the noncurrent version of an object, use a uri like s3://[BUCKET_NAME]/[OBJECT_NAME]?versionId=[VERSION_ID].
For example,

	resp, err := c.Get("s3://shogo82148-s3protocol/example.txt?versionId=null")
*/
package s3protocol
