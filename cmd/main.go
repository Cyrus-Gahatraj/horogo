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

	var gender string
	for {
		genderInput := strings.ToLower(strings.TrimSpace(getInput("Your gender (M/F): ")))
		if genderInput == "m" {
			gender = "male"
			break
		}
		if genderInput == "f" {
			gender = "female"
			break
		}
		fmt.Println("Wrong value, please try again.")
	}

	var year, month, day int
	for {
		dob := getInput("Your DOB (format: yyyy-mm-dd): ")
		splitDOB := strings.Split(dob, "-")
		if len(splitDOB) != 3 {
			fmt.Println("Wrong value, please try again.")
			continue
		}
		y, err1 := strconv.Atoi(splitDOB[0])
		m, err2 := strconv.Atoi(splitDOB[1])
		d, err3 := strconv.Atoi(splitDOB[2])
		if err1 != nil || err2 != nil || err3 != nil || m < 1 || m > 12 || d < 1 || d > 31 {
			fmt.Println("Wrong value, please try again.")
			continue
		}
		year, month, day = y, m, d
		break
	}

	var hour, minute int
	for {
		tob := getInput("Your birth time (format: HH:MM, 24-hour): ")
		splitTOB := strings.Split(tob, ":")
		if len(splitTOB) != 2 {
			fmt.Println("Wrong value, please try again.")
			continue
		}
		h, err1 := strconv.Atoi(splitTOB[0])
		m, err2 := strconv.Atoi(splitTOB[1])
		if err1 != nil || err2 != nil || h < 0 || h > 23 || m < 0 || m > 59 {
			fmt.Println("Wrong value, please try again.")
			continue
		}
		hour, minute = h, m
		break
	}

	var place string
	var lat, lon float64
	var tzOffset int
	for {
		place = getInput("Your birth place (format: city, country): ")
		var err error
		lat, lon, err = internal.GeocodePlace(place)
		if err != nil {
			fmt.Println("Wrong value, please try again.")
			continue
		}
		fmt.Printf("Location: %.2f°N, %.2f°E\n", lat, lon)

		tzOffset, err = internal.GetTimezoneOffset(lat, lon)
		if err != nil {
			fmt.Println("Timezone lookup failed:", err)
			continue
		}
		break
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
	chart.Gender = gender

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

	// Generate the wiki files by calling the server
	fmt.Println("\nGenerating your personalized astrology wiki... (this may take a minute)")
	generateWiki(name, place, string(byte))
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

		// Check if we can answer from the local wiki first!
		matchedWikiFile := matchWikiCategory(cleanPrompt)
		if matchedWikiFile != "" {
			wikiPath := fmt.Sprintf("%s/%s/wiki/%s", dataDir, selectedProfile, matchedWikiFile)
			content, err := os.ReadFile(wikiPath)
			if err == nil {
				fmt.Printf("\n[Reading from local Wiki: %s]\n\n", matchedWikiFile)
				fmt.Println(string(content))
				continue
			}
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

func matchWikiCategory(prompt string) string {
	// Personality keywords
	personalityKeywords := []string{
		"personality", "ascendant", "lagna", "myself", "who am i",
		"temperament", "behavior", "krittika", "nakshatra", "constitution",
	}
	for _, kw := range personalityKeywords {
		if strings.Contains(prompt, kw) {
			return "ascendant_personality.md"
		}
	}

	// Planets and houses keywords
	planetKeywords := []string{
		"sun", "moon", "mars", "mercury", "jupiter", "saturn",
		"venus", "uranus", "neptune", "pluto", "house", "degree",
		"placement", "aspect", "conjunction",
	}
	for _, kw := range planetKeywords {
		if strings.Contains(prompt, kw) {
			return "planets_houses.md"
		}
	}

	// Career and wealth keywords
	careerKeywords := []string{
		"career", "job", "work", "wealth", "money", "finance",
		"profession", "business", "occupation", "earning", "rich",
	}
	for _, kw := range careerKeywords {
		if strings.Contains(prompt, kw) {
			return "career_wealth.md"
		}
	}

	// Relationships and spirituality keywords
	relKeywords := []string{
		"relationship", "love", "marriage", "spouse", "partner",
		"wife", "husband", "spirituality", "growth", "inner",
		"occult", "mystic", "meditation", "god", "soul",
	}
	for _, kw := range relKeywords {
		if strings.Contains(prompt, kw) {
			return "relationships_spirituality.md"
		}
	}

	return ""
}

func generateWiki(name, place, chartData string) {
	nameDir := dataDir + "/" + name
	wikiDir := nameDir + "/wiki"
	os.MkdirAll(wikiDir, os.ModePerm)

	// Topics to generate
	topics := []struct {
		filename string
		title    string
		prompt   string
	}{
		{
			filename: "ascendant_personality.md",
			title:    "Ascendant & Personality Profile",
			prompt:   "Write a detailed, professional analysis of my Ascendant (Lagna) and its Nakshatra based on my chart. Include how the Ascendant sign and Nakshatra shape my core personality, physical traits, and general life outlook. Keep it narrative and insightful, without raw degrees/coordinates tables.",
		},
		{
			filename: "planets_houses.md",
			title:    "Planetary Placements & Houses",
			prompt:   "Write a comprehensive breakdown of all the planetary placements in my chart. Detail what each planet represents and how it manifests in its specific zodiac sign and house. Highlight any conjunctions or key relationships between planets. Keep it narrative and insightful, without raw degrees/coordinates tables.",
		},
		{
			filename: "career_wealth.md",
			title:    "Career & Wealth Profile",
			prompt:   "Analyze my professional inclinations, career strengths, and financial growth patterns based on my chart. Focus on the houses and planets related to work, daily service, intellect, and assets. Highlight the best industries or roles and provide a long-term career outlook. Keep it narrative and insightful, without raw degrees/coordinates tables.",
		},
		{
			filename: "relationships_spirituality.md",
			title:    "Relationships & Spirituality",
			prompt:   "Analyze my relationship patterns, partnership potential, and spiritual or inner growth markers based on my chart. Focus on the 7th house, the Moon, and spiritual/outer planets. Describe what I seek in relationships, how I process emotions, and my spiritual potential. Keep it narrative and insightful, without raw degrees/coordinates tables.",
		},
	}

	for _, t := range topics {
		fmt.Printf("- Generating %s... ", t.title)
		resp, err := internal.AskSilent(t.prompt, chartData)
		if err != nil {
			fmt.Printf("failed: %v\n", err)
			continue
		}
		
		content := fmt.Sprintf("# %s: %s\n\n%s\n", t.title, name, resp)
		err = os.WriteFile(wikiDir+"/"+t.filename, []byte(content), 0644)
		if err != nil {
			fmt.Printf("failed writing: %v\n", err)
			continue
		}
		fmt.Println("done.")
	}

	// Generate index.md
	fmt.Print("- Generating index.md... ")
	indexContent := fmt.Sprintf(`# Astrological Wiki: %[1]s

Welcome to your personalized astrological wiki. This system compiles and analyzes the data of your birth chart to provide structured insights into your personality, career, relationships, and spiritual path.

---

## 🌌 Natal Chart Summary (At a Glance)
*   **Name:** %[1]s
*   **Birth Place:** %[2]s

---

## 📚 Wiki Sections

Click on the sections below to read detailed analyses:

1.  **[Ascendant & Personality Profile](wiki/ascendant_personality.md)**
    *   *Lagna, Nakshatras, and core constitution.*
2.  **[Planetary Placements & Houses](wiki/planets_houses.md)**
    *   *Every planet's placement and sign manifestation.*
3.  **[Career & Wealth Profile](wiki/career_wealth.md)**
    *   *Professional strengths and financial outlook.*
4.  **[Relationships & Spirituality](wiki/relationships_spirituality.md)**
    *   *Interpersonal patterns and inner growth.*

---

## 🗺️ Raw Chart Data
The mathematical computations of your planetary positions are stored in:
*   [chart.json](raw/chart.json)
`, name, place)

	err := os.WriteFile(nameDir+"/index.md", []byte(indexContent), 0644)
	if err != nil {
		fmt.Printf("failed writing index.md: %v\n", err)
	} else {
		fmt.Println("done.")
	}
}
