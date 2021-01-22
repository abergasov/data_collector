## Testing 
```shell
# containers up
bash dev.sh
```

test 
```shell
wrk -s bench.lua -t12 -c400 -d30s http://localhost:29115/collect && curl http://localhost:29115/state
hey -m POST -d '{"id":123,"label":"view"}' -z 10s http://localhost:29115/collect && curl http://localhost:29115/state
```
![result](data_collector.jpg)