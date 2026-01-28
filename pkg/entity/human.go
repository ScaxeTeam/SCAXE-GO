package entity

type Human struct {
	*Living

	UUID     string
	Username string

	SkinData string
	SkinName string

	Food       float64
	MaxFood    float64
	Saturation float64
	Exhaustion float64

	TotalXP    int
	XPLevel    int
	XPProgress float64

	Absorption int

	HeldItemSlot int
}

func NewHuman() *Human {
	h := &Human{
		Living:       NewLiving(),
		UUID:         "",
		Username:     "",
		SkinData:     "",
		SkinName:     "",
		Food:         20,
		MaxFood:      20,
		Saturation:   20,
		Exhaustion:   0,
		TotalXP:      0,
		XPLevel:      0,
		XPProgress:   0,
		Absorption:   0,
		HeldItemSlot: 0,
	}
	h.initHumanAttributes()
	h.Height = 1.8
	h.Width = 0.6
	h.EyeHeight = 1.62
	return h
}

func (h *Human) initHumanAttributes() {
	h.Attributes.AddAttribute(GetDefaultAttribute(AttributeHunger))
	h.Attributes.AddAttribute(GetDefaultAttribute(AttributeSaturation))
	h.Attributes.AddAttribute(GetDefaultAttribute(AttributeExhaustion))
	h.Attributes.AddAttribute(GetDefaultAttribute(AttributeExperienceLevel))
	h.Attributes.AddAttribute(GetDefaultAttribute(AttributeExperience))
	h.Attributes.AddAttribute(GetDefaultAttribute(AttributeAbsorption))
}

func (h *Human) GetName() string {
	return h.Username
}

func (h *Human) GetSkinData() string {
	return h.SkinData
}

func (h *Human) GetSkinName() string {
	return h.SkinName
}

func (h *Human) SetSkin(data, name string) {
	h.SkinData = data
	h.SkinName = name
}

func (h *Human) GetUniqueID() string {
	return h.UUID
}

func (h *Human) GetFood() float64 {
	return h.Food
}

func (h *Human) SetFood(food float64) {
	if food < 0 {
		food = 0
	}
	if food > h.MaxFood {
		food = h.MaxFood
	}
	h.Food = food

	if attr := h.Attributes.GetAttribute(AttributeHunger); attr != nil {
		attr.SetValue(food)
	}
}

func (h *Human) AddFood(amount float64) {
	h.SetFood(h.Food + amount)
}

func (h *Human) GetMaxFood() float64 {
	return h.MaxFood
}

func (h *Human) GetSaturation() float64 {
	return h.Saturation
}

func (h *Human) SetSaturation(saturation float64) {
	if saturation < 0 {
		saturation = 0
	}
	if saturation > h.Food {
		saturation = h.Food
	}
	h.Saturation = saturation

	if attr := h.Attributes.GetAttribute(AttributeSaturation); attr != nil {
		attr.SetValue(saturation)
	}
}

func (h *Human) AddSaturation(amount float64) {
	h.SetSaturation(h.Saturation + amount)
}

func (h *Human) GetExhaustion() float64 {
	return h.Exhaustion
}

func (h *Human) SetExhaustion(exhaustion float64) {
	if exhaustion < 0 {
		exhaustion = 0
	}
	h.Exhaustion = exhaustion

	if attr := h.Attributes.GetAttribute(AttributeExhaustion); attr != nil {
		attr.SetValue(exhaustion)
	}
}

func (h *Human) Exhaust(amount float64) float64 {
	h.Exhaustion += amount

	for h.Exhaustion >= 4 {
		h.Exhaustion -= 4
		if h.Saturation > 0 {
			h.SetSaturation(h.Saturation - 1)
		} else if h.Food > 0 {
			h.SetFood(h.Food - 1)
		}
	}

	h.SetExhaustion(h.Exhaustion)
	return amount
}

func (h *Human) GetXPLevel() int {
	return h.XPLevel
}

func (h *Human) SetXPLevel(level int) {
	if level < 0 {
		level = 0
	}
	h.XPLevel = level

	if attr := h.Attributes.GetAttribute(AttributeExperienceLevel); attr != nil {
		attr.SetValue(float64(level))
	}
}

func (h *Human) AddXPLevel(levels int) {
	h.SetXPLevel(h.XPLevel + levels)
}

func (h *Human) GetXPProgress() float64 {
	return h.XPProgress
}

func (h *Human) SetXPProgress(progress float64) {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	h.XPProgress = progress

	if attr := h.Attributes.GetAttribute(AttributeExperience); attr != nil {
		attr.SetValue(progress)
	}
}

func (h *Human) GetTotalXP() int {
	return h.TotalXP
}

func (h *Human) SetTotalXP(xp int) {
	if xp < 0 {
		xp = 0
	}
	h.TotalXP = xp
}

func (h *Human) AddXP(xp int) {
	h.TotalXP += xp

}

func (h *Human) GetAbsorption() int {
	return h.Absorption
}

func (h *Human) SetAbsorption(absorption int) {
	if absorption < 0 {
		absorption = 0
	}
	h.Absorption = absorption
}
