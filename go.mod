module github.com/https-whoyan/MafiaBot

go 1.22

toolchain go1.22.4

require (
	github.com/bwmarrin/discordgo v0.28.1
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.6.1
	go.mongodb.org/mongo-driver v1.17.1
)

require (
	github.com/https-whoyan/MafiaCore v0.0.1
	github.com/samber/lo v1.47.0
)

require (
	github.com/LastPossum/kamino v0.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/klauspost/compress v1.17.10 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)

replace github.com/https-whoyan/MafiaCore v0.0.1 => ../MafiaCore
