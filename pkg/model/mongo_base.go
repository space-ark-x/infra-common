package model

type MongoBaseModel struct {
	Id        string `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt int64  `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt int64  `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
