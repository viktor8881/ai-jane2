package utils

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

func ReadFileAndUnmarshal(filename string, out interface{}) error {
	extension := strings.ToLower(filepath.Ext(filename))

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	switch extension {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, out)
	case ".json":
		err = json.Unmarshal(data, out)
	default:
		return fmt.Errorf("unsupported file extension: %s", extension)
	}

	return err
}
