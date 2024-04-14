package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Task represents a todo task
type Task struct {
	ID          int    `json:"id,omitempty"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
	Day         string `json:"-"`
}

var currentDay string

func main() {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create tasks table if it does not exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			description TEXT,
			done BOOLEAN,
			day TEXT
		)`)
	if err != nil {
		log.Fatal(err)
	}

	// Get the current day
	currentDay = time.Now().Format("Monday")

	// Start a goroutine to continuously update progress bar
	go updateProgressBar(db)

	r := mux.NewRouter()

	// Define routes
	fs := http.FileServer(http.Dir("./web/build"))

	//http.Handle("/", fs)
	r.HandleFunc("/tasks", getTasks(db)).Methods("GET")
	r.HandleFunc("/addTask", addTaskHandler(db)).Methods("POST")
	r.HandleFunc("/toggleTask", toggleTaskHandler(db)).Methods("POST")
	r.HandleFunc("/editTask", editTaskHandler(db)).Methods("POST")
	r.HandleFunc("/removeTask", removeTaskHandler(db)).Methods("POST")
	r.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	// Start the server
	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", r)
}

// getTasks returns the tasks for a given day
func getTasks(db *sql.DB) http.HandlerFunc {
	fmt.Printf("current day")
	return func(w http.ResponseWriter, r *http.Request) {
		day := time.Now().Format("Monday")
		fmt.Printf("current day %s", day)

		rows, err := db.Query("SELECT id, description, done FROM tasks WHERE day=?", day)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var task Task
			err := rows.Scan(&task.ID, &task.Description, &task.Done)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}

		json.NewEncoder(w).Encode(tasks)
	}
}

// addTask adds a new task for a given day
func addTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		description := r.FormValue("description")
		day := r.FormValue("day")
		_, err := db.Exec("INSERT INTO tasks(description, done, day) VALUES(?, ?, ?)", description, false, day)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// toggleTaskHandler handles the submission of the toggle task form
func toggleTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		description := r.FormValue("description")
		day := r.FormValue("day")
		_, err := db.Exec("UPDATE tasks SET done = NOT done WHERE description = ? AND day = ?", description, day)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// editTaskHandler handles the submission of the edit task form
func editTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oldDescription := r.FormValue("oldDescription")
		newDescription := r.FormValue("newDescription")
		day := r.FormValue("day")
		_, err := db.Exec("UPDATE tasks SET description = ? WHERE description = ? AND day = ?", newDescription, oldDescription, day)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// removeTaskHandler handles the submission of the remove task form
func removeTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		description := r.FormValue("description")
		day := r.FormValue("day")
		_, err := db.Exec("DELETE FROM tasks WHERE description = ? AND day = ?", description, day)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// updateProgressBar continuously updates the progress bar
func updateProgressBar(db *sql.DB) {
	//prevTaskCount := 0

	for {
		// Get the count of completed tasks for the current day
		var taskCount int
		err := db.QueryRow("SELECT COUNT(*) FROM tasks WHERE done=? AND day=?", true, currentDay).Scan(&taskCount)
		if err != nil {
			log.Println("Error getting task count:", err)
		}
		progress := float64(taskCount) / float64(len(TasksForCurrentDay(db))) * 100

		// Clear the console
		fmt.Print("\033[H\033[2J")

		// Print the progress bar
		fmt.Printf("Progress for %s:\n", currentDay)
		fmt.Printf("[")
		for i := 0; i < 50; i++ {
			if float64(i)*2 < progress {
				fmt.Printf("=")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("] %.2f%%\n", progress)

		// Sleep for some time before updating again
		time.Sleep(5 * time.Second)
	}
}

// TasksForCurrentDay returns the tasks for the current day
func TasksForCurrentDay(db *sql.DB) []Task {
	rows, err := db.Query("SELECT id, description, done FROM tasks WHERE day=?", currentDay)
	if err != nil {
		log.Println("Error getting tasks for current day:", err)
		return nil
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Description, &task.Done)
		if err != nil {
			log.Println("Error scanning task row:", err)
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks
}

// Update progress bar
