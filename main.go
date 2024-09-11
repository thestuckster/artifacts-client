package main

import (
	"artifacts-client/internal"
	"artifacts-client/internal/devHelpers"
	"artifacts-client/internal/planner"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/thestuckster/gopherfacts/pkg/clients"
	"github.com/thestuckster/gopherfacts/pkg/items"
	"github.com/thestuckster/gopherfacts/pkg/maps"
	"github.com/thestuckster/gopherfacts/pkg/monsters"
	"os"
)

var logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

func main() {
	logger.Info().Msg("Starting artifacts client")

	token := os.Getenv("TOKEN")
	sdk := clients.NewClient(&token)

	checkArtifactsServerStatus(sdk)

	itemData, _, _ := preLoadData()
	logger.Debug().Msg("Parsing item information into hashmaps for faster lookups")
	itemsByCode := internal.BuildItemByCode(itemData)
	_ = internal.BuildItemBySkill(itemData)

	//Throw away testing code
	mostRequiredMaterials := devHelpers.FindItemsWithMostUniqueCraftingRequirements(itemData)
	fmt.Printf("Item with most required materials: %+v\n", mostRequiredMaterials)

	planner.BuildPlanForItem("multislimes_sword", itemsByCode)

}

func checkArtifactsServerStatus(sdk *clients.GopherFactClient) {
	logger.Info().Msg("Checking Artifacts server status")
	statusInfo, err := sdk.CheckServerStatus()
	if err != nil {
		logger.Error().Err(err).Msg("Artifacts Server isnt healthy, check back later.")
		os.Exit(1)
	}

	logger.Info().Msg("Artifacts server status is healthy")
	logger.Debug().Msgf("Server status details: %+v", statusInfo)
}

func preLoadData() (items []items.ItemMetaData, maps []maps.MapData, monsters []monsters.Monster) {
	logger.Info().Msg("Preloading important game data...")
	items = loadItemData()
	maps = loadMapData()
	monsters = loadMonsterData()

	logger.Info().Msg("Preloading finished!")
	return items, maps, monsters
}

func loadItemData() []items.ItemMetaData {
	logger.Info().Msg("Loading ALL item data")
	itemData, err := items.GetAllItemData()
	if err != nil {
		logger.Error().Err(err).Msg("Error loading item data")
		os.Exit(1)
	}

	logger.Info().Msgf("Loaded %d items", len(itemData))
	return itemData
}

func loadMapData() []maps.MapData {
	logger.Info().Msg("Loading ALL map data")
	mapData, err := maps.GetAllMapData()
	if err != nil {
		logger.Error().Err(err).Msg("Error loading map data")
		os.Exit(1)
	}

	logger.Info().Msgf("Loaded %d map tiles", len(mapData))
	return mapData
}

func loadMonsterData() []monsters.Monster {
	logger.Info().Msg("Loading ALL monsters")
	monsterData, err := monsters.GetAllMonsterData()
	if err != nil {
		logger.Error().Err(err).Msg("Error loading monster data")
		os.Exit(1)
	}

	logger.Info().Msgf("Loaded %d monsters", len(monsterData))
	return monsterData
}
