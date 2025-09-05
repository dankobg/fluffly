package dbcustom

import (
	"database/sql/driver"
	"encoding/json"
)

type JsonType map[string]any

func (m *JsonType) Scan(value any) error {
	return json.Unmarshal(value.([]byte), m)
}

func (m JsonType) Value() (driver.Value, error) {
	return json.Marshal(m)
}
