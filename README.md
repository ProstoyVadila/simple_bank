# Simple bank app

## Description:
This is a simple bank app with API.

## Purpose
Refresh my knowledge of Go and also get more understanding and practice of building a web service in Go.

## How to run
1. Via docker-compose:
```bash
docker-compose up
```
2. Via makefile:
```bash
make run
```

### TODO:
1. Refactor an error handling in some places.
2. Refactor integration db tests (add full cleanup, suits).
3. Add throttle middleware by ip.
4. Change primary key from username to uuid id for users.
5. Remove custom UUIDString type to uuid.UUID. (figure out with gin bug)
6. Improve docker compose file.
7. Add graceful shutdown.
8. Add more handlers and tests.
9. Figure out with trusted proxies.
10. Check sqlc for a more complicated queries.
11. Add monitoring.
12. Add swagger.
13. Add fieldalignment as git pre-hook?
