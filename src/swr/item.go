/*  Star Wars Role-Playing Mud
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
	Id         uint     `yaml:"id"`                      // instance id of the item
	OId        uint     `yaml:"itemId,omitempty"`        // item type id.
	Filename   string   `yaml:"-"`                       // filename for this item
	Name       string   `yaml:"name"`                    // name of the item
	Desc       string   `yaml:"desc"`                    // description of the item
	Keywords   []string `yaml:"keywords,flow"`           // keywords for the item
	Type       string   `yaml:"type"`                    // item type, a value of ITEM_TYPE_* const.
	Value      int      `yaml:"value"`                   // how much is this item generally worth?
	Weight     int      `yaml:"weight"`                  // how much does this item weigh?
	AC         int      `yaml:"ac,omitempty"`            // If armor, what's the AC (common AC values are 1-8 for torso, 2-3 for hands/head/feet, 0-1 for waist)
	WearLoc    *string  `yaml:"wearLoc,omitempty"`       // where is this item worn? nil means it's not wearable.
	WeaponType *string  `yaml:"weaponType,omitempty"`    // weapon type from ITEM_WEAPON_TYPE_* const, nil means it's not a weapon.
	Dmg        *string  `yaml:"dmgRoll,omitempty"`       // Damage roll represented by a D20 compatible string. Weapons do damage.
	Items      []Item   `yaml:"contains,omitempty,flow"` // If item type is "container", then this is the list of stored items.
}

type Item interface {
	GetId() uint           // instance id
	GetTypeId() uint       // item type id
	GetData() *ItemData    // underlying core [ItemData] struct pointer
	GetKeywords() []string // keywords for this item
	IsWeapon() bool        // is item a weapon?
	IsContainer() bool     // is item a container? (corpses are containers)
	IsWearable() bool      // is item wearable?
	IsCorpse() bool        // is item a corpse? (corpses are containers but containers aren't always corpses)
	GetWeight() int        // item weight in kg.
}

func (i *ItemData) GetData() *ItemData {
	return i
}

func (i *ItemData) GetId() uint {
	return i.Id
}

func (i *ItemData) GetTypeId() uint {
	if i.OId == 0 {
		return i.Id
	}
	return i.OId
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
		Id:         gen_item_id(),
		OId:        i.Id,
		Name:       i.Name,
		Filename:   i.Filename,
		Desc:       i.Desc,
		Keywords:   make([]string, 0),
		Type:       i.Type,
		Value:      i.Value,
		Weight:     i.Weight,
		AC:         i.AC,
		WearLoc:    i.WearLoc,
		WeaponType: i.WeaponType,
		Dmg:        i.Dmg,
		Items:      make([]Item, 0),
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
	return i.Type == ITEM_TYPE_CONTAINER || i.IsCorpse()
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

func item_get_weapon_skill(item Item) string {
	if item == nil {
		return "martial-arts"
	}
	data := item.GetData()
	if data.Type == ITEM_TYPE_1H_WEAPON || data.Type == ITEM_TYPE_2H_WEAPON {
		weaponType := "vibro-blades"
		if data.WeaponType != nil {
			weaponType = *data.WeaponType
		}
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
