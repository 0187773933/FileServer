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
}

// type PostTypes  struct {
// 	Text string `yaml:"text"`
// 	HTML string `yaml:"html"`
// 	Bytes []byte `yaml:"bytes"`
// }

type Post struct {
	UUID string `json:"uuid"`
	ULID string `json:"ulid"`
	// SeqID int `json:"seq_id"`
	Date string `json:"date"`
	Type string `json:"type"`
	HTML string `json:"html"`
	Text string `json:"text"`
	MD []string `json:"mark_down"`
}

type FileData struct {
	FileName string
	Data []byte
}

type Page struct {
	UUID string `json:"uuid"`
	// ULID string `json:"ulid"`
	// Date string `json:"date"`
	// Type string `json:"type"`
	// HTML string `json:"html"`
	HTMLB64 string `json:"html_b64"`
	URL string `json:"url"`
	// Text string `json:"text"`
	// MD []string `json:"mark_down"`
}
