package storage

import (
	"errors"
	"fmt"
	"testing"
)

// migrationMockStorage is a minimal storage backend for testing migrations.
type migrationMockStorage struct {
	data     map[string]map[string]any
	getErr   error
	setErr   error
	unsetErr error
	sets     map[string]string
	unsets   []string
}

func newMigrationMockStorage() *migrationMockStorage {
	return &migrationMockStorage{
		data: make(map[string]map[string]any),
		sets: make(map[string]string),
	}
}

func (m *migrationMockStorage) Get(key string) (map[string]any, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	val, ok := m.data[key]
	if !ok {
		return nil, ErrorNotFound
	}
	return val, nil
}

func (m *migrationMockStorage) Set(key, value string) error {
	if m.setErr != nil {
		return m.setErr
	}
	m.sets[key] = value
	return nil
}

func (m *migrationMockStorage) SetDocument(key string, value any) error { return nil }

func (m *migrationMockStorage) Unset(key string) error {
	if m.unsetErr != nil {
		return m.unsetErr
	}
	m.unsets = append(m.unsets, key)
	return nil
}

//
//	Test cases
//

func TestMigrate_RunsAllMigrations(t *testing.T) {
	ms := newMigrationMockStorage()
	c := &config{storage: ms}
	ms.data[c.nestKeys(UserConfig, "passthrough.env")] = map[string]any{
		"MY_VAR": "hello",
	}

	if err := c.Migrate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMigrate_PropagatesError(t *testing.T) {
	ms := newMigrationMockStorage()
	ms.getErr = fmt.Errorf("storage unavailable")
	c := &config{storage: ms}

	if err := c.Migrate(); err == nil {
		t.Fatal("expected error from Migrate, got nil")
	}
}

func TestMigratePassthroughEnv_NoConfig(t *testing.T) {
	c := &config{storage: newMigrationMockStorage()}

	if err := c.migratePassthroughEnv(); err != nil {
		t.Fatalf("expected no error when no config exists, got %v", err)
	}
}

func TestMigratePassthroughEnv_MigratesVariables(t *testing.T) {
	ms := newMigrationMockStorage()
	c := &config{storage: ms}
	ms.data[c.nestKeys(UserConfig, "passthrough.env")] = map[string]any{
		"MY_VAR":    "hello",
		"OTHER_VAR": "world",
	}

	if err := c.migratePassthroughEnv(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify new env keys were set
	for varName, wantVal := range map[string]string{"MY_VAR": "hello", "OTHER_VAR": "world"} {
		newKey := c.nestKeys(UserConfig, "env."+varName)
		if got, ok := ms.sets[newKey]; !ok {
			t.Errorf("expected Set(%q) to be called", newKey)
		} else if got != wantVal {
			t.Errorf("Set(%q): got %q, want %q", newKey, got, wantVal)
		}
	}

	// Verify old passthrough keys were unset
	unsetSet := make(map[string]bool, len(ms.unsets))
	for _, u := range ms.unsets {
		unsetSet[u] = true
	}
	for _, varName := range []string{"MY_VAR", "OTHER_VAR"} {
		oldKey := c.nestKeys(UserConfig, "passthrough.env."+varName)
		if !unsetSet[oldKey] {
			t.Errorf("expected Unset(%q) to be called", oldKey)
		}
	}
}

func TestMigratePassthroughEnv_StorageGetError(t *testing.T) {
	ms := newMigrationMockStorage()
	ms.getErr = fmt.Errorf("storage unavailable")
	c := &config{storage: ms}

	err := c.migratePassthroughEnv()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if errors.Is(err, ErrorNotFound) {
		t.Fatal("expected non-NotFound error to propagate")
	}
}

func TestMigratePassthroughEnv_SetError(t *testing.T) {
	ms := newMigrationMockStorage()
	c := &config{storage: ms}
	ms.data[c.nestKeys(UserConfig, "passthrough.env")] = map[string]any{
		"MY_VAR": "hello",
	}
	ms.setErr = fmt.Errorf("set failed")

	if err := c.migratePassthroughEnv(); err == nil {
		t.Fatal("expected error from Set, got nil")
	}
}

func TestMigratePassthroughEnv_UnsetError(t *testing.T) {
	ms := newMigrationMockStorage()
	c := &config{storage: ms}
	ms.data[c.nestKeys(UserConfig, "passthrough.env")] = map[string]any{
		"MY_VAR": "hello",
	}
	ms.unsetErr = fmt.Errorf("unset failed")

	if err := c.migratePassthroughEnv(); err == nil {
		t.Fatal("expected error from Unset, got nil")
	}
}
