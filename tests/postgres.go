package tests

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/testcontainers/testcontainers-go"
)

const (
	defaultPGUsername    = "test"
	defaultPGDatabase    = "test"
	defaultPGDockerImage = "postgres:12"

	defaultStartupTimeout = 10 * time.Minute
)

var defaultPGCmdLineOptions = []string{"-c", "fsync=off", "-c", "log_statement=all"}

var postgresURL = ""

type postgresLauncherConfig struct {
	image string

	databaseName string
	databaseUser string

	cmdLineOptions []string
}

func GetPostgresURL() string {
	return postgresURL
}

func PostgresMain(tests func() int) error {
	postgresURL = os.Getenv("TEST_INT_POSTGRES_DSN")
	if postgresURL != "" {
		_ = tests()
		return nil
	}

	launcherConfig := &postgresLauncherConfig{
		image: defaultPGDockerImage,

		databaseName: defaultPGDatabase,
		databaseUser: defaultPGUsername,

		cmdLineOptions: defaultPGCmdLineOptions,
	}

	stop, err := startPostgres(launcherConfig)
	if err != nil {
		return fmt.Errorf("failed to run integration tests: unable to start postgres container: %v", err)
	}
	defer stop()

	_ = tests()
	return nil
}

func startPostgres(launcherConfig *postgresLauncherConfig) (func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: launcherConfig.image,
			Env: map[string]string{
				"POSTGRES_DB":   launcherConfig.databaseName,
				"POSTGRES_USER": launcherConfig.databaseUser,

				"POSTGRES_HOST_AUTH_METHOD": "trust",
			},
			ExposedPorts: []string{"5432/tcp"},
			Cmd:          launcherConfig.cmdLineOptions,
			WaitingFor:   forPostgres(launcherConfig.databaseName, launcherConfig.databaseUser),
		},
		Started: false,
	})
	if err != nil {
		return nil, err
	}

	startCtx, startCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer startCancel()
	if err := postgresC.Start(startCtx); err != nil {
		return nil, err
	}

	stop := func() {
		_ = postgresC.Terminate(context.Background())
	}

	host, err := postgresC.Host(context.TODO())
	if err != nil {
		stop()
		return nil, err
	}

	port, err := postgresC.MappedPort(context.TODO(), "5432")
	if err != nil {
		stop()
		return nil, err
	}

	postgresURL = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port.Int(), launcherConfig.databaseUser, launcherConfig.databaseName)
	return stop, nil
}

var _ wait.Strategy = (*postgresWaitStrategy)(nil)

type postgresWaitStrategy struct {
	DatabaseUser string
	DatabaseName string
}

func forPostgres(dbUser, dbName string) wait.Strategy {
	return &postgresWaitStrategy{
		DatabaseName: dbName,
		DatabaseUser: dbUser,
	}
}

func (ws *postgresWaitStrategy) WaitUntilReady(ctx context.Context, target wait.StrategyTarget) error {
	ctx, cancel := context.WithTimeout(ctx, defaultStartupTimeout)
	defer cancel()

	host, err := target.Host(ctx)
	if err != nil {
		return err
	}

	port, err := target.MappedPort(ctx, "5432")
	if err != nil {
		return err
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port.Int(), ws.DatabaseUser, ws.DatabaseName)
	connConf, err := pgx.ParseConfig(connStr)
	if err != nil {
		return err
	}

LOOP:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		conn, err := pgx.ConnectConfig(ctx, connConf)
		if err != nil {
			continue
		}

		if err := conn.Ping(ctx); err == nil {
			break LOOP
		}
	}

	return nil
}
