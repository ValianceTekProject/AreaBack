package initializers

import (
	// "context"
	// "encoding/json"
	"fmt"

	// adapt "demo" to your module name if it differs
	"github.com/ValianceTekProject/AreaBack/db"
)

func ConnectDB() error {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		fmt.Println("❌ Erreur connexion DB:", err)
		return err
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()
	fmt.Println("✅ Connecté à la base de données !")
	return nil
}
