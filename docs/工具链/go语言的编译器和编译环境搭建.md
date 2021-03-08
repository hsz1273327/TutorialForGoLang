# go语言的编译器和编译环境搭建

就像c语言需要`gcc`,同样作为静态语言,go语言也需要现有编译器才可以编译源码.

我们可以去[go语言的官网](https://golang.org/doc/install)下载安装,如果和我一样是mac用户,可以借助[Homeberry](http://blog.hszofficial.site/recommend/2016/06/28/mac%E7%9A%84%E5%8C%85%E7%AE%A1%E7%90%86%E5%B7%A5%E5%85%B7homebrew/)安装

```shell
brew install golang
```

装好后并不是直接可以用,go语言有这么几个比较重要的环境变量需要设置

+ `GOPATH` 用于标明GO语言的默认全局编译和安装工作目录,类似python中的`site-packages`,其中会有`src`,`pkg`和`bin`三个目录,分别是项目的源码,库文件和可执行文件.
+ `GOROOT` 用于标明Go语言的安装目录

其他的还有一些我们用到的时候再说,如果要查看当前安装环境下与go语言有关的全部环境变量,可以使用命令
`go env`查看.

同时为了方便起见,也应该向环境变量中的`PATH`中添加`$GOPATH/bin`,这样全局安装的go程序就可以直接执行了.

## 编译和解释

从python,js这类动态语言过来的同学可能有些不理解什么是编译器.在动态语言中有解释器用于解释执行源码,即将源码先翻译为二进制码,然后解释器将这些二进制码通过调用内部的对应虚拟机执行.像java虽然说自己是静态语言,但实质上也是这个流程,只是它是按模块读取二进制码不是逐行读取,并且不允许运行时修改而已.

而类似C语言,Go语言这类则是完全不一样的路数--编译器会先读取源码,做好类型检测,语法检测,然后将其直接编译为机器可以执行的机器码,然后通过链接的方式将不容模块链接合并成动态库,静态库或者可执行文件.

因此python执行总是`python xxxx.py`,js执行总是`node xxx.js`,java执行也总是`java xxxx.jar`,毕竟他们其实是使用解释器解释程序让他在虚拟机中执行;而C,go这类则是编译出个叫`hello`的程序后,直接使用`./hello`就可以执行了,因为编译好的程序是操作系统本身就可以执行的程序.

## helloworld

我们惯例的从一个helloworld开始这一系列.一个典型的go语言项目结构中包括以下几个部分:

1. `go.mod`文件,用于描述项目的依赖关系

    ```text
    module helloworld //描述项目名

    go 1.15 //描述项目最低支持的go语言版本
    ```

2. 入口文件,用于作为程序入口

    ```go
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello, 世界 Golang!")
    }
    ```

### 编译可执行文件

go编译程序使用命令[`go build [-o output] [-i] [build flags] [入口文件]`](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies)

### 执行程序

和c语言一样,我们直接执行编译出来的可执行文件即可.

## 静态库

go语言早期只支持[静态库](https://baike.baidu.com/item/%E9%9D%99%E6%80%81%E5%BA%93/8955694),我们的第二个例子就来构建一个静态库.这个静态库实现[牛顿法求平方根](https://tour.go-zh.org/flowcontrol/8).

go在早期使用`GOPATH`来控制包依赖,现在`GOPATH`已经基本不再被推荐使用,官方现在推的是`go mod`模式,也就是我们上面的模式.但而即便到了现在我们本地构造静态库依然只能使用`GOPATH`模式.这种模式说白了就是把项目全部放在环境变量`GOPATH`指定的根目录下.比如我们要构造一个静态链接库`mymath`,就在`$GOPATH/src/mymath`目录下写源码即可

+ sqrt.go

```go
package mymath

func Sqrt(x float64) float64 {
	z:=1.0
	for i:=0;i<10;i++ {
		z -= (z*z-x)/(2*z)
	}
	return z
}
```

go中判断一个函数或结构体是否模块外可见使用的是首字母来区分,如果是大写则为可见,否则不可见.像上面的代码我们的`Sqrt`函数就是可见的.

我们将这个文件命名为`mymath.go`,放在`$GOPATH/src/my/mymath`目录下.

### 编译静态库

编译静态库一样是两种方式:

1. 在`$GOPATH`下使用`go install my/mymath`命令
2. 在`$GOPATH/src/my/mymath`目录下使用`go install`命令
   
这两种方式都可以在`$GOPATH/pkg/{platform}`目录下的`my`文件夹下生成`mymath.a`静态库.这个`platform`取决于你使用的是什么操作系统什么cpu,比如我是mac,那就是`darwin_amd64`.

### 调用静态库

调用静态库也很简单,只需要使用`import`语句即可

```go
package main

import (
    "fmt"
    mymath "my/mymath"
)

func main() {
	fmt.Println(mymath.Sqrt(2))
}
```

调用可用的静态库,其搜索路劲按顺序依次为

1. 本地`vendor`文件夹(这个在下一节`依赖控制`部分讲)
2. 环境变量`GOPATH`下
3. 环境变量`GOROOT`下

如果我们给库定义了别名,那么我们就可以在后面的代码中使用这个别名代表这个库,这就有点像python中的`import xxx as yyy`

### go mod模式下的"本地静态库"

如果我们希望本地构造静态库,又希望不用`GOPATH`而是使用`go mod`,那目前看是不行的,但我们还是可以通过迂回的方式使用`本地静态库`.
就是在要使用本地静态库的项目的`go.mod`中使用如下设置将依赖的本地库指定给调用方即可

```text
require mymath v0.0.0

replace mymath => ../mymath

```

这种方式本质上并没有使用静态库,而是使用了其源码,golang的最大优势就是编译快,利用这一优势实际上静态库并不是很有必要.这也是为啥golang模块的分发方式是源码分发.

## 动态库 (golang 1.8+)

在go 1.8版本加入了对[动态库](https://baike.baidu.com/item/%E5%8A%A8%E6%80%81%E9%93%BE%E6%8E%A5%E5%BA%93/100352)的支持,使用标准库[plugin](https://golang.org/pkg/plugin/)来实现.目前它还只支持linux和mac.是的,动态库在go语境下叫做插件,这主要是为了和c语言的动态库做区分,事实上go对C语言相当亲和,go代码可以直接编译为c语言的动态库.这个后面再说.

### 编译动态库

我们还是以这个牛顿法求开根号的项目为例子.这回我们不能直接用使用上面的代码,必须将package改为main

```go
package main

func Sqrt(x float64) float64 {
	z:=1.0
	for i:=0;i<10;i++ {
		z -= (z*z-x)/(2*z)
	}
	return z
}
```
这个文件我们放在`$GOPATH/src/calculsqrt_plugin`文件夹下,使用如下命令编译

```bash
go build -buildmode=plugin -o sqrt.so
```

这样就可以编译一个插件了,执行上面的命令后我们可以再这个文件夹下获得动态库文件`sqrt.so`

### 使用动态库

在相同文件夹下我们写一个入口模块`main.go`

```go
package main

import (
	"fmt"
	"plugin"
)

func main() {
    // 加载动态库
	module, err := plugin.Open("./sqrt.so")
	if err != nil {
		fmt.Println("plugin load error")
	} else {
        //在动态库中查找方法
		Sqrt, err := module.Lookup("Sqrt")
		if err != nil {
			fmt.Println("no Sqrt in plugin")
		} else {
            //为函数赋予类型后再调用
			fmt.Println(Sqrt.(func(float64) float64)(2))
		}
	}
}
```

可以看到动态库的使用和python中使用ctype调用c语言的动态库很像,需要先加载,再找到对应的函数,然后还要为这个函数赋予一个类型,最后才能调用.

之后就像之前的程序编译一样直接使用`go install`安装就好,可以看到程序正常编译和安装了,执行`calculsqrt_plugin`也可以获得正确答案.

相对而言我还是更加推荐使用静态库的,go的插件写法复杂不说,编译后项目的大小也比不用的大,除非希望借助插件实现动态热更新,否则完全找不到使用的理由.而热更新其实完全可以借助docker swarm这样的集群工具实现

## 交叉编译

go的另一大卖点是交叉编译,也就说我在mac上可以直接编译windows可以执行的程序.目前这一功能在各种语言中都属于相当先进的特性.

而使用的方式也是相当简单只要在编译时指定特定环境变量即可:

+ GOOS：目标操作系统
+ GOARCH：目标操作系统的架构

下面是目前支持的交叉编译组合:

OS|ARCH|OS version
---|---|---
`linux`|`386 / amd64 / arm`|`>= Linux 2.6`
`darwin`|`386 / amd64`|`OS X (Snow Leopard + Lion)`
`freebsd`|`386 / amd64`|`>= FreeBSD 7`
`windows`|`386 / amd64`|`>= Windows 2000`

无论是`go build`还是`go install`都可以通过指定这两个环境变量编译为对应平台的程序.它会在对应的文件夹下生成一个平台名+cpu命令集名的文件夹,用于存放跨平台的库或者程序,这其实在编译静态库的部分我们已经见识过了.

我们以上面的`calculsqrt`为例,在项目根目录执行命令`GOOS=windows GOARCH=amd64 go build`就可以在项目根目录获得windows上可执行的exe文件了.

## 减小可执行文件体积

go的编译速度很快,算是它的一大卖点,但比较让人诟病的就是它编译出来的可执行文件比较大,虽说硬盘是计算机上最不值钱的东西,但毕竟还是有一些场景我们不得不考虑可执行文件的大小比如在arm单片机上做些服务.我们有两种途径减小文件大小

+ 编译时指定一些参数

    我们可以指定`-ldflags "-s -w"`来减小编译出来的可执行文件大小,代价就是失去一些功能和信息.`-s`的作用是去掉符号信息.`-w`的作用是去掉调试信息.因此这种方式减小文件大小比较适合在已经经过充分测试的发行版上使用.这种方式经过我测试压缩我们的`calculsqrt`项目,从原始的2.25m压缩到了1.61m

+ 使用[upx](https://github.com/upx/upx)

    丑话说前面,upx并不保证可以在各个平台上正常执行,因此这个方案并不稳.upx是一个压缩工具,可以直接压缩可执行文件并且多数情况下不会影响其执行(但执行开始时会有一个解压过程).我们可以去项目的[release栏目下](https://github.com/upx/upx/releases)下按自己的开发平台下载这个工具,并将其加入到PATH环境变量下.

    ```bash
    upx -9 -o $output $target
    ```
    其中9是压缩等级,压缩等级为1到9,9是压缩率最大的等级,压缩原始可执行文件后的可执行文件大小为1.12m;压缩经过`-ldflags`缩减过的可执行文件后其大小为576k.无论如何upx都是一个值得一试的工具,它确实可以解决问题.

### Linux/mac下借助bash脚本实现选择平台编译

为了不用每次都敲一遍相同的代码我们可以使用[bash](http://c.biancheng.net/shell/)来简化这个操作.

+ make.sh

```bash
ASSETS="bin"
GOARCHS=("386" "amd64")
GOOSS=("linux" "darwin" "windows")
export GO111MODULE="on"
# Set the GOPROXY environment variable
export GOPROXY="https://goproxy.io"

case $(uname) in
Darwin)
    case $(uname -m) in
    x86_64)
        cmd="mac"
        ;;
    *)
        cmd="mac32"
        ;;
    esac
    ;;
*)
    case $(uname -m) in
    x86_64)
        cmd="linux64"
        ;;
    *)
        cmd="linux32"
        ;;
    esac
    ;;
esac

cmd="mac"
name="calculsqrt"
if test $# -eq 0; then
    cmd="mac"
elif test $# -eq 1; then
    cmd=$1
elif test $# -eq 2; then
    cmd=$1
    name=$2
else
    echo "args too much"
    exit 0
fi

if ! test -d $ASSETS; then
    mkdir $ASSETS
fi

case $cmd in
all)
    for goarch in ${GOARCHS[@]}; do
        for goos in ${GOOSS[@]}; do
            export GOARCH=$goarch
            export GOOS=$goos
            target="$ASSETS/$GOOS-$GOARCH"
            echo "---------$target----------------"
            if ! test -d $target; then
                mkdir $target
            fi
            case $goos in
            windows)
                go build -o $target/$name.exe
                ;;
            *)
                go build -o $target/$name
                ;;
            esac
        done
    done
    ;;
win32)
    export GOARCH="386"
    export GOOS="windows"
    target="$ASSETS/$GOOS-$GOARCH"
    if ! test -d $target; then
        mkdir $target
    fi
    go build -o $target/$name.exe
    ;;
win64)
    export GOARCH="amd64"
    export GOOS="windows"
    target="$ASSETS/$GOOS-$GOARCH"
    if ! test -d $target; then
        mkdir $target
    fi
    go build -o $target/$name.exe
    ;;
mac)
    export GOARCH="amd64"
    export GOOS="darwin"
    target="$ASSETS/$GOOS-$GOARCH"
    if ! test -d $target; then
        mkdir $target
    fi
    go build -o $target/$name
    ;;
mac32)
    export GOARCH="386"
    export GOOS="darwin"
    target="$ASSETS/$GOOS-$GOARCH"
    if ! test -d $target; then
        mkdir $target
    fi
    go build -o $target/$name
    ;;
linux32)
    export GOARCH="386"
    export GOOS="linux"
    target="$ASSETS/$GOOS-$GOARCH"
    if ! test -d $target; then
        mkdir $target
    fi
    go build -o $target/$name
    ;;
linux64)
    export GOARCH="amd64"
    export GOOS="linux"
    target="$ASSETS/$GOOS-$GOARCH"
    if ! test -d $target; then
        mkdir $target
    fi
    go build -o $target/$name
    ;;
linuxarm)
    export GOARCH="arm"
    export GOOS="linux"
    target="$ASSETS/$GOOS-$GOARCH"
    if ! test -d $target; then
        mkdir $target
    fi
    go build -o $target/$name
    ;;
*)
    echo "unknown cmd $cmd"
    ;;
esac
```

这个脚本允许带两个参数==平台和编译后的名字

### windows下借助powershell选择编译的平台

+ make.ps1

```ps1
$ASSETS = "bin"
$GOARCHS = "386", "amd64"
$GOOSS = "linux", "darwin", "windows"
$env:GO111MODULE="on"
# Set the GOPROXY environment variable
$env:GOPROXY="https://goproxy.io"


$cmd = "win64"
$name = "calculsqrt"
if ($args.Count -eq 0){
    $cmd = "win64"
}elseif ($args.Count -eq 1){
    $cmd = $args[0]
}elseif ($args.Count -eq 2){
    $cmd = $args[0]
    $name = $args[1]
}else{
    echo "args too much"
    exit
}
 
if (!(Test-Path $ASSETS)) {
    mkdir $ASSETS
} 

if ($cmd -eq "all"){
    foreach ($env:GOARCH in $GOARCHS) {
        foreach ($env:GOOS in $GOOSS){
            $target = "$ASSETS/$env:GOOS-$env:GOARCH"
            if (!(Test-Path $target)){
                mkdir $target
            }
            if ($env:GOOS -eq "windows"){
                go build -o $target/$name.exe
            }else {
                go build -o $target/$name
            }
            
        }
    }
}elseif ($cmd -eq "win32") {
    $env:GOARCH="386"
    $env:GOOS="windows"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -o $target/$name.exe
}elseif ($cmd -eq "win64") {
    $env:GOARCH="amd64"
    $env:GOOS="windows"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -o $target/$name.exe
}elseif ($cmd -eq "mac") {
    $env:GOARCH="amd64"
    $env:GOOS="darwin"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -o $target/$name
    
}elseif ($cmd -eq "mac32") {
    $env:GOARCH="386"
    $env:GOOS="darwin"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -o $target/$name
}elseif ($cmd -eq "linux32") {
    $env:GOARCH="386"
    $env:GOOS="linux"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -o $target/$name
}elseif ($cmd -eq "linux64") {
    $env:GOARCH="amd64"
    $env:GOOS="linux"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -o $target/$name
}elseif ($cmd -eq "linuxarm") {
    $env:GOARCH="arm"
    $env:GOOS="linux"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -o $target/$name
}else{
    echo "unknown cmd $cmd"
}
```

