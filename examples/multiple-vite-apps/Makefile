.PHONY: start-frontend-app start-frontend-admin start-frontend-all start-backend start-all build-frontend-app build-frontend-admin build-frontend-all build-backend build-all

start-frontend-app:	
	cd frontend/app && npm run dev

start-frontend-admin:
	cd frontend/admin-app && VITE_BASE_URL=/admin/ npm run dev

start-frontend-all:
	$(MAKE) --no-print-directory -j2 start-frontend-app start-frontend-admin

start-backend:
	go run main.go config.go render.go

start-all:
	$(MAKE) --no-print-directory -j2 start-frontend-all start-backend

build-frontend-app:
	cd frontend/app && npm run build

build-frontend-admin:
	cd frontend/admin-app && npm run build

build-frontend-all:
	$(MAKE) --no-print-directory -j2 build-frontend-app build-frontend-admin

build-backend:
	go build -o backend main.go config.go render.go

build-all:
	$(MAKE) --no-print-directory -j2 build-frontend-all build-backend