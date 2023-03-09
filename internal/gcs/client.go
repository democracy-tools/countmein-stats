package gcs

import (
	"bytes"
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/democracy-tools/countmein-stats/internal/env"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

const EnvKeyGCSToken = "GCS_KEY"

type Client interface {
	Update(name string, data []byte) error
	Close() error
}

type ClientWrapper struct {
	store     *storage.Client
	projectId string
	bucket    string
}

func NewClientWrapper(project string, bucketName string) Client {

	if key := env.GetEnvSensitive(EnvKeyGCSToken); key != "" {
		conf, err := google.JWTConfigFromJSON([]byte(key), storage.ScopeReadWrite)
		if err != nil {
			log.Fatalf("failed to config storage JWT with '%v'", err)
		}

		ctx := context.Background()
		client, err := storage.NewClient(ctx, option.WithTokenSource(conf.TokenSource(ctx)))
		if err != nil {
			log.Fatalf("failed to create GCS client with '%v'", err)
		}

		return newClientWrapper(client, project, bucketName)
	}

	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create GCS client without token with '%v'", err)
	}

	return newClientWrapper(client, project, bucketName)
}

func newClientWrapper(client *storage.Client, project string, bucketName string) Client {

	return &ClientWrapper{store: client, projectId: project, bucket: bucketName}
}

func (client *ClientWrapper) Update(name string, data []byte) error {

	writer := client.store.Bucket(client.bucket).Object(name).NewWriter(context.Background())
	writer.ObjectAttrs.CacheControl = "no-cache, max-age=0"
	_, err := io.Copy(writer, bytes.NewReader(data))
	if err != nil {
		log.Errorf("failed to copy file '%s' into bucket '%s' with '%v'", name, client.bucket, err)
		return err
	}
	if err = writer.Close(); err != nil {
		log.Errorf("failed to close GCS file writer. bucket '%s' file '%s'  with '%v'", client.bucket, name, err)
		return err
	}

	return nil
}

func (client *ClientWrapper) Close() error {

	return client.store.Close()
}
