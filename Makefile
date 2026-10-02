.PHONY: help dev engine-migrate engine-bootstrap engine-dev engine-test web-dev web-build web-test web-lint test codegen

help:
	@echo "Base Hub 通用底座项目快捷命令："
	@echo "  make dev              - 提示前后端启动方式"
	@echo "  make engine-migrate   - 运行后端数据库迁移 (SQLite/MySQL)"
	@echo "  make engine-bootstrap - 初始化初始超级管理员账号"
	@echo "  make engine-dev       - 启动后端服务 (http://localhost:8085)"
	@echo "  make engine-test      - 运行后端所有单元与集成测试"
	@echo "  make web-dev          - 启动前端开发服务 (http://localhost:4321)"
	@echo "  make web-build        - 编译前端静态产物"
	@echo "  make web-test         - 运行前端单元测试"
	@echo "  make web-lint         - 运行前端代码规范检查"
	@echo "  make test             - 运行前后端全量测试"
	@echo "  make codegen          - 重新生成前后端 GraphQL 代码"

dev:
	@echo "请在两个独立终端分别执行："
	@echo "  终端 1: cd base-engine && go run . start --cors"
	@echo "  终端 2: cd base-web && pnpm dev"

engine-migrate:
	cd base-engine && go run . migrate

engine-bootstrap:
	cd base-engine && go run . bootstrap-admin

engine-dev:
	cd base-engine && go run . start --cors

engine-test:
	cd base-engine && go test ./...

web-dev:
	cd base-web && pnpm dev

web-build:
	cd base-web && pnpm build

web-test:
	cd base-web && pnpm test

web-lint:
	cd base-web && pnpm lint

test: engine-test web-test

codegen:
	cd base-engine && GO111MODULE=on go run github.com/sj-distributor/dolphin && go run ./tools/patchresolver
	cd base-web && pnpm codegen
