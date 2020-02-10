module c3/serverhanddler

require (
	c3/errors v0.0.0
	c3/event v0.0.0
	c3/logger v0.0.0
	c3/room v0.0.0
	c3/wsexchange v0.0.0
	github.com/gorilla/websocket v1.4.1-0.20190306004257-0ec3d1bd7fe5
)

replace (
	c3/errors v0.0.0 => ../errors
	c3/event v0.0.0 => ../event
	c3/logger v0.0.0 => ../logger
	c3/room v0.0.0 => ../room
	c3/wsexchange v0.0.0 => ../wsexchange
)

go 1.12
