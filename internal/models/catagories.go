package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Catagory struct {
	ID   int    `json:"-"`
	Name string `json:"name" validate:"required"`
}

type CatagoryModel struct {
	DB *pgxpool.Pool
}

func (cm *CatagoryModel) Insert(c *Catagory) error {
	statement := `
  INSERT INTO categories (name)
  VALUES ($1)
  `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := cm.DB.Exec(ctx, statement, &c.Name)
	if err != nil {
		return err
	}

	return nil
}

func (cm *CatagoryModel) DeleteByName(name string) error {
	statement := `
  DELETE FROM categories
  WHERE name = ($1)
  `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := cm.DB.Exec(ctx, statement, name)
	if err != nil {
		return err
	}

	return nil
}

func (cm *CatagoryModel) EditOnUser(userName string, categories []string, activate bool) error {
	activateTemplate := `
  INSERT INTO user_categories (user_id, category_id)
  VALUES
  ((SELECT id FROM users WHERE name = $1), (SELECT id FROM categories WHERE name = $2))
  `

	deactivateTemplate := `
  DELETE FROM user_categories
  WHERE
  user_id = (SELECT id FROM users WHERE name = $1)
  AND
  category_id = (SELECT id FROM categories WHERE name = $2)
  `

	batch := &pgx.Batch{}

	for _, category := range categories {
		switch activate {
		case true:
			batch.Queue(activateTemplate, userName, category)
		case false:
			batch.Queue(deactivateTemplate, userName, category)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	br := cm.DB.SendBatch(ctx, batch)
	_, err := br.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (um *CatagoryModel) Exists(id int) error {
	statement := `
  SELECT null FROM categories
  WHERE id = $1
  `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRow(ctx, statement).Scan()
	if err != nil {
		return err
	}

	return nil
}

func (um *CatagoryModel) GetAll(filters Filters) ([]*Catagory, Metadata, error) {
	statement := fmt.Sprintf(`
  SELECT count(*) OVER(), name FROM categories
  ORDER BY name %s, id ASC
  LIMIT %d OFFSET %d `, filters.sortDirection(), filters.limit(), filters.offset())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := um.DB.Query(ctx, statement)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	categories := []*Catagory{}

	for rows.Next() {
		var cat Catagory

		err := rows.Scan(
			&totalRecords,
			&cat.Name,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		categories = append(categories, &cat)
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return categories, metadata, nil
}

func (um *CatagoryModel) GetAllActive(userID int, filters Filters) ([]*Catagory, Metadata, error) {
	// we use Sprintf because we can't use variables in the some of the paramaters
	// statement := fmt.Sprintf(`
	// SELECT count(*) OVER(), categories.name FROM categories
	// JOIN user_categories
	// ON categories.id = user_categories.category_id
	// WHERE user_categories.user_id = %d
	// ORDER BY name %s, id ASC
	// LIMIT %d OFFSET %d `, userID, filters.sortDirection(), filters.limit(), filters.offset())

	statement := `
  SELECT count(*) OVER(), categories.name FROM categories
  JOIN user_categories
  ON categories.id = user_categories.category_id
  JOIN users 
  ON user_categories.user_id = users.id 
  WHERE user_categories.user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Println("ALKDSJF", userID)
	rows, err := um.DB.Query(ctx, statement, userID)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	categories := []*Catagory{}

	for rows.Next() {
		var cat Catagory

		err := rows.Scan(
			&totalRecords,
			&cat.Name,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		categories = append(categories, &cat)
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return categories, metadata, nil
}
