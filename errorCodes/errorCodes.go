package errorCodes

type ErrorCodeMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

/*
1001 - 1049 ==> Server Errors
1050 - 1249 ==> Database Errors
>=1250 ==> App Errors
*/
var BAD_REQUEST = ErrorCodeMessage{
	Code:    1001,
	Message: "Bad Request",
}
var NOT_FOUND = ErrorCodeMessage{
	Code:    1002,
	Message: "Not Found",
}
var INTERNAL_SERVER_ERROR = ErrorCodeMessage{
	Code:    1001,
	Message: "Internal Server Error",
}
