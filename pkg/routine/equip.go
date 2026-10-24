package routine

import (
	"fmt"
	"strings"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
)

// EquipOptions configures movement and utility handling.
type EquipOptions struct {
	Movement MoveOptions
	// Utilities gives desired values; only changed slots are reset.
	Utilities map[string]string
	// ClearUtilities retains Equip's legacy behavior.
	ClearUtilities bool
}

func Equip(name string, equipments []schemas.EquipSchema, options EquipOptions) (schemas.CharacterSchema, error) {
	// +gocover:ignore:block production wrapper over tested implementation
	return equip(defaultDeps, name, equipments, options)
}

func equip(d deps, name string, equipments []schemas.EquipSchema, options EquipOptions) (schemas.CharacterSchema, error) {
	err := validateEquipments(equipments, options.Utilities)
	if err != nil {
		return schemas.CharacterSchema{}, err
	}
	character, err := d.characters(name)
	if err != nil {
		return schemas.CharacterSchema{}, err
	}
	normalizeArtifactSlots(character, equipments)
	normalizeRingPermutation(character, equipments)
	err = checkLevelRequirements(character, equipments)
	if err != nil {
		return schemas.CharacterSchema{}, err
	}
	needed := filterNeededEquipments(character, equipments)
	hasEquipmentChanges := len(needed) > 0
	utilitySlots := changedUtilitySlots(character, options.Utilities)
	if options.ClearUtilities && (character.Utility1Slot != "" || character.Utility2Slot != "") {
		utilitySlots = []string{"utility1", "utility2"}
	}
	hasUtilityChanges := len(utilitySlots) > 0
	if !hasEquipmentChanges && !hasUtilityChanges {
		return character, nil
	}
	if hasEquipmentChanges {
		err = validateTotalStock(d, character, needed)
		if err != nil {
			return schemas.CharacterSchema{}, err
		}
	}

	Cooldown(character)
	if hasEquipmentChanges {
		hpLoss := equipmentHPLoss(character, needed)
		if hpLoss > 0 && hpLoss >= character.Hp {
			character, err = hp(d, character, HpOptions{
				MinHP:   hpLoss + 1,
				UseFood: true,
			})
			if err != nil {
				return character, err
			}
		}
	}

	missing := calculateMissingItems(character, needed)
	shouldPrepare := hasUtilityChanges || len(missing) > 0
	if shouldPrepare {
		character, err = deposit(d, character, DepositOptions{
			Movement: options.Movement,
		})
		if err != nil {
			return character, err
		}
		if hasUtilityChanges {
			character, err = clearUtilities(d, character, utilitySlots)
			if err != nil {
				return character, err
			}
		}
		if hasEquipmentChanges {
			withdraw := equipmentItems(needed)
			transaction, err := d.myActionBankWithdrawItem(name, withdraw)
			if err != nil {
				return character, err
			}
			character = transaction.Character
		}
	}
	if hasEquipmentChanges {
		result, err := d.myActionEquip(name, needed)
		if err != nil {
			return character, err
		}
		character = result.Character
	}
	if shouldPrepare {
		character, err = deposit(d, character, DepositOptions{
			Movement: options.Movement,
		})
		if err != nil {
			return character, err
		}
	}
	return character, nil
}

func changedUtilitySlots(character schemas.CharacterSchema, changes map[string]string) []string {
	var slots []string
	for _, slot := range []string{"utility1", "utility2"} {
		desired, changed := changes[slot]
		if !changed {
			continue
		}
		current := character.Utility1Slot
		if slot == "utility2" {
			current = character.Utility2Slot
		}
		if desired != current {
			slots = append(slots, slot)
		}
	}
	return slots
}

func equipmentHPLoss(character schemas.CharacterSchema, needed []schemas.EquipSchema) int {
	current := map[string]string{
		"amulet":     character.AmuletSlot,
		"artifact1":  character.Artifact1Slot,
		"artifact2":  character.Artifact2Slot,
		"artifact3":  character.Artifact3Slot,
		"bag":        character.BagSlot,
		"body_armor": character.BodyArmorSlot,
		"boots":      character.BootsSlot,
		"helmet":     character.HelmetSlot,
		"leg_armor":  character.LegArmorSlot,
		"ring1":      character.Ring1Slot,
		"ring2":      character.Ring2Slot,
		"rune":       character.RuneSlot,
		"shield":     character.ShieldSlot,
		"weapon":     character.WeaponSlot,
	}
	loss := 0
	for _, equipment := range needed {
		if strings.HasPrefix(equipment.Slot, "utility") {
			continue
		}
		oldItem, exists := catalog.Items().Get(current[equipment.Slot])
		if !exists {
			continue
		}
		loss += itemHPEffect(oldItem)
	}
	return loss
}

