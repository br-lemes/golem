package stock

import (
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
)

type ItemStock struct {
	Current int `json:"current"`
	Needed  int `json:"needed"`
	Demand  int `json:"demand"`
	Safety  int `json:"safety"`
	Target  int `json:"target"`
}

type goal struct {
	Quantity int
	Expand   bool
}

func Overview(stocks []models.Stock, codes []string, target bool) (map[string]ItemStock, error) {
	characters, err := api.AccountsCharacters("")
	if err != nil {
		return nil, err
	}
	bankItems, err := api.MyBankItems()
	if err != nil {
		return nil, err
	}
	bank, err := api.MyBank()
	if err != nil {
		return nil, err
	}
	quantities := itemQuantities(bankItems, characters)
	quantities["gold"] = goldQuantity(bank.Gold, characters)
	return Calculate(stocks, quantities, characters, codes, target)
}

func Calculate(stocks []models.Stock, quantities map[string]int, characters []schemas.CharacterSchema, codes []string, target bool) (map[string]ItemStock, error) {
	configured := configured(stocks)
	reachable := taskSafeties(catalog.Tasks().All(), nil, maxSkillLevels(characters))
	if err := initialize(configured, potionSafeties()); err != nil {
		return nil, err
	}
	if err := initialize(configured, allTaskSafeties(catalog.Tasks().All())); err != nil {
		return nil, err
	}
	if err := updateReplenishing(configured, quantities); err != nil {
		return nil, err
	}
	return requirements(configured, quantities, characters, taskCodes(), reachable, codes, target), nil
}

func requirements(stocks map[string]models.Stock, quantities map[string]int, characters []schemas.CharacterSchema, taskCodes map[string]bool, reachable map[string]int, codes []string, target bool) map[string]ItemStock {
	direct := directRequirements(stocks, quantities, characters, taskCodes, reachable, target)
	required, needed := recipeRequirements(catalogItems(), npcOffers(), direct, quantities, characters)
	result := map[string]ItemStock{}
	for code, stock := range stocks {
		if !stock.Enabled || stock.Mode == "" || (len(codes) > 0 && !slices.Contains(codes, code)) {
			continue
		}
		_, reachableTask := reachable[code]
		current := quantities[code]
		if current == 0 && !itemAvailable(code, characters) {
			continue
		}
		if taskCodes[code] && !reachableTask && current == 0 {
			continue
		}
		item := itemStock(stock, current)
		item.Demand = demand(stock, required[code])
		item.Needed = needed[code]
		if !reachableTask && taskCodes[code] && current > 0 {
			result[code] = item
			continue
		}
		if item.Needed == 0 {
			continue
		}
		result[code] = item
	}
	for code, quantity := range needed {
		if quantity == 0 || (len(codes) > 0 && !slices.Contains(codes, code)) {
			continue
		}
		stock, configured := stocks[code]
		if configured && !stock.Enabled {
			continue
		}
		item := ItemStock{
			Current: quantities[code],
			Needed:  quantity,
			Demand:  required[code],
		}
		if configured && stock.Mode != "" {
			item = itemStock(stock, quantities[code])
			item.Needed = quantity
			item.Demand = demand(stock, required[code])
		}
		result[code] = item
	}
	return result
}

func directRequirements(stocks map[string]models.Stock, quantities map[string]int, characters []schemas.CharacterSchema, taskCodes map[string]bool, reachable map[string]int, target bool) map[string]goal {
	direct := map[string]goal{}
	for code, stock := range stocks {
		if !stock.Enabled || stock.Mode == "" {
			continue
		}
		current := quantities[code]
		if current == 0 && !itemAvailable(code, characters) {
			continue
		}
		_, reachableTask := reachable[code]
		if taskCodes[code] && !reachableTask && current == 0 {
			continue
		}
		if stock.Mode == models.StockModeExact {
			direct[code] = goal{
				Quantity: stock.Quantity,
				Expand:   current < stock.Quantity,
			}
			continue
		}
		threshold := stock.Quantity
		desired := threshold
		if stock.Replenishing {
			desired += threshold / 3
		}
		if target {
			threshold = stock.Quantity + stock.Quantity/3
			desired = threshold
		}
		direct[code] = goal{Quantity: desired, Expand: current < desired}
	}
	return direct
}

