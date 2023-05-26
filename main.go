package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"io/ioutil"
)

type Message struct {
	Sdp string `json:"sdp"`
}

func main() {

	errlf := godotenv.Load(".env")
	if errlf != nil {
		log.Fatal("Error loading .env file")
	}

	gin.SetMode(os.Getenv("GIM_MODE"))

	// Logging to a file.
	if os.Getenv("GIM_MODE") == "release" {
		f, _ := os.Create("gin.log")
		gin.DefaultWriter = io.MultiWriter(f)
	}

	serverKey := os.Getenv("SENTRY_KEY")
	serverPem := os.Getenv("SENTRY_PEM")

	// init API
	fmt.Println("Starting the application ...<3")

	// Start Gin
	port := os.Getenv("SENTRY_PORT")
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 31536000000, //12 * time.Hour,
	}))

	// if port is not define
	if port == "" {
		port = "8000"
	}

	router.GET("/getMessage", GetMessage)
	router.GET("/getServerMessage", GetServerMessage)
	router.POST("/setMessage", SetMessage)
	router.POST("/SetServerMessage", SetServerMessage)

	if len(serverPem) > 0 {
		if err := http.ListenAndServeTLS(":"+port, serverPem, serverKey, router); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := http.ListenAndServe(":"+port, router); err != nil {
			log.Fatal(err)
		}
	}

}

func GetMessage(c *gin.Context) {
	// read the whole file at once
	body, err := ioutil.ReadFile("sdp.txt")
	if err != nil {
		c.IndentedJSON(http.StatusFailedDependency, gin.H{"error": err})
	} else {

		if len(body) > 0 {
			// truncate file
			err2 := ioutil.WriteFile("sdp.txt", []byte(""), 0644)
			if err2 != nil {
				c.IndentedJSON(http.StatusFailedDependency, gin.H{"error": err})
			}
			// truncate file

			c.JSON(http.StatusOK, string(body))
		} else {
			c.JSON(http.StatusFailedDependency, "Error")
		}

	}
}

func SetMessage(c *gin.Context) {
	var messageRequest Message

	if err := c.BindJSON(&messageRequest); err != nil {
		c.JSON(400, gin.H{"Error": "Wrong data"})
		return
	}

	messageUpdate := Message{
		Sdp: messageRequest.Sdp,
	}

	err := ioutil.WriteFile("sdp.txt", []byte(messageUpdate.Sdp), 0644)

	if err != nil {
		c.IndentedJSON(http.StatusFailedDependency, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, "ok")
	}
}

func GetServerMessage(c *gin.Context) {
	// read the whole file at once
	body, err := ioutil.ReadFile("serverResponse.txt")
	if err != nil {
		c.IndentedJSON(http.StatusFailedDependency, gin.H{"error": err})
	} else {

		if len(body) > 0 {
			// truncate file
			err2 := ioutil.WriteFile("serverResponse.txt", []byte(""), 0644)
			if err2 != nil {
				c.IndentedJSON(http.StatusFailedDependency, gin.H{"error": err})
			}
			// truncate file

			c.JSON(http.StatusOK, string(body))
		} else {
			c.JSON(http.StatusFailedDependency, "Error")
		}

	}
}

func SetServerMessage(c *gin.Context) {
	var messageRequest Message

	if err := c.BindJSON(&messageRequest); err != nil {
		c.JSON(400, gin.H{"Error": "Wrong data"})
		return
	}

	messageUpdate := Message{
		Sdp: messageRequest.Sdp,
	}

	err := ioutil.WriteFile("serverResponse.txt", []byte(messageUpdate.Sdp), 0644)

	if err != nil {
		c.IndentedJSON(http.StatusFailedDependency, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, "ok")
	}
}
