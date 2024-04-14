package main

const (
	BindAddress               = "0.0.0.0:4221"
	HttpOKResponse            = "HTTP/1.1 200 OK" + Separator
	HttpCreatedResponse       = "HTTP/1.1 201 Created" + Separator
	HttpNotFoundResponse      = "HTTP/1.1 404 Not Found" + Separator
	HttpInternalErrorResponse = "HTTP/1.1 500 Internal Server Error" + Separator
	Separator                 = "\r\n"
)
