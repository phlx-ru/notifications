package conf

import (
	"bytes"
	"os"
	"regexp"

	"github.com/go-kratos/kratos/v2/config"
	"gopkg.in/yaml.v3"
)

const (
	envKeywordRegex = `\${((.+?)(:.+?)*?)}`
)

func EnvDecoder(kv *config.KeyValue, v map[string]any) error {
	configData := replaceEnv(kv.Value)
	return yaml.Unmarshal(configData, v)
}

func replaceEnv(configData []byte) []byte {
	for _, match := range regexp.MustCompile(envKeywordRegex).FindAllSubmatch(configData, -1) {
		key := string(match[2])
		value := os.Getenv(key)
		if value != "" {
			configData = bytes.Replace(configData, match[0], []byte(value), 1)
		}
	}

	return configData
}
