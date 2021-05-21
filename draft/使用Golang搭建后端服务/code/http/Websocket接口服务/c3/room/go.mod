module c3/room

require (
	c3/errors v0.0.0
	c3/wsexchange v0.0.0
	c3/logger v0.0.0
)

replace (
	c3/errors v0.0.0 => ../errors
	c3/wsexchange v0.0.0 => ../wsexchange
	c3/logger v0.0.0 => ../logger
)

go 1.12
