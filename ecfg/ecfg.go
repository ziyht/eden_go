package ecfg

import (
	"fmt"
	"os"
	"path/filepath"
)

func ParsingFromCfgFile(path string, key string, dest interface{}) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("get abs path failed for '%s':\n %s", path, err)
	}

	ext := filepath.Ext(path)
	if len(ext) > 1 {
		ext = ext[1:]
	} else {
		return fmt.Errorf("can not found ext in file '%s' like .yml .ini .toml", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read '%s' failed: %s", path, err)
	}

	err = parsingCfgFromStr(string(data), ext, key, dest)
	if err != nil {
		return fmt.Errorf("parsing failed for path '%s': %s", path, err)
	}

	return nil
}

