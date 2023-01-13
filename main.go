package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "log"
    "time"

    "sync"
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
    region := "your_region"
    sd, _ := CreateDriver(&region)
    t := time.Now()
    buckArr, _ := sd.BucketList(context.TODO())
    fmt.Println("Time Taken:", time.Since(t))
    out, _ := json.MarshalIndent(buckArr, "", "   ")
    err := os.WriteFile("metadata.json", out, 0644)
    if err != nil {
        log.Fatal(err)
    }
}

func CreateDriver(region *string) (*AwsAdapter, error) {
    AKID := "your_access_key"
    SECRET_KEY := "your_secret_key"
    creds := credentials.NewStaticCredentials(AKID, SECRET_KEY, "")

    sess, err := session.NewSession(&aws.Config{
        Region:      region,
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
    numBuckets := len(output.Buckets)
    bucketArray := make([]model.MetaBucket, numBuckets)
    wg := sync.WaitGroup{}
    for i, bucket := range output.Buckets {
        wg.Add(1)
        go func (i int, bucket *s3.Bucket) {
            defer wg.Done()
            buck := model.MetaBucket{}
            buck.CreationDate = bucket.CreationDate
            buck.Name = *bucket.Name
            loc, _ := svc.GetBucketLocation(&s3.GetBucketLocationInput{Bucket: bucket.Name})
            buck.Region = *loc.LocationConstraint
            tags, _ := svc.GetBucketTagging(&s3.GetBucketTaggingInput{Bucket: bucket.Name})
            tagset := make(map[string]string)
            for _, tag := range tags.TagSet {
                tagset[*tag.Key] = *tag.Value
            } 
            buck.BucketTags = tagset
            newAdap, _ := CreateDriver(&buck.Region)
            err := newAdap.ObjectList(ctx, &buck)
            if err != nil {
                log.Fatal(err)
            }
            bucketArray[i] = buck
        } (i, bucket)
    }
    wg.Wait()
    return bucketArray, err
}

func (ad *AwsAdapter) ObjectList(ctx context.Context, bucket *model.MetaBucket) (error) {
    svc := s3.New(ad.Session)
    output, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: &bucket.Name})
    if err != nil {
        return err
    }

    numObjects := len(output.Contents)
    
    objectArray := make([]model.MetaObject, numObjects)
    wg := sync.WaitGroup{}
    var totSize int64 = 0
    for i, object := range output.Contents {
        wg.Add(1)
        go func (i int, object *s3.Object) {
            defer wg.Done()
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
            objectArray[i] = obj
        } (i, object)
    }
    wg.Wait()
    bucket.NumberOfObjects = int64(numObjects)
    bucket.TotalSize = totSize
    bucket.Objects = objectArray
    return err
}

