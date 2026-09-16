EVENT_NETWORK := event_network
APP_NETWORK := app_network
HISTORY_NETWORK := history_network

GATEWAY_NETWORK := gateway_network
TODODB_NETWORK := tododb_network

.PHONY: network nats appdb historydb app history down clean webup webdown webclean webnetwork webpers webtodo webgateway

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

webup:
	docker network inspect $(EVENT_NETWORK) >/dev/null 2>&1 || docker network create $(EVENT_NETWORK)
	docker network inspect $(GATEWAY_NETWORK) >/dev/null 2>&1 || docker network create $(GATEWAY_NETWORK)
	docker network inspect $(TODODB_NETWORK) >/dev/null 2>&1 || docker network create $(TODODB_NETWORK)
	docker compose -f docker-compose.kafka.yml up -d
	docker compose -f docker-compose.todo-db.yml up -d
	docker compose -f docker-compose.todo.yml up -d
	docker compose -f docker-compose.gateway.yml up

webdown:
	docker compose -f docker-compose.gateway.yml down
	docker compose -f docker-compose.todo.yml down
	docker compose -f docker-compose.todo-db.yml down
	docker compose -f docker-compose.kafka.yml down

webclean: webdown
	docker compose -f docker-compose.todo.yml build --no-cache
	docker compose -f docker-compose.gateway.yml build --no-cache
	docker network rm $(GATEWAY_NETWORK) 2>/dev/null || true
	docker network rm $(TODODB_NETWORK) 2>/dev/null || true
	docker network rm $(EVENT_NETWORK) 2>/dev/null || true

webnetwork:
	docker network inspect $(EVENT_NETWORK) >/dev/null 2>&1 || docker network create $(EVENT_NETWORK)
	docker network inspect $(GATEWAY_NETWORK) >/dev/null 2>&1 || docker network create $(GATEWAY_NETWORK)
	docker network inspect $(TODODB_NETWORK) >/dev/null 2>&1 || docker network create $(TODODB_NETWORK)

webpers:
	docker compose -f docker-compose.kafka.yml up -d
	docker compose -f docker-compose.todo-db.yml up -d

webtodo:
	docker compose -f docker-compose.todo.yml up

webgateway:
	docker compose -f docker-compose.gateway.yml up










# -------- The Order of Execution -------- #
# make network then make nats -> terminal 1
# make appdb -> terminal 2
# make historydb -> terminal3
# make app -> terminal4
# make history -> terminal5


# ----------- To Kill ---------- #
# make clean
