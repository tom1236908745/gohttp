package rel_post

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func firebaseInit(ctx context.Context) (*firestore.Client, error) {
	// Use a service account
	sa := option.WithCredentialsFile("../path/to/serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return client, nil
}

// firebase にデータを追加
func dataAdd(name string, age string, address string) error {

	ctx := context.Background()
	client, err := firebaseInit(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// データ追加
	_ , err = client.Collection("users").Doc(name).Set(ctx, map[string]interface{}{
		"age":     age,
		"address": address,
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}

	// データ読み込み
	allData := client.Collection("users").Documents(ctx)
	// 全ドキュメント取得
	docs, err := allData.GetAll()
	if err != nil {
		log.Fatalf("Failed adding getAll: %v", err)
	}

	// 配列の初期化
	users := make([]*User, 0)
	for _, doc := range docs {
		// 構造体の初期化
		u := new(User)
		// 構造体にFirestoreのデータをセット
		mapToStruct(doc.Data(), &u)
		// ドキュメント名を取得してnameにセット
		u.Name = doc.Ref.ID
		// 配列に構造体をセット
		users = append(users, u)
	}

	// 切断
	defer client.Close()

	// エラーなしは成功
	return err
}

// map -> 構造体の変換
func mapToStruct(m map[string]interface{}, val interface{}) error {
	tmp, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = json.Unmarshal(tmp, val)
	if err != nil {
		return err
	}
	return nil
}
