package main

import (
	"artifacts-client/internal"
	"artifacts-client/internal/planner"
	"github.com/rs/zerolog"
	"github.com/thestuckster/gopherfacts/pkg/clients"
	"github.com/thestuckster/gopherfacts/pkg/items"
	"github.com/thestuckster/gopherfacts/pkg/maps"
	"github.com/thestuckster/gopherfacts/pkg/monsters"
	"github.com/thestuckster/gopherfacts/pkg/resources"
	"os"
	"time"
)

var logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

func main() {
	logger.Info().Msg("Starting artifacts client")

	token := os.Getenv("TOKEN")
	sdk := clients.NewClient(&token)

	checkArtifactsServerStatus(sdk)

	itemData, _, _, resourceData := preLoadData()
	logger.Debug().Msg("Parsing item information")
	itemsByCode := internal.BuildItemByCode(itemData)
	_ = internal.BuildItemBySkill(itemData)

	logger.Debug().Msg("Parsing resource information")
	_ = internal.BuildResourceMapByDropItemCode(resourceData)

	planner.BuildPlanForItem("multislimes_sword", itemsByCode)

	//wait for all cooldowns to process
	waitForCooldowns()

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

func preLoadData() (items []items.ItemMetaData, maps []maps.MapData, monsters []monsters.Monster, resources []resources.Resource) {
	logger.Info().Msg("Preloading important game data...")
	items = loadItemData()
	maps = loadMapData()
	monsters = loadMonsterData()
	resources = loadResourceData()

	logger.Info().Msg("Preloading finished!")
	return items, maps, monsters, resources
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

func loadResourceData() []resources.Resource {
	logger.Info().Msg("Loading ALL resource data")
	resourceData, err := resources.GetAllResources()
	if err != nil {
		logger.Error().Err(err).Msg("Error loading resource data")
		os.Exit(1)
	}

	logger.Info().Msgf("Loaded %d resources", len(resourceData))
	return resourceData
}

func waitForCooldowns() {
	//TODO: we just assume 58 seconds for now, later we need to make this smarter by preloading character data and waiting in each thread
	cooldown := 58
	logger.Info().Msgf("Waiting for %d seconds before starting characters for all cool downs to finish", cooldown)
	time.Sleep(time.Duration(cooldown) * time.Second)
}
