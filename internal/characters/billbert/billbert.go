package billbert

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/thestuckster/gopherfacts/pkg/clients"
	"github.com/thestuckster/gopherfacts/pkg/exchange"
)

const name = "Billbert"
const copperBarMin = 8 //it takes 8 ores to make 1 bar

func GameLoop(sdk *clients.GopherFactClient, wg *sync.WaitGroup) {
	defer wg.Done()

	logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Str("character", name).Logger()

	for {
		err := mineCopper(sdk, &logger)
		if err != nil {
			logger.Error().Err(err).Msg("Something went wrong, stopping loop")
			break
		}
	}
}

func mineCopper(sdk *clients.GopherFactClient, logger *zerolog.Logger) error {
	res, err := sdk.EasyClient.MineCopper(name)

	if err != nil {
		var ex *clients.CharacterInventoryFullException
		if errors.As(err, &ex) {
			logger.Warn().Msg("Inventory is full, cannot mine more copper")
			smeltAndSellCopper(sdk, logger)
		} else {
			logger.Error().Err(err).Msg("Failed to mine copper")
			return err
		}
	}

	msgFormat := fmt.Sprintf("Mined copper! got %d total items: ", len(res.Details.Items))
	for _, item := range res.Details.Items {
		msgFormat += fmt.Sprintf("%s: %d, ", item.Code, item.Quantity)
	}

	return nil
}

func smeltAndSellCopper(sdk *clients.GopherFactClient, logger *zerolog.Logger) error {
	logger.Info().Msg("Smelting copper")
	_, err := sdk.EasyClient.MoveToForge(name)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to move to forge")
		return err
	}

	inv, err := getCharacterInventory(sdk, logger)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get inventory")
		return err
	}

	if inv["copper_ore"] < copperBarMin {
		logger.Warn().Msgf("Not enough copper ore to smelt, need %d, have %d", copperBarMin, inv["copper_ore"])
		return errors.New("not enough copper ore to smelt")
	}

	barsToProduce := inv["copper_ore"] / copperBarMin
	logger.Debug().Msgf("Crafting %d copper bars", barsToProduce)

	craftRes, err := sdk.EasyClient.Craft(name, "copper", barsToProduce)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to smelt copper")
		return err
	}

	logger.Info().Msgf("Smelted copper! got %d total items: ", len(craftRes.Details.Items))
	for _, item := range craftRes.Details.Items {
		logger.Info().Msgf("%s: %d, ", item.Code, item.Quantity)
	}

	sellItem(sdk, logger, "copper", barsToProduce)

	return nil
}

func getCharacterInventory(sdk *clients.GopherFactClient, logger *zerolog.Logger) (map[string]int, error) {
	charInfo, err := sdk.CharacterClient.GetCharacterInfo(name)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get inventory")
		return nil, err
	}

	itemMap := make(map[string]int)
	for _, item := range charInfo.Inventory {
		itemMap[item.Code] += item.Quantity
	}

	return itemMap, nil
}

func sellItem(sdk *clients.GopherFactClient, logger *zerolog.Logger, itemCode string, quantity int) error {
	logger.Info().Msgf("Selling %d %s", quantity, itemCode)

	_, err := sdk.EasyClient.MoveToExchange(name)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to move to exchange")
		return err
	}

	currentQuantity := quantity
	maxSellAtOnce := 25

	for currentQuantity > 0 {
		currentPrice, err := getCurrentSellPrice(sdk, logger, itemCode)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get current sell price")
			return err
		}

		amountToSell := maxSellAtOnce
		if currentQuantity < maxSellAtOnce {
			amountToSell = currentQuantity
		}

		sellRes, err := sdk.EasyClient.SellItem(name, itemCode, amountToSell, currentPrice)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to sell item")
			return err
		}

		logger.Info().Msgf("Sold %d %s and earned %d gold", amountToSell, itemCode, sellRes.Transaction.TotalPrice)
		currentQuantity -= amountToSell
	}
	return nil
}

func getCurrentSellPrice(sdk *clients.GopherFactClient, logger *zerolog.Logger, itemCode string) (int, error) {
	exchangeData, err := exchange.GetItemExchangeData(itemCode)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get exchange data")
		return 0, err
	}

	return exchangeData.SellPrice, nil
}
