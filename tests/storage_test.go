package tests

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
)

func TestStore_FindUserByPhone(t *testing.T) {
	testCtx := prepareTestContext(t)

	_, err := testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь1', '+7983');`)
	_, err = testCtx.store.Exec(context.Background(), `insert into users (name, phone_number) values ('Пользователь2', '+7984');`)

	user, err := testCtx.store.FindUserByPhone(context.Background(), "+7983")
	if err != nil {
		t.Fatal(err)
	}

	if user.Name != "Пользователь1" {
		t.Fatal("должен был вернуться первый пользователь")
	}

}
