package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

type SensorData struct {
	Temp     float64 `json:"temp"`
	Dust     float64     `json:"dust"`
	Humidity float64     `json:"humidity"`
}

var sensorDataFile = "sensordata.json"
var mu sync.Mutex // Mutex for concurrent access to the file

func main() {
	router := gin.Default()

    router.POST("/api/sensor", func(c *gin.Context) {
        var sensorData SensorData

        // Bind JSON to the SensorData struct
        if err := c.ShouldBindJSON(&sensorData); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Lock the mutex before writing to the file
        mu.Lock()
        defer mu.Unlock()

        // Read existing data from file
        existingData, err := readSensorDataFromFile()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read sensor data file"})
            return
        }

        // Append the new data
        existingData = append(existingData, sensorData)

        // Write updated data back to the file
        err = writeSensorDataToFile(existingData)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write sensor data to file"})
            return
        }

        // Print the received sensor data
        fmt.Printf("Received Sensor Data: %+v\n", sensorData)

        // Print the updated data
        fmt.Printf("Updated Sensor Data: %+v\n", existingData)

        c.JSON(http.StatusOK, gin.H{"message": "Sensor data received successfully"})
    })



	// Define a GET endpoint for retrieving sensor data
	router.GET("/api/sensor", func(c *gin.Context) {
		// Lock the mutex before reading from the file
		mu.Lock()
		defer mu.Unlock()

		// Read data from file
		sensorData, err := readSensorDataFromFile()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read sensor data file"})
			return
		}

		c.JSON(http.StatusOK, sensorData)
	})

	// Run the server on port 8080
	router.Run(":8080")
}

// readSensorDataFromFile reads sensor data from the file
func readSensorDataFromFile() ([]SensorData, error) {
	data, err := ioutil.ReadFile(sensorDataFile)
	if err != nil {
		// If the file doesn't exist yet, return an empty slice
		if os.IsNotExist(err) {
			return []SensorData{}, nil
		}
		return nil, err
	}

	var sensorData []SensorData
	err = json.Unmarshal(data, &sensorData)
	if err != nil {
		return nil, err
	}

	return sensorData, nil
}

// writeSensorDataToFile writes sensor data to the file
func writeSensorDataToFile(sensorData []SensorData) error {
	data, err := json.Marshal(sensorData)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(sensorDataFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

