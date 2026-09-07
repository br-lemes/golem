package schemas

import (
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type AccessSchema struct {
	Conditions *[]ConditionSchema `json:"conditions,omitempty"`
	Type       string             `json:"type"`
}

type AccountAchievementObjectiveSchema struct {
	Progress *int    `json:"progress,omitempty"`
	Target   *string `json:"target,omitempty"`
	Total    int     `json:"total"`
	Type     string  `json:"type"`
}

type AccountAchievementSchema struct {
	Code        string                              `json:"code"`
	CompletedAt *time.Time                          `json:"completed_at,omitempty"`
	Description string                              `json:"description"`
	Name        string                              `json:"name"`
	Objectives  []AccountAchievementObjectiveSchema `json:"objectives"`
	Points      int                                 `json:"points"`
	Rewards     AchievementRewardsSchema            `json:"rewards"`
}

type AccountDetails struct {
	AchievementsPoints int       `json:"achievements_points"`
	Badges             *[]string `json:"badges,omitempty"`
	BanReason          *string   `json:"ban_reason,omitempty"`
	Banned             bool      `json:"banned"`
	Member             bool      `json:"member"`
	Skins              []string  `json:"skins"`
	Status             string    `json:"status"`
	Username           string    `json:"username"`
}

type AccountDetailsSchema struct {
	Data AccountDetails `json:"data"`
}

type AccountLeaderboardSchema struct {
	Account            string     `json:"account"`
	AchievementsPoints int        `json:"achievements_points"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	Gold               int        `json:"gold"`
	Member             bool       `json:"member"`
	Position           int        `json:"position"`
}

type AchievementObjectiveSchema struct {
	Target *string `json:"target,omitempty"`
	Total  int     `json:"total"`
	Type   string  `json:"type"`
}

type AchievementResponseSchema struct {
	Data AchievementSchema `json:"data"`
}

type AchievementRewardsSchema struct {
	Gold  *int                `json:"gold,omitempty"`
	Items *[]RewardItemSchema `json:"items,omitempty"`
}

type AchievementSchema struct {
	Code        string                       `json:"code"`
	Description string                       `json:"description"`
	Name        string                       `json:"name"`
	Objectives  []AchievementObjectiveSchema `json:"objectives"`
	Points      int                          `json:"points"`
	Rewards     AchievementRewardsSchema     `json:"rewards"`
}

type ActiveCharacterSchema struct {
	Account string `json:"account"`
	Layer   string `json:"layer"`
	MapId   int    `json:"map_id"`
	Name    string `json:"name"`
	Skin    string `json:"skin"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
}

type ActiveEventResponseSchema struct {
	Data ActiveEventSchema `json:"data"`
}

type ActiveEventSchema struct {
	Code        string    `json:"code"`
	CreatedAt   time.Time `json:"created_at"`
	Duration    int       `json:"duration"`
	Expiration  time.Time `json:"expiration"`
	Map         MapSchema `json:"map"`
	Name        string    `json:"name"`
	PreviousMap MapSchema `json:"previous_map"`
}

type AddAccountSchema struct {
	Email    openapi_types.Email `json:"email"`
	Password string              `json:"password"`
	Username string              `json:"username"`
}

type AddCharacterSchema struct {
	Name string `json:"name"`
	Skin string `json:"skin"`
}

type AssistantAnswerDataSchema struct {
	Answer       string               `json:"answer"`
	Assistant    RateLimitScopeSchema `json:"assistant"`
	PaidWithGems bool                 `json:"paid_with_gems"`
}

type AssistantAnswerSchema struct {
	Data AssistantAnswerDataSchema `json:"data"`
}

type AssistantQuestionSchema struct {
	PayWithGems *bool  `json:"pay_with_gems,omitempty"`
	Question    string `json:"question"`
}

type BadgeResponseSchema struct {
	Data BadgeSchema `json:"data"`
}

type BadgeSchema struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Season      *int   `json:"season,omitempty"`
}

type BankExtensionSchema struct {
	Price int `json:"price"`
}

type BankExtensionTransactionResponseSchema struct {
	Data BankExtensionTransactionSchema `json:"data"`
}

type BankExtensionTransactionSchema struct {
	Character   CharacterSchema     `json:"character"`
	Cooldown    CooldownSchema      `json:"cooldown"`
	Transaction BankExtensionSchema `json:"transaction"`
}

type BankGoldTransactionResponseSchema struct {
	Data BankGoldTransactionSchema `json:"data"`
}

type BankGoldTransactionSchema struct {
	Bank      GoldSchema      `json:"bank"`
	Character CharacterSchema `json:"character"`
	Cooldown  CooldownSchema  `json:"cooldown"`
}

type BankItemTransactionResponseSchema struct {
	Data BankItemTransactionSchema `json:"data"`
}

type BankItemTransactionSchema struct {
	Bank      []SimpleItemSchema `json:"bank"`
	Character CharacterSchema    `json:"character"`
	Cooldown  CooldownSchema     `json:"cooldown"`
	Items     []SimpleItemSchema `json:"items"`
}

type BankResponseSchema struct {
	Data BankSchema `json:"data"`
}

type BankSchema struct {
	Expansions        int `json:"expansions"`
	Gold              int `json:"gold"`
	NextExpansionCost int `json:"next_expansion_cost"`
	Slots             int `json:"slots"`
}

type BuyCustomDesignRequestSchema struct {
	Code string `json:"code"`
}

type BuySkinRequestSchema struct {
	Code string `json:"code"`
}

type BuySkinResponseDataSchema struct {
	Gems  int      `json:"gems"`
	Skin  string   `json:"skin"`
	Skins []string `json:"skins"`
}

type BuySkinResponseSchema struct {
	Data BuySkinResponseDataSchema `json:"data"`
}

type ChangeEmailSchema struct {
	CurrentEmail openapi_types.Email `json:"current_email"`
	NewEmail     openapi_types.Email `json:"new_email"`
}

type ChangePasswordSchema struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ChangeSkinCharacterDataSchema struct {
	Character CharacterSchema `json:"character"`
	Cooldown  CooldownSchema  `json:"cooldown"`
	Skin      string          `json:"skin"`
}

type ChangeSkinCharacterSchema struct {
	Skin string `json:"skin"`
}

type ChangeSkinResponseSchema struct {
	Data ChangeSkinCharacterDataSchema `json:"data"`
}

type CharacterFightDataSchema struct {
	Characters []CharacterSchema    `json:"characters"`
	Cooldown   CooldownSchema       `json:"cooldown"`
	Fight      CharacterFightSchema `json:"fight"`
}

type CharacterFightResponseSchema struct {
	Data CharacterFightDataSchema `json:"data"`
}

