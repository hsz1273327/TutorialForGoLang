module c2

require (
	c2/logger v0.0.0
	c2/sources v0.0.0
	github.com/gin-gonic/gin v1.3.0
	github.com/labstack/gommon v0.2.8
	github.com/mattn/go-colorable v0.1.1 // indirect
	github.com/sirupsen/logrus v1.4.0
	github.com/toorop/gin-logrus v0.0.0-20190324082946-8887861896bb
	github.com/valyala/fasttemplate v1.0.1 // indirect
)

replace (
	c2/logger v0.0.0 => ./logger
	c2/sources v0.0.0 => ./sources
)

go 1.12
