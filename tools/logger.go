// tools 工具类,日志工具
package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//Logger 日志
type Logger struct {
	logLevelMap map[string]string
}

var logPath string
var showlog bool

const (
	logDebug   = "DEBUG"
	logInfo    = "INFO"
	logWarring = "WARRING"
	logError   = "ERROR"
)

func (log *Logger) readConfig() {
	configInfo, err := ReadConfigFile("./health_config.yml")
	var loglevel string
	if nil != err {
		log.Errorf("读配置文件错误,错误信息:%v", err)
		loglevel = logInfo
		showlog = false
	} else {
		loglevel = configInfo["log.loglevel"].(string)
		showlog = configInfo["log.showlog"].(bool)
	}
	log.logLevelMap = make(map[string]string, 4)

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
		str := fmt.Sprintf("[%s][%s]%s\n", GetNow(), level, logInfo)
		baseDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			baseDir = "./log/"
		} else {
			baseDir = fmt.Sprintf("%s/log/", baseDir)
		}
		if ok := PathExists(baseDir); ok == false {
			err = os.Mkdir(baseDir, os.ModePerm)
			if err != nil {
				fmt.Printf("创建文件夹失败,错误信息:%v", err)
				return
			}
		}
		logPath = fmt.Sprintf("%s%s_%s.log", baseDir, level, time.Now().Format("2006-01-02_15"))
		data := []byte(str)
		if err = ioutil.WriteFile(logPath, data, 0644); err != nil {
			fmt.Printf("写日志失败,错误信息:%v", err)
		}
		if showlog {
			fmt.Println(str)
		}
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
