package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"toll-booth/internal/auth"
	"toll-booth/internal/db"
	"toll-booth/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type redirectResponseInfo struct {
	TargetLink   string           `json:"targetLink" bson:"targetLink"`
	ShortCode    string           `json:"shortCode" bson:"shortCode"`
	Status       model.StatusEnum `json:"status" bson:"status"`
	PinProtected bool             `json:"pinProtected" bson:"pinProtected"`
}

type redirectRequestInfo struct {
	ShortCode string  `json:"shortCode"  bson:"shortCode"`
	Pin       *string `json:"PIN,omitempty"`
}

func HandleRedirectRequest(w http.ResponseWriter, r *http.Request) {
	var retrievedLink redirectResponseInfo
	shortCode := r.PathValue("shortCode")

	collection := db.Client.Database("toll-booth-db").Collection("links")

	filter := bson.M{
		"shortCode": shortCode,
	}

	result := collection.FindOne(context.TODO(), filter)

	err := result.Decode(&retrievedLink)

	if err == mongo.ErrNoDocuments {
		http.Error(w, "Link doesn't exist", http.StatusBadRequest)
	}
	if err != nil {
		http.Error(w, "Error retreiving data", http.StatusInternalServerError)
		return
	}

	if retrievedLink.Status != "ACTIVE" {
		http.Error(w, "Link Not Active", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&retrievedLink)

}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {

	var redirectRequestInfo redirectRequestInfo

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&redirectRequestInfo)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var retrievedLink model.Link

	collection := db.Client.Database("toll-booth-db").Collection("links")

	filter := bson.M{
		"shortCode": redirectRequestInfo.ShortCode,
	}

	result := collection.FindOne(context.TODO(), filter)

	err = result.Decode(&retrievedLink)
	if err == mongo.ErrNoDocuments {
		http.Error(w, "Error, no link found", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Error retreiving data", http.StatusInternalServerError)
		return
	}

	if retrievedLink.PinProtected {
		err := auth.ComparePassword(*retrievedLink.PinHash, *redirectRequestInfo.Pin)
		if err != nil {
			http.Error(w, "Invalid PIN", http.StatusBadRequest)
			return
		}
	}
	filter = bson.M{
		"shortCode": redirectRequestInfo.ShortCode,
	}

	update := bson.M{
		"$inc": bson.M{
			"visitCount": 1,
		},
	}

	collection = db.Client.Database("toll-booth-db").Collection("links")

	result = collection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
	)

	if result.Err() != nil {
		http.Error(w, "Error updating db"+result.Err().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"targetLink": retrievedLink.TargetLink,
		"pinMatch":   true,
	})
}
