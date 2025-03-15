package shop

type Shop struct {
	Id          string `bson:"_id" json:"_id"`
	Version     int    `bson:"version" json:"version"`
	Name        string `bson:"name" json:"name"`
	Location    string `bson:"location" json:"location"`
	Description string `bson:"description" json:"description"`
}
