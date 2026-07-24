package config

import (
	_ "embed"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

const (
	prefix     = ""
	oldDelim   = "_"
	newDelim   = "."
	sliceDelim = ","
)

var knownKeys = []string{
	"app",
	"kratos",
	"keto",
	"server",
	"cors",
	"postgres",
	"redis",
	"email",
	"logger",
	"rustfs",
	"petfinder",
	"rustfs",
}

// AppConfig contains fluffly app settings
type AppConfig struct {
	ENV            string `koanf:"env"`
	Host           string `koanf:"host"`
	Port           int    `koanf:"port"`
	BaseURL        string `koanf:"base_url"`
	WebsiteURL     string `koanf:"website_url"`
	FileStorage    string `koanf:"file_storage"`
	UploadDir      string `koanf:"upload_dir"`
	OpenapiSpecURL string `koanf:"openapi_spec_url"`
}

// ServerConfig contains http server settings
type ServerConfig struct {
	ReadHeaderTimeout time.Duration `koanf:"read_header_timeout"`
	ReadTimeout       time.Duration `koanf:"read_timeout"`
	WriteTimeout      time.Duration `koanf:"write_timeout"`
	IdleTimeout       time.Duration `koanf:"idle_timeout"`
	GracefulTimeout   time.Duration `koanf:"graceful_timeout"`
	UseTLS            bool          `koanf:"use_tls"`
	CERT_FILE         string        `koanf:"cert_file"`
	KEY_FILE          string        `koanf:"key_file"`
}

// KratosConfig contains kratos settings
type KratosConfig struct {
	ServePublicBaseURL string `koanf:"serve_public_base_url"`
	ServeAdminBaseURL  string `koanf:"serve_admin_base_url"`
	ApiKey             string `koanf:"api_key"`
}

// KetoConfig contains keto settings
type KetoConfig struct {
	ServeReadURL  string `koanf:"serve_read_url"`
	ServeWriteURL string `koanf:"serve_write_url"`
	ApiKey        string `koanf:"api_key"`
}

// CorsConfig contains CORS settings
type CorsConfig struct {
	AllowOrigins     []string `koanf:"allow_origins"`
	AllowMethods     []string `koanf:"allow_methods"`
	AllowHeaders     []string `koanf:"allow_headers"`
	ExposeHeaders    []string `koanf:"expose_headers"`
	AllowCredentials bool     `koanf:"allow_credentials"`
	MaxAge           int      `koanf:"max_age"`
	Debug            bool     `koanf:"debug"`
}

// PostgresConfig contains postgres settings
type PostgresConfig struct {
	Host         string        `koanf:"host"`
	Port         int           `koanf:"port"`
	DB           string        `koanf:"db"`
	User         string        `koanf:"user"`
	Password     string        `koanf:"password"`
	Scheme       string        `koanf:"scheme"`
	SSLMode      string        `koanf:"ssl_mode"`
	RetriesNum   int           `koanf:"retries_num"`
	RetriesDelay time.Duration `koanf:"retries_delay"`
}

// RustfsConfig contains rustfs settings
type RustfsConfig struct {
	Host                      string   `koanf:"host"`
	Port                      int      `koanf:"port"`
	ConsolePort               int      `koanf:"console_port"`
	ConsoleEnable             bool     `koanf:"console_enable"`
	CorsAllowedOrigins        []string `koanf:"cors_allowed_origins"`
	ConsoleCorsAllowedOrigins []string `koanf:"console_cors_allowed_origins"`
	AccessKey                 string   `koanf:"access_key"`
	SecretKey                 string   `koanf:"secret_key"`
	DefaultBucket             string   `koanf:"default_bucket"`
	ServerDomains             []string `koanf:"server_domains"`
	UseSSL                    bool     `koanf:"use_ssl"`
	Token                     string   `koanf:"token"`
}

// RedisConfig contains redis settings
type RedisConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	DB       int    `koanf:"db"`
	Password string `koanf:"password"`
}