type CharacterFightSchema struct {
	Characters []CharacterMultiFightResultSchema `json:"characters"`
	Logs       []string                          `json:"logs"`
	Opponent   string                            `json:"opponent"`
	Result     string                            `json:"result"`
	Turns      int                               `json:"turns"`
}

type CharacterLeaderboardSchema struct {
	Account                string `json:"account"`
	AlchemyLevel           int    `json:"alchemy_level"`
	AlchemyTotalXp         int    `json:"alchemy_total_xp"`
	CookingLevel           int    `json:"cooking_level"`
	CookingTotalXp         int    `json:"cooking_total_xp"`
	FishingLevel           int    `json:"fishing_level"`
	FishingTotalXp         int    `json:"fishing_total_xp"`
	GearcraftingLevel      int    `json:"gearcrafting_level"`
	GearcraftingTotalXp    int    `json:"gearcrafting_total_xp"`
	Gold                   int    `json:"gold"`
	JewelrycraftingLevel   int    `json:"jewelrycrafting_level"`
	JewelrycraftingTotalXp int    `json:"jewelrycrafting_total_xp"`
	Level                  int    `json:"level"`
	Member                 bool   `json:"member"`
	MiningLevel            int    `json:"mining_level"`
	MiningTotalXp          int    `json:"mining_total_xp"`
	Name                   string `json:"name"`
	Position               int    `json:"position"`
	Skin                   string `json:"skin"`
	TotalXp                int    `json:"total_xp"`
	WeaponcraftingLevel    int    `json:"weaponcrafting_level"`
	WeaponcraftingTotalXp  int    `json:"weaponcrafting_total_xp"`
	WoodcuttingLevel       int    `json:"woodcutting_level"`
	WoodcuttingTotalXp     int    `json:"woodcutting_total_xp"`
}

type CharacterMovementDataSchema struct {
	Character   CharacterSchema `json:"character"`
	Cooldown    CooldownSchema  `json:"cooldown"`
	Destination MapSchema       `json:"destination"`
	Path        [][2]int        `json:"path"`
}

type CharacterMovementResponseSchema struct {
	Data CharacterMovementDataSchema `json:"data"`
}

type CharacterMultiFightResultSchema struct {
	CharacterName string       `json:"character_name"`
	Drops         []DropSchema `json:"drops"`
	FinalHp       int          `json:"final_hp"`
	Gold          int          `json:"gold"`
	Xp            int          `json:"xp"`
}

type CharacterResponseSchema struct {
	Data CharacterSchema `json:"data"`
}

type CharacterRestDataSchema struct {
	Character  CharacterSchema `json:"character"`
	Cooldown   CooldownSchema  `json:"cooldown"`
	HpRestored int             `json:"hp_restored"`
}

type CharacterRestResponseSchema struct {
	Data CharacterRestDataSchema `json:"data"`
}

type CharacterSchema struct {
	Account              string                 `json:"account"`
	AlchemyLevel         int                    `json:"alchemy_level"`
	AlchemyMaxXp         int                    `json:"alchemy_max_xp"`
	AlchemyXp            int                    `json:"alchemy_xp"`
	AmuletSlot           string                 `json:"amulet_slot"`
	Artifact1Slot        string                 `json:"artifact1_slot"`
	Artifact2Slot        string                 `json:"artifact2_slot"`
	Artifact3Slot        string                 `json:"artifact3_slot"`
	AttackAir            int                    `json:"attack_air"`
	AttackEarth          int                    `json:"attack_earth"`
	AttackFire           int                    `json:"attack_fire"`
	AttackWater          int                    `json:"attack_water"`
	BagSlot              string                 `json:"bag_slot"`
	BodyArmorSlot        string                 `json:"body_armor_slot"`
	BootsSlot            string                 `json:"boots_slot"`
	CookingLevel         int                    `json:"cooking_level"`
	CookingMaxXp         int                    `json:"cooking_max_xp"`
	CookingXp            int                    `json:"cooking_xp"`
	Cooldown             int                    `json:"cooldown"`
	CooldownExpiration   *time.Time             `json:"cooldown_expiration,omitempty"`
	CriticalStrike       int                    `json:"critical_strike"`
	Dmg                  int                    `json:"dmg"`
	DmgAir               int                    `json:"dmg_air"`
	DmgEarth             int                    `json:"dmg_earth"`
	DmgFire              int                    `json:"dmg_fire"`
	DmgWater             int                    `json:"dmg_water"`
	Effects              *[]StorageEffectSchema `json:"effects,omitempty"`
	FishingLevel         int                    `json:"fishing_level"`
	FishingMaxXp         int                    `json:"fishing_max_xp"`
	FishingXp            int                    `json:"fishing_xp"`
	GearcraftingLevel    int                    `json:"gearcrafting_level"`
	GearcraftingMaxXp    int                    `json:"gearcrafting_max_xp"`
	GearcraftingXp       int                    `json:"gearcrafting_xp"`
	Gold                 int                    `json:"gold"`
	Haste                int                    `json:"haste"`
	HelmetSlot           string                 `json:"helmet_slot"`
	Hp                   int                    `json:"hp"`
	Initiative           int                    `json:"initiative"`
	Inventory            *[]InventorySlotSchema `json:"inventory,omitempty"`
	InventoryMaxItems    int                    `json:"inventory_max_items"`
	JewelrycraftingLevel int                    `json:"jewelrycrafting_level"`
	JewelrycraftingMaxXp int                    `json:"jewelrycrafting_max_xp"`
	JewelrycraftingXp    int                    `json:"jewelrycrafting_xp"`
	Layer                string                 `json:"layer"`
	LegArmorSlot         string                 `json:"leg_armor_slot"`
	Level                int                    `json:"level"`
	MapId                int                    `json:"map_id"`
	MaxHp                int                    `json:"max_hp"`
	MaxXp                int                    `json:"max_xp"`
	MiningLevel          int                    `json:"mining_level"`
	MiningMaxXp          int                    `json:"mining_max_xp"`
	MiningXp             int                    `json:"mining_xp"`
	Name                 string                 `json:"name"`
	Prospecting          int                    `json:"prospecting"`
	ResAir               int                    `json:"res_air"`
	ResEarth             int                    `json:"res_earth"`
	ResFire              int                    `json:"res_fire"`
	ResWater             int                    `json:"res_water"`
	Ring1Slot            string                 `json:"ring1_slot"`
	Ring2Slot            string                 `json:"ring2_slot"`
	RuneSlot             string                 `json:"rune_slot"`
	ShieldSlot           string                 `json:"shield_slot"`
	Skin                 string                 `json:"skin"`
	Speed                int                    `json:"speed"`
	Task                 string                 `json:"task"`
	TaskProgress         int                    `json:"task_progress"`
	TaskTotal            int                    `json:"task_total"`
	TaskType             string                 `json:"task_type"`
	Threat               int                    `json:"threat"`
	Utility1Slot         string                 `json:"utility1_slot"`
	Utility1SlotQuantity int                    `json:"utility1_slot_quantity"`
	Utility2Slot         string                 `json:"utility2_slot"`
	Utility2SlotQuantity int                    `json:"utility2_slot_quantity"`
	WeaponSlot           string                 `json:"weapon_slot"`
	WeaponcraftingLevel  int                    `json:"weaponcrafting_level"`
	WeaponcraftingMaxXp  int                    `json:"weaponcrafting_max_xp"`
	WeaponcraftingXp     int                    `json:"weaponcrafting_xp"`
	Wisdom               int                    `json:"wisdom"`
	WoodcuttingLevel     int                    `json:"woodcutting_level"`
	WoodcuttingMaxXp     int                    `json:"woodcutting_max_xp"`
	WoodcuttingXp        int                    `json:"woodcutting_xp"`
	X                    int                    `json:"x"`
	Xp                   int                    `json:"xp"`
	Y                    int                    `json:"y"`
}

