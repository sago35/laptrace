package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/fsnotify/fsnotify"
	"google.golang.org/api/option"
)

const dbURL = "https://laptrace-live01-default-rtdb.asia-southeast1.firebasedatabase.app"

//const dbURL = "https://laptrace-live01.firebaseio.com"

func main() {
	oneshot := flag.Bool("oneshot", false, "oneshot mode")
	flag.Parse()

	args := flag.Args()
	if len(args) < 3 {
		fmt.Println("Usage: go run main.go <dbPath> <jsonFile> [debugN]")
		os.Exit(1)
	}

	dbPath := args[0]
	filePath := args[1]
	credFile := `serviceAccount.json`

	var debugN int
	if len(args) >= 3 {
		n, err := strconv.Atoi(args[2])
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

	// 初回送信
	sendJSON(ctx, ref, filePath, debugN)

	if *oneshot {
		return
	}

	// ファイル監視
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	err = watcher.Add(filepath.Dir(filePath))
	if err != nil {
		panic(err)
	}

	fmt.Println("Watching file:", filePath)
	fmt.Println("Press Ctrl+C to exit")

	// Ctrl+C 捕捉
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case event := <-watcher.Events:
			fmt.Printf("%#v\n", event)
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
				fmt.Println("File changed, reloading...")
				sendJSON(ctx, ref, filePath, debugN)
			}

		case err := <-watcher.Errors:
			fmt.Println("Watcher error:", err)

		case <-sigChan:
			fmt.Println("Exit")
			return
		}
	}
}

func sendJSON(ctx context.Context, ref *db.Ref, filePath string, debugN int) {

	raw, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	var obj interface{}
	if err := json.Unmarshal(raw, &obj); err != nil {
		fmt.Println("JSON error:", err)
		return
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
		fmt.Println("Write error:", err)
		return
	}

	fmt.Println("Write success")
}
