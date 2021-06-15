module c2/clienthanddler

require (
	c2/errors v0.0.0
	c2/event v0.0.0
	c2/logger v0.0.0
	github.com/gorilla/websocket v1.4.1-0.20190306004257-0ec3d1bd7fe5
)

replace (
	c2/errors v0.0.0 => ../errors
	c2/event v0.0.0 => ../event
	c2/logger v0.0.0 => ../logger
)

go 1.12
