package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
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
	/*
		m =

		if m == nil || len(m) == 0 {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.IndentedJSON(http.StatusOK, m)
		}
	*/
}

func getMeetingByID(c *gin.Context) {}

func postMeeting(c *gin.Context) {}

func putMeeting(c *gin.Context) {}

func deleteMeeting(c *gin.Context) {}

func getStaff(c *gin.Context) {}

func getLocation(c *gin.Context) {}
