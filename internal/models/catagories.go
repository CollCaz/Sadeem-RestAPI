package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Catagory struct {
	ID      int    `json:"-"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled,omitempty"`
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

func (cm *CatagoryModel) Delete(c *Catagory) error {
	statement := `
  DELETE FROM categories
  WHERE name = ($1)
  `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := cm.DB.Exec(ctx, statement, &c.Name)
	if err != nil {
		return err
	}

	return nil
}

func (cm *CatagoryModel) GetByName(name string) (*Catagory, error) {
	cat := &Catagory{}
	statement := `
  SELECT id, name, activated
  FROM categories
  WHERE LOWER(name) = LOWER($1)
  `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := cm.DB.QueryRow(ctx, statement, name).Scan(&cat.ID, &cat.Name, &cat.Enabled)
	if err != nil {
		return nil, err
	}

	return cat, nil
}

func (cm *CatagoryModel) Activate(c *Catagory) error {
	statement := `
  UPDATE catagories
  SET activated = TRUE
  WHERE name = ($1)
  `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := cm.DB.Exec(ctx, statement, &c.Name)
	if err != nil {
		return err
	}

	return nil
}

func (um *CatagoryModel) GetAll(filters Filters) ([]*Catagory, Metadata, error) {
	statement := fmt.Sprintf(`
  SELECT count(*) OVER(), name, activated FROM categories
  ORDER BY %s %s, id ASC
  LIMIT %d OFFSET %d `, filters.sortColumn(), filters.sortDirection(), filters.limit(), filters.offset())

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
			&cat.Enabled,
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
	statement := fmt.Sprintf(`
  SELECT count(*) OVER(), name FROM categories
  JOIN user_categories
  ON categories.id = user_categories.category_id
  WHERE user_categories.user_id = %d
  ORDER BY name %s, id ASC
  LIMIT %d OFFSET %d `, userID, filters.sortDirection(), filters.limit(), filters.offset())

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
