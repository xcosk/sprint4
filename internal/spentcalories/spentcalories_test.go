package spentcalories

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SpentCaloriesTestSuite struct {
	suite.Suite
}

func TestSpentCaloriesSuite(t *testing.T) {
	suite.Run(t, new(SpentCaloriesTestSuite))
}

func (suite *SpentCaloriesTestSuite) TestParseTraining() {
	tests := []struct {
		name         string
		input        string
		wantSteps    int
		wantDuration time.Duration
		wantErr      bool
	}{
		{
			name:         "корректный ввод с часами и минутами",
			input:        "3456,Ходьба,3h00m",
			wantSteps:    3456,
			wantDuration: 3 * time.Hour,
			wantErr:      false,
		},
		{
			name:         "корректный ввод с минутами",
			input:        "678,Бег,5m",
			wantSteps:    678,
			wantDuration: 5 * time.Minute,
			wantErr:      false,
		},
		{
			name:         "положительное число с плюсом",
			input:        "+12345,Ходьба,1h30m",
			wantSteps:    12345,
			wantDuration: 90 * time.Minute,
			wantErr:      false,
		},
		{
			name:         "продолжительность - только минуты",
			input:        "1000,Бег,30m",
			wantSteps:    1000,
			wantDuration: 30 * time.Minute,
			wantErr:      false,
		},
		{
			name:         "продолжительность - только часы",
			input:        "1000,Ходьба,2h",
			wantSteps:    1000,
			wantDuration: 2 * time.Hour,
			wantErr:      false,
		},
		{
			name:         "продолжительность - дробные часы",
			input:        "1000,Бег,1.5h",
			wantSteps:    1000,
			wantDuration: 90 * time.Minute,
			wantErr:      false,
		},
		{
			name:         "продолжительность - дробные минуты",
			input:        "1000,Ходьба,30.5m",
			wantSteps:    1000,
			wantDuration: 30*time.Minute + 30*time.Second,
			wantErr:      false,
		},
		{
			name:         "неверный формат - неправильное количество параметров",
			input:        "678,Ходьба",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверный формат - четыре параметра",
			input:        "678,Ходьба,1h30m,extra",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "пустой ввод",
			input:        "",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверные шаги - не числовое значение",
			input:        "abc,Ходьба,1h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверные шаги - ноль",
			input:        "0,Ходьба,1h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверные шаги - отрицательное значение",
			input:        "-100,Ходьба,1h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверные шаги - только знак минус",
			input:        "-,Ходьба,1h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверные шаги - только знак плюс",
			input:        "+,Ходьба,1h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверный формат продолжительности",
			input:        "678,Ходьба,invalid",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверная продолжительность - ноль",
			input:        "678,Бег,0h0m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверная продолжительность - отрицательное значение",
			input:        "678,Ходьба,-1h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверная продолжительность - отрицательные минуты",
			input:        "678,Бег,1h-30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверная продолжительность - неверная единица измерения",
			input:        "678,Ходьба,1.5d",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверная продолжительность - пробел между числом и единицей",
			input:        "678,Бег,1 h30m",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
		{
			name:         "неверная продолжительность - пропущена единица измерения",
			input:        "678,Ходьба,30",
			wantSteps:    0,
			wantDuration: 0,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			gotSteps, _, gotDuration, err := parseTraining(tt.input)

			if tt.wantErr {
				assert.Error(suite.T(), err)
				assert.Equal(suite.T(), 0, gotSteps)
				assert.Equal(suite.T(), time.Duration(0), gotDuration)
				return
			}

			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tt.wantSteps, gotSteps)
			assert.Equal(suite.T(), tt.wantDuration, gotDuration)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestDistance() {
	tests := []struct {
		name     string
		steps    int
		height   float64
		wantDist float64
	}{
		{
			name:     "нормальное количество шагов",
			steps:    1000,
			height:   1.75,
			wantDist: 0.7875,
		},
		{
			name:     "большое количество шагов",
			steps:    10000,
			height:   1.75,
			wantDist: 7.875,
		},
		{
			name:     "маленькое количество шагов",
			steps:    100,
			height:   1.75,
			wantDist: 0.07875,
		},
		{
			name:     "ноль шагов",
			steps:    0,
			height:   1.75,
			wantDist: 0,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got := distance(tt.steps, tt.height)
			assert.Equal(suite.T(), tt.wantDist, got)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestMeanSpeed() {
	tests := []struct {
		name      string
		steps     int
		height    float64
		duration  time.Duration
		wantSpeed float64
	}{
		{
			name:      "нормальная скорость - один час",
			steps:     6000,
			height:    1.75,
			duration:  1 * time.Hour,
			wantSpeed: 4.725,
		},
		{
			name:      "нормальная скорость - полчаса",
			steps:     3000,
			height:    1.75,
			duration:  30 * time.Minute,
			wantSpeed: 4.725,
		},
		{
			name:      "нормальная скорость - два часа",
			steps:     12000,
			height:    1.75,
			duration:  2 * time.Hour,
			wantSpeed: 4.725,
		},
		{
			name:      "маленькая скорость",
			steps:     1000,
			height:    1.75,
			duration:  2 * time.Hour,
			wantSpeed: 0.39375,
		},
		{
			name:      "большая скорость",
			steps:     20000,
			height:    1.75,
			duration:  1 * time.Hour,
			wantSpeed: 15.75,
		},
		{
			name:      "нулевая продолжительность",
			steps:     1000,
			height:    1.75,
			duration:  0,
			wantSpeed: 0,
		},
		{
			name:      "отрицательная продолжительность",
			steps:     1000,
			height:    1.75,
			duration:  -1 * time.Hour,
			wantSpeed: 0,
		},
		{
			name:      "ноль шагов",
			steps:     0,
			height:    1.75,
			duration:  1 * time.Hour,
			wantSpeed: 0,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got := meanSpeed(tt.steps, tt.height, tt.duration)
			assert.Equal(suite.T(), tt.wantSpeed, got)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestRunningSpentCalories() {
	tests := []struct {
		name     string
		steps    int
		weight   float64
		height   float64
		duration time.Duration
		wantCal  float64
		wantErr  bool
	}{
		{
			name:     "нормальная нагрузка - один час",
			steps:    6000,
			weight:   75.0,
			height:   1.75,
			duration: 1 * time.Hour,
			wantCal:  354.375,
			wantErr:  false,
		},
		{
			name:     "нормальная нагрузка - полчаса",
			steps:    3000,
			weight:   75.0,
			height:   1.75,
			duration: 30 * time.Minute,
			wantCal:  177.1875,
			wantErr:  false,
		},
		{
			name:     "высокая скорость",
			steps:    20000,
			weight:   75.0,
			height:   1.75,
			duration: 1 * time.Hour,
			wantCal:  1181.25,
			wantErr:  false,
		},
		{
			name:     "низкая скорость",
			steps:    1000,
			weight:   75.0,
			height:   1.75,
			duration: 2 * time.Hour,
			wantCal:  59.0625,
			wantErr:  false,
		},
		{
			name:     "другой вес",
			steps:    6000,
			weight:   60.0,
			height:   1.75,
			duration: 1 * time.Hour,
			wantCal:  283.5,
			wantErr:  false,
		},
		{
			name:     "нулевая продолжительность",
			steps:    1000,
			weight:   75.0,
			height:   1.75,
			duration: 0,
			wantCal:  0,
			wantErr:  true,
		},
		{
			name:     "отрицательная продолжительность",
			steps:    1000,
			weight:   75.0,
			height:   1.75,
			duration: -1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
		{
			name:     "ноль шагов",
			steps:    0,
			weight:   75.0,
			height:   1.75,
			duration: 1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
		{
			name:     "отрицательные шаги",
			steps:    -1000,
			weight:   75.0,
			height:   1.75,
			duration: 1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
		{
			name:     "нулевой вес",
			steps:    1000,
			weight:   0,
			height:   1.75,
			duration: 1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
		{
			name:     "отрицательный вес",
			steps:    1000,
			weight:   -75.0,
			height:   1.75,
			duration: 1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			gotCal, gotErr := RunningSpentCalories(tt.steps, tt.weight, tt.height, tt.duration)

			if tt.wantErr {
				assert.Error(suite.T(), gotErr)
				assert.Equal(suite.T(), 0.0, gotCal)
				return
			}

			assert.NoError(suite.T(), gotErr)
			assert.InDelta(suite.T(), tt.wantCal, gotCal, 0.1)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestWalkingSpentCalories() {
	tests := []struct {
		name     string
		steps    int
		weight   float64
		height   float64
		duration time.Duration
		wantCal  float64
		wantErr  bool
	}{
		{
			name:     "нормальная нагрузка",
			steps:    6000,
			weight:   75.0,
			height:   1.75,
			duration: 1 * time.Hour,
			wantCal:  177.19,
			wantErr:  false,
		},
		{
			name:     "меньше шагов",
			steps:    3000,
			weight:   75.0,
			height:   1.75,
			duration: 30 * time.Minute,
			wantCal:  88.594,
			wantErr:  false,
		},
		{
			name:     "больше шагов",
			steps:    20000,
			weight:   75.0,
			height:   1.75,
			duration: 1 * time.Hour,
			wantCal:  590.62,
			wantErr:  false,
		},
		{
			name:     "другой вес",
			steps:    6000,
			weight:   60.0,
			height:   1.75,
			duration: 1 * time.Hour,
			wantCal:  141.75,
			wantErr:  false,
		},
		{
			name:     "другой рост",
			steps:    6000,
			weight:   75.0,
			height:   1.85,
			duration: 1 * time.Hour,
			wantCal:  187.313,
			wantErr:  false,
		},
		{
			name:     "нулевые шаги",
			steps:    0,
			weight:   75.0,
			height:   1.75,
			duration: 1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
		{
			name:     "отрицательные шаги",
			steps:    -1000,
			weight:   75.0,
			height:   1.75,
			duration: 1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
		{
			name:     "нулевой вес",
			steps:    6000,
			weight:   0,
			height:   1.75,
			duration: 1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
		{
			name:     "отрицательный вес",
			steps:    6000,
			weight:   -75.0,
			height:   1.75,
			duration: 1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
		{
			name:     "нулевой рост",
			steps:    6000,
			weight:   75.0,
			height:   0,
			duration: 1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
		{
			name:     "отрицательный рост",
			steps:    6000,
			weight:   75.0,
			height:   -1.75,
			duration: 1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			gotCal, gotErr := WalkingSpentCalories(tt.steps, tt.weight, tt.height, tt.duration)

			if tt.wantErr {
				assert.Error(suite.T(), gotErr)
				assert.Equal(suite.T(), 0.0, gotCal)
				return
			}

			assert.NoError(suite.T(), gotErr)
			assert.InDelta(suite.T(), tt.wantCal, gotCal, 0.1)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestTrainingInfo() {
	tests := []struct {
		name    string
		input   string
		weight  float64
		height  float64
		want    string
		wantErr bool
	}{
		{
			name:    "ходьба - нормальная нагрузка",
			input:   "6000,Ходьба,1h00m",
			weight:  75.0,
			height:  1.75,
			want:    "Тип тренировки: Ходьба\nДлительность: 1.00 ч.\nДистанция: 4.72 км.\nСкорость: 4.72 км/ч\nСожгли калорий: 177.19\n",
			wantErr: false,
		},
		{
			name:    "бег - нормальная нагрузка",
			input:   "6000,Бег,1h00m",
			weight:  75.0,
			height:  1.75,
			want:    "Тип тренировки: Бег\nДлительность: 1.00 ч.\nДистанция: 4.72 км.\nСкорость: 4.72 км/ч\nСожгли калорий: 354.38\n",
			wantErr: false,
		},
		{
			name:    "ходьба - высокая скорость",
			input:   "20000,Ходьба,1h00m",
			weight:  75.0,
			height:  1.75,
			want:    "Тип тренировки: Ходьба\nДлительность: 1.00 ч.\nДистанция: 15.75 км.\nСкорость: 15.75 км/ч\nСожгли калорий: 590.62\n",
			wantErr: false,
		},
		{
			name:    "бег - высокая скорость",
			input:   "20000,Бег,1h00m",
			weight:  75.0,
			height:  1.75,
			want:    "Тип тренировки: Бег\nДлительность: 1.00 ч.\nДистанция: 15.75 км.\nСкорость: 15.75 км/ч\nСожгли калорий: 1181.25\n",
			wantErr: false,
		},
		{
			name:    "ходьба - другой вес и рост",
			input:   "6000,Ходьба,1h00m",
			weight:  60.0,
			height:  1.85,
			want:    "Тип тренировки: Ходьба\nДлительность: 1.00 ч.\nДистанция: 5.00 км.\nСкорость: 5.00 км/ч\nСожгли калорий: 149.85\n",
			wantErr: false,
		},
		{
			name:    "бег - другой вес",
			input:   "6000,Бег,1h00m",
			weight:  60.0,
			height:  1.75,
			want:    "Тип тренировки: Бег\nДлительность: 1.00 ч.\nДистанция: 4.72 км.\nСкорость: 4.72 км/ч\nСожгли калорий: 283.50\n",
			wantErr: false,
		},
		{
			name:    "ходьба - полчаса",
			input:   "3000,Ходьба,30m",
			weight:  75.0,
			height:  1.75,
			want:    "Тип тренировки: Ходьба\nДлительность: 0.50 ч.\nДистанция: 2.36 км.\nСкорость: 4.72 км/ч\nСожгли калорий: 88.59\n",
			wantErr: false,
		},
		{
			name:    "бег - полчаса",
			input:   "3000,Бег,30m",
			weight:  75.0,
			height:  1.75,
			want:    "Тип тренировки: Бег\nДлительность: 0.50 ч.\nДистанция: 2.36 км.\nСкорость: 4.72 км/ч\nСожгли калорий: 177.19\n",
			wantErr: false,
		},
		{
			name:    "неизвестный тип тренировки",
			input:   "6000,Плавание,1h00m",
			weight:  75.0,
			height:  1.75,
			want:    "",
			wantErr: true,
		},
		{
			name:    "неизвестный тип тренировки - проверка текста ошибки",
			input:   "6000,Плавание,1h00m",
			weight:  75.0,
			height:  1.75,
			want:    "",
			wantErr: true,
		},
		{
			name:    "некорректный формат данных",
			input:   "6000,Ходьба",
			weight:  75.0,
			height:  1.75,
			want:    "",
			wantErr: true,
		},
		{
			name:    "некорректное количество шагов",
			input:   "0,Ходьба,1h00m",
			weight:  75.0,
			height:  1.75,
			want:    "",
			wantErr: true,
		},
		{
			name:    "некорректная продолжительность",
			input:   "6000,Ходьба,0h00m",
			weight:  75.0,
			height:  1.75,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got, err := TrainingInfo(tt.input, tt.weight, tt.height)

			if tt.wantErr {
				assert.Error(suite.T(), err)
				assert.Empty(suite.T(), got)
				if tt.name == "неизвестный тип тренировки - проверка текста ошибки" {
					assert.Contains(suite.T(), err.Error(), "неизвестный тип тренировки")
				}
				return
			}

			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tt.want, got)
		})
	}
}
