package k8s

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("SearchResult", func() {

	Context("ToStarlarkValue", func() {

		Context("ListKind", func() {
			sr := SearchResult{ListKind: "PodList"}

			It("creates value object with ListKind value", func() {
				_ = sr.ToStarlarkValue()
			})
		})

		Context("For ResourceName", func() {

		})

		Context("For ResourceKind", func() {

		})

		Context("For Namespaced", func() {

		})

		Context("For Namespace", func() {

		})
	})
})