func itemHPEffect(item *schemas.ItemSchema) int {
	if item.Effects == nil {
		return 0
	}
	for _, effect := range *item.Effects {
		if effect.Code == "hp" {
			return effect.Value
		}
	}
	return 0
}

func equipmentItems(equipments []schemas.EquipSchema) []schemas.SimpleItemSchema {
	items := make([]schemas.SimpleItemSchema, 0, len(equipments))
	for _, equipment := range equipments {
		quantity := 1
		if equipment.Quantity != nil {
			quantity = *equipment.Quantity
		}
		items = append(items, schemas.SimpleItemSchema{
			Code:     equipment.Code,
			Quantity: quantity,
		})
	}
	return items
}

func validateEquipments(equipments []schemas.EquipSchema, utilityChanges map[string]string) error {
	slotsReserved := map[string]string{}
	artifacts := map[string]bool{}
	utilities := map[string]bool{}
	for _, slot := range []string{"utility1", "utility2"} {
		code := utilityChanges[slot]
		if code == "" {
			continue
		}
		if utilities[code] {
			return fmt.Errorf("utility cannot be equipped in multiple slots: %s", code)
		}
		utilities[code] = true
	}
	for _, equipment := range equipments {
		if equipment.Code == "" || equipment.Slot == "" {
			return fmt.Errorf("invalid equipment request: missing code or slot")
		}
		item, exists := catalog.Items().Get(equipment.Code)
		if !exists {
			return fmt.Errorf("item not found in database: %s", equipment.Code)
		}
		slots, hasSlot := catalog.EquipmentTypeToSlots[item.Type]
		if !hasSlot {
			return fmt.Errorf("item type cannot be equipped: %s", item.Type)
		}
		if item.Type == "artifact" {
			if artifacts[equipment.Code] {
				return fmt.Errorf("artifact cannot be equipped in multiple slots: %s", equipment.Code)
			}
			artifacts[equipment.Code] = true
		}
		if item.Type == "utility" {
			if utilities[equipment.Code] {
				return fmt.Errorf("utility cannot be equipped in multiple slots: %s", equipment.Code)
			}
			utilities[equipment.Code] = true
		}
		validSlot := false
		for _, s := range slots {
			if s == equipment.Slot {
				validSlot = true
				break
			}
		}
		if !validSlot {
			return fmt.Errorf("item %s cannot be equipped in slot %s", equipment.Code, equipment.Slot)
		}
		if equipment.Quantity != nil && *equipment.Quantity > 1 && !strings.HasPrefix(equipment.Slot, "utility") {
			return fmt.Errorf("cannot specify quantity for non-utility slot: %s", equipment.Slot)
		}
		conflictingItem, reserved := slotsReserved[equipment.Slot]
		if reserved {
			return fmt.Errorf("slot conflict: both %s and %s are targeting %s", conflictingItem, equipment.Code, equipment.Slot)
		}
		slotsReserved[equipment.Slot] = equipment.Code
	}
	return nil
}

func normalizeArtifactSlots(character schemas.CharacterSchema, equipments []schemas.EquipSchema) {
	current := map[string]string{
		"artifact1": character.Artifact1Slot,
		"artifact2": character.Artifact2Slot,
		"artifact3": character.Artifact3Slot,
	}
	usedSlots := map[string]bool{}
	assigned := map[int]bool{}

	for index, equipment := range equipments {
		if !isArtifactSlot(equipment.Slot) {
			continue
		}
		for slot, code := range current {
			if code == equipment.Code && !usedSlots[slot] {
				equipments[index].Slot = slot
				usedSlots[slot] = true
				assigned[index] = true
				break
			}
		}
	}

	for index, equipment := range equipments {
		if !isArtifactSlot(equipment.Slot) || assigned[index] {
			continue
		}
		slot := equipment.Slot
		if usedSlots[slot] {
			slot = nextFreeArtifactSlot(usedSlots)
		}
		equipments[index].Slot = slot
		usedSlots[slot] = true
	}
}

func isArtifactSlot(slot string) bool {
	return slot == "artifact1" || slot == "artifact2" || slot == "artifact3"
}

func nextFreeArtifactSlot(usedSlots map[string]bool) string {
	for _, slot := range []string{"artifact1", "artifact2", "artifact3"} {
		if !usedSlots[slot] {
			return slot
		}
	}
	return ""
}

