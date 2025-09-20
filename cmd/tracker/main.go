package main

import (
	"fmt"
	"log"

	"github.com/Yandex-Practicum/tracker/internal/daysteps"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

func main() {
	weight := 84.6
	height := 1.87

	// дневная активность
	input := []string{
		"678,0h50m",
		"792,1h14m",
		"1078,1h30m",
		"7830,2h40m",
		",3456",
		"12:40:00, 3456",
		"something is wrong",
	}

	fmt.Println("Активность в течение дня")

	var (
		dayActionsInfo string
		dayActionsLog  []string
	)

	for _, v := range input {
		dayActionsInfo = daysteps.DayActionInfo(v, weight, height)
		dayActionsLog = append(dayActionsLog, dayActionsInfo)
	}

	for _, v := range dayActionsLog {
		fmt.Println(v)
	}

	// тренировки
	trainings := []string{
		"3456,Ходьба,3h00m",
		"something is wrong",
		"678,Бег,0h5m",
		"1078,Бег,0h10m",
		",3456 Ходьба",
		"7892,Ходьба,3h10m",
		"15392,Бег,0h45m",
	}

	var trainingLog []string

	for _, v := range trainings {
		trainingInfo, err := spentcalories.TrainingInfo(v, weight, height)
		if err != nil {
			log.Printf("не получилось получить информацию о тренировке: %v", err)
			continue
		}
		trainingLog = append(trainingLog, trainingInfo)
	}

	fmt.Println("Журнал тренировок")

	for _, v := range trainingLog {
		fmt.Println(v)
	}
}
