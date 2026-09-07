package schemas

type GetAccountAchievementsAccountsAccountAchievementsGetParams struct {
	Type      *string `form:"type,omitempty" json:"type,omitempty"`
	Completed *bool   `form:"completed,omitempty" json:"completed,omitempty"`
	Page      *int    `form:"page,omitempty" json:"page,omitempty"`
	Size      *int    `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllAchievementsAchievementsGetParams struct {
	Type *string `form:"type,omitempty" json:"type,omitempty"`
	Page *int    `form:"page,omitempty" json:"page,omitempty"`
	Size *int    `form:"size,omitempty" json:"size,omitempty"`
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
	Type *string `form:"type,omitempty" json:"type,omitempty"`
	Page *int    `form:"page,omitempty" json:"page,omitempty"`
	Size *int    `form:"size,omitempty" json:"size,omitempty"`
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
	Code     *string `form:"code,omitempty" json:"code,omitempty"`
	Account  *string `form:"account,omitempty" json:"account,omitempty"`
	Type     *string `form:"type,omitempty" json:"type,omitempty"`
	ItemType *string `form:"item_type,omitempty" json:"item_type,omitempty"`
	Page     *int    `form:"page,omitempty" json:"page,omitempty"`
	Size     *int    `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllItemsItemsGetParams struct {
	Name          *string `form:"name,omitempty" json:"name,omitempty"`
	MinLevel      *int    `form:"min_level,omitempty" json:"min_level,omitempty"`
	MaxLevel      *int    `form:"max_level,omitempty" json:"max_level,omitempty"`
	Type          *string `form:"type,omitempty" json:"type,omitempty"`
	CraftSkill    *string `form:"craft_skill,omitempty" json:"craft_skill,omitempty"`
	CraftMaterial *string `form:"craft_material,omitempty" json:"craft_material,omitempty"`
	Page          *int    `form:"page,omitempty" json:"page,omitempty"`
	Size          *int    `form:"size,omitempty" json:"size,omitempty"`
}

type GetAccountsLeaderboardLeaderboardAccountsGetParams struct {
	Sort *string `form:"sort,omitempty" json:"sort,omitempty"`
	Name *string `form:"name,omitempty" json:"name,omitempty"`
	Page *int    `form:"page,omitempty" json:"page,omitempty"`
	Size *int    `form:"size,omitempty" json:"size,omitempty"`
}

type GetCharactersLeaderboardLeaderboardCharactersGetParams struct {
	Sort *string `form:"sort,omitempty" json:"sort,omitempty"`
	Name *string `form:"name,omitempty" json:"name,omitempty"`
	Page *int    `form:"page,omitempty" json:"page,omitempty"`
	Size *int    `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllMapsMapsGetParams struct {
	Layer           *string `form:"layer,omitempty" json:"layer,omitempty"`
	ContentType     *string `form:"content_type,omitempty" json:"content_type,omitempty"`
	ContentCode     *string `form:"content_code,omitempty" json:"content_code,omitempty"`
	HideBlockedMaps *bool   `form:"hide_blocked_maps,omitempty" json:"hide_blocked_maps,omitempty"`
	HideEvent       *bool   `form:"hide_event,omitempty" json:"hide_event,omitempty"`
	Transition      *bool   `form:"transition,omitempty" json:"transition,omitempty"`
	Page            *int    `form:"page,omitempty" json:"page,omitempty"`
	Size            *int    `form:"size,omitempty" json:"size,omitempty"`
}

type GetLayerMapsMapsLayerGetParams struct {
	ContentType     *string `form:"content_type,omitempty" json:"content_type,omitempty"`
	ContentCode     *string `form:"content_code,omitempty" json:"content_code,omitempty"`
	HideBlockedMaps *bool   `form:"hide_blocked_maps,omitempty" json:"hide_blocked_maps,omitempty"`
	HideEvent       *bool   `form:"hide_event,omitempty" json:"hide_event,omitempty"`
	Transition      *bool   `form:"transition,omitempty" json:"transition,omitempty"`
	Page            *int    `form:"page,omitempty" json:"page,omitempty"`
	Size            *int    `form:"size,omitempty" json:"size,omitempty"`
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
	Code *string `form:"code,omitempty" json:"code,omitempty"`
	Type *string `form:"type,omitempty" json:"type,omitempty"`
	Page *int    `form:"page,omitempty" json:"page,omitempty"`
	Size *int    `form:"size,omitempty" json:"size,omitempty"`
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

type GetAllNpcsNpcsDetailsGetParams struct {
	Name     *string `form:"name,omitempty" json:"name,omitempty"`
	Type     *string `form:"type,omitempty" json:"type,omitempty"`
	Currency *string `form:"currency,omitempty" json:"currency,omitempty"`
	Item     *string `form:"item,omitempty" json:"item,omitempty"`
	Page     *int    `form:"page,omitempty" json:"page,omitempty"`
	Size     *int    `form:"size,omitempty" json:"size,omitempty"`
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

type GetAllRaidsRaidsGetParams struct {
	Name   *string `form:"name,omitempty" json:"name,omitempty"`
	Active *bool   `form:"active,omitempty" json:"active,omitempty"`
	Page   *int    `form:"page,omitempty" json:"page,omitempty"`
	Size   *int    `form:"size,omitempty" json:"size,omitempty"`
}

type GetRaidLeaderboardRaidsCodeLeaderboardGetParams struct {
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	Size *int `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllResourcesResourcesGetParams struct {
	MinLevel *int    `form:"min_level,omitempty" json:"min_level,omitempty"`
	MaxLevel *int    `form:"max_level,omitempty" json:"max_level,omitempty"`
	Skill    *string `form:"skill,omitempty" json:"skill,omitempty"`
	Drop     *string `form:"drop,omitempty" json:"drop,omitempty"`
	Page     *int    `form:"page,omitempty" json:"page,omitempty"`
	Size     *int    `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllSeasonRewardsSeasonRewardsGetParams struct {
	Type *string `form:"type,omitempty" json:"type,omitempty"`
	Page *int    `form:"page,omitempty" json:"page,omitempty"`
	Size *int    `form:"size,omitempty" json:"size,omitempty"`
}

type GetSeasonRewardsByCodeSeasonRewardsCodeGetParams struct {
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	Size *int `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllSkinsSkinsGetParams struct {
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	Size *int `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllTasksTasksListGetParams struct {
	MinLevel *int    `form:"min_level,omitempty" json:"min_level,omitempty"`
	MaxLevel *int    `form:"max_level,omitempty" json:"max_level,omitempty"`
	Skill    *string `form:"skill,omitempty" json:"skill,omitempty"`
	Type     *string `form:"type,omitempty" json:"type,omitempty"`
	Page     *int    `form:"page,omitempty" json:"page,omitempty"`
	Size     *int    `form:"size,omitempty" json:"size,omitempty"`
}

type GetAllTasksRewardsTasksRewardsGetParams struct {
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	Size *int `form:"size,omitempty" json:"size,omitempty"`
}
