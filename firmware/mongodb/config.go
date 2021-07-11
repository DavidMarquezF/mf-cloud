package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

type FirmwareInfo struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Date    primitive.DateTime `bson:"date"`
	Version string             `bson:"version"`
	Elf     primitive.ObjectID `bson:"elf"`
}

type FirmwareExec struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Exec primitive.Binary   `bson:"exec"`
}
