package main

import (
	"context"
	"log"

	"github.com/ValianceTekProject/AreaBack/db"
)

func Create_github_pr_to_discord_message(client *db.PrismaClient) {
	// ctx := context.Background()
	//
	// githubService, err := client.Services.FindFirst(
	// 	db.Services.Name.Equals("Github"),
	// ).Exec(ctx)
	// if err != nil {
	// 	log.Fatalf("Github service not found: %s", err)
	// }
	//
	// discordService, err := client.Services.FindFirst(
	// 	db.Services.Name.Equals("Discord"),
	// ).Exec(ctx)
	// if err != nil {
	// 	log.Fatalf("Discord service not found: %s", err)
	// }
	//
	// area, err := client.Areas.UpsertOne(
	// 	db.Areas.Name.Equals("Github_pr_to_discord"),
	// ).Create(
	// 	db.Areas.Name.Set("Github_pr_to_discord"),
	// ).Update().Exec(ctx)
	// if err != nil {
	// 	log.Fatalf("Error while upserting area: %s", err)
	// }
	//
	// existingAction, _ := client.Actions.FindFirst(
	// 	db.Actions.And(
	// 		db.Actions.AreaID.Equals(area.ID),
	// 		db.Actions.ServiceID.Equals(githubService.ID),
	// 	),
	// ).Exec(ctx)
	//
	// if existingAction == nil {
	// 	_, err = client.Actions.CreateOne(
	// 		db.Actions.Triggered.Set(false),
	// 		db.Actions.Area.Link(
	// 			db.Areas.ID.Equals(area.ID),
	// 		),
	// 		db.Actions.Service.Link(
	// 			db.Services.ID.Equals(githubService.ID),
	// 		),
	// 	).Exec(ctx)
	// 	if err != nil {
	// 		log.Fatalf("Error creating action: %s", err)
	// 	}
	// }
	// existingReaction, _ := client.Reactions.FindFirst(
	// 	db.Reactions.And(
	// 		db.Reactions.AreaID.Equals(area.ID),
	// 		db.Reactions.ServiceID.Equals(discordService.ID),
	// 	),
	// ).Exec(ctx)
	// 	if existingReaction == nil {
	// 	_, err = client.Reactions.CreateOne(
	// 		db.Reactions.Area.Link(
	// 			db.Areas.ID.Equals(area.ID),
	// 		),
	// 		db.Reactions.Service.Link(
	// 			db.Services.ID.Equals(discordService.ID),
	// 		),
	// 	).Exec(ctx)
	// 	if err != nil {
	// 		log.Fatalf("Error creating reaction: %s", err)
	// 	}
	// }
}

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

func create_discord_seed(client *db.PrismaClient) {
	ctx := context.Background()

	_, err := client.Services.UpsertOne(
		db.Services.ID.Equals(3),
	).Create(
		db.Services.Name.Set("Discord"),
	).Update(
		db.Services.Name.Set("Discord"),
	).Exec(ctx)

	if err != nil {
		log.Fatalf("Error while seeding: %s", err)
	}
}

func create_steam_seed(client *db.PrismaClient) {
	ctx := context.Background()

	_, err := client.Services.UpsertOne(
		db.Services.ID.Equals(4),
	).Create(
		db.Services.Name.Set("Steam"),
	).Update(
		db.Services.Name.Set("Steam"),
	).Exec(ctx)

	if err != nil {
		log.Fatalf("Error while seeding: %s", err)
	}
}

func create_twitch_seed(client *db.PrismaClient) {
	ctx := context.Background()

	_, err := client.Services.UpsertOne(
		db.Services.ID.Equals(4),
	).Create(
		db.Services.Name.Set("Twitch"),
	).Update(
		db.Services.Name.Set("Twitch"),
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
	create_discord_seed(client)
	create_steam_seed(client)
	create_twitch_seed(client)
	// Create_github_pr_to_discord_message(client)
}
