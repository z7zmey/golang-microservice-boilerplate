package repository

type CarRepository interface {
	Find(id int) ([]Car, error)
}

type carRepo struct{}

func NewCarRepo() *carRepo {
	return &carRepo{}
}

func (carRepo) Find(_ int) ([]Car, error) {
	return []Car{
		{1, "hatchback", "DescriptionPeugeot"},
		{2, "sedan", "volkswagen"},
	}, nil
}