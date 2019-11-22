// tools 工具类,日志工具
package tools

import (
	"fmt"
	"strings"
)

//Logger 日志
type Logger struct {
	logLevelMap map[string]string
}

const (
	logDebug   = "DEBUG"
	logInfo    = "INFO"
	logWarring = "WARRING"
	logError   = "ERROR"
)

func (log *Logger) readConfig() {
	configInfo, err := ReadConfigFile("./config.yml")
	var loglevel string
	log.logLevelMap = make(map[string]string, 4)
	if nil != err {
		log.Errorf("读配置文件错误,错误信息:%v", err)
		loglevel = logInfo
	} else {
		loglevel = configInfo["log.loglevel"].(string)
	}
	loglevel = strings.ToUpper(loglevel)
	log.logLevelMap[logDebug] = logDebug
	log.logLevelMap[logInfo] = logInfo
	log.logLevelMap[logWarring] = logWarring
	log.logLevelMap[logError] = logError
	if loglevel == logInfo {
		delete(log.logLevelMap, logDebug)
	} else if loglevel == logWarring {
		delete(log.logLevelMap, logDebug)
		delete(log.logLevelMap, logInfo)
	} else if loglevel == logError {
		delete(log.logLevelMap, logDebug)
		delete(log.logLevelMap, logInfo)
		delete(log.logLevelMap, logWarring)
	}

}

func (log *Logger) writeLog(level, logInfo string) {
	if log.logLevelMap == nil || len(log.logLevelMap) == 0 {
		log.readConfig()
	}
	if _, ok := log.logLevelMap[level]; ok == true {
		fmt.Printf("[%s][%s]%s\n", GetNow(), level, logInfo)
	}
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
