package remote

import (
	"bytes"
	"io"
	"io/ioutil"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"mime"
	"os"
	"path/filepath"
)

type S3 struct {
	Bucket    string
	AccessKey string
	SecretKey string
	Region    string
}

func (self *S3) getBucket() *s3.Bucket {
	auth := aws.Auth{
		AccessKey: self.AccessKey,
		SecretKey: self.SecretKey,
	}
	region := aws.Regions[self.Region]
	connection := s3.New(auth, region)
	return connection.Bucket(self.Bucket)
}

func (self *S3) FilesList(prefix string) ([]string, error) {
	bucket := self.getBucket()
	res, err := bucket.List(prefix, "", "", 1000)
	if err != nil {
		return nil, err
	}
	files := []string{}
	for _, v := range res.Contents {
		files = append(files, v.Key)
	}
	return files, nil
}

func (self *S3) Push(path, destination string) error {
	bucket := self.getBucket()
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	perm := s3.BucketOwnerFull
	contentType := mime.TypeByExtension(filepath.Ext(path))
	return bucket.Put(destination, content, contentType, perm)
}

func (self *S3) Pull(filepath, destination string) error {
	bucket := self.getBucket()
	content, err := bucket.Get(filepath)
	if err != nil {
		return err
	}
	destinationFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	if _, err := io.Copy(destinationFile, bytes.NewReader(content)); err != nil {
		destinationFile.Close()
		return err
	}
	return destinationFile.Close()
}

func (self *S3) Remove(filepath string) error {
	bucket := self.getBucket()
	return bucket.Del(filepath)
}
