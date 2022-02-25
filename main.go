package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"tecace/cache"
	"tecace/googleapi"
	"tecace/response"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

const (
	SERVER_PORT = "SERVER_PORT"
)

const (
	ROUTE_PATH_GET_ALL         = "/all"
	ROUTE_PATH_POST_DATA       = "/data"
	ROUTE_PATH_DELETE_DATA_KEY = "/data/:key"
)

const (
	DEV_ENV_FILE = "dev.env"
	KEY          = "key"
)

const (
	ACCEPT_APPLICATION_JSON = "application/json"
)

func main() {
	// Load environment file
	// For the api to work, must have an env file with SERVER_POST=<port number>
	envFile := DEV_ENV_FILE
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("error loading %s file\n", envFile)
	}

	// Establish APIs
	app := fiber.New()

	app.Get(ROUTE_PATH_GET_ALL, getAll)
	app.Post(ROUTE_PATH_POST_DATA, postData)
	app.Delete(ROUTE_PATH_DELETE_DATA_KEY, deleteData)

	app.Listen(fmt.Sprintf(":%s", os.Getenv(SERVER_PORT)))
}

func getAll(c *fiber.Ctx) error {
	result, err := googleapi.Get()
	if err != nil {
		return c.JSON(response.Response{Result: fiber.StatusInternalServerError, Description: "unable to retrieve data from sheet"})
	}

	return c.JSON(response.ResponseData{Result: fiber.StatusOK, Data: result})
}

func postData(c *fiber.Ctx) error {
	c.Accepts(ACCEPT_APPLICATION_JSON)
	var request = make(map[string]string)
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		return c.JSON(response.Response{Result: fiber.StatusInternalServerError, Description: "unable to parse body"})
	}

	for k, v := range request {
		if cache.HasKey(k) {
			err := googleapi.Update(k, v)
			if err != nil {
				return c.JSON(response.Response{Result: fiber.StatusInternalServerError, Description: "unable to update key value pair"})
			}
			cache.AddKey(k)
		} else {
			err := googleapi.Post(k, v)
			if err != nil {
				return c.JSON(response.Response{Result: fiber.StatusInternalServerError, Description: "unable to insert key value pair"})
			}
			cache.AddKey(k)
		}
	}

	return c.JSON(response.Response{Result: fiber.StatusOK, Description: response.OK})
}

func deleteData(c *fiber.Ctx) error {
	key := c.Params(KEY)
	if cache.HasKey(key) {
		err := googleapi.Delete(key)
		if err != nil {
			return c.JSON(response.Response{Result: fiber.StatusNotFound, Description: err.Error()})
		}
		cache.RemoveKey(key)
		return c.JSON(response.Response{Result: fiber.StatusOK, Description: response.OK})
	}
	return c.JSON(response.Response{Result: fiber.StatusNotFound, Description: response.KEY_NOT_FOUND})
}
