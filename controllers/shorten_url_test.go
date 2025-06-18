package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	. "github.com/hotrungnhan/surl/controllers"
	services_mocks "github.com/hotrungnhan/surl/generated/mocks/services"
	"github.com/hotrungnhan/surl/models"
	"github.com/hotrungnhan/surl/serializers"
	"github.com/hotrungnhan/surl/services"
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"

	"github.com/gofiber/fiber/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/do/v2"
)

var _ = Describe("ShortenUrlController", func() {
	var (
		app                   *fiber.App
		router                fiber.Router
		mockShortenUrlService *services_mocks.MockShortenUrlService
		injector              do.Injector
	)

	BeforeEach(func() {
		injector = do.New()
		mockShortenUrlService = &services_mocks.MockShortenUrlService{}
		// Provide mocks to the injector
		do.ProvideValue[services.ShortenUrlService](injector, mockShortenUrlService)

		var errNew error
		controller, errNew := NewShortenUrlController(injector)

		app, router = AppTest()

		controller.RegisterEndpoint(router)

		Expect(errNew).To(BeNil())
	})

	Describe("Get Public - GET /shortlinks/:id", func() {
		It("bind invalid params", func() {
			req, err := http.NewRequest("GET", "/shortlinks/123456", nil)

			mockShortenUrlService.On("Get", mock.Anything, &services.GetShortenUrlParams{ID: "test-id"}).Return(nil, nil)

			res, err := app.Test(req)

			resBodyBytes, _ := io.ReadAll(res.Body)

			Expect(err).To(BeNil())

			Expect(res.StatusCode).To(Equal(400))

			Expect(gjson.GetBytes(resBodyBytes, "message").String()).To(ContainSubstring("len"))

		})
		It("Correct", func() {
			url := "https://example.com"
			req, err := http.NewRequest("GET", "/shortlinks/12345678", nil)

			record := models.ShortenUrl{ID: "12345678", OriginalUrl: url}
			mockShortenUrlService.On("Get", mock.Anything, &services.GetShortenUrlParams{ID: "12345678"}).Return(
				serializers.NewShortUrlDetailSerializer(&record),
				nil,
			)

			res, err := app.Test(req)

			Expect(err).To(BeNil())

			Expect(res.StatusCode).To(Equal(302))

			Expect(res.Header.Get(fiber.HeaderLocation)).To(Equal(url))
		})
	})

	Describe("Get - GET /api/shortlinks/:id", func() {
		It("bind invalid params", func() {
			req, err := http.NewRequest("GET", "/api/shortlinks/123456", nil)

			mockShortenUrlService.On("Get", mock.Anything, &services.GetShortenUrlParams{ID: "test-id"}).Return(nil, nil)

			res, err := app.Test(req)

			resBodyBytes, _ := io.ReadAll(res.Body)

			Expect(err).To(BeNil())

			Expect(res.StatusCode).To(Equal(400))
			Expect(gjson.GetBytes(resBodyBytes, "message").String()).To(ContainSubstring("len"))

		})
		It("Correct", func() {
			id := "12345678"
			url := "https://example.com"
			req, err := http.NewRequest("GET", fmt.Sprintf("/api/shortlinks/%s", id), nil)

			record := models.ShortenUrl{ID: id, OriginalUrl: url}
			mockShortenUrlService.On("Get", mock.Anything, &services.GetShortenUrlParams{ID: id}).Return(
				serializers.NewShortUrlDetailSerializer(&record),
				nil,
			)

			res, err := app.Test(req)

			Expect(err).To(BeNil())

			Expect(res.StatusCode).To(Equal(200))

			var resp map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&resp)
			Expect(err).To(BeNil())

			Expect(resp).To(HaveKeyWithValue("original_url", url))
			Expect(resp).To(HaveKeyWithValue("id", id))
			Expect(resp).To(HaveKey("created_at"))
		})
	})

	Describe("Add - POST /api/shortlinks", func() {
		It("bind invalid params", func() {
			jsonBody, err := json.Marshal(map[string]string{
				"original_url": "",
			})
			Expect(err).To(BeNil())

			req, err := http.NewRequest("POST", "/api/shortlinks", bytes.NewBuffer(jsonBody))

			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			mockShortenUrlService.On("Get", mock.Anything, &services.AddShortenUrlParams{OriginalUrl: ""}).Return(nil, nil)

			res, err := app.Test(req)

			resBodyBytes, _ := io.ReadAll(res.Body)

			Expect(err).To(BeNil())

			Expect(res.StatusCode).To(Equal(400))

			Expect(gjson.GetBytes(resBodyBytes, "message").String()).To(ContainSubstring("required"))

		})
		It("bind invalid params #2", func() {
			jsonBody, err := json.Marshal(map[string]string{
				"original_url": "not google",
			})
			Expect(err).To(BeNil())

			req, err := http.NewRequest("POST", "/api/shortlinks", bytes.NewBuffer(jsonBody))

			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			mockShortenUrlService.On("Get", mock.Anything, &services.AddShortenUrlParams{OriginalUrl: ""}).Return(nil, nil)

			res, err := app.Test(req)

			Expect(err).To(BeNil())

			Expect(res.StatusCode).To(Equal(400))

		})

		It("Correct", func() {
			url := "https://example.com"
			jsonBody, err := json.Marshal(map[string]string{
				"original_url": url,
			})
			Expect(err).To(BeNil())

			req, err := http.NewRequest("POST", "/api/shortlinks", bytes.NewBuffer(jsonBody))
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			record := models.ShortenUrl{ID: "12345678", OriginalUrl: url}

			mockShortenUrlService.On("Add", mock.Anything, &services.AddShortenUrlParams{OriginalUrl: url}).Return(
				serializers.NewShortUrlDetailSerializer(&record),
				nil,
			)

			res, err := app.Test(req)

			Expect(err).To(BeNil())

			Expect(res.StatusCode).To(Equal(201))
		})
	})

})
