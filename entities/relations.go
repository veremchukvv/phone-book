package entities

type RelationList []*Relation

type Relation struct {
	UserID         int
	RelationUserID int
}
