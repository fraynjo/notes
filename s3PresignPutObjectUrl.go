package main


import (
	uuid "github.com/satori/go.uuid"
	"path"
	"regexp"
	"strings"
	"time"
  	"github.com/aws/aws-sdk-go/aws/credentials"
  	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
  region = ""
  appId = ""
  secret = ""
  uploadRoot = ""
  bucket = ""
  expired = 300
  fileDomain="https://xxxxxxx.com/"
)

// filename: xxxxxxxxx.png
func GetAWSS3PredesignedUrl(filename string) (putUrl string, fileUrl string, err error) {
	newFileName := uuid.NewV4().String()
	sess, _ := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				appId, secret, "",
			),
		},
	)
	// Create S3 service client
	svc := s3.New(sess)
	key := uploadRoot + newFileName + "_" + getFilename(filename)
	resp, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	fileUrl = fileDomain + key
	putUrl, err = resp.Presign(5 * time.Minute)
	return
}

func getFilename(filename string) string {
	//获取文件名带后缀
	filenameWithSuffix := path.Base(filename)
	//获取文件后缀
	fileSuffix := path.Ext(filenameWithSuffix)
	//获取文件名
	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	if filenameOnly != "" {
		reg := regexp.MustCompile(`\\s+`)
		reg2 := regexp.MustCompile(`\\\\`)
		return reg2.ReplaceAllLiteralString(reg.ReplaceAllString(filenameOnly, ""), "") + fileSuffix
	} else {
		return "randomNumber" + fileSuffix
	}
}
