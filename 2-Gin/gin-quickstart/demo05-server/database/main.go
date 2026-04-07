package main

import (
  "database/sql"
  "net/http"

  "github.com/gin-gonic/gin"
)

func listUsers(db *sql.DB) gin.HandlerFunc {
  return func(c *gin.Context) {
    rows, err := db.QueryContext(c.Request.Context(), "SELECT id, name FROM users")
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
      return
    }
    defer rows.Close()

    type User struct {
      ID   int    `json:"id"`
      Name string `json:"name"`
    }
    var users []User
    for rows.Next() {
      var u User
      if err := rows.Scan(&u.ID, &u.Name); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "scan failed"})
        return
      }
      users = append(users, u)
    }
    c.JSON(http.StatusOK, users)
  }
}

func main() {
  db, _ := sql.Open("postgres", "your-dsn")
  defer db.Close()

  r := gin.Default()
  r.GET("/users", listUsers(db))
  r.Run(":8080")
}