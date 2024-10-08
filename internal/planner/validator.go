package planner

import (
	"artifacts-client/internal/crafting"
	"github.com/thestuckster/gopherfacts/pkg/clients"
	"github.com/thestuckster/gopherfacts/pkg/exchange"
	"github.com/thestuckster/gopherfacts/pkg/items"
	"github.com/thestuckster/gopherfacts/pkg/resources"
)

type CraftingPlanValidator struct {
	itemsByCode, itemsBySkill map[string]items.ItemMetaData
	resourcesByCode           map[string]resources.Resource
}

func (v *CraftingPlanValidator) IsItemObtainable(characterName string, sdk *clients.GopherFactClient, node crafting.CraftTreeNode) bool {

	itemToCraft := node.Value

	if itemToCraft.Type == "resource" {
		ok, err := v.validateResource(characterName, sdk, itemToCraft)
		if err != nil {
			return false
		}

		return ok
	}

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

func (v *CraftingPlanValidator) validateResource(characterName string, sdk *clients.GopherFactClient, item *items.ItemMetaData) (bool, error) {

	characterInfo, err := v.getCharacterInfo(characterName, sdk)
	if err != nil {
		return false, err
	}

	if v.isItemInCharacterInventory(characterInfo, item) {
		return true, nil
	}

	if v.isItemInBank(characterInfo, sdk, item) {
		return true, nil
	}

	if v.isItemEquippedToCharacter(characterInfo, item) {
		return true, nil
	}

	if v.isItemPurchasable(characterInfo, sdk, item) {
		return true, nil
	}

	if v.isItemGatherable(characterInfo, item) {
		return true, nil
	}

	return false, nil
}

func (v *CraftingPlanValidator) validateCraftable(characterName string, sdk *clients.GopherFactClient, node crafting.CraftTreeNode) (bool, error) {
	return false, nil
}

func (v *CraftingPlanValidator) isItemInCharacterInventory(characterInfo *clients.CharacterSchema, item *items.ItemMetaData) bool {
	for _, inventory := range characterInfo.Inventory {
		if inventory.Code == item.Code {
			return true
		}
	}

	return false
}

func (v *CraftingPlanValidator) isItemEquippedToCharacter(characterInfo *clients.CharacterSchema, item *items.ItemMetaData) bool {

	shield := characterInfo.ShieldSlot
	boots := characterInfo.BootsSlot
	body := characterInfo.BodyArmorSlot
	amulet := characterInfo.AmuletSlot
	helmet := characterInfo.HelmetSlot
	leg := characterInfo.LegArmorSlot
	weapon := characterInfo.WeaponSlot
	a1 := characterInfo.Artifact1Slot
	a2 := characterInfo.Artifact2Slot
	a3 := characterInfo.Artifact3Slot
	c1 := characterInfo.Consumable1Slot
	c2 := characterInfo.Consumable2Slot
	r1 := characterInfo.Ring1Slot
	r2 := characterInfo.Ring2Slot

	equippedItems := []string{
		shield, boots, body, amulet, helmet, leg, weapon,
		a1, a2, a3, c1, c2, r1, r2,
	}

	for _, equippedItem := range equippedItems {
		//TODO: could be name, not sure right now
		if equippedItem == item.Code {
			return true
		}
	}

	return false
}

func (v *CraftingPlanValidator) isItemInBank(characterInfo *clients.CharacterSchema, sdk *clients.GopherFactClient, item *items.ItemMetaData) bool {

	//TODO: this hasn't been implemented in the sdk yet...
	bankItems, err := sdk.AccountClient.GetAllBankItems()
	if err != nil {
		return false
	}

	for _, i := range bankItems {
		if i.Code == item.Code {
			return true
		}
	}

	return false
}

func (v *CraftingPlanValidator) isItemPurchasable(characterInfo *clients.CharacterSchema, sdk *clients.GopherFactClient, item *items.ItemMetaData) bool {

	exchangeData, err := exchange.GetItemExchangeData(item.Code)
	if err != nil {
		//TODO: need to add logging so I know what went wrong if this fails
		return false
	}
	cost := exchangeData.BuyPrice

	haveGoldOnHand := characterInfo.Gold <= cost
	bankBalance, err := sdk.AccountClient.GetBankGold()
	if err != nil {
		return false
	}

	haveGoldInBank := bankBalance.Quantity <= cost

	if haveGoldOnHand == false && haveGoldInBank == false {
		return false
	}

	return true
}

func (v *CraftingPlanValidator) isItemGatherable(characterInfo *clients.CharacterSchema, item *items.ItemMetaData) bool {
	resource, ok := v.resourcesByCode[item.Code]
	if !ok {
		return false
	}

	switch resource.Skill {
	case "mining":
		return characterInfo.MiningLevel >= resource.Level
	case "woodcutting":
		return characterInfo.WoodcuttingLevel >= resource.Level
	case "fishing":
		return characterInfo.FishingLevel >= resource.Level
	}

	return false
}

func (v *CraftingPlanValidator) getCharacterInfo(characterName string, sdk *clients.GopherFactClient) (*clients.CharacterSchema, error) {
	return sdk.CharacterClient.GetCharacterInfo(characterName)
}
