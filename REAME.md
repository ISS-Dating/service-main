## Running
- `docker-compose up postgres`
- `migrate -database "postgres://postgres:12345@127.0.0.1:5433/postgres?sslmode=disable" -path "<path/to/migrations>" up 1`
- `docker-compose up`

### Accessing main service through `0.0.0.0:8090`
- See the list of endpoints in `web/server.go`
- See the list of payloads in `web/request.go` and `model/model.go`
- Requests should contain json payload
- Send files in binary

### Accessing matcher service through `0.0.0.0:8091`
- See the list of endpoints in `web/matcher.go`
- See the list of payloads in `web/request.go` and `model/model.go`
- Requests should contain json payload
- Sequence of calls should be:
  - `/match`
  - `/answer`
  - `/status`

### Functions
- [x] Register
- [x] Login
- [x] Update
- [x] Search (mod)
- [x] Ban
- [x] Mod
- [x] Meet
- [x] ListFriends
- [ ] Report (TBD)
- [ ] Search report (TBD)