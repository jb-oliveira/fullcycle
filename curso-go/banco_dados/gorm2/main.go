package main

type Category struct {
	ID       string `gorm:"primaryKey"`
	Name     string
	Produtos []Product
}

type Product struct {
	ID           string `gorm:"primaryKey"`
	Name         string
	Price        float64
	CategoryID   int
	Category     Category
	SerialNumber SerialNumber
}

type SerialNumber struct {
	ID        string `gorm:"primaryKey"`
	Number    string
	ProductID int
}

func main() {

}
