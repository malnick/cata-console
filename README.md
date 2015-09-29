# Kestrel Console 
Listens for POSTs from kestrel agents and displays host/docker information. Manages alerts and external services to alert to. 

## Development
### Requires

1. InfluxDB on branch 'influx'
  1. brew update && brew install influxdb
  1. See brew output for running influx DB 
1. Elasticsearch on branch 'elasticsearch' and 'master'
  1. brew install elasticsearch
  1. run it with your config.yml
  1. Also requires java 7x 

#### To run

1. git clone git@github.com/malnick/cata-console
1. cd cata-console
1. go run *.go
1. Open browser to localhost:9000

See config.go for available options (always changing) regarding runtime env. Most settings changed via env vars. 


#### Agent

1. git clone git@github.com/malnick/cata-agent
1. cd cata-agent
1. go run *.go

You'll see some info output and it'll connect to a console on localhost:9000 by default. Agent API available at :8080. 

#### Options

1. Both agent and console support -v for verbose output. Running in this mode is important since most useful debugging output occurs there. 
1. Agent also supports -p so you can switch the port it runs on at runtime. 
1. Always check config.go in both apps to see the latest env var config options
