package controllers

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"glucose-meter-backend/database"
	"glucose-meter-backend/utils"
	"log"
	"net/http"
	"strings"
	"time"
)

// AddData receives a single string of data, splits it, and stores it in the database.
func AddData(c *gin.Context) {
	// Parse the request body to get the data string
	var requestBody struct {
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utils.JsonErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	dataString := requestBody.Data

	// Split the data string
	dataParts := strings.Split(dataString, "|")
	if len(dataParts) != 9 {
		utils.JsonErrorResponse(c, http.StatusBadRequest, "Invalid data format")
		return
	}

	// Extract data fields
	patientName := strings.TrimSpace(dataParts[0])
	day := strings.TrimSpace(dataParts[1])
	month := strings.TrimSpace(dataParts[2])
	year := strings.TrimSpace(dataParts[3])
	patientAge := strings.TrimSpace(dataParts[4])
	patientAddress := strings.TrimSpace(dataParts[5])
	glucoseTime := strings.TrimSpace(dataParts[6])
	glucoseValue := strings.TrimSpace(dataParts[7])
	glucoseStatus := strings.TrimSpace(dataParts[8])

	// Combine and parse the date of birth
	patientDOBStr := fmt.Sprintf("%s %s %s", day, month, year)
	patientDOB, err := time.Parse("2 January 2006", patientDOBStr)
	if err != nil {
		utils.JsonErrorResponse(c, http.StatusBadRequest, "Invalid date of birth format")
		return
	}

	// Parse glucose time
	glucoseTimestamp, err := time.Parse("02/01/2006 15:04:05", glucoseTime)
	if err != nil {
		utils.JsonErrorResponse(c, http.StatusBadRequest, "Invalid glucose time format")
		return
	}

	// Insert data into the database
	query := `INSERT INTO glucose_data (patient_name, patient_date_of_birth, patient_age, patient_address, glucose_time, glucose_value, glucose_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = database.DbPool.Exec(context.Background(), query, patientName, patientDOB, patientAge, patientAddress, glucoseTimestamp, glucoseValue, glucoseStatus)
	if err != nil {
		log.Println("Failed to insert data", err)
		utils.JsonErrorResponse(c, http.StatusInternalServerError, "Failed to insert data")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Data added successfully"})
}

// DownloadData retrieves all data from the database and returns it as a CSV file.
func DownloadData(c *gin.Context) {
	// Query data from the database
	query := `SELECT id, patient_name, patient_date_of_birth, patient_age, patient_address, glucose_time, glucose_value, glucose_status FROM glucose_data`
	rows, err := database.DbPool.Query(context.Background(), query)
	if err != nil {
		log.Println("Failed to query data", err)
		utils.JsonErrorResponse(c, http.StatusInternalServerError, "Failed to query data")
		return
	}
	defer rows.Close()

	// Generate dynamic file name
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("glucose_data_%s.csv", currentTime)

	// Prepare CSV writer
	c.Writer.Header().Set("Content-Type", "text/csv")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
	csvWriter := csv.NewWriter(c.Writer)
	defer csvWriter.Flush()

	// Write CSV header
	csvWriter.Write([]string{"ID", "Patient Name", "Date of Birth", "Age", "Address", "Glucose Time", "Glucose Value", "Glucose Status"})

	// Write rows to CSV
	for rows.Next() {
		var id, patientName, patientAge, patientAddress, glucoseValue, glucoseStatus string
		var patientDOB, glucoseTime time.Time
		if err := rows.Scan(&id, &patientName, &patientDOB, &patientAge, &patientAddress, &glucoseTime, &glucoseValue, &glucoseStatus); err != nil {
			log.Println("Failed to scan data", err)
			utils.JsonErrorResponse(c, http.StatusInternalServerError, "Failed to scan data")
			return
		}
		csvWriter.Write([]string{
			id,
			patientName,
			patientDOB.Format("2006-01-02"),
			patientAge,
			patientAddress,
			glucoseTime.Format("2006-01-02 15:04:05"),
			glucoseValue,
			glucoseStatus,
		})
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows", err)
		utils.JsonErrorResponse(c, http.StatusInternalServerError, "Error iterating rows")
		return
	}
}
