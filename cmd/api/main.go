package main

import (
	"expvar"
	"runtime"
	"time"

	"github.com/kuluruvineeth/social-go/internal/auth"
	"github.com/kuluruvineeth/social-go/internal/db"
	"github.com/kuluruvineeth/social-go/internal/env"
	"github.com/kuluruvineeth/social-go/internal/mailer"
	"github.com/kuluruvineeth/social-go/internal/ratelimiter"
	"github.com/kuluruvineeth/social-go/internal/store"
	"github.com/kuluruvineeth/social-go/internal/store/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Social Go API
//	@description	This is a sample server Social Go server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				The token for the user
func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social_go?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdletime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:         env.GetString("ENV", "development"),
		apiURL:      env.GetString("ADDR", "0.0.0.0:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		mail: mailConfig{
			fromEmail: env.GetString("FROM_EMAIL", ""),
			exp:       time.Hour * 24 * 3, //3 days
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAILTRAP_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicAuthConfig{
				user:     env.GetString("BASIC_AUTH_USER", "admin"),
				password: env.GetString("BASIC_AUTH_PASSWORD", "password"),
			},
			token: tokenConfig{
				secret: env.GetString("TOKEN_SECRET", ""),
				exp:    time.Hour * 24 * 3, //3 days
				iss:    env.GetString("TOKEN_ISS", "social-go"),
				aud:    env.GetString("TOKEN_AUD", "social-go"),
			},
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			db:      env.GetInt("REDIS_DB", 0),
			pw:      env.GetString("REDIS_PW", ""),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATE_LIMITER_REQUESTS_PER_TIME_FRAME", 20),
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", false),
		},
	}

	//Logger
	var logger *zap.SugaredLogger
	if cfg.env == "development" {
		// Use development logger for better debugging
		zapLogger := zap.Must(zap.NewDevelopment())
		logger = zapLogger.Sugar()
	} else {
		// Use production logger for production
		logger = zap.Must(zap.NewProduction()).Sugar()
	}
	defer logger.Sync()

	// Log the configuration for debugging
	logger.Infow("Starting server",
		"env", cfg.env,
		"addr", cfg.addr,
		"basic_auth_user", cfg.auth.basic.user,
		"redis_enabled", cfg.redisCfg.enabled,
	)

	//Database
	db, err := db.NewDB(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdletime)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("Connected to the database")

	//cache
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Info("Connected to the redis")

		defer rdb.Close()
	}

	//rate limiter
	rateLimiter := ratelimiter.NewFixedWindowRateLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	store := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)
	mailer := mailer.NewSendGrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	// Uncomment this to use Mailtrap
	// mailtrap, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.aud, cfg.auth.token.iss)

	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
		cache:         cacheStorage,
		rateLimiter:   rateLimiter,
	}

	// Metrics collected
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()

	logger.Fatal(app.run(mux))

}
