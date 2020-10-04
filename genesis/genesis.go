package genesis

import "github.com/spf13/viper"

//初始化配置文件，将配置文件中的内容读取到map中
func InitConfig() {
	// 设置配置文件名称
	viper.SetConfigName("config_down")
	// 设置配置文件与可执行二进制文件(main.go，因为编译生成的可执行二进制文件与main.go在同一路径)的相对路径
	viper.AddConfigPath("./conf")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			Logger.Println("no such config file")
		} else {
			// Config file was found but another error was produced
			Logger.Println("read config error")
		}
		Logger.Fatal(err) // 读取配置文件失败致命错误
	}
}
