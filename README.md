
```
go-auth-backend-api
├─ .air.toml
├─ cmd
│  ├─ api
│  │  └─ main.go
│  └─ migrate
│     └─ main.go
├─ dockerfile
├─ go.mod
├─ go.sum
├─ go.work
├─ go.work.sum
├─ index.html
├─ internal
│  ├─ config
│  │  ├─ authconfig
│  │  │  └─ GoogleOAuthConfig.go
│  │  ├─ env
│  │  │  └─ env.go
│  │  ├─ mailConfig
│  │  │  └─ mailConfig.go
│  │  └─ redisConfig
│  │     └─ redisConfig.go
│  ├─ handler
│  │  ├─ authHandler
│  │  │  ├─ authRequests.go
│  │  │  ├─ loginHandler.go
│  │  │  ├─ passwordHandler.go
│  │  │  ├─ registerHandler.go
│  │  │  └─ sessionHandler.go
│  │  └─ oAuthHandler
│  │     ├─ googleOAuthHandler.go
│  │     └─ googleRefreshHandler.go
│  ├─ middleware
│  │  ├─ cors
│  │  │  └─ cors.go
│  │  ├─ logger
│  │  │  └─ access_log.go
│  │  ├─ rateLimiter
│  │  │  └─ rateLimiter.go
│  │  └─ setup.go
│  ├─ model
│  │  └─ AuthModel
│  │     └─ authModel.go
│  ├─ repository
│  │  └─ authRepository.go
│  ├─ router
│  │  ├─ authRoutes.go
│  │  └─ router.go
│  └─ service
│     ├─ authService
│     │  ├─ authTypes.go
│     │  ├─ forgotPassword.go
│     │  ├─ helpers.go
│     │  ├─ login.go
│     │  ├─ password.go
│     │  └─ register.go
│     ├─ oAuthService
│     │  ├─ oAuthlogin.go
│     │  └─ oAuthTypes.go
│     └─ rateLimiterService
│        └─ rateLimiterService.go
├─ migrations
│  ├─ authModelMigration.go
│  └─ migrate.go
├─ pkg
│  ├─ cache
│  ├─ database
│  │  └─ connection.go
│  ├─ logger
│  │  └─ logger.go
│  ├─ mailer
│  │  └─ mailer.go
│  ├─ redis
│  │  └─ redis.go
│  ├─ tokenJWT
│  │  └─ jwt.go
│  └─ utils
│     ├─ crypto.go
│     ├─ fomatValidationErrors.go
│     └─ passwordHashing.go
└─ README.md

```