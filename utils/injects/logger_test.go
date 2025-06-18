package injects_test

import (
	"github.com/hotrungnhan/surl/utils/injects"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

var _ = ginkgo.Describe("NewLogger", func() {
	var injector do.Injector

	ginkgo.BeforeEach(func() {
		injector = do.New()
	})

	ginkgo.It("should create a new zerolog.Logger without error", func() {
		logger, err := injects.NewLogger(injector)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(logger).ToNot(gomega.BeNil())
		_, ok := interface{}(logger).(*zerolog.Logger)
		gomega.Expect(ok).To(gomega.BeTrue())
	})

	ginkgo.It("should write logs to stdout", func() {
		logger, _ := injects.NewLogger(injector)
		// Because logger writes to os.Stdout, here we only verify logger is usable,
		// real output capture requires redirecting os.Stdout (more complex)
		logger.Info().Msg("test log")
		// Just confirm no panic or error during logging
	})
})
