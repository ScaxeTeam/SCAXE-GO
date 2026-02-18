# SCAXE-GO

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![MCPE Version](https://img.shields.io/badge/MCPE-0.14.3-green?style=flat)
![Platform](https://img.shields.io/badge/Platform-Windows%20|%20Linux-blue?style=flat)

**A high-fidelity Minecraft Pocket Edition 0.14.3 server core rewritten in Go**

*World generation algorithms aligned with Minecraft Java Edition 1.12.2 at bit-level precision*

</div>

---

## Features

- **High-Fidelity World Generation** - 93.77% terrain consistency with Minecraft Java 1.12.2
- **Bit-Level GenLayer Precision** - Biome system achieves 99.9% bit-level accuracy
- **High-Performance Concurrency** - Thread-safe chunk generation powered by Go goroutines
- **Full Protocol Implementation** - Complete MCPE 0.14.3 (Protocol 70) support
- **Lua Plugin System** - Extensible plugin architecture with Lua scripting and hot-reload

---

## Implemented Features

### Core Systems

#### World Generation Engine (Gorigional)

A precise port of the Overworld generator from Minecraft Java 1.12.2 source:

| Module               | Status   | Accuracy          |
| -------------------- | -------- | ----------------- |
| Density Grid Terrain | Complete | 93.77%            |
| GenLayer Biomes      | Complete | 99.9% bit-level   |
| Villages             | Complete | 100%              |
| Desert Temple        | Complete | 100%              |
| Jungle Temple        | Complete | 100%              |
| Witch Hut            | Complete | 100%              |
| Abandoned Mineshaft  | Complete | 100%              |
| Stronghold           | Complete | 100%              |
| Caves                | Complete | SinTable-aligned  |
| Ravines              | Complete | SinTable-aligned  |
| 128-Height Squash    | Complete | MCPE 0.14 adapted |

#### Biome System

Fully ported biome and decorator system:

**Major Biomes:**
- Plains / Sunflower Plains
- Forest / Flower Forest
- Taiga / Cold Taiga
- Jungle / Jungle Edge
- Desert / Beach
- Savanna / Plateau
- Mesa / Mesa Plateau
- Roofed Forest
- Extreme Hills
- Swamp

**Decoration Generation:**
- Trees: Oak, Birch, Spruce, Pine, Acacia, Dark Oak, Mega Jungle, Mega Pine
- Ores: Coal, Iron, Gold, Redstone, Diamond, Lapis Lazuli (precise RNG sequences)
- Vegetation: Flowers, Grass, Mushrooms, Cacti, Sugar Cane, Lily Pads
- Terrain Features: Lakes, Dungeons, Ice Spikes

#### Network Protocol Layer

Full implementation of MCPE 0.14.3 (Protocol 70):

| Packet Category   | Count | Status   |
| ----------------- | ----- | -------- |
| Login/Auth        | 6     | Complete |
| Chunk Data        | 4     | Complete |
| Entity Management | 12    | Complete |
| Player Actions    | 10    | Complete |
| Items/Inventory   | 8     | Complete |
| World Events      | 8     | Complete |
| Other             | 15+   | Complete |

**Protocol Features:**
- BatchPacket (0x92) compression
- StartGame (0x95) full initialization
- RakNet reliable transport layer
- NBT little-endian serialization

#### Command System

45+ admin and player commands implemented:

| Category          | Commands                                                            |
| ----------------- | ------------------------------------------------------------------- |
| **Player Mgmt**   | `/ban`, `/ban-ip`, `/kick`, `/op`, `/deop`, `/whitelist`, `/pardon` |
| **Game Mode**     | `/gamemode`, `/defaultgamemode`, `/difficulty`                      |
| **Teleport/Loc**  | `/tp`, `/spawnpoint`, `/setworldspawn`                              |
| **Items/Effects** | `/give`, `/enchant`, `/effect`, `/xp`                               |
| **World Edit**    | `/setblock`, `/fill`, `/world_edit`                                 |
| **Server Mgmt**   | `/stop`, `/save`, `/time`, `/weather`, `/seed`                      |
| **Information**   | `/help`, `/list`, `/status`, `/version`, `/tps`, `/ping`            |
| **Communication** | `/say`, `/tell`, `/me`                                              |
| **Other**         | `/kill`, `/summon`, `/particle`, `/biome_find`, `/mw`               |

#### Lua Plugin System

Built-in Lua scripting engine for server extensibility:

- YAML-based plugin descriptors (`plugin.yml`)
- Event listener API (`events.listen`)
- Command registration API (`commands.register`)
- Player, Server, Level, Logger, and Scheduler APIs
- Plugin management commands (`/plugins`, `/luaplugin`)

**Example plugin structure:**
```
plugins/
  example/
    plugin.yml      # Plugin metadata
    main.lua        # Plugin entry point
```

#### Block and Item System

- Full block ID system (MCPE 0.14 compatible)
- Item metadata and NBT support
- Crafting recipe system framework

#### Entity System

- Entity base class (Entity)
- Living entities (Living/Mob)
- Human entities (Human)
- Item drop entities (ItemEntity)
- AABB collision detection
- Entity attributes system (Attributes)
- Entity metadata (Metadata)
- AI behavior framework

---

## Technical Verification

### World Generation Accuracy

Verified against seed `114514` across 1280 chunks (approximately 33.8 million blocks):

| Metric            | Result | Notes                          |
| ----------------- | ------ | ------------------------------ |
| **Block ID Diff** | 7.03%  | Primarily floating-point drift |
| **Biome Diff**    | 0.12%  | Near bit-level precision       |
| **Structure Pos** | 100%   | Exact match                    |

### Key Issues Resolved

- Ore RNG consumption (NextDouble to NextFloat)
- Surface depth coordinate scaling (16x multiplication fix)
- Chunk seeding alignment (Java 1.12 magic constants)
- Ravine math functions (65536-entry SinTable)
- Population seeding fix (dynamic multipliers)
- BFS light propagation engine
- Biome color alpha channel

---

## Quick Start

### Requirements

- Go 1.22 or later
- Windows / Linux

### Build and Run

```bash
# Clone the repository
git clone https://github.com/ScaxeTeam/SCAXE-GO.git
cd SCAXE-GO

# Build
go build -o scaxe-go ./cmd/server

# Run
./scaxe-go
```

### Command-Line Options

```
--version           Show version information and exit
--help              Show help message and exit
--config PATH       Path to server.properties (default: server.properties)
--debug             Enable debug logging (shows packet details)
--no-color          Disable colored output
```

### Configuration

`server.properties` key settings:

```properties
server-name=Scaxe Go Server
server-port=19132
server-ip=0.0.0.0
max-players=20
motd=A Scaxe Go Server
gamemode=0
difficulty=1
level-name=world
level-seed=
level-type=gorigional
online-mode=false
white-list=false
view-distance=8
pvp=true
```

---

## Project Structure

```
SCAXE-GO/
  cmd/
    server/               # Application entry point
  internal/
    version/              # Version constants
    wizard/               # First-run setup wizard
  pkg/
    block/                # Block system
    command/              # Command system
      defaults/           # Built-in commands (45+)
    config/               # Configuration loader
    crafting/             # Crafting recipe system
    entity/               # Entity system
      ai/                 # AI behaviors
      attribute/          # Attribute system
      effect/             # Status effects
    event/                # Event system
    inventory/            # Inventory system
    item/                 # Item system
    level/                # World / Level
      anvil/              # Anvil format I/O
      generator/          # World generator
        biome/            # Biome definitions (18)
        biomegrid/        # Biome grid mapping
        gorigional/       # Java 1.12 exact-port engine
          layer/          # GenLayer pipeline (24 layers)
          noise/          # Noise generators
          structure/      # Structure generation
        ground/           # Ground populators
        object/           # Decorators (35+)
        objects/          # Additional object types
        populator/        # Populator base
        populators/       # Populator implementations
    logger/               # Logging system
    lua/                  # Lua plugin engine
    math/                 # Math utilities
    nbt/                  # NBT serialization
    network/              # Network layer
    permission/           # Permission system
    player/               # Player management
    protocol/             # MCPE protocol (63 packet types)
    raknet/               # RakNet implementation
    scheduler/            # Task scheduler
    server/               # Server core
    world/                # World management
  plugins/
    example/              # Example Lua plugin
  logs/                   # Server log output
```

---

## Technical Highlights

### Hybrid Truth Model

To achieve bit-level accuracy, the project employs dual-source verification:

- **World Generation Algorithms** -- Minecraft Java 1.12.2 vanilla logic
- **Protocol and Physics** -- SCAXE PHP (MCPE 0.14 branch)

### 128-Height Squash Strategy

Adapted for the MCPE 0.14 128-block height limit:

| Parameter          | Vanilla (256h) | Squashed (128h)    |
| ------------------ | -------------- | ------------------ |
| Noise Segments (Y) | 33             | 17                 |
| StretchY           | 12.0           | 24.0               |
| Base Height        | base + noise   | (base + noise) / 2 |

### Concurrency Safety

- Fully stateless chunk generation
- Local buffer noise generation
- Thread-safe random number instances

---

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).

---

## Acknowledgments

- Minecraft Java Edition 1.12.2 - World generation logic reference
- SCAXE PHP / PocketMine-MP - MCPE protocol reference
- Go Community - Excellent toolchain support

---

<div align="center">

**Made by SCAXE Team**

</div>
