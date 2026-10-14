package repository

// UserRepository хранит пользователей в памяти
type UserRepository struct {
	// Карта, где ключ - ID (int), значение - имя (string)
	data map[int]string
}

// NewUserRepository - конструктор репозитория (в Go это обычная функция)
func NewUserRepository() *UserRepository {
	return &UserRepository{
		data: map[int]string{
			1: "Алексей",
			2: "Иван",
			3: "Мария",
		},
	}
}

// GetByID ищет пользователя по ID.
// Возвращает имя и булево значение (найден или нет).
func (r *UserRepository) GetByID(id int) (string, bool) {
	name, exists := r.data[id]
	return name, exists
}
