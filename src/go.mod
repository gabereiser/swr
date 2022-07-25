module server

go 1.18

require github.com/gabereiser/swr v0.0.0

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gabereiser/swr => ./swr
