# Gormongo

Gormongo is a lightweight Go package designed to simplify MongoDB integration for applications built with the Gin framework.

It provides a clean and scalable structure for:

* MongoDB connection management
* Collection handling
* CRUD operations
* Context support with Gin
* Modular API architecture
* High-performance REST APIs

---

# Installation

```bash
go get github.com/mohametdiatta/gormongo
```

---

# Requirements

* Go 1.21+
* MongoDB
* Gin Framework

---

# Getting Started

## Import Package

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/mohametdiatta/gormongo"
)
```

---

## Initialize MongoDB Connection

```go
package main

import (
    "log"

    "github.com/gin-gonic/gin"
    "github.com/mohametdiatta/gormongo"
)

func main() {
	url := os.Getenv("MONGO_DB_URL")
	client, err := gormongo.Connect(url)
	if err != nil {
		panic(err)
	}

	db := client.Database("sample_mflix")
	registry := gormongo.NewRegistry(context.Background(), db)
	registry.Register("comments", &models.Commentschema{}, "comments")

	app := gormongo.NewApp(registry)
	
 	app.Router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	app.Run(":3000")
}
```

---

# Define a Model

```go
package models

import "github.com/mohametdiatta/gormongo"

type Commentschema struct {
	gormongo.Model
	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Text  string `bson:"text" json:"text"`
}

var Comment = &gormongo.Model{}

```

---

# Access a Collection

```go
	userModel, err := gormongo.GetModel[*models.Userschema](a.Registry, "users")
```

---

# CRUD Examples

## Create

```go
data := User{
    Name:  "Mohamet",
    Email: "mohamet@example.com",
}

_, err := userModel.Create( user)
```

---

## Find One

```go
var user User

err := userModel.FindOne( bson.M{
    "email": "mohamet@example.com",
})
```

---

## Update

```go
_, err := userModel.UpdateOne(
    bson.M{"email": "mohamet@example.com"},
    bson.M{
        "$set": bson.M{
            "name": "Mohamet Diatta",
        },
    },
)
```

---

## Delete

```go
_, err := userModel.DeleteOne( bson.M{
    "email": "mohamet@example.com",
})
```

---

# Gin Integration Example

```go
package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/mohametdiatta/gormongo"
    "go.mongodb.org/mongo-driver/bson"
)

func GetUsers(a *gormongo.App) {
	return func(c *gin.Context) {
	comment, err := gormongo.GetModel[*models.Commentschema](a.Registry, "comments")

   var comments []models.Commentschema
 
		err := comment.FindAll(bson.D{}, &comments)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
 
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"count":   len(comments),
			"data":    comments,
		})
}}
```

---


---

# Features

* Native MongoDB driver support
* Seamless Gin integration
* Simplified collection management
* HTTP context support
* Fast CRUD operations
* Modular architecture
* Lightweight and performant

---

# Roadmap

* [ ] Built-in validation
* [ ] Automatic pagination
* [ ] Soft delete support
* [ ] Lifecycle hooks
* [ ] Query builder
* [ ] Redis cache integration
* [ ] Simplified transactions

---

# Contributing

Contributions are welcome.

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push your branch
5. Open a Pull Request

---

# License

MIT License

---

# Author

Mohamet Diatta
