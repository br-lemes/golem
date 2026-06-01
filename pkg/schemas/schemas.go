package schemas

import (
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	HTTPBasicScopes = "HTTPBasic.Scopes"
	JWTBearerScopes = "JWTBearer.Scopes"
)

const (
	AccountLeaderboardTypeAchievementsPoints AccountLeaderboardType = "achievements_points"
	AccountLeaderboardTypeGold               AccountLeaderboardType = "gold"
)

const (
	AccountStatusFounder     AccountStatus = "founder"
	AccountStatusGoldFounder AccountStatus = "gold_founder"
	AccountStatusStandard    AccountStatus = "standard"
	AccountStatusVipFounder  AccountStatus = "vip_founder"
)

const (
	AchievementTypeCombatDrop  AchievementType = "combat_drop"
	AchievementTypeCombatKill  AchievementType = "combat_kill"
	AchievementTypeCombatLevel AchievementType = "combat_level"
	AchievementTypeCrafting    AchievementType = "crafting"
	AchievementTypeGathering   AchievementType = "gathering"
	AchievementTypeNpcBuy      AchievementType = "npc_buy"
	AchievementTypeNpcSell     AchievementType = "npc_sell"
	AchievementTypeOther       AchievementType = "other"
	AchievementTypeRecycling   AchievementType = "recycling"
	AchievementTypeTask        AchievementType = "task"
	AchievementTypeUse         AchievementType = "use"
)

const (
	ActionTypeBuyBankExpansion ActionType = "buy_bank_expansion"
	ActionTypeBuyGe            ActionType = "buy_ge"
	ActionTypeBuyNpc           ActionType = "buy_npc"
	ActionTypeCancelGe         ActionType = "cancel_ge"
	ActionTypeChangeSkin       ActionType = "change_skin"
	ActionTypeClaimItem        ActionType = "claim_item"
	ActionTypeCrafting         ActionType = "crafting"
	ActionTypeCreateBuyOrderGe ActionType = "create_buy_order_ge"
	ActionTypeDeleteItem       ActionType = "delete_item"
	ActionTypeDepositGold      ActionType = "deposit_gold"
	ActionTypeDepositItem      ActionType = "deposit_item"
	ActionTypeEquip            ActionType = "equip"
	ActionTypeFight            ActionType = "fight"
	ActionTypeFillBuyOrderGe   ActionType = "fill_buy_order_ge"
	ActionTypeGathering        ActionType = "gathering"
	ActionTypeGiveGold         ActionType = "give_gold"
	ActionTypeGiveItem         ActionType = "give_item"
	ActionTypeMovement         ActionType = "movement"
	ActionTypeMultiFight       ActionType = "multi_fight"
	ActionTypeRecycling        ActionType = "recycling"
	ActionTypeRename           ActionType = "rename"
	ActionTypeRest             ActionType = "rest"
	ActionTypeSandboxGiveGold  ActionType = "sandbox_give_gold"
	ActionTypeSandboxGiveItem  ActionType = "sandbox_give_item"
	ActionTypeSandboxGiveXp    ActionType = "sandbox_give_xp"
	ActionTypeSellGe           ActionType = "sell_ge"
	ActionTypeSellNpc          ActionType = "sell_npc"
	ActionTypeTask             ActionType = "task"
	ActionTypeTransition       ActionType = "transition"
	ActionTypeUnequip          ActionType = "unequip"
	ActionTypeUse              ActionType = "use"
	ActionTypeWithdrawGold     ActionType = "withdraw_gold"
	ActionTypeWithdrawItem     ActionType = "withdraw_item"
)

const (
	CharacterLeaderboardTypeAlchemy         CharacterLeaderboardType = "alchemy"
	CharacterLeaderboardTypeCombat          CharacterLeaderboardType = "combat"
	CharacterLeaderboardTypeCooking         CharacterLeaderboardType = "cooking"
	CharacterLeaderboardTypeFishing         CharacterLeaderboardType = "fishing"
	CharacterLeaderboardTypeGearcrafting    CharacterLeaderboardType = "gearcrafting"
	CharacterLeaderboardTypeJewelrycrafting CharacterLeaderboardType = "jewelrycrafting"
	CharacterLeaderboardTypeMining          CharacterLeaderboardType = "mining"
	CharacterLeaderboardTypeWeaponcrafting  CharacterLeaderboardType = "weaponcrafting"
	CharacterLeaderboardTypeWoodcutting     CharacterLeaderboardType = "woodcutting"
)

const (
	Corrupted1 CharacterSkin = "corrupted1"
	Goblin1    CharacterSkin = "goblin1"
	Marauder1  CharacterSkin = "marauder1"
	Men1       CharacterSkin = "men1"
	Men2       CharacterSkin = "men2"
	Men3       CharacterSkin = "men3"
	Women1     CharacterSkin = "women1"
	Women2     CharacterSkin = "women2"
	Women3     CharacterSkin = "women3"
	Zombie1    CharacterSkin = "zombie1"
)

const (
	AchievementUnlocked ConditionOperator = "achievement_unlocked"
	Cost                ConditionOperator = "cost"
	Eq                  ConditionOperator = "eq"
	Gt                  ConditionOperator = "gt"
	HasItem             ConditionOperator = "has_item"
	Lt                  ConditionOperator = "lt"
	Ne                  ConditionOperator = "ne"
)

const (
	CraftSkillAlchemy         CraftSkill = "alchemy"
	CraftSkillCooking         CraftSkill = "cooking"
	CraftSkillGearcrafting    CraftSkill = "gearcrafting"
	CraftSkillJewelrycrafting CraftSkill = "jewelrycrafting"
	CraftSkillMining          CraftSkill = "mining"
	CraftSkillWeaponcrafting  CraftSkill = "weaponcrafting"
	CraftSkillWoodcutting     CraftSkill = "woodcutting"
)