func recipeRequirements(items map[string]*schemas.ItemSchema, offers map[string][]schemas.SimpleNPCItemSchema, direct map[string]goal, quantities map[string]int, characters []schemas.CharacterSchema) (map[string]int, map[string]int) {
	available := map[string]int{}
	for code, quantity := range quantities {
		available[code] = quantity
	}
	required := map[string]int{}
	needed := map[string]int{}
	missing := map[string]int{}
	for code, goal := range direct {
		required[code] += goal.Quantity
		used := min(goal.Quantity, available[code])
		available[code] -= used
		if goal.Expand && goal.Quantity > used {
			missing[code] = goal.Quantity - used
			needed[code] += goal.Quantity - used
		}
	}
	for code, quantity := range missing {
		expandRecipe(items, offers, code, quantity, available, required, needed, characters, map[string]bool{})
	}
	return required, needed
}

func expandRecipe(items map[string]*schemas.ItemSchema, offers map[string][]schemas.SimpleNPCItemSchema, code string, quantity int, available, required, needed map[string]int, characters []schemas.CharacterSchema, visiting map[string]bool) {
	item, exists := items[code]
	if visiting[code] {
		return
	}
	if !exists || !craftable(*item, characters) {
		expandNPCOffer(items, offers, code, quantity, available, required, needed, characters, visiting)
		return
	}
	craft := *item.Craft
	if craft.Items == nil {
		return
	}
	output := 1
	if craft.Quantity != nil {
		output = *craft.Quantity
	}
	if output <= 0 {
		return
	}
	visiting[code] = true
	batches := (quantity + output - 1) / output
	for _, ingredient := range *craft.Items {
		demand := ingredient.Quantity * batches
		required[ingredient.Code] += demand
		used := min(demand, available[ingredient.Code])
		available[ingredient.Code] -= used
		if demand > used {
			missing := demand - used
			needed[ingredient.Code] += missing
			expandRecipe(items, offers, ingredient.Code, missing, available, required, needed, characters, visiting)
		}
	}
	available[code] += batches*output - quantity
	delete(visiting, code)
}

func expandNPCOffer(items map[string]*schemas.ItemSchema, offers map[string][]schemas.SimpleNPCItemSchema, code string, quantity int, available, required, needed map[string]int, characters []schemas.CharacterSchema, visiting map[string]bool) {
	if len(offers[code]) == 0 {
		return
	}
	visiting[code] = true
	defer delete(visiting, code)
	item := offers[code][0]
	if item.Currency == "gold" && !npcGoldRequired(items, code) {
		return
	}
	demand := *item.BuyPrice * quantity
	required[item.Currency] += demand
	used := min(demand, available[item.Currency])
	available[item.Currency] -= used
	if demand == used {
		return
	}
	missing := demand - used
	needed[item.Currency] += missing
	expandRecipe(items, offers, item.Currency, missing, available, required, needed, characters, visiting)
}

func updateReplenishing(stocks map[string]models.Stock, quantities map[string]int) error {
	for code, stock := range stocks {
		if !stock.Enabled || stock.Mode != models.StockModeSafety {
			continue
		}
		target := stock.Quantity + stock.Quantity/3
		replenishing := stock.Replenishing
		if replenishing && quantities[code] >= target {
			replenishing = false
		} else if !replenishing && quantities[code] < stock.Quantity {
			replenishing = true
		}
		if replenishing == stock.Replenishing {
			continue
		}
		err := cache.SetStockReplenishing(code, replenishing)
		if err != nil {
			return err
		}
		stock.Replenishing = replenishing
		stocks[code] = stock
	}
	return nil
}

func initialize(stocks map[string]models.Stock, safeties map[string]int) error {
	for code, safety := range safeties {
		stock, exists := stocks[code]
		if exists && (!stock.Enabled || stock.Mode != "") {
			continue
		}
		err := cache.SetStock(code, safety, models.StockModeSafety)
		if err != nil {
			return err
		}
		stocks[code] = models.Stock{
			Code:     code,
			Quantity: safety,
			Mode:     models.StockModeSafety,
			Enabled:  true,
		}
	}
	return nil
}

func taskSafeties(tasks []*schemas.TaskFullSchema, codes []string, maxSkillLevel map[string]int) map[string]int {
	safeties := map[string]int{}
	for _, task := range tasks {
		if task.Skill == nil || task.Type != "items" || task.Level > maxSkillLevel[*task.Skill] || (len(codes) > 0 && !slices.Contains(codes, task.Code)) {
			continue
		}
		safeties[task.Code] += taskSafety(task)
	}
	return safeties
}

