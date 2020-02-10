package errors

import (
	"errors"
)

var WsConnectError error = errors.New("Connect Error")
var GetMessageError error = errors.New("Get Message Error")
var WriteMessageError error = errors.New("Write Message Error")