const (
	EffectSubtypeBuff      EffectSubtype = "buff"
	EffectSubtypeDebuff    EffectSubtype = "debuff"
	EffectSubtypeGathering EffectSubtype = "gathering"
	EffectSubtypeGold      EffectSubtype = "gold"
	EffectSubtypeHeal      EffectSubtype = "heal"
	EffectSubtypeOther     EffectSubtype = "other"
	EffectSubtypeSpecial   EffectSubtype = "special"
	EffectSubtypeStat      EffectSubtype = "stat"
	EffectSubtypeTeleport  EffectSubtype = "teleport"
)

const (
	EffectTypeCombat     EffectType = "combat"
	EffectTypeConsumable EffectType = "consumable"
	EffectTypeEquipment  EffectType = "equipment"
)

const (
	Loss FightResult = "loss"
	Win  FightResult = "win"
)

const (
	Buy  GEOrderType = "buy"
	Sell GEOrderType = "sell"
)

const (
	GatheringSkillAlchemy     GatheringSkill = "alchemy"
	GatheringSkillFishing     GatheringSkill = "fishing"
	GatheringSkillMining      GatheringSkill = "mining"
	GatheringSkillWoodcutting GatheringSkill = "woodcutting"
)

const (
	Amulet    ItemSlot = "amulet"
	Artifact1 ItemSlot = "artifact1"
	Artifact2 ItemSlot = "artifact2"
	Artifact3 ItemSlot = "artifact3"
	Bag       ItemSlot = "bag"
	BodyArmor ItemSlot = "body_armor"
	Boots     ItemSlot = "boots"
	Helmet    ItemSlot = "helmet"
	LegArmor  ItemSlot = "leg_armor"
	Ring1     ItemSlot = "ring1"
	Ring2     ItemSlot = "ring2"
	Rune      ItemSlot = "rune"
	Shield    ItemSlot = "shield"
	Utility1  ItemSlot = "utility1"
	Utility2  ItemSlot = "utility2"
	Weapon    ItemSlot = "weapon"
)

const (
	ItemTypeAmulet     ItemType = "amulet"
	ItemTypeArtifact   ItemType = "artifact"
	ItemTypeBag        ItemType = "bag"
	ItemTypeBodyArmor  ItemType = "body_armor"
	ItemTypeBoots      ItemType = "boots"
	ItemTypeConsumable ItemType = "consumable"
	ItemTypeCurrency   ItemType = "currency"
	ItemTypeHelmet     ItemType = "helmet"
	ItemTypeLegArmor   ItemType = "leg_armor"
	ItemTypeResource   ItemType = "resource"
	ItemTypeRing       ItemType = "ring"
	ItemTypeRune       ItemType = "rune"
	ItemTypeShield     ItemType = "shield"
	ItemTypeUtility    ItemType = "utility"
	ItemTypeWeapon     ItemType = "weapon"
)

const (
	LogTypeAchievement         LogType = "achievement"
	LogTypeBuyBankExpansion    LogType = "buy_bank_expansion"
	LogTypeBuyGe               LogType = "buy_ge"
	LogTypeBuyNpc              LogType = "buy_npc"
	LogTypeCancelGe            LogType = "cancel_ge"
	LogTypeChangeSkin          LogType = "change_skin"
	LogTypeClaimItem           LogType = "claim_item"
	LogTypeCrafting            LogType = "crafting"
	LogTypeCreateBuyOrderGe    LogType = "create_buy_order_ge"
	LogTypeDeleteCharacter     LogType = "delete_character"
	LogTypeDeleteItem          LogType = "delete_item"
	LogTypeDepositGold         LogType = "deposit_gold"
	LogTypeDepositItem         LogType = "deposit_item"
	LogTypeEquip               LogType = "equip"
	LogTypeFight               LogType = "fight"
	LogTypeFillBuyOrderGe      LogType = "fill_buy_order_ge"
	LogTypeGathering           LogType = "gathering"
	LogTypeGiveGold            LogType = "give_gold"
	LogTypeGiveItem            LogType = "give_item"
	LogTypeMovement            LogType = "movement"
	LogTypeMultiFight          LogType = "multi_fight"
	LogTypeNewTask             LogType = "new_task"
	LogTypeReceiveGold         LogType = "receive_gold"
	LogTypeReceiveItem         LogType = "receive_item"
	LogTypeRecycling           LogType = "recycling"
	LogTypeRename              LogType = "rename"
	LogTypeRest                LogType = "rest"
	LogTypeSandboxGiveGold     LogType = "sandbox_give_gold"
	LogTypeSandboxGiveItem     LogType = "sandbox_give_item"
	LogTypeSandboxGiveXp       LogType = "sandbox_give_xp"
	LogTypeSandboxResetAccount LogType = "sandbox_reset_account"
	LogTypeSellGe              LogType = "sell_ge"
	LogTypeSellNpc             LogType = "sell_npc"
	LogTypeSpawn               LogType = "spawn"
	LogTypeTaskCancelled       LogType = "task_cancelled"
	LogTypeTaskCompleted       LogType = "task_completed"
	LogTypeTaskExchange        LogType = "task_exchange"
	LogTypeTaskTrade           LogType = "task_trade"
	LogTypeTransition          LogType = "transition"
	LogTypeUnequip             LogType = "unequip"
	LogTypeUse                 LogType = "use"
	LogTypeWithdrawGold        LogType = "withdraw_gold"
	LogTypeWithdrawItem        LogType = "withdraw_item"
)

const (
	MapAccessTypeBlocked       MapAccessType = "blocked"
	MapAccessTypeConditional   MapAccessType = "conditional"
	MapAccessTypeStandard      MapAccessType = "standard"
	MapAccessTypeTeleportation MapAccessType = "teleportation"
)

const (
	Bank          MapContentType = "bank"
	GrandExchange MapContentType = "grand_exchange"
	Monster       MapContentType = "monster"
	Npc           MapContentType = "npc"
	Resource      MapContentType = "resource"
	TasksMaster   MapContentType = "tasks_master"
	Workshop      MapContentType = "workshop"
)

