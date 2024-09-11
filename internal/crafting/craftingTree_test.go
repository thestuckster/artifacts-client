package crafting_test

import (
	"artifacts-client/internal"
	"artifacts-client/internal/crafting"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/thestuckster/gopherfacts/pkg/items"
)

var _ = Describe("CraftingTree", func() {
	var itemData map[string]items.ItemMetaData

	Describe("Building a tree", func() {
		rawItems, err := items.GetAllItemData()
		gomega.Expect(err).To(gomega.BeNil())
		itemData = internal.BuildItemByCode(rawItems)

		It("Should build a simple tree", func() {
			root := crafting.NewCraftingTree("copper", itemData)
			gomega.Expect(root).ToNot(gomega.BeNil())
			gomega.Expect(root.Value.Code).To(gomega.Equal("copper"))
			gomega.Expect(len(root.Leafs)).To(gomega.Equal(1))
		})

		It("Should build a complex tree", func() {
			root := crafting.NewCraftingTree("multislimes_sword", itemData)
			gomega.Expect(root).ToNot(gomega.BeNil())
			gomega.Expect(root.Value.Code).To(gomega.Equal("multislimes_sword"))
			gomega.Expect(len(root.Leafs)).To(gomega.Equal(6))

			var subtree crafting.CraftTreeNode
			for _, item := range root.Leafs {
				if item.Value.Code == "iron" {
					subtree = item
				}
			}

			gomega.Expect(subtree).ToNot(gomega.BeNil())
			gomega.Expect(subtree.Value.Code).To(gomega.Equal("iron"))
			gomega.Expect(len(subtree.Leafs)).To(gomega.Equal(1))
		})
	})
})
