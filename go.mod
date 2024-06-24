module github.com/https-whoyan/MafiaBot

go 1.22

toolchain go1.22.4

require github.com/https-whoyan/MafiaBot/core v0.0.0

require (
	github.com/bwmarrin/discordgo v0.28.1
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.5.3
	github.com/stretchr/testify v1.9.0
	go.mongodb.org/mongo-driver v1.15.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240424034433-3c2c7870ae76 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect

)

replace github.com/https-whoyan/MafiaBot/core => ./pkg/core