func normalizeRingPermutation(character schemas.CharacterSchema, equipments []schemas.EquipSchema) {
	if character.Ring1Slot == "" || character.Ring2Slot == "" {
		return
	}
	var ring1Index, ring2Index = -1, -1
	for index, equipment := range equipments {
		switch equipment.Slot {
		case "ring1":
			ring1Index = index
		case "ring2":
			ring2Index = index
		}
	}
	if ring1Index < 0 || ring2Index < 0 || character.Ring1Slot == character.Ring2Slot {
		return
	}
	if equipments[ring1Index].Code != character.Ring2Slot || equipments[ring2Index].Code != character.Ring1Slot {
		return
	}
	equipments[ring1Index].Slot = "ring2"
	equipments[ring2Index].Slot = "ring1"
}

func checkLevelRequirements(character schemas.CharacterSchema, equipments []schemas.EquipSchema) error {
	for _, equipment := range equipments {
		item, _ := catalog.Items().Get(equipment.Code)
		if !utils.MeetsItemConditions(character, *item) {
			return fmt.Errorf("does not meet requirement for %s", item.Name)
		}
	}
	return nil
}

func filterNeededEquipments(character schemas.CharacterSchema, equipments []schemas.EquipSchema) []schemas.EquipSchema {
	currentSlots := map[string]string{
		"amulet":     character.AmuletSlot,
		"artifact1":  character.Artifact1Slot,
		"artifact2":  character.Artifact2Slot,
		"artifact3":  character.Artifact3Slot,
		"bag":        character.BagSlot,
		"body_armor": character.BodyArmorSlot,
		"boots":      character.BootsSlot,
		"helmet":     character.HelmetSlot,
		"leg_armor":  character.LegArmorSlot,
		"ring1":      character.Ring1Slot,
		"ring2":      character.Ring2Slot,
		"rune":       character.RuneSlot,
		"shield":     character.ShieldSlot,
		"utility1":   character.Utility1Slot,
		"utility2":   character.Utility2Slot,
		"weapon":     character.WeaponSlot,
	}
	currentQuantities := map[string]int{
		"utility1": character.Utility1SlotQuantity,
		"utility2": character.Utility2SlotQuantity,
	}

	var needed []schemas.EquipSchema
	for _, equipment := range equipments {
		currentQty := 1
		currentItem := currentSlots[equipment.Slot]
		q, ok := currentQuantities[equipment.Slot]
		if ok {
			currentQty = q
		}
		targetQty := 1
		if equipment.Quantity != nil {
			targetQty = *equipment.Quantity
		}
		if currentItem == equipment.Code && currentQty == targetQty {
			continue
		}
		needed = append(needed, equipment)
	}
	return needed
}

func validateTotalStock(d deps, character schemas.CharacterSchema, needed []schemas.EquipSchema) error {
	bankItems, err := d.myBankItems()
	if err != nil {
		return err
	}
	totalRequired := map[string]int{}
	for _, equipment := range needed {
		targetQty := 1
		if equipment.Quantity != nil {
			targetQty = *equipment.Quantity
		}
		totalRequired[equipment.Code] += targetQty
	}
	for code, requiredQty := range totalRequired {
		totalAvailable := 0
		if character.Inventory != nil {
			for _, invSlot := range *character.Inventory {
				if invSlot.Code == code {
					totalAvailable += invSlot.Quantity
				}
			}
		}
		for _, bankItem := range bankItems {
			if bankItem.Code == code {
				totalAvailable += bankItem.Quantity
			}
		}
		if totalAvailable < requiredQty {
			return fmt.Errorf("insufficient items: %s requires %d, but you only have %d", code, requiredQty, totalAvailable)
		}
	}
	return nil
}

func calculateMissingItems(character schemas.CharacterSchema, needed []schemas.EquipSchema) []schemas.SimpleItemSchema {
	totalRequired := map[string]int{}
	for _, equipment := range needed {
		targetQty := 1
		if equipment.Quantity != nil {
			targetQty = *equipment.Quantity
		}
		totalRequired[equipment.Code] += targetQty
	}
	var missing []schemas.SimpleItemSchema
	for code, requiredQty := range totalRequired {
		qtyInInventory := 0
		if character.Inventory != nil {
			for _, invSlot := range *character.Inventory {
				if invSlot.Code == code {
					qtyInInventory = invSlot.Quantity
					break
				}
			}
		}
		if qtyInInventory < requiredQty {
			missing = append(missing, schemas.SimpleItemSchema{
				Code:     code,
				Quantity: requiredQty - qtyInInventory,
			})
		}
	}
	return missing
}
