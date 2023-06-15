package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

const totalWorker = 10

var jobs = make(chan []interface{}, 0)
var wg = new(sync.WaitGroup)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = []Book{
	{ID: "1", Title: "Harry Potter", Author: "J. K. Rowling"},
	{ID: "2", Title: "The Lord of the Rings", Author: "J. R. R. Tolkien"},
	{ID: "3", Title: "The Wizard of Oz", Author: "L. Frank Baum"},
}

func main() {
	r := gin.New()

	go dispatchWorkers(jobs, wg)

	r.GET("/books", listBooksHandler)
	r.POST("/books", createBookHandler)
	r.DELETE("/books/:id", deleteBookHandler)

	wg.Wait()

	r.Run()
}

func dispatchWorkers(jobs <-chan []interface{}, wg *sync.WaitGroup) {
	for workerIndex := 0; workerIndex <= totalWorker; workerIndex++ {
		go func(workerIndex int, jobs <-chan []interface{}, wg *sync.WaitGroup) {
			counter := 0

			for job := range jobs {
				doTheJob(workerIndex, counter, job)
				wg.Done()
				counter++
			}
		}(workerIndex, jobs, wg)
	}
}

func doTheJob(workerIndex, counter int, values []interface{}) {
	for {
		var outerError error
		func(outerError *error) {
			defer func() {
				if err := recover(); err != nil {
					*outerError = fmt.Errorf("%v", err)
				}
			}()

			//log.Println(values)
		}(&outerError)
		if outerError == nil {
			break
		}
	}

	if counter%10 == 0 {
		log.Println("=> worker", workerIndex, "inserted", counter, "data")
	}
}

func listBooksHandler(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}

func createBookHandler(c *gin.Context) {
	var book Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newTransaction := make([]interface{}, 0)
	newTransaction = append(newTransaction, book)

	wg.Add(1)
	jobs <- newTransaction
	//books = append(books, book)

	c.JSON(http.StatusCreated, nil)
}

func deleteBookHandler(c *gin.Context) {
	id := c.Param("id")

	for i, a := range books {
		if a.ID == id {
			books = append(books[:i], books[i+1:]...)
			break
		}
	}

	c.Status(http.StatusNoContent)
}
