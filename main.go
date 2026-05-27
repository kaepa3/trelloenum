package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// TrelloからエクスポートしたJSONの構造に合わせた定義
type TrelloExport struct {
	Name  string `json:"name"`  // ボード名
	Lists []List `json:"lists"` // リスト一覧
	Cards []Card `json:"cards"` // カード一覧
}

type List struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Card struct {
	ID     string `json:"id"`
	IDList string `json:"idList"` // どのリストに所属しているか
	Name   string `json:"name"`
	Desc   string `json:"desc"`
}

func main() {
	// 1. 引数の数をチェック (プログラム名 + 第一引数 なので、2未満ならエラー)
	if len(os.Args) < 2 {
		fmt.Println("❌ 使い方: go run main.go [JSONファイル名]")
		fmt.Println("例: go run main.go trello.json")
		return
	}

	// 2. 第一引数からファイル名を取得
	filename := os.Args[1]
	// 1. エクスポートしたJSONファイルを開く (DB接続のようなもの)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("ファイルが開けませんでした: %v\n", err)
		return
	}
	defer file.Close()

	// 2. 構造体にデコードする
	var data TrelloExport
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		fmt.Printf("JSONのパースに失敗しました: %v\n", err)
		return
	}

	fmt.Printf("📦 ボード名: %s のデータを読み込みました\n", data.Name)

	// 3. リストごとにループを回し、所属するカードを紐づけて表示
	for _, list := range data.Lists {
		fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n📂 リスト: %s\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n", list.Name)

		hasCard := false
		for _, card := range data.Cards {
			// カードの所属リストIDが、現在のリストIDと一致するか（SQLのJOINのような処理）
			if card.IDList == list.ID {
				hasCard = true
				fmt.Printf("📌 [%s]\n", card.Name)
				if card.Desc != "" {
					fmt.Printf("   説明: %s\n", card.Desc)
				}
				fmt.Println()
			}
		}

		if !hasCard {
			fmt.Println("  (カードはありません)")
		}
	}
}
