package service

import (
	"context"
	"toll-booth/internal/db"

	"github.com/jxskiss/base62"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func LinkShortner() (shortenedLink string, err error) {

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	filter := bson.M{
		"_id": "shortcode",
	}

	update := bson.M{
		"$inc": bson.M{
			"seq": 1,
		},
	}

	collection := db.Client.Database("toll-booth-db").Collection("counter")

	result := collection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		opts,
	)

	var counter struct {
		Seq int64 `bson:"seq"`
	}

	err = result.Decode(&counter)

	num := counter.Seq
	encodedBytes := base62.FormatInt(num)

	shortCode := string(encodedBytes)

	return shortCode, err
}
