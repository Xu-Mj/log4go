# log4go
根据官方的log包进行修改。由于官方的log包在配置了log输出文件后就无法继续在控制台进行输出，
所以在官方包的基础之上添加了对二者的支持，可以同时输出到日志文件以及在控制台进行打印。
而且支持日志级别配置。
目前支持四个级别：DEBUG、INFO、WARN、ERROR。默认的输出级别DEBUG
## 配置项
1. 输出流(io.Writer)-->默认为nil
2. 日志级别(Level uint8)-->默认为DEBUG(0)
3. 日志时间日期格式(string)-->默认为"2006-01-02 15:04:05"
4. 是否打印日志到控制台(bool)-->默认为true
## 使用方式
1. 以公共方法的方式直接调用(包名.方法名)
- 全部采用默认配置
    ```go
    log4go.DEBUG("this is debug log")
    ```
  ```log
  [2021-03-12 13:31:56] [DEBUG] main.go:7: this is debug log
  ```
  
- 指定输出文件
  ```go
  logFile, err := os.OpenFile("logs/error.log",os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
  if err != nil {
  	panic(err.Error())
  }
  log4go.SetOutput(logFile)
  log4go.Debug("this is debug log")
  ```

  控制台和/logs/error.log中都会打印：

  ```log
  [2021-03-12 13:40:06] [DEBUG] main.go:15: this is debug log
  ```

- 指定日志级别

    ```go
    logFile, err := os.OpenFile("logs/error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
    	panic(err.Error())
    }
    log4go.SetOutput(logFile)
    log4go.SetLevel(log.INFO)
    log4go.Debug("this is debug log")
    log4go.Info("this is info log")
    ```

    控制台打印

    ```log
    [2021-03-12 13:47:05] [DEBUG] main.go:16: this is debug log
    [2021-03-12 13:47:05] [INFO] main.go:17: this is info log
    ```

    /logs/error.log文件

    ```log
    [2021-03-12 13:47:05] [INFO] main.go:17: this is info log
    ```

- 设置时间格式

    ```go
    log4go.SetDateFormat("2006 01 02 15:04:05")
    log4go.Debug("this is debug log")
    log4go.Info("this is info log")
    ```

    ```log
    [2021 03 12 13:51:54] [DEBUG] main.go:17: this is debug log
    [2021 03 12 13:51:54] [INFO] main.go:18: this is info log
    ```

- 关闭控制台打印

    ```go
    log4go.SetFlag(false)
    ```

2. 使用创建对象的方式，直接使用New方法创建同时配置Logger对象，原理相同

   ```go
   var logger log4go.Logger
   logger = log4go.New(nil, "", log.INFO, true)
   logger.Info("this is info log")
   ```

   