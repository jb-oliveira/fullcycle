package product

type ProductInputDto struct {
	ID    string
	Name  string
	Price float64
}

type InsertProductUseCase interface {
	Execute(product *ProductInputDto) error
}

type productUseCaseImpl struct {
	productRepository ProductRepository
}

func NewProductUseCaseImpl(productRepository ProductRepository) InsertProductUseCase {
	return &productUseCaseImpl{productRepository: productRepository}
}

func (u *productUseCaseImpl) Execute(product *ProductInputDto) error {
	err := u.productRepository.Insert(&Product{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	})
	if err != nil {
		return err
	}
	return nil
}
