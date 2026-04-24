-- lua/get-and-delete.lua

local val = redis.call("GET", KEYS[1])

if not val then
    return nil
end

redis.call("DEL", KEYS[1])

return val