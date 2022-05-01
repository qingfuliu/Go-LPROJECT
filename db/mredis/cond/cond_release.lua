--KEYS   1               2
--      zSetKey        unique id
local rs = tonumber(redis.call("zRem", KEYS[1], KEYS[2]))
if (rs == nil) or (rs~=1) then
    return "FAILED"
end
return "OK"