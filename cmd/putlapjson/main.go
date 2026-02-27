package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const dbURL = "https://laptrace-live01-default-rtdb.asia-southeast1.firebasedatabase.app"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <dbPath> <jsonFile> [debugN]")
		os.Exit(1)
	}

	dbPath := os.Args[1]
	filePath := os.Args[2]

	var debugN int
	if len(os.Args) >= 4 {
		n, err := strconv.Atoi(os.Args[3])
		if err == nil && n > 0 {
			debugN = n
		}
	}

	// ---- JSON 読み込み ----
	raw, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	if debugN > 0 {
		var obj map[string]interface{}
		if err := json.Unmarshal(raw, &obj); err != nil {
			panic(err)
		}

		if records, ok := obj["Records"].([]interface{}); ok {
			if debugN < len(records) {
				obj["Records"] = records[:debugN]
			}
		}

		raw, err = json.Marshal(obj)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Debug mode: sending first %d records\n", debugN)
	}

	// ---- PUT ----
	url := fmt.Sprintf("%s/%s.json", dbURL, dbPath)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(raw))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
}
