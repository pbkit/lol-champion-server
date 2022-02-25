package main

import (
	"regexp"
	"strconv"
	"strings"

	dd "github.com/pbkit/lol-champion-server/datadragon"
	proto "github.com/pbkit/lol-champion-server/gen"
)

var statToNameMap = map[string]string{
	"armor":        "방어력",
	"attackdamage": "공격력",
	"attackrange":  "사거리",
	"attackspeed":  "공격속도",
	"hp":           "체력",
	"hpregen":      "체력회복",
	"movespeed":    "이동속도",
	"mp":           "마나",
	"mpregen":      "마나회복",
	"spellblock":   "마법방어",
}

var statToIsProportional = map[string]bool{
	"attackspeed": true,
}

var tagToTypeMap = map[string]string{
	"Assassin": "암살자",
	"Fighter":  "전사",
	"Mage":     "마법사",
	"Marksman": "원거리",
	"Support":  "서포터",
	"Tank":     "탱커",
}

func championTypes(tags []string) []string {
	var ret []string
	for _, tag := range tags {
		if val, ok := tagToTypeMap[tag]; ok {
			ret = append(ret, val)
		}
	}
	return ret
}

func ListChampions(championFull dd.ChampionFull) *proto.ListChampionsResponse {
	champions := make([]*proto.ListChampionsResponse_Champion, len(championFull.Data))

	ptr := 0
	for key, v := range championFull.Data {
		champions[ptr] = &proto.ListChampionsResponse_Champion{
			Name:            v.Name,
			Key:             key,
			Title:           v.Title,
			ProfileImageUrl: dd.ResolveImageURL(v.Image),
			Types:           championTypes(v.Tags),
		}
		ptr += 1
	}

	return &proto.ListChampionsResponse{
		Champions: champions,
	}
}

func GetChampionStory(championFull dd.ChampionFull, key string) *proto.GetChampionStoryResponse {
	data, ok := championFull.Data[key]
	if !ok {
		return nil
	}

	return &proto.GetChampionStoryResponse{
		Story:              data.Blurb,
		BackgroundImageUrl: dd.ResolveChampionTileImageURL(key),
	}
}

func GetChampionStats(championFull dd.ChampionFull, key string) *proto.GetChampionStatsResponse {
	data, ok := championFull.Data[key]
	if !ok {
		return nil
	}

	statsMap := map[string]*proto.GetChampionStatsResponse_Stat{}
	for key, val := range data.Stats {
		statKey := strings.TrimSuffix(key, "perlevel")
		statName := statKey
		if name, ok := statToNameMap[statKey]; ok {
			statName = name
		} else {
			continue
		}

		if _, ok := statsMap[statKey]; !ok {
			statsMap[statKey] = &proto.GetChampionStatsResponse_Stat{
				Name: statName,
			}
		} else {
		}
		if strings.HasSuffix(key, "perlevel") {
			if _, ok := statToIsProportional[statKey]; ok {
				statsMap[statKey].PerLevel = &proto.GetChampionStatsResponse_Stat_Proportional{
					Proportional: val,
				}
			} else {
				statsMap[statKey].PerLevel = &proto.GetChampionStatsResponse_Stat_Exact{
					Exact: val,
				}
			}
		} else {
			statsMap[statKey].Base = val
		}
	}

	stats := make([]*proto.GetChampionStatsResponse_Stat, 0, len(statsMap))
	for _, stat := range statsMap {
		stats = append(stats, stat)
	}

	return &proto.GetChampionStatsResponse{
		Stats: stats,
	}
}

func getRegexParams(expr *regexp.Regexp, s string) map[string]string {
	params := expr.FindStringSubmatch(s)
	if params == nil {
		return nil
	}

	ret := map[string]string{}
	for i, name := range expr.SubexpNames() {
		if i > 0 && i <= len(params) {
			ret[name] = params[i]
		}
	}

	return ret
}

func getSpellType(spellId string) proto.GetChampionSkillsResponse_SpellType {
	switch {
	case strings.HasSuffix(spellId, "Q"):
		return proto.GetChampionSkillsResponse_ACTIVE_Q
	case strings.HasSuffix(spellId, "W"):
		return proto.GetChampionSkillsResponse_ACTIVE_W
	case strings.HasSuffix(spellId, "E"):
		return proto.GetChampionSkillsResponse_ACTIVE_E
	case strings.HasSuffix(spellId, "R"):
		return proto.GetChampionSkillsResponse_ACTIVE_R
	}
	return proto.GetChampionSkillsResponse_UNSPECIFIED
}

func GetChampionSkills(championFull dd.ChampionFull, key string) *proto.GetChampionSkillsResponse {
	data, ok := championFull.Data[key]
	if !ok {
		return nil
	}

	// Matches only trivial cases like {{ e0 }} -> {{ e0NL }}
	// Interpreting other cases requires CommunityDragon data.
	effectPattern := regexp.MustCompile(`{{\s*e(?P<effectIdx>[0-9])\s*}}\s*->\s*{{\s*e[0-9]NL\s*}}(?P<isProportional>%)?`)

	spells := []*proto.GetChampionSkillsResponse_Spell{}
	spells = append(spells, &proto.GetChampionSkillsResponse_Spell{
		Name:        data.Passive.Name,
		Description: data.Passive.Description,
		ImageUrl:    dd.ResolveImageURL(data.Passive.Image),
		Type:        proto.GetChampionSkillsResponse_PASSIVE,
	})

	for _, spell := range data.Spells {
		spellEntry := &proto.GetChampionSkillsResponse_Spell{
			Name:        spell.Name,
			Description: spell.Description,
			ImageUrl:    dd.ResolveImageURL(spell.Image),
			MaxLevel:    spell.Maxrank,
			Type:        getSpellType(spell.Id),
			Stats:       []*proto.GetChampionSkillsResponse_SpellStat{},
		}
		spellEntry.Stats = append(spellEntry.Stats, &proto.GetChampionSkillsResponse_SpellStat{
			Name:   "쿨다운",
			Values: spell.Cooldown,
		})
		spellEntry.Stats = append(spellEntry.Stats, &proto.GetChampionSkillsResponse_SpellStat{
			Name:   "사거리",
			Values: spell.Range,
		})
		for i, levelTipName := range spell.Leveltip.Label {
			if strings.Contains(levelTipName, "{{") {
				continue
			}
			effectName := levelTipName
			effectParams := getRegexParams(effectPattern, spell.Leveltip.Effect[i])
			if effectParams == nil {
				continue
			} else if _, ok := effectParams["isProportional"]; ok {
				effectName += "(%)"
			}
			effectValueId, _ := strconv.ParseInt(effectParams["effectIdx"], 10, 32)
			spellStat := &proto.GetChampionSkillsResponse_SpellStat{
				Name:   effectName,
				Values: spell.Effect[effectValueId],
			}
			spellEntry.Stats = append(spellEntry.Stats, spellStat)
		}
		spells = append(spells, spellEntry)
	}

	return &proto.GetChampionSkillsResponse{
		Spells: spells,
	}
}
