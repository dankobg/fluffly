package config

import (
	"os"
	"slices"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	prefix     = "FLUFFLY"
	oldDelim   = "_"
	newDelim   = "."
	sliceDelim = ","
)

var knownKeys = []string{"server", "cors", "postgres", "redis", "email", "logger"}

// Fluffly contains common fluffly app settings
type Fluffly struct {
	ENV             string `koanf:"env"`
	Host            string `koanf:"host"`
	Port            int    `koanf:"port"`
	BaseURL         string `koanf:"base_url"`
	WebsiteURL      string `koanf:"website_url"`
	FileStorage     string `koanf:"file_storage"`
	UploadDir       string `koanf:"upload_dir"`
	OpenapiSpecURL  string `koanf:"openapi_spec_url"`
	KratosPublicURL string `koanf:"kratos_public_url"`
	KratosAdminURL  string `koanf:"kratos_admin_url"`
	KratosAPIKey    string `koanf:"kratos_api_key"`
	KetoAPIKey      string `koanf:"keto_api_key"`
}

// ServerConfig contains the http server settings
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

// CorsConfig contains the CORS settings
type CorsConfig struct {
	AllowOrigins     []string `koanf:"allow_origins"`
	AllowMethods     []string `koanf:"allow_methods"`
	AllowHeaders     []string `koanf:"allow_headers"`
	ExposeHeaders    []string `koanf:"expose_headers"`
	AllowCredentials bool     `koanf:"allow_credentials"`
	MaxAge           int      `koanf:"max_age"`
	Debug            bool     `koanf:"debug"`
}

// DatabaseConfig contains DB settings
type DatabaseConfig struct {
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

// MinioConfig contains minio settings
type MinioConfig struct {
	RootUser      string `koanf:"root_user"`
	RootPassword  string `koanf:"root_password"`
	Host          string `koanf:"host"`
	ApiPort       int    `koanf:"api_port"`
	ConsolePort   int    `koanf:"console_port"`
	DefaultBucket string `koanf:"default_bucket"`
	Token         string `koanf:"token"`
	UseSSL        bool   `koanf:"use_ssl"`
}

// RedisConfig contains Redis db settings
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

// LoggerConfig contains the logger settings
type LoggerConfig struct {
	Level  string `koanf:"level"`
	Pretty bool   `koanf:"pretty"`
}

// Config represents the app config
type Config struct {
	Fluffly  `koanf:",squash"`
	Server   ServerConfig   `koanf:"server"`
	Cors     CorsConfig     `koanf:"cors"`
	Database DatabaseConfig `koanf:"postgres"`
	Redis    RedisConfig    `koanf:"redis"`
	Email    EmailConfig    `koanf:"email"`
	Logger   LoggerConfig   `koanf:"logger"`
	Minio    MinioConfig    `koanf:"minio"`
}

// loadEnv loads env files by convention: https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use
func loadEnv() error {
	env := os.Getenv("FLUFFLY_ENV")
	if env == "" {
		env = "development"
	}
	godotenv.Load(".env." + env + ".local")
	if env != "test" {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load(".env")

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
	return k.Load(file.Provider("config/defaults.yaml"), yaml.Parser())
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
