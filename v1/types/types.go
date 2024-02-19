package types

type ConfigFile struct {
	ServerBaseUrl string `yaml:"server_base_url"`
	ServerLiveUrl string `yaml:"server_live_url"`
	ServerPort string `yaml:"server_port"`
	ServerAPIKey string `yaml:"server_api_key"`
	ServerCookieName string `yaml:"server_cookie_name"`
	ServerCookieSecret string `yaml:"server_cookie_secret"`
	ServerCookieAdminSecretMessage string `yaml:"server_cookie_admin_secret_message"`
	ServerCookieSecretMessage string `yaml:"server_cookie_secret_message"`
	AdminUsername string `yaml:"admin_username"`
	AdminPassword string `yaml:"admin_password"`
	TimeZone string `yaml:"time_zone"`
	BoltDBPath string `yaml:"bolt_db_path"`
	BoltDBEncryptionKey string `yaml:"bolt_db_encryption_key"`
	RedisHost string `yaml:"redis_host"`
	RedisPort string `yaml:"redis_port"`
	RedisDBNumber int `yaml:"redis_db_number"`
	RedisPassword string `yaml:"redis_password"`
	ServeDirectory string `yaml:"serve_directory"`
	ServeBrowsable bool `yaml:"serve_browsable"`
	ServeIndexFile string `yaml:"serve_index_file"`
	PublicLimiterMax int `yaml:"public_limiter_max"`
	PublicLimiterSeconds int `yaml:"public_limiter_seconds"`
	PublicLimiterMaxLimitCount int64 `yaml:"public_limiter_max_limit_count"`
}