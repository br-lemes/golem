package cmd

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type xpFlags struct {
	Group string   `flag:"group" shorthand:"g" desc:"Group by \"skill\" or \"character\" (defaults to \"skill\")"`
	Skill []string `flag:"skill" shorthand:"k" desc:"Show XP for specific skills"`
	Name  []string `flag:"name" shorthand:"n" desc:"Calculate from specific characters"`
	From  int      `flag:"from" default:"-1" desc:"Calculate from this starting level (0 for zero XP)"`
}

var xpGroups = []string{"skill", "character"}

var xpCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(2),
	Use:   "xp <level> [account]",
	Short: "Show the XP needed to reach a level",
	Long: `Show the XP needed to reach a level

Arguments:
  level   The target level (between 1 and 50).
  account The name of the account (optional).`,
	ValidArgsFunction: completion.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		account := ""
		if len(args) > 1 {
			account = args[1]
		}
		options, err := utils.ReadFlags[xpFlags](cmd)
		if err != nil {
			return err
		}
		groupChanged := cmd.Flags().Changed("group")
		if !groupChanged {
			options.Group = "skill"
		}
		err = xpValidate(target, options, groupChanged)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return xpRun(target, account, options)
	},
}

func xpValidate(target string, options xpFlags, groupChanged bool) error {
	if !slices.Contains(xpGroups, options.Group) {
		return fmt.Errorf("invalid group %q: allowed values are %v", options.Group, xpGroups)
	}
	level, err := strconv.Atoi(target)
	if err != nil || level < 1 || level > 50 {
		return fmt.Errorf("invalid target level %q: must be between 1 and 50", target)
	}
	validSkills := database.Enum("CharacterLeaderboardType")
	for _, skill := range options.Skill {
		if !slices.Contains(validSkills, skill) {
			return fmt.Errorf("invalid skill %q: allowed values are %v", skill, validSkills)
		}
	}
	if options.From < -1 || options.From > 50 {
		return fmt.Errorf("invalid starting level %d: must be between 0 and 50", options.From)
	}
	if options.From >= 0 {
		if len(options.Skill) > 0 || len(options.Name) > 0 || groupChanged {
			return fmt.Errorf("--from cannot be combined with --skill, --name or --group")
		}
	}
	return nil
}

func xpRun(targetString, account string, xpOptions xpFlags) error {
	target, _ := strconv.Atoi(targetString)
	skills := xpOptions.Skill
	if len(skills) == 0 {
		skills = database.Enum("CharacterLeaderboardType")
	}
	if xpOptions.From >= 0 {
		startLevel := max(1, xpOptions.From)
		return console.Auto(xpNeeded(startLevel, 0, target))
	}
	characters, err := api.AccountsCharacters(account)
	if err != nil {
		return err
	}
	names := xpOptions.Name
	if len(names) == 0 {
		for _, character := range characters {
			names = append(names, character.Name)
		}
	}
	bySkill := map[string]map[string]int{}
	for _, name := range names {
		found := false
		for _, character := range characters {
			if character.Name != name {
				continue
			}
			found = true
			for _, skill := range skills {
				level, xp := characterSkillProgress(character, skill)
				if bySkill[skill] == nil {
					bySkill[skill] = map[string]int{}
				}
				bySkill[skill][name] = xpNeeded(level, xp, target)
			}
			break
		}
		if !found {
			return fmt.Errorf("character %q not found in account %q", name, account)
		}
	}
	if xpOptions.Group == "skill" {
		return console.Auto(bySkill)
	}
	return console.Auto(xpByCharacter(bySkill))
}

func xpByCharacter(bySkill map[string]map[string]int) map[string]map[string]int {
	result := map[string]map[string]int{}
	for skill, characterXP := range bySkill {
		for character, xp := range characterXP {
			if result[character] == nil {
				result[character] = map[string]int{}
			}
			result[character][skill] = xp
		}
	}
	return result
}

func xpNeeded(startLevel, currentXP, targetLevel int) int {
	if targetLevel <= startLevel {
		return 0
	}
	needed := -currentXP
	for level := startLevel; level < targetLevel; level++ {
		needed += xpToNextLevel[level]
	}
	return max(0, needed)
}

var xpToNextLevel = [...]int{
	0,
	150,
	250,
	350,
	450,
	700,
	950,
	1200,
	1450,
	1700,
	2100,
	2500,
	2900,
	3300,
	3700,
	4400,
	5100,
	5800,
	6500,
	7200,
	8200,
	9200,
	10200,
	11200,
	12200,
	13400,
	14600,
	15800,
	17000,
	18200,
	19700,
	21200,
	22700,
	24200,
	25700,
	27500,
	29300,
	31100,
	32900,
	34700,
	36500,
	38600,
	40700,
	42800,
	44900,
	47000,
	48800,
	50600,
	52400,
	54200,
}

func characterSkillProgress(c schemas.CharacterSchema, skill string) (int, int) {
	level, _ := utils.GetCharacterSkillLevel(c, skill)
	if skill == "combat" {
		return level, c.Xp
	}
	values := map[string]int{
		"alchemy":         c.AlchemyXp,
		"cooking":         c.CookingXp,
		"fishing":         c.FishingXp,
		"gearcrafting":    c.GearcraftingXp,
		"jewelrycrafting": c.JewelrycraftingXp,
		"mining":          c.MiningXp,
		"weaponcrafting":  c.WeaponcraftingXp,
		"woodcutting":     c.WoodcuttingXp,
	}
	return level, values[skill]
}

func init() {
	rootCmd.AddCommand(xpCmd)
	err := utils.RegisterFlags[xpFlags](xpCmd)
	if err != nil {
		panic(err)
	}
	err = xpCmd.RegisterFlagCompletionFunc("group", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return xpGroups, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
	err = xpCmd.RegisterFlagCompletionFunc("skill", completion.StringSlice(func() []string {
		return database.Enum("CharacterLeaderboardType")
	}))
	if err != nil {
		panic(err)
	}
	err = xpCmd.RegisterFlagCompletionFunc("name", completion.StringSlice(utils.GetCharacters))
	if err != nil {
		panic(err)
	}
}