// EmailConfig contains email settings
type EmailConfig struct {
	Enabled         bool   `koanf:"enabled"`
	TLS             bool   `koanf:"tls"`
	FromName        string `koanf:"from_name"`
	FromAddress     string `koanf:"from_address"`
	SMTPHost        string `koanf:"smtp_host"`
	SMTPPort        int    `koanf:"smtp_port"`
	SMTPUsername    string `koanf:"smtp_username"`
	SMTPPassword    string `koanf:"smtp_password"`
	DevSMTPHost     string `koanf:"dev_smtp_host"`
	DevSMTPPort     int    `koanf:"dev_smtp_port"`
	DevSMTPUsername string `koanf:"dev_smtp_username"`
	DevSMTPPassword string `koanf:"dev_smtp_password"`
}

// LoggerConfig contains logger settings
type LoggerConfig struct {
	Level  string `koanf:"level"`
	Pretty bool   `koanf:"pretty"`
}

// PetfinderConfig contains petfinder settings
type PetfinderConfig struct {
	ApiKey      string `koanf:"api_key"`
	ApiSecret   string `koanf:"api_secret"`
	AccessToken string `koanf:"access_token"`
}

// GeocodingConfig contains geocoding settings
type GeocodingConfig struct {
	Enabled bool `koanf:"enabled"`
}

// Config represents the app config
type Config struct {
	App       AppConfig       `koanf:"app"`
	Server    ServerConfig    `koanf:"server"`
	Kratos    KratosConfig    `koanf:"kratos"`
	Keto      KetoConfig      `koanf:"keto"`
	Cors      CorsConfig      `koanf:"cors"`
	Postgres  PostgresConfig  `koanf:"postgres"`
	Redis     RedisConfig     `koanf:"redis"`
	Email     EmailConfig     `koanf:"email"`
	Logger    LoggerConfig    `koanf:"logger"`
	Rustfs    RustfsConfig    `koanf:"rustfs"`
	Petfinder PetfinderConfig `koanf:"petfinder"`
	Geocoding GeocodingConfig `koanf:"geocoding"`
}

// loadEnv loads env files by convention: https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use
func loadEnv() error {
	env := os.Getenv("FLUFFLY_ENV")
	if env == "" {
		env = "development"
	}

	_ = godotenv.Load(".env." + env + ".local")

	if env != "test" {
		_ = godotenv.Load(".env.local")
	}

	_ = godotenv.Load(".env." + env)
	_ = godotenv.Load(".env")

	return nil
}

func loadFromEnv(k *koanf.Koanf) error {
	p := env.ProviderWithValue(prefix, newDelim, func(s string, v string) (string, any) {
		str := strings.ToLower(strings.TrimPrefix(s, prefix+oldDelim))
		key := strings.Replace(str, oldDelim, newDelim, 1)

		if !slices.Contains(knownKeys, strings.Split(key, newDelim)[0]) {
			key = strings.Replace(key, newDelim, oldDelim, 1)
		}

		if strings.Contains(v, sliceDelim) {
			return key, strings.Split(v, sliceDelim)
		}

		return key, v
	})

	return k.Load(p, nil)
}

func loadDefaults(k *koanf.Koanf) error {
	defaultConfig, err := getDefaultConfig()
	if err != nil {
		return err
	}

	return k.Load(structs.Provider(defaultConfig, "koanf"), nil)
}

func getConfig(k *koanf.Koanf) (*Config, error) {
	var cfg Config
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// New creates the new config
func New() (*Config, *koanf.Koanf, error) {
	k := koanf.New(".")

	if err := loadEnv(); err != nil {
		return nil, nil, err
	}

	if err := loadDefaults(k); err != nil {
		return nil, nil, err
	}

	if err := loadFromEnv(k); err != nil {
		return nil, nil, err
	}

	cfg, err := getConfig(k)
	if err != nil {
		return nil, nil, err
	}

	return cfg, k, nil
}
