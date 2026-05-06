package blogx

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

// 日志。
var Logger *slog.Logger

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
		fmt.Printf("  key= %s value= %s \n", k, v)
	}
}

func CheckErr(title string, err2 error) {
	if err2 != nil {
		panic(title + " error : " + err2.Error())
	}
}

func ToJsonString(anyx any) string {
	bytex, err := json.Marshal(anyx)
	CheckErr("转json", err)
	return string(bytex)
}

// 密码加个密。
func PasswordEncode(password string) string {
	salt := "86571JFJWELSA" // 加盐。
	tmp := password + salt
	bytex := sha256.Sum256([]byte(tmp))                // 定长的字节数组
	str := base64.StdEncoding.EncodeToString(bytex[:]) // 转base64
	return str
}
