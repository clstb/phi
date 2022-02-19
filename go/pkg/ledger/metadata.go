package ledger

import "fmt"

type Metadata struct {
	Key   string
	Value string
}

func parseMetadata(s string) Metadata {
	matches := MetadataRE.FindStringSubmatch(s)

	return Metadata{
		Key:   matches[1],
		Value: matches[2],
	}
}

func (m Metadata) String() string {
	return fmt.Sprintf(
		"  %s: %s\n",
		m.Key,
		m.Value,
	)
}
