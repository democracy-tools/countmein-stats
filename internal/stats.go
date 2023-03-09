package internal

import (
	"encoding/json"
	"time"

	"github.com/democracy-tools/countmein-stats/internal/bq"
	"github.com/democracy-tools/countmein-stats/internal/gcs"
	log "github.com/sirupsen/logrus"
)

func Run(bqc bq.Client, store gcs.Client) error {

	count, err := bqc.GetAnnouncementCount(time.Now().Add(time.Hour * (-7)).Unix())
	if err != nil {
		return err
	}
	
	return persist(store, count)
}

func persist(client gcs.Client, count int64) error {

	demonstrations, err := json.Marshal(map[string]interface{}{"demonstrations": struct {
		Count int64 `json:"count"`
	}{Count: count}})
	if err != nil {
		log.Errorf("failed to marshal demonstrations with '%v'", err)
		return err
	}

	return client.Update("demonstrations", demonstrations)
}
