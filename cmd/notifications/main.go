package main

import (
	"bytes"
	"flag"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"regexp"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"notifications/internal/conf"
)

const (
	dotenv = `.env`
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	ppp := os.Getenv("POSTGRES_PASS")
	println(ppp)
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	//absconf, err := filepath.Abs(flagconf)
	//if err != nil {
	//	panic(err)
	//}
	err := godotenv.Overload(path.Join(flagconf, dotenv))
	if err != nil {
		panic(err)
	}
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			configData := ReplaceEnv(kv.Value)
			return yaml.Unmarshal(configData, v)
		}),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func ReplaceEnv(configData []byte) []byte {
	r, _ := regexp.Compile(`\${((.+?):?\w+?)}`)

	for _, match := range r.FindAllSubmatch(configData, -1) {
		key := os.Getenv(string(match[2]))
		if key != "" {
			configData = bytes.Replace(configData, match[0], []byte(key), 1)
		}
	}

	return configData
}
