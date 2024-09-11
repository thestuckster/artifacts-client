package main

import (
	"artifacts-client/internal"
	"artifacts-client/internal/devHelpers"
	"artifacts-client/internal/planner"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/thestuckster/gopherfacts/pkg/clients"
	"github.com/thestuckster/gopherfacts/pkg/items"
	"os"
)

var logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

func main() {
	logger.Info().Msg("Starting artifacts client")

	token := os.Getenv("TOKEN")
	sdk := clients.NewClient(&token)

	checkArtifactsServerStatus(sdk)

	itemData := loadItemData()
	logger.Debug().Msg("Parsing item information into hashmaps for faster lookups")
	itemsByCode := internal.BuildItemByCode(itemData)
	_ = internal.BuildItemBySkill(itemData)

	//Throw away testing code
	mostRequiredMaterials := devHelpers.FindItemsWithMostUniqueCraftingRequirements(itemData)
	fmt.Printf("Item with most required materials: %+v\n", mostRequiredMaterials)

	//tree := crafting.NewCraftingTree("multislimes_sword", itemsByCode)
	planner.BuildPlanForItem("multislimes_sword", itemsByCode)

	//fmt.Printf("%+v", tree)
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
