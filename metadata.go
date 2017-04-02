package spirali

import (
	"bytes"
	"encoding/json"
	"io"
)

// MetaDataFileName ...
const MetaDataFileName = "metadata.json"

// MetaData of migration.
type MetaData struct {
	Migrations Migrations `json:"migrations"`
}

// ReadMetaData reads the metadata of migration from io.Reader.
func ReadMetaData(r io.Reader) (*MetaData, error) {
	decoder := json.NewDecoder(r)
	var m MetaData
	if err := decoder.Decode(&m); err != nil {
		return nil, err
	}
	m.Migrations.sort()
	return &m, nil
}

// Save updates the metadata file.
func (m *MetaData) Save(w io.Writer) error {
	b, err := json.Marshal(&m)
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	if err := json.Indent(&buffer, b, "", "  "); err != nil {
		return err
	}
	if _, err := buffer.WriteTo(w); err != nil {
		return err
	}
	return nil
}
