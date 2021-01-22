math.randomseed(os.time())
number = math.random(10, 100)

wrk.method = "POST"
-- wrk.body = '{"id":' .. math.random(10, 100) .. ',"label":"' .. tostring(math.random(10, 100)) .. '"}'
wrk.headers["Content-Type"] = "application/json"


request = function()
number = math.random(10, 100)
wrk.body = '{"id":' .. math.random(10, 100) .. ',"label":"' .. tostring(math.random(10, 100)) .. '"}'
return wrk.format(nil, path)
end