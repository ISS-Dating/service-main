## Running
- `docker-compose up postgres`
- `migrate -database "postgres://postgres:12345@127.0.0.1:5433/postgres?sslmode=disable" -path "<path/to/migrations>" up 1`
- `docker-compose up`

### Accessing
- `0.0.0.0:8090`

### Functions
- [x] Register
- [x] Login
- [x] Update
- [ ] Search (user)
- [x] Search (mod)
- [ ] Ban
- [ ] Mod
- [ ] Meet
- [ ] Report
- [ ] Search report