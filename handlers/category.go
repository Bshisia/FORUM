package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := GetCategories()
	if err != nil {
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		log.Println("Error fetching categories:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func GetCategories() ([]string, error) {
    rows, err := db.Query("SELECT name FROM post_category")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []string
    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            return nil, err
        }
        categories = append(categories, name)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return categories, nil
}
