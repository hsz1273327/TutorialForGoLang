module testlog

require (
	github.com/sirupsen/logrus v1.1.1
	testlog/logger v0.0.0
)

replace testlog/logger v0.0.0 => ./logger

go 1.12
