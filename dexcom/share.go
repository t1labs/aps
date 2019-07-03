package dexcom

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	shareApplicationID    = "d89443d2-327c-4a6f-89e5-496bbb0317db"
	shareURLLogin         = "https://share1.dexcom.com/ShareWebServices/Services/General/LoginPublisherAccountByName"
	shareURLLatestGlucose = "https://share1.dexcom.com/ShareWebServices/Services/Publisher/ReadPublisherLatestGlucoseValues"
	shareUserAgent        = "Dexcom Share/3.0.2.11 CFNetwork/711.2.23 Darwin/14.0.0"
	shareInterval         = time.Second * 60
)

// ShareConfig holds information around connecting to the Dexcom API.
// You should keep your username and password secret.
type ShareConfig struct {
	Username string
	Password string
}

// Share will connect to the Dexcom API and management the underlying
// http connections.
type Share struct {
	config ShareConfig
	c      http.Client
}

// NewShare will create a new instance of Share, initialized with
// some default values.
func NewShare(config ShareConfig) *Share {
	sh := Share{
		config: config,
		c: http.Client{
			Timeout: time.Second * 3,
		},
	}

	return &sh
}

// ListenForGlucoses will check every shareInterval for a new Glucose
// and send it onto the channel when it is found. This function will not
// stop in the face of an error, it will send it on the error channel.
func (sh *Share) ListenForGlucoses(ctx context.Context, gs chan Glucose, errs chan error) {
	sessionID, err := sh.login(ctx)
	if err != nil {
		errs <- err
	}
	g, err := sh.getLatestGlucose(ctx, sessionID)
	if err != nil {
		errs <- err
	}
	gs <- g

	ticker := time.NewTicker(shareInterval).C

	for {
		select {
		case <-ticker:
			sessionID, err := sh.login(ctx)
			if err != nil {
				errs <- err
				break
			}
			g, err := sh.getLatestGlucose(ctx, sessionID)
			if err != nil {
				errs <- err
				break
			}
			gs <- g
		}
	}
}

type shareGlucose struct {
	Value     int    `json:"Value"`
	SampledAt string `json:"WT"`
}

func (sh *Share) getLatestGlucose(ctx context.Context, sessionID string) (Glucose, error) {
	q := url.Values{}
	q.Set("sessionID", sessionID)
	q.Set("maxCount", "1")
	q.Set("minutes", "1440")

	uri := fmt.Sprintf("%s?%s", shareURLLatestGlucose, q.Encode())
	req, err := http.NewRequest(http.MethodPost, uri, nil)
	if err != nil {
		return Glucose{}, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("user-agent", shareUserAgent)
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-length", "0")

	resp, err := sh.c.Do(req)
	if err != nil {
		return Glucose{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return Glucose{}, fmt.Errorf("dexcom: bad status code: %s", resp.Status)
	}

	// Example response
	// [
	// 	{
	// 		"DT": "/Date(1562162271000+0000)/",
	// 		"ST": "/Date(1562176671000)/",
	// 		"Trend": 4,
	// 		"Value": 148,
	// 		"WT": "/Date(1562176671000)/"
	// 	}
	// ]
	var sgs []shareGlucose
	err = json.NewDecoder(resp.Body).Decode(&sgs)
	if err != nil {
		return Glucose{}, err
	}

	for _, sg := range sgs {
		g := Glucose{
			Value:     sg.Value,
			Unit: "mg/dl",
			SampledAt: time.Now(),
		}
		return g, nil
	}

	return Glucose{}, fmt.Errorf("dexcom: no latest glucose")
}

type shareLoginRequest struct {
	AccountName   string `json:"accountName"`
	ApplicationID string `json:"applicationId"`
	Password      string `json:"password"`
}

func (sh *Share) login(ctx context.Context) (string, error) {
	loginReq := shareLoginRequest{
		AccountName:   sh.config.Username,
		ApplicationID: shareApplicationID,
		Password:      sh.config.Password,
	}
	b, err := json.Marshal(loginReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, shareURLLogin, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("user-agent", shareUserAgent)
	req.Header.Set("accept", "application/json")

	resp, err := sh.c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("dexcom: bad status code: %s", resp.Status)
	}

	// Response should be
	// "009afa03-5350-4c87-9b43-39eecb16b10f"
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	sessionID := strings.Replace(string(rb), `"`, "", -1)

	return sessionID, nil
}
