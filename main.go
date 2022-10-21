package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
	"integ/api"
	"integ/config"
	"integ/service"
	"integ/storage"
	"integ/storage/migration"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
)

const (
	componentName = "phone"
)

func main() {
	var (
		help       = pflag.BoolP("help", "h", false, "show help message")
		configFile = pflag.StringP("config", "c", "", "name config file only name without extension")
		debug      = pflag.BoolP("debug", "d", false, "enable debug logging")
	)

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if *help {
		pflag.CommandLine.PrintDefaults()
		os.Exit(0)
	}

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	log := logrus.StandardLogger().WithField("component", componentName)

	log.Info("Starting...")
	conf, err := config.LoadConfig(*configFile)
	if err != nil {
		log.WithError(err).Fatalln("Failed to load configuration")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := evolution(conf.Database.ToDataSourceName()); err != nil {
		log.WithError(err).Fatalln("Failed migrate")
	}

	if err := run(ctx, log, conf); err != nil {
		log.WithError(err).Fatalln("Program terminated unexpectedly.")
	}
}

func run(ctx context.Context, log *logrus.Entry, conf *config.Config) error {
	connConf, err := pgx.ParseConfig(conf.Database.ToDataSourceName())
	if err != nil {
		log.WithError(err).Fatalln("Failed to parse database connection string")
	}
	store, err := storage.New(ctx, connConf, log)
	if err != nil {
		log.WithError(err).Fatalln("Failed to create connection pool to database")
	}
	defer store.CloseFn(ctx)

	svc := service.NewContactService(store)

	httpHandler := api.NewHTTPHandler(svc, log)

	router := gin.Default()

	router.POST("/user/:userID/contact", httpHandler.AddContacts)
	router.GET("/user/:userID/friends", httpHandler.Friends)
	router.GET("/user/:userID/contact/name", httpHandler.Name)

	var g errgroup.Group

	g.Go(func() error {
		log.Infof("Listening Phone service at %s", conf.ListenAddr)
		return router.Run(conf.ListenAddr)
	})

	return g.Wait()
}

func evolution(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	tooling := migration.New(db)

	return tooling.Run()
}