func allTaskSafeties(tasks []*schemas.TaskFullSchema) map[string]int {
	safeties := map[string]int{}
	for _, task := range tasks {
		if task.Skill == nil || task.Type != "items" {
			continue
		}
		safeties[task.Code] += taskSafety(task)
	}
	return safeties
}

func taskSafety(task *schemas.TaskFullSchema) int {
	safety := task.MaxQuantity * 10
	if task.MaxQuantity > 100 {
		safety = task.MaxQuantity * 5
	}
	return safety
}

func potionSafeties() map[string]int {
	safeties := map[string]int{}
	for _, potion := range catalog.Items().Potions().All() {
		safeties[potion.Code] = 500
	}
	return safeties
}

func taskCodes() map[string]bool {
	codes := map[string]bool{}
	potions := map[string]bool{}
	for _, potion := range catalog.Items().Potions().All() {
		potions[potion.Code] = true
	}
	for _, code := range catalog.Tasks().Items() {
		if potions[code] {
			continue
		}
		codes[code] = true
	}
	return codes
}

func maxSkillLevels(characters []schemas.CharacterSchema) map[string]int {
	maxSkillLevel := map[string]int{}
	for _, character := range characters {
		for _, skill := range catalog.Tasks().Skills() {
			level, _ := utils.GetCharacterSkillLevel(character, skill)
			maxSkillLevel[skill] = max(maxSkillLevel[skill], level)
		}
	}
	return maxSkillLevel
}

func itemAvailable(code string, characters []schemas.CharacterSchema) bool {
	item, ok := catalog.Items().Get(code)
	if !ok {
		return true
	}
	return itemUsable(*item, characters)
}

func itemUsable(item schemas.ItemSchema, characters []schemas.CharacterSchema) bool {
	for _, character := range characters {
		if utils.MeetsItemConditions(character, item) {
			return true
		}
	}
	return false
}

func craftable(item schemas.ItemSchema, characters []schemas.CharacterSchema) bool {
	if item.Craft == nil || item.Craft.Items == nil {
		return false
	}
	for _, character := range characters {
		if !utils.MeetsItemConditions(character, item) {
			continue
		}
		if item.Craft.Skill == nil || item.Craft.Level == nil {
			return true
		}
		level, ok := utils.GetCharacterSkillLevel(character, *item.Craft.Skill)
		if ok && level >= *item.Craft.Level {
			return true
		}
	}
	return false
}

func npcGoldRequired(items map[string]*schemas.ItemSchema, code string) bool {
	item := items[code]
	if item != nil && item.Craft != nil {
		return false
	}
	for _, resource := range catalog.Resources.All() {
		for _, drop := range resource.Drops {
			if drop.Code == code {
				return false
			}
		}
	}
	for _, monster := range catalog.Monsters.All() {
		for _, drop := range monster.Drops {
			if drop.Code == code {
				return false
			}
		}
	}
	return true
}

func configured(items []models.Stock) map[string]models.Stock {
	stocks := make(map[string]models.Stock, len(items))
	for _, stock := range items {
		stocks[stock.Code] = stock
	}
	return stocks
}

func demand(stock models.Stock, required int) int {
	goal := stock.Quantity
	if stock.Mode == models.StockModeSafety {
		goal += goal / 3
	}
	return max(0, required-goal)
}

func itemStock(stock models.Stock, current int) ItemStock {
	item := ItemStock{Current: current}
	if stock.Mode == models.StockModeExact {
		item.Target = stock.Quantity
		if current < stock.Quantity {
			item.Needed = stock.Quantity - current
		}
		return item
	}
	target := stock.Quantity + stock.Quantity/3
	item.Safety = stock.Quantity
	item.Target = target
	if current < target {
		item.Needed = target - current
	}
	return item
}

func catalogItems() map[string]*schemas.ItemSchema {
	items := map[string]*schemas.ItemSchema{}
	for _, item := range catalog.Items().All() {
		items[item.Code] = item
	}
	return items
}

func npcOffers() map[string][]schemas.SimpleNPCItemSchema {
	offers := map[string][]schemas.SimpleNPCItemSchema{}
	for _, npc := range catalog.NpcsDetails.All() {
		if npc.Items == nil {
			continue
		}
		for _, item := range *npc.Items {
			if item.BuyPrice != nil && *item.BuyPrice > 0 {
				offers[item.Code] = append(offers[item.Code], item)
			}
		}
	}
	return offers
}
