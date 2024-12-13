package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// create a table according the go struct
	//err = db.AutoMigrate(&Product{})
	//if err != nil {
	//	panic("failed to migrate table")
	//}

	// Create
	//db.Create(&Product{Code: "D42", Price: 100})
	var product Product

	// fetch product by id
	err = db.First(&product, 1).Error
	// db.First(&product, 1) // find product with integer primary key
	//  db.First(&product, "code = ?", "D42") // find product with code D42
	if err != nil {
		log.Fatal("Failed to find product:", err)
	}
	fmt.Println(product)
}

func createProduct(db *gorm.DB, code string, price uint) {
	product := Product{Code: code, Price: price}

	err := db.Create(&product).Error
	if err != nil {
		log.Fatal("Failed to create product:", err)
	}
	log.Printf("Created product: %+v\n", product)
}

func findProductByID(db *gorm.DB, id uint) {
	var product Product

	err := db.First(&product, id).Error
	if err != nil {
		log.Fatal("Failed to find product by ID:", err)
	}
	log.Printf("Found product by ID: %+v\n", product)
}

func findAllProducts(db *gorm.DB) {
	var products []Product
	err := db.Find(&products).Error
	if err != nil {
		log.Fatal("Failed to find all products:", err)
	}
	log.Printf("Found all products: %+v\n", products)
}

func updateProduct(db *gorm.DB, id uint, newCode string, newPrice uint) {
	var product Product
	err := db.First(&product, id).Error
	if err != nil {
		log.Fatal("Failed to find product for update:", err)
	}
	err = db.Model(&product).Updates(Product{Code: newCode, Price: newPrice}).Error
	if err != nil {
		log.Fatal("Failed to update product:", err)
	}
	log.Printf("Updated product: %+v\n", product)
}

func deleteProduct(db *gorm.DB, id uint) {
	err := db.Delete(&Product{}, id).Error
	if err != nil {
		log.Fatal("Failed to delete product:", err)
	}
	log.Printf("Deleted product with ID: %d\n", id)
}
