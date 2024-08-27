.PHONY: dev

install:
	@echo "Installing all dependencies..."
	@go mod tidy

dev:
	@echo "Starting all services with hot reloading..."
	@trap 'kill %1; kill %2; kill %3' SIGINT; \
	templ generate --watch & \
	tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch & \
	air & \
	wait

# A command to stop all background processes if needed
stop:
	@echo "Stopping all background processes..."
	@pkill -f "templ generate --watch"
	@pkill -f "tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch"
	@pkill -f "air"
