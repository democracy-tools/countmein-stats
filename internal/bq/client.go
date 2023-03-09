package bq

import (
	"context"
	"fmt"
	"reflect"

	"cloud.google.com/go/bigquery"
	"github.com/democracy-tools/countmein-stats/internal/env"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

const (
	EnvKeyBQToken = "BIGQUERY_KEY"

	TableAnnouncement = "announcement"
)

type Client interface {
	GetAnnouncementCount(from int64) (int64, error)
	Close() error
}

type ClientWrapper struct {
	bqc     *bigquery.Client
	dataset string
}

func NewClientWrapper(project string) Client {

	if key := env.GetEnvSensitive(EnvKeyBQToken); key != "" {
		conf, err := google.JWTConfigFromJSON([]byte(key), bigquery.Scope)
		if err != nil {
			log.Fatalf("failed to config bigquery JWT with %q", err)
		}

		ctx := context.Background()
		client, err := bigquery.NewClient(ctx, project, option.WithTokenSource(conf.TokenSource(ctx)))
		if err != nil {
			log.Fatalf("failed to create bigquery client with %q", err)
		}

		return newClientWrapper(client)
	}

	client, err := bigquery.NewClient(context.Background(), project)
	if err != nil {
		log.Fatalf("failed to create bigquery client without token with %q", err)
	}

	return newClientWrapper(client)
}

func newClientWrapper(client *bigquery.Client) Client {

	return &ClientWrapper{bqc: client, dataset: env.GetBQDataset()}
}

func (c *ClientWrapper) GetAnnouncementCount(from int64) (int64, error) {

	query := c.bqc.Query(`SELECT count(DISTINCT user_id) FROM ` + getTableFullName(c.dataset, TableAnnouncement) + ` WHERE user_time > @time`)
	query.Parameters = []bigquery.QueryParameter{{
		Name:  "time",
		Value: from,
	}}

	iterator, err := query.Read(context.Background())
	if err != nil {
		log.Errorf("failed to execute get announcement count query from time '%d' with '%v'", from, err)
		return -1, err
	}

	var values []bigquery.Value
	err = iterator.Next(&values)
	if err != nil {
		log.Errorf("failed to get announcement count from time '%d' with '%v'", from, err)
		return -1, err
	}
	if len(values) != 1 {
		log.Errorf("unexpected bigquery value '%v' when getting announcement count", values)
		return -1, err
	}
	count, ok := values[0].(int64)
	if !ok {
		log.Errorf("unexpected bigquery value of announcement count '%s'", reflect.TypeOf(values[0]).String())
		return -1, err
	}

	return count, nil
}

func (client *ClientWrapper) Close() error {

	return client.bqc.Close()
}

func getTableFullName(dataset string, table string) string {

	return fmt.Sprintf("democracy-tools.%s.%s", dataset, table)
}
