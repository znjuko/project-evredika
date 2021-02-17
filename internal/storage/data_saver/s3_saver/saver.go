package s3_saver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"

	"project-evredika/internal/storage/data_saver"
)

const (
	// StatusNotFound ...
	StatusNotFound = "NotFound"
)

type saver struct {
	client *s3.S3
}

func (s *saver) CreateData(ctx context.Context, data *data_saver.Data) (err error) {
	err = s.validateDataExist(data.Metadata)
	if err == nil {
		return fmt.Errorf("file with key '%s', bucket '%s' already exists", data.Key, data.Bucket)
	}

	var awsErr awserr.Error
	if !errors.As(err, &awsErr) || (awsErr.Code() != StatusNotFound && awsErr.Code() != s3.ErrCodeNoSuchKey) {
		return
	}

	if _, err = s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Body:   bytes.NewReader(data.B),
		Bucket: aws.String(data.Bucket),
		Key:    aws.String(data.Key),
	}); err != nil {
		return
	}

	return nil
}

func (s *saver) UpdateData(ctx context.Context, data *data_saver.Data) (err error) {
	if err = s.validateDataExist(data.Metadata); err != nil {
		return
	}

	if _, err = s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Body:   bytes.NewReader(data.B),
		Bucket: aws.String(data.Bucket),
		Key:    aws.String(data.Key),
	}); err != nil {
		return
	}

	return nil
}

func (s *saver) ReadData(ctx context.Context, info *data_saver.Metadata) (data []byte, err error) {
	var response *s3.GetObjectOutput
	if response, err = s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(info.Bucket),
		Key:    aws.String(info.Key),
	}); err != nil {
		return
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

func (s *saver) DeleteData(ctx context.Context, info *data_saver.Metadata) (err error) {
	if _, err = s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(info.Bucket),
		Key:    aws.String(info.Key),
	}); err != nil {
		return
	}

	return
}

func (s *saver) ListData(ctx context.Context, info *data_saver.Metadata) (data []*data_saver.Data, err error) {
	var listResponse *s3.ListObjectsOutput
	if listResponse, err = s.client.ListObjectsWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(info.Bucket),
		Prefix: aws.String(info.Key),
	}); err != nil {
		return
	}

	for _, object := range listResponse.Contents {
		if object == nil {
			continue
		}

		var objectData []byte
		if objectData, err = s.ReadData(ctx, &data_saver.Metadata{
			Key:    aws.StringValue(object.Key),
			Bucket: info.Bucket,
		}); err != nil {
			return
		}

		data = append(data, &data_saver.Data{
			Metadata: data_saver.Metadata{Key: aws.StringValue(object.Key)},
			B:        objectData,
		})
	}

	return data, nil
}

func (s *saver) Initiate(ctx context.Context, bucket string) {
	resp, err := s.client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
		ACL:    aws.String("public-read | public-read-write"),
		Bucket: aws.String(bucket),
	})
	fmt.Println("initiate resp :" ,resp)
	fmt.Println("initiate err :", err)
}

func (s *saver) validateDataExist(info data_saver.Metadata) (err error) {
	if _, err = s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(info.Bucket),
		Key:    aws.String(info.Key),
	}); err != nil {
		return
	}

	return nil
}

// NewS3Saver ...
func NewS3Saver(client *s3.S3) data_saver.DataSaver { return &saver{client: client} }
