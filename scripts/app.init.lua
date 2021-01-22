#!/usr/bin/env tarantool

require('strict').on()
os = require('os')
log = require('log')
fiber = require('fiber')

box.cfg{listen = 3301}
-- Create a user, then init storage, create functions and then revoke ALL priveleges from user
local function init_storage(init_func, interface)
    init_func()

    for _, v in pairs(interface) do
        box.schema.func.create(v, {setuid = true, if_not_exists = true})
    end

    box.session.su('admin')
end

local function init()
    local space = box.schema.create_space('data_collector', {
        if_not_exists = true,
    })

    space:format({
        {name = 'event_id', type = 'unsigned'},
        {name = 'event_label', type = 'string'},
        {name = 'counter', type = 'unsigned'}
    })

    space:create_index('primary', {
        if_not_exists = true,
        type = 'TREE',
        parts = {
            'event_id',
        }
    })

    space:create_index('cccombo', {
        type = 'hash',
        unique = true,
        parts = {
            'event_id',
            'event_label',
        },
        if_not_exists = true,
    })
end

-- box.space.data_collector:select()
-- box.space.data_collector.index.cccombo:select({1, "a"})
function increment_counter (event_id, event_name)
    box.begin()
    local data = box.space.data_collector.index.cccombo:select({event_id, event_name})
    if len(data) == 0 then
        box.space.data_collector:insert({event_id, event_name, 1})
    else
        box.space.data_collector:update({event_id, event_name}, {"=", "=", 'value' + 1})
    end
    box.commit()
    return "ok"
end

-- box.space.data_collector.index.cccombo:pairs({1, "a"})
function load_counter ()
    return "123_load_counter"
end

local interface = {
    'increment_counter',
    'load_counter',
}

init_storage(init, interface)


log.info('Started tarantool for data_collector')
log.info('Started tarantool for data_collector')
log.info('Started tarantool for data_collector')
log.info('Started tarantool for data_collector')
log.info('Started tarantool for data_collector')
log.info('Started tarantool for data_collector')