type CharacterStatsResponseSchema struct {
	Data CharacterStatsSchema `json:"data"`
}

type CharacterStatsSchema struct {
	ActionCounts      *map[string]int `json:"action_counts,omitempty"`
	Deaths            *int            `json:"deaths,omitempty"`
	ItemsCrafted      *map[string]int `json:"items_crafted,omitempty"`
	MonstersKilled    *map[string]int `json:"monsters_killed,omitempty"`
	ResourcesGathered *map[string]int `json:"resources_gathered,omitempty"`
}

type CharacterTransitionDataSchema struct {
	Character   CharacterSchema  `json:"character"`
	Cooldown    CooldownSchema   `json:"cooldown"`
	Destination MapSchema        `json:"destination"`
	Transition  TransitionSchema `json:"transition"`
}

type CharacterTransitionResponseSchema struct {
	Data CharacterTransitionDataSchema `json:"data"`
}

type CharactersListSchema struct {
	Data []CharacterSchema `json:"data"`
}

type CheckoutResponseSchema struct {
	CheckoutUrl string `json:"checkout_url"`
	SessionId   string `json:"session_id"`
}

type CheckoutResponseWrapperSchema struct {
	Data CheckoutResponseSchema `json:"data"`
}

type ClaimPendingItemDataSchema struct {
	Character CharacterSchema   `json:"character"`
	Cooldown  CooldownSchema    `json:"cooldown"`
	Item      PendingItemSchema `json:"item"`
}

type ClaimPendingItemResponseSchema struct {
	Data ClaimPendingItemDataSchema `json:"data"`
}

type CombatResultSchema struct {
	CharacterResults []map[string]interface{} `json:"character_results"`
	Logs             []string                 `json:"logs"`
	Result           string                   `json:"result"`
	Turns            int                      `json:"turns"`
}

type CombatSimulationDataSchema struct {
	Losses  int                  `json:"losses"`
	Results []CombatResultSchema `json:"results"`
	Winrate float32              `json:"winrate"`
	Wins    int                  `json:"wins"`
}

type CombatSimulationRequestSchema struct {
	Characters []FakeCharacterSchema `json:"characters"`
	Iterations int                   `json:"iterations"`
	Monster    string                `json:"monster"`
}

type CombatSimulationResponseSchema struct {
	Data CombatSimulationDataSchema `json:"data"`
}

type ConditionSchema struct {
	Code     string `json:"code"`
	Operator string `json:"operator"`
	Value    int    `json:"value"`
}

type CooldownSchema struct {
	Expiration       time.Time `json:"expiration"`
	Reason           string    `json:"reason"`
	RemainingSeconds int       `json:"remaining_seconds"`
	StartedAt        time.Time `json:"started_at"`
	TotalSeconds     int       `json:"total_seconds"`
}

type CraftSchema struct {
	Items    *[]SimpleItemSchema `json:"items,omitempty"`
	Level    *int                `json:"level,omitempty"`
	Quantity *int                `json:"quantity,omitempty"`
	Skill    *string             `json:"skill,omitempty"`
}

type CraftingSchema struct {
	Code     string `json:"code"`
	Quantity *int   `json:"quantity,omitempty"`
}

type DataPageAccountAchievementSchema struct {
	Data  []AccountAchievementSchema `json:"data"`
	Page  int                        `json:"page"`
	Pages int                        `json:"pages"`
	Size  int                        `json:"size"`
	Total int                        `json:"total"`
}

type DataPageAccountLeaderboardSchema struct {
	Data  []AccountLeaderboardSchema `json:"data"`
	Page  int                        `json:"page"`
	Pages int                        `json:"pages"`
	Size  int                        `json:"size"`
	Total int                        `json:"total"`
}

type DataPageActiveCharacterSchema struct {
	Data  []ActiveCharacterSchema `json:"data"`
	Page  int                     `json:"page"`
	Pages int                     `json:"pages"`
	Size  int                     `json:"size"`
	Total int                     `json:"total"`
}

type DataPageCharacterLeaderboardSchema struct {
	Data  []CharacterLeaderboardSchema `json:"data"`
	Page  int                          `json:"page"`
	Pages int                          `json:"pages"`
	Size  int                          `json:"size"`
	Total int                          `json:"total"`
}

type DataPageGEOrderHistorySchema struct {
	Data  []GEOrderHistorySchema `json:"data"`
	Page  int                    `json:"page"`
	Pages int                    `json:"pages"`
	Size  int                    `json:"size"`
	Total int                    `json:"total"`
}

type DataPageGEOrderSchema struct {
	Data  []GEOrderSchema `json:"data"`
	Page  int             `json:"page"`
	Pages int             `json:"pages"`
	Size  int             `json:"size"`
	Total int             `json:"total"`
}

type DataPageLogSchema struct {
	Data  []LogSchema `json:"data"`
	Page  int         `json:"page"`
	Pages int         `json:"pages"`
	Size  int         `json:"size"`
	Total int         `json:"total"`
}

type DataPagePendingItemSchema struct {
	Data  []PendingItemSchema `json:"data"`
	Page  int                 `json:"page"`
	Pages int                 `json:"pages"`
	Size  int                 `json:"size"`
	Total int                 `json:"total"`
}

type DataPageRaidLeaderboardEntrySchema struct {
	Data  []RaidLeaderboardEntrySchema `json:"data"`
	Page  int                          `json:"page"`
	Pages int                          `json:"pages"`
	Size  int                          `json:"size"`
	Total int                          `json:"total"`
}

