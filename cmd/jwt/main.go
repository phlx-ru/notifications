package main

import (
	"flag"
	"fmt"
	"path"

	"notifications/internal/auth"
	"notifications/internal/conf"
	"notifications/internal/pkg/logger"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/joho/godotenv"
)

var (
	// Name is the name of the compiled software.
	Name = `notifications_jwt`
	// flagconf is the config flag.
	flagconf string
	// dotenv is loaded from config path .env file
	dotenv string
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&dotenv, "dotenv", ".env", ".env file, eg: -dotenv .env.local")
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flag.Parse()

	var err error

	envPath := path.Join(flagconf, dotenv)
	err = godotenv.Overload(envPath)
	if err != nil {
		return err
	}

	log.SetLogger(logger.New(``, Name, ``, log.LevelFatal.String()))

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(conf.EnvDecoder),
	)
	defer func() {
		_ = c.Close()
	}()

	if err = c.Load(); err != nil {
		return err
	}

	var bc conf.Bootstrap
	if err = c.Scan(&bc); err != nil {
		return err
	}

	token := auth.MakeJWT(bc.Auth.Jwt.Secret)

	colors := map[string]string{
		`reset`:  "\u001B[0m",
		`red`:    "\u001B[31m",
		`green`:  "\u001B[32m",
		`yellow`: "\u001B[33m",
		`blue`:   "\u001B[34m",
		`purple`: "\u001B[35m",
		`cyan`:   "\u001B[36m",
		`white`:  "\u001B[37m",
	}

	fmt.Println(colors["green"] + `JWT token generated:` + colors["reset"])
	fmt.Println(colors["blue"] + token + colors["reset"])
	return nil
}
