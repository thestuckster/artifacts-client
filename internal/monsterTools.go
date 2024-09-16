package internal

import "github.com/thestuckster/gopherfacts/pkg/monsters"

func BuildMonsterMapByDropItemCode(m []monsters.Monster) map[string]monsters.Monster {
	dropMap := make(map[string]monsters.Monster)
	for _, monster := range m {
		drops := monster.Drops
		for _, drop := range drops {
			dropMap[drop.Code] = monster
		}
	}

	return dropMap
}
