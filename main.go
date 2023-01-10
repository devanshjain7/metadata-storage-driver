package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	// "sync"
	"v2/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AwsAdapter struct {
	Session *session.Session
}

func main() {

    sd, _ := CreateDriver("ap-south-1")
    buckArr, _ := sd.BucketList(context.TODO())
	fmt.Print(len(buckArr))
	out, _ := json.MarshalIndent(buckArr, "", "   ")
    err := ioutil.WriteFile("metadata.json", out, 0644)
    if err != nil {
        log.Fatal(err)
    }
}

func CreateDriver(region string) (*AwsAdapter, error) {
	AKID := "your_access_key"
	SECRET_KEY := "your_secret_key"
	creds := credentials.NewStaticCredentials(AKID, SECRET_KEY, "")

	sess, err := session.NewSession(&aws.Config{
		Region:      &region,
		Credentials: creds,
	})
	if err != nil {
		return nil, err
	}

	adap := &AwsAdapter{Session: sess}

	return adap, nil
}

func (ad *AwsAdapter) BucketList(ctx context.Context) ([]model.MetaBucket, error) {
	svc := s3.New(ad.Session)

	output, err := svc.ListBuckets(&s3.ListBucketsInput{})

	if err != nil {
		log.Fatal(err)
	}
	bucketArray := make([]model.MetaBucket, 0)
	// var wg sync.WaitGroup
	// var lock sync.Mutex
	for _, bucket := range output.Buckets {
		// wg.Add(1)
		buck := model.MetaBucket{}
		buck.CreationDate = bucket.CreationDate
		buck.Name = *bucket.Name
		loc, _ := svc.GetBucketLocation(&s3.GetBucketLocationInput{Bucket: bucket.Name})
		buck.Region = *loc.LocationConstraint
		ad.Session.Config.Region = loc.LocationConstraint
		svc = s3.New(ad.Session)
		tags, _ := svc.GetBucketTagging(&s3.GetBucketTaggingInput{Bucket: bucket.Name})
		tagset := make(map[string]string)
		for _, tag := range tags.TagSet {
			tagset[*tag.Key] = *tag.Value
		} 
		buck.BucketTags = tagset
		objects, err := ad.ObjectList(ctx, &buck)
		buck.Objects = objects
		if err == nil {
			// lock.Lock()
			// defer lock.Unlock()
			bucketArray = append(bucketArray, buck)
		}
	}
	// wg.Wait()
	return bucketArray, err
}

func (ad *AwsAdapter) ObjectList(ctx context.Context, bucket *model.MetaBucket) ([]model.MetaObject, error) {
	svc := s3.New(ad.Session)

	output, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: &bucket.Name})
	if err != nil {
        return nil, err
	}
	objectArray := make([]model.MetaObject, 0)
	// var wg sync.WaitGroup
	// var lock sync.Mutex
	var numObjs, totSize int64 = 0, 0
	for _, object := range output.Contents {
		// wg.Add(1)
		numObjs += 1
		obj := model.MetaObject{}
		obj.LastModifiedDate = object.LastModified
		obj.ObjectName = *object.Key
		obj.Size = *object.Size
		totSize += obj.Size
		obj.StorageClass = *object.StorageClass
		
		meta, _ := svc.HeadObject(&s3.HeadObjectInput{Bucket: &bucket.Name, Key: object.Key})
		if meta.ServerSideEncryption != nil {obj.ServerSideEncryption = *meta.ServerSideEncryption}
		if meta.VersionId != nil {obj.VersionId = *meta.VersionId}
		obj.ObjectType = *meta.ContentType 
		if meta.Expires != nil {obj.ExpiresDate = *meta.Expires}
		if meta.ReplicationStatus != nil {obj.ReplicationStatus = *meta.ReplicationStatus}
		// lock.Lock()
		// defer lock.Unlock()
		objectArray = append(objectArray, obj)
		// wg.Done()
	}
	// wg.Wait()
	bucket.NumberOfObjects = numObjs
	bucket.TotalSize = totSize
	return objectArray, err
}

