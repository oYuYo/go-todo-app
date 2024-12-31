package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type task struct {
	Id                 string    `json:"task_id"`
	TaskName           string    `json:"task_name"`
	Description        string    `json:"description"`
	Status             string    `json:"status"`
	DueDate            time.Time `json:"duedate"`
	CompletionDateTime time.Time `json:"completion_datetime"`
	CreatedDateTime    time.Time `json:"created_datetime"`
	ModifiedDateTIme   time.Time `json:"modified_datetime"`
}

func dbConn() *sql.DB {

	sqldb, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal("Couldn't open database :", err)
	}
	log.Print("Opened database.")

	return sqldb
}

func dbPing(db *sql.DB) {
	if err := db.Ping(); err != nil {
		log.Printf("Database ping error :%v", err)
	}
}

func getTasks(c *gin.Context) {
	var tasks []task

	db := dbConn()
	defer db.Close()

	//dbPing(db)

	rows, err := db.Query("select * from tasks;")
	if err != nil {
		log.Print("Couldn't get tasks :", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Couldn't get tasks"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, taskName, description, status string
		var dudate, createdDatetime time.Time
		var completionDatetime, modifiedDatetime sql.NullTime

		if err := rows.Scan(&id, &taskName, &description, &status, &dudate, &completionDatetime, &createdDatetime, &modifiedDatetime); err != nil {
			log.Print("Failed to row scan :", err)
		}
		fmt.Printf("ID:%s, task:%s, desc:%s, stas:%s, duedate:%s, completion:%s, created:%s, modified:%s", id, taskName, description, status, dudate, completionDatetime.Time, createdDatetime, modifiedDatetime.Time)

		task := task{id, taskName, description, status, dudate, completionDatetime.Time, createdDatetime, modifiedDatetime.Time}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		log.Print("Error occurred in row scanning :", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error occurred in row scanning"})
		return
	}

	c.IndentedJSON(http.StatusOK, tasks)
}

func getTaskById(c *gin.Context) {
	id := c.Param("id")

	db := dbConn()
	defer db.Close()

	row := db.QueryRow("select * from tasks where task_id=$1;", id)

	var taskName, description, status string
	var dudate, createdDatetime time.Time
	var completionDatetime, modifiedDatetime sql.NullTime

	err := row.Scan(&id, &taskName, &description, &status, &dudate, &completionDatetime, &createdDatetime, &modifiedDatetime)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("Couldn't find specified task :", err)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Couldn't find specified task"})
			return
		}
		log.Print("Error occurred in row scanning :", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error occurred in row scanning"})
		return
	}
	task := task{id, taskName, description, status, dudate, completionDatetime.Time, createdDatetime, modifiedDatetime.Time}
	c.IndentedJSON(http.StatusOK, task)
}

func postTask(c *gin.Context) {
	var newTask task

	if err := c.BindJSON(&newTask); err != nil {
		log.Print("Invalid task params :", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid task parameters"})
		return
	}

	db := dbConn()
	defer db.Close()

	_, err := db.Exec("insert into tasks(task_name, description, duedate) values($1, $2, $3);", newTask.TaskName, newTask.Description, newTask.DueDate)
	if err != nil {
		log.Print("Error occurred in insert task :", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error occurred in insert task"})
		return
	}

	c.IndentedJSON(http.StatusOK, newTask)
}

func deleteTaskById(c *gin.Context) {
	id := c.Param("id")

	db := dbConn()
	defer db.Close()

	_, err := db.Exec("delete from tasks where task_id=$1;", id)
	if err != nil {
		log.Print("Error occurred in deleting task :", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error occurred in deleting task"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Succeed to delete task"})
}

func updateTaskById(c *gin.Context) {
	id := c.Param("id")

	var updatedTask task

	if err := c.BindJSON(&updatedTask); err != nil {
		log.Print("Invalid task params :", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid task parameters"})
		return
	}

	db := dbConn()
	defer db.Close()

	_, err := db.Exec(`update tasks 
							set task_name=$1
							, description=$2
							, status=$3
							, duedate=$4
							, completion_datetime=$5
							, modified_datetime=current_timestamp
						where task_id=$6;`,
		updatedTask.TaskName,
		updatedTask.Description,
		updatedTask.Status,
		updatedTask.DueDate,
		updatedTask.CompletionDateTime,
		id)

	if err != nil {
		log.Print("Error occurred in updating task :", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error occurred in updating task"})
		return
	}

	row := db.QueryRow("select * from tasks where task_id=$1;", updatedTask.Id)

	var taskName, description, status string
	var dudate, createdDatetime time.Time
	var completionDatetime, modifiedDatetime sql.NullTime

	err = row.Scan(&id, &taskName, &description, &status, &dudate, &completionDatetime, &createdDatetime, &modifiedDatetime)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("Couldn't find specified task :", err)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Couldn't find specified task"})
			return
		}
		log.Print("Error occurred in row scanning :", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error occurred in row scanning"})
		return
	}
	newTask := task{id, taskName, description, status, dudate, completionDatetime.Time, createdDatetime, modifiedDatetime.Time}
	c.IndentedJSON(http.StatusOK, newTask)
}

func main() {

	router := gin.Default()
	router.GET("/tasks", getTasks)
	router.GET("/task/:id", getTaskById)
	router.POST("/task", postTask)
	router.DELETE("/task/:id", deleteTaskById)
	router.PUT("/task/:id", updateTaskById)

	router.Run(":8080")

}
