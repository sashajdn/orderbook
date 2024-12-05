PROJ_DIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
TOOL_DIR=${PROJ_DIR}/tools
BIN_DIR=${PROJ_DIR}/bin
DOCKER_COMPOSE_CMD=$(shell command -v docker-compose >/dev/null 2>&1 && echo "docker-compose" || echo "docker compose")


.PHONY:help
help: ## show self documenting help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
