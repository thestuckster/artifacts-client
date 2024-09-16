package planner

import (
	"artifacts-client/internal"
	"artifacts-client/internal/crafting"
	"errors"
	"fmt"
	pq "github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/emirpasic/gods/utils"
	"github.com/thestuckster/gopherfacts/pkg/items"
	"slices"
)

type ItemCraftingPlan struct {
}

func BuildPlanForItem(itemCode string, itemData map[string]items.ItemMetaData) (*ItemCraftingPlan, error) {

	_, ok := itemData[itemCode]
	if !ok {
		return nil, errors.New(fmt.Sprintf("No data for item %s in supplied itemData map", itemCode))
	}

	craftingTree := crafting.NewCraftingTree(itemCode, itemData)
	_ = depthFirstTraversal(craftingTree)

	return nil, nil
}

func buildTaskQueue(items []crafting.CraftTreeNode) *pq.Queue {
	queue := pq.NewWith(sortByPriority)

	p := len(items)
	for idx, _ := range items {

		//item.Value.Code
		queue.Enqueue(internal.ArtifactsTask{
			Priority: p - idx,
		})
	}

	return queue
}

func sortByPriority(a, b any) int {
	pA := a.(internal.ArtifactsTask).Priority
	pb := b.(internal.ArtifactsTask).Priority
	return -utils.IntComparator(pA, pb)
}

func validateItemIsObtainable(item crafting.CraftTreeNode) bool {
	//TODO: theres lots of associations to check here... we might want to move this to its own file because its going
	// to end up being a monster.

	//steps to decide if something is obtainable
	// is it a drop?
	// yes - what level is the monster?
	// is it a gatherable?
	// yes - what level is the resource node?
	// required > ours = how big is the gap? can we grind to get there?
	// does it have to be crafted?
	// yes - is our crafting skill high enough?
	// required > ours = how big is the gap? can we grind to get there?
	// no to all 3 - Can we buy it?
	// yes - is it in stock? do we have enough gold on our person? do we have enough gold in the bank?
	// no to everything! bail out, this item can't currently be obtained by us.

	return true
}

func depthFirstTraversal(node *crafting.CraftTreeNode) []crafting.CraftTreeNode {
	queue := make([]crafting.CraftTreeNode, 0)
	dfs(node, queue)
	return queue
}

func dfs(node *crafting.CraftTreeNode, l []crafting.CraftTreeNode) {
	if node == nil {
		return
	}

	for _, leaf := range node.Leafs {
		dfs(&leaf, l)
	}

	slices.Insert(l, 0, *node)
}
