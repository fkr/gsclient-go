package gsclient

import (
	"encoding/json"
	"time"
)

const gsTimeLayout = "2006-01-02T15:04:05Z"

//GSTime is custom time type of Gridscale
type GSTime struct {
	time.Time
}

//UnmarshalJSON custom unmarshaller for GSTime
func (t *GSTime) UnmarshalJSON(b []byte) error {
	var tstring string
	if err := json.Unmarshal(b, &tstring); err != nil {
		return err
	}
	parsedTime, err := time.Parse(gsTimeLayout, tstring)
	*t = GSTime{parsedTime}
	return err
}

//MarshalJSON custom marshaller for GSTime
func (t GSTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Format(gsTimeLayout))
}

type serverHardwareProfile struct {
	string
}

//MarshalJSON custom marshal for serverHardwareProfile
func (s serverHardwareProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.string)
}

type storageType struct {
	string
}

//MarshalJSON custom marshal for storageType
func (s storageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.string)
}

type ipAddressType struct {
	int
}

//MarshalJSON custom marshal for ipAddressType
func (i ipAddressType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.int)
}

type loadbalancerAlgorithm struct {
	string
}

//MarshalJSON custom marshal for loadbalancerAlgorithm
func (l loadbalancerAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.string)
}
