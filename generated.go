// Code generated by codegen.go; DO NOT EDIT

package s3protocol

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func newGetObjectInput(req *http.Request) *s3.GetObjectInput {
	var in s3.GetObjectInput
	header := req.Header
	if header == nil {
		header = make(http.Header)
	}
	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		query = make(url.Values)
	}
	if v, ok := header["X-Amz-Checksum-Mode"]; ok && len(v) > 0 {
		in.ChecksumMode = aws.String(v[0])
	}
	if v, ok := header["X-Amz-Expected-Bucket-Owner"]; ok && len(v) > 0 {
		in.ExpectedBucketOwner = aws.String(v[0])
	}
	if v, ok := header["If-Match"]; ok && len(v) > 0 {
		in.IfMatch = aws.String(v[0])
	}
	if v, ok := header["If-Modified-Since"]; ok && len(v) > 0 {
		t, err := http.ParseTime(v[0])
		if err == nil {
			in.IfModifiedSince = aws.Time(t)
		}
	}
	if v, ok := header["If-None-Match"]; ok && len(v) > 0 {
		in.IfNoneMatch = aws.String(v[0])
	}
	if v, ok := header["If-Unmodified-Since"]; ok && len(v) > 0 {
		t, err := http.ParseTime(v[0])
		if err == nil {
			in.IfUnmodifiedSince = aws.Time(t)
		}
	}
	if v, ok := query["partNumber"]; ok && len(v) > 0 {
		i, err := strconv.ParseInt(v[0], 10, 64)
		if err == nil {
			in.PartNumber = aws.Int64(i)
		}
	}
	if v, ok := header["Range"]; ok && len(v) > 0 {
		in.Range = aws.String(v[0])
	}
	if v, ok := header["X-Amz-Request-Payer"]; ok && len(v) > 0 {
		in.RequestPayer = aws.String(v[0])
	}
	if v, ok := query["response-cache-control"]; ok && len(v) > 0 {
		in.ResponseCacheControl = aws.String(v[0])
	}
	if v, ok := query["response-content-disposition"]; ok && len(v) > 0 {
		in.ResponseContentDisposition = aws.String(v[0])
	}
	if v, ok := query["response-content-encoding"]; ok && len(v) > 0 {
		in.ResponseContentEncoding = aws.String(v[0])
	}
	if v, ok := query["response-content-language"]; ok && len(v) > 0 {
		in.ResponseContentLanguage = aws.String(v[0])
	}
	if v, ok := query["response-content-type"]; ok && len(v) > 0 {
		in.ResponseContentType = aws.String(v[0])
	}
	if v, ok := query["response-expires"]; ok && len(v) > 0 {
		t, err := http.ParseTime(v[0])
		if err == nil {
			in.ResponseExpires = aws.Time(t)
		}
	}
	if v, ok := header["X-Amz-Server-Side-Encryption-Customer-Algorithm"]; ok && len(v) > 0 {
		in.SSECustomerAlgorithm = aws.String(v[0])
	}
	if v, ok := header["X-Amz-Server-Side-Encryption-Customer-Key"]; ok && len(v) > 0 {
		in.SSECustomerKey = aws.String(v[0])
	}
	if v, ok := header["X-Amz-Server-Side-Encryption-Customer-Key-Md5"]; ok && len(v) > 0 {
		in.SSECustomerKeyMD5 = aws.String(v[0])
	}
	if v, ok := query["versionId"]; ok && len(v) > 0 {
		in.VersionId = aws.String(v[0])
	}
	return &in
}

