package db

import "github.com/gofiber/fiber/v2"

type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// AllowOrigin defines a list of origins that may access the resource.
	//
	// Optional. Default value "*"
	AllowOrigins string

	// AllowMethods defines a list methods allowed when accessing the resource.
	// This is used in response to a preflight request.
	//
	// Optional. Default value "GET,POST,HEAD,PUT,DELETE,PATCH"
	AllowMethods string

	// AllowHeaders defines a list of request headers that can be used when
	// making the actual request. This is in response to a preflight request.
	//
	// Optional. Default value "".
	AllowHeaders string

	// AllowCredentials indicates whether or not the response to the request
	// can be exposed when the credentials flag is true. When used as part of
	// a response to a preflight request, this indicates whether or not the
	// actual request can be made using credentials.
	//
	// Optional. Default value false.
	AllowCredentials bool

	// ExposeHeaders defines a whitelist headers that clients are allowed to
	// access.
	//
	// Optional. Default value "".
	ExposeHeaders string

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached.
	//
	// Optional. Default value 0.
	MaxAge int
}

var ConfigDefault = Config{
	Next:             nil,
	AllowOrigins:     "http://localhost:3000, https://zoom.icetoast.cloud",
	AllowMethods:     "GET,POST,PUT,DELETE",
	AllowHeaders:     "Origin, X-Requested-With, Content-Type, Accept, Authorization",
	AllowCredentials: true,
	ExposeHeaders:    "",
	MaxAge:           0,
}
