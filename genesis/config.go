package genesis

import (
	"github.com/spf13/viper"
)

var (
	RetryCount    int64
	SubFilePrefix string
)

func InitConfigInfo() {
	RetryCount = viper.GetInt64("Defaul.RetryCount")
	SubFilePrefix = viper.GetString("SubFile.Prefix")
}
