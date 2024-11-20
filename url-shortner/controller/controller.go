package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"urlshortner/cache"
	"urlshortner/db"
	"urlshortner/helper"
	"urlshortner/model"
	"urlshortner/ratelimiter"

	"go.mongodb.org/mongo-driver/bson"
)

const dbName = "url-shortner"
const collName = "url2"

var snowFlake *helper.Snowflake

// Initialize Redis client

func snowFlakeGenerator() int64 {
	if snowFlake == nil {
		s, err := helper.NewSnowflake(1)
		if err != nil {
			fmt.Println("Error creating snowflake", err)
			panic(err)
		}
		snowFlake = s
	}
	return snowFlake.GenerateID()
}

// momgo helpers
func InsertOne(user model.User) {
	// Generate ID

	fmt.Println("Generated ID:", user.ID)

	collection, err := db.ConnectDb(dbName, collName)

	if err != nil {
		fmt.Println("Error connecting to MongoDB", err)
		panic(err)
	}

	inseted, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		panic(err)
	}
	fmt.Println(inseted)
	fmt.Println("Inserted ID:", inseted.InsertedID)
}

func FindOne(ShortUrl string) (model.User, error) {
	collection, err := db.ConnectDb(dbName, collName)

	if err != nil {
		fmt.Println("Error connecting to MongoDB", err)
		panic(err)
	}

	var user model.User
	err = collection.FindOne(context.Background(), bson.M{"short_url": ShortUrl}).Decode(&user)
	if err != nil {
		fmt.Println("Error finding one", err)
		return model.User{}, err
	}

	return user, nil
}

func FindOneByLong(LongUrl string) (model.User, error) {
	collection, err := db.ConnectDb(dbName, collName)

	if err != nil {
		fmt.Println("Error connecting to MongoDB", err)
		panic(err)
	}

	var user model.User
	err = collection.FindOne(context.Background(), bson.M{"long_url": LongUrl}).Decode(&user)
	if err != nil {
		fmt.Println("Error finding one", err)
		return model.User{}, err
	}

	return user, nil
}

func FindOneByIndex(ID int64) model.User {
	collection, err := db.ConnectDb(dbName, collName)

	if err != nil {
		fmt.Println("Error connecting to MongoDB", err)
		panic(err)
	}

	var user model.User
	err = collection.FindOne(context.Background(), bson.M{"_id": ID}).Decode(&user)
	if err != nil {
		fmt.Println("Error finding one", err)
		panic(err)
	}
	fmt.Println("User:", user.ShortUrl)
	return user
}

func GetClicks(LongUrl string) int {
	collection, err := db.ConnectDb(dbName, collName)

	if err != nil {
		fmt.Println("Error connecting to MongoDB", err)
		panic(err)
	}

	var user model.User
	err = collection.FindOne(context.Background(), model.User{LongUrl: LongUrl}).Decode(&user)
	if err != nil {
		fmt.Println("Error finding one", err)
		panic(err)
	}

	return user.Clicks
}

func updateClicks(LongUrl string) {
	collection, err := db.ConnectDb(dbName, collName)

	if err != nil {
		fmt.Println("Error connecting to MongoDB", err)
		panic(err)
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"long_url": LongUrl}, model.User{Clicks: GetClicks(LongUrl) + 1})
	if err != nil {
		fmt.Println("Error updating clicks", err)
		panic(err)

	}
}

// Handling all requests

func ConvertUrl(w http.ResponseWriter, r *http.Request) {
	clientIP := helper.GetClientIP(r)
	fmt.Println(clientIP)
	// Getting client IP
	rate, er := ratelimiter.RateLimiting(r)

	if er != nil {
		helper.WriteResponse(w, "Error in rate limiting")
		return
	}

	if !rate {
		helper.WriteResponse(w, "Rate limit exceeded")
		return
	}
	fmt.Println("Hello i am here ", rate)

	fmt.Println("Hello i am here ")
	w.Header().Set("Content-Type", "application/json")

	// Extracting url for body and coverting it
	if r.Body == nil {
		helper.WriteResponse(w, "Please provide a URL")
		return
	}

	var temp model.User

	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		helper.WriteResponse(w, "Error decoding request body")
		return
	}

	if temp.LongUrl == "" {
		helper.WriteResponse(w, "Please provide a URL")
		return
	}

	fmt.Println("Long URL:", temp.LongUrl)

	// Check if long url already exists
	user, err := FindOneByLong(temp.LongUrl)
	if err == nil {
		fmt.Println("Long URL already exists")
		cache.Set(model.User(user))
		helper.WriteResponse(w, model.User(user))
		return
	}

	// Convert to snowFlake and short url
	temp.ID = snowFlakeGenerator()
	temp.ShortUrl = helper.ToBase64(temp.ID)

	temp.Clicks = 0

	// Insert into DB
	InsertOne(temp)
	cache.Set(model.User(temp))
	newUser := FindOneByIndex(temp.ID)
	helper.WriteResponse(w, model.User(newUser))

	return
}

func GetLongUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extracting short url from request
	shortUrl := r.URL.Path[1:]
	fmt.Println("Short URL:", shortUrl)

	user, err := cache.Get(shortUrl)
	if err != nil {
		// Find in DB
		user, err := FindOne(shortUrl)
		if err != nil {
			helper.WriteResponse(w, "Short URL not found")
			return
		}
		// Cache the user
		cache.Set(model.User(user))
		fmt.Println("User:", user.LongUrl)

		// Update clicks
		// updateClicks(user.LongUrl)

		// Redirect to long url
		http.Redirect(w, r, user.LongUrl, http.StatusMovedPermanently)
	}
	fmt.Println("Cache hit for short URL")
	http.Redirect(w, r, user.LongUrl, http.StatusMovedPermanently)

}
