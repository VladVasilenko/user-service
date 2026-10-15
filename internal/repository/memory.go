package repository

type UserRepository struct {
	data map[int]string
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		data: map[int]string{
			1: "Алексей",
			2: "Иван",
			3: "Мария",
		},
	}
}

func (r *UserRepository) GetByID(id int) (string, bool) {
	name, exists := r.data[id]
	return name, exists
}
