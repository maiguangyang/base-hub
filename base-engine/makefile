# DATABASE = sqlite3://test.db
DATABASE = mysql://'root:123456@localhost:3306/base_engine?charset=utf8mb4&parseTime=True&loc=Local'

PORT = 8085

mysql:
	docker run --restart=always -itd --name mysql -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql

redis:
	docker run --restart=always -itd --name redis -p 6379:6379 redis

docker:
	docker run --rm -it \
	--link mysql:db \
	--name golang \
	-e GOPATH=/data/go \
	-e GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
	-v ~/docker/go:/data/go \
	-v ~/wwwroot:/data/wwwroot \
	-p $(PORT):$(PORT) \
	golang bash

generate:
	GO111MODULE=on go run github.com/sj-distributor/dolphin
	go run ./tools/patchresolver

# 只读取手写 Go 源码；检查器跳过 gen/，本目标不触发代码生成。
check_structure:
	go run ./tools/sourcecheck -root .

# 接口清单、可调用工具、模型输入、授权与受保护业务路径的一致性检查。
check_ai_contracts:
	go test . ./src/services/ai ./src/services/ai/tools ./src/integration -count=1

# 只读加载当前已启用的后台模型配置，用合成工具验证五次连续流式调用。
ai_live_test:
	@set -a; . ./.env.local; set +a; DATABASE_URL=$(DATABASE) AI_LIVE_TEST_FROM_STORE=true go test ./src/services/ai -run '^TestChatCompletionsLive$$' -count=1 -v

ai_live_preview_test:
	@set -a; . ./.env.local; set +a; DATABASE_URL=$(DATABASE) AI_LIVE_TEST_FROM_STORE=true go test ./src/integration -run '^TestAILiveProtectedBusinessRead$$' -count=1 -v

ai_live_permissions_test:
	@set -a; . ./.env.local; set +a; DATABASE_URL=$(DATABASE) AI_LIVE_TEST_FROM_STORE=true go test ./src/integration -run '^TestAILiveAccountPermissionsQuestion$$' -count=1 -v

migrate:
	DATABASE_URL=$(DATABASE) PORT=$(PORT) go run . migrate

start:
	@if [ -f .env.local ]; then set -a; . ./.env.local; set +a; fi; AUTH_TOKEN_ISSUER="$${AUTH_TOKEN_ISSUER:-base-engine-local}" AUTH_COOKIE_SECURE="$${AUTH_COOKIE_SECURE:-false}" ALLOWED_ORIGINS="$${ALLOWED_ORIGINS:-http://localhost:4321,http://localhost:4322,http://127.0.0.1:4321,http://127.0.0.1:4322}" DATABASE_URL=$(DATABASE) PORT=$(PORT) DEBUG="true" go run . start --cors
