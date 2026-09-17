GATEWAY_NETWORK := gateway_network
TODODB_NETWORK := tododb_network
EVENT_NETWORK := event_network
HISTORYDB_NETWORK := historydb_network

.PHONY: webup webdown webclean webnetwork webstorage webtodo webgateway webhistory



webdown:
	docker compose -f docker-compose.gateway.yml down
	docker compose -f docker-compose.todo.yml down
	docker compose -f docker-compose.todo-db.yml down
	docker compose -f docker-compose.history.yml down
	docker compose -f docker-compose.history-db.yml down
	docker compose -f docker-compose.kafka.yml down

webclean: webdown
	docker compose -f docker-compose.history.yml build --no-cache
	docker compose -f docker-compose.todo.yml build --no-cache
	docker compose -f docker-compose.gateway.yml build --no-cache
	docker network rm $(GATEWAY_NETWORK) 2>/dev/null || true
	docker network rm $(TODODB_NETWORK) 2>/dev/null || true
	docker network rm $(EVENT_NETWORK) 2>/dev/null || true
	docker network rm $(HISTORYDB_NETWORK) 2>/dev/null || true

webnetwork:
	docker network inspect $(EVENT_NETWORK) >/dev/null 2>&1 || docker network create $(EVENT_NETWORK)
	docker network inspect $(GATEWAY_NETWORK) >/dev/null 2>&1 || docker network create $(GATEWAY_NETWORK)
	docker network inspect $(TODODB_NETWORK) >/dev/null 2>&1 || docker network create $(TODODB_NETWORK)
	docker network inspect $(HISTORYDB_NETWORK) >/dev/null 2>&1 || docker network create $(HISTORYDB_NETWORK)

webstorage:
	docker compose -f docker-compose.kafka.yml up -d
	docker compose -f docker-compose.todo-db.yml up -d
	docker compose -f docker-compose.history-db.yml up -d

webhistory:
	docker compose -f docker-compose.history.yml up

webtodo:
	docker compose -f docker-compose.todo.yml up

webgateway:
	docker compose -f docker-compose.gateway.yml up


webup: webnetwork webstorage
	docker compose -f docker-compose.history.yml up -d
	docker compose -f docker-compose.todo.yml up -d
	docker compose -f docker-compose.gateway.yml up





