package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Staff struct {
	StaffID    int    `json:"staff_id"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	PositionID int    `json:"position_id"`
}

type Location struct {
	LocationID   int    `json:"location_id"`
	LocationName string `json:"location_name"`
	Address      string `json:"address"`
	Floor        string `json:"floor"`
}

type Meeting struct {
	MeetingID   int    `json:"meeting_id"`
	LocationID  int    `json:"location_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MeetingDate string `json:"meeting_date"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	MeetingType string `json:"meeting_type"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/rapat")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	r := gin.Default()

	r.GET("/api/meetings", getMeetings)
	r.GET("/api/meetings/:id", getMeetingByID)
	r.POST("/api/meetings", postMeeting)
	r.PUT("/api/meetings/:id", putMeeting)
	r.DELETE("/api/meetings/:id", deleteMeeting)

	r.GET("/api/staff", getStaff)

	r.GET("/api/location", getLocation)

	r.Run("localhost:8080")
}

func getMeetings(c *gin.Context) {
	var meetings []Meeting
	var query = "SELECT Meeting_ID, Location_ID, Title, Description, Meeting_Date, Start_Time, End_Time, Meeting_Type FROM Meeting"

	rows, err := db.Query(query)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	for rows.Next() {
		var m Meeting

		if err := rows.Scan(&m.MeetingID, &m.LocationID, &m.Title, &m.Description, &m.MeetingDate, &m.StartTime, &m.EndTime, &m.MeetingType); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}

		meetings = append(meetings, m)
	}

	c.IndentedJSON(http.StatusOK, meetings)
}

func getMeetingByID(c *gin.Context) {}

func postMeeting(c *gin.Context) {}

func putMeeting(c *gin.Context) {}

func deleteMeeting(c *gin.Context) {}

func getStaff(c *gin.Context) {
	staffs := []Staff{}
	var query = "SELECT Staff_ID, Full_Name, Email, Position_ID FROM Staff"

	rows, err := db.Query(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer rows.Close()

	for rows.Next() {
		var s Staff

		if err := rows.Scan(&s.StaffID, &s.FullName, &s.Email, &s.PositionID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		staffs = append(staffs, s)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	if len(staffs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "staffs table is empty"})
	} else {
		c.JSON(http.StatusOK, staffs)
	}
}

func getLocation(c *gin.Context) {
	locations := []Location{}
	var query = "SELECT Location_ID, Location_Name, Address, Floor FROM Location"

	rows, err := db.Query(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer rows.Close()

	for rows.Next() {
		var l Location

		if err := rows.Scan(&l.LocationID, &l.LocationName, &l.Address, &l.Floor); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		locations = append(locations, l)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	if len(locations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "staffs table is empty"})
	} else {
		c.JSON(http.StatusOK, locations)
	}
}
