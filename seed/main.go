package main

import (
	"context"
	"log"

	"github.com/ValianceTekProject/AreaBack/db"
)

func create_google_seed(client *db.PrismaClient) {
	ctx := context.Background()

	_, err := client.Services.UpsertOne(
		db.Services.ID.Equals(1),
	).Create(
		db.Services.Name.Set("Google"),
	).Update(
		db.Services.Name.Set("Google"),
	).Exec(ctx)

	if err != nil {
		log.Fatalf("Error while seeding: %s", err)
	}
}

func create_github_seed(client *db.PrismaClient) {
	ctx := context.Background()

	_, err := client.Services.UpsertOne(
		db.Services.ID.Equals(2),
	).Create(
		db.Services.Name.Set("Github"),
	).Update(
		db.Services.Name.Set("Github"),
	).Exec(ctx)

	if err != nil {
		log.Fatalf("Error while seeding: %s", err)
	}
}

func main() {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Prisma.Disconnect()

	create_google_seed(client)
	create_github_seed(client)
}
