module github.com/https-whoyan/MafiaBot

go 1.21.1

require (
	github.com/https-whoyan/MafiaBot/core v0.0.0
)

require (
	github.com/bwmarrin/discordgo v0.28.1
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.5.1
	go.mongodb.org/mongo-driver v1.15.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect

)

replace github.com/https-whoyan/MafiaBot/core => ./pkg/core