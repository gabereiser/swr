module server

go 1.18

replace github.com/gabereiser/swr => ./swr

require github.com/gabereiser/swr v0.0.0-00010101000000-000000000000

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/robertkrimen/otto v0.0.0-20211024170158-b87d35c0b86f // indirect
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
