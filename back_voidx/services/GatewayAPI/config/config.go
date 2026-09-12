package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// serv
	ServerPort string

	// gRPC servses addrs
	AuthServiceAddr         string
	UserServiceAddr         string
	ChatServiceAddr         string
	ConfigServiceAddr       string
	NotificationServiceAddr string

	// Redis
	RedisAddr string
	RedisDB   int

	// JWT
	JWTSecret string
	JWTExpire string

	// Timeouts (in seconds)
	GRPCTimeout  int
	ReadTimeout  int
	WriteTimeout int

	// Env
	Environment string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("[ERROR] - config load error")
	}

	redisDb, err := strconv.Atoi(os.Getenv("RedisDB"))
	if err != nil {
		log.Fatal("[ERROR] - convert str to int")
	}
	grpctmt, err := strconv.Atoi(os.Getenv("GRPCTimeout"))
	if err != nil {
		log.Fatal("[ERROR] - convert str to int")
	}
	readtmt, err := strconv.Atoi(os.Getenv("ReadTimeout"))
	if err != nil {
		log.Fatal("[ERROR] - convert str to int")
	}
	writetmt, err := strconv.Atoi(os.Getenv("WriteTimeout"))
	if err != nil {
		log.Fatal("[ERROR] - convert str to int")
	}

	return Config{
		ServerPort: os.Getenv("PORT"),

		AuthServiceAddr:         os.Getenv("AuthServiceAddr"),
		UserServiceAddr:         os.Getenv("UserServiceAddr"),
		ChatServiceAddr:         os.Getenv("ChatServiceAddr"),
		ConfigServiceAddr:       os.Getenv("ConfigServiceAddr"),
		NotificationServiceAddr: os.Getenv("NotificationServiceAddr"),

		RedisAddr: os.Getenv("RedisAddr"),
		RedisDB:   redisDb,
		JWTSecret: os.Getenv("JWTSecret"),
		JWTExpire: os.Getenv("JWTExpire"),

		GRPCTimeout:  grpctmt,
		ReadTimeout:  readtmt,
		WriteTimeout: writetmt,

		Environment: os.Getenv("Environment"),
	}
}
