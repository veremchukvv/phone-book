package service

import (
	"context"
	"integ/entities"
)

type ContactStore interface {
	FindUserByPhone(ctx context.Context, number string) (*entities.User, error)
	SaveRelation(ctx context.Context, relation *entities.Relation) error
	FindFriends(ctx context.Context, uid int) ([]*entities.Friend, error)
	GetName(ctx context.Context, uid int) (string, error)
}

type ContactService struct {
	store ContactStore
}

func NewContactService(store ContactStore) *ContactService {
	return &ContactService{store: store}
}

func (c *ContactService) SaveContacts(ctx context.Context, userID int, contacts entities.ContactList) error {
	relations := make(entities.RelationList, 0)

	// tx, err := c.store.StartTX(ctx, pgx.TxOptions{})
	// if err != nil {
	// 	return nil
	// }

	// ctx = context.WithValue(ctx, "tx", tx)

	for _, contact := range contacts {
		user, err := c.store.FindUserByPhone(ctx, contact.PhoneNumber)
		if err != nil {
			return err
		}
		if user == nil {
			continue
		}

		relations = append(relations, &entities.Relation{
			UserID:         userID,
			RelationUserID: user.UserID,
		})
	}

	for _, relation := range relations {
		err := c.store.SaveRelation(ctx, relation)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ContactService) FindFriends(ctx context.Context, userID int) (entities.FriendsList, error) {
	friends := make(entities.FriendsList, 0)

	resp, err := c.store.FindFriends(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, f := range resp {
		friends = append(friends, *f)
	}

	return friends, nil
}

func (c *ContactService) GetName(ctx context.Context, userID int) (string, error) {

	resp, err := c.store.GetName(ctx, userID)
	if err != nil {
		return "", err
	}

	return resp, nil
}

// func (c *ContactService) ExecWithTX(ctx context.Context, txOptions pgx.TxOptions) (Result, error) {

// 	tx, err := c.db.ExecWithTransaction(ctx, pgx.TxOptions{})

// 	return tx, nil
// }
