.PHONY: run migrate-down migrate-down-all help

# Chạy server bình thường
run:
	go run .

# Rollback tất cả migrations
migrate-down-all:
	go run main.go -migrate-down -limit=100

# Rollback số lượng cụ thể migrations (dùng: make migrate-down-n LIMIT=3)
migrate-down-n:
	go run main.go -migrate-down -limit=$(LIMIT)

# Hiển thị hướng dẫn
help:
	@echo "Available commands:"
	@echo "  make run              - Chạy server"
	@echo "  make migrate-down-all - Rollback tất cả migrations"
	@echo "  make migrate-down-n LIMIT=3 - Rollback số lượng cụ thể"
