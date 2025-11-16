package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // Import the driver, aliased with _ to run its init function
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UuidModel struct {
	// ID        string `gorm:"primarykey;default:uuidv7()::varchar"`
	ID        uuid.UUID `gorm:"type:uuid;primarykey;default:uuidv7()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Categoria struct {
	UuidModel
	Name string
}

func (Categoria) TableName() string {
	return "catalogo.categorias"
}

type Produto struct {
	UuidModel
	Name         string
	Price        float64
	CategoriaID  uuid.UUID
	Categoria    Categoria
	NumeroSerial SerialNumber
}

type SerialNumber struct {
	gorm.Model
	Number    string
	ProdutoID uuid.UUID
}

func (SerialNumber) TableName() string {
	return "catalogo.numero_seriais"
}

// func (u *Produto) BeforeCreate(tx *gorm.DB) error {
// 	// Verifica se o ID já está preenchido (útil para testes ou inserts manuais)
// 	if u.ID == "" {
// 		uuid, err := uuid.NewV7()
// 		if err != nil {
// 			return err
// 		}
// 		// Gera um UUID v7 e o converte para string
// 		u.ID = uuid.String()
// 		// Se você estiver usando uma biblioteca específica para V7:
// 		// u.ID = uuid.NewV7().String()
// 	}
// 	return nil
// }

func (Produto) TableName() string {
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
		//Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}

	ctx := context.Background()

	// Migrate the schema
	err = db.AutoMigrate(&Categoria{}, &Produto{}, &SerialNumber{})
	if err != nil {
		panic(err)
	}

	// prod := criaCategoriaEProduto(db, ctx, "Categoria 1", "Produto 1", 34.87)

	// gorm.G[SerialNumber](db).Create(ctx, &SerialNumber{
	// 	Number:    "12345",
	// 	ProdutoID: prod.ID,
	// })

	testaProduto(db, ctx)

	// produto := Produto{Name: "Teste 5", Price: 1234.78}

	// err = gorm.G[Produto](db).Create(ctx, &produto)
	// if err != nil {
	// 	panic(err)
	// }

	// produtos := []Produto{
	// 	{Name: "Teste 5", Price: 1234.78},
	// 	{Name: "Teste 6", Price: 234.78},
	// 	{Name: "Teste 7", Price: 110.78},
	// 	{Name: "Teste 8", Price: 9756.78},
	// }

	// produto := Produto{
	// 	Name:  "teste",
	// 	Price: 13.89,
	// }

	// err = gorm.G[Produto](db).CreateInBatches(ctx, &produtos, 100)
	// if err != nil {
	// 	panic(err)
	// }

	// var produtos []Produto
	// min := 1.0
	// max := 1000.0
	// for i := range 505 {
	// 	produtos = append(produtos, Produto{Name: fmt.Sprintf("Produto: %d", i), Price: min + rand.Float64()*(max-min)})
	// }
	// gorm.G[Produto](db).CreateInBatches(ctx, &produtos, 100)

	// batchSize := 50
	// batchCount := 0
	// err = gorm.G[Produto](db).FindInBatches(ctx, batchSize, func(data []Produto, batch int) error {
	// 	batchCount = batch // Keep track of the current batch number

	// 	fmt.Printf("--- Processing Batch %d (Records %d to %d) ---\n", batch, (batch-1)*batchSize+1, ((batch-1)*batchSize)+len(data))

	// 	// Iterate over the records in the current batch
	// 	for _, p := range data {
	// 		fmt.Printf("Produto ID: %s, Name: %s\n", p.ID, p.Name)
	// 		// You can perform your processing logic here (e.g., update a field, send an email)
	// 	}

	// 	// Important: If you return an error, the process will stop.
	// 	// return errors.New("stop processing")

	// 	// If everything is okay, return nil to continue to the next batch
	// 	return nil
	// })
	// fmt.Printf("Batch count: %d\n", batchCount)

	// Read
	// product, err := gorm.G[Produto](db).Where("id = ?", "019a88e9-fdd3-7cf1-91d2-b4a2d3f7ca8e").First(ctx) // find product with integer primary key
	// if err == gorm.ErrRecordNotFound {
	// 	fmt.Println("Record not found for the given ID.")
	// 	// You can implement specific logic here, like returning a 404 error in an API,
	// 	// or prompting the user to create a new record.
	// } else if err != nil {
	// 	// Handle other potential database errors
	// 	fmt.Printf("An unexpected error occurred: %v\n", err)
	// } else {
	// 	// Record found, proceed with using the 'user' object
	// 	fmt.Printf("Found produto: %s (Nome: %d)\n", product.ID, product.Name)
	// }

	//   products, err := gorm.G[Product](db).Where("code = ?", "D42").Find(ctx) // find product with code D42

	//   // Update - update product's price to 200
	//   err = gorm.G[Product](db).Where("id = ?", product.ID).Update(ctx, "Price", 200)
	//   // Update - update multiple fields
	//   err = gorm.G[Product](db).Where("id = ?", product.ID).Updates(ctx, Product{Code: "D42", Price: 100})

	//   // Delete - delete product
	// rowCount, err := gorm.G[Produto](db).Where("id = ?", product.ID).Delete(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// if rowCount > 0 {
	// 	fmt.Println("Produto existia")
	// } else {
	// 	fmt.Println("Não achou")
	// }

}

func testaProduto(db *gorm.DB, ctx context.Context) {

	// carrega a categoria eager
	produtos, err := gorm.G[Produto](db).Preload("Categoria", nil).Preload("NumeroSerial", nil).Find(ctx)
	if err != nil {
		panic(err)
	}
	for _, p := range produtos {
		fmt.Printf("Produto: %v\n", p)
	}
}

func criaCategoriaEProduto(db *gorm.DB, ctx context.Context, nomeCategoria, nomeProduto string, preco float64) *Produto {
	categoria := Categoria{Name: nomeCategoria}

	gorm.G[Categoria](db).Create(ctx, &categoria)

	produto := Produto{
		Name:        nomeProduto,
		Price:       preco,
		CategoriaID: categoria.ID,
	}
	gorm.G[Produto](db).Create(ctx, &produto)
	return &produto
}
