package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Cyrus-Gahatraj/horogo/internal"
	"github.com/spf13/cobra"
)

var reader = bufio.NewReader(os.Stdin)

func getInput(prompt string) string {
	fmt.Print(prompt)
	val, _ := reader.ReadString('\n')
	return strings.TrimSpace(val)
}

func runCmd(cmd *cobra.Command, args []string) {
	name := getInput("Your name (default: User): ")
	if name == "" {
		name = "User"
	}

	dob := getInput("Your DOB (format: yyyy-mm-dd): ")
	splitDOB := strings.Split(dob, "-")
	if len(splitDOB) != 3 {
		fmt.Println("Invalid date format, expected yyyy-mm-dd")
		os.Exit(1)
	}
	year, err := strconv.Atoi(splitDOB[0])
	if err != nil {
		fmt.Println("Invalid year:", splitDOB[0])
		os.Exit(1)
	}
	month, err := strconv.Atoi(splitDOB[1])
	if err != nil {
		fmt.Println("Invalid month:", splitDOB[1])
		os.Exit(1)
	}
	day, err := strconv.Atoi(splitDOB[2])
	if err != nil {
		fmt.Println("Invalid day:", splitDOB[2])
		os.Exit(1)
	}

	tob := getInput("Your birth time (format: HH:MM, 24-hour): ")
	splitTOB := strings.Split(tob, ":")
	if len(splitTOB) != 2 {
		fmt.Println("Invalid time format, expected HH:MM")
		os.Exit(1)
	}
	hour, err := strconv.Atoi(splitTOB[0])
	if err != nil || hour < 0 || hour > 23 {
		fmt.Println("Invalid hour:", splitTOB[0])
		os.Exit(1)
	}
	minute, err := strconv.Atoi(splitTOB[1])
	if err != nil || minute < 0 || minute > 59 {
		fmt.Println("Invalid minute:", splitTOB[1])
		os.Exit(1)
	}

	place := getInput("Your birth place (format: city, country): ")

	lat, lon, err := internal.GeocodePlace(place)
	if err != nil {
		fmt.Println("Geocoding failed:", err)
		os.Exit(1)
	}
	fmt.Printf("Location: %.2f°N, %.2f°E\n", lat, lon)

	tzOffset, err := internal.GetTimezoneOffset(lat, lon)
	if err != nil {
		fmt.Println("Timezone lookup failed:", err)
		os.Exit(1)
	}

	person := internal.Person{
		Name:     name,
		Year:     year,
		Month:    month,
		Day:      day,
		Hour:     hour,
		Minute:   minute,
		Second:   0,
		Lat:      lat,
		Lon:      lon,
		TZOffset: tzOffset,
	}

	person.GetPlanetryPosition()
}

var rootCmd = &cobra.Command{
	Use:   "horogo",
	Short: "Analyze birth chart from CLI",
	Long:  `Horogo is a CLI tool for analyzing birth charts.`,
	Run:   runCmd,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
