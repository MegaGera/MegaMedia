package handlers

import (
	"MegaMedia/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Image struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name,omitempty" json:"name,omitempty"`
	URL      string             `bson:"url,omitempty" json:"url,omitempty"`
	HeightPX int                `bson:"heightPX,omitempty" json:"heightPX,omitempty"`
	WidthPX  int                `bson:"widthPX,omitempty" json:"widthPX,omitempty"`
	Previous []string           `bson:"previous,omitempty" json:"previous,omitempty"` // Optional
}

// FetchImages from a collection
func FetchImages(w http.ResponseWriter, r *http.Request, databaseName string, collectionName string) error {

	// Create MongoDB client
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Cfg.MongoDBURI))
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return fmt.Errorf("database connection error: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Access database and collection
	collection := client.Database(databaseName).Collection(collectionName)

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Fetch all images
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		// http.Error(w, "Error fetching images", http.StatusInternalServerError)
		return fmt.Errorf("error fetching images: %v", http.StatusInternalServerError)
	}
	defer cursor.Close(context.TODO())

	// Decode cursor results into a slice
	var images []Image
	if err := cursor.All(context.TODO(), &images); err != nil {
		// http.Error(w, "Error decoding images", http.StatusInternalServerError)
		return fmt.Errorf("error decoding images: %v", http.StatusInternalServerError)
	}

	// Convert to JSON
	imagesJSON, err := json.Marshal(images)
	if err != nil {
		return fmt.Errorf("error encoding images to JSON: %v", http.StatusInternalServerError)
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	w.Write(imagesJSON)
	return nil
}

// FetchImageByID fetches a single image by its ID
func FetchImageByID(w http.ResponseWriter, r *http.Request, databaseName string, collectionName string, id string) (imageReturn Image, err error) {

	// Create MongoDB client
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Cfg.MongoDBURI))
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return Image{}, fmt.Errorf("database connection error: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Access database and collection
	collection := client.Database(databaseName).Collection(collectionName)

	// Convert id string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return Image{}, fmt.Errorf("invalid ID format: %v", err)
	}

	// Find the image by ID
	var image Image
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&image)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Image not found", http.StatusNotFound)
			return image, fmt.Errorf("image not found: %v", err)
		}
		http.Error(w, "Error fetching image", http.StatusInternalServerError)
		return Image{}, fmt.Errorf("error fetching image: %v", err)
	}

	return image, nil
}