func newHeadObjectInput(req *http.Request) *s3.HeadObjectInput {
	var in s3.HeadObjectInput
	header := req.Header
	if header == nil {
		header = make(http.Header)
	}
	query, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		query = make(url.Values)
	}
	if v, ok := header["X-Amz-Checksum-Mode"]; ok && len(v) > 0 {
		in.ChecksumMode = aws.String(v[0])
	}
	if v, ok := header["X-Amz-Expected-Bucket-Owner"]; ok && len(v) > 0 {
		in.ExpectedBucketOwner = aws.String(v[0])
	}
	if v, ok := header["If-Match"]; ok && len(v) > 0 {
		in.IfMatch = aws.String(v[0])
	}
	if v, ok := header["If-Modified-Since"]; ok && len(v) > 0 {
		t, err := http.ParseTime(v[0])
		if err == nil {
			in.IfModifiedSince = aws.Time(t)
		}
	}
	if v, ok := header["If-None-Match"]; ok && len(v) > 0 {
		in.IfNoneMatch = aws.String(v[0])
	}
	if v, ok := header["If-Unmodified-Since"]; ok && len(v) > 0 {
		t, err := http.ParseTime(v[0])
		if err == nil {
			in.IfUnmodifiedSince = aws.Time(t)
		}
	}
	if v, ok := query["partNumber"]; ok && len(v) > 0 {
		i, err := strconv.ParseInt(v[0], 10, 64)
		if err == nil {
			in.PartNumber = aws.Int64(i)
		}
	}
	if v, ok := header["Range"]; ok && len(v) > 0 {
		in.Range = aws.String(v[0])
	}
	if v, ok := header["X-Amz-Request-Payer"]; ok && len(v) > 0 {
		in.RequestPayer = aws.String(v[0])
	}
	if v, ok := header["X-Amz-Server-Side-Encryption-Customer-Algorithm"]; ok && len(v) > 0 {
		in.SSECustomerAlgorithm = aws.String(v[0])
	}
	if v, ok := header["X-Amz-Server-Side-Encryption-Customer-Key"]; ok && len(v) > 0 {
		in.SSECustomerKey = aws.String(v[0])
	}
	if v, ok := header["X-Amz-Server-Side-Encryption-Customer-Key-Md5"]; ok && len(v) > 0 {
		in.SSECustomerKeyMD5 = aws.String(v[0])
	}
	if v, ok := query["versionId"]; ok && len(v) > 0 {
		in.VersionId = aws.String(v[0])
	}
	return &in
}

