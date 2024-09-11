package internal

import "github.com/thestuckster/gopherfacts/pkg/items"

func BuildItemBySkill(data []items.ItemMetaData) map[string][]items.ItemMetaData {
	m := make(map[string][]items.ItemMetaData)

	for _, item := range data {
		skill := item.Craft.Skill
		m[skill] = append(m[skill], item)
	}

	return m
}

func BuildItemByCode(data []items.ItemMetaData) map[string]items.ItemMetaData {
	m := make(map[string]items.ItemMetaData)
	for _, item := range data {
		m[item.Code] = item
	}

	return m
}
