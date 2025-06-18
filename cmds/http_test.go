package cmds_test

import (
	"github.com/hotrungnhan/surl/cmds"
	"github.com/hotrungnhan/surl/utils/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HttpServer", func() {
	var server types.Service
	var opts cmds.HttpServerOptions

	BeforeEach(func() {
		opts = cmds.HttpServerOptions{
			Address: "127.0.0.1",
			Port:    8089, // use a non-privileged port for testing
		}
		server = cmds.NewHttpServer(opts)
	})

	It("should start and respond to a DNS query", func() {
		go func() {
			err := server.Serve()
			Expect(err).ToNot(HaveOccurred())
		}()
	})
})
