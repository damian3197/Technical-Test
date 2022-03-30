package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

type Student struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Variable to save Student struct
var students = []Student{}

func main() {
	r := mux.NewRouter()
	studentsR := r.PathPrefix("/student").Subrouter()
	studentsR.Path("").Methods(http.MethodGet).HandlerFunc(GetAllStudents)
	studentsR.Path("").Methods(http.MethodPost).HandlerFunc(RegisterStudent)
	studentsR.Path("/{id}").Methods(http.MethodGet).HandlerFunc(GetStudent)
	studentsR.Path("/{id}").Methods(http.MethodPut).HandlerFunc(UpdateStudent)
	studentsR.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(DeleteStudent)

	fmt.Println("Welcome to Mentari Kindergarten School ")
	fmt.Println(http.ListenAndServe(":8080", r))
}

// Function to get all data in student slices
func GetAllStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(students); err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

// Function to loop over the slices of students and returns the index with a matching field
func indexByID(students []Student, id string) int {
	for i := 0; i < len(students); i++ {
		if students[i].ID == id {
			return i
		}
	}
	return -1
}

// REGISTER STUDENT INTO THE SYSTEM
func RegisterStudent(w http.ResponseWriter, r *http.Request) {
	u := Student{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		fmt.Println(err)
		http.Error(w, "Error decoding response object", http.StatusBadRequest)
		return
	}
	students = append(students, u)
	response, err := json.Marshal(&u)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// UPDATE STUDENT BY ID
func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	index := indexByID(students, id)
	if index < 0 {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	u := Student{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		fmt.Println(err)
		http.Error(w, "Error decoding response object", http.StatusBadRequest)
		return
	}
	students[index] = u
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&u); err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

// GET STUDENT BY ID
func GetStudent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	index := indexByID(students, id)
	if index < 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(students[index]); err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding response object", http.StatusInternalServerError)
	}
}

// DELETE STUDENT BY ID
func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	index := indexByID(students, id)
	if index < 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	students = append(students[:index], students[index+1:]...)
	w.WriteHeader(http.StatusOK)
}