const (
	Interior    MapLayer = "interior"
	Overworld   MapLayer = "overworld"
	Underground MapLayer = "underground"
)

const (
	Boss   MonsterType = "boss"
	Elite  MonsterType = "elite"
	Normal MonsterType = "normal"
)

const (
	Merchant NPCType = "merchant"
	Trader   NPCType = "trader"
)

const (
	PendingItemSourceAchievement   PendingItemSource = "achievement"
	PendingItemSourceAdmin         PendingItemSource = "admin"
	PendingItemSourceEvent         PendingItemSource = "event"
	PendingItemSourceGrandExchange PendingItemSource = "grand_exchange"
	PendingItemSourceOther         PendingItemSource = "other"
)

const (
	SkillAlchemy         Skill = "alchemy"
	SkillCooking         Skill = "cooking"
	SkillFishing         Skill = "fishing"
	SkillGearcrafting    Skill = "gearcrafting"
	SkillJewelrycrafting Skill = "jewelrycrafting"
	SkillMining          Skill = "mining"
	SkillWeaponcrafting  Skill = "weaponcrafting"
	SkillWoodcutting     Skill = "woodcutting"
)

const (
	Items    TaskType = "items"
	Monsters TaskType = "monsters"
)

type AccessSchema struct {
	Conditions *[]ConditionSchema `json:"conditions"`
	Type       MapAccessType      `json:"type"`
}

type AccountAchievementObjectiveSchema struct {
	Progress *int            `json:"progress,omitempty"`
	Target   *string         `json:"target"`
	Total    int             `json:"total"`
	Type     AchievementType `json:"type"`
}

type AccountAchievementSchema struct {
	Code        string                              `json:"code"`
	CompletedAt *time.Time                          `json:"completed_at"`
	Description string                              `json:"description"`
	Name        string                              `json:"name"`
	Objectives  []AccountAchievementObjectiveSchema `json:"objectives"`
	Points      int                                 `json:"points"`
	Rewards     AchievementRewardsSchema            `json:"rewards"`
}

type AccountDetails struct {
	AchievementsPoints int           `json:"achievements_points"`
	Badges             *[]string     `json:"badges,omitempty"`
	BanReason          *string       `json:"ban_reason,omitempty"`
	Banned             bool          `json:"banned"`
	Member             bool          `json:"member"`
	Skins              []string      `json:"skins"`
	Status             AccountStatus `json:"status"`
	Username           string        `json:"username"`
}

type AccountDetailsSchema struct {
	Data AccountDetails `json:"data"`
}

type AccountLeaderboardSchema struct {
	Account            string        `json:"account"`
	AchievementsPoints int           `json:"achievements_points"`
	Gold               int           `json:"gold"`
	Position           int           `json:"position"`
	Status             AccountStatus `json:"status"`
}

type AccountLeaderboardType string

type AccountStatus string

type AchievementObjectiveSchema struct {
	Target *string         `json:"target"`
	Total  int             `json:"total"`
	Type   AchievementType `json:"type"`
}

type AchievementResponseSchema struct {
	Data AchievementSchema `json:"data"`
}

type AchievementRewardsSchema struct {
	Gold  *int                `json:"gold,omitempty"`
	Items *[]RewardItemSchema `json:"items"`
}

type AchievementSchema struct {
	Code        string                       `json:"code"`
	Description string                       `json:"description"`
	Name        string                       `json:"name"`
	Objectives  []AchievementObjectiveSchema `json:"objectives"`
	Points      int                          `json:"points"`
	Rewards     AchievementRewardsSchema     `json:"rewards"`
}

type AchievementType string

type ActionType string

type ActiveCharacterSchema struct {
	Account string        `json:"account"`
	Layer   MapLayer      `json:"layer"`
	MapId   int           `json:"map_id"`
	Name    string        `json:"name"`
	Skin    CharacterSkin `json:"skin"`
	X       int           `json:"x"`
	Y       int           `json:"y"`
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
	Name string        `json:"name"`
	Skin CharacterSkin `json:"skin"`
}

type BadgeConditionSchema struct {
	Code     string `json:"code"`
	Quantity *int   `json:"quantity"`
}

type BadgeResponseSchema struct {
	Data BadgeSchema `json:"data"`
}

