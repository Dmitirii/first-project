package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	// Настройка режима и роутера
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// Логирование запросов
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s\n",
			param.TimeStamp.Format("15:04:05"),
			param.Method,
			param.Path,
		)
	}))

	// Обработчик GET /predict
	r.GET("/predict", func(c *gin.Context) {
		data := map[string]interface{}{
			"uid":    c.Query("uid"),
			"age":    parseInt(c.Query("age")),
			"gender": parseInt(c.Query("gender")),
			"rdw":    parseFloat(c.Query("rdw")),
			"wbc":    parseFloat(c.Query("wbc")),
			"rbc":    parseFloat(c.Query("rbc")),
			"hgb":    parseFloat(c.Query("hgb")),
			"hct":    parseFloat(c.Query("hct")),
			"mcv":    parseFloat(c.Query("mcv")),
			"mch":    parseFloat(c.Query("mch")),
			"mchc":   parseFloat(c.Query("mchc")),
			"plt":    parseFloat(c.Query("plt")),
			"neu":    parseFloat(c.Query("neu")),
			"eos":    parseFloat(c.Query("eos")),
			"bas":    parseFloat(c.Query("bas")),
			"lym":    parseFloat(c.Query("lym")),
			"mon":    parseFloat(c.Query("mon")),
			"soe":    parseFloat(c.Query("soe")),
			"chol":   parseFloat(c.Query("chol")),
			"glu":    parseFloat(c.Query("glu")),
		}

		response, err := sendToLabhub(data)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error":   err.Error(),
				"details": "Ошибка при обращении к целевому API",
			})
			return
		}

		c.JSON(http.StatusOK, response)
	})

	// Запуск сервера
	fmt.Println("Сервер запущен на http://localhost:8080")
	r.Run(":8080")
}

// Вспомогательные функции
func parseInt(s string) int {
	if s == "" {
		return 0
	}
	i, _ := strconv.Atoi(s)
	return i
}

func parseFloat(s string) float64 {
	if s == "" {
		return 0.0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func sendToLabhub(data map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("ошибка JSON: %v", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://apiml.labhub.online/api/v1/predict/hba1c",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %v", err)
	}

	// Заголовки как в ТЗ
	req.Header.Set("Authorization", "Bearer 0l62<EJi/zJx]a?")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка соединения: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	return result, nil
}
