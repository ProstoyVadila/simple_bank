# Simple bank app
### TODO:
1. Refactor an error handling in some places.
2. Refactor integration db tests (add full cleanup, suits).
3. Add throttle middleware by ip.
4. Change primary key from username to uuid id for users.
5. Remove custom UUIDString type to uuid.UUID. (figure out with gin bug)
6. Add graceful shutdown.
7. Add more handlers and tests.
8. Check sqlc for a more complicated queries (with join for example).
9. Add monitoring feature.
