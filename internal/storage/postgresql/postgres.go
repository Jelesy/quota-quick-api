package postgresql

import (
	"api.quota-quick/api/internal/storage"
	"database/sql"
	"errors"
	"fmt"

	"api.quota-quick/api/internal/config"
	"api.quota-quick/api/internal/models"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func (st *Storage) GetDb() *sql.DB {
	return st.db
}

func New(connStr string) (*Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", connStr, err)
	}

	// TODO: Инициализация базы 38:00
	err = initializeSchema(db)
	if err != nil {
		return nil, fmt.Errorf("initializeSchema - %s: %w", connStr, err)
	}

	return &Storage{db: db}, nil
}

func GetConnStr(cfg *config.Config) (string, error) {

	return fmt.Sprintf(
		"host=localhost port=5433 user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Dbname, cfg.Database.Sslmode,
	), nil
}

// Если не будет таблиц, до нужно будет их создать
func initializeSchema(db *sql.DB) error {
	// TODO: Реализовать инициализацию таблиц
	//db.Prepare()
	return nil
}

// TODO: Аргумент строка (json) или структура, где делать конвертацию из json?
func (s *Storage) SaveContainer(cntr models.Container) error {
	const op = "storage.postgres.SaveContainer"

	stmt, err := s.db.Prepare("INSERT INTO containers(title, description, owner_id, created_at, changed_at) VALUES($1, $2, $3, NOW(), NOW())")
	if err != nil {
		return fmt.Errorf("%s: %w", op+"1", err)
	}
	fmt.Printf("%#v\n", cntr)
	_, err = stmt.Exec(cntr.Title, cntr.Description, cntr.OwnerId)
	if err != nil {
		// TODO: Обработка ошибок базы
		//if pqErr, ok := err.(pq.Error); ok && pqErr.Code == pq.ErrorCode() {
		// TODO: Общая ошибка для хендлеров
		//	return fmt.Errorf("%s: %w", op, storage.)
		//}
		return fmt.Errorf("%s: %w", op+"2", err)
	}

	return nil
}

// TODO: Аргумент строка (json) или структура, где делать конвертацию из json?
func (s *Storage) GetContainerById(id uint64) (models.Container, error) {
	const op = "storage.postgres.GetContainer"

	stmt, err := s.db.Prepare("SELECT * FROM containers WHERE id=$1")
	if err != nil {
		return models.Container{}, fmt.Errorf("%s: %w", op, err)
	}

	var cont models.Container
	err = stmt.QueryRow(id).Scan(&cont.ID, &cont.Title, &cont.Description, &cont.OwnerId, &cont.CreatedAt, &cont.ChangedAt)

	if err != nil {
		// TODO: Обработка ошибок базы
		//if pqErr, ok := err.(pq.Error); ok && pqErr.Code == pq.ErrorCode() {
		// TODO: Общая ошибка для хендлеров
		//	return fmt.Errorf("%s: %w", op, storage.Err...)
		//}

		// Проверка ошибки на соответствие
		if errors.Is(err, sql.ErrNoRows) {
			return models.Container{}, storage.ErrContainerNotFound
		}

		return models.Container{}, fmt.Errorf("%s: %w %T", op, err, err)
	}

	//pq.CopyInSchema()

	return cont, nil
}

func (s *Storage) GetAllContainers() ([]models.Container, error) {
	const op = "storage.postgres.GetAllContainers"

	stmt, err := s.db.Prepare("SELECT * FROM containers")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var conts []models.Container
	for rows.Next() {
		var cont models.Container
		err = rows.Scan(&cont.ID, &cont.Title, &cont.Description, &cont.OwnerId, &cont.CreatedAt, &cont.ChangedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		conts = append(conts, cont)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	//pq.CopyInSchema()

	return conts, nil
}

func (s *Storage) DeleteContainerById(id uint64) error {
	const op = "storage.postgres.DeleteContainerById"

	stmt, err := s.db.Prepare("DELETE FROM containers WHERE id=$1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id)

	if err != nil {
		// TODO: Обработка ошибок базы
		//if pqErr, ok := err.(pq.Error); ok && pqErr.Code == pq.ErrorCode() {
		// TODO: Общая ошибка для хендлеров
		//	return fmt.Errorf("%s: %w", op, storage.Err...)
		//}

		return fmt.Errorf("%s: %w", op, err)
	}

	//pq.CopyInSchema()

	return nil
}
