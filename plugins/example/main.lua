-- Example Plugin for SCAXE-GO
-- Demonstrates basic Lua plugin API usage

function onEnable()
    logger.info("Example plugin enabled!")
end

function onDisable()
    logger.info("Example plugin disabled!")
end

-- Broadcast when a player joins
events.listen("PlayerJoinEvent", function(e)
    server.broadcast("§a+ " .. e.playerName .. " joined the server")
end)

-- Broadcast when a player quits
events.listen("PlayerQuitEvent", function(e)
    server.broadcast("§c- " .. e.playerName .. " left the server")
end)

-- Register /online command
commands.register({
    name = "online",
    description = "Show online player count",
    usage = "/online",
    callback = function(sender, args)
        local count = server.getOnlineCount()
        local max = server.getMaxPlayers()
        player.sendMessage(sender.name, "§aOnline: §f" .. count .. "/" .. max)
    end
})

-- Register /tps command via Lua
commands.register({
    name = "luatps",
    description = "Show server TPS (Lua)",
    usage = "/luatps",
    callback = function(sender, args)
        local tps = server.getTPS()
        player.sendMessage(sender.name, string.format("§aTPS: §f%.2f", tps))
    end
})
