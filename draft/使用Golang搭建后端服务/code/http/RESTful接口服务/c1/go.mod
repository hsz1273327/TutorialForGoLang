module c1

require (
	c1/middleware v0.0.0
	github.com/gin-gonic/gin v1.3.0
)

replace c1/middleware v0.0.0 => ./middleware

go 1.12
