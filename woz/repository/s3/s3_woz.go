package s3

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/datahuys/scraperv2/domain"
)

type s3WozRepository struct {
	bucket string
	upl    *s3manager.Uploader
	cl     client.ConfigProvider
}

func NewS3WozRepository(c client.ConfigProvider, bucket string) domain.WozRepository {
	return &s3WozRepository{
		bucket: bucket,
		upl:    s3manager.NewUploader(c),
		cl:     c,
	}
}

func (r *s3WozRepository) GetByID(id int64) (woz domain.Woz, err error) {
	return
}

func (r *s3WozRepository) Store(woz domain.Woz) (err error) {
	key := fmt.Sprintf("%d.json", woz.ID)

	_, err = r.upl.Upload(&s3manager.UploadInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(woz.Payload),

		ContentType: aws.String("application/json"),
		Metadata: map[string]*string{
			"Mtime":  aws.String(fmt.Sprintf("%d", woz.ScrapedAt.Unix())),
			"Status": aws.String(fmt.Sprintf("%d", woz.Status)),
		},
	})
	if err != nil {
		return
	}

	return
}
