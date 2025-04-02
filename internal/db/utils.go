package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONB map[string]interface{}

// Implement `driver.Valuer` for saving the JSON field
func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Implement `sql.Scanner` for reading the JSON field
func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}
