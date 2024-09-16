package internal

import "github.com/thestuckster/gopherfacts/pkg/resources"

func BuildResourceMapByDropItemCode(r []resources.Resource) map[string]resources.Resource {
	dropMap := make(map[string]resources.Resource)

	for _, res := range r {
		drops := res.Drops
		for _, drop := range drops {
			itemCode := drop.Code

			dropMap[itemCode] = res
		}
	}

	return dropMap
}
