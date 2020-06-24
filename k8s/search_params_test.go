package k8s

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var _ = Describe("SearchParams", func() {

	var searchParams SearchParams

	Context("Building a new instance from a Starlark struct", func() {

		var (
			input *starlarkstruct.Struct
			args  starlark.StringDict
		)

		It("returns a new instance of the SearchParams type", func() {
			input = starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{})
			searchParams = NewSearchParams(input)
			Expect(searchParams).To(BeAssignableToTypeOf(SearchParams{}))
		})

		Context("With kinds", func() {

			Context("In the input struct", func() {

				It("returns a new instance with kinds struct member populated", func() {
					args = starlark.StringDict{
						"kinds": starlark.String("deployments"),
					}
					input = starlarkstruct.FromStringDict(starlarkstruct.Default, args)
					searchParams = NewSearchParams(input)
					Expect(searchParams).To(BeAssignableToTypeOf(SearchParams{}))
					Expect(searchParams.kinds).To(HaveLen(1))
					Expect(searchParams.Kinds()).To(Equal("deployments"))
				})

				It("returns a new instance with kinds struct member populated", func() {
					args = starlark.StringDict{
						"kinds": starlark.NewList([]starlark.Value{starlark.String("deployments"), starlark.String("replicasets")}),
					}
					input = starlarkstruct.FromStringDict(starlarkstruct.Default, args)
					searchParams = NewSearchParams(input)
					Expect(searchParams).To(BeAssignableToTypeOf(SearchParams{}))
					Expect(searchParams.kinds).To(HaveLen(2))
					Expect(searchParams.Kinds()).To(Equal("deployments replicasets"))
				})
			})

			Context("not in the input struct", func() {

				It("returns a new instance with default value of kinds struct member populated", func() {
					input = starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{})
					searchParams = NewSearchParams(input)
					Expect(searchParams).To(BeAssignableToTypeOf(SearchParams{}))
					Expect(searchParams.kinds).To(HaveLen(0))
					Expect(searchParams.Kinds()).To(Equal(""))
				})
			})
		})

		Context("With namespaces", func() {

			Context("In the input struct", func() {

				It("returns a new instance with namespaces struct member populated", func() {
					args = starlark.StringDict{
						"namespaces": starlark.String("foo"),
					}
					input = starlarkstruct.FromStringDict(starlarkstruct.Default, args)
					searchParams = NewSearchParams(input)
					Expect(searchParams).To(BeAssignableToTypeOf(SearchParams{}))
					Expect(searchParams.namespaces).To(HaveLen(1))
					Expect(searchParams.Namespaces()).To(Equal("foo"))
				})

				It("returns a new instance with namespaces struct member populated", func() {
					args = starlark.StringDict{
						"namespaces": starlark.NewList([]starlark.Value{starlark.String("foo"), starlark.String("bar")}),
					}
					input = starlarkstruct.FromStringDict(starlarkstruct.Default, args)
					searchParams = NewSearchParams(input)
					Expect(searchParams).To(BeAssignableToTypeOf(SearchParams{}))
					Expect(searchParams.namespaces).To(HaveLen(2))
					Expect(searchParams.Namespaces()).To(Equal("foo bar"))
				})
			})

			Context("not in the input struct", func() {

				It("returns a new instance with default value of namespaces struct member populated", func() {
					input = starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{})
					searchParams = NewSearchParams(input)
					Expect(searchParams).To(BeAssignableToTypeOf(SearchParams{}))
					Expect(searchParams.namespaces).To(HaveLen(1))
					Expect(searchParams.Namespaces()).To(Equal("default"))
				})
			})
		})
	})
})
