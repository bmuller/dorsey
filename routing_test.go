package dorsey

import (
	//	. "dorsey"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routing", func() {
	Describe("A string route", func() {
		Context("With a named path", func() {
			It("should make a new path", func() {
				Expect(1).To(Equal(1))
			})
		})
	})
})
