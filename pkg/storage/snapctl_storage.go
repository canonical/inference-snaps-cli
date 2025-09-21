package storage

import (
	"encoding/json"

	"github.com/canonical/go-snapctl"
)

type SnapctlStorage struct{}

func NewSnapctlStorage() *SnapctlStorage {
	return &SnapctlStorage{}
}

func (s *SnapctlStorage) Set(key, value string) error {
	return snapctl.Set(key, string(value)).Run()
}

func (s *SnapctlStorage) SetDocument(key string, value any) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return snapctl.Set(key, string(b)).Document().Run()
}

func (s *SnapctlStorage) Get(key string) ([]byte, error) {
	val, err := snapctl.Get(key).Run()
	if err != nil {
		return nil, err
	}
	// TODO: query as document to distinguish between empty and not found
	if val == "" {
		return nil, ErrorNotFound
	}
	return []byte(val), nil
}

func (s *SnapctlStorage) Unset(key string) error {
	err := snapctl.Unset(key).Run()
	if err != nil {
		return err
	}
	return nil
}
