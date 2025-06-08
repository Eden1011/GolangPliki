package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("nie można otworzyć pliku konfiguracji: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("błąd parsowania pliku konfiguracji: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("nieprawidłowa konfiguracja: %w", err)
	}

	return &config, nil
}

func SaveConfig(config *Config, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("nie można utworzyć pliku konfiguracji: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(config)
	if err != nil {
		return fmt.Errorf("błąd zapisywania konfiguracji: %w", err)
	}

	return nil
}

func GetDefaultConfig() *Config {
	config := &Config{}

	config.Threats.HighTemperature = 30.0
	config.Threats.LowTemperature = -15.0
	config.Threats.HighWindSpeed = 60.0
	config.Threats.HighPrecipitation = 20.0
	config.Threats.HighUVIndex = 8.0

	return config
}

func validateConfig(config *Config) error {
	if config.Threats.HighTemperature <= config.Threats.LowTemperature {
		return fmt.Errorf("próg wysokiej temperatury (%.1f) musi być wyższy od progu niskiej temperatury (%.1f)",
			config.Threats.HighTemperature, config.Threats.LowTemperature)
	}

	if config.Threats.HighTemperature < -50 || config.Threats.HighTemperature > 60 {
		return fmt.Errorf("próg wysokiej temperatury (%.1f) jest poza rozsądnym zakresem (-50 do 60°C)",
			config.Threats.HighTemperature)
	}

	if config.Threats.LowTemperature < -60 || config.Threats.LowTemperature > 40 {
		return fmt.Errorf("próg niskiej temperatury (%.1f) jest poza rozsądnym zakresem (-60 do 40°C)",
			config.Threats.LowTemperature)
	}

	if config.Threats.HighWindSpeed <= 0 || config.Threats.HighWindSpeed > 300 {
		return fmt.Errorf("próg silnego wiatru (%.1f) musi być między 0 a 300 km/h",
			config.Threats.HighWindSpeed)
	}

	if config.Threats.HighPrecipitation <= 0 || config.Threats.HighPrecipitation > 500 {
		return fmt.Errorf("próg intensywnych opadów (%.1f) musi być między 0 a 500 mm",
			config.Threats.HighPrecipitation)
	}

	if config.Threats.HighUVIndex <= 0 || config.Threats.HighUVIndex > 15 {
		return fmt.Errorf("próg wysokiego UV (%.1f) musi być między 0 a 15",
			config.Threats.HighUVIndex)
	}

	return nil
}

func CreateDefaultConfigFile(filename string) error {
	config := GetDefaultConfig()
	return SaveConfig(config, filename)
}

func DisplayConfig(config *Config) {
	fmt.Println("AKTUALNA KONFIGURACJA:")
	fmt.Printf("Próg wysokiej temperatury: %.1f°C\n", config.Threats.HighTemperature)
	fmt.Printf("Próg niskiej temperatury: %.1f°C\n", config.Threats.LowTemperature)
	fmt.Printf("Próg silnego wiatru: %.1f km/h\n", config.Threats.HighWindSpeed)
	fmt.Printf("Próg intensywnych opadów: %.1f mm\n", config.Threats.HighPrecipitation)
	fmt.Printf("Próg wysokiego UV: %.1f\n", config.Threats.HighUVIndex)
	fmt.Println()
}

func UpdateConfig(config *Config) error {
	fmt.Println("Aktualizacja konfiguracji - naciśnij Enter aby zachować aktualną wartość:")

	getValue := func(prompt string, currentValue float64) float64 {
		fmt.Printf("%s (aktualna: %.1f): ", prompt, currentValue)
		var input string
		fmt.Scanln(&input)
		if input == "" {
			return currentValue
		}

		var newValue float64
		if _, err := fmt.Sscanf(input, "%f", &newValue); err != nil {
			fmt.Printf("Nieprawidłowa wartość, zachowuję %.1f\n", currentValue)
			return currentValue
		}
		return newValue
	}

	config.Threats.HighTemperature = getValue("Próg wysokiej temperatury (°C)", config.Threats.HighTemperature)
	config.Threats.LowTemperature = getValue("Próg niskiej temperatury (°C)", config.Threats.LowTemperature)
	config.Threats.HighWindSpeed = getValue("Próg silnego wiatru (km/h)", config.Threats.HighWindSpeed)
	config.Threats.HighPrecipitation = getValue("Próg intensywnych opadów (mm)", config.Threats.HighPrecipitation)
	config.Threats.HighUVIndex = getValue("Próg wysokiego UV", config.Threats.HighUVIndex)

	if err := validateConfig(config); err != nil {
		return fmt.Errorf("nieprawidłowa konfiguracja: %w", err)
	}

	return nil
}
