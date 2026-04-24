package handlers

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CompanyEDRResponse mirrors fields returned by the real Ukrainian EDR/USR registry.
// Field names match what the frontend expects for its Ukrainian-label mapping.
type CompanyEDRResponse struct {
	EDRPOU           string `json:"edrpou"`
	Name             string `json:"name"`              // Повна назва
	ShortName        string `json:"short_name"`        // Скорочена назва
	Boss             string `json:"boss"`              // Директор
	Address          string `json:"address"`           // Адреса
	RegistrationDate string `json:"registration_date"` // Дата реєстрації
	KVED             string `json:"kved"`              // КВЕД
	State            string `json:"state"`             // Статус
	Phone            string `json:"phone"`             // Телефон
	Email            string `json:"email"`             // Email
	IBAN             string `json:"iban"`              // IBAN
	Bank             string `json:"bank"`              // Банк
	MFO              string `json:"mfo"`               // МФО
}

// GetCompanyByEDRPOU returns mock company data deterministic per ЄДРПОУ code.
// TODO: replace mock with real EDR API call (https://usr.minjust.gov.ua or ring.org.ua).
func GetCompanyByEDRPOU(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "EDRPOU code is required"})
		return
	}

	slog.Info("Fetching company mock data", "edrpou", code)

	seed := codeToSeed(code)
	rng := rand.New(rand.NewSource(seed))

	surnames := []string{"Шевченко", "Коваленко", "Бондаренко", "Мельник", "Кравченко",
		"Олійник", "Ткаченко", "Лисенко", "Марченко", "Сидоренко"}
	firstNames := []string{"Олексій", "Іван", "Микола", "Сергій", "Петро",
		"Андрій", "Василь", "Михайло", "Юрій", "Дмитро"}
	patronymics := []string{"Олексійович", "Іванович", "Миколайович", "Сергійович",
		"Петрович", "Андрійович", "Васильович", "Михайлович"}

	companyTypes := []string{"ТОВ", "ПП", "ПАТ", "ТОВ", "ТОВ"}
	companyWords := [][]string{
		{"Агро", "Транс", "Лого", "Промо", "Техно", "Укр"},
		{"Сервіс", "Груп", "Трейд", "Постач", "Інвест", "Партнер"},
	}
	cities := []string{"Київ", "Харків", "Одеса", "Дніпро", "Львів",
		"Запоріжжя", "Кривий Ріг", "Миколаїв", "Херсон", "Полтава"}
	streets := []string{"вул. Центральна", "вул. Шевченка", "вул. Франка",
		"просп. Незалежності", "вул. Грушевського", "вул. Соборна",
		"просп. Перемоги", "вул. Хрещатик"}
	kvedCodes := []string{
		"46.12 Діяльність посередників у торгівлі паливом",
		"49.41 Вантажний автомобільний транспорт",
		"52.29 Інша допоміжна діяльність у сфері транспорту",
		"46.71 Оптова торгівля твердим, рідким та газоподібним паливом",
		"38.11 Збирання безпечних відходів",
		"01.11 Вирощування зернових культур",
		"46.61 Оптова торгівля с/г технікою",
		"33.12 Ремонт машин і устаткування",
	}
	banks := []struct{ Name, MFO string }{
		{"АТ КБ «ПриватБанк»", "305299"},
		{"АТ «Ощадбанк»", "300465"},
		{"АТ «Укрексімбанк»", "322313"},
		{"АТ «Райффайзен Банк»", "380805"},
		{"АТ «ПУМБ»", "334851"},
		{"АТ «Укргазбанк»", "320478"},
	}

	surname := surnames[rng.Intn(len(surnames))]
	firstName := firstNames[rng.Intn(len(firstNames))]
	patronymic := patronymics[rng.Intn(len(patronymics))]
	directorName := fmt.Sprintf("%s %s %s", surname, firstName, patronymic)

	compType := companyTypes[rng.Intn(len(companyTypes))]
	word1 := companyWords[0][rng.Intn(len(companyWords[0]))]
	word2 := companyWords[1][rng.Intn(len(companyWords[1]))]
	companyName := fmt.Sprintf("%s «%s%s»", compType, word1, word2)
	shortName := fmt.Sprintf("%s «%s%s»", compType, word1, word2)

	city := cities[rng.Intn(len(cities))]
	street := streets[rng.Intn(len(streets))]
	building := rng.Intn(120) + 1
	address := fmt.Sprintf("%s, %s, %d", city, street, building)

	// Generate a registration date 5-25 years ago
	yearsAgo := 5 + rng.Intn(20)
	month := 1 + rng.Intn(12)
	day := 1 + rng.Intn(28)
	regDate := time.Date(time.Now().Year()-yearsAgo, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	registrationDate := regDate.Format("02.01.2006")

	kved := kvedCodes[rng.Intn(len(kvedCodes))]
	phone := fmt.Sprintf("+380%02d%03d%02d%02d", rng.Intn(99)+10, rng.Intn(999), rng.Intn(99), rng.Intn(99))
	email := fmt.Sprintf("info@%s%s.ua", strings.ToLower(word1), strings.ToLower(word2))

	bank := banks[rng.Intn(len(banks))]
	iban := fmt.Sprintf("UA%02d%s%019d", rng.Intn(90)+10, bank.MFO, seed)

	resp := CompanyEDRResponse{
		EDRPOU:           code,
		Name:             companyName,
		ShortName:        shortName,
		Boss:             directorName,
		Address:          address,
		RegistrationDate: registrationDate,
		KVED:             kved,
		State:            "зареєстровано",
		Phone:            phone,
		Email:            email,
		IBAN:             iban,
		Bank:             bank.Name,
		MFO:              bank.MFO,
	}

	// Simulate network latency (100-400ms)
	time.Sleep(time.Duration(100+rng.Intn(300)) * time.Millisecond)

	c.JSON(http.StatusOK, resp)
}

func codeToSeed(code string) int64 {
	var s int64
	for _, ch := range code {
		s = s*31 + int64(ch)
	}
	return s
}
