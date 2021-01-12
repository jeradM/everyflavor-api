package core

// AppConfig holder for application configuration
type AppConfig struct {
	// DbURL is the url of the main OLTP MySQL database.
	// Required
	DbURL string

	// LogfileDir is the directory to write logfiles. If not set
	// logging will only be written to os.Stdout
	LogfileDir string

	// Production is set to false during development.
	// Defaults to false
	Production bool

	// RedisURL is the URL of the Redis server for sessions/rate limiter.
	// Required
	RedisURL string

	// RedisDB is the db number on the redis server.
	// Defaults to 0
	RedisDB int

	// RedisKeys are the authentication and encryption keys for
	// the Redis instance. Will be intialized with a static authentication
	// key and no encryption keys if none are provided here
	RedisKeys []string

	// ShowSQL controls whether or not SQL statements are logged.
	// Defaults to false
	ShowSQL bool

	// ServerAddr is the listen address for the HTTP server.
	// Defaults to localhost:8099
	ServerAddr string

	// CorsAllowedOrigins controls which origins will be passed to
	// the CORS middleware as allowed. Don't pass '*' or cross-origin
	// credentials won't work when using the Fetch API and user sessions
	// will break
	CorsAllowedOrigins []string
}