func makeHeaderFromGetObjectOutput(out *s3.GetObjectOutput) http.Header {
	header := make(http.Header)
	if out == nil {
		return header
	}
	if out.AcceptRanges != nil {
		header.Set("Accept-Ranges", aws.StringValue(out.AcceptRanges))
	}
	if out.BucketKeyEnabled != nil {
		header.Set("X-Amz-Server-Side-Encryption-Bucket-Key-Enabled", strconv.FormatBool(aws.BoolValue(out.BucketKeyEnabled)))
	}
	if out.CacheControl != nil {
		header.Set("Cache-Control", aws.StringValue(out.CacheControl))
	}
	if out.ChecksumCRC32 != nil {
		header.Set("X-Amz-Checksum-Crc32", aws.StringValue(out.ChecksumCRC32))
	}
	if out.ChecksumCRC32C != nil {
		header.Set("X-Amz-Checksum-Crc32c", aws.StringValue(out.ChecksumCRC32C))
	}
	if out.ChecksumSHA1 != nil {
		header.Set("X-Amz-Checksum-Sha1", aws.StringValue(out.ChecksumSHA1))
	}
	if out.ChecksumSHA256 != nil {
		header.Set("X-Amz-Checksum-Sha256", aws.StringValue(out.ChecksumSHA256))
	}
	if out.ContentDisposition != nil {
		header.Set("Content-Disposition", aws.StringValue(out.ContentDisposition))
	}
	if out.ContentEncoding != nil {
		header.Set("Content-Encoding", aws.StringValue(out.ContentEncoding))
	}
	if out.ContentLanguage != nil {
		header.Set("Content-Language", aws.StringValue(out.ContentLanguage))
	}
	if out.ContentLength != nil {
		header.Set("Content-Length", strconv.FormatInt(aws.Int64Value(out.ContentLength), 10))
	}
	if out.ContentRange != nil {
		header.Set("Content-Range", aws.StringValue(out.ContentRange))
	}
	if out.ContentType != nil {
		header.Set("Content-Type", aws.StringValue(out.ContentType))
	}
	if out.DeleteMarker != nil {
		header.Set("X-Amz-Delete-Marker", strconv.FormatBool(aws.BoolValue(out.DeleteMarker)))
	}
	if out.ETag != nil {
		header.Set("Etag", aws.StringValue(out.ETag))
	}
	if out.Expiration != nil {
		header.Set("X-Amz-Expiration", aws.StringValue(out.Expiration))
	}
	if out.Expires != nil {
		header.Set("Expires", aws.StringValue(out.Expires))
	}
	if out.LastModified != nil {
		header.Set("Last-Modified", out.LastModified.Format(http.TimeFormat))
	}
	if out.MissingMeta != nil {
		header.Set("X-Amz-Missing-Meta", strconv.FormatInt(aws.Int64Value(out.MissingMeta), 10))
	}
	if out.ObjectLockLegalHoldStatus != nil {
		header.Set("X-Amz-Object-Lock-Legal-Hold", aws.StringValue(out.ObjectLockLegalHoldStatus))
	}
	if out.ObjectLockMode != nil {
		header.Set("X-Amz-Object-Lock-Mode", aws.StringValue(out.ObjectLockMode))
	}
	if out.ObjectLockRetainUntilDate != nil {
		header.Set("X-Amz-Object-Lock-Retain-Until-Date", out.ObjectLockRetainUntilDate.Format(http.TimeFormat))
	}
	if out.PartsCount != nil {
		header.Set("X-Amz-Mp-Parts-Count", strconv.FormatInt(aws.Int64Value(out.PartsCount), 10))
	}
	if out.ReplicationStatus != nil {
		header.Set("X-Amz-Replication-Status", aws.StringValue(out.ReplicationStatus))
	}
	if out.RequestCharged != nil {
		header.Set("X-Amz-Request-Charged", aws.StringValue(out.RequestCharged))
	}
	if out.Restore != nil {
		header.Set("X-Amz-Restore", aws.StringValue(out.Restore))
	}
	if out.SSECustomerAlgorithm != nil {
		header.Set("X-Amz-Server-Side-Encryption-Customer-Algorithm", aws.StringValue(out.SSECustomerAlgorithm))
	}
	if out.SSECustomerKeyMD5 != nil {
		header.Set("X-Amz-Server-Side-Encryption-Customer-Key-Md5", aws.StringValue(out.SSECustomerKeyMD5))
	}
	if out.SSEKMSKeyId != nil {
		header.Set("X-Amz-Server-Side-Encryption-Aws-Kms-Key-Id", aws.StringValue(out.SSEKMSKeyId))
	}
	if out.ServerSideEncryption != nil {
		header.Set("X-Amz-Server-Side-Encryption", aws.StringValue(out.ServerSideEncryption))
	}
	if out.StorageClass != nil {
		header.Set("X-Amz-Storage-Class", aws.StringValue(out.StorageClass))
	}
	if out.TagCount != nil {
		header.Set("X-Amz-Tagging-Count", strconv.FormatInt(aws.Int64Value(out.TagCount), 10))
	}
	if out.VersionId != nil {
		header.Set("X-Amz-Version-Id", aws.StringValue(out.VersionId))
	}
	if out.WebsiteRedirectLocation != nil {
		header.Set("X-Amz-Website-Redirect-Location", aws.StringValue(out.WebsiteRedirectLocation))
	}
	return header
}

