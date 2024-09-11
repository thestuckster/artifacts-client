package planner

import (
	"artifacts-client/internal/crafting"
	"container/list"
	"fmt"
	"github.com/thestuckster/gopherfacts/pkg/items"
)

func BuildPlanForItem(itemCode string, itemData map[string]items.ItemMetaData) {

	craftingTree := crafting.NewCraftingTree(itemCode, itemData)
	buildQueue := depthFirstTraversal(craftingTree)
	fmt.Printf("", buildQueue)
}

func depthFirstTraversal(node *crafting.CraftTreeNode) *list.List {
	queue := list.New()
	dfs(node, queue)
	return queue
}

func dfs(node *crafting.CraftTreeNode, l *list.List) {
	if node == nil {
		return
	}

	for _, leaf := range node.Leafs {
		dfs(&leaf, l)
	}

	l.PushFront(node)
}
