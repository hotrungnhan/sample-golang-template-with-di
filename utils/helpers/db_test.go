package helpers_test

import (
	"context"
	"github.com/hotrungnhan/surl/utils/injects"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/samber/do/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var _ = ginkgo.Describe("NewDatabase with real Postgres container", ginkgo.Ordered, func() {
	var (
		injector          do.Injector
		ctx               context.Context
		cancel            context.CancelFunc
		postgresContainer testcontainers.Container
	)
	var _ = ginkgo.BeforeAll(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)

		req := testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_DB":       "postgres",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp").
				WithStartupTimeout(60 * time.Second),
		}

		var err error
		postgresContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		host, err := postgresContainer.Host(ctx)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		injector = do.New()
		cfg := injects.ApplicationConfig{
			MasterDB: injects.DBConfig{
				Host:     host,
				Port:     int(mappedPort.Int()),
				Username: "postgres",
				Password: "postgres",
				Name:     "postgres",
				SslMode:  "disable",
			},
			SlaveDB: injects.DBConfig{
				Host:     host,
				Port:     int(mappedPort.Int()),
				Username: "postgres",
				Password: "postgres",
				Name:     "postgres",
				SslMode:  "disable",
			},
		}
		do.ProvideValue(injector, cfg)
	})

	var _ = ginkgo.AfterAll(func() {
		if postgresContainer != nil {
			_ = postgresContainer.Terminate(ctx)
		}
		cancel()
	})

	ginkgo.It("successfully creates DB with master and slave connections", func() {
		db, err := injects.NewDatabase(injector)

		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(db).ToNot(gomega.BeNil())
		gomega.Expect(db.Master).ToNot(gomega.BeNil())
		gomega.Expect(db.Slave).ToNot(gomega.BeNil())
	})

	ginkgo.It("returns error if master DB config DSN is nil", func() {
		cfg := injects.ApplicationConfig{
			MasterDB: injects.DBConfig{}, // empty means DSN is nil
		}
		do.OverrideValue(injector, cfg)

		db, err := injects.NewDatabase(injector)

		gomega.Expect(db).To(gomega.BeNil())
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(err.Error()).To(gomega.ContainSubstring("Master DB configuration is not set"))
	})

	ginkgo.It("uses master DSN for slave if slave DSN is nil", func() {
		cfg := injects.ApplicationConfig{
			MasterDB: injects.DBConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "postgres",
				Password: "postgres",
				Name:     "postgres",
				SslMode:  "disable",
			},
			SlaveDB: injects.DBConfig{}, // empty slave means DSN nil
		}
		do.OverrideValue(injector, cfg)

		db, err := injects.NewDatabase(injector)

		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(db).ToNot(gomega.BeNil())
		// Slave will fallback to master DSN internally
	})
})

var _ = ginkgo.Describe("NewDatabase using Postgres containers", ginkgo.Ordered, func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		injector        do.Injector
		masterContainer testcontainers.Container
		slaveContainer  testcontainers.Container
		masterHost      string
		masterPort      int
		slaveHost       string
		slavePort       int
	)

	ginkgo.BeforeAll(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)

		createPostgresContainer := func(dbName string) (testcontainers.Container, string, int, error) {
			req := testcontainers.ContainerRequest{
				Image:        "postgres:15-alpine",
				ExposedPorts: []string{"5432/tcp"},
				Env: map[string]string{
					"POSTGRES_USER":     "postgres",
					"POSTGRES_PASSWORD": "postgres",
					"POSTGRES_DB":       dbName,
				},
				WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
			}

			container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
				ContainerRequest: req,
				Started:          true,
			})
			if err != nil {
				return nil, "", 0, err
			}

			host, err := container.Host(ctx)
			if err != nil {
				return nil, "", 0, err
			}

			port, err := container.MappedPort(ctx, "5432")
			if err != nil {
				return nil, "", 0, err
			}

			return container, host, port.Int(), nil
		}

		var err error
		masterContainer, masterHost, masterPort, err = createPostgresContainer("master")
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		slaveContainer, slaveHost, slavePort, err = createPostgresContainer("slave")
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	ginkgo.AfterAll(func() {
		if masterContainer != nil {
			_ = masterContainer.Terminate(ctx)
		}
		if slaveContainer != nil {
			_ = slaveContainer.Terminate(ctx)
		}
		cancel()
	})

	ginkgo.BeforeEach(func() {
		injector = do.New()
	})

	ginkgo.It("successfully creates DB with both master and slave", func() {
		cfg := injects.ApplicationConfig{
			MasterDB: injects.DBConfig{
				Host:     masterHost,
				Port:     masterPort,
				Username: "postgres",
				Password: "postgres",
				Name:     "master",
				SslMode:  "disable",
			},
			SlaveDB: injects.DBConfig{
				Host:     slaveHost,
				Port:     slavePort,
				Username: "postgres",
				Password: "postgres",
				Name:     "slave",
				SslMode:  "disable",
			},
		}
		do.ProvideValue(injector, cfg)

		db, err := injects.NewDatabase(injector)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(db).ToNot(gomega.BeNil())

		masterSQL, err := db.Master.DB()
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(masterSQL.Ping()).To(gomega.Succeed())

		slaveSQL, err := db.Slave.DB()
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(slaveSQL.Ping()).To(gomega.Succeed())
	})
})
