
BINARY_NAME=pogoda
MAIN_FILE=.

.PHONY: all
all: build


.PHONY: build
build:
	go build -o $(BINARY_NAME) $(MAIN_FILE)


.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f *.png


.PHONY: deps
deps:
	go mod download
	go mod tidy


.PHONY: config
config:
	@if [ ! -f config.json ]; then \
		echo "Tworzenie domyślnego pliku konfiguracji..."; \
		echo '{\n  "threats": {\n    "high_temperature": 30.0,\n    "low_temperature": -15.0,\n    "high_wind_speed": 60.0,\n    "high_precipitation": 20.0,\n    "high_uv_index": 8.0\n  }\n}' > config.json; \
		echo "Utworzono config.json"; \
	else \
		echo "Plik config.json już istnieje"; \
	fi


.PHONY: fmt
fmt:
	go fmt ./...


.PHONY: vet
vet:
	go vet ./...


.PHONY: test
test:
	go test ./...


.PHONY: setup
setup: deps config
	@echo "Aplikacja gotowa do użycia!"


.PHONY: examples
examples:
	@echo "Przykłady użycia:"
	@echo "  make run aktualna Warszawa"
	@echo "  make run prognoza Kraków 5"
	@echo "  make run historia Gdańsk 2024-01-01,2024-01-07"


.PHONY: help
help:
	@echo "Dostępne komendy:"
	@echo "  make build    - Budowanie aplikacji"
	@echo "  make run      - Uruchamianie aplikacji"
	@echo "  make clean    - Czyszczenie plików"
	@echo "  make deps     - Instalowanie dependencji"
	@echo "  make config   - Tworzenie domyślnej konfiguracji"
	@echo "  make fmt      - Formatowanie kodu"
	@echo "  make vet      - Sprawdzanie kodu"
	@echo "  make test     - Uruchamianie testów"
	@echo "  make setup    - Pełne przygotowanie aplikacji"
	@echo "  make examples - Pokazanie przykładów użycia"
	@echo "  make install  - Instalacja w systemie"
	@echo "  make help     - Ta pomoc"
