module c0

require (
	c0/squarerpc_service v0.0.0
	github.com/golang/protobuf v1.3.1
	google.golang.org/grpc v1.19.0
)

replace c0/squarerpc_service v0.0.0 => ./squarerpc_service

go 1.12
