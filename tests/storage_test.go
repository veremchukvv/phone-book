package tests

import (
	"context"
	"integ/entities"
	"testing"

	_ "github.com/lib/pq"
)

func TestStore_FindUserByPhone(t *testing.T) {
	testCtx := prepareTestContext(t)

	_, err := testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь1', '+7983');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь2', '+7984');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь3', '+7985');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь4', '+7986');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь5', '+7987');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь2', '+7988');`)

	user, err := testCtx.store.FindUserByPhone(context.Background(), "+7983")
	if err != nil {
		t.Fatal(err)
	}

	if user.Name != "Пользователь1" {
		t.Fatal("должен был вернуться первый пользователь")
	}

}

func TestStore_FindFriends(t *testing.T) {
	testCtx := prepareTestContext(t)

	_, err := testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь1', '+7983');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь2', '+7984');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь3', '+7985');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь4', '+7986');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь5', '+7987');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь2', '+7988');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь2', '+7988');`)

	// _, err = testCtx.store.Exec(context.Background(), `insert into relations (user_id, relation_user_id) values (1, 2);`)
	// _, err = testCtx.store.Exec(context.Background(), `insert into relations (user_id, relation_user_id) values (1, 3);`)
	// _, err = testCtx.store.Exec(context.Background(), `insert into relations (user_id, relation_user_id) values (1, 4);`)

	err = testCtx.store.SaveRelation(context.Background(), &entities.Relation{UserID: 1, RelationUserID: 2})
	err = testCtx.store.SaveRelation(context.Background(), &entities.Relation{UserID: 1, RelationUserID: 3})
	err = testCtx.store.SaveRelation(context.Background(), &entities.Relation{UserID: 1, RelationUserID: 4})

	for {

	}

	rel, err := testCtx.store.FindFriends(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	if (rel[0].PhoneNumber != "+7984") && (rel[1].PhoneNumber != "+7985") && (rel[2].PhoneNumber != "+7986") {
		t.Fatal("вернулся неправильный список друзей")
	}

}

func TestStore_GetName(t *testing.T) {
	testCtx := prepareTestContext(t)

	_, err := testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь1', '+7983');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь2', '+7984');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь3', '+7985');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь4', '+7986');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь5', '+7987');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь2', '+7988');`)

	name, err := testCtx.store.GetName(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	if name != "Пользователь1" {
		t.Fatal("должен был вернуться первый пользователь")
	}

}

// func TestStore_SaveRelations(t *testing.T) {
// 	testCtx := prepareTestContext(t)

// 	_, err := testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь1', '+7983');`)
// 	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь2', '+7984');`)
// 	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь3', '+7985');`)
// 	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь4', '+7986');`)
// 	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь5', '+7987');`)
// 	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь2', '+7988');`)

// 	err = testCtx.store.SaveRelation(context.Background(), &entities.Relation{1, 2})

// }
