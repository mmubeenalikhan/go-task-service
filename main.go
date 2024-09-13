package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    _"github.com/lib/pq"
    "github.com/joho/godotenv"
)

var db *sql.DB

// Task represents the structure of a task
type Task struct {
    Id          int    `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Priority    string `json:"priority"`
    Status      string `json:"status"`
    Deadline    string `json:"deadline"`
}

// Initialize the database connection
func initDB() {
    var err error

    // Load environment variables from the .env file
    err = godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Get environment variables
    host := os.Getenv("HOST")
    user := os.Getenv("USER")
    password := os.Getenv("PASSWORD")
    dbname := os.Getenv("DBNAME")
    port := os.Getenv("DBPORT")
    sslmode := os.Getenv("SSLMODE")

    // Create the connection string
    connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        host, user, password, dbname, port, sslmode)

    // Connect to the database
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to the database:", err)
    }

    // Check if the database is reachable
    if err := db.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }

    log.Println("Connected to the database successfully")
}

// Create tasks in bulk
func createBulkTasks(w http.ResponseWriter, r *http.Request) {
    var tasks []Task
    if err := json.NewDecoder(r.Body).Decode(&tasks); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    tx, err := db.Begin()
    if err != nil {
        http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    stmt, err := tx.Prepare(`
        INSERT INTO "Tasks" (title, description, priority, status, deadline, "createdAt", "updatedAt")
        VALUES ($1, $2, $3, $4, $5, now(), now())
    `)
    if err != nil {
        http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    for _, task := range tasks {
        if _, err := stmt.Exec(task.Title, task.Description, task.Priority, task.Status, task.Deadline); err != nil {
            http.Error(w, "Failed to execute statement", http.StatusInternalServerError)
            return
        }
    }

    if err := tx.Commit(); err != nil {
        http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
        return
    }
	
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(tasks)
}

func main() {
    initDB()
	appPort := os.Getenv("PORT")
    if appPort == "" {
        appPort = "8081" 
    }

    http.HandleFunc("/tasks/bulk", createBulkTasks)

    log.Printf("Go service running on port %s\n", appPort)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", appPort), nil))
}
