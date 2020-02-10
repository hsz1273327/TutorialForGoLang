$ASSETS = "bin"
$GOARCHS = "386", "amd64"
$GOOSS = "linux", "darwin", "windows"
$env:GO111MODULE="on"
# Set the GOPROXY environment variable
$env:GOPROXY="https://goproxy.io"


$cmd = "win64"
$name = "server"

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
                go build -o $target/client.exe client/client.go
            }else {
                go build -o $target/$name
                go build -o  $target/client client/client.go
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
    go build -o $target/client.exe client/client.go
}elseif ($cmd -eq "win64") {
    $env:GOARCH="amd64"
    $env:GOOS="windows"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -o $target/$name.exe
    go build -o $target/client.exe client/client.go
}elseif ($cmd -eq "mac") {
    $env:GOARCH="amd64"
    $env:GOOS="darwin"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -o $target/$name
    go build -o $target/client client/client.go
}elseif ($cmd -eq "mac32") {
    $env:GOARCH="386"
    $env:GOOS="darwin"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -o $target/$name
    go build -o $target/client client/client.go
}elseif ($cmd -eq "linux32") {
    $env:GOARCH="386"
    $env:GOOS="linux"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -o $target/$name
    go build -o $target/client client/client.go
}elseif ($cmd -eq "linux64") {
    $env:GOARCH="amd64"
    $env:GOOS="linux"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -o $target/$name
    go build -o $target/client client/client.go
}elseif ($cmd -eq "linuxarm") {
    $env:GOARCH="arm"
    $env:GOOS="linux"
    $target = "$ASSETS/$env:GOOS-$env:GOARCH"
    if (!(Test-Path $target)){
        mkdir $target
    }
    go build -o $target/$name
    go build -o $target/client client/client.go
}else{
    echo "unknown cmd $cmd"
}