type DataPageSimpleItemSchema struct {
	Data  []SimpleItemSchema `json:"data"`
	Page  int                `json:"page"`
	Pages int                `json:"pages"`
	Size  int                `json:"size"`
	Total int                `json:"total"`
}

type DeleteCharacterSchema struct {
	Name string `json:"name"`
}

type DeleteItemResponseSchema struct {
	Data DeleteItemSchema `json:"data"`
}

type DeleteItemSchema struct {
	Character CharacterSchema  `json:"character"`
	Cooldown  CooldownSchema   `json:"cooldown"`
	Item      SimpleItemSchema `json:"item"`
}

type DepositWithdrawGoldSchema struct {
	Quantity int `json:"quantity"`
}

type DestinationSchema struct {
	MapId *int `json:"map_id,omitempty"`
	X     *int `json:"x,omitempty"`
	Y     *int `json:"y,omitempty"`
}

type DropRateSchema struct {
	Code        string `json:"code"`
	MaxQuantity int    `json:"max_quantity"`
	MinQuantity int    `json:"min_quantity"`
	Rate        int    `json:"rate"`
}

type DropSchema struct {
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
}

type EffectResponseSchema struct {
	Data EffectSchema `json:"data"`
}

type EffectSchema struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Subtype     string `json:"subtype"`
	Type        string `json:"type"`
}

type EquipSchema struct {
	Code     string `json:"code"`
	Quantity *int   `json:"quantity,omitempty"`
	Slot     string `json:"slot"`
}

type EquipmentItemSchema struct {
	Item ItemSchema `json:"item"`
	Slot string     `json:"slot"`
}

type EquipmentResponseSchema struct {
	Data EquipmentTransactionSchema `json:"data"`
}

type EquipmentTransactionSchema struct {
	Character CharacterSchema       `json:"character"`
	Cooldown  CooldownSchema        `json:"cooldown"`
	Items     []EquipmentItemSchema `json:"items"`
}

type ErrorResponseSchema struct {
	Error ErrorSchema `json:"error"`
}

type ErrorSchema struct {
	Code    int                     `json:"code"`
	Data    *map[string]interface{} `json:"data,omitempty"`
	Message string                  `json:"message"`
}

type EventContentSchema struct {
	Code string `json:"code"`
	Type string `json:"type"`
}

type EventMapSchema struct {
	Layer string `json:"layer"`
	MapId int    `json:"map_id"`
	Skin  string `json:"skin"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
}

type EventSchema struct {
	Code               string              `json:"code"`
	Content            *EventContentSchema `json:"content,omitempty"`
	Cooldown           *int                `json:"cooldown,omitempty"`
	CooldownExpiration *time.Time          `json:"cooldown_expiration,omitempty"`
	Duration           int                 `json:"duration"`
	Maps               []EventMapSchema    `json:"maps"`
	Name               string              `json:"name"`
	Price              *int                `json:"price,omitempty"`
	Rate               int                 `json:"rate"`
	Transition         *TransitionSchema   `json:"transition,omitempty"`
}

type FakeCharacterSchema struct {
	AmuletSlot           *string `json:"amulet_slot,omitempty"`
	Artifact1Slot        *string `json:"artifact1_slot,omitempty"`
	Artifact2Slot        *string `json:"artifact2_slot,omitempty"`
	Artifact3Slot        *string `json:"artifact3_slot,omitempty"`
	BodyArmorSlot        *string `json:"body_armor_slot,omitempty"`
	BootsSlot            *string `json:"boots_slot,omitempty"`
	HelmetSlot           *string `json:"helmet_slot,omitempty"`
	LegArmorSlot         *string `json:"leg_armor_slot,omitempty"`
	Level                int     `json:"level"`
	Ring1Slot            *string `json:"ring1_slot,omitempty"`
	Ring2Slot            *string `json:"ring2_slot,omitempty"`
	RuneSlot             *string `json:"rune_slot,omitempty"`
	ShieldSlot           *string `json:"shield_slot,omitempty"`
	Utility1Slot         *string `json:"utility1_slot,omitempty"`
	Utility1SlotQuantity *int    `json:"utility1_slot_quantity,omitempty"`
	Utility2Slot         *string `json:"utility2_slot,omitempty"`
	Utility2SlotQuantity *int    `json:"utility2_slot_quantity,omitempty"`
	WeaponSlot           *string `json:"weapon_slot,omitempty"`
}

type FightRequestSchema struct {
	Participants *[]string `json:"participants,omitempty"`
}

type GEBuyOrderCreationSchema struct {
	Code     string `json:"code"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type GEBuyOrderSchema struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type GECancelOrderSchema struct {
	Id string `json:"id"`
}

type GECreateOrderTransactionResponseSchema struct {
	Data GEOrderTransactionSchema `json:"data"`
}

type GEFillBuyOrderSchema struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type GEOrderCreatedSchema struct {
	Code       string    `json:"code"`
	CreatedAt  time.Time `json:"created_at"`
	Id         string    `json:"id"`
	Price      int       `json:"price"`
	Quantity   int       `json:"quantity"`
	TotalPrice int       `json:"total_price"`
}

type GEOrderCreationSchema struct {
	Code     string `json:"code"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type GEOrderHistorySchema struct {
	Buyer    string    `json:"buyer"`
	Code     string    `json:"code"`
	OrderId  string    `json:"order_id"`
	Price    int       `json:"price"`
	Quantity int       `json:"quantity"`
	Seller   string    `json:"seller"`
	SoldAt   time.Time `json:"sold_at"`
}

type GEOrderResponseSchema struct {
	Data GEOrderSchema `json:"data"`
}

type GEOrderSchema struct {
	Account   *string   `json:"account,omitempty"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	Id        string    `json:"id"`
	Price     int       `json:"price"`
	Quantity  int       `json:"quantity"`
	Type      string    `json:"type"`
}

type GEOrderTransactionSchema struct {
	Character CharacterSchema      `json:"character"`
	Cooldown  CooldownSchema       `json:"cooldown"`
	Order     GEOrderCreatedSchema `json:"order"`
}

type GETransactionListSchema struct {
	Character CharacterSchema     `json:"character"`
	Cooldown  CooldownSchema      `json:"cooldown"`
	Order     GETransactionSchema `json:"order"`
}

type GETransactionResponseSchema struct {
	Data GETransactionListSchema `json:"data"`
}

type GETransactionSchema struct {
	Code       string `json:"code"`
	Id         string `json:"id"`
	Price      int    `json:"price"`
	Quantity   int    `json:"quantity"`
	TotalPrice int    `json:"total_price"`
}

type GemShopCatalogDataSchema struct {
	CustomDesigns []GemShopCustomDesignCatalogItemSchema `json:"custom_designs"`
	Skins         []GemShopSkinCatalogItemSchema         `json:"skins"`
	SpawnEvents   []GemShopSpawnEventCatalogItemSchema   `json:"spawn_events"`
	Subscriptions []GemShopSubscriptionCatalogItemSchema `json:"subscriptions"`
}

type GemShopCatalogResponseSchema struct {
	Data GemShopCatalogDataSchema `json:"data"`
}

type GemShopCustomDesignCatalogItemSchema struct {
	Category        string `json:"category"`
	Code            string `json:"code"`
	Description     string `json:"description"`
	Name            string `json:"name"`
	Price           int    `json:"price"`
	UniqueToAccount bool   `json:"unique_to_account"`
}

type GemShopCustomDesignPurchaseResponseDataSchema struct {
	Code string `json:"code"`
	Cost int    `json:"cost"`
	Gems int    `json:"gems"`
	Name string `json:"name"`
}

type GemShopCustomDesignPurchaseResponseSchema struct {
	Data GemShopCustomDesignPurchaseResponseDataSchema `json:"data"`
}

type GemShopSkinCatalogItemSchema struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
}

