module testconfigload

require (
	github.com/sirupsen/logrus v1.2.0
	github.com/small-tk/pathlib v0.0.0-20190601032836-742166d9b695 // indirect
	github.com/tutorialforgolang/calculsqrt v0.0.3-0.20190718085353-a1338c608320
	testconfigload/config v0.0.0
	testconfigload/logger v0.0.0
)

replace (
	testconfigload/config v0.0.0 => ./config
	testconfigload/logger v0.0.0 => ./logger
)

go 1.12
