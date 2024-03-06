module server

go 1.22

replace github.com/gabereiser/swr => ./swr

require github.com/gabereiser/swr v0.0.0-00010101000000-000000000000

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.14 // indirect
	github.com/robertkrimen/otto v0.0.0-20211024170158-b87d35c0b86f // indirect
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/sqlite v1.3.6 // indirect
	gorm.io/gorm v1.23.8 // indirect
)