type GemShopSpawnEventCatalogItemSchema struct {
	Code        string `json:"code"`
	ContentCode string `json:"content_code"`
	ContentType string `json:"content_type"`
	Duration    int    `json:"duration"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
}

type GemShopSubscriptionCatalogItemSchema struct {
	Code         string `json:"code"`
	DurationDays int    `json:"duration_days"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
}

type GemShopSubscriptionResponseDataSchema struct {
	Cost             int       `json:"cost"`
	Gems             int       `json:"gems"`
	Member           bool      `json:"member"`
	MemberExpiration time.Time `json:"member_expiration"`
}

type GemShopSubscriptionResponseSchema struct {
	Data GemShopSubscriptionResponseDataSchema `json:"data"`
}

type GemTransactionListResponseSchema struct {
	Data []GemTransactionSchema `json:"data"`
}

type GemTransactionSchema struct {
	CreatedAt   time.Time              `json:"created_at"`
	Description string                 `json:"description"`
	Gems        int                    `json:"gems"`
	Metadata    map[string]interface{} `json:"metadata"`
	Type        string                 `json:"type"`
}

type GiveGoldDataSchema struct {
	Character         CharacterSchema `json:"character"`
	Cooldown          CooldownSchema  `json:"cooldown"`
	Quantity          int             `json:"quantity"`
	ReceiverCharacter CharacterSchema `json:"receiver_character"`
}

type GiveGoldResponseSchema struct {
	Data GiveGoldDataSchema `json:"data"`
}

type GiveGoldSchema struct {
	Character string `json:"character"`
	Quantity  int    `json:"quantity"`
}

type GiveItemDataSchema struct {
	Character         CharacterSchema    `json:"character"`
	Cooldown          CooldownSchema     `json:"cooldown"`
	Items             []SimpleItemSchema `json:"items"`
	ReceiverCharacter CharacterSchema    `json:"receiver_character"`
}

type GiveItemResponseSchema struct {
	Data GiveItemDataSchema `json:"data"`
}

type GiveItemsSchema struct {
	Character string             `json:"character"`
	Items     []SimpleItemSchema `json:"items"`
}

type GoldSchema struct {
	Quantity int `json:"quantity"`
}

type InteractionSchema struct {
	Content    *MapContentSchema `json:"content,omitempty"`
	Transition *TransitionSchema `json:"transition,omitempty"`
}

type InventorySlotSchema struct {
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
	Slot     int    `json:"slot"`
}

type ItemResponseSchema struct {
	Data ItemSchema `json:"data"`
}

type ItemSchema struct {
	Code        string                `json:"code"`
	Conditions  *[]ConditionSchema    `json:"conditions,omitempty"`
	Craft       *CraftSchema          `json:"craft,omitempty"`
	Description string                `json:"description"`
	Effects     *[]SimpleEffectSchema `json:"effects,omitempty"`
	Level       int                   `json:"level"`
	Name        string                `json:"name"`
	Recyclable  *bool                 `json:"recyclable,omitempty"`
	Subtype     string                `json:"subtype"`
	Tradeable   bool                  `json:"tradeable"`
	Type        string                `json:"type"`
}

type LogSchema struct {
	Account            string      `json:"account"`
	Character          string      `json:"character"`
	Content            interface{} `json:"content"`
	Cooldown           int         `json:"cooldown"`
	CooldownExpiration *time.Time  `json:"cooldown_expiration,omitempty"`
	CreatedAt          time.Time   `json:"created_at"`
	Description        string      `json:"description"`
	Type               string      `json:"type"`
}

type MapContentSchema struct {
	Code string `json:"code"`
	Type string `json:"type"`
}

type MapResponseSchema struct {
	Data MapSchema `json:"data"`
}

type MapSchema struct {
	Access       AccessSchema      `json:"access"`
	Interactions InteractionSchema `json:"interactions"`
	Layer        string            `json:"layer"`
	MapId        int               `json:"map_id"`
	Name         string            `json:"name"`
	Skin         string            `json:"skin"`
	X            int               `json:"x"`
	Y            int               `json:"y"`
}

type MemberTokenSubscriptionResponseDataSchema struct {
	Member           bool      `json:"member"`
	MemberExpiration time.Time `json:"member_expiration"`
	MemberToken      int       `json:"member_token"`
}

type MemberTokenSubscriptionResponseSchema struct {
	Data MemberTokenSubscriptionResponseDataSchema `json:"data"`
}

type MonsterResponseSchema struct {
	Data MonsterSchema `json:"data"`
}

type MonsterSchema struct {
	AttackAir      int                   `json:"attack_air"`
	AttackEarth    int                   `json:"attack_earth"`
	AttackFire     int                   `json:"attack_fire"`
	AttackWater    int                   `json:"attack_water"`
	Code           string                `json:"code"`
	CriticalStrike int                   `json:"critical_strike"`
	Drops          []DropRateSchema      `json:"drops"`
	Effects        *[]SimpleEffectSchema `json:"effects,omitempty"`
	Hp             int                   `json:"hp"`
	Initiative     int                   `json:"initiative"`
	Level          int                   `json:"level"`
	MaxGold        int                   `json:"max_gold"`
	MinGold        int                   `json:"min_gold"`
	Name           string                `json:"name"`
	ResAir         int                   `json:"res_air"`
	ResEarth       int                   `json:"res_earth"`
	ResFire        int                   `json:"res_fire"`
	ResWater       int                   `json:"res_water"`
	Type           string                `json:"type"`
}

