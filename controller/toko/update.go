package toko

import (
	"dungeonSnackBE/model"
	"encoding/json"
	"net/http"
	"time"
	"dungeonSnackBE/helper/slug"
	"dungeonSnackBE/config"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateToko(w http.ResponseWriter, r *http.Request) {
	var toko model.Toko
	err := json.NewDecoder(r.Body).Decode(&toko)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	toko.Slug = slug.GenerateSlug(toko.NamaToko)

	for i := range toko.Produk {
		toko.Produk[i].ID = primitive.NewObjectID() // Generate new ObjectID
		toko.Produk[i].CreatedAt = time.Now()       // Set createdAt to current time
		toko.Produk[i].UpdatedAt = time.Now()       // Set updatedAt to current time
	}

	collection := config.Mongoconn.Collection("toko")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.UpdateOne(ctx, primitive.D{primitive.E{Key: "slug", Value: toko.Slug}}, primitive.D{primitive.E{Key: "$set", Value: toko}})
	if err != nil {
		http.Error(w, "Failed to update toko", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(toko)
}