// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	"encoding/json"
	"io/ioutil"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var populateSearchResults = func() []SearchResult {
	content, err := ioutil.ReadFile("../testing/search_results.json")
	Expect(err).NotTo(HaveOccurred())
	Expect(len(content)).NotTo(Equal(0))

	var lists []unstructured.UnstructuredList
	err = json.Unmarshal(content, &lists)
	Expect(err).NotTo(HaveOccurred())

	var results []SearchResult
	for index, list := range lists {
		Expect(list.Items).To(HaveLen(index + 1))
		results = append(results, SearchResult{
			List: list.DeepCopy(),
		})
	}
	return results
}

var _ = Describe("SearchResult", func() {

	Context("ToStarlarkValue", func() {

		It("returns a dictionary of size equal to number of struct elements", func() {
			sr := SearchResult{ListKind: "PodList"}
			val := sr.ToStarlarkValue()
			Expect(val).To(BeAssignableToTypeOf(&starlarkstruct.Struct{}))
			Expect(val.AttrNames()).To(HaveLen(7))
		})

		sr := SearchResult{
			ListKind:     "NodeList",
			ResourceName: "nodes",
			ResourceKind: "Node",
			Namespace:    "",
			Namespaced:   false,
		}

		DescribeTable("String types", func(typeDescription, stringVal string) {
			structVal := sr.ToStarlarkValue()
			val, err := structVal.Attr(typeDescription)
			Expect(err).NotTo(HaveOccurred())

			strVal, _ := val.(starlark.String)
			Expect(strVal.GoString()).To(Equal(stringVal))
		},
			Entry("", "ListKind", "NodeList"),
			Entry("", "ResourceName", "nodes"),
			Entry("", "ResourceKind", "Node"),
			Entry("", "Namespace", ""),
		)

		Context("For Namespaced", func() {
			It("creates a dictionary with Namespaced value", func() {
				dict := sr.ToStarlarkValue()
				val, err := dict.Attr("Namespaced")
				Expect(err).NotTo(HaveOccurred())

				boolVal, _ := val.(starlark.Bool)
				Expect(boolVal.Truth()).To(Equal(starlark.False))
			})
		})

		Context("For List", func() {

			It("returns a starlark struct", func() {
				sr = searchResults[0]
				structVal := sr.ToStarlarkValue()
				Expect(structVal).To(BeAssignableToTypeOf(&starlarkstruct.Struct{}))

				val, err := structVal.Attr("List")
				Expect(err).NotTo(HaveOccurred())

				_, ok := val.(*starlarkstruct.Struct)
				Expect(ok).To(BeTrue())
			})

			It("contains a starlark struct with the Object key", func() {
				sr = searchResults[0]
				structVal := sr.ToStarlarkValue()
				val, _ := structVal.Attr("List")
				listVal, _ := val.(*starlarkstruct.Struct)

				objVal, err := listVal.Attr("Object")
				Expect(err).NotTo(HaveOccurred())
				objStructVal, ok := objVal.(*starlarkstruct.Struct)
				Expect(ok).To(BeTrue())
				Expect(objStructVal).To(BeAssignableToTypeOf(&starlarkstruct.Struct{}))
			})

			It("contains a starlark list with the Items key", func() {
				sr = searchResults[0]
				structVal := sr.ToStarlarkValue()
				val, _ := structVal.Attr("List")
				listVal, _ := val.(*starlarkstruct.Struct)

				itemsVal, err := listVal.Attr("Items")
				Expect(err).NotTo(HaveOccurred())
				itemsListVal, ok := itemsVal.(*starlark.List)
				Expect(ok).To(BeTrue())
				Expect(itemsListVal).To(BeAssignableToTypeOf(&starlark.List{}))

				Expect(itemsListVal.Len()).To(Equal(1))
			})

			Context("For each list entry", func() {

				var listStructVal *starlarkstruct.Struct

				BeforeEach(func() {
					sr := searchResults[0]
					structVal := sr.ToStarlarkValue()
					val, _ := structVal.Attr("List")
					listVal, _ := val.(*starlarkstruct.Struct)
					itemsVal, _ := listVal.Attr("Items")
					itemsListVal, _ := itemsVal.(*starlark.List)
					listStructVal, _ = itemsListVal.Index(0).(*starlarkstruct.Struct)
				})

				It("returns a starlark string for a string value", func() {
					kindAttrVal, err := listStructVal.Attr("kind")
					Expect(err).NotTo(HaveOccurred())
					if kind, ok := kindAttrVal.(starlark.String); !ok {
						Expect(kind.GoString()).To(Equal("Service"))
					} else {
						Expect(ok).To(BeTrue())
					}

					apiVersionVal, err := listStructVal.Attr("apiVersion")
					Expect(err).NotTo(HaveOccurred())
					if version, ok := apiVersionVal.(starlark.String); ok {
						Expect(version.GoString()).To(Equal("v1"))
					} else {
						Expect(ok).To(BeTrue())
					}
				})

				It("returns a starlark struct for a map value", func() {
					metadataAttrVal, err := listStructVal.Attr("metadata")
					Expect(err).NotTo(HaveOccurred())
					metadata, ok := metadataAttrVal.(*starlarkstruct.Struct)
					Expect(ok).To(BeTrue())

					labelVal, err := metadata.Attr("labels")
					Expect(err).NotTo(HaveOccurred())
					Expect(labelVal).To(BeAssignableToTypeOf(&starlarkstruct.Struct{}))
				})

				It("returns a starlark list for an array value", func() {
					specAttrVal, err := listStructVal.Attr("spec")
					Expect(err).NotTo(HaveOccurred())
					spec, ok := specAttrVal.(*starlarkstruct.Struct)
					Expect(ok).To(BeTrue())

					portsVal, err := spec.Attr("ports")
					Expect(err).NotTo(HaveOccurred())
					Expect(portsVal).To(BeAssignableToTypeOf(&starlark.List{}))

					ports, ok := portsVal.(*starlark.List)
					Expect(ok).To(BeTrue())
					Expect(ports.Len()).To(Equal(3))
					Expect(ports.Index(0)).To(BeAssignableToTypeOf(&starlarkstruct.Struct{}))
				})
			})
		})
	})
})
