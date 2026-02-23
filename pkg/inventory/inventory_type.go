package inventory

type InventoryType struct {
	id          int
	defaultSize int
	title       string
	networkType byte
}

const (
	TypeChest          = 0
	TypeDoubleChest    = 1
	TypePlayer         = 2
	TypeFurnace        = 3
	TypeCrafting       = 4
	TypeWorkbench      = 5
	TypeBrewingStand   = 7
	TypeAnvil          = 8
	TypeEnchantTable   = 9
	TypeDispenser      = 10
	TypeDropper        = 11
	TypeHopper         = 12
	TypeMinecartChest  = 13
	TypeMinecartHopper = 14
	TypeMob            = 253
)

const (
	NetworkInventory    byte = 0xFF
	NetworkContainer    byte = 0
	NetworkWorkbench    byte = 1
	NetworkFurnace      byte = 2
	NetworkEnchantment  byte = 3
	NetworkBrewingStand byte = 4
	NetworkAnvil        byte = 5
	NetworkDispenser    byte = 6
	NetworkDropper      byte = 7
	NetworkHopper       byte = 8
)

var inventoryTypes map[int]*InventoryType

func init() {
	inventoryTypes = map[int]*InventoryType{
		TypeChest:          {id: TypeChest, defaultSize: 27, title: "Chest", networkType: NetworkContainer},
		TypeDoubleChest:    {id: TypeDoubleChest, defaultSize: 54, title: "Double Chest", networkType: NetworkContainer},
		TypePlayer:         {id: TypePlayer, defaultSize: 40, title: "Player", networkType: NetworkInventory},
		TypeFurnace:        {id: TypeFurnace, defaultSize: 3, title: "Furnace", networkType: NetworkFurnace},
		TypeCrafting:       {id: TypeCrafting, defaultSize: 5, title: "Crafting", networkType: NetworkInventory},
		TypeWorkbench:      {id: TypeWorkbench, defaultSize: 10, title: "Crafting", networkType: NetworkWorkbench},
		TypeEnchantTable:   {id: TypeEnchantTable, defaultSize: 2, title: "Enchant", networkType: NetworkEnchantment},
		TypeBrewingStand:   {id: TypeBrewingStand, defaultSize: 4, title: "Brewing", networkType: NetworkBrewingStand},
		TypeAnvil:          {id: TypeAnvil, defaultSize: 3, title: "Anvil", networkType: NetworkAnvil},
		TypeDispenser:      {id: TypeDispenser, defaultSize: 9, title: "Dispenser", networkType: NetworkDispenser},
		TypeDropper:        {id: TypeDropper, defaultSize: 9, title: "Dropper", networkType: NetworkDropper},
		TypeHopper:         {id: TypeHopper, defaultSize: 5, title: "Hopper", networkType: NetworkHopper},
		TypeMinecartChest:  {id: TypeMinecartChest, defaultSize: 27, title: "Minecart with Chest", networkType: NetworkContainer},
		TypeMinecartHopper: {id: TypeMinecartHopper, defaultSize: 5, title: "Minecart with Hopper", networkType: NetworkHopper},
		TypeMob:            {id: TypeMob, defaultSize: 5, title: "Mob", networkType: NetworkContainer},
	}
}

func GetInventoryType(typeID int) *InventoryType {
	return inventoryTypes[typeID]
}

func (t *InventoryType) GetDefaultSize() int     { return t.defaultSize }
func (t *InventoryType) GetDefaultTitle() string { return t.title }
func (t *InventoryType) GetNetworkType() byte    { return t.networkType }
func (t *InventoryType) GetID() int              { return t.id }