type MyAccountDetails struct {
	AchievementsPoints int                 `json:"achievements_points"`
	Badges             *[]string           `json:"badges,omitempty"`
	BanReason          *string             `json:"ban_reason,omitempty"`
	Banned             bool                `json:"banned"`
	Email              openapi_types.Email `json:"email"`
	Gems               int                 `json:"gems"`
	Member             bool                `json:"member"`
	MemberExpiration   *time.Time          `json:"member_expiration,omitempty"`
	MemberToken        *int                `json:"member_token,omitempty"`
	Skins              []string            `json:"skins"`
	Status             string              `json:"status"`
	Username           string              `json:"username"`
}

type MyAccountDetailsSchema struct {
	Data MyAccountDetails `json:"data"`
}

type MyCharactersListSchema struct {
	Data []CharacterSchema `json:"data"`
}

type NPCItemSchema struct {
	BuyPrice  *int   `json:"buy_price,omitempty"`
	Code      string `json:"code"`
	Currency  string `json:"currency"`
	Npc       string `json:"npc"`
	SellPrice *int   `json:"sell_price,omitempty"`
}

type NPCResponseSchema struct {
	Data NPCSchema `json:"data"`
}

type NPCSchema struct {
	Code        string                 `json:"code"`
	Description string                 `json:"description"`
	Items       *[]SimpleNPCItemSchema `json:"items,omitempty"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
}

type NpcItemTransactionSchema struct {
	Code       string `json:"code"`
	Currency   string `json:"currency"`
	Price      int    `json:"price"`
	Quantity   int    `json:"quantity"`
	TotalPrice int    `json:"total_price"`
}

type NpcMerchantBuySchema struct {
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
}

type NpcMerchantTransactionResponseSchema struct {
	Data NpcMerchantTransactionSchema `json:"data"`
}

type NpcMerchantTransactionSchema struct {
	Character   CharacterSchema          `json:"character"`
	Cooldown    CooldownSchema           `json:"cooldown"`
	Transaction NpcItemTransactionSchema `json:"transaction"`
}

type PasswordResetConfirmSchema struct {
	NewPassword string `json:"new_password"`
	Token       string `json:"token"`
}

type PasswordResetRequestSchema struct {
	Email openapi_types.Email `json:"email"`
}

type PasswordResetResponseSchema struct {
	Message string `json:"message"`
}

type PendingItemSchema struct {
	Account     string              `json:"account"`
	ClaimedAt   *time.Time          `json:"claimed_at,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	Description string              `json:"description"`
	Gold        *int                `json:"gold,omitempty"`
	Id          string              `json:"id"`
	Items       *[]SimpleItemSchema `json:"items,omitempty"`
	Source      string              `json:"source"`
	SourceId    *string             `json:"source_id,omitempty"`
}

type PurchaseGemsRequestSchema struct {
	Quantity int `json:"quantity"`
}

type PurchaseHistoryListResponseSchema struct {
	Data []PurchaseHistorySchema `json:"data"`
}

type PurchaseHistorySchema struct {
	Amount       int       `json:"amount"`
	CreatedAt    time.Time `json:"created_at"`
	Description  string    `json:"description"`
	GemsCredited *int      `json:"gems_credited,omitempty"`
	Type         string    `json:"type"`
}

type RaidDamageRewardSchema struct {
	DamagePerReward int                 `json:"damage_per_reward"`
	Items           *[]SimpleItemSchema `json:"items,omitempty"`
	MaxRewards      *int                `json:"max_rewards,omitempty"`
}

type RaidInstanceSchema struct {
	EndedAt              *time.Time `json:"ended_at,omitempty"`
	EndsAt               time.Time  `json:"ends_at"`
	ParticipantCount     *int       `json:"participant_count,omitempty"`
	RemainingHp          int        `json:"remaining_hp"`
	Result               *string    `json:"result,omitempty"`
	RewardsDistributedAt *time.Time `json:"rewards_distributed_at,omitempty"`
	StartsAt             time.Time  `json:"starts_at"`
	Status               string     `json:"status"`
	TotalHp              int        `json:"total_hp"`
}

type RaidLeaderboardEntrySchema struct {
	Account  string `json:"account"`
	Points   int    `json:"points"`
	Position int    `json:"position"`
}

type RaidRankRewardSchema struct {
	Items   *[]SimpleItemSchema `json:"items,omitempty"`
	MaxRank int                 `json:"max_rank"`
	MinRank int                 `json:"min_rank"`
}

type RaidResponseSchema struct {
	Data RaidSchema `json:"data"`
}

type RaidRewardsSchema struct {
	DamageRewards *[]RaidDamageRewardSchema `json:"damage_rewards,omitempty"`
	Leaderboard   *[]RaidRankRewardSchema   `json:"leaderboard,omitempty"`
}

type RaidScheduleSchema struct {
	DurationHours  *int     `json:"duration_hours,omitempty"`
	StartHourUtc   *int     `json:"start_hour_utc,omitempty"`
	StartMinuteUtc *int     `json:"start_minute_utc,omitempty"`
	Weekdays       []string `json:"weekdays"`
}

type RaidSchema struct {
	ActiveInstance   *RaidInstanceSchema `json:"active_instance,omitempty"`
	Code             string              `json:"code"`
	Description      *string             `json:"description,omitempty"`
	LatestInstance   *RaidInstanceSchema `json:"latest_instance,omitempty"`
	Monster          string              `json:"monster"`
	Name             string              `json:"name"`
	NextStartAt      time.Time           `json:"next_start_at"`
	ParticipantCount *int                `json:"participant_count,omitempty"`
	Rewards          *RaidRewardsSchema  `json:"rewards,omitempty"`
	Schedule         RaidScheduleSchema  `json:"schedule"`
	Status           string              `json:"status"`
}

type RateLimitSchema struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type RateLimitScopeSchema struct {
	Day    *RateLimitWindowSchema `json:"day,omitempty"`
	Hour   *RateLimitWindowSchema `json:"hour,omitempty"`
	Minute *RateLimitWindowSchema `json:"minute,omitempty"`
	Second *RateLimitWindowSchema `json:"second,omitempty"`
}

type RateLimitWindowSchema struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Reset     time.Time `json:"reset"`
}

type RateLimitsDataSchema struct {
	Account    RateLimitScopeSchema  `json:"account"`
	Action     RateLimitScopeSchema  `json:"action"`
	Assistant  *RateLimitScopeSchema `json:"assistant,omitempty"`
	Data       RateLimitScopeSchema  `json:"data"`
	Simulation RateLimitScopeSchema  `json:"simulation"`
}

type RateLimitsSchema struct {
	Data RateLimitsDataSchema `json:"data"`
}

