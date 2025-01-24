package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/api/meetings", getMeetings)
	r.GET("/api/meetings/:id", getMeetingByID)
	r.POST("/api/meetings", postMeeting)
	r.PUT("/api/meetings/:id", putMeeting)
	r.DELETE("/api/meetings/:id", deleteMeeting)

	r.GET("/api/staffs", getStaffs)
	r.GET("/api/staffs/:id", getStaffByID)
	r.POST("/api/staffs", postStaff)
	r.PUT("/api/staffs/:id", putStaff)
	r.DELETE("/api/staffs/:id", deleteStaff)

	r.GET("/api/locations", getLocations)
	r.GET("/api/locations/:id", getLocationByID)
	r.POST("/api/locations", postLocation)
	r.PUT("/api/locations/:id", putLocation)
	r.DELETE("/api/locations/:id", deleteLocation)

	r.Run("localhost:8080")
}

func getMeetings(c *gin.Context) {
	var meetings []Meeting
	var query = "SELECT Meeting_ID, Location_ID, Title, Description, Meeting_Date, Start_Time, End_Time, Meeting_Type FROM Meeting"

	rows, err := db.Query(query)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer rows.Close()

	for rows.Next() {
		var m Meeting

		if err := rows.Scan(&m.MeetingID, &m.LocationID, &m.Title, &m.Description, &m.MeetingDate, &m.StartTime, &m.EndTime, &m.MeetingType); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}

		meetings = append(meetings, m)
	}

	c.IndentedJSON(http.StatusOK, meetings)
}

func getMeetingByID(c *gin.Context) {
	id := c.Param("id")

	var query = "SELECT Meeting_ID, Location_ID, Title, Description, Meeting_Date, Start_Time, End_Time, Meeting_Type FROM Meeting WHERE Meeting_ID = ?"

	res, err := db.Query(query, id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer res.Close()

	if res.Next() {
		var m Meeting

		if err := res.Scan(&m.MeetingID, &m.LocationID, &m.Title, &m.Description, &m.MeetingDate, &m.StartTime, &m.EndTime, &m.MeetingType); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		}

		c.IndentedJSON(http.StatusOK, m)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "meeting not found"})
	}
}

func postMeeting(c *gin.Context) {
	var newMeeting Meeting

	if err := c.BindJSON(&newMeeting); err != nil {
		return
	}

	var query = "INSERT INTO meeting (Location_ID, Title, Description, Meeting_Date, Start_Time, End_Time, Meeting_Type) VALUES (?, ?, ?, ?, ?, ?, ?)"

	ins, err := db.Query(query, newMeeting.LocationID, newMeeting.Title, newMeeting.Description, newMeeting.MeetingDate, newMeeting.StartTime, newMeeting.EndTime, newMeeting.MeetingType)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "meeting creation failed."})
	} else {
		c.IndentedJSON(http.StatusCreated, newMeeting)
	}

	defer ins.Close()
}

func putMeeting(c *gin.Context) {
	var updatedMeeting Meeting
	updatedMeeting.LocationID = -1

	if err := c.BindJSON(&updatedMeeting); err != nil {
		return
	}

	updated := 0

	id := c.Param("id")

	var query = "UPDATE meeting SET "

	if updatedMeeting.LocationID != -1 {
		query += "Location_ID = " + strconv.Itoa(updatedMeeting.LocationID)
		updated++
	}

	if updatedMeeting.Title != "" {
		if updated > 0 {
			query += ", Title = '" + updatedMeeting.Title + "'"
		} else {
			query += "Title = '" + updatedMeeting.Title + "'"
		}

		updated++
	}

	if updatedMeeting.Description != "" {
		if updated > 0 {
			query += ", Description = '" + updatedMeeting.Description + "'"
		} else {
			query += "Description = '" + updatedMeeting.Description + "'"
		}

		updated++
	}

	if updatedMeeting.MeetingDate != "" {
		if updated > 0 {
			query += ", Meeting_Date = '" + updatedMeeting.MeetingDate + "'"
		} else {
			query += "Meeting_Date = '" + updatedMeeting.MeetingDate + "'"
		}

		updated++
	}

	if updatedMeeting.StartTime != "" {
		if updated > 0 {
			query += ", Start_Time = '" + updatedMeeting.StartTime + "'"
		} else {
			query += "Start_Time = '" + updatedMeeting.StartTime + "'"
		}
	}

	if updatedMeeting.EndTime != "" {
		if updated > 0 {
			query += ", End_Time = '" + updatedMeeting.EndTime + "'"
		} else {
			query += "End_Time = '" + updatedMeeting.EndTime + "'"
		}

		updated++
	}

	if updatedMeeting.MeetingType != "" {
		if updated > 0 {
			query += ", Meeting_Type = '" + updatedMeeting.MeetingType + "'"
		} else {
			query += "Meeting_Type = '" + updatedMeeting.MeetingType + "'"
		}

		updated++
	}

	query += " WHERE Meeting_ID = " + id

	res, err := db.Query(query)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "meeting " + id + " updated."})
	}

	defer res.Close()
}

