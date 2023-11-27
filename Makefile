start-services:
	@echo "Starting services..."
	docker-compose up -d
	@echo "Services started successfully."

stop-services:
	@echo "Stopping services..."
	docker-compose down
	@echo "Services stopped successfully."
