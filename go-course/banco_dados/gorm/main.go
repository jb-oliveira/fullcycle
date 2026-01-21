package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // Import the driver, aliased with _ to run its init function
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type UuidModel struct {
	// ID        string `gorm:"primarykey;default:uuidv7()::varchar"`
	ID        uuid.UUID `gorm:"type:uuid;primarykey;default:uuidv7()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Category struct {
	UuidModel
	Name     string
	Products []Product
}

func (Category) TableName() string {
	return "catalogo.categorias"
}

type Product struct {
	UuidModel
	Name         string
	Price        float64
	CategoryID   uuid.UUID
	Category     Category
	SerialNumber SerialNumber
}

type SerialNumber struct {
	gorm.Model
	Number    string
	ProductID uuid.UUID
}

func (SerialNumber) TableName() string {
	return "catalogo.numero_seriais"
}

// func (u *Produto) BeforeCreate(tx *gorm.DB) error {
// 	// Checks if the ID is already filled (useful for tests or manual inserts)
// 	if u.ID == "" {
// 		uuid, err := uuid.NewV7()
// 		if err != nil {
// 			return err
// 		}
// 		// Generates a UUID v7 and converts it to string
// 		u.ID = uuid.String()
// 		// If you are using a specific library for V7:
// 		// u.ID = uuid.NewV7().String()
// 	}
// 	return nil
// }

func (Product) TableName() string {
	return "catalogo.produtos"
}

func main() {

	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	// 	logger.Config{
	// 		SlowThreshold:             time.Second, // Slow SQL threshold
	// 		LogLevel:                  logger.Info, // Log level
	// 		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
	// 		Colorful:                  true,        // Disable color
	// 	},
	// )

	connStr := "user=postgres password=password dbname=myapp host=localhost port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info, // Must be Info to see all SQLs, including associations
			},
		),
	})
	if err != nil {
		panic("failed to connect database")
	}

	ctx := context.Background()

	// Migrate the schema
	// err = db.AutoMigrate(&Category{}, &Product{}, &SerialNumber{})
	// if err != nil {
	// 	panic(err)
	// }

	// prod := createCategoryAndProduct(db, ctx, "Category 1", "Product 1", 34.87)

	// gorm.G[SerialNumber](db).Create(ctx, &SerialNumber{
	// 	Number:    "12345",
	// 	ProductID: prod.ID,
	// })

	//testProduct(db, ctx)

	// cats, err := gorm.G[Category](db).Find(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// for i := range 10 {
	// 	product := Product{Name: fmt.Sprintf("Product : %d", i), Price: 100.45, CategoryID: cats[0].ID}
	// 	err = gorm.G[Product](db).Create(ctx, &product)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// categories, err := gorm.G[Category](db).
	// 	Preload("Products", nil).
	// 	Preload("Products.SerialNumber", nil).
	// 	Find(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, cat := range categories {
	// 	fmt.Printf("Category: %s", cat.Name)
	// 	for _, p := range cat.Products {
	// 		fmt.Printf("-- Product = %s, SerialNumber = %s \n", p.Name, p.SerialNumber.Number)
	// 	}
	// }

	prods, err := gorm.G[Product](db).
		Joins(clause.Has("Category"), nil).
		Joins(clause.Has("SerialNumber"), nil).
		Find(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total products %d", len(prods))

	// err = gorm.G[Product](db).
	// 	Joins(clause.InnerJoin.Association("Category"), nil).
	// 	Joins(clause.InnerJoin.Association("SerialNumber"), nil).
	// 	FindInBatches(ctx, 10000, func(data []Product, batch int) error {
	// 		fmt.Printf("Total products: %d", len(data))
	// 		for _, prod := range data {
	// 			str := prod.ID.String()
	// 			if str == "" {
	// 				return errors.ErrUnsupported
	// 			}
	// 		}
	// 		return nil
	// 	})
	// if err != nil {
	// 	panic(err)
	// }

	// product := Product{Name: "Test 5", Price: 1234.78}

	// err = gorm.G[Product](db).Create(ctx, &product)
	// if err != nil {
	// 	panic(err)
	// }

	// products := []Product{
	// 	{Name: "Test 5", Price: 1234.78},
	// 	{Name: "Test 6", Price: 234.78},
	// 	{Name: "Test 7", Price: 110.78},
	// 	{Name: "Test 8", Price: 9756.78},
	// }

	// product := Product{
	// 	Name:  "test",
	// 	Price: 13.89,
	// }

	// err = gorm.G[Product](db).CreateInBatches(ctx, &products, 100)
	// if err != nil {
	// 	panic(err)
	// }

	// var products []Product
	// min := 1.0
	// max := 1000.0
	// for i := range 505 {
	// 	products = append(products, Product{Name: fmt.Sprintf("Product: %d", i), Price: min + rand.Float64()*(max-min)})
	// }
	// gorm.G[Product](db).CreateInBatches(ctx, &products, 100)

	// batchSize := 50
	// batchCount := 0
	// err = gorm.G[Product](db).FindInBatches(ctx, batchSize, func(data []Product, batch int) error {
	// 	batchCount = batch // Keep track of the current batch number

	// 	fmt.Printf("--- Processing Batch %d (Records %d to %d) ---\n", batch, (batch-1)*batchSize+1, ((batch-1)*batchSize)+len(data))

	// 	// Iterate over the records in the current batch
	// 	for _, p := range data {
	// 		fmt.Printf("Product ID: %s, Name: %s\n", p.ID, p.Name)
	// 		// You can perform your processing logic here (e.g., update a field, send an email)
	// 	}

	// 	// Important: If you return an error, the process will stop.
	// 	// return errors.New("stop processing")

	// 	// If everything is okay, return nil to continue to the next batch
	// 	return nil
	// })
	// fmt.Printf("Batch count: %d\n", batchCount)

	// Read
	// product, err := gorm.G[Product](db).Where("id = ?", "019a88e9-fdd3-7cf1-91d2-b4a2d3f7ca8e").First(ctx) // find product with integer primary key
	// if err == gorm.ErrRecordNotFound {
	// 	fmt.Println("Record not found for the given ID.")
	// 	// You can implement specific logic here, like returning a 404 error in an API,
	// 	// or prompting the user to create a new record.
	// } else if err != nil {
	// 	// Handle other potential database errors
	// 	fmt.Printf("An unexpected error occurred: %v\n", err)
	// } else {
	// 	// Record found, proceed with using the 'user' object
	// 	fmt.Printf("Found product: %s (Name: %d)\n", product.ID, product.Name)
	// }

	//   products, err := gorm.G[Product](db).Where("code = ?", "D42").Find(ctx) // find product with code D42

	//   // Update - update product's price to 200
	//   err = gorm.G[Product](db).Where("id = ?", product.ID).Update(ctx, "Price", 200)
	//   // Update - update multiple fields
	//   err = gorm.G[Product](db).Where("id = ?", product.ID).Updates(ctx, Product{Code: "D42", Price: 100})

	//   // Delete - delete product
	// rowCount, err := gorm.G[Product](db).Where("id = ?", product.ID).Delete(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// if rowCount > 0 {
	// 	fmt.Println("Product existed")
	// } else {
	// 	fmt.Println("Not found")
	// }

}

func testProduct(db *gorm.DB, ctx context.Context) {

	// load the category eager
	products, err := gorm.G[Product](db).Preload("Category", nil).Preload("SerialNumber", nil).Find(ctx)
	if err != nil {
		panic(err)
	}
	for _, p := range products {
		fmt.Printf("Category: %s\n", p.Name)
	}
}

func createCategoryAndProduct(db *gorm.DB, ctx context.Context, categoryName, productName string, price float64) *Product {
	category := Category{Name: categoryName}

	gorm.G[Category](db).Create(ctx, &category)

	product := Product{
		Name:       productName,
		Price:      price,
		CategoryID: category.ID,
	}
	gorm.G[Product](db).Create(ctx, &product)
	return &product
}
