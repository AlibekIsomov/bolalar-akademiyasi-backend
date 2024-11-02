package models

import (
	"encoding/json"
	"fmt"
)

type Source string

const (
	Telegram  Source = "telegram"
	Website   Source = "website"
	Telephone Source = "telephone"
	Instagram Source = "instagram"
)

var sourceToString = map[Source]string{
	Telegram:  "telegram",
	Website:   "website",
	Telephone: "telephone",
	Instagram: "instagram",
}

var stringToSource = map[string]Source{
	"telegram":  Telegram,
	"website":   Website,
	"telephone": Telephone,
	"instagram": Instagram,
}

func (s Source) String() string {
	return sourceToString[s]
}

func (s Source) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Source) UnmarshalJSON(data []byte) error {
	var sourceStr string
	if err := json.Unmarshal(data, &sourceStr); err != nil {
		return err
	}

	source, ok := stringToSource[sourceStr]
	if !ok {
		return fmt.Errorf("invalid source: %s", sourceStr)
	}

	*s = source
	return nil
}
