package ecfg

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func parsingCfgFromStr(content string, ext string, key string, dest interface{})(err error){

	v := viper.New()
	v.SetConfigType(ext)
	if err := v.ReadConfig(strings.NewReader(content)); err != nil {
		return fmt.Errorf("parsing failed: %s", err)
	}

	if key == "" {
		err = v.Unmarshal(dest)
	} else {
		err = v.UnmarshalKey(key, dest)
	}
	if err != nil {
		return fmt.Errorf("unmarshaling failed: %s", err)
	}

	return  nil
}
