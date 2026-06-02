package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Cyrus-Gahatraj/horogo/internal"
	"github.com/spf13/cobra"
)

var reader = bufio.NewReader(os.Stdin)
const dataDir = "data"

func getInput(prompt string) string {
	fmt.Print(prompt)
	val, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			trimmed := strings.TrimSpace(val)
			if trimmed != "" {
				return trimmed
			}
			fmt.Println("\nGoodbye!")
			os.Exit(0)
		}
		return ""
	}
	return strings.TrimSpace(val)
}

func getProfiles() ([]string, error) {
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, err
	}

	var profiles []string
	for _, file := range files {
		if file.IsDir() {
			profiles = append(profiles, file.Name())
		}
	}
	return profiles, nil
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

	chart := person.GetPlanetryPosition()
	chart.Place = place

	outputDir := dataDir
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	nameDir := outputDir + "/" + name
	rawDir := nameDir + "/raw"
	if _, err := os.Stat(rawDir); os.IsNotExist(err) {
		os.MkdirAll(rawDir, os.ModePerm)
	}

	byte, err := json.MarshalIndent(chart, "", " ")
	if err != nil {
		panic(err)
	}

	wholePath := rawDir + "/" + "chart.json"
	err = os.WriteFile(wholePath, byte, 0644)
	if err != nil {
		panic(err)
	}
}

func runLsCmd(cmd *cobra.Command, args []string) {
	profiles, err := getProfiles()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Total profiles found: 0 (data directory does not exist)")
			return
		}
		panic(err)
	}

	fmt.Printf("Total profiles found: %d\n", len(profiles))
	for _, p := range profiles {
		fmt.Printf("- %s\n", p)
	}
}

func runAskCmd(cmd *cobra.Command, args []string) {
	var selectedProfile string

	if len(args) > 0 {
		selectedProfile = args[0]
	} else {
		profiles, err := getProfiles()
		if err != nil || len(profiles) == 0 {
			fmt.Println("No profiles available to analyze.")
			return
		}

		fmt.Println("Available Profiles:")
		for i, profile := range profiles {
			fmt.Printf("[%d] %s\n", i+1, profile)
		}

		choiceStr := getInput("Select a profile number: ")
		choice, err := strconv.Atoi(choiceStr)
		if err != nil || choice < 1 || choice > len(profiles) {
			fmt.Println("Invalid profile selection.")
			return
		}
		selectedProfile = profiles[choice-1]
	}

	chartPath := fmt.Sprintf("%s/%s/raw/chart.json", dataDir, selectedProfile)
	chartBytes, err := os.ReadFile(chartPath)
	if err != nil {
		// Fallback to legacy path for backward compatibility
		legacyPath := fmt.Sprintf("%s/%s/chart.json", dataDir, selectedProfile)
		var legacyErr error
		chartBytes, legacyErr = os.ReadFile(legacyPath)
		if legacyErr != nil {
			fmt.Printf("Profile configuration error: Could not read file at %s or %s\n", chartPath, legacyPath)
			return
		}
	}

	fmt.Printf("\nLoaded chart context for: %s\n", selectedProfile)

	fmt.Println(`Commands:
  help or /help    show this help
  exit or /exit    exit the session
`)

	for {
		userPrompt := getInput("\nAsk a question about this chart: ")
		cleanPrompt := strings.ToLower(strings.TrimSpace(userPrompt))
		if cleanPrompt == "" {
			continue
		}
		if cleanPrompt == "/exit" || cleanPrompt == "\\exit" || cleanPrompt == "exit" {
			fmt.Println("Goodbye!")
			break
		}
		if cleanPrompt == "/help" || cleanPrompt == "\\help" || cleanPrompt == "help" {
			fmt.Println(`Commands:
  help or /help    show this help
  exit or /exit    exit the session`)
			continue
		}

		err := internal.Ask(userPrompt, string(chartBytes))
		if err != nil {
			fmt.Printf("\nError: %v\n", err)
			continue
		}
		fmt.Println()
	}
}

var rootCmd = &cobra.Command{
	Use:   "horogo",
	Short: "Analyze birth chart from CLI",
	Long:  `Horogo is a CLI tool for analyzing birth charts.`,
	Run:   runCmd,
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List total profile structures inside data folder",
	Run:   runLsCmd,
}

var askCmd = &cobra.Command{
	Use:   "ask [name]",
	Short: "Query internal models regarding generated chart data profiles",
	Run:   runAskCmd,
}

func init() {
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(askCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
