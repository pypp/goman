package maps

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func removeEmptyLines(lines []string) []string {
	var filteredLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			filteredLines = append(filteredLines, line)
		}
	}
	return filteredLines
}

func LoadMap(mapFile string) (gameMap [][]int, mapHeight int, mapWidth int) {
	data, err := os.ReadFile(mapFile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lines := removeEmptyLines(strings.Split(string(data), "\n"))
	mapHeight = len(lines)
	mapWidth = len(strings.Fields(lines[0]))
	gameMap = make([][]int, mapHeight)

	for i, line := range lines {
		values := strings.Fields(line)
		gameMap[i] = make([]int, len(values))

		for j, val := range values {
			num, _ := strconv.ParseInt(val, 10, 64)
			gameMap[i][j] = int(num)
		}
	}

	return
}
