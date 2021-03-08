module c2/room

require (
	c2/errors v0.0.0
	c2/wsexchange v0.0.0
	c2/logger v0.0.0
	github.com/rfyiamcool/syncmap v0.0.0-20181227021732-a2eaf358f89c
)

replace (
	c2/errors v0.0.0 => ../errors
	c2/wsexchange v0.0.0 => ../wsexchange
	c2/logger v0.0.0 => ../logger
)

go 1.12