type BadgeSchema struct {
	Code        string                 `json:"code"`
	Conditions  []BadgeConditionSchema `json:"conditions"`
	Description string                 `json:"description"`
	Season      *int                   `json:"season"`
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

type ChangePassword struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ChangeSkinCharacterDataSchema struct {
	Character CharacterSchema `json:"character"`
	Cooldown  CooldownSchema  `json:"cooldown"`
	Skin      string          `json:"skin"`
}

type ChangeSkinCharacterSchema struct {
	Skin CharacterSkin `json:"skin"`
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
	Result     FightResult                       `json:"result"`
	Turns      int                               `json:"turns"`
}

type CharacterLeaderboardSchema struct {
	Account                string        `json:"account"`
	AlchemyLevel           int           `json:"alchemy_level"`
	AlchemyTotalXp         int           `json:"alchemy_total_xp"`
	CookingLevel           int           `json:"cooking_level"`
	CookingTotalXp         int           `json:"cooking_total_xp"`
	FishingLevel           int           `json:"fishing_level"`
	FishingTotalXp         int           `json:"fishing_total_xp"`
	GearcraftingLevel      int           `json:"gearcrafting_level"`
	GearcraftingTotalXp    int           `json:"gearcrafting_total_xp"`
	Gold                   int           `json:"gold"`
	JewelrycraftingLevel   int           `json:"jewelrycrafting_level"`
	JewelrycraftingTotalXp int           `json:"jewelrycrafting_total_xp"`
	Level                  int           `json:"level"`
	MiningLevel            int           `json:"mining_level"`
	MiningTotalXp          int           `json:"mining_total_xp"`
	Name                   string        `json:"name"`
	Position               int           `json:"position"`
	Skin                   string        `json:"skin"`
	Status                 AccountStatus `json:"status"`
	TotalXp                int           `json:"total_xp"`
	WeaponcraftingLevel    int           `json:"weaponcrafting_level"`
	WeaponcraftingTotalXp  int           `json:"weaponcrafting_total_xp"`
	WoodcuttingLevel       int           `json:"woodcutting_level"`
	WoodcuttingTotalXp     int           `json:"woodcutting_total_xp"`
}

type CharacterLeaderboardType string

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
	Inventory            *[]InventorySlot       `json:"inventory,omitempty"`
	InventoryMaxItems    int                    `json:"inventory_max_items"`
	JewelrycraftingLevel int                    `json:"jewelrycrafting_level"`
	JewelrycraftingMaxXp int                    `json:"jewelrycrafting_max_xp"`
	JewelrycraftingXp    int                    `json:"jewelrycrafting_xp"`
	Layer                MapLayer               `json:"layer"`
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
	Skin                 CharacterSkin          `json:"skin"`
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

type CharacterSkin string

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

type ConditionOperator string

type ConditionSchema struct {
	Code     string            `json:"code"`
	Operator ConditionOperator `json:"operator"`
	Value    int               `json:"value"`
}

type CooldownSchema struct {
	Expiration       time.Time  `json:"expiration"`
	Reason           ActionType `json:"reason"`
	RemainingSeconds int        `json:"remaining_seconds"`
	StartedAt        time.Time  `json:"started_at"`
	TotalSeconds     int        `json:"total_seconds"`
}

type CraftSchema struct {
	Items    *[]SimpleItemSchema `json:"items,omitempty"`
	Level    *int                `json:"level,omitempty"`
	Quantity *int                `json:"quantity,omitempty"`
	Skill    *CraftSkill         `json:"skill,omitempty"`
}

type CraftSkill string

type CraftingSchema struct {
	Code     string `json:"code"`
	Quantity *int   `json:"quantity,omitempty"`
}

type DataPageAccountAchievementSchema struct {
	Data  []AccountAchievementSchema `json:"data"`
	Page  *int                       `json:"page"`
	Pages *int                       `json:"pages"`
	Size  *int                       `json:"size"`
	Total *int                       `json:"total"`
}

type DataPageAccountLeaderboardSchema struct {
	Data  []AccountLeaderboardSchema `json:"data"`
	Page  *int                       `json:"page"`
	Pages *int                       `json:"pages"`
	Size  *int                       `json:"size"`
	Total *int                       `json:"total"`
}

type DataPageActiveCharacterSchema struct {
	Data  []ActiveCharacterSchema `json:"data"`
	Page  *int                    `json:"page"`
	Pages *int                    `json:"pages"`
	Size  *int                    `json:"size"`
	Total *int                    `json:"total"`
}

type DataPageCharacterLeaderboardSchema struct {
	Data  []CharacterLeaderboardSchema `json:"data"`
	Page  *int                         `json:"page"`
	Pages *int                         `json:"pages"`
	Size  *int                         `json:"size"`
	Total *int                         `json:"total"`
}

type DataPageGEOrderSchema struct {
	Data  []GEOrderSchema `json:"data"`
	Page  *int            `json:"page"`
	Pages *int            `json:"pages"`
	Size  *int            `json:"size"`
	Total *int            `json:"total"`
}

type DataPageGeOrderHistorySchema struct {
	Data  []GeOrderHistorySchema `json:"data"`
	Page  *int                   `json:"page"`
	Pages *int                   `json:"pages"`
	Size  *int                   `json:"size"`
	Total *int                   `json:"total"`
}

type DataPageLogSchema struct {
	Data  []LogSchema `json:"data"`
	Page  *int        `json:"page"`
	Pages *int        `json:"pages"`
	Size  *int        `json:"size"`
	Total *int        `json:"total"`
}

type DataPagePendingItemSchema struct {
	Data  []PendingItemSchema `json:"data"`
	Page  *int                `json:"page"`
	Pages *int                `json:"pages"`
	Size  *int                `json:"size"`
	Total *int                `json:"total"`
}

type DataPageSimpleItemSchema struct {
	Data  []SimpleItemSchema `json:"data"`
	Page  *int               `json:"page"`
	Pages *int               `json:"pages"`
	Size  *int               `json:"size"`
	Total *int               `json:"total"`
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
	Code        string        `json:"code"`
	Description string        `json:"description"`
	Name        string        `json:"name"`
	Subtype     EffectSubtype `json:"subtype"`
	Type        EffectType    `json:"type"`
}

type EffectSubtype string

type EffectType string

type EquipRequestSchema struct {
	Character CharacterSchema `json:"character"`
	Cooldown  CooldownSchema  `json:"cooldown"`
	Item      ItemSchema      `json:"item"`
	Slot      ItemSlot        `json:"slot"`
}

type EquipSchema struct {
	Code     string   `json:"code"`
	Quantity *int     `json:"quantity,omitempty"`
	Slot     ItemSlot `json:"slot"`
}

type EquipmentResponseSchema struct {
	Data EquipRequestSchema `json:"data"`
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
	Code string         `json:"code"`
	Type MapContentType `json:"type"`
}

type EventMapSchema struct {
	Layer string `json:"layer"`
	MapId int    `json:"map_id"`
	Skin  string `json:"skin"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
}

type EventSchema struct {
	Code     string             `json:"code"`
	Content  EventContentSchema `json:"content"`
	Duration int                `json:"duration"`
	Maps     []EventMapSchema   `json:"maps"`
	Name     string             `json:"name"`
	Rate     int                `json:"rate"`
}

type FakeCharacterSchema struct {
	AmuletSlot           *string `json:"amulet_slot"`
	Artifact1Slot        *string `json:"artifact1_slot"`
	Artifact2Slot        *string `json:"artifact2_slot"`
	Artifact3Slot        *string `json:"artifact3_slot"`
	BodyArmorSlot        *string `json:"body_armor_slot"`
	BootsSlot            *string `json:"boots_slot"`
	HelmetSlot           *string `json:"helmet_slot"`
	LegArmorSlot         *string `json:"leg_armor_slot"`
	Level                int     `json:"level"`
	Ring1Slot            *string `json:"ring1_slot"`
	Ring2Slot            *string `json:"ring2_slot"`
	RuneSlot             *string `json:"rune_slot"`
	ShieldSlot           *string `json:"shield_slot"`
	Utility1Slot         *string `json:"utility1_slot"`
	Utility1SlotQuantity *int    `json:"utility1_slot_quantity,omitempty"`
	Utility2Slot         *string `json:"utility2_slot"`
	Utility2SlotQuantity *int    `json:"utility2_slot_quantity,omitempty"`
	WeaponSlot           *string `json:"weapon_slot"`
}

type FightRequestSchema struct {
	Participants *[]string `json:"participants,omitempty"`
}

type FightResult string

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

type GEOrderCreationrSchema struct {
	Code     string `json:"code"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type GEOrderResponseSchema struct {
	Data GEOrderSchema `json:"data"`
}

type GEOrderSchema struct {
	Account   *string     `json:"account"`
	Code      string      `json:"code"`
	CreatedAt time.Time   `json:"created_at"`
	Id        string      `json:"id"`
	Price     int         `json:"price"`
	Quantity  int         `json:"quantity"`
	Type      GEOrderType `json:"type"`
}

type GEOrderTransactionSchema struct {
	Character CharacterSchema      `json:"character"`
	Cooldown  CooldownSchema       `json:"cooldown"`
	Order     GEOrderCreatedSchema `json:"order"`
}

type GEOrderType string

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

type GatheringSkill string

type GeOrderHistorySchema struct {
	Buyer    string    `json:"buyer"`
	Code     string    `json:"code"`
	OrderId  string    `json:"order_id"`
	Price    int       `json:"price"`
	Quantity int       `json:"quantity"`
	Seller   string    `json:"seller"`
	SoldAt   time.Time `json:"sold_at"`
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

type InventorySlot struct {
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
	Subtype     string                `json:"subtype"`
	Tradeable   bool                  `json:"tradeable"`
	Type        string                `json:"type"`
}

type ItemSlot string

type ItemType string

type LogSchema struct {
	Account            string      `json:"account"`
	Character          string      `json:"character"`
	Content            interface{} `json:"content"`
	Cooldown           int         `json:"cooldown"`
	CooldownExpiration *time.Time  `json:"cooldown_expiration"`
	CreatedAt          time.Time   `json:"created_at"`
	Description        string      `json:"description"`
	Type               LogType     `json:"type"`
}

type LogType string

type MapAccessType string

type MapContentSchema struct {
	Code string         `json:"code"`
	Type MapContentType `json:"type"`
}

type MapContentType string

type MapLayer string

type MapResponseSchema struct {
	Data MapSchema `json:"data"`
}

type MapSchema struct {
	Access       AccessSchema      `json:"access"`
	Interactions InteractionSchema `json:"interactions"`
	Layer        MapLayer          `json:"layer"`
	MapId        int               `json:"map_id"`
	Name         string            `json:"name"`
	Skin         string            `json:"skin"`
	X            int               `json:"x"`
	Y            int               `json:"y"`
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
	Type           MonsterType           `json:"type"`
}

type MonsterType string

type MyAccountDetails struct {
	AchievementsPoints int                 `json:"achievements_points"`
	Badges             *[]string           `json:"badges,omitempty"`
	BanReason          *string             `json:"ban_reason,omitempty"`
	Banned             bool                `json:"banned"`
	Email              openapi_types.Email `json:"email"`
	EventToken         int                 `json:"event_token"`
	Gems               int                 `json:"gems"`
	Member             bool                `json:"member"`
	MemberExpiration   *time.Time          `json:"member_expiration"`
	Skins              []string            `json:"skins"`
	Status             AccountStatus       `json:"status"`
	Username           string              `json:"username"`
}

type MyAccountDetailsSchema struct {
	Data MyAccountDetails `json:"data"`
}

type MyCharactersListSchema struct {
	Data []CharacterSchema `json:"data"`
}

type NPCItem struct {
	BuyPrice  *int   `json:"buy_price"`
	Code      string `json:"code"`
	Currency  string `json:"currency"`
	Npc       string `json:"npc"`
	SellPrice *int   `json:"sell_price"`
}

type NPCResponseSchema struct {
	Data NPCSchema `json:"data"`
}

type NPCSchema struct {
	Code        string           `json:"code"`
	Description string           `json:"description"`
	Items       *[]SimpleNPCItem `json:"items,omitempty"`
	Name        string           `json:"name"`
	Type        NPCType          `json:"type"`
}

type NPCType string

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
	ClaimedAt   *time.Time          `json:"claimed_at"`
	CreatedAt   time.Time           `json:"created_at"`
	Description string              `json:"description"`
	Gold        *int                `json:"gold,omitempty"`
	Id          string              `json:"id"`
	Items       *[]SimpleItemSchema `json:"items,omitempty"`
	Source      PendingItemSource   `json:"source"`
	SourceId    *string             `json:"source_id"`
}

type PendingItemSource string

type RateLimitSchema struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type RecyclingDataSchema struct {
	Character CharacterSchema      `json:"character"`
	Cooldown  CooldownSchema       `json:"cooldown"`
	Details   RecyclingItemsSchema `json:"details"`
}

type RecyclingItemsSchema struct {
	Items []DropSchema `json:"items"`
}

type RecyclingResponseSchema struct {
	Data RecyclingDataSchema `json:"data"`
}

type RecyclingSchema struct {
	Code     string `json:"code"`
	Quantity *int   `json:"quantity,omitempty"`
}

type ResourceResponseSchema struct {
	Data ResourceSchema `json:"data"`
}

type ResourceSchema struct {
	Code  string           `json:"code"`
	Drops []DropRateSchema `json:"drops"`
	Level int              `json:"level"`
	Name  string           `json:"name"`
	Skill GatheringSkill   `json:"skill"`
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

type SeasonBadgeSchema struct {
	Code           string `json:"code"`
	Description    string `json:"description"`
	RequiredPoints int    `json:"required_points"`
}

type SeasonSchema struct {
	Badges    []SeasonBadgeSchema `json:"badges"`
	Name      *string             `json:"name,omitempty"`
	Number    *int                `json:"number,omitempty"`
	Skins     []SeasonSkinSchema  `json:"skins"`
	StartDate *time.Time          `json:"start_date,omitempty"`
}

type SeasonSkinSchema struct {
	Code           string `json:"code"`
	Description    string `json:"description"`
	RequiredPoints int    `json:"required_points"`
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

type SimpleNPCItem struct {
	BuyPrice  *int   `json:"buy_price"`
	Code      string `json:"code"`
	Currency  string `json:"currency"`
	SellPrice *int   `json:"sell_price"`
}

type Skill string

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

type SpawnEventRequest struct {
	Code string `json:"code"`
}

type StaticDataPageAchievementSchema struct {
	Data  []AchievementSchema `json:"data"`
	Page  *int                `json:"page"`
	Pages *int                `json:"pages"`
	Size  *int                `json:"size"`
	Total *int                `json:"total"`
}

type StaticDataPageActiveEventSchema struct {
	Data  []ActiveEventSchema `json:"data"`
	Page  *int                `json:"page"`
	Pages *int                `json:"pages"`
	Size  *int                `json:"size"`
	Total *int                `json:"total"`
}

type StaticDataPageBadgeSchema struct {
	Data  []BadgeSchema `json:"data"`
	Page  *int          `json:"page"`
	Pages *int          `json:"pages"`
	Size  *int          `json:"size"`
	Total *int          `json:"total"`
}

type StaticDataPageDropRateSchema struct {
	Data  []DropRateSchema `json:"data"`
	Page  *int             `json:"page"`
	Pages *int             `json:"pages"`
	Size  *int             `json:"size"`
	Total *int             `json:"total"`
}

type StaticDataPageEffectSchema struct {
	Data  []EffectSchema `json:"data"`
	Page  *int           `json:"page"`
	Pages *int           `json:"pages"`
	Size  *int           `json:"size"`
	Total *int           `json:"total"`
}

type StaticDataPageEventSchema struct {
	Data  []EventSchema `json:"data"`
	Page  *int          `json:"page"`
	Pages *int          `json:"pages"`
	Size  *int          `json:"size"`
	Total *int          `json:"total"`
}

type StaticDataPageItemSchema struct {
	Data  []ItemSchema `json:"data"`
	Page  *int         `json:"page"`
	Pages *int         `json:"pages"`
	Size  *int         `json:"size"`
	Total *int         `json:"total"`
}

type StaticDataPageMapSchema struct {
	Data  []MapSchema `json:"data"`
	Page  *int        `json:"page"`
	Pages *int        `json:"pages"`
	Size  *int        `json:"size"`
	Total *int        `json:"total"`
}

type StaticDataPageMonsterSchema struct {
	Data  []MonsterSchema `json:"data"`
	Page  *int            `json:"page"`
	Pages *int            `json:"pages"`
	Size  *int            `json:"size"`
	Total *int            `json:"total"`
}

type StaticDataPageNPCItem struct {
	Data  []NPCItem `json:"data"`
	Page  *int      `json:"page"`
	Pages *int      `json:"pages"`
	Size  *int      `json:"size"`
	Total *int      `json:"total"`
}

type StaticDataPageNPCSchema struct {
	Data  []NPCSchema `json:"data"`
	Page  *int        `json:"page"`
	Pages *int        `json:"pages"`
	Size  *int        `json:"size"`
	Total *int        `json:"total"`
}

type StaticDataPageResourceSchema struct {
	Data  []ResourceSchema `json:"data"`
	Page  *int             `json:"page"`
	Pages *int             `json:"pages"`
	Size  *int             `json:"size"`
	Total *int             `json:"total"`
}

type StaticDataPageTaskFullSchema struct {
	Data  []TaskFullSchema `json:"data"`
	Page  *int             `json:"page"`
	Pages *int             `json:"pages"`
	Size  *int             `json:"size"`
	Total *int             `json:"total"`
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

type StorageEffectSchema struct {
	Code  string `json:"code"`
	Value int    `json:"value"`
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
	Skill       *string       `json:"skill"`
	Type        TaskType      `json:"type"`
}

type TaskResponseSchema struct {
	Data TaskDataSchema `json:"data"`
}

type TaskSchema struct {
	Code    string        `json:"code"`
	Rewards RewardsSchema `json:"rewards"`
	Total   int           `json:"total"`
	Type    TaskType      `json:"type"`
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

type TaskType string

type TokenResponseSchema struct {
	Token string `json:"token"`
}

type TransitionSchema struct {
	Conditions *[]ConditionSchema `json:"conditions"`
	Layer      MapLayer           `json:"layer"`
	MapId      int                `json:"map_id"`
	X          int                `json:"x"`
	Y          int                `json:"y"`
}

type UnequipSchema struct {
	Quantity *int     `json:"quantity,omitempty"`
	Slot     ItemSlot `json:"slot"`
}

type UseItemResponseSchema struct {
	Data UseItemSchema `json:"data"`
}

type UseItemSchema struct {
	Character CharacterSchema `json:"character"`
	Cooldown  CooldownSchema  `json:"cooldown"`
	Item      ItemSchema      `json:"item"`
}

type GetAccountAchievementsAccountsAccountAchievementsGetParams struct {
	Type      *AchievementType `form:"type,omitempty" json:"type,omitempty"`
	Completed *bool            `form:"completed,omitempty" json:"completed,omitempty"`
	Page      *int             `form:"page,omitempty" json:"page,omitempty"`
	Size      *int             `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllAchievementsAchievementsGetParams struct {
	Type *AchievementType `form:"type,omitempty" json:"type,omitempty"`
	Page *int             `form:"page,omitempty" json:"page,omitempty"`
	Size *int             `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllBadgesBadgesGetParams struct {
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	Size *int `form:"size,omitempty" json:"size,omitempty"`
}

type GetActiveCharactersCharactersActiveGetParams struct {
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	Size *int `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllEffectsEffectsGetParams struct {
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	Size *int `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllEventsEventsGetParams struct {
	Type *MapContentType `form:"type,omitempty" json:"type,omitempty"`
	Page *int            `form:"page,omitempty" json:"page,omitempty"`
	Size *int            `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllActiveEventsEventsActiveGetParams struct {
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	Size *int `form:"size,omitempty" json:"size,omitempty"`
}

type GetGeHistoryGrandexchangeHistoryCodeGetParams struct {
	Account *string `form:"account,omitempty" json:"account,omitempty"`
	Page    *int    `form:"page,omitempty" json:"page,omitempty"`
	Size    *int    `form:"size,omitempty" json:"size,omitempty"`
}

type GetGeOrdersGrandexchangeOrdersGetParams struct {
	Code    *string      `form:"code,omitempty" json:"code,omitempty"`
	Account *string      `form:"account,omitempty" json:"account,omitempty"`
	Type    *GEOrderType `form:"type,omitempty" json:"type,omitempty"`
	Page    *int         `form:"page,omitempty" json:"page,omitempty"`
	Size    *int         `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllItemsItemsGetParams struct {
	Name          *string     `form:"name,omitempty" json:"name,omitempty"`
	MinLevel      *int        `form:"min_level,omitempty" json:"min_level,omitempty"`
	MaxLevel      *int        `form:"max_level,omitempty" json:"max_level,omitempty"`
	Type          *ItemType   `form:"type,omitempty" json:"type,omitempty"`
	CraftSkill    *CraftSkill `form:"craft_skill,omitempty" json:"craft_skill,omitempty"`
	CraftMaterial *string     `form:"craft_material,omitempty" json:"craft_material,omitempty"`
	Page          *int        `form:"page,omitempty" json:"page,omitempty"`
	Size          *int        `form:"size,omitempty" json:"size,omitempty"`
}

type GetAccountsLeaderboardLeaderboardAccountsGetParams struct {
	Sort *AccountLeaderboardType `form:"sort,omitempty" json:"sort,omitempty"`
	Name *string                 `form:"name,omitempty" json:"name,omitempty"`
	Page *int                    `form:"page,omitempty" json:"page,omitempty"`
	Size *int                    `form:"size,omitempty" json:"size,omitempty"`
}

type GetCharactersLeaderboardLeaderboardCharactersGetParams struct {
	Sort *CharacterLeaderboardType `form:"sort,omitempty" json:"sort,omitempty"`
	Name *string                   `form:"name,omitempty" json:"name,omitempty"`
	Page *int                      `form:"page,omitempty" json:"page,omitempty"`
	Size *int                      `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllMapsMapsGetParams struct {
	Layer           *MapLayer       `form:"layer,omitempty" json:"layer,omitempty"`
	ContentType     *MapContentType `form:"content_type,omitempty" json:"content_type,omitempty"`
	ContentCode     *string         `form:"content_code,omitempty" json:"content_code,omitempty"`
	HideBlockedMaps *bool           `form:"hide_blocked_maps,omitempty" json:"hide_blocked_maps,omitempty"`
	Page            *int            `form:"page,omitempty" json:"page,omitempty"`
	Size            *int            `form:"size,omitempty" json:"size,omitempty"`
}

type GetLayerMapsMapsLayerGetParams struct {
	ContentType     *MapContentType `form:"content_type,omitempty" json:"content_type,omitempty"`
	ContentCode     *string         `form:"content_code,omitempty" json:"content_code,omitempty"`
	HideBlockedMaps *bool           `form:"hide_blocked_maps,omitempty" json:"hide_blocked_maps,omitempty"`
	Page            *int            `form:"page,omitempty" json:"page,omitempty"`
	Size            *int            `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllMonstersMonstersGetParams struct {
	Name     *string `form:"name,omitempty" json:"name,omitempty"`
	MinLevel *int    `form:"min_level,omitempty" json:"min_level,omitempty"`
	MaxLevel *int    `form:"max_level,omitempty" json:"max_level,omitempty"`
	Drop     *string `form:"drop,omitempty" json:"drop,omitempty"`
	Page     *int    `form:"page,omitempty" json:"page,omitempty"`
	Size     *int    `form:"size,omitempty" json:"size,omitempty"`
}

type GetBankItemsMyBankItemsGetParams struct {
	ItemCode *string `form:"item_code,omitempty" json:"item_code,omitempty"`
	Page     *int    `form:"page,omitempty" json:"page,omitempty"`
	Size     *int    `form:"size,omitempty" json:"size,omitempty"`
}

type GetGeHistoryMyGrandexchangeHistoryGetParams struct {
	Id   *string `form:"id,omitempty" json:"id,omitempty"`
	Code *string `form:"code,omitempty" json:"code,omitempty"`
	Page *int    `form:"page,omitempty" json:"page,omitempty"`
	Size *int    `form:"size,omitempty" json:"size,omitempty"`
}

type GetGeOrdersMyGrandexchangeOrdersGetParams struct {
	Code *string      `form:"code,omitempty" json:"code,omitempty"`
	Type *GEOrderType `form:"type,omitempty" json:"type,omitempty"`
	Page *int         `form:"page,omitempty" json:"page,omitempty"`
	Size *int         `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllCharactersLogsMyLogsGetParams struct {
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	Size *int `form:"size,omitempty" json:"size,omitempty"`
}

type GetCharacterLogsMyLogsNameGetParams struct {
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	Size *int `form:"size,omitempty" json:"size,omitempty"`
}

type GetPendingItemsMyPendingItemsGetParams struct {
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	Size *int `form:"size,omitempty" json:"size,omitempty"`
}

type ActionDepositBankItemMyNameActionBankDepositItemPostJSONBody = []SimpleItemSchema

type ActionWithdrawBankItemMyNameActionBankWithdrawItemPostJSONBody = []SimpleItemSchema

type ActionFightMyNameActionFightPostJSONBody = FightRequestSchema

type GetAllNpcsNpcsDetailsGetParams struct {
	Name     *string  `form:"name,omitempty" json:"name,omitempty"`
	Type     *NPCType `form:"type,omitempty" json:"type,omitempty"`
	Currency *string  `form:"currency,omitempty" json:"currency,omitempty"`
	Item     *string  `form:"item,omitempty" json:"item,omitempty"`
	Page     *int     `form:"page,omitempty" json:"page,omitempty"`
	Size     *int     `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllNpcsItemsNpcsItemsGetParams struct {
	Code     *string `form:"code,omitempty" json:"code,omitempty"`
	Npc      *string `form:"npc,omitempty" json:"npc,omitempty"`
	Currency *string `form:"currency,omitempty" json:"currency,omitempty"`
	Page     *int    `form:"page,omitempty" json:"page,omitempty"`
	Size     *int    `form:"size,omitempty" json:"size,omitempty"`
}

type GetNpcItemsNpcsItemsCodeGetParams struct {
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	Size *int `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllResourcesResourcesGetParams struct {
	MinLevel *int            `form:"min_level,omitempty" json:"min_level,omitempty"`
	MaxLevel *int            `form:"max_level,omitempty" json:"max_level,omitempty"`
	Skill    *GatheringSkill `form:"skill,omitempty" json:"skill,omitempty"`
	Drop     *string         `form:"drop,omitempty" json:"drop,omitempty"`
	Page     *int            `form:"page,omitempty" json:"page,omitempty"`
	Size     *int            `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllTasksTasksListGetParams struct {
	MinLevel *int      `form:"min_level,omitempty" json:"min_level,omitempty"`
	MaxLevel *int      `form:"max_level,omitempty" json:"max_level,omitempty"`
	Skill    *Skill    `form:"skill,omitempty" json:"skill,omitempty"`
	Type     *TaskType `form:"type,omitempty" json:"type,omitempty"`
	Page     *int      `form:"page,omitempty" json:"page,omitempty"`
	Size     *int      `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllTasksRewardsTasksRewardsGetParams struct {
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	Size *int `form:"size,omitempty" json:"size,omitempty"`
}

type CreateAccountAccountsCreatePostJSONRequestBody = AddAccountSchema

type ForgotPasswordAccountsForgotPasswordPostJSONRequestBody = PasswordResetRequestSchema

type ResetPasswordAccountsResetPasswordPostJSONRequestBody = PasswordResetConfirmSchema

type CreateCharacterCharactersCreatePostJSONRequestBody = AddCharacterSchema

type DeleteCharacterCharactersDeletePostJSONRequestBody = DeleteCharacterSchema

type SpawnEventEventsSpawnPostJSONRequestBody = SpawnEventRequest

type ChangePasswordMyChangePasswordPostJSONRequestBody = ChangePassword

type ActionDepositBankGoldMyNameActionBankDepositGoldPostJSONRequestBody = DepositWithdrawGoldSchema

type ActionDepositBankItemMyNameActionBankDepositItemPostJSONRequestBody = ActionDepositBankItemMyNameActionBankDepositItemPostJSONBody

type ActionWithdrawBankGoldMyNameActionBankWithdrawGoldPostJSONRequestBody = DepositWithdrawGoldSchema

type ActionWithdrawBankItemMyNameActionBankWithdrawItemPostJSONRequestBody = ActionWithdrawBankItemMyNameActionBankWithdrawItemPostJSONBody

type ActionChangeSkinMyNameActionChangeSkinPostJSONRequestBody = ChangeSkinCharacterSchema

type ActionCraftingMyNameActionCraftingPostJSONRequestBody = CraftingSchema

type ActionDeleteItemMyNameActionDeletePostJSONRequestBody = SimpleItemSchema

type ActionEquipItemMyNameActionEquipPostJSONRequestBody = EquipSchema

type ActionFightMyNameActionFightPostJSONRequestBody = ActionFightMyNameActionFightPostJSONBody

type ActionGiveGoldMyNameActionGiveGoldPostJSONRequestBody = GiveGoldSchema

type ActionGiveItemsMyNameActionGiveItemPostJSONRequestBody = GiveItemsSchema

type ActionGeBuyItemMyNameActionGrandexchangeBuyPostJSONRequestBody = GEBuyOrderSchema

type ActionGeCancelOrderMyNameActionGrandexchangeCancelPostJSONRequestBody = GECancelOrderSchema

type ActionGeCreateBuyOrderMyNameActionGrandexchangeCreateBuyOrderPostJSONRequestBody = GEBuyOrderCreationSchema

type ActionGeCreateSellOrderMyNameActionGrandexchangeCreateSellOrderPostJSONRequestBody = GEOrderCreationrSchema

type ActionGeFillMyNameActionGrandexchangeFillPostJSONRequestBody = GEFillBuyOrderSchema

type ActionMoveMyNameActionMovePostJSONRequestBody = DestinationSchema

type ActionNpcBuyItemMyNameActionNpcBuyPostJSONRequestBody = NpcMerchantBuySchema

type ActionNpcSellItemMyNameActionNpcSellPostJSONRequestBody = NpcMerchantBuySchema

type ActionRecyclingMyNameActionRecyclingPostJSONRequestBody = RecyclingSchema

type ActionTaskTradeMyNameActionTaskTradePostJSONRequestBody = SimpleItemSchema

type ActionUnequipItemMyNameActionUnequipPostJSONRequestBody = UnequipSchema

type ActionUseItemMyNameActionUsePostJSONRequestBody = SimpleItemSchema

type FightSimulationSimulationFightSimulationPostJSONRequestBody = CombatSimulationRequestSchema
