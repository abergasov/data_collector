create table counters
(
    event_id int null,
    event_label char(50) null,
    counter int default 0
);
create index conters_event_id_index on counters (event_id);
create index conters_event_label_index on counters (event_label);
create unique index conters_event_id_event_label_uindex on counters (event_id, event_label);
