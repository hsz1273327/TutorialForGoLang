package errors

import (
	"errors"
)

var WsConnectError error = errors.New("Connect Error")
var GetMessageError error = errors.New("Get Message Error")
var WriteMessageError error = errors.New("Write Message Error")
var WsExistError error = errors.New("Ws Exist Error")
var RoomNotExistError error = errors.New("Room Not Exist Error")
var RoomExistError error = errors.New("Room Exist Error")