func makeHeaderFromHeadObjectOutput(out *s3.HeadObjectOutput) http.Header {
	header := make(http.Header)
	if out == nil {
		return header
	}
	if out.AcceptRanges != nil {
		header.Set("Accept-Ranges", aws.StringValue(out.AcceptRanges))
	}
	if out.ArchiveStatus != nil {
		header.Set("X-Amz-Archive-Status", aws.StringValue(out.ArchiveStatus))
	}
	if out.BucketKeyEnabled != nil {
		header.Set("X-Amz-Server-Side-Encryption-Bucket-Key-Enabled", strconv.FormatBool(aws.BoolValue(out.BucketKeyEnabled)))
	}
	if out.CacheControl != nil {
		header.Set("Cache-Control", aws.StringValue(out.CacheControl))
	}
	if out.ChecksumCRC32 != nil {
		header.Set("X-Amz-Checksum-Crc32", aws.StringValue(out.ChecksumCRC32))
	}
	if out.ChecksumCRC32C != nil {
		header.Set("X-Amz-Checksum-Crc32c", aws.StringValue(out.ChecksumCRC32C))
	}
	if out.ChecksumSHA1 != nil {
		header.Set("X-Amz-Checksum-Sha1", aws.StringValue(out.ChecksumSHA1))
	}
	if out.ChecksumSHA256 != nil {
		header.Set("X-Amz-Checksum-Sha256", aws.StringValue(out.ChecksumSHA256))
	}
	if out.ContentDisposition != nil {
		header.Set("Content-Disposition", aws.StringValue(out.ContentDisposition))
	}
	if out.ContentEncoding != nil {
		header.Set("Content-Encoding", aws.StringValue(out.ContentEncoding))
	}
	if out.ContentLanguage != nil {
		header.Set("Content-Language", aws.StringValue(out.ContentLanguage))
	}
	if out.ContentLength != nil {
		header.Set("Content-Length", strconv.FormatInt(aws.Int64Value(out.ContentLength), 10))
	}
	if out.ContentType != nil {
		header.Set("Content-Type", aws.StringValue(out.ContentType))
	}
	if out.DeleteMarker != nil {
		header.Set("X-Amz-Delete-Marker", strconv.FormatBool(aws.BoolValue(out.DeleteMarker)))
	}
	if out.ETag != nil {
		header.Set("Etag", aws.StringValue(out.ETag))
	}
	if out.Expiration != nil {
		header.Set("X-Amz-Expiration", aws.StringValue(out.Expiration))
	}
	if out.Expires != nil {
		header.Set("Expires", aws.StringValue(out.Expires))
	}
	if out.LastModified != nil {
		header.Set("Last-Modified", out.LastModified.Format(http.TimeFormat))
	}
	if out.MissingMeta != nil {
		header.Set("X-Amz-Missing-Meta", strconv.FormatInt(aws.Int64Value(out.MissingMeta), 10))
	}
	if out.ObjectLockLegalHoldStatus != nil {
		header.Set("X-Amz-Object-Lock-Legal-Hold", aws.StringValue(out.ObjectLockLegalHoldStatus))
	}
	if out.ObjectLockMode != nil {
		header.Set("X-Amz-Object-Lock-Mode", aws.StringValue(out.ObjectLockMode))
	}
	if out.ObjectLockRetainUntilDate != nil {
		header.Set("X-Amz-Object-Lock-Retain-Until-Date", out.ObjectLockRetainUntilDate.Format(http.TimeFormat))
	}
	if out.PartsCount != nil {
		header.Set("X-Amz-Mp-Parts-Count", strconv.FormatInt(aws.Int64Value(out.PartsCount), 10))
	}
	if out.ReplicationStatus != nil {
		header.Set("X-Amz-Replication-Status", aws.StringValue(out.ReplicationStatus))
	}
	if out.RequestCharged != nil {
		header.Set("X-Amz-Request-Charged", aws.StringValue(out.RequestCharged))
	}
	if out.Restore != nil {
		header.Set("X-Amz-Restore", aws.StringValue(out.Restore))
	}
	if out.SSECustomerAlgorithm != nil {
		header.Set("X-Amz-Server-Side-Encryption-Customer-Algorithm", aws.StringValue(out.SSECustomerAlgorithm))
	}
	if out.SSECustomerKeyMD5 != nil {
		header.Set("X-Amz-Server-Side-Encryption-Customer-Key-Md5", aws.StringValue(out.SSECustomerKeyMD5))
	}
	if out.SSEKMSKeyId != nil {
		header.Set("X-Amz-Server-Side-Encryption-Aws-Kms-Key-Id", aws.StringValue(out.SSEKMSKeyId))
	}
	if out.ServerSideEncryption != nil {
		header.Set("X-Amz-Server-Side-Encryption", aws.StringValue(out.ServerSideEncryption))
	}
	if out.StorageClass != nil {
		header.Set("X-Amz-Storage-Class", aws.StringValue(out.StorageClass))
	}
	if out.VersionId != nil {
		header.Set("X-Amz-Version-Id", aws.StringValue(out.VersionId))
	}
	if out.WebsiteRedirectLocation != nil {
		header.Set("X-Amz-Website-Redirect-Location", aws.StringValue(out.WebsiteRedirectLocation))
	}
	return header
}
