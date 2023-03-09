package main

import (
	"github.com/democracy-tools/countmein-stats/internal"
	"github.com/democracy-tools/countmein-stats/internal/bq"
	"github.com/democracy-tools/countmein-stats/internal/env"
	"github.com/democracy-tools/countmein-stats/internal/gcs"
)

func main() {

	env.Load()
	project := env.GetProject()

	bqc := bq.NewClientWrapper(project)
	defer bqc.Close()

	store := gcs.NewClientWrapper(project, env.GetBucket())
	defer store.Close()

	_ = internal.Run(bqc, store)
}
