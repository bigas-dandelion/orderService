package generator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"l0/prod/entity"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func GenerateData() ([]byte, error) {
	gofakeit.Seed(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())

	orderUID := uuid.New().String()
	trackNumber := "WBILMTESTTRACK" + gofakeit.DigitN(6)

	price := int(gofakeit.Price(100, 1000))
	sale := rand.Intn(51)
	totalPrice := price - (price*sale)/100

	ridI := gofakeit.IntRange(1, 1000000)
	ridJ := gofakeit.IntRange(1, 100)

	item := entity.Item{
		ChrtID:      gofakeit.Int64(),
		TrackNumber: trackNumber,
		Price:       price,
		Rid:         fmt.Sprintf("rid_%d_%d", ridI, ridJ),
		Name:        gofakeit.ProductName(),
		Sale:        sale,
		Size:        strconv.Itoa(38 + rand.Intn(10)),
		TotalPrice:  totalPrice,
		NmID:        int64(1000000 + rand.Intn(9000000)),
		Brand:       gofakeit.Company(),
		Status:      202,
	}

	items := []entity.Item{item}

	amount := item.TotalPrice
	if amount == 0 {
		amount = 1
	}

	var phone, email string
	switch rand.Intn(3) {
	case 0:
		phone = gofakeit.Phone()
	case 1:
		email = gofakeit.Email()
	default:
		phone = gofakeit.Phone()
		email = gofakeit.Email()
	}

	// 5% битых заказов
	isBroken := rand.Intn(100) < 5
	if isBroken {
		switch rand.Intn(5) {
		case 0:
			orderUID = ""
		case 1:
			phone = ""
		case 2:
			email = ""
		case 3:
			amount = -100
		case 4:
			items = []entity.Item{}
		}
	}

	locales := []string{"en", "ru", "es", "fr", "de", "it", "pt", "ja", "ko", "zh"}
	locale := locales[rand.Intn(len(locales))]

	order := entity.Order{
		OrderUID:    orderUID,
		TrackNumber: trackNumber,
		Entry:       "WBIL",
		Delivery: entity.Delivery{
			Name:    gofakeit.Name(),
			Phone:   phone,
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Street(),
			Region:  gofakeit.State(),
			Email:   email,
		},
		Payment: entity.Payment{
			Transaction:  gofakeit.LetterN(15) + "test",
			RequestID:    gofakeit.DigitN(8),
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       amount,
			PaymentDT:    time.Now().Unix(),
			Bank:         "discount",
			DeliveryCost: 500,
			GoodsTotal:   amount,
			CustomFee:    0,
		},
		Items:             items,
		Locale:            locale,
		InternalSignature: "",
		CustomerID:        gofakeit.LetterN(15),
		DeliveryService:   "dhl",
		ShardKey:          strconv.Itoa(int(rand.Int31n(10))),
		SmID:              100 + ridI,
		DateCreated:       time.Now().UTC().Format(time.RFC3339),
		OofShard:          strconv.Itoa(int(rand.Int31n(10))),
	}

	data, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("не удалось преобразовать заказ в JSON: %w", err)
	}

	return data, nil
}
