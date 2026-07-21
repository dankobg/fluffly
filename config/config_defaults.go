package config

import "time"

func getDefaultConfig() (Config, error) {
	defaultApp := AppConfig{
		ENV:             "development",
		Host:            "localhost",
		Port:            6969,
		BaseURL:         "https://fluffly-dev.xyz",
		WebsiteURL:      "https://fluffly-dev.xyz",
		FileStorage:     "rustfs",
		UploadDir:       "./uploads",
		OpenapiSpecURL:  "https://fluffly-dev.xyz/spec",
		KratosPublicURL: "http://localhost:4433",
		KratosAdminURL:  "http://localhost:4434",
		KratosAPIKey:    "v3Ry_s3Cr3t_tExT_kr4t0s",
		KetoReadURL:     "localhost:4466",
		KetoWriteURL:    "localhost:4467",
		KetoAPIKey:      "v3Ry_s3Cr3t_tExT_k3t0",
	}

	defaultServer := ServerConfig{
		ReadHeaderTimeout: time.Second * 5,
		ReadTimeout:       time.Second * 15,
		WriteTimeout:      time.Second * 15,
		IdleTimeout:       time.Second * 120,
		GracefulTimeout:   time.Second * 30,
		UseTLS:            false,
		CERT_FILE:         "",
		KEY_FILE:          "",
	}

	defaultCors := CorsConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://fluffly-dev.xyz",
			"https://fluffly-dev.xyz",
			"https://client.scalar.com",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"HEAD",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-CSRF-Token",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Cache-Control",
			"Content-Language",
			"Content-Type",
			"Content-Range",
			"Expires",
			"Last-Modified",
			"Pragma",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           86400,
		Debug:            false,
	}

	defaultDatabase := DatabaseConfig{
		Host:         "localhost",
		Port:         5432,
		DB:           "test",
		User:         "test",
		Password:     "test",
		Scheme:       "postgres",
		SSLMode:      "disable",
		RetriesNum:   5,
		RetriesDelay: time.Second * 2,
	}

	defaultRedis := RedisConfig{
		Host:     "localhost",
		Port:     6379,
		DB:       0,
		Password: "",
	}

	defaultEmail := EmailConfig{
		Enabled:         false,
		TLS:             false,
		FromName:        "fluffly",
		FromAddress:     "fluffly-dev.xyz",
		SMTPHost:        "",
		SMTPPort:        587,
		SMTPUsername:    "",
		SMTPPassword:    "",
		DevSMTPHost:     "mailpit",
		DevSMTPPort:     1025,
		DevSMTPUsername: "test",
		DevSMTPPassword: "test",
	}

	defaultLogger := LoggerConfig{
		Level:  "debug",
		Pretty: true,
	}

	defaultRustfs := RustfsConfig{
		Host:                      "localhost",
		Address:                   9000,
		ConsoleAddress:            9001,
		ConsoleEnable:             true,
		AccessKey:                 "test",
		SecretKey:                 "test",
		DefaultBucket:             "fluffly",
		ServerDomains:             []string{"localhost"},
		CorsAllowedOrigins:        []string{"*"},
		ConsoleCorsAllowedOrigins: []string{"*"},
		UseSSL:                    false,
		Token:                     "",
	}

	defaultConfig := Config{
		App:      defaultApp,
		Server:   defaultServer,
		Cors:     defaultCors,
		Database: defaultDatabase,
		Redis:    defaultRedis,
		Email:    defaultEmail,
		Logger:   defaultLogger,
		Rustfs:   defaultRustfs,
	}

	return defaultConfig, nil
}
