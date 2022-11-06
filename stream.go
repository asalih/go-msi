package msi

import "github.com/asalih/go-mscfb"

type Streams struct {
	Entries *mscfb.Entries
}

func NewStreams(entries *mscfb.Entries) *Streams {
	return &Streams{
		Entries: entries,
	}
}

func (s *Streams) Next() string {
	for {
		entry := s.Entries.Next()
		if entry == nil {
			return ""
		}

		if !entry.IsStream() ||
			entry.Name == DIGITAL_SIGNATURE_STREAM_NAME ||
			entry.Name == SUMMARY_INFO_STREAM_NAME ||
			entry.Name == MSI_DIGITAL_SIGNATURE_EX_STREAM_NAME {
			continue
		}

		name, isTable := NameDecode(entry.Name)
		if !isTable {
			return name
		}
	}
}
