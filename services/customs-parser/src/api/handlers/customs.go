package handlers

import (
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DeclarationResponse struct {
	DeclarationNumber string `json:"declaration_number"`
	Declarant         string `json:"declarant"`
	Goods             string `json:"goods"`
	Sender            string `json:"sender"`
	Receiver          string `json:"receiver"`
	VMDNumber         string `json:"vmd_number"`
}

func GetDeclarationMock(c *gin.Context) {
	number := c.Param("number")

	if number == "" {
		slog.Warn("Attempted to get declaration without number")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Declaration number is required"})
		return
	}

	slog.Info("Fetching declaration mock data", "number", number)

	// Мокані компанії та товари для правдоподібної емуляції
	companies := []string{"ТОВ 'ЛОГІСТИКА Україна'", "Корпорація 'ТРАНС-ЕКСПОРТ'", "ПП 'Агро-торг Плюс'", "ФОП Шевченко О.В.", "ДП 'Укрпромпостач'"}
	goodsList := []string{"Автозапчастини для легкових авто", "Побутова техніка (холодильники)", "Меблі дерев'яні розібрані", "Продукти харчування (макаронні вироби)", "Одяг та текстильні вироби", "Сільськогосподарська техніка"}
	vmdPrefixes := []string{"UA100000", "UA200000", "UA300000"}

	// Seed random based on declaration number to get consistent results for same number
	seed := int64(0)
	for _, char := range number {
		seed += int64(char)
	}
	rng := rand.New(rand.NewSource(seed + int64(time.Now().Day()))) // Змінюється щодня

	mockData := DeclarationResponse{
		DeclarationNumber: number,
		Declarant:         companies[rng.Intn(len(companies))],
		Goods:             goodsList[rng.Intn(len(goodsList))],
		Sender:            companies[rng.Intn(len(companies))],
		Receiver:          companies[rng.Intn(len(companies))],
		VMDNumber:         vmdPrefixes[rng.Intn(len(vmdPrefixes))] + "/" + time.Now().Format("2006") + "/" + generateDigits(rng, 6),
	}

	// Штучна затримка для симуляції відповіді митниці (300-800ms)
	time.Sleep(time.Duration(300+rng.Intn(500)) * time.Millisecond)

	slog.Debug("Generated mock declaration", "number", number, "vmd", mockData.VMDNumber)
	c.JSON(http.StatusOK, mockData)
}

func generateDigits(rng *rand.Rand, length int) string {
	digits := ""
	for i := 0; i < length; i++ {
		digits += string(rune('0' + rng.Intn(10)))
	}
	return digits
}
