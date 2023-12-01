module c3

require (
	c3/clienthanddler v0.0.0
	c3/config v0.0.0
	c3/errors v0.0.0
	c3/event v0.0.0
	c3/logger v0.0.0
	c3/room v0.0.0
	c3/serverhanddler v0.0.0
	c3/wsexchange v0.0.0
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.1-0.20190306004257-0ec3d1bd7fe5

	github.com/sirupsen/logrus v1.4.0
)

replace (
	c3/clienthanddler v0.0.0 => ./clienthanddler
	c3/config v0.0.0 => ./config
	c3/errors v0.0.0 => ./errors
	c3/event v0.0.0 => ./event
	c3/logger v0.0.0 => ./logger
	c3/room v0.0.0 => ./room
	c3/serverhanddler v0.0.0 => ./serverhanddler
	c3/wsexchange v0.0.0 => ./wsexchange
)

go 1.12
