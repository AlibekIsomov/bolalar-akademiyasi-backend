package models

import (
	"encoding/json"
	"fmt"
)

type Status string

const (
	Active   Status = "active"
	Inactive Status = "inactive"
	Pending  Status = "pending"
	Agree    Status = "agree"
)

var statusToString = map[Status]string{
	Active:   "active",
	Inactive: "inactive",
	Pending:  "pending",
	Agree:    "agree",
}

var stringToStatus = map[string]Status{
	"active":   Active,
	"inactive": Inactive,
	"pending":  Pending,
	"agree":    Agree,
}

func (s Status) String() string {
	return statusToString[s]
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var statusStr string
	if err := json.Unmarshal(data, &statusStr); err != nil {
		return err
	}

	status, ok := stringToStatus[statusStr]
	if !ok {
		return fmt.Errorf("invalid status: %s", statusStr)
	}

	*s = status
	return nil
}