type RecyclingDataSchema struct {
	Character CharacterSchema      `json:"character"`
	Cooldown  CooldownSchema       `json:"cooldown"`
	Details   RecyclingItemsSchema `json:"details"`
}

type RecyclingItemsSchema struct {
	Enhanced *bool        `json:"enhanced,omitempty"`
	Gold     *int         `json:"gold,omitempty"`
	Items    []DropSchema `json:"items"`
}

type RecyclingResponseSchema struct {
	Data RecyclingDataSchema `json:"data"`
}

type RecyclingSchema struct {
	Code     string `json:"code"`
	Enhanced *bool  `json:"enhanced,omitempty"`
	Quantity *int   `json:"quantity,omitempty"`
}

type RenameCharacterDataSchema struct {
	Character CharacterSchema `json:"character"`
	Cooldown  CooldownSchema  `json:"cooldown"`
	NewName   string          `json:"new_name"`
	OldName   string          `json:"old_name"`
}

type RenameCharacterSchema struct {
	Name string `json:"name"`
}

type RenameResponseSchema struct {
	Data RenameCharacterDataSchema `json:"data"`
}

type ResourceResponseSchema struct {
	Data ResourceSchema `json:"data"`
}

type ResourceSchema struct {
	Code  string           `json:"code"`
	Drops []DropRateSchema `json:"drops"`
	Level int              `json:"level"`
	Name  string           `json:"name"`
	Skill string           `json:"skill"`
}

type ResponseSchema struct {
	Message string `json:"message"`
}

type RewardDataResponseSchema struct {
	Data RewardDataSchema `json:"data"`
}

type RewardDataSchema struct {
	Character CharacterSchema `json:"character"`
	Cooldown  CooldownSchema  `json:"cooldown"`
	Rewards   RewardsSchema   `json:"rewards"`
}

type RewardItemSchema struct {
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
}

type RewardResponseSchema struct {
	Data DropRateSchema `json:"data"`
}

type RewardsSchema struct {
	Gold  int                `json:"gold"`
	Items []SimpleItemSchema `json:"items"`
}

type SandboxCharacterActionSchema struct {
	Character string `json:"character"`
}

type SandboxGiveGoldDataSchema struct {
	Character CharacterSchema `json:"character"`
	Cooldown  CooldownSchema  `json:"cooldown"`
	Quantity  int             `json:"quantity"`
}

type SandboxGiveGoldResponseSchema struct {
	Data SandboxGiveGoldDataSchema `json:"data"`
}

type SandboxGiveItemDataSchema struct {
	Character CharacterSchema  `json:"character"`
	Cooldown  CooldownSchema   `json:"cooldown"`
	Item      SimpleItemSchema `json:"item"`
}

type SandboxGiveItemResponseSchema struct {
	Data SandboxGiveItemDataSchema `json:"data"`
}

type SandboxGiveItemSchema struct {
	Character string `json:"character"`
	Code      string `json:"code"`
	Quantity  int    `json:"quantity"`
}

type SandboxGiveXPDataSchema struct {
	Amount    int             `json:"amount"`
	Character CharacterSchema `json:"character"`
	Cooldown  CooldownSchema  `json:"cooldown"`
	Type      string          `json:"type"`
}

type SandboxGiveXPResponseSchema struct {
	Data SandboxGiveXPDataSchema `json:"data"`
}

type SandboxGiveXPSchema struct {
	Amount    int    `json:"amount"`
	Character string `json:"character"`
	Type      string `json:"type"`
}

type SandboxResponseSchema struct {
	Data SandboxSchema `json:"data"`
}

type SandboxSchema struct {
	Character CharacterSchema `json:"character"`
}

type SandboxTeleportDataSchema struct {
	Character   CharacterSchema `json:"character"`
	Destination MapSchema       `json:"destination"`
}

type SandboxTeleportResponseSchema struct {
	Data SandboxTeleportDataSchema `json:"data"`
}

type SandboxTeleportSchema struct {
	Character string `json:"character"`
	MapId     int    `json:"map_id"`
}

type SeasonRewardSchema struct {
	Code           string `json:"code"`
	Description    string `json:"description"`
	FirstOnly      *bool  `json:"first_only,omitempty"`
	MemberRequired *bool  `json:"member_required,omitempty"`
	Quantity       *int   `json:"quantity,omitempty"`
	RequiredPoints int    `json:"required_points"`
	Type           string `json:"type"`
}

type SeasonSchema struct {
	Name      *string                    `json:"name,omitempty"`
	Number    *int                       `json:"number,omitempty"`
	Rewards   []StatusSeasonRewardSchema `json:"rewards"`
	StartDate *time.Time                 `json:"start_date,omitempty"`
}

type SimpleEffectSchema struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Value       int    `json:"value"`
}

type SimpleItemSchema struct {
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
}

type SimpleNPCItemSchema struct {
	BuyPrice  *int   `json:"buy_price,omitempty"`
	Code      string `json:"code"`
	Currency  string `json:"currency"`
	SellPrice *int   `json:"sell_price,omitempty"`
}

type SkillDataSchema struct {
	Character CharacterSchema `json:"character"`
	Cooldown  CooldownSchema  `json:"cooldown"`
	Details   SkillInfoSchema `json:"details"`
}

type SkillInfoSchema struct {
	Items []DropSchema `json:"items"`
	Xp    int          `json:"xp"`
}

type SkillResponseSchema struct {
	Data SkillDataSchema `json:"data"`
}

type SkinResponseSchema struct {
	Data SkinSchema `json:"data"`
}

