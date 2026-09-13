EVENT_NETWORK := event_network
APP_NETWORK := app_network
HISTORY_NETWORK := history_network

.PHONY: network nats appdb historydb app history down clean

network:
	docker network inspect $(EVENT_NETWORK) >/dev/null 2>&1 || docker network create $(EVENT_NETWORK)
	docker network inspect $(APP_NETWORK) >/dev/null 2>&1 || docker network create $(APP_NETWORK)
	docker network inspect $(HISTORY_NETWORK) >/dev/null 2>&1 || docker network create $(HISTORY_NETWORK)

nats:
	docker compose -f docker-compose.nats.yml up

appdb:
	docker compose -f docker-compose.app-db.yml up

historydb:
	docker compose -f docker-compose.history-db.yml up

app:
	docker compose -f docker-compose.app.yml up

history:
	docker compose -f docker-compose.history.yml up

down:
	docker compose -f docker-compose.nats.yml down
	docker compose -f docker-compose.app-db.yml down
	docker compose -f docker-compose.history-db.yml down
	docker compose -f docker-compose.history.yml down
	docker compose -f docker-compose.app.yml down

clean: down
	docker compose -f docker-compose.app.yml build --no-cache
	docker compose -f docker-compose.history.yml build --no-cache
	docker network rm $(EVENT_NETWORK) $(APP_NETWORK) $(HISTORY_NETWORK) 2>/dev/null || true



# -------- The Order of Execution -------- #
# make network then make nats -> terminal 1
# make appdb -> terminal 2
# make historydb -> terminal3
# make app -> terminal4
# make history -> terminal5


# ----------- To Kill ---------- #
# make clean
