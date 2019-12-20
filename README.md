# system_detect

go code

## 系统功能

  本代码采用go语言编写  
  获取系统的相关信息并通过ActiveMQ发送出去;  
  目前系统中获取了CPU使用情况,内存使用情况,进程情况
  > 第一版实现了linux上的功能

## 发送格式  

```
{
    "H":
          {
               "F":"10001",
               "V":"V1.0.1",
               "E":"false",
               "K":"",
               "R":"true",
               "M":"0",
               "S":"1",
               "T":"0",
               "I":"0"
          },
    "B":  
    [  
        {
         "ServerIP":"192.168.0.100",</br> 
            "ServerName":"DatabaseServer",  
            "CollectTime":"2019-08-08 13:14:21",  
            "ProcessCount":10,  
            "CPU":10,  
            "MEM":13,  
            "ThreadCount":10,  
            "HandlerCount":10,  
            "NetWork":10,  
            "DiskFreeSpace":  
            [  
                {  
                    "Path":"/",  
                    "TotalSpace":"1024GB",  
                    "FreeSpace":"240GB"  
                }  
            ],  
            "ProcessInfo":  
            [  
                {  
                    "PID":"1542",  
                    "ProcessName":"wps",   
                    "ProcessPath":"c:\programe\wps.exe",  
                    "CPU":10,  
                    "MEM":13,  
                    "ThreadCount":10,  
                    "HandlerCount":10,  
                    "NetWork":10,  
                    "PortInfo":"80,8080"  
                }  
            ]  
        }  
    ] 
} 
```
## 格式说明
### H部分说明
|标识符|功能|说明|
|:----|:-----|:------|
|F|请求功能码|请求的唯一功能编号|
|E|是否加密|true加密数据,false不加密数据,如果是加密传输则将B部分综合为一个属性(Data)|
|V|协议版本|1.0.1(第一个协议版本号)|
|K|Token|服务器返回的授权识别号,可根据实际情况考虑是否需要|
|R|执行结果|true成功,否则为失败(当被动响应时传入)|
|M|返回的消息|执行失败之后返回的消息(被动响应时传入)|
|S|顺序号|编号从1开始直到FFFFFFFF结束,超过之后回到1,主动方负责生成,被动方原样返回|
|T|总批次|当传输记录超过100条的时候需要进行分批传输;0表示未分批;|
|I|批次号|当出现多批次传输数据是需要|

### B部分说明
| 参数名称| 中文名 |必须| 说明 |
| :------------|:-----------|:------|:---------|
|ServerIP|服务器IP|是||
|ServerName|服务器名称|是||
|CollectTime|数据采集时间|是||
|ProcessCount|当前进程数量|是||
|CPU|当前服务器CPU使用率|是||
|MEM|当前服务器内存使用|是||
|ThreadCount|当前服务器线程总数量|是||
|HandlerCount|当前服务器句柄总数量|是||
|NetWork|当前服务器网络使用率|是||
|DiskFreeSpace|硬盘空间|是|对象数组|
|Path|硬盘路径|是||
|TotalSpace|总空间|是||
|FreeSpace|剩余空间|是||
|ProcessInfo|进程信息|是|对象数组|
|PID|进程PID|是||
|ProcessName|进程名称|是||
|ProcessPath|进程路径|是||
|CPU|进程CPU使用率|是||
|MEM|进程内存使用率|是||
|ThreadCount|进程线程数量|是||
|HandlerCount|进程句柄数量|是||
|NetWork|进程网络使用率|是||
|PortInfo|进程打开的端口|是|多个端口逗号分开|

# 配置说明
health_config.yml文件放置在同级目录即可.才用yml格式配置.
