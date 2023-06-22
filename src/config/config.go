package config

type MySQLEmpty struct {
	Address  string `env:"EMPTY_MYSQL_ADDRESS" envDefault:"localhost:3306"`
	Username string `env:"EMPTY_MYSQL_USERNAME" envDefault:"root"`
	Password string `env:"EMPTY_MYSQL_PASSWORD" envDefault:"root"`
	Database string `env:"EMPTY_MYSQL_DATABASE" envDefault:"boilerplate"`
}

type MongoEmpty struct {
	Host       string `env:"AUTH_MONGO_HOST" envDefault:"localhost"`
	Port       string `env:"AUTH_MONGO_PORT" envDefault:"27017"`
	Username   string `env:"AUTH_MONGO_USERNAME" envDefault:""`
	Password   string `env:"AUTH_MONGO_PASSWORD" envDefault:""`
	Database   string `env:"AUTH_MONGO_DATABASE" envDefault:"empty"`
	Collection string `env:"AUTH_MONGO_COLLECTION" envDefault:"empties"`
}

type I18n struct {
	Fallback string   `env:"I18N_FALLBACK_LANGUAGE" envDefault:"en"`
	Dir      string   `env:"I18N_DIR" envDefault:"./src/locales"`
	Locales  []string `env:"I18N_LOCALES" envDefault:"en,tr"`
}

type Server struct {
	Host  string `env:"SERVER_HOST" envDefault:"localhost"`
	Port  int    `env:"SERVER_PORT" envDefault:"3000"`
	Group string `env:"SERVER_GROUP" envDefault:"account"`
}

type Redis struct {
	Host string `env:"REDIS_HOST"`
	Port string `env:"REDIS_PORT"`
	Pw   string `env:"REDIS_PASSWORD"`
	Db   int    `env:"REDIS_DB"`
}

type HttpHeaders struct {
	AllowedOrigins   string `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`
	AllowedMethods   string `env:"CORS_ALLOWED_METHODS" envDefault:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowedHeaders   string `env:"CORS_ALLOWED_HEADERS" envDefault:"*"`
	AllowCredentials bool   `env:"CORS_ALLOW_CREDENTIALS" envDefault:"true"`
	Domain           string `env:"HTTP_HEADER_DOMAIN" envDefault:"*"`
}

type TokenSrv struct {
	Expiration int    `env:"TOKEN_EXPIRATION" envDefault:"3600"`
	Project    string `env:"TOKEN_PROJECT" envDefault:"empty"`
}

type Session struct {
	Topic string `env:"SESSION_TOPIC"`
}

type Topics struct {
	Empty EmptyTopics
}

type EmptyTopics struct {
	Created string `env:"STREAMING_TOPIC_EXAMPLE_CREATED"`
	Updated string `env:"STREAMING_TOPIC_EXAMPLE_UPDATED"`
}

type Nats struct {
	Url     string   `env:"NATS_URL" envDefault:"nats://localhost:4222"`
	Streams []string `env:"NATS_STREAMS" envDefault:""`
}

type CDN struct {
	Host        string `env:"CDN_HOST"`
	UploadHost  string `env:"CDN_UPLOAD_HOST"`
	StorageZone string `env:"CDN_STORAGE_ZONE"`
	ApiKey      string `env:"CDN_API_KEY"`
}

type App struct {
	Protocol    string `env:"PROTOCOL" envDefault:"http"`
	MySQLEmpty  MySQLEmpty
	MongoEmpty  MongoEmpty
	Server      Server
	HttpHeaders HttpHeaders
	I18n        I18n
	Topics      Topics
	Session     Session
	Nats        Nats
	Redis       Redis
	TokenSrv    TokenSrv
	CDN         CDN
}
