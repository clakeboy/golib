package utils

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"fmt"
)

//Get bucket
func GetBucket(bucketName string,endpoint string, accessID string, accessKey string) (*oss.Bucket, error) {
	// New Client
	client, err := oss.New(endpoint, accessID, accessKey)
	if err != nil {
		fmt.Println("init oss field")
		return nil, err
	}
	// check bucket exist
	flag,err := client.IsBucketExist(bucketName)

	if err != nil {
		fmt.Println("check bucket field")
		return nil,err
	}

	// if not exist bucket then Create Bucket
	if !flag {
		err = client.CreateBucket(bucketName)
		if err != nil {
			fmt.Println("create bucket field")
			return nil, err
		}
	}


	// Get Bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("get bucket field")
		return nil, err
	}

	return bucket, nil
}