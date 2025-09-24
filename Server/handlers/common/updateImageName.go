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

type UpdateNameRequest struct {
	Name string `json:"name"`
}

// UpdateImageName updates the name of an image in MongoDB
func UpdateImageName(w http.ResponseWriter, r *http.Request, databaseName string, collectionName string, imageID string) error {
	// Parse the request body
	var updateRequest UpdateNameRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		return fmt.Errorf("failed to decode request body: %v", err)
	}

	if updateRequest.Name == "" {
		return fmt.Errorf("name is required")
	}

	// Create MongoDB client
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Cfg.MongoDBURI))
	if err != nil {
		return fmt.Errorf("database connection error: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Access database and collection
	collection := client.Database(databaseName).Collection(collectionName)

	// Convert id string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(imageID)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	// Update the image name
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"name": updateRequest.Name}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update image name: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("image not found")
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"success": true,
		"message": "Image name updated successfully",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}

	w.Write(jsonResponse)
	return nil
}
