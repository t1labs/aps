package main

import (
	"context"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/kelseyhightower/envconfig"
	"github.com/t1labs/aps/dexcom"
)

type config struct {
	DexcomShareUsername string `required:"true" split_words:"true"`
	DexcomSharePassword string `required:"true" split_words:"true"`
}

func main() {
	l := log.NewJSONLogger(os.Stdout)
	l = log.WithPrefix(l, "date", log.DefaultTimestampUTC)
	info := log.WithPrefix(l, "level", "info")
	danger := log.WithPrefix(l, "level", "error")

	var c config
	err := envconfig.Process("", &c)
	if err != nil {
		danger.Log("msg", "could not process env", "err", err.Error())
		os.Exit(1)
	}

	sh := dexcom.NewShare(dexcom.ShareConfig{
		Username: c.DexcomShareUsername,
		Password: c.DexcomSharePassword,
	})

	l.Log("msg", "starting aps")

	var gs = make(chan dexcom.Glucose)
	var errs = make(chan error)
	go sh.ListenForGlucoses(context.Background(), gs, errs)

	for {
		select {
		case g := <-gs:
			info.Log("glucose", g.Value, "unit", g.Unit, "sampledAt", g.SampledAt)
		case err := <-errs:
			danger.Log("err", err.Error())
		}
	}
}
