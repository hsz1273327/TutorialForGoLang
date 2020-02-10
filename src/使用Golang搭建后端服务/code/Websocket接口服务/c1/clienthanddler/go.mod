module c1/clienthanddler

require (
	c1/errors v0.0.0
	c1/event v0.0.0
	c1/logger v0.0.0
	github.com/gorilla/websocket v1.4.1-0.20190306004257-0ec3d1bd7fe5
)

replace (
	c1/errors v0.0.0 => ../errors
	c1/event v0.0.0 => ../event
	c1/logger v0.0.0 => ../logger
)

go 1.12
