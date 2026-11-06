package main

import (
	"context"
	"fmt"
	"os"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/config"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/database"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/repository"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/service"
	"gorm.io/gorm"
)

func main() {
	run()
}

func run() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	db, err := database.NewDB(cfg)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), *cfg.DefaultDuration)
	defer cancel()

	switch os.Args[1] {
	case "migrate":
		migrate(db)
	case "seed":
		seed(ctx, db)
	case "drop-tables":
		dropTables(db)
	default:
		usage()
	}
}

func usage() {
	fmt.Println("args: migrate, seed, drop-tables")
}

func migrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		new(model.User),
		new(model.Role),
		new(model.Item),
		new(model.Tag),
	); err != nil {
		panic(fmt.Sprint("migrate error:", err))
	}
	fmt.Println("migrate success")
}

func seed(ctx context.Context, db *gorm.DB) {
	if err := db.Transaction(func(tx *gorm.DB) (err error) {
		roleRepo := repository.NewRoleRepo(tx)
		roleService := service.NewRoleService(roleRepo)
		roleIns := []model.RoleIn{
			{
				Name: new("crew"),
				Users: []model.UserIn{
					{
						Name:  new("crew"),
						Email: new(model.Email("crew@example.com")),
						Pw:    new(model.Pw("Crew123!")),
					},
					{
						Name:  new("crew2"),
						Email: new(model.Email("crew2@example.com")),
						Pw:    new(model.Pw("Crew123!")),
					},
				},
			},
			{
				Name: new("admin"),
				Users: []model.UserIn{
					{
						Name:  new("admin"),
						Email: new(model.Email("admin@example.com")),
						Pw:    new(model.Pw("Admin123!")),
					},
				},
			},
		}
		if _, err = roleService.Creates(ctx, roleIns); err != nil {
			return
		}

		tagRepo := repository.NewTagRepo(db)
		tagService := service.NewTagService(tagRepo)
		tagIns := []model.TagIn{
			{
				Name: new("consumables"),
				Items: []model.ItemIn{
					{Name: new("Printer Ink"), Price: new(uint(125000)), Quantity: new(uint(20))},
					{Name: new("Printer Toner"), Price: new(uint(350000)), Quantity: new(uint(10))},
					{Name: new("AA Batteries"), Price: new(uint(30000)), Quantity: new(uint(50))},
					{Name: new("Cleaning Solution"), Price: new(uint(45000)), Quantity: new(uint(15))},
					{Name: new("Hand Sanitizer"), Price: new(uint(25000)), Quantity: new(uint(30))},
				},
			},
			{
				Name: new("furniture"),
				Items: []model.ItemIn{
					{Name: new("Office Chair"), Price: new(uint(750000)), Quantity: new(uint(12))},
					{Name: new("Work Desk"), Price: new(uint(1500000)), Quantity: new(uint(8))},
					{Name: new("Filing Cabinet"), Price: new(uint(1200000)), Quantity: new(uint(6))},
					{Name: new("Bookshelf"), Price: new(uint(900000)), Quantity: new(uint(5))},
					{Name: new("Meeting Table"), Price: new(uint(2500000)), Quantity: new(uint(2))},
				},
			},
			{
				Name: new("peripheral"),
				Items: []model.ItemIn{
					{Name: new("Keyboard"), Price: new(uint(250000)), Quantity: new(uint(15))},
					{Name: new("Mouse"), Price: new(uint(150000)), Quantity: new(uint(20))},
					{Name: new("Monitor"), Price: new(uint(2200000)), Quantity: new(uint(7))},
					{Name: new("Webcam"), Price: new(uint(500000)), Quantity: new(uint(8))},
					{Name: new("USB Hub"), Price: new(uint(175000)), Quantity: new(uint(12))},
				},
			},
			{
				Name: new("stationery"),
				Items: []model.ItemIn{
					{Name: new("A4 Paper"), Price: new(uint(55000)), Quantity: new(uint(40))},
					{Name: new("Ballpoint Pen"), Price: new(uint(5000)), Quantity: new(uint(100))},
					{Name: new("Notebook"), Price: new(uint(15000)), Quantity: new(uint(60))},
					{Name: new("Sticky Notes"), Price: new(uint(12000)), Quantity: new(uint(35))},
					{Name: new("Stapler"), Price: new(uint(45000)), Quantity: new(uint(10))},
				},
			},
			{
				Name: new("vehicle"),
				Items: []model.ItemIn{
					{Name: new("Company Motorcycle"), Price: new(uint(18000000)), Quantity: new(uint(2))},
					{Name: new("Company Car"), Price: new(uint(250000000)), Quantity: new(uint(1))},
					{Name: new("Motorcycle Helmet"), Price: new(uint(350000)), Quantity: new(uint(4))},
					{Name: new("Spare Tire"), Price: new(uint(1200000)), Quantity: new(uint(4))},
					{Name: new("Vehicle Battery"), Price: new(uint(950000)), Quantity: new(uint(3))},
				},
			},
		}
		_, err = tagService.Creates(ctx, tagIns)
		return
	}); err != nil {
		panic(fmt.Sprintln("seed error:", err))
	}
	fmt.Println("seed success")
}

func dropTables(db *gorm.DB) {
	if err := db.Migrator().DropTable(
		new(model.User),
		new(model.Role),
		new(model.Item),
		new(model.Tag),
		"item_tags",
	); err != nil {
		panic(fmt.Sprintln("drop-tables err:", err))
	}
	fmt.Println("drop-tables success")
}
