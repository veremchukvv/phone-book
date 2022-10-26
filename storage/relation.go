package storage

import (
	"context"
	"fmt"
	"integ/entities"
	"log"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
)

var (
	dialect = goqu.Dialect("postgres")

	userTable = dialect.From("users").Prepared(true)
	userCols  = []interface{}{
		"user_id",
		"name",
		"phone_number",
	}

	// relationTable = dialect.From("relations").Prepared(true)
	// relationCols  = []interface{}{
	// 	"relation_id",
	// 	"user_id",
	// 	"relation_user_id",
	// }

	friendsReq  = dialect.From("relations").InnerJoin(goqu.T("users"), goqu.On(goqu.Ex{"relations.relation_user_id": goqu.I("users.user_id")})).Prepared(true)
	friendsCols = []interface{}{
		"relations.user_id",
		"relations.relation_user_id",
		"users.phone_number",
	}
)

type UserRaw struct {
	UserID      int
	Name        string
	PhoneNumber string
}

type RelationRaw struct {
	RelationID     int
	UserID         int
	RelationUserID int
}

type FriendsRaw struct {
	UserID         int
	RelationUserID int
	PhoneNumber    string
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

	if ctx.Value("tx") != nil {
		t := ctx.Value("tx").(*pgx.Tx)
		tx := *t

		rows, err := tx.Query(ctx, selectSQL, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to execute a query user")
		}

		defer rows.Close()

		log.Print("in transaction")

		userRaws := make([]*UserRaw, 0)

		for rows.Next() {
			userRaw, err := scanUser(rows)
			if err != nil {
				return nil, fmt.Errorf("failed to read user from database")
			}
			userRaws = append(userRaws, userRaw)
		}

		if len(userRaws) == 0 {
			return nil, err
		}

		user := buildUser(userRaws[0])

		return user, nil

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
	createRelationSQL, args, err := goqu.
		Insert("relations").
		Cols("user_id", "relation_user_id").
		Vals(goqu.Vals{relation.UserID, relation.RelationUserID}).
		ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build a query insert relation")
	}

	if ctx.Value("tx") != nil {
		t := ctx.Value("tx").(*pgx.Tx)
		tx := *t

		_, err = tx.Exec(ctx, createRelationSQL, args...)
		if err != nil {
			return fmt.Errorf("failed to build a query insert relation")
		}
	}

	// реализация без транзакций
	_, err = s.Exec(ctx, createRelationSQL, args...)
	if err != nil {
		return fmt.Errorf("failed to build a query insert relation")
	}

	//реализация с транзакциями
	// _, err = s.db.ExecWithTransaction(ctx, createRelationSQL, pgx.TxOptions{}, args...)
	// if err != nil {
	// 	return fmt.Errorf("failed to build a query insert relation")
	// }

	return nil
}

func (s *Store) FindFriends(ctx context.Context, uid int) ([]*entities.Friend, error) {
	selectSQL, args, err := friendsReq.
		Select(friendsCols...).
		Where(goqu.T("relations").Col("user_id").Eq(uid)).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("failed to build a friends request")
	}
	rows, err := s.Query(ctx, selectSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute a friends request")
	}
	friendsRaw := make([]*FriendsRaw, 0)
	for rows.Next() {
		friendRaw, err := scanFriends(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to read friends from database")
		}
		friendsRaw = append(friendsRaw, friendRaw)
	}

	friends := make([]*entities.Friend, 0)

	for _, frnd := range friendsRaw {
		friend := buildFriends(frnd)
		friends = append(friends, friend)
	}

	return friends, nil
}

func (s *Store) GetName(ctx context.Context, uid int) (string, error) {
	selectSQL, args, err := userTable.
		Select(userCols...).
		Where(goqu.C("user_id").Eq(uid)).
		Limit(uint(1)).
		ToSQL()

	if err != nil {
		return "", fmt.Errorf("failed to build a name request")
	}

	rows, err := s.Query(ctx, selectSQL, args...)
	if err != nil {
		return "", fmt.Errorf("failed to execute a name request")
	}

	userRaws := make([]*UserRaw, 0)
	for rows.Next() {
		userRaw, err := scanUser(rows)
		if err != nil {
			return "", fmt.Errorf("failed to read user from database")
		}
		userRaws = append(userRaws, userRaw)
	}

	user := buildUser(userRaws[0])
	return user.Name, nil
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

// func scanRelation(rows Row) (*RelationRaw, error) {
// 	relationRaw := RelationRaw{}
// 	err := rows.Scan(
// 		&relationRaw.RelationID,
// 		&relationRaw.UserID,
// 		&relationRaw.RelationUserID,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &relationRaw, nil
// }

func scanFriends(rows Row) (*FriendsRaw, error) {
	friendsRaw := FriendsRaw{}
	err := rows.Scan(
		&friendsRaw.UserID,
		&friendsRaw.RelationUserID,
		&friendsRaw.PhoneNumber,
	)

	if err != nil {
		return nil, err
	}

	return &friendsRaw, nil
}

func buildUser(userRaw *UserRaw) *entities.User {
	return &entities.User{
		UserID:      userRaw.UserID,
		Name:        userRaw.Name,
		PhoneNumber: userRaw.PhoneNumber,
	}
}

// func buildRelations(relationRaw *RelationRaw) *entities.Relation {
// 	return &entities.Relation{
// 		UserID:         relationRaw.UserID,
// 		RelationUserID: relationRaw.RelationUserID,
// 	}
// }

func buildFriends(friendRaw *FriendsRaw) *entities.Friend {
	return &entities.Friend{
		UserID:         friendRaw.UserID,
		RelationUserID: friendRaw.RelationUserID,
		PhoneNumber:    friendRaw.PhoneNumber,
	}
}
