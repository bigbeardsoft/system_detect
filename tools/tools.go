package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fwhezfwhez/errorx"
	"github.com/spf13/viper"
)

// IsNumeric 判断字符串是否为数字
func IsNumeric(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
	case float32, float64, complex64, complex128:
		return true
	case string:
		str := val.(string)
		if str == "" {
			return false
		}
		// Trim any whitespace
		str = strings.Trim(str, " \\t\\n\\r\\v\\f")
		if str[0] == '-' || str[0] == '+' {
			if len(str) == 1 {
				return false
			}
			str = str[1:]
		}
		// hex
		if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
			for _, h := range str[2:] {
				if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
					return false
				}
			}
			return true
		}
		// 0-9,Point,Scientific
		p, s, l := 0, 0, len(str)
		for i, v := range str {
			if v == '.' { // Point
				if p > 0 || s > 0 || i+1 == l {
					return false
				}
				p = i
			} else if v == 'e' || v == 'E' { // Scientific
				if i == 0 || s > 0 || i+1 == l {
					return false
				}
				s = i
			} else if v < '0' || v > '9' {
				return false
			}
		}
		return true
	}

	return false
}

var logger Logger

// ExecuteCommand 执行linux命令
func ExecuteCommand(cmdString string) (string, error) {
	cmdwithpath, err := exec.LookPath("bash")
	if err != nil {
		logger.Error("not find bash.")
		return "", errorx.New(err)
	}

	cmd := exec.Command(cmdwithpath, "-c", cmdString)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errorx.New(fmt.Errorf("执行命令:%s,发生错误,错误信息:%v", cmdString, err))
	}
	s := string(out)
	return s, nil
}

// ConvertToJSON 结构体转成json
func ConvertToJSON(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	jsonStr := string(b)
	return jsonStr, err
}

// ReadConfigFile 读配置文件
func ReadConfigFile(path string) (map[string]interface{}, error) {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if nil != err {
		return nil, errorx.New(err)
	}
	allkey := viper.AllKeys()

	allmap := make(map[string]interface{}, len(allkey))

	for _, key := range allkey {
		allmap[key] = viper.Get(key)
	}
	return allmap, nil
}

// GetNow 获取当期日期和时间,格式:yyyy-MM-dd hh:mm:ss
func GetNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//  PathExists 判断文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
