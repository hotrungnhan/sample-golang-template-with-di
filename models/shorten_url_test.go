package models_test

import (
	"context"
	"github.com/hotrungnhan/surl/models"
	"github.com/hotrungnhan/surl/utils/helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ShortenUrl", func() {
	model := models.ShortenUrl{}

	Describe("Validate", func() {
		It("No error", func(ctx context.Context) {
			err := helpers.Validate(model)
			Expect(err).To(BeNil())

			// dnsServiceMock := services_mocks.NewMockDnsService(GinkgoT())
			// dnsServiceMock.On("Resolve", mock.Anything, mock.Anything).Return(nil, nil)
			// dnsServiceMock.Resolve(ctx, []dns.Question{})
		})
	})
	Describe("Set Default", func() {
		It("No error", func(ctx context.Context) {
			err := helpers.SetDefaults(&model)
			Expect(err).To(BeNil())
		})
	})
})
