// tools 工具类,日志工具
package tools

import "fmt"

//Logger 日志
type Logger struct{}

const (
	logDebug   = "DEBUG"
	logError   = "ERROR"
	logInfo    = "INFO"
	logWarring = "WARRING"
)

func (log *Logger) writeLog(level, logInfo string) {
	fmt.Printf("[%s][%s]%s\n", GetNow(), level, logInfo)
}

//Debug 调试日志
func (log *Logger) Debug(value string) {
	log.writeLog(logDebug, value)
}

//Debugf 格式化输出
func (log *Logger) Debugf(formater string, value ...interface{}) {
	str := fmt.Sprintf(formater, value...)
	log.Debug(str)
}

// Info 信息日志
func (log *Logger) Info(value string) {
	log.writeLog(logInfo, value)
}

// Infof 格式化输出
func (log *Logger) Infof(formater string, value ...interface{}) {
	str := fmt.Sprintf(formater, value...)
	log.Info(str)
}

//Error 错误日志
func (log *Logger) Error(value string) {
	log.writeLog(logError, value)
}

//Errorf 格式化输出
func (log *Logger) Errorf(formater string, value ...interface{}) {
	str := fmt.Sprintf(formater, value...)
	log.Error(str)
}

// Warring  警告日志
func (log *Logger) Warring(value string) {
	log.writeLog(logWarring, value)
}

// Warringf 格式化输出
func (log *Logger) Warringf(formater string, value ...interface{}) {
	str := fmt.Sprintf(formater, value...)
	log.Warring(str)
}
