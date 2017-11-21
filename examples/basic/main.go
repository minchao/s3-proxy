package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/minchao/s3-proxy"
)

func main() {
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	svc := s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_S3_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_S3_ID"),
			os.Getenv("AWS_S3_SECRET"),
			"",
		),
	})))

	auth := s3proxy.SimpleBasicAuth{
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}

	http.Handle("/", s3proxy.BasicAuthHandler(s3proxy.ProxyHandler(final, svc, os.Getenv("S3_BUCKET")), &auth))
	http.ListenAndServe(":8080", nil)
}
