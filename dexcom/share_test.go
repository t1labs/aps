package dexcom

import (
	"context"
	"testing"

	"github.com/kelseyhightower/envconfig"
)

func TestShare_getLatestGlucose(t *testing.T) {
	type config struct {
		DexcomShareUsername string `required:"true" split_words:"true"`
		DexcomSharePassword string `required:"true" split_words:"true"`
	}
	var c config
	err := envconfig.Process("", &c)
	if err != nil {
		t.Skip(err.Error())
	}

	sh := NewShare(ShareConfig{
		Username: c.DexcomShareUsername,
		Password: c.DexcomSharePassword,
	})

	sessionID, err := sh.login(context.Background())
	if err != nil {
		t.Error(err)
	}

	g, err := sh.getLatestGlucose(context.Background(), sessionID)
	if err != nil {
		t.Error(err)
	}

	if g.Value == 0 {
		t.Error("expected glucose to be returned")
	}
	if g.SampledAt.IsZero() {
		t.Error("expected sampled at to be set")
	}
}

func TestShare_login(t *testing.T) {
	type config struct {
		DexcomShareUsername string `required:"true" split_words:"true"`
		DexcomSharePassword string `required:"true" split_words:"true"`
	}
	var c config
	err := envconfig.Process("", &c)
	if err != nil {
		t.Skip(err.Error())
	}

	sh := NewShare(ShareConfig{
		Username: c.DexcomShareUsername,
		Password: c.DexcomSharePassword,
	})

	_, err = sh.login(context.Background())
	if err != nil {
		t.Error(err)
	}
}
