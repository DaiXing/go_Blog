package blogx

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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

// post 发送json，接收json
func PostJson[RespBean any](
	url string,
	reqBody any,
	headers map[string]string,
) (status int, resp *RespBean) {
	bufx := "PostJson >> "
	bufx += "\n url = " + url
	defer func() { // 最后执行。
		bufx += "\n"
		fmt.Println(bufx)
	}()

	// 转成流。
	jsonx := ToJsonString(reqBody)
	bufx += "\n 请求json = " + jsonx
	bodyx := strings.NewReader(jsonx)
	reqx, err := http.NewRequest("POST", url, bodyx)
	CheckErr("NewRequest", err)

	// 设置默认Content-Type
	reqx.Header.Set("Content-Type", "application/json")

	// 自定义header
	if headers != nil {
		for k, v := range headers {
			reqx.Header.Set(k, v)
		}
	}

	respx, err2 := http.DefaultClient.Do(reqx)
	CheckErr("http.DefaultClient.Do", err2)
	defer respx.Body.Close()

	statusx := respx.StatusCode
	bytex, err3 := io.ReadAll(respx.Body)
	CheckErr("ReadAll", err3)
	str := string(bytex)
	bufx += "\n 响应json = " + str

	// json2 := string(bytex)
	var resp2 RespBean
	err5 := json.Unmarshal(bytex, &resp2)
	CheckErr("ReadAll", err5)

	// respJson := ToJsonString(&resp2)
	return statusx, &resp2
}

// post 发送json，接收json。带上 token
func PostJsonWithToken[RespBean any](
	url string,
	reqBody any,
	token string,
) (status int, resp *RespBean) {
	headers := make(map[string]string)
	headers[Key_usertoken] = token
	return PostJson[RespBean](url, reqBody, headers)
}

// get 接收json
func GetJson[RespBean any](url string) *RespBean {
	resp, err := http.Get(url)
	CheckErr("http.Get", err)

	bufx := "GetJson >> "
	bufx += "\n url = " + url
	defer func() {
		bufx += "\n"
		fmt.Println(bufx)
	}()

	bytex, err2 := io.ReadAll(resp.Body)
	CheckErr("ReadAll", err2)
	bufx += "\n 响应json = " + string(bytex)

	var respx RespBean
	err3 := json.Unmarshal(bytex, &respx)
	CheckErr("Unmarshal", err3)
	// respJson := ToJsonString(&respx)
	return &respx
}

// 返回错误消息。
func WebReturnErrorMsg(ctx *gin.Context, status int, msg string) {
	ctx.JSON(status, &PxBaseResp{
		Error: msg,
	})
}

// 请求错误。
func WebBadRequest(ctx *gin.Context, errMsg string) {
	WebReturnErrorMsg(ctx, http.StatusBadRequest, errMsg)
}

// 禁止。
func WebForbidden(ctx *gin.Context, errMsg string) {
	WebReturnErrorMsg(ctx, http.StatusForbidden, errMsg)
}

// 未认证
func WebUnauthorized(ctx *gin.Context, errMsg string) {
	WebReturnErrorMsg(ctx, http.StatusUnauthorized, errMsg)
}

// 404
func WebNotFound(ctx *gin.Context, errMsg string) {
	WebReturnErrorMsg(ctx, http.StatusNotFound, errMsg)
}
