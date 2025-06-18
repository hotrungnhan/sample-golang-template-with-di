package helpers_test

import (
	"github.com/hotrungnhan/surl/utils/helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type TestStruct struct {
	Name  string `validate:"required" default:"John Doe"`
	Email string `validate:"required,email" default:"johndoe@example.com"`
}

var _ = Describe("Utils Validation and Defaults", func() {
	var testObj TestStruct

	BeforeEach(func() {
		testObj = TestStruct{}
	})

	Describe("Validate", func() {
		It("should pass validation with valid fields", func() {
			testObj.Name = "Alice"
			testObj.Email = "alice@example.com"
			err := helpers.Validate(testObj)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail validation with invalid fields", func() {
			testObj.Name = ""
			testObj.Email = "invalid-email"
			err := helpers.Validate(testObj)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("SetDefaults", func() {
		It("should set default values on empty struct", func() {
			err := helpers.SetDefaults(&testObj)
			Expect(err).ToNot(HaveOccurred())
			Expect(testObj.Name).To(Equal("John Doe"))
			Expect(testObj.Email).To(Equal("johndoe@example.com"))
		})
	})

	Describe("ValidateAndDefault", func() {
		It("should fail validation if struct is invalid", func() {
			err := helpers.ValidateAndDefault(&testObj)
			Expect(err).To(HaveOccurred())
		})

		It("should pass validation and set defaults if struct is valid", func() {
			testObj.Name = "Bob"
			testObj.Email = "bob@example.com"
			err := helpers.ValidateAndDefault(&testObj)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
