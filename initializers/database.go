package initializers

import (
	"fmt"

	"github.com/ValianceTekProject/AreaBack/db"
)

var DB *db.PrismaClient

func ConnectDB() error {
	DB = db.NewClient()
	if err := DB.Prisma.Connect(); err != nil {
		fmt.Println("Failed to connect to database", err)
		return err
	}
	fmt.Println("Connected to database !")
	return nil
}
