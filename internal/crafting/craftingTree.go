package crafting

import "github.com/thestuckster/gopherfacts/pkg/items"

type CraftTreeNode struct {
	Value *items.ItemMetaData
	Leafs []CraftTreeNode
}

func NewCraftingTree(itemCodeToCraft string, itemData map[string]items.ItemMetaData) *CraftTreeNode {

	item, ok := itemData[itemCodeToCraft]
	if !ok {
		return nil
	}

	node := CraftTreeNode{
		Value: &item,
	}

	leafs := make([]CraftTreeNode, 0)
	for _, item := range item.Craft.Items {
		leafs = append(leafs, *NewCraftingTree(item.Code, itemData))
	}

	node.Leafs = leafs
	return &node
}
