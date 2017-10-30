# telemetry
Telemetry is tiny telemetry collector and buffer to InfluxDB

#### Run influxdb, chronograf and grafana
  ```
  $cd example
  $docker-compose up -d
  ```
#### Connect Chronograf and Grafana to Influxdb
  ```
  use this url: http://influxdb:8086
  
  on Grafana you need to choose 'proxy' connection option to connect to Influxdb  
  ```
  
#### note 1: 
Wait for Grafana run migrations (it could take some time!) before try to open it on browser
Check if Grafana is ON by executing the following command
```
$docker logs grafana
```

#### note 2:
To login into Grafana:
``` 
 user: admin
 password: admin
```

#### note 3:
You can use the example.go to test all the stack.
But before you need to create a Database test

#### note 4:
On Browser
``` 
 chronograf: http://localhost:8888
 grafana: http://localhost:8081
```