func deleteMeeting(c *gin.Context) {
	id := c.Param("id")

	var query = "DELETE FROM meeting WHERE Meeting_ID = ?"

	res, err := db.Query(query, id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer res.Close()

	c.IndentedJSON(http.StatusOK, gin.H{"message": "meeting " + id + " deleted."})
}

func getStaffs(c *gin.Context) {
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

func getStaffByID(c *gin.Context) {
	id := c.Param("id")

	var query = "SELECT Staff_ID, Full_Name, Email, Position_ID FROM Staff WHERE Staff_ID = ?"

	res, err := db.Query(query, id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer res.Close()

	if res.Next() {
		var s Staff

		if err := res.Scan(&s.StaffID, &s.FullName, &s.Email, &s.PositionID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		c.IndentedJSON(http.StatusOK, s)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "staff not found"})
	}
}

func postStaff(c *gin.Context) {
	var newStaff Staff

	if err := c.BindJSON(&newStaff); err != nil {
		return
	}

	var query = "INSERT INTO staff (Full_Name, Email, Position_ID) VALUES (?, ?, ?)"

	ins, err := db.Query(query, newStaff.FullName, newStaff.Email, newStaff.PositionID)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "staff creation failed."})
	} else {
		c.IndentedJSON(http.StatusCreated, newStaff)
	}

	defer ins.Close()
}

func putStaff(c *gin.Context) {
	var updatedStaff Staff
	updatedStaff.PositionID = -1

	if err := c.BindJSON(&updatedStaff); err != nil {
		return
	}

	updated := 0

	id := c.Param("id")

	var query = "UPDATE staff SET "

	if updatedStaff.FullName != "" {
		query += "Full_Name = '" + updatedStaff.FullName + "'"

		updated++
		fmt.Println("updated name")
	}

	if updatedStaff.Email != "" {
		if updated > 0 {
			query += ", Email = '" + updatedStaff.Email + "'"
		} else {
			query += "Email = '" + updatedStaff.Email + "'"
		}

		updated++
		fmt.Println("updated email")
	}

	if updatedStaff.PositionID != -1 {
		if updated > 0 {
			query += ", Position_ID = " + strconv.Itoa(updatedStaff.PositionID)
		} else {
			query += "Position_ID = " + strconv.Itoa(updatedStaff.PositionID)
		}

		updated++
		fmt.Println("updated position")
	}

	query += " WHERE Staff_ID = " + id

	res, err := db.Query(query)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	defer res.Close()

	c.IndentedJSON(http.StatusOK, gin.H{"message": "staff " + id + " updated."})
}

func deleteStaff(c *gin.Context) {
	id := c.Param("id")

	var query = "DELETE FROM staff WHERE Staff_ID = ?"

	res, err := db.Query(query, id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	defer res.Close()

	c.IndentedJSON(http.StatusOK, gin.H{"message": "staff " + id + " deleted."})
}

func getLocations(c *gin.Context) {
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

func getLocationByID(c *gin.Context) {
	id := c.Param("id")

	var query = "SELECT Location_ID, Location_Name, Address, Floor FROM Location WHERE Location_ID = ?"

	res, err := db.Query(query, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer res.Close()

	if res.Next() {
		var l Location

		if err := res.Scan(&l.LocationID, &l.LocationName, &l.Address, &l.Floor); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		c.IndentedJSON(http.StatusOK, l)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "location not found"})
	}
}

func postLocation(c *gin.Context) {
	var newLocation Location

	if err := c.BindJSON(&newLocation); err != nil {
		return
	}

	var query = "INSERT INTO location (Location_Name, Address, Floor) VALUES (?, ?, ?)"

	ins, err := db.Query(query, newLocation.LocationName, newLocation.Address, newLocation.Floor)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "location creation failed."})
	} else {
		c.IndentedJSON(http.StatusCreated, newLocation)
	}

	defer ins.Close()
}

func putLocation(c *gin.Context) {
	var updatedLocation Location

	if err := c.BindJSON(&updatedLocation); err != nil {
		return
	}

	updated := 0

	id := c.Param("id")

	var query = "UPDATE location SET "

	if updatedLocation.LocationName != "" {
		query += "Location_Name = '" + updatedLocation.LocationName + "'"

		updated++
	}

	if updatedLocation.Address != "" {
		if updated > 0 {
			query += ", Address = '" + updatedLocation.Address + "'"
		} else {
			query += "Address = '" + updatedLocation.Address + "'"
		}

		updated++
	}

	if updatedLocation.Floor != "" {
		if updated > 0 {
			query += ", Floor = '" + updatedLocation.Floor + "'"
		} else {
			query += "Floor = '" + updatedLocation.Floor + "'"
		}

		updated++
	}

	query += " WHERE Location_ID = " + id

	res, err := db.Query(query)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "location " + id + " updated."})
	}

	defer res.Close()
}

func deleteLocation(c *gin.Context) {
	id := c.Param("id")

	var query = "DELETE FROM location WHERE Location_ID = ?"

	res, err := db.Query(query, id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}

	defer res.Close()

	c.IndentedJSON(http.StatusOK, gin.H{"message": "location " + id + " deleted."})
}
