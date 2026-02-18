# SCAXE-GO

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![MCPE Version](https://img.shields.io/badge/MCPE-0.14.3-green?style=flat)
![Platform](https://img.shields.io/badge/Platform-Windows%20|%20Linux-blue?style=flat)

**A high-performance server core compatible with MCPE 0.14.3, built from scratch in Go**

*High-precision world generation engine based on Overworld core algorithm logic*

</div>

---

## Features

- **High-Precision World Generation** - Based on Overworld core algorithm logic, 93.77% terrain consistency
- **Bit-Level GenLayer Precision** - Biome system achieves 99.9% bit-level accuracy
- **1:1 Block Property Parity** - 182 registered blocks with properties matching PHP core exactly
- **High-Performance Concurrency** - Thread-safe chunk generation powered by Go goroutines
- **Full Protocol Implementation** - Complete MCPE 0.14.3 (Protocol 70) support
- **Lua Plugin System** - Extensible plugin architecture with Lua scripting and hot-reload
- **EULA Acceptance** - AGPL-3.0 license agreement on first startup

---

## Implemented Features

### Core Systems

#### Block System (v0.3.0 New)

Complete block property system with automated parity verification:

| Metric               | Result         |
| -------------------- | -------------- |
| Registered Blocks    | 182 / 256      |
| Properties per Block | 14             |
| Parity Test Result   | 182/182 (100%) |

**14 verified properties per block:**
Name, Hardness, BlastResistance, LightLevel, LightFilter, Solid, Transparent, Replaceable, ToolType, FlammableChance, BurnChance, DiffusesSkyLight, FuelTime, Flowable

#### World Generation Engine (Gorigional)

World generation engine based on Overworld core algorithm logic:

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

Complete biome and decorator system:

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

#### Entity System

- Entity base class (Entity)
- Living entities (Living/Mob)
- Human entities (Human)
- Item drop entities (ItemEntity)
- AABB collision detection
- Entity attributes system (Attributes)
- Entity metadata (Metadata)
- AI behavior framework

#### EULA System (v0.3.0 New)

- First-run AGPL-3.0 license display and acceptance prompt
- Acceptance state persisted to `eula.txt`
- Server refuses to start without license acceptance

---

## Technical Verification

### World Generation Accuracy

Verified against seed `114514` across 1280 chunks (approximately 33.8 million blocks):

| Metric            | Result | Notes                          |
| ----------------- | ------ | ------------------------------ |
| **Block ID Diff** | 7.03%  | Primarily floating-point drift |
| **Biome Diff**    | 0.12%  | Near bit-level precision       |
| **Structure Pos** | 100%   | Exact match                    |

### Block Property Parity

Automated parity test results:

| Metric              | Result    |
| ------------------- | --------- |
| **Blocks Tested**   | 182 / 182 |
| **Properties Each** | 14        |
| **Failures**        | 0         |

### Key Issues Resolved

- Ore RNG consumption (NextDouble to NextFloat)
- Surface depth coordinate scaling (16x multiplication fix)
- Chunk seeding alignment (magic constants)
- Ravine math functions (65536-entry SinTable)
- Population seeding fix (dynamic multipliers)
- BFS light propagation engine
- Biome color alpha channel

---

## Quick Start

### Requirements

- Windows / Linux

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

## Technical Highlights

### Algorithm Architecture

Built from scratch based on Overworld core algorithm logic:

- **World Generation** -- Based on Overworld core density grid and GenLayer pipeline
- **Protocol & Physics** -- Compatible with MCPE 0.14 protocol specification

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



<div align="center">

**Made by SCAXE Team**

</div>
