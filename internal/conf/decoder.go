package conf

import (
	"bytes"
	"os"
	"regexp"

	"github.com/go-kratos/kratos/v2/config"
	"gopkg.in/yaml.v3"
)

func EnvDecoder(kv *config.KeyValue, v map[string]any) error {
	configData := replaceEnv(kv.Value)
	return yaml.Unmarshal(configData, v)
}

func replaceEnv(configData []byte) []byte {
	r, _ := regexp.Compile(`\${((.+?)(:\w+)*?)}`)

	for _, match := range r.FindAllSubmatch(configData, -1) {
		key := string(match[2])
		value := os.Getenv(key)
		if value != "" {
			configData = bytes.Replace(configData, match[0], []byte(value), 1)
		}
	}

	return configData
}
