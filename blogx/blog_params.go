package blogx

import (
	"fmt"

	"github.com/spf13/viper"
)

// 参数。
var ConfigParams ConfigParamsPo

// 加载配置项。
func LoadParams() {
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

// 配置项。
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
