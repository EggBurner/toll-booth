package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"toll-booth/internal/db"
	"toll-booth/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type GetLinksRequest struct {
	UserID string `json:"userID"`
}

type GetLinkReuquestLinkComponent struct {
	LinkID          bson.ObjectID    `json:"_id"             bson:"_id"`
	TargetLink      string           `json:"targetLink"      bson:"targetLink"`
	ShortCode       string           `json:"shortCode"       bson:"shortCode"`
	LinkDateCreated time.Time        `json:"linkDateCreated" bson:"linkDateCreated"`
	OwnerID         bson.ObjectID    `json:"ownerID"         bson:"ownerID"`
	Status          model.StatusEnum `json:"status"          bson:"status"`
	VisitCount      uint32           `json:"visitCount"      bson:"visitCount"`
	PinProtected    bool             `json:"pinProtected"    bson:"pinProtected"`
}
type GetLinksResponse struct {
	Links []GetLinkReuquestLinkComponent `json:"links"`
}

func GetAllLinks(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req GetLinksRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	var res GetLinksResponse
	res.Links = []GetLinkReuquestLinkComponent{}

	collection := db.Client.Database("toll-booth-db").Collection("links")

	ownerID, err := bson.ObjectIDFromHex(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	filter := bson.M{
		"ownerID": ownerID,
	}

	cursor, err := collection.Find(context.TODO(), filter)

	if err != nil {
		http.Error(w, "Error accesing db", http.StatusInternalServerError)
		return
	}

	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &res.Links); err != nil {
		http.Error(w, "error parsing", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(&res)

}
