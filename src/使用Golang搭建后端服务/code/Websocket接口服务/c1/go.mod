module c1

require (
	c1/clienthanddler v0.0.0
	c1/config v0.0.0
	c1/errors v0.0.0
	c1/event v0.0.0
	c1/logger v0.0.0
	c1/serverhanddler v0.0.0
	github.com/gorilla/websocket v1.4.1-0.20190306004257-0ec3d1bd7fe5

	github.com/sirupsen/logrus v1.4.0
)

replace (
	c1/clienthanddler v0.0.0 => ./clienthanddler
	c1/config v0.0.0 => ./config
	c1/errors v0.0.0 => ./errors
	c1/event v0.0.0 => ./event
	c1/logger v0.0.0 => ./logger
	c1/serverhanddler v0.0.0 => ./serverhanddler
)

go 1.12
