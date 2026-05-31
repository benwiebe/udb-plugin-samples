PLUGIN_NAME := udb-plugin-samples
OUTPUT := $(PLUGIN_NAME).so

.PHONY: build clean

build:
	go build -buildmode=plugin -o $(OUTPUT) .

clean:
	rm -f $(OUTPUT)
