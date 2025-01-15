package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Staff struct {
	staffID    int    `json:"staff_id" db:"Staff_ID"`
	fullName   string `json:"full_name" db:"Full_Name"`
	email      string `json:"email" db:"Email"`
	positionID int    `json:"position_id" db:"Position_ID"`
}

type Location struct {
	locationID   int    `json:"location_id"`
	locationName string `json:"location_name"`
	address      string `json:"address"`
	floor        string `json:"floor"`
}

type Meeting struct {
	meetingID   int    `json:"meeting_id"`
	locationID  int    `json:"location_id"`
	title       string `json:"title"`
	description string `json:"description"`
	meetingDate string `json:"meeting_date"`
	startTime   string `json:"start_time"`
	endTime     string `json:"end_time"`
	meetingType string `json:"meeting_type"`
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

		if err := rows.Scan(&m.meetingID, &m.locationID, &m.title, &m.description, &m.meetingDate, &m.startTime, &m.endTime, &m.meetingType); err != nil {
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
	var staffs []Staff
	var query = "SELECT Staff_ID, Full_Name, Email, Position_ID FROM Staff"

	rows, err := db.Query(query)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer rows.Close()

	for rows.Next() {
		var s Staff

		if err := rows.Scan(&s.staffID, &s.fullName, &s.email, &s.positionID); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		fmt.Println(s)
		staffs = append(staffs, s)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, staffs)
}

func getLocation(c *gin.Context) {
	var locations []Location
	var query = "SELECT Location_ID, Location_Name, Address, Floor FROM Location"

	rows, err := db.Query(query)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	for rows.Next() {
		var l Location

		if err := rows.Scan(&l.locationID, &l.locationName, &l.address, &l.floor); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}

		locations = append(locations, l)
	}

	c.JSON(http.StatusOK, locations)
}
