package main

import (
	"aws/config"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"

	//"github.com/aws/aws-sdk-go/aws/session"
	"io/ioutil"
	//"os"
	//
	"github.com/aws/aws-sdk-go/aws"
)

func main() {

	if err := config.Init(); err != nil {
		log.Fatal(err)
		return
	}
	env := os.Getenv("ENVIRON")
	uploader, err := InitAws(env)
	if err != nil {
		log.Fatal(err)
		return
	}

	bucket := aws.String(config.Configuration.GetString(env + ".sbucket"))
	key := config.Configuration.GetString("path")
	path := "/Users/lyw/workspace/aws/source/share"
	dir_list, e := ioutil.ReadDir(path)
	if e != nil {
		fmt.Println("read dir error")
		return
	}
	for _, v := range dir_list {
		dirName := v.Name()
		var tmpPathKey string
		if len(key) > 0 {
			tmpPathKey = key + "/" + dirName
		} else {
			tmpPathKey = dirName
		}

		tmp_file_list, _ := ioutil.ReadDir(path + "/" + dirName)
		for _, v := range tmp_file_list {
			fileKey := tmpPathKey + "/" + v.Name()
			filePath := path + "/" + dirName + "/" + v.Name()
			Upload(uploader, bucket, fileKey, filePath)
		}
	}
	fmt.Println("upload finish")

}

func InitAws(env string) (*s3manager.Uploader, error) {
	id := config.Configuration.GetString(env + ".id")
	secret := config.Configuration.GetString(env + ".secret")
	region := config.Configuration.GetString(env + ".region")
	fmt.Printf("id:%s, secret:%s, region:%s \n", id, secret, region)
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(id, secret, ""),
		Region:      aws.String(region),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		fmt.Println("session error", err)
		return nil, err
	}
	return s3manager.NewUploader(newSession), nil
}

func Upload(uploader *s3manager.Uploader, bucket *string, fileKey, filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("failed to open file %q, %v", filePath, err)
		return
	}
	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      bucket,
		Key:         aws.String(fileKey),
		Body:        f,
		ContentType: aws.String("png"),
		ACL:         aws.String("public-read"),
	}, func(u *s3manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024 // 分块大小,当文件体积超过10M开始进行分块上传
		u.LeavePartsOnError = true
		u.Concurrency = 3
	}) //并发数
	if err != nil {
		fmt.Printf("Failed to upload data to %s/%s, %s\n", *bucket, filePath, err.Error())
		return
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)
}
