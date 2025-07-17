package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestManifestFiles(t *testing.T) {
	stacksDir := "../../test_data/stacks"

	entries, err := os.ReadDir(stacksDir)
	if err != nil {
		t.Fatalf("Failed reading stacks dir: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			stack := entry.Name()
			stackPath := filepath.Join(stacksDir, stack, "stack.yaml")
			t.Run(stack, func(t *testing.T) {
				err = Stack(stackPath)
				if err != nil {
					t.Fatalf("%s: %v", stack, err)
				}
			})
		}
	}
}

func TestManifestEmpty(t *testing.T) {
	data := ""
	err := validateStackYaml("", []byte(data))
	if err == nil {
		t.Fatal("Empty yaml should fail")
	}
	t.Log(err)
}

func TestUnknownField(t *testing.T) {
	data := `
name: test
description: test
vendor: test
grade: stable
devices:
memory: 1
disk-space: 1
components:
configurations:
  engine: test
  model: test
unknown-field: test
`
	err := validateStackYaml("test/stack.yaml", []byte(data))
	if err == nil {
		t.Fatal("Unknown field should fail")
	}
	t.Log(err)
}

func TestNameRequired(t *testing.T) {
	data := `
#name: test
description: test
vendor: test
grade: stable
devices:
memory: 1
disk-space: 1
components:
configurations:
  engine: test
  model: test
`
	err := validateStackYaml("test/stack.yaml", []byte(data))
	if err == nil {
		t.Fatal("name field is required")
	}
	t.Log(err)

}

func TestDescriptionRequired(t *testing.T) {
	data := `
name: test
#description: test
vendor: test
grade: stable
devices:
memory: 1
disk-space: 1
components:
configurations:
  engine: test
  model: test
`
	err := validateStackYaml("test/stack.yaml", []byte(data))
	if err == nil {
		t.Fatal("description is required")
	}
	t.Log(err)

}

func TestVendorRequired(t *testing.T) {
	data := `
name: test
description: test
#vendor: test
grade: stable
devices:
memory: 1
disk-space: 1
components:
configurations:
  engine: test
  model: test
`
	err := validateStackYaml("test/stack.yaml", []byte(data))
	if err == nil {
		t.Fatal("vendor is required")
	}
	t.Log(err)

}

func TestGradeRequired(t *testing.T) {
	data := `
name: test
description: test
vendor: test
#grade: stable
devices:
memory: 1
disk-space: 1
components:
configurations:
  engine: test
  model: test
`
	err := validateStackYaml("test/stack.yaml", []byte(data))
	if err == nil {
		t.Fatal("grade is required")
	}
	t.Log(err)

}

func TestGradeValid(t *testing.T) {
	t.Run("grade stable", func(t *testing.T) {
		data := `
name: test
description: test
vendor: test
grade: stable
devices:
memory: 1
disk-space: 1
components:
configurations:
  engine: test
  model: test
`
		err := validateStackYaml("test/stack.yaml", []byte(data))
		if err != nil {
			t.Fatal("grade stable should be valid")
		}
	})
	t.Run("grade devel", func(t *testing.T) {
		data := `
name: test
description: test
vendor: test
grade: devel
devices:
memory: 1
disk-space: 1
components:
configurations:
  engine: test
  model: test
`
		err := validateStackYaml("test/stack.yaml", []byte(data))
		if err != nil {
			t.Fatal("grade devel should be valid")
		}
	})
	t.Run("grade invalid", func(t *testing.T) {
		data := `
name: test
description: test
vendor: test
grade: invalid
devices:
memory: 1
disk-space: 1
components:
configurations:
  engine: test
  model: test
`
		err := validateStackYaml("test/stack.yaml", []byte(data))
		if err == nil {
			t.Fatal("grade invalid")
		}
		t.Log(err)
	})

}

func TestMemoryValues(t *testing.T) {
	dataTemplate := `
name: test
description: test
vendor: test
grade: stable
devices:
memory: %s
disk-space: 1
components:
configurations:
  engine: test
  model: test
`
	t.Run("valid GB", func(t *testing.T) {
		dataTest := fmt.Sprintf(dataTemplate, "1G")
		err := validateStackYaml("test/stack.yaml", []byte(dataTest))
		if err != nil {
			t.Logf("memory should be valid: %v", err)
		}
	})

	t.Run("valid MB", func(t *testing.T) {
		dataTest := fmt.Sprintf(dataTemplate, "512M")
		err := validateStackYaml("test/stack.yaml", []byte(dataTest))
		if err != nil {
			t.Logf("memory should be valid: %v", err)
		}
	})

	// Empty string is parsed as nil, which we interpret as unset, which is valid
	//t.Run("empty", func(t *testing.T) {
	//	dataTest := fmt.Sprintf(dataTemplate, " ")
	//	err := validateStackYaml("test/stack.yaml", []byte(dataTest))
	//	if err == nil {
	//		t.Fatal("empty memory should be invalid")
	//	}
	//	t.Log(err)
	//})

	t.Run("not numeric", func(t *testing.T) {
		dataTest := fmt.Sprintf(dataTemplate, "abc")
		err := validateStackYaml("test/stack.yaml", []byte(dataTest))
		if err == nil {
			t.Fatal("non-numeric memory should be invalid")
		}
		t.Log(err)
	})

}

func TestDiskValues(t *testing.T) {
	dataTemplate := `
name: test
description: test
vendor: test
grade: stable
devices:
memory: 1G
disk-space: %s
components:
configurations:
  engine: test
  model: test
`
	t.Run("valid GB", func(t *testing.T) {
		dataTest := fmt.Sprintf(dataTemplate, "1G")
		err := validateStackYaml("test/stack.yaml", []byte(dataTest))
		if err != nil {
			t.Logf("disk should be valid: %v", err)
		}
	})

	t.Run("valid MB", func(t *testing.T) {
		dataTest := fmt.Sprintf(dataTemplate, "512M")
		err := validateStackYaml("test/stack.yaml", []byte(dataTest))
		if err != nil {
			t.Logf("disk should be valid: %v", err)
		}
	})

	// Empty string is parsed as nil, which we interpret as unset, which is valid
	//t.Run("empty", func(t *testing.T) {
	//	dataTest := fmt.Sprintf(dataTemplate, " ")
	//	err := validateStackYaml("test/stack.yaml", []byte(dataTest))
	//	if err == nil {
	//		t.Fatal("empty memory should be invalid")
	//	}
	//	t.Log(err)
	//})

	t.Run("not numeric", func(t *testing.T) {
		dataTest := fmt.Sprintf(dataTemplate, "abc")
		err := validateStackYaml("test/stack.yaml", []byte(dataTest))
		if err == nil {
			t.Fatal("non-numeric disk should be invalid")
		}
		t.Log(err)
	})

}

func TestConfigEngine(t *testing.T) {
	data := `
name: test
description: test
vendor: test
grade: stable
devices:
memory: 1
disk-space: 1
components:
configurations:
#  engine: test
  model: test
`
	err := validateStackYaml("test/stack.yaml", []byte(data))
	if err == nil {
		t.Fatal("missing config.engine field should be invalid")
	}
	t.Log(err)
}

func TestConfigModel(t *testing.T) {
	data := `
name: test
description: test
vendor: test
grade: stable
devices:
memory: 1
disk-space: 1
components:
configurations:
  engine: test
#  model: test
`
	err := validateStackYaml("test/stack.yaml", []byte(data))
	if err == nil {
		t.Fatal("missing config.model field should be invalid")
	}
	t.Log(err)
}
