package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ADDRESS         = "0.0.0.0:5000"
	Serv_Cert_Path  = ""
	Serv_Key_Path   = ""
	Ca_Cert_Path    = ""
	Client_Crl_Path = ""
)

//UselessMiddlewareFactory 没什么用的中间件
func UselessMiddlewareFactory(name string) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		log.Println("Hello ", name)
		ctx.Writer.Header().Set("Author", name)
		ctx.Next()
		log.Println("bye ", name)
	}
}

func runserv(router *gin.Engine) {
	srv := &http.Server{
		Addr:    ADDRESS,
		Handler: router,
	}
	usetls := false
	if Serv_Cert_Path != "" && Serv_Key_Path != "" {
		usetls = true
		//双向认证
		if Ca_Cert_Path != "" {
			capool := x509.NewCertPool()
			caCrt, err := ioutil.ReadFile(Ca_Cert_Path)
			if err != nil {
				log.Println("read pem file error")
				os.Exit(2)
			}
			capool.AppendCertsFromPEM(caCrt)
			tlsconf := &tls.Config{
				RootCAs:    capool,
				ClientAuth: tls.RequireAndVerifyClientCert, // 检验客户端证书
			}
			//指定client名单
			if Client_Crl_Path != "" {
				clipool := x509.NewCertPool()
				cliCrt, err := ioutil.ReadFile(Client_Crl_Path)
				if err != nil {
					log.Println("read pem file error")
					os.Exit(2)
				}
				clipool.AppendCertsFromPEM(cliCrt)
				tlsconf.ClientCAs = clipool
			}
			srv.TLSConfig = tlsconf
		}
	}

	//启动服务器
	go func() {
		// 服务连接
		log.Println("servrt start")
		var err error
		if usetls {
			err = srv.ListenAndServeTLS(Serv_Cert_Path, Serv_Key_Path)

		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			log.Println("listen error")
			os.Exit(2)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func main() {
	router := gin.Default()
	router.Use(UselessMiddlewareFactory("router step1"))
	router.Use(UselessMiddlewareFactory("router step2"))
	group1 := router.Group("/group1")
	group1.Use(UselessMiddlewareFactory("group1 step1"))
	group1.Use(UselessMiddlewareFactory("group1 step2"))
	{
		group1.GET("/ping",
			UselessMiddlewareFactory("ping step1"),
			UselessMiddlewareFactory("ping step2"),
			func(ctx *gin.Context) {
				ctx.JSON(200, gin.H{
					"message": "pong",
				})
			})
		group1.GET("/pong",
			UselessMiddlewareFactory("pong step1"),
			UselessMiddlewareFactory("pong step2"),
			func(ctx *gin.Context) {
				ctx.JSON(200, gin.H{
					"message": "ping",
				})
			})
	}
	runserv(router) // listen and serve on 0.0.0.0:5000
}
