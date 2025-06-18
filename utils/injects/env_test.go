package injects_test

import (
	"github.com/hotrungnhan/surl/utils/injects"
	"os"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/samber/do/v2"
)

var _ = ginkgo.Describe("ApplicationConfig and DBConfig", func() {
	ginkgo.BeforeEach(func() {
		// Clear env vars to isolate tests
		os.Unsetenv("MASTER_DB_HOST")
		os.Unsetenv("MASTER_DB_PORT")
		os.Unsetenv("MASTER_DB_USERNAME")
		os.Unsetenv("MASTER_DB_PASSWORD")
		os.Unsetenv("MASTER_DB_NAME")
		os.Unsetenv("MASTER_DB_SSL_MODE")
		os.Unsetenv("GO_ENV")
		os.Unsetenv("FEATURE_FLAGS_ENABLE_FEATURE_X")
	})

	ginkgo.It("should return nil DSN if all fields empty", func() {
		cfg := injects.DBConfig{}
		gomega.Expect(cfg.GetDSN()).To(gomega.BeNil())
	})

	ginkgo.It("should construct correct DSN string", func() {
		cfg := injects.DBConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "user",
			Password: "pass",
			Name:     "db",
			SslMode:  "disable",
		}

		dsn := cfg.GetDSN()
		gomega.Expect(dsn).ToNot(gomega.BeNil())
		gomega.Expect(*dsn).To(gomega.ContainSubstring("host=localhost"))
		gomega.Expect(*dsn).To(gomega.ContainSubstring("user=user"))
		gomega.Expect(*dsn).To(gomega.ContainSubstring("password=pass"))
		gomega.Expect(*dsn).To(gomega.ContainSubstring("dbname=db"))
		gomega.Expect(*dsn).To(gomega.ContainSubstring("port=5432"))
		gomega.Expect(*dsn).To(gomega.ContainSubstring("sslmode=disable"))
	})

	ginkgo.It("should load defaults and env overrides in NewAppConfig", func() {
		// Set some environment variables to override defaults
		os.Setenv("MASTER_DB_HOST", "dbhost")
		os.Setenv("MASTER_DB_PORT", "5433")
		os.Setenv("GO_ENV", "production")
		os.Setenv("FEATURE_FLAGS_ENABLE_FEATURE_X", "true")

		injector := do.New()

		cfg, err := injects.NewAppConfig(injector)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(cfg).ToNot(gomega.BeNil())

		// Defaults applied and env vars override some values
		gomega.Expect(cfg.MasterDB.Host).To(gomega.Equal("dbhost"))
		gomega.Expect(cfg.MasterDB.Port).To(gomega.Equal(5433))
		gomega.Expect(string(cfg.GoEnv)).To(gomega.Equal("production"))

		// Feature flag parsed from env string (should parse "true" to boolean true)
		gomega.Expect(cfg.FeatureFlags.ENABLE_FEATURE_X).To(gomega.BeTrue())

		// Defaults remain for other fields
		gomega.Expect(cfg.MasterDB.Username).To(gomega.Equal("postgres"))
		gomega.Expect(cfg.MasterDB.Password).To(gomega.Equal("postgres"))
		gomega.Expect(cfg.MasterDB.Name).To(gomega.Equal("example_db"))
		gomega.Expect(cfg.MasterDB.SslMode).To(gomega.Equal("disable"))
	})
})
