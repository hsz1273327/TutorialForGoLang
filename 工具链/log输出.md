# log输出

虽然go有一个log模块,但太过简陋了,现在比较常见的使用[github.com/sirupsen/logrus](https://github.com/sirupsen/logrus)这个包优势在于全面,支持结构化log,支持log分级,也支持写入文件.

当然在使用docker的条件下log输出到文件并不是一个必要的事情.这个部分的例子子我放在[项目的code文件夹下](https://github.com/hsz1273327/TutorialForGoLang/tree/master/%E5%B7%A5%E5%85%B7%E9%93%BE/code/testlog)

我通常会为log建立一个单独的子模块logger,这样的好处是不容易冲突.项目的结构是:

```bash

testlog\
       |-go.mod
       |-main.go
       |-logger\
               |-go.mod
               |-logger.go

```

+ logger.go

    ```golang
    package logger

    import (
        logrus "github.com/sirupsen/logrus"
    )

    func Init() *logrus.Logger {
        log := logrus.New()
        log.SetFormatter(&logrus.JSONFormatter{})
        return log
    }

    //Logger 默认的logger
    var Logger = Init()

    //Log 有默认字段的log
    var Log = Logger.WithFields(logrus.Fields{
        "common": "this is a common field",
    })
    ```

+ main.go

```golang
package main

import (
    "testlog/logger"

    logrus "github.com/sirupsen/logrus"
)

func main() {
    logger.Logger.SetLevel(logrus.InfoLevel)
    logger.Logger.Info("测试")
    logger.Log.Info("测试")
    logger.Logger.WithFields(logrus.Fields{
        "event": "field",
    }).Info("测试 field")
    logger.Log.WithFields(logrus.Fields{
        "event": "field",
    }).Info("测试 field")

    logger.Logger.SetLevel(logrus.WarnLevel)
    logger.Logger.Info("测试 INFO")
    logger.Log.Info("测试 INFO")
    logger.Logger.Warn("测试 warn")
    logger.Log.Warn("测试 warn")
}

```

## 结构化输出

我们只需要设置log对象的`.SetFormatter(&logrus.JSONFormatter{})`就可以输出json格式的log了.而格式化输出中默认的字段有

+ `level`:log等级
+ `time`:log的打印时间
+ `msg`:`Info`等方法中填入的字符串

### 增加字段

我们可以为log对象使用方法

```golang
log.WithFields(logrus.Fields{
        "common": "this is a common field",
    })
```

添加字段,这个新对象我们可以拿它再执行`Info`这类方法来输出log,这样我们就有了两种方法

+ 在需要打log时使用

```golang
log.WithFields(logrus.Fields{
        "common": "this is a common field",
    }).Info("msg")
```

+ 先定义好默认的字段,再在需要打印log时直接使用这个对象来打印

```golang
var logger = log.WithFields(logrus.Fields{
        "common": "this is a common field",
    })

...

logger.Info("msg)
```

## 设置level

我们只要使用log实例的`.SetLevel(logrus.InfoLevel)`方法就可以设置log等级.通常我们会设置在`logrus.WarnLevel`这一级

## 另一个log库

在github上另一个以快为特色的log库是[go.uber.org/zap](https://github.com/uber-go/zap),在对性能有更高要求的场景下建议改用这个

## 使用[loggerhelper](https://github.com/Golang-Tools/loggerhelper)来打log

这部分也是私货,我做了一个`logrus`的代理工具,可以在代码中先使用它打log,后通过`Init`或者`InitWithOutput`接口初始化定义log的具体全局行为.

其用法如下:

```go
package main

import (
    log "github.com/Golang-Tools/loggerhelper"
    "io/ioutil"
    "os"
    "github.com/sirupsen/logrus"
    "github.com/sirupsen/logrus/hooks/writer"
)
func main() {
    hook := writer.Hook{ // Send logs with level higher than warning to stderr
        Writer: os.Stderr,
        LogLevels: []logrus.Level{
            logrus.PanicLevel,
            logrus.FatalLevel,
            logrus.ErrorLevel,
            logrus.WarnLevel,
        },
    }
    log.InitWithOutput("WARN", log.Dict{"d": 3}, ioutil.Discard, &hook)
    log.Info("test")
    log.Warn("qweqwr", log.Dict{"a": 1})
}
```