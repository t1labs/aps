package dexcom

import (
	"time"
)

// Glucose is a single measurement from your CGM, it will
// have a value (often estimated) and a sampled at time.
type Glucose struct {
	Value     int       `json:"value"`
	SampledAt time.Time `json:"sampledAt"`
}
