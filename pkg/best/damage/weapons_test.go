package damage

import (
	"slices"
	"testing"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestFindWeaponsExcludesTools(t *testing.T) {
	character := schemas.CharacterSchema{Level: 50}
	weapons := findWeapons(character, Options{})
	if rawWeaponIndex(weapons, "copper_axe") != -1 {
		t.Fatal("tool weapon included")
	}
	if rawWeaponIndex(weapons, "iron_sword") == -1 {
		t.Fatal("combat weapon missing")
	}
}

func TestFindWeaponsFiltersByCharacterLevel(t *testing.T) {
	low := findWeapons(schemas.CharacterSchema{Level: 5}, Options{})
	if rawWeaponIndex(low, "iron_sword") != -1 {
		t.Fatal("weapon above character level included")
	}
	if rawWeaponIndex(low, "sticky_sword") == -1 {
		t.Fatal("weapon within character level missing")
	}
	high := findWeapons(schemas.CharacterSchema{Level: 10}, Options{})
	if rawWeaponIndex(high, "iron_sword") == -1 {
		t.Fatal("weapon at character level missing")
	}
}

func TestFindWeaponsFiltersByAvailability(t *testing.T) {
	monster, found := catalogMonster("mushmush")
	if !found {
		t.Fatal("mushmush missing from catalog")
	}
	options := Options{Available: map[string]int{"forest_whip": 1}}
	weapons := FindWeapons(monster, schemas.CharacterSchema{Level: 26}, options)
	if len(weapons) != 1 || weapons[0].Code != "forest_whip" {
		t.Fatalf("weapons = %v, want only forest_whip", weapons)
	}
}

func TestFindWeaponsExcludesWithout(t *testing.T) {
	monster, found := catalogMonster("mushmush")
	if !found {
		t.Fatal("mushmush missing from catalog")
	}
	options := Options{Without: map[string]bool{"forest_whip": true}}
	weapons := FindWeapons(monster, schemas.CharacterSchema{Level: 26}, options)
	if weaponIndex(weapons, "forest_whip") != -1 {
		t.Fatal("excluded weapon included")
	}
	if weaponIndex(weapons, "highwayman_dagger") == -1 {
		t.Fatal("weapon not excluded missing")
	}
}

func TestFindWeaponRejectsWeaponOutsideAccount(t *testing.T) {
	monster, found := catalogMonster("mushmush")
	if !found {
		t.Fatal("mushmush missing from catalog")
	}
	options := Options{Available: map[string]int{}}
	_, err := FindWeapon(monster, schemas.CharacterSchema{Level: 26}, options, "forest_whip")
	if err == nil {
		t.Fatal("weapon outside the account accepted")
	}
}

func TestFindWeaponRejectsExcludedWeapon(t *testing.T) {
	monster, found := catalogMonster("mushmush")
	if !found {
		t.Fatal("mushmush missing from catalog")
	}
	options := Options{Without: map[string]bool{"forest_whip": true}}
	_, err := FindWeapon(monster, schemas.CharacterSchema{Level: 26}, options, "forest_whip")
	if err == nil {
		t.Fatal("excluded weapon accepted")
	}
}

func TestFindWeaponsPrunesDominatedWeapons(t *testing.T) {
	monster, found := catalogMonster("mushmush")
	if !found {
		t.Fatal("mushmush missing from catalog")
	}
	weapons := FindWeapons(monster, schemas.CharacterSchema{Level: 26}, Options{})
	if weaponIndex(weapons, "forest_whip") == -1 {
		t.Fatal("surviving weapon missing")
	}
	if weaponIndex(weapons, "iron_dagger") != -1 {
		t.Fatal("dominated weapon included")
	}
}

func TestFindWeaponsReportsDamageRange(t *testing.T) {
	monster, found := catalogMonster("mushmush")
	if !found {
		t.Fatal("mushmush missing from catalog")
	}
	weapons := FindWeapons(monster, schemas.CharacterSchema{Level: 26}, Options{})
	index := weaponIndex(weapons, "forest_whip")
	if index == -1 {
		t.Fatal("forest_whip missing from result")
	}
	weapon := weapons[index]
	if weapon.Min != 52 {
		t.Fatalf("forest_whip min = %d, want 52", weapon.Min)
	}
	if weapon.Max != 104 {
		t.Fatalf("forest_whip max = %d, want 104", weapon.Max)
	}
}

func TestFindWeaponsOrdersByMaximumDamage(t *testing.T) {
	monster, found := catalogMonster("mushmush")
	if !found {
		t.Fatal("mushmush missing from catalog")
	}
	weapons := FindWeapons(monster, schemas.CharacterSchema{Level: 26}, Options{})
	for index := 1; index < len(weapons); index++ {
		if weapons[index-1].Max < weapons[index].Max {
			t.Fatalf("weapon %s ordered after lower max damage", weapons[index-1].Code)
		}
	}
}

func TestFindWeaponReportsCombinations(t *testing.T) {
	monster, found := catalogMonster("mushmush")
	if !found {
		t.Fatal("mushmush missing from catalog")
	}
	weapon, err := FindWeapon(monster, schemas.CharacterSchema{Level: 26}, Options{}, "forest_whip")
	if err != nil {
		t.Fatalf("FindWeapon: %v", err)
	}
	if len(weapon.Combinations) != 2 {
		t.Fatalf("combinations = %d, want 2", len(weapon.Combinations))
	}
	best := weapon.Combinations[0]
	if best.Damage != weapon.Max {
		t.Fatalf("best damage = %d, want %d", best.Damage, weapon.Max)
	}
	if best.Damage != 104 {
		t.Fatalf("best damage = %d, want 104", best.Damage)
	}
	if best.Bonuses.Air != 95 {
		t.Fatalf("best air = %d, want 95", best.Bonuses.Air)
	}
	if best.Critical != 5 {
		t.Fatalf("best critical = %d, want 5", best.Critical)
	}
	if len(best.Items) != 11 {
		t.Fatalf("best items = %d, want 11", len(best.Items))
	}
	tradeOff := weapon.Combinations[len(weapon.Combinations)-1]
	if tradeOff.Damage != 75 {
		t.Fatalf("trade-off damage = %d, want 75", tradeOff.Damage)
	}
	if tradeOff.Critical != 11 {
		t.Fatalf("trade-off critical = %d, want 11", tradeOff.Critical)
	}
	if len(tradeOff.Items) != 11 {
		t.Fatalf("trade-off items = %d, want 11", len(tradeOff.Items))
	}
}

func TestPruneCombinationsKeepsTradeOffs(t *testing.T) {
	got := pruneCombinations([]Combination{
		{Damage: 100, Critical: 10},
		{Damage: 90, Critical: 20},
		{Damage: 80, Critical: 5},
		{Damage: 100, Critical: 10},
	})
	if len(got) != 3 {
		t.Fatalf("combinations = %d, want 3", len(got))
	}
	if got[0].Damage != 100 || got[0].Critical != 10 {
		t.Fatalf("first = %+v, want 100/10", got[0])
	}
	if got[1].Damage != 90 || got[1].Critical != 20 {
		t.Fatalf("second = %+v, want 90/20", got[1])
	}
}

func TestFindWeaponRejectsUnavailableWeapon(t *testing.T) {
	monster, found := catalogMonster("mushmush")
	if !found {
		t.Fatal("mushmush missing from catalog")
	}
	_, err := FindWeapon(monster, schemas.CharacterSchema{Level: 1}, Options{}, "adamantite_sword")
	if err == nil {
		t.Fatal("weapon above character level accepted")
	}
}

func TestFindWeaponKeepsCriticalTradeOff(t *testing.T) {
	monster, found := catalogMonster("mushmush")
	if !found {
		t.Fatal("mushmush missing from catalog")
	}
	weapon, err := FindWeapon(monster, schemas.CharacterSchema{Level: 26}, Options{}, "highwayman_dagger")
	if err != nil {
		t.Fatalf("FindWeapon: %v", err)
	}
	if weapon.Critical != 35 {
		t.Fatalf("highwayman_dagger critical = %d, want 35", weapon.Critical)
	}
}

func TestFindWeaponRejectsDominatedWeapon(t *testing.T) {
	monster, found := catalogMonster("mushmush")
	if !found {
		t.Fatal("mushmush missing from catalog")
	}
	_, err := FindWeapon(monster, schemas.CharacterSchema{Level: 26}, Options{}, "iron_dagger")
	if err == nil {
		t.Fatal("dominated weapon accepted")
	}
}

func catalogMonster(code string) (schemas.MonsterSchema, bool) {
	monster, found := catalog.Monsters.Get(code)
	if !found {
		return schemas.MonsterSchema{}, false
	}
	return *monster, true
}

func weaponIndex(weapons []Weapon, code string) int {
	return slices.IndexFunc(weapons, func(weapon Weapon) bool {
		return weapon.Code == code
	})
}

func rawWeaponIndex(weapons []weapon, code string) int {
	return slices.IndexFunc(weapons, func(candidate weapon) bool {
		return candidate.code == code
	})
}
