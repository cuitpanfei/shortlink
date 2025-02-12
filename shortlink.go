package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var url_cache = make(map[string]string, 1000)

func updateCache(data map[string]string) {
	url_cache = data
}

func handler(w http.ResponseWriter, r *http.Request) {
	target := getRedirectTarget(r.URL.Path)
	if target != "" {
		http.Redirect(w, r, target, http.StatusFound)
	} else {
		http.NotFound(w, r)
	}
}

func getRedirectTarget(abbr string) string {
	abbr = strings.TrimPrefix(abbr, "/")
	if url, ok := url_cache[abbr]; ok {
		return url
	} else {
		log.Print("Unable to find mappings at: ", getMappingFile())
		return ""
	}
}

func getMappingFile() string {
	return getEnv("SL_MAPPING", DefaultMapping)
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func startServer() {
	port := fmt.Sprintf(":%d", getEnvInt("SL_PORT", DefaultPort))
	log.Print("Starting shortlink server on ", port)
	log.Print("Using mappings: ", getMappingFile())
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(port, nil))
}
