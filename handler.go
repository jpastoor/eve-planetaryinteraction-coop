package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math"
)

type Handler struct {
	inv           *Inventory
	ledger        *Ledger
	playerFetcher PlayerFetcher
	typeFetcher   TypeFetcher
}

/**
TODO This method does not correctly take into account the order in which transactions are done, so the output can vary with the same input.
     Think this comes due to the unsortedness of the data map.
 */
func (h *Handler) Process(ts []Transaction) (invCreditMuts []InventoryMutation, invDebitMuts []InventoryMutation, err error) {

	data := map[int]map[string]int{}
	for _, t := range ts {

		// First we update the playerName of each transaction
		who, err := h.playerFetcher.getOrCreatePlayerByName(t.PlayerName)
		if err != nil {
			return nil, nil, err
		}

		if who.Main != "" {
			t.PlayerName = who.Main
		}

		if t.MarkedForCorp {
			t.PlayerName = "ADHC"
		}

		// Then we make a sum of all the item locks and unlocks per player
		transType, err := h.typeFetcher.getTypeByName(t.TypeName)
		if err != nil {
			return nil, nil, err
		}

		if _, exists := data[transType.TypeID]; !exists {
			data[transType.TypeID] = map[string]int{}
		}

		if _, exists := data[transType.TypeID][t.PlayerName]; !exists {
			data[transType.TypeID][t.PlayerName] = 0
		}

		if t.Action == ACTION_LOCK {
			data[transType.TypeID][t.PlayerName] += t.Quantity
		}

		if t.Action == ACTION_UNLOCK {
			data[transType.TypeID][t.PlayerName] -= t.Quantity
		}
	}

	var creditMutations []InventoryMutation
	var debitMutations []InventoryMutation

	// We then process all the positive amounts
	for typeId, playerMap := range data {
		for playerName, amount := range playerMap {
			if amount > 0 {
				invMut := InventoryMutation{
					PlayerName: playerName,
					TypeId:     typeId,
					Change:     amount,
				}

				h.inv.Add(invMut)
			}
		}
	}

	// Then we process all the negative amounts
	for typeId, playerMap := range data {
		for playerName, amount := range playerMap {
			if amount < 0 {
				invMut := InventoryMutation{
					PlayerName: playerName,
					TypeId:     typeId,
					Change:     amount,
				}

				nCreditMuts, nDebitMuts := h.inv.Sub(invMut)
				creditMutations = append(creditMutations, nCreditMuts...)
				debitMutations = append(debitMutations, nDebitMuts...)
			}
		}
	}

	return creditMutations, debitMutations, nil
}

func (h *Handler) CalculateSchematicProfit() ([]SchematicProfit, error) {

	output := []SchematicProfit{}

	invFlat := h.inv.getInventoryFlat()

	planetSchematicsBytes, err := ioutil.ReadFile("./sde/planetSchematics.yaml")
	if err != nil {
		return nil, err
	}

	var schematics []PlanetSchematic
	if err := yaml.Unmarshal(planetSchematicsBytes, &schematics); err != nil {
		return nil, err
	}

	planetSchematicsTypeMapBytes, err := ioutil.ReadFile("./sde/planetSchematicsTypeMap.yaml")
	if err != nil {
		return nil, err
	}

	var schematicTypes []PlanetSchematicTypeMap
	if err := yaml.Unmarshal(planetSchematicsTypeMapBytes, &schematicTypes); err != nil {
		return nil, err
	}

	for _, schema := range schematics {

		inputMats := []PlanetSchematicTypeMap{}
		var outputMat PlanetSchematicTypeMap

		var typeIDs []int
		for _, schemaMat := range schematicTypes {
			if schemaMat.SchematicID == schema.SchematicID {
				if !schemaMat.IsInput {
					outputMat = schemaMat
				} else {
					inputMats = append(inputMats, schemaMat)
				}

				typeIDs = append(typeIDs, schemaMat.TypeID)
			}
		}

		// Check how much cycles we can make
		var maxCycles float64
		maxCycles = math.MaxFloat64
		for _, inputMat := range inputMats {
			amountOfMat, _ := invFlat[inputMat.TypeID]
			cycles := float64(amountOfMat) / float64(inputMat.Quantity)
			maxCycles = math.Min(maxCycles, cycles)
		}

		if maxCycles > 0 {
			typePrices, _ := h.ledger.priceApi.FetchPrices(typeIDs)

			// Calculate input volume and isk
			var inputVolume float64
			var inputISK float64
			var inputMatsDisplay []SchematicInput
			for _, inputMat := range inputMats {
				inputType, _ := h.typeFetcher.getTypeById(inputMat.TypeID)
				inputISK += float64(typePrices[inputMat.TypeID]) * float64(inputMat.Quantity) * maxCycles
				inputVolume += float64(inputType.Volume) * maxCycles * float64(inputMat.Quantity)
				inputMatsDisplay = append(inputMatsDisplay, SchematicInput{
					TypeName: inputType.TypeName,
					Quantity: int(maxCycles) * inputMat.Quantity,
				})
			}

			outputType, _ := h.typeFetcher.getTypeById(outputMat.TypeID)
			outputVolume := float64(outputType.Volume) * maxCycles * float64(outputMat.Quantity)
			outputISK := float64(typePrices[outputMat.TypeID]) * maxCycles * float64(outputMat.Quantity)

			output = append(output, SchematicProfit{
				OutputTypeName:  outputType.TypeName,
				Cycles:          int(maxCycles),
				VolumeReduction: math.Round(inputVolume - outputVolume),
				ISK:             math.Round(outputISK - inputISK),
				TotalTime:       int(maxCycles) * schema.CycleTime,
				InputMats:       inputMatsDisplay,
			})
		}
	}

	return output, nil
}

type SchematicProfit struct {
	OutputTypeName  string
	ISK             float64
	VolumeReduction float64
	Cycles          int
	TotalTime       int
	InputMats       []SchematicInput
}

type SchematicInput struct {
	TypeName string
	Quantity int
}
