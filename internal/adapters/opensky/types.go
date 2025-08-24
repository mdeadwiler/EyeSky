package opensky

import (
	"time" // OpenSky API contains timestamp fields
)

type StateVectorResponse struct {
	Time int64 `json:"time"`
	States [][]interface{} `json:"states"`
}
