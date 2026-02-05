package product

type ProductInputDto struct {
	ID    string
	Name  string
	Price float64
}

type InsertProductUseCase interface {
	Execute(product *ProductInputDto) error
}

type ProductUseCaseImpl struct {
	ProductRepository ProductRepository
}

func NewProductUseCaseImpl(productRepository ProductRepository) InsertProductUseCase {
	return &ProductUseCaseImpl{ProductRepository: productRepository}
}

func (u *ProductUseCaseImpl) Execute(product *ProductInputDto) error {
	err := u.ProductRepository.Insert(&Product{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	})
	if err != nil {
		return err
	}
	return nil
}
