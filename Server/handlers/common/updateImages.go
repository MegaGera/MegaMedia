package handlers

import (
	"MegaMedia/internal/config"
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateImageName finds an image by name and updates its "previous" field with a new name
func UpdateImagesPrevious(w http.ResponseWriter, r *http.Request, databaseName string, collectionName string, imageID primitive.ObjectID, newName string, operation string) error {

	// Create MongoDB client
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Cfg.MongoDBURI))
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return fmt.Errorf("database connection error: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Access database and collection
	collection := client.Database(databaseName).Collection(collectionName)

	filter := bson.M{"_id": imageID}
	var image Image
	err = collection.FindOne(context.TODO(), filter).Decode(&image)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Image not found", http.StatusNotFound)
			return fmt.Errorf("image not found: %v", err)
		}
		http.Error(w, "Error finding image", http.StatusInternalServerError)
		return fmt.Errorf("error finding image: %v", err)
	}
	var update bson.M

	if operation == "add" {
		// Update the "previous" field by adding newName
		update = bson.M{
			"$push": bson.M{"previous": newName},
		}
	} else if operation == "delete" {
		// Update the "previous" field by removing newName
		update = bson.M{
			"$pull": bson.M{"previous": newName},
		}
	} else {
		http.Error(w, "Invalid operation", http.StatusBadRequest)
		return fmt.Errorf("invalid operation: %s", operation)
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, "Error updating image", http.StatusInternalServerError)
		return fmt.Errorf("error updating image: %v", err)
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Image %s updated successfully", imageID)))
	return nil
}
