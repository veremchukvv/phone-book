package storage

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"integ/entities"
)

var (
	userTable = goqu.Dialect("postgres").From("users").Prepared(true)
	userCols  = []interface{}{
		"user_id",
		"name",
		"phone_number",
	}

	relationTable = goqu.Dialect("postgres").From("relations").Prepared(true)
	relationCols  = []interface{}{
		"relation_id",
		"user_id",
		"relation_user_id",
	}
)

type UserRaw struct {
	UserID      int
	Name        string
	PhoneNumber string
}

func (s *Store) FindUserByPhone(ctx context.Context, number string) (*entities.User, error) {
	selectSQL, args, err := userTable.
		Select(userCols...).
		Where(goqu.C("phone_number").Eq(number)).
		Limit(uint(1)).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build a query user")
	}
	rows, err := s.Query(ctx, selectSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute a query user")
	}
	userRaws := make([]*UserRaw, 0)
	for rows.Next() {
		userRaw, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to read user from database")
		}
		userRaws = append(userRaws, userRaw)
	}

	user := buildUser(userRaws[0])

	return user, nil
}

func (s *Store) SaveRelation(ctx context.Context, relation *entities.Relation) error {
	createRelationSQL, args, err := relationTable.
		Insert().
		Rows(goqu.Record{
			"user_id":          relation.UserID,
			"relation_user_id": relation.RelationUserID,
		}).
		ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build a query insert relation")
	}

	_, err = s.Exec(ctx, createRelationSQL, args)
	if err != nil {
		return fmt.Errorf("failed to build a query insert relation")
	}

	return nil
}

func scanUser(rows Row) (*UserRaw, error) {
	userRaw := UserRaw{}
	err := rows.Scan(
		&userRaw.UserID,
		&userRaw.Name,
		&userRaw.PhoneNumber,
	)

	if err != nil {
		return nil, err
	}

	return &userRaw, nil
}

func buildUser(userRaw *UserRaw) *entities.User {
	return &entities.User{
		UserID:      userRaw.UserID,
		Name:        userRaw.Name,
		PhoneNumber: userRaw.PhoneNumber,
	}
}
