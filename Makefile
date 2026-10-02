.PHONY: help init update pull reset status clean

help:
	@echo "make init     - 初始化并拉取所有 submodules（首次使用）"
	@echo "make update   - 更新所有 submodules（拉取最新代码）"
	@echo "make pull     - 拉取主仓库 + 更新 submodules"
	@echo "make reset    - 强制重置 submodules（解决异常）"
	@echo "make status   - 查看 submodule 状态"
	@echo "make clean    - 清理 submodule（危险操作）"

# 🚀 初始化（首次 clone 后使用）
init:
	git submodule update --init --recursive

# 🔄 更新 submodule（保持最新）
update:
	git submodule update --remote --merge --recursive

# 🔄 拉主仓库 + 更新子模块
pull:
	git pull
	git submodule update --init --recursive
	git submodule update --remote --merge --recursive

# 🔧 强制修复（submodule 状态错乱时用）
reset:
	git submodule foreach --recursive git reset --hard
	git submodule foreach --recursive git clean -fd
	git submodule update --init --recursive

# 📊 查看状态
status:
	git submodule status

# ☠️ 清理（慎用！会删除 submodule 目录）
clean:
	git submodule deinit -f .
	rm -rf .git/modules/*
	rm -rf base-app base-engine base-web
