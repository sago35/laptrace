package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

const dbURL = "https://laptrace-live01-default-rtdb.asia-southeast1.firebasedatabase.app"

//const dbURL = "https://laptrace-live01.firebaseio.com"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <dbPath> <jsonFile> [debugN]")
		os.Exit(1)
	}

	dbPath := os.Args[1]
	filePath := os.Args[2]
	credFile := `serviceAccount.json`

	var debugN int
	if len(os.Args) >= 4 {
		n, err := strconv.Atoi(os.Args[3])
		if err == nil && n > 0 {
			debugN = n
		}
	}

	ctx := context.Background()

	// Firebase 初期化
	opt := option.WithCredentialsFile(credFile)
	app, err := firebase.NewApp(ctx, &firebase.Config{DatabaseURL: dbURL}, opt)
	if err != nil {
		panic(err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		panic(err)
	}

	ref := client.NewRef(dbPath)

	// JSON 読み込み
	raw, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var obj interface{}
	if err := json.Unmarshal(raw, &obj); err != nil {
		panic(err)
	}

	// debug record 切り詰め
	if debugN > 0 {
		if m, ok := obj.(map[string]interface{}); ok {
			if records, ok := m["Records"].([]interface{}); ok {
				if debugN < len(records) {
					m["Records"] = records[:debugN]
				}
			}
		}
		fmt.Printf("Debug mode: sending first %d records\n", debugN)
	}

	// 書き込み
	if err := ref.Set(ctx, obj); err != nil {
		panic(err)
	}

	fmt.Println("Write success")
}
