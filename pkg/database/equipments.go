package database

var (
	EquipmentTypes = []string{
		"amulet",
		"artifact",
		"bag",
		"body_armor",
		"boots",
		"helmet",
		"leg_armor",
		"ring",
		"rune",
		"shield",
		"tool",
		"utility",
		"weapon",
	}

	EquipmentTypeToSlots = map[string][]string{
		"amulet":     {"amulet"},
		"artifact":   {"artifact1", "artifact2", "artifact3"},
		"bag":        {"bag"},
		"body_armor": {"body_armor"},
		"boots":      {"boots"},
		"helmet":     {"helmet"},
		"leg_armor":  {"leg_armor"},
		"ring":       {"ring1", "ring2"},
		"rune":       {"rune"},
		"shield":     {"shield"},
		"utility":    {"utility1", "utility2"},
		"weapon":     {"weapon"},
	}

	EquipmentSlotToTypes = map[string]string{
		"amulet":     "amulet",
		"artifact1":  "artifact",
		"artifact2":  "artifact",
		"artifact3":  "artifact",
		"bag":        "bag",
		"body_armor": "body_armor",
		"boots":      "boots",
		"helmet":     "helmet",
		"leg_armor":  "leg_armor",
		"ring1":      "ring",
		"ring2":      "ring",
		"rune":       "rune",
		"shield":     "shield",
		"utility1":   "utility",
		"utility2":   "utility",
		"weapon":     "weapon",
	}
)
