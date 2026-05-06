package blogx

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

// 日志。
var Logger *slog.Logger

// 参数。
var ConfigParams ConfigParamsPo

// 配置日志。
func ConfigLogger() {
	// 配置 JSON 格式，开启 caller 信息便于定位行号
	opts := &slog.HandlerOptions{
		AddSource: true,           // 开启 caller
		Level:     slog.LevelInfo, // 生产环境设置为 Info 级别
	}
	// handler := slog.NewJSONHandler(os.Stdout, opts)// json格式。
	handler := slog.NewTextHandler(os.Stdout, opts) // 文本格式
	logger := slog.New(handler)

	// 替换全局默认 logger
	slog.SetDefault(logger)
	Logger = logger
	Logger.Info("logger 初始化完成")
}

// 加载配置项。
func ConfigLoadParams() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.SetConfigName("blog")
	viper.SetConfigType("yaml")
	err1 := viper.ReadInConfig()
	CheckErr("ReadInConfig", err1)
	Logger.Info("yaml 读配置完成")

	fmt.Println("viper 的全部参数：")
	keys := viper.AllKeys()
	for _, k := range keys {
		v := viper.GetString(k)
		fmt.Printf("  key= %-30s value= %-30s \n", k, v)
	}

	// 转对象。
	err3 := viper.Unmarshal(&ConfigParams)
	CheckErr("viper.Unmarshal", err3)
	fmt.Println("viper 对象参数： ", ToJsonString(&ConfigParams))
}

type ConfigParamsPo struct {
	Datasource ConfigDataSourcePo `mapstructure:"datasource"`
	Init       ConfigInitPo       `mapstructure:"init"`
}
type ConfigDataSourcePo struct {
	Driver      string `mapstructure:"driver"`
	Connection  string `mapstructure:"connection"`
	ConnTimeout uint   `mapstructure:"connTimeout"`
}
type ConfigInitPo struct {
	DbDropTableEnable  bool `mapstructure:"dbDropTableEnable"`
	DbInsertRowsEnable bool `mapstructure:"dbInsertRowsEnable"`
}
