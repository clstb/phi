package parser

import (
	"fmt"
	"regexp"
)

var MetadataFieldRE = regexp.MustCompile(`^  ([a-z].*): (\S*)$`)

type MetadataField struct {
	Key   string
	Value string
}

func parseMetadata(s string) MetadataField {
	matches := MetadataFieldRE.FindStringSubmatch(s)

	return MetadataField{
		Key:   matches[1],
		Value: matches[2],
	}
}

func (m MetadataField) String() string {
	return fmt.Sprintf(
		"  %s: %s\n",
		m.Key,
		m.Value,
	)
}
