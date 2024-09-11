// Package devHelpers are one off functions used to find and report data for development purposes. They are not meant to be used
// for client behavior
package devHelpers

import "github.com/thestuckster/gopherfacts/pkg/items"

// FindItemsWithMostUniqueCraftingRequirements Find the item that has the most unique items required for crafting
func FindItemsWithMostUniqueCraftingRequirements(data []items.ItemMetaData) items.ItemMetaData {

	var mostUnique items.ItemMetaData

	count := 0
	for _, item := range data {
		if count < len(item.Craft.Items) {
			count = len(item.Craft.Items)
			mostUnique = item
		}
	}

	return mostUnique
}
