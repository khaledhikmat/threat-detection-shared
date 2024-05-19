package utils

import (
	"math/rand"
	"os"
	"time"
)

var humanTags = []string{"person", "male", "female", "young", "old", "child", "adult"}
var clothingTags = []string{"clothing", "shirt", "t-shirt", "jacket", "coat", "trousers", "jeans", "shorts", "skirt", "dress", "suit", "tie", "scarf", "gloves", "hat", "cap", "shoes", "boots", "trainers", "sandals", "flip-flops", "slippers", "socks", "tights", "underwear", "bra", "pants", "knickers", "boxers", "vest", "nightwear", "pyjamas", "dressing-gown", "slippers", "onesie", "swimwear", "bikini", "swimming-trunks", "wetsuit", "accessories", "belt", "scarf", "gloves", "hat", "cap", "sunglasses", "umbrella", "bag", "backpack", "handbag", "briefcase", "suitcase", "wallet", "purse", "watch", "bracelet", "necklace", "earrings", "ring", "brooch", "tie-pin", "cufflinks", "hairband", "hairclip", "hairgrip", "hairslide", "hairpin", "hair-tie", "hair-bobble", "hair-clip", "hair-slide", "hair-pin", "hair-tie", "hair-bobble", "hairband", "hairband", "hairclip", "hairgrip", "hairslide", "hairpin", "hair-tie", "hair-bobble", "hair-clip", "hair-slide", "hair-pin", "hair-tie", "hair-bobble", "hairband", "hairband", "hairclip", "hairgrip", "hairslide", "hairpin", "hair-tie", "hair-bobble", "hair-clip", "hair-slide", "hair-pin", "hair-tie", "hair-bobble"}
var fireTags = []string{"fire", "flame", "smoke", "burning", "heat", "inferno", "blaze", "embers", "ignition", "combustion", "scorch", "char", "ash", "soot", "pyre", "incineration", "firestorm", "fireball", "fireplace", "firewood", "firefighter", "fireproof", "firebreak", "firebug", "firecracker", "firefighting", "firefly", "fireman", "firepower", "firetrap", "firewater", "firework", "firearm", "firebrand", "fireclay", "firecracker", "firehouse", "firelight", "fireplace", "fireplug", "firepower", "fireproof", "fireside", "firetrap", "firewater", "firewood", "fireworks", "firearm", "fireball", "firebrand", "firebreak", "firebug", "fireclay", "firefly", "fireman", "fireplug", "firetrap", "firewater", "firewood", "fireworks", "firearm", "fireball", "firebrand", "firebreak", "firebug", "fireclay", "firefly", "fireman", "fireplug"}
var weaponTags = []string{"weapon", "gun", "rifle", "pistol", "firearm", "blade", "knife", "sword", "dagger", "spear", "lance", "pike", "halberd", "axe", "mace", "club", "staff", "baton", "bow", "crossbow", "sling", "javelin", "harpoon", "trident"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Contains[T comparable](slice []T, item T) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func RandWeaponTags(n int) []string {
	var tags []string
	for i := 0; i < n; i++ {
		tags = append(tags, RandWeaponTag())
		tags = append(tags, RandHumanTag())
		tags = append(tags, RandClothingTag())
	}
	return tags
}

func RandFireTags(n int) []string {
	var tags []string
	for i := 0; i < n; i++ {
		tags = append(tags, RandFireTag())
		tags = append(tags, RandHumanTag())
		tags = append(tags, RandClothingTag())
	}
	return tags
}

func RandHumanTag() string {
	return humanTags[rand.Intn(len(humanTags))]
}

func RandClothingTag() string {
	return clothingTags[rand.Intn(len(clothingTags))]
}

func RandFireTag() string {
	return fireTags[rand.Intn(len(fireTags))]
}

func RandWeaponTag() string {
	return weaponTags[rand.Intn(len(weaponTags))]
}
