/*  Space Wars Rebellion Mud
 *  Copyright (C) 2022 @{See Authors}
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */
package swr

import (
	"strings"
)

const (
	ITEM_TYPE_GENERIC   = "generic"
	ITEM_TYPE_COMS      = "comlink"
	ITEM_TYPE_1H_WEAPON = "weapon"
	ITEM_TYPE_2H_WEAPON = "weapon-2h"
	ITEM_TYPE_CONTAINER = "container"
	ITEM_TYPE_ARMOR     = "armor"
	ITEM_TYPE_TRASH_BIN = "bin"
	ITEM_TYPE_KEY       = "key"
	ITEM_TYPE_CORPSE    = "corpse"
	ITEM_TYPE_MATERIAL  = "material"
)

const (
	ITEM_WEAPON_TYPE_KNIFE      = "vibro-blade"
	ITEM_WEAPON_TYPE_BLASTER    = "blaster"
	ITEM_WEAPON_TYPE_RIFLE      = "rifle"
	ITEM_WEAPON_TYPE_REPEATER   = "repeater"
	ITEM_WEAPON_TYPE_BOWCASTER  = "bowcaster"
	ITEM_WEAPON_TYPE_FORCEPIKE  = "force-pike"
	ITEM_WEAPON_TYPE_GRENADE    = "grenade"
	ITEM_WEAPON_TYPE_MINE       = "mine"
	ITEM_WEAPON_TYPE_CLAYMORE   = "claymore"
	ITEM_WEAPON_TYPE_LIGHTSABER = "lightsaber"
)

type ItemData struct {
	Id         uint     `yaml:"id"`
	Name       string   `yaml:"name"`
	Desc       string   `yaml:"desc"`
	Keywords   []string `yaml:"keywords,flow"`
	Type       string   `yaml:"type"`
	Value      int      `yaml:"value"`
	Weight     int      `yaml:"weight"`
	AC         int      `yaml:"ac,omitempty"`
	WearLoc    *string  `yaml:"wearLoc,omitempty"`
	WeaponType *string  `yaml:"weaponType,omitempty"`
	Dmg        *string  `yaml:"dmgRoll,omitempty"`
	Items      []Item   `yaml:"contains,omitempty,flow"`
}

type Item interface {
	GetId() uint
	GetData() *ItemData
	GetKeywords() []string
	IsWeapon() bool
	IsContainer() bool
	IsWearable() bool
	IsCorpse() bool
	GetWeight() int
}

func (i *ItemData) GetData() *ItemData {
	return i
}

func (i *ItemData) GetId() uint {
	return i.Id
}

func (i *ItemData) GetKeywords() []string {
	return i.Keywords
}

func (i *ItemData) GetWeight() int {
	weight := i.Weight
	if i.IsContainer() || i.IsCorpse() {
		for id := range i.Items {
			weight += i.Items[id].GetWeight()
		}
	}
	return weight
}

func item_clone(item Item) Item {
	i := item.GetData()
	c := &ItemData{
		Id:       i.Id,
		Name:     i.Name,
		Desc:     i.Desc,
		Keywords: make([]string, 0),
		Type:     i.Type,
		Value:    i.Value,
		Weight:   i.Weight,
		AC:       i.AC,
		WearLoc:  i.WearLoc,
		Dmg:      i.Dmg,
		Items:    make([]Item, 0),
	}
	for idx := range i.Items {
		con_item := i.Items[idx]
		if con_item == nil {
			continue
		}
		c.Items = append(c.Items, item_clone(con_item))
	}
	for idx := range i.Keywords {
		k := i.Keywords[idx]
		c.Keywords = append(c.Keywords, k)
	}
	return c
}

func (i *ItemData) IsWeapon() bool {
	return i.Type == ITEM_TYPE_1H_WEAPON || i.Type == ITEM_TYPE_2H_WEAPON
}

func (i *ItemData) IsWearable() bool {
	return i.WearLoc != nil
}

func (i *ItemData) IsContainer() bool {
	return i.Type == ITEM_TYPE_CONTAINER || i.Type == ITEM_TYPE_CORPSE
}

func (i *ItemData) IsCorpse() bool {
	return i.Type == ITEM_TYPE_CORPSE
}

func (i *ItemData) FindItemInContainer(keyword string) Item {
	if i.IsContainer() || i.IsCorpse() {
		for id := range i.Items {
			keys := i.Items[id].GetKeywords()
			for k := range keys {
				key := keys[k]
				if strings.HasPrefix(key, keyword) {
					return i.Items[id]
				}
			}
		}
	}
	return nil
}

func (i *ItemData) RemoveItem(item Item) {
	if !i.IsContainer() {
		return
	}
	idx := -1
	for id := range i.Items {
		if i.Items[id] == item {
			idx = id
		}
	}
	if idx > -1 {
		ret := make([]Item, len(i.Items)-1)
		ret = append(ret, i.Items[:idx]...)
		ret = append(ret, i.Items[idx+1:]...)
		i.Items = ret
	}
}

func (i *ItemData) AddItem(item Item) {
	if !i.IsContainer() {
		ErrorCheck(Err("Can't add item because i* is not a container!"))
		return
	}
	i.Items = append(i.Items, item)
}

func get_weapon_skill(item Item) string {
	if item == nil {
		return "martial-arts"
	}
	data := item.GetData()
	if data.Type == ITEM_TYPE_1H_WEAPON || data.Type == ITEM_TYPE_2H_WEAPON {
		weaponType := *data.WeaponType
		return weaponType
	} else {
		return "martial-arts"
	}
}

func get_weapon_skill_stat(weaponType string, str uint, dex uint) uint {
	switch weaponType {
	case ITEM_WEAPON_TYPE_KNIFE:
	case ITEM_WEAPON_TYPE_BLASTER:
	case ITEM_WEAPON_TYPE_BOWCASTER:
	case ITEM_WEAPON_TYPE_CLAYMORE:
	case ITEM_WEAPON_TYPE_LIGHTSABER:
		return (dex / 10)
	case ITEM_WEAPON_TYPE_FORCEPIKE:
	case ITEM_WEAPON_TYPE_RIFLE:
	case ITEM_WEAPON_TYPE_REPEATER:
	case ITEM_WEAPON_TYPE_GRENADE:
	case ITEM_WEAPON_TYPE_MINE:
		return (str / 10)
	default:
		return umin((str / 10), (dex / 10))
	}
	return 0
}
