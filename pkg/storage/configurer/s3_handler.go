package configurer

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kelseyhightower/envconfig"

	"project-evredika/internal/storage/data_saver"
	"project-evredika/internal/storage/data_saver/s3_saver"
)

type S3Cfg struct {
	S3Endpoint       string `envconfig:"S3_ENDPOINT" required:"true"`
	S3Region         string `envconfig:"S3_REGION" required:"true"`
	AccessKeyID      string `envconfig:"ACCESS_KEY_ID" required:"true"`
	SecretAccessKey  string `envconfig:"SECRET_ACCESS_KEY" required:"true"`
	DisableSSL       bool   `envconfig:"DISABLE_SSL" required:"true"`
	S3ForcePathStyle bool   `envconfig:"S3_FORCE_PATH_STYLE" required:"true"`
}

func CreateS3Storage(ctx context.Context, bucket string) (st data_saver.DataSaver, err error) {
	var cfg S3Cfg
	if err = envconfig.Process("", &cfg); err != nil {
		return
	}

	svc := s3.New(
		session.Must(session.NewSession()),
		aws.NewConfig().WithEndpoint(cfg.S3Endpoint),
		aws.NewConfig().WithRegion(cfg.S3Region),
		aws.NewConfig().WithCredentials(
			credentials.NewStaticCredentials(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			),
		),
		aws.NewConfig().WithDisableSSL(cfg.DisableSSL),
		aws.NewConfig().WithS3ForcePathStyle(cfg.S3ForcePathStyle),
	)

	st = s3_saver.NewS3Saver(svc)
	st.Initiate(ctx, bucket)
	return
}
