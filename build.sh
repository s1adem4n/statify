cd frontend
bun install
bun run build

cd ..
go build -o build/statify main.go
