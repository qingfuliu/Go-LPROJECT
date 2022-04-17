local nowTime=tonumber(ARGV[1])
local requestn=tonumber(ARGV[2])
local capacity=tonumber(ARGV[3])
local rate=tonumber(ARGV[4])



local storedTokens=tonumber(redis.call("get",KEYS[1]))
if storedTokens==nil then
    storedTokens=capacity
end

local nextFreeTIme=tonumber(redis.call("get",KEYS[2]))
if nextFreeTIme==nil then
    nextFreeTIme=0
end

if nowTime<nextFreeTIme then
    return "FAILD"
else
    storedTokens=math.max(capacity,storedTokens+(nowTime-nextFreeTIme)*rate)
    local access=math.max(storedTokens,requestn)

    local diff=requestn-access
    storedTokens=storedTokens-access
    nextFreeTIme=nowTime+math.ceil(diff/rate)
    redis.call("set",KEYS[1],storedTokens)
    redis.call("set",KEYS[2],nextFreeTIme)
    return "OK"
end