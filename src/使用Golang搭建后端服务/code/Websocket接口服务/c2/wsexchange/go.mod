module c2/wsexchange

require (
	c2/errors v0.0.0
	c2/event v0.0.0
	c2/logger v0.0.0
	github.com/gorilla/websocket v1.4.1-0.20190306004257-0ec3d1bd7fe5
	github.com/rfyiamcool/syncmap v0.0.0-20181227021732-a2eaf358f89c
	github.com/satori/go.uuid v1.2.0
)

replace (
	c2/errors v0.0.0 => ../errors
	c2/event v0.0.0 => ../event
	c2/logger v0.0.0 => ../logger
)

go 1.12
