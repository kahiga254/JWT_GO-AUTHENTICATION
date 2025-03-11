package controllers

import (
	"go_jwt/helpers"
	"go_jwt/models"
	"go_jwt/database"
	"context"

	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err :=bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func verifyPassword(userPassword string, providePassword string)( bool, string ){
	err := bcrypt.CompareHashAndPassword([]byte(providePassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "email or Password do not match"
		check = false
	}
	return check, msg
}

func Signup()gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}


		_, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while checking for the email"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password
		count,  err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while checking for the phone number"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email or phone number already exists"})
			return
		}

		user.Created_at, _= time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _= time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		userId := user.ID.Hex()
		user.User_id = &userId 
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name , *user.User_type, *user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := "User item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, resultInsertionNumber)
	}

} 

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel() // Ensure context is canceled before returning

		var user models.User
		var foundUser models.User

		// Bind request JSON to user struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate user struct
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if the user exists in the database
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
			return
		}

		// Verify password
		passwordIsValid, msg := verifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		// Check if email is nil
		if foundUser.Email == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		// Generate JWT token and refresh token
		token, refreshToken, _ := helpers.GenerateAllTokens(
			*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, *foundUser.User_id,
		)

		// Update the tokens in the database
		helpers.UpdateAllTokens(token, refreshToken, *foundUser.User_id)

		// Retrieve updated user data
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return user data in response
		c.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers()gin.HandlerFunc{
	return func(c *gin.Context){
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordpage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		matchStage := bson.D{{Key: "$match", Value: bson.D{}}}
		 
		groupStage := bson.D{{
			Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
			},
		}}

		projectStage := bson.D{{
			Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
			},
		}}

	// Executing aggregation query
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Fetching all users
		var allUsers []bson.M
		if err := result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}

		// Sending response
		c.JSON(http.StatusOK, allUsers)
	}
}

func GetUser()gin.HandlerFunc{
	return func(c *gin.Context){
		userId := c.Param("user_id")
		
		if err := helpers.CheckUserType(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}


