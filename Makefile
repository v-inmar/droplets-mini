GATEWAY_NETWORK := gateway_network
TODODB_NETWORK := tododb_network
EVENT_NETWORK := event_network
HISTORYDB_NETWORK := historydb_network

.PHONY: webup-network webup-storage webup-todo webup-history webup-gateway webup webdown webstatus webclean

# -- Individual / multi-terminal --

webup-network:
	docker network inspect $(EVENT_NETWORK) >/dev/null 2>&1 || docker network create $(EVENT_NETWORK)
	docker network inspect $(GATEWAY_NETWORK) >/dev/null 2>&1 || docker network create $(GATEWAY_NETWORK)
	docker network inspect $(TODODB_NETWORK) >/dev/null 2>&1 || docker network create $(TODODB_NETWORK)
	docker network inspect $(HISTORYDB_NETWORK) >/dev/null 2>&1 || docker network create $(HISTORYDB_NETWORK)

webup-storage:
	docker compose -f docker-compose.storage.yml up

webup-todo:
	docker compose -f services/todo/docker-compose.yml build --no-cache
	docker compose -f services/todo/docker-compose.yml up

webup-history:
	docker compose -f services/history/docker-compose.yml build --no-cache
	docker compose -f services/history/docker-compose.yml up

webup-gateway:
	docker compose -f services/gateway/docker-compose.yml build --no-cache
	docker compose -f services/gateway/docker-compose.yml up

# -- Single-terminal --

webup: webup-network
	docker compose -f docker-compose.storage.yml up -d
	docker compose -f services/todo/docker-compose.yml build --no-cache
	docker compose -f services/todo/docker-compose.yml up -d
	docker compose -f services/history/docker-compose.yml build --no-cache
	docker compose -f services/history/docker-compose.yml up -d
	docker compose -f services/gateway/docker-compose.yml build --no-cache
	docker compose -f services/gateway/docker-compose.yml up -d

webdown:
	docker compose -f services/gateway/docker-compose.yml down
	docker compose -f services/history/docker-compose.yml down
	docker compose -f services/todo/docker-compose.yml down
	docker compose -f docker-compose.storage.yml down

webstatus:
	docker compose -f docker-compose.storage.yml ps
	docker compose -f services/todo/docker-compose.yml ps
	docker compose -f services/history/docker-compose.yml ps
	docker compose -f services/gateway/docker-compose.yml ps

webclean: webdown
	docker compose -f services/todo/docker-compose.yml down --rmi local
	docker compose -f services/history/docker-compose.yml down --rmi local
	docker compose -f services/gateway/docker-compose.yml down --rmi local
	docker network rm $(EVENT_NETWORK) 2>/dev/null || true
	docker network rm $(GATEWAY_NETWORK) 2>/dev/null || true
	docker network rm $(TODODB_NETWORK) 2>/dev/null || true
	docker network rm $(HISTORYDB_NETWORK) 2>/dev/null || true