package initializers

import (
	"context"
	// "encoding/json"
	"fmt"

	// adapt "demo" to your module name if it differs
	"github.com/ValianceTekProject/AreaBack/db"
)

func ConnectDB() error {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		fmt.Println("Failed to connect to database", err)
		return err
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()
	fmt.Println("Connected to database !")
	ctx := context.Background()
	client.Todo.CreateOne(
		db.Todo.Item.Set("Example"),
		db.Todo.Completed.Set(false),
	).Exec(ctx)
	return nil
}
