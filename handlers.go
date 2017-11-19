package s3proxy

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

func ProxyHandler(next http.Handler, svc *s3.S3, bucket string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		obj, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(path),
		})
		if err != nil {
			if e, ok := err.(awserr.Error); ok {
				switch e.Code() {
				case s3.ErrCodeNoSuchBucket, s3.ErrCodeNoSuchKey:
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		setHeader(w, "Cache-Control", aws.StringValue(obj.CacheControl))
		setHeader(w, "Content-Disposition", aws.StringValue(obj.ContentDisposition))
		setHeader(w, "Content-Encoding", aws.StringValue(obj.ContentEncoding))
		setHeader(w, "Content-Language", aws.StringValue(obj.ContentLanguage))
		setHeader(w, "Content-Length", strconv.FormatInt(aws.Int64Value(obj.ContentLength), 10))
		setHeader(w, "Content-Range", aws.StringValue(obj.ContentRange))
		setHeader(w, "Content-Type", aws.StringValue(obj.ContentType))
		setHeader(w, "ETag", aws.StringValue(obj.ETag))
		setHeader(w, "Expires", aws.StringValue(obj.Expires))
		setHeader(w, "Last-Modified", timeToString(obj.LastModified))

		io.Copy(w, obj.Body)

		next.ServeHTTP(w, r)
	})
}

func timeToString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format(http.TimeFormat)
}

func setHeader(w http.ResponseWriter, key, value string) {
	if value != "" {
		w.Header().Add(key, value)
	}
}

func BasicAuthHandler(next http.Handler, auth BasicAuth) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if !ok {
			unauthorized(w, r)
			return
		}

		if !auth.authenticate(username, password) {
			unauthorized(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+r.Host+`"`)
	http.Error(w, "", http.StatusUnauthorized)
}
