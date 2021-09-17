package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func main() {
	// 初期化
	ctx := context.Background()
	sa := option.WithCredentialsFile("../path/to/serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}


	// フィールド削除
	//_, errorDelete := client.Collection("users").Doc("user1").Update(ctx, []firestore.Update{
	//	{
	//		Path:  "first",
	//		Value: firestore.Delete,
	//	},
	//})
	//if errorDelete != nil {
	//	// Handle any errors in an appropriate way, such as returning them.
	//	log.Printf("An error has occurred: %s", err)
	//}

	// ドキュメント削除
	//_, errorDelete := client.Collection("users").Doc("uesr2").Delete(ctx)
	//if errorDelete != nil {
	//	// Handle any errors in an appropriate way, such as returning them.
	//	log.Printf("An error has occurred: %s", err)
	//}

	// コレクション削除
	ref := client.Collection("users")
	deleteCollection(ctx, client, ref, 10)

	// データ読み取り
	iter := client.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}
	// 切断
	defer client.Close()
}
func deleteCollection(ctx context.Context, client *firestore.Client,
	ref *firestore.CollectionRef, batchSize int) error {

	for {
		// バッチを参照し、イテレーターを取得
		iter := ref.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		// ドキュメントを再帰的に各ドキュメントの削除操作を行う
		// (Batchに追加)
		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		// 削除するドキュメントが無い場合、
		// プロセス終了
		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
}
