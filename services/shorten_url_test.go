package services_test

import (
	"context"
	"errors"

	mapper "github.com/hotrungnhan/go-automapper"
	repositories_mocks "github.com/hotrungnhan/surl/generated/mocks/repositories"
	"github.com/hotrungnhan/surl/models"
	"github.com/hotrungnhan/surl/repositories"
	"github.com/hotrungnhan/surl/serializers"
	. "github.com/hotrungnhan/surl/services"
	"github.com/hotrungnhan/surl/utils/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/do/v2"
)

var _ = Describe("ShortenUrlService", func() {
	var (
		service               ShortenUrlService
		mockShortenRepository *repositories_mocks.MockShortenUrlRepository
		injector              do.Injector
	)

	BeforeEach(func() {
		injector = do.New()
		mockShortenRepository = &repositories_mocks.MockShortenUrlRepository{}
		// Provide mocks to the injector
		do.ProvideValue[repositories.ShortenUrlRepository](injector, mockShortenRepository)

		var errNew error
		service, errNew = NewShortenUrlService(injector)
		Expect(errNew).To(BeNil())
	})

	Describe("Get", func() {
		var (
			ctx    context.Context
			params *GetShortenUrlParams
			record *models.ShortenUrl
		)

		BeforeEach(func() {
			ctx = context.Background()
			params = &GetShortenUrlParams{ID: "test-id"}
			record = &models.ShortenUrl{ID: "test-id", OriginalUrl: "https://example.com"}
		})

		It("should return the record if found", func() {
			mockShortenRepository.On("Get", ctx, &models.ShortenUrlFilterParams{ID: &params.ID}).Return(record, nil)
			data, err := service.Get(ctx, params)
			Expect(err).To(BeNil())
			Expect(data).To(BeAssignableToTypeOf(mapper.MustMap[*models.ShortenUrl, *serializers.ShortUrlDetailSerializer](mapper.Global, record)))
		})

		It("should return InternalServerError if repo returns an error", func() {
			mockShortenRepository.On("Get", ctx, &models.ShortenUrlFilterParams{ID: &params.ID}).Return(nil, errors.New("repo error"))
			_, err := service.Get(ctx, params)
			Expect(err).To(Equal(types.InternalServerError))
		})

		It("should return NoContentError if record is not found", func() {
			mockShortenRepository.On("Get", ctx, &models.ShortenUrlFilterParams{ID: &params.ID}).Return(nil, nil)
			_, err := service.Get(ctx, params)
			Expect(err).To(Equal(types.NoContentError))
		})
	})

	Describe("Add", func() {
		var (
			ctx    context.Context
			params *AddShortenUrlParams
			record *models.ShortenUrl
		)

		BeforeEach(func() {
			ctx = context.Background()
			params = &AddShortenUrlParams{OriginalUrl: "https://example.com"}
			record = &models.ShortenUrl{ID: "test-id", OriginalUrl: "https://example.com"}
		})

		It("should return the record if already exists", func() {
			mockShortenRepository.On("Get", ctx, params.ToFilter()).Return(record, nil)
			data, err := service.Add(ctx, params)
			Expect(err).To(BeNil())
			Expect(data).To(BeAssignableToTypeOf(mapper.MustMap[*models.ShortenUrl, *serializers.ShortUrlSerializer](mapper.Global, record)))
		})

		It("should create a new record if not exists", func() {
			mockShortenRepository.On("Get", ctx, params.ToFilter()).Return(nil, nil)
			mockShortenRepository.On("Add", ctx, params.ToCreateModel()).Return(record, nil)
			data, err := service.Add(ctx, params)
			Expect(err).To(BeNil())
			Expect(data).To(BeAssignableToTypeOf(mapper.MustMap[*models.ShortenUrl, *serializers.ShortUrlSerializer](mapper.Global, record)))
		})

		It("should return InternalServerError if repo returns an error", func() {
			mockShortenRepository.On("Get", ctx, params.ToFilter()).Return(nil, errors.New("repo error"))
			_, err := service.Add(ctx, params)
			Expect(err).To(Equal(types.InternalServerError))
		})
	})
})
