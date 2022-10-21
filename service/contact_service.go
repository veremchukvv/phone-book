package service

import (
	"context"
	"integ/entities"
)

type ContactStore interface {
	FindUserByPhone(ctx context.Context, number string) (*entities.User, error)
	SaveRelation(ctx context.Context, relation *entities.Relation) error
}

type ContactService struct {
	store ContactStore
}

func NewContactService(store ContactStore) *ContactService {
	return &ContactService{store: store}
}

func (c *ContactService) SaveContacts(ctx context.Context, userID int, contacts entities.ContactList) error {
	relations := make(entities.RelationList, 0)
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