type SkinSchema struct {
	Code        string `json:"code"`
	Default     bool   `json:"default"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Price       *int   `json:"price,omitempty"`
}

type SpawnEventRequestSchema struct {
	Code string `json:"code"`
}

type StaticDataPageAchievementSchema struct {
	Data  []AchievementSchema `json:"data"`
	Page  int                 `json:"page"`
	Pages int                 `json:"pages"`
	Size  int                 `json:"size"`
	Total int                 `json:"total"`
}

type StaticDataPageActiveEventSchema struct {
	Data  []ActiveEventSchema `json:"data"`
	Page  int                 `json:"page"`
	Pages int                 `json:"pages"`
	Size  int                 `json:"size"`
	Total int                 `json:"total"`
}

type StaticDataPageBadgeSchema struct {
	Data  []BadgeSchema `json:"data"`
	Page  int           `json:"page"`
	Pages int           `json:"pages"`
	Size  int           `json:"size"`
	Total int           `json:"total"`
}

type StaticDataPageDropRateSchema struct {
	Data  []DropRateSchema `json:"data"`
	Page  int              `json:"page"`
	Pages int              `json:"pages"`
	Size  int              `json:"size"`
	Total int              `json:"total"`
}

type StaticDataPageEffectSchema struct {
	Data  []EffectSchema `json:"data"`
	Page  int            `json:"page"`
	Pages int            `json:"pages"`
	Size  int            `json:"size"`
	Total int            `json:"total"`
}

type StaticDataPageEventSchema struct {
	Data  []EventSchema `json:"data"`
	Page  int           `json:"page"`
	Pages int           `json:"pages"`
	Size  int           `json:"size"`
	Total int           `json:"total"`
}

type StaticDataPageItemSchema struct {
	Data  []ItemSchema `json:"data"`
	Page  int          `json:"page"`
	Pages int          `json:"pages"`
	Size  int          `json:"size"`
	Total int          `json:"total"`
}

type StaticDataPageMapSchema struct {
	Data  []MapSchema `json:"data"`
	Page  int         `json:"page"`
	Pages int         `json:"pages"`
	Size  int         `json:"size"`
	Total int         `json:"total"`
}

type StaticDataPageMonsterSchema struct {
	Data  []MonsterSchema `json:"data"`
	Page  int             `json:"page"`
	Pages int             `json:"pages"`
	Size  int             `json:"size"`
	Total int             `json:"total"`
}

type StaticDataPageNPCItemSchema struct {
	Data  []NPCItemSchema `json:"data"`
	Page  int             `json:"page"`
	Pages int             `json:"pages"`
	Size  int             `json:"size"`
	Total int             `json:"total"`
}

type StaticDataPageNPCSchema struct {
	Data  []NPCSchema `json:"data"`
	Page  int         `json:"page"`
	Pages int         `json:"pages"`
	Size  int         `json:"size"`
	Total int         `json:"total"`
}

type StaticDataPageRaidSchema struct {
	Data  []RaidSchema `json:"data"`
	Page  int          `json:"page"`
	Pages int          `json:"pages"`
	Size  int          `json:"size"`
	Total int          `json:"total"`
}

type StaticDataPageResourceSchema struct {
	Data  []ResourceSchema `json:"data"`
	Page  int              `json:"page"`
	Pages int              `json:"pages"`
	Size  int              `json:"size"`
	Total int              `json:"total"`
}

type StaticDataPageSeasonRewardSchema struct {
	Data  []SeasonRewardSchema `json:"data"`
	Page  int                  `json:"page"`
	Pages int                  `json:"pages"`
	Size  int                  `json:"size"`
	Total int                  `json:"total"`
}

type StaticDataPageSkinSchema struct {
	Data  []SkinSchema `json:"data"`
	Page  int          `json:"page"`
	Pages int          `json:"pages"`
	Size  int          `json:"size"`
	Total int          `json:"total"`
}

type StaticDataPageTaskFullSchema struct {
	Data  []TaskFullSchema `json:"data"`
	Page  int              `json:"page"`
	Pages int              `json:"pages"`
	Size  int              `json:"size"`
	Total int              `json:"total"`
}

type StatusResponseSchema struct {
	Data StatusSchema `json:"data"`
}

type StatusSchema struct {
	CharactersOnline int               `json:"characters_online"`
	MaxLevel         int               `json:"max_level"`
	MaxSkillLevel    int               `json:"max_skill_level"`
	RateLimits       []RateLimitSchema `json:"rate_limits"`
	Season           *SeasonSchema     `json:"season,omitempty"`
	ServerTime       time.Time         `json:"server_time"`
	Version          string            `json:"version"`
}

type StatusSeasonRewardSchema struct {
	Code           string `json:"code"`
	Description    string `json:"description"`
	FirstOnly      *bool  `json:"first_only,omitempty"`
	MemberRequired *bool  `json:"member_required,omitempty"`
	Quantity       *int   `json:"quantity,omitempty"`
	RequiredPoints int    `json:"required_points"`
	Type           string `json:"type"`
}

type StorageEffectSchema struct {
	Code  string `json:"code"`
	Value int    `json:"value"`
}

type SubscribeRequestSchema struct {
	Plan      string `json:"plan"`
	Recurring *bool  `json:"recurring,omitempty"`
}

type SubscriptionResponseSchema struct {
	Data SubscriptionSchema `json:"data"`
}

type SubscriptionSchema struct {
	CancelledAt        *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	CurrentPeriodEnd   time.Time  `json:"current_period_end"`
	CurrentPeriodStart time.Time  `json:"current_period_start"`
	Plan               string     `json:"plan"`
	PurchaseSource     string     `json:"purchase_source"`
	Status             string     `json:"status"`
}

type TaskCancelledResponseSchema struct {
	Data TaskCancelledSchema `json:"data"`
}

type TaskCancelledSchema struct {
	Character CharacterSchema `json:"character"`
	Cooldown  CooldownSchema  `json:"cooldown"`
}

type TaskDataSchema struct {
	Character CharacterSchema `json:"character"`
	Cooldown  CooldownSchema  `json:"cooldown"`
	Task      TaskSchema      `json:"task"`
}

type TaskFullResponseSchema struct {
	Data TaskFullSchema `json:"data"`
}

type TaskFullSchema struct {
	Code        string        `json:"code"`
	Level       int           `json:"level"`
	MaxQuantity int           `json:"max_quantity"`
	MinQuantity int           `json:"min_quantity"`
	Rewards     RewardsSchema `json:"rewards"`
	Skill       *string       `json:"skill,omitempty"`
	Type        string        `json:"type"`
}

type TaskResponseSchema struct {
	Data TaskDataSchema `json:"data"`
}

type TaskSchema struct {
	Code    string        `json:"code"`
	Rewards RewardsSchema `json:"rewards"`
	Total   int           `json:"total"`
	Type    string        `json:"type"`
}

type TaskTradeDataSchema struct {
	Character CharacterSchema `json:"character"`
	Cooldown  CooldownSchema  `json:"cooldown"`
	Trade     TaskTradeSchema `json:"trade"`
}

type TaskTradeResponseSchema struct {
	Data TaskTradeDataSchema `json:"data"`
}

type TaskTradeSchema struct {
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
}

type TokenResponseSchema struct {
	Token string `json:"token"`
}

type TransitionSchema struct {
	Conditions *[]ConditionSchema `json:"conditions,omitempty"`
	Layer      string             `json:"layer"`
	MapId      int                `json:"map_id"`
	X          int                `json:"x"`
	Y          int                `json:"y"`
}

type UnequipSchema struct {
	Quantity *int   `json:"quantity,omitempty"`
	Slot     string `json:"slot"`
}

type UseItemResponseSchema struct {
	Data UseItemSchema `json:"data"`
}

type UseItemSchema struct {
	Character CharacterSchema `json:"character"`
	Cooldown  CooldownSchema  `json:"cooldown"`
	Item      ItemSchema      `json:"item"`
}
