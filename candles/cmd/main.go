package main

import (
	"context"
	"de-reasearch-project-csu/candles/models"
	"de-reasearch-project-csu/candles/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var parsedShares = []models.Stock{
	{
		ISIN:         "RU0009029540",
		Ticker:       "SBER",
		MarketSector: "Банк",
	},
	{
		ISIN:         "RU0007661625",
		Ticker:       "GAZP",
		MarketSector: "Добыча ресурсов (газ, нефть)",
	},
	{
		ISIN:         "RU0007288411",
		Ticker:       "GMKN",
		MarketSector: "Добыча ресурсов (газ, нефть)",
	},
	{
		ISIN:         "RU0009024270",
		Ticker:       "LKOH",
		MarketSector: "Добыча ресурсов (газ, нефть)",
	},
	{
		ISIN:         "RU000A0J2Q06",
		Ticker:       "ROSN",
		MarketSector: "Добыча ресурсов (газ, нефть)",
	},
	{
		ISIN:         "RU0009062276",
		Ticker:       "AFLT",
		MarketSector: "Транспорт",
	},
	{
		ISIN:         "RU0008926258",
		Ticker:       "SNGS",
		MarketSector: "Транспорт",
	},
	{
		ISIN:         "RU000A0JKQU8",
		Ticker:       "MGNT",
		MarketSector: "Ритейл",
	},
	{
		ISIN:         "RU0009033597",
		Ticker:       "TATN",
		MarketSector: "Добыча ресурсов (газ, нефть)",
	},
	{
		ISIN:         "RU000A0HGZ25",
		Ticker:       "PLZL",
		MarketSector: "Добыча ресурсов (металлы)",
	},
	{
		ISIN:         "RU0009067001",
		Ticker:       "MAGN",
		MarketSector: "Добыча ресурсов (металлы)",
	},
	{
		ISIN:         "RU000A0JKQU1",
		Ticker:       "PHOR",
		MarketSector: "Сельское хозяйство",
	},
	{
		ISIN:         "RU000A0JP6E2",
		Ticker:       "HYDR",
		MarketSector: "Энерго сектор",
	},
	{
		ISIN:         "RU000A0D6866",
		Ticker:       "URKA",
		MarketSector: "Добыча ресурсов (металлы)",
	},
	{
		ISIN:         "RU000A0D9F79",
		Ticker:       "NLMK",
		MarketSector: "Добыча ресурсов (металлы)",
	},
	{
		ISIN:         "RU000A0D9G20",
		Ticker:       "MTL",
		MarketSector: "Добыча ресурсов (металлы)",
	},
	{
		ISIN:         "RU0009046510",
		Ticker:       "CHMF",
		MarketSector: "Добыча ресурсов (металлы)",
	},
	{
		ISIN:         "RU0009091938",
		Ticker:       "TRNFP",
		MarketSector: "Добыча ресурсов (газ, нефть)",
	},
	{
		ISIN:         "RU000A0JPNC0",
		Ticker:       "KMAZ",
		MarketSector: "Транспорт",
	},
	{
		ISIN:         "RU000A0D76M4",
		Ticker:       "RSTI",
		MarketSector: "Энерго сектор",
	},
	{
		ISIN:         "RU000A0DQ8G7",
		Ticker:       "AFKS",
		MarketSector: "Фармацевтика",
	},
	{
		ISIN:         "RU000A0JSQZ6",
		Ticker:       "LSRG",
		MarketSector: "Недвижимость / Строительство",
	},
	{
		ISIN:         "RU000A0JQRU9",
		Ticker:       "ENRU",
		MarketSector: "Энерго сектор",
	},
	{
		ISIN:         "RU000A0DQ2Z8",
		Ticker:       "TGKD",
		MarketSector: "Энерго сектор",
	},
	{
		ISIN:         "RU000A0JPJT0",
		Ticker:       "CPTC",
		MarketSector: "Транспорт",
	},
	{
		ISIN:         "RU000A0JKTY2",
		Ticker:       "AFKG",
		MarketSector: "Сельское хозяйство",
	},
	{
		ISIN:         "RU000A0JPPJ0",
		Ticker:       "RUAL",
		MarketSector: "Добыча ресурсов (металлы)",
	},
	{
		ISIN:         "RU000A0JPDH8",
		Ticker:       "OGKB",
		MarketSector: "Энерго сектор",
	},
	{
		ISIN:         "RU000A0JPGA4",
		Ticker:       "PIKK",
		MarketSector: "Недвижимость / Строительство",
	},
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file %s", err.Error())
	}

	config := investgo.Config{
		EndPoint:                      "sandbox-invest-public-api.tinkoff.ru:443",
		Token:                         os.Getenv("API_TOKEN"),
		AppName:                       "invest-api-go-sdk",
		AccountId:                     "",
		DisableResourceExhaustedRetry: false,
		DisableAllRetry:               false,
		MaxRetries:                    3,
	}

	connStr := os.Getenv("CONN_STR")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	logger := logrus.New()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		logger.Fatalf("error connect to DB %s", err.Error())
	}

	db := storage.NewStorage(pool)

	if err := db.Init(ctx); err != nil {
		logger.Errorf("error init sql %s", err.Error())
	}

	err = db.InsertCandleTypes(ctx, []models.CandleType{
		models.TypeMinute,
		models.TypeHour,
		models.TypeDay,
		models.TypeWeek,
	})
	if err != nil {
		logger.Errorf("error inserting candle types %s", err.Error())
	}

	err = db.InsertStocks(ctx, parsedShares)
	if err != nil {
		logger.Errorf("error insert Stocks %s", err.Error())
	}

	logger.Infof("Inserted all stocks")
	client, err := investgo.NewClient(ctx, config, logger)
	if err != nil {
		logger.Fatalf("client creating error %v", err.Error())
	}
	defer func() {
		if err := client.Stop(); err != nil {
			logger.Errorf("client stopping fatal error: %s", err.Error())
		}
	}()

	// создаем клиента для сервиса инструментов
	instrumentsService := client.NewInstrumentsServiceClient()

	marketDataService := client.NewMarketDataServiceClient()
	wg := &sync.WaitGroup{}
	for i := range parsedShares {
		if i%5 == 4 {
			time.Sleep(2 * time.Minute)
		}

		go parse(ctx, instrumentsService, db, marketDataService, parsedShares[i], wg, logger)
	}
}

func parse(ctx context.Context, isc *investgo.InstrumentsServiceClient, db *storage.Storage, marketDataServiceClient *investgo.MarketDataServiceClient, stock models.Stock, wg *sync.WaitGroup, logger *logrus.Logger) {
	wg.Add(1)
	defer wg.Done()

	instrument, err := isc.FindInstrument(stock.Ticker)
	if err != nil {
		logger.Errorf("could not found instrument %s for ticker %s", err.Error(), stock.Ticker)
		return
	}

	logger.Infof("GET CANDLES FOR TICKER %s", stock.Ticker)
	bestIndex := 0
	allCandles := make([][]models.Candle, len(instrument.Instruments))
	for j := range instrument.Instruments {
		candles, err := marketDataServiceClient.GetCandles(
			instrument.Instruments[j].Uid,
			investapi.CandleInterval_CANDLE_INTERVAL_DAY,
			time.Now().Add(-12*30*24*time.Hour),
			time.Now(),
			investapi.GetCandlesRequest_CANDLE_SOURCE_UNSPECIFIED,
			20000000,
		)
		if err != nil {
			logger.Errorf("error get all historical candles for TICKER: %s, %s", stock.Ticker, err.Error())
			return
		}

		logger.Infof("candles get successful for Ticker: %s, count: %d ", stock.Ticker, len(candles.Candles))
		parsedCandles := make([]models.Candle, len(candles.Candles))
		for i, candle := range candles.Candles {
			parsedCandles[i] = models.Candle{
				StockISIN: stock.Ticker,
				Open:      parseMoney(candle.Open.String()),
				Close:     parseMoney(candle.Close.String()),
				Time:      candle.Time.AsTime(),
				Type:      models.TypeHour,
			}
		}

		allCandles = append(allCandles, parsedCandles)
	}

	for x := range allCandles {
		if len(allCandles[x]) > bestIndex {
			bestIndex = x
		}
	}

	logger.Infof("parsed candles len: %d for Ticker %s", len(allCandles[bestIndex]), stock.Ticker)

	if err := db.InsertCandles(ctx, allCandles[bestIndex]); err != nil {
		logger.Errorf("error inserting candles %s", err.Error())
	}
}

func parseMoney(money string) string {
	money = strings.Replace(money, " ", "", -1)
	money = strings.Replace(money, "units:", "", -1)
	money = strings.Replace(money, "nano:", " ", -1)

	parts := strings.Split(money, " ")
	nano := "0"
	if len(parts) > 1 {
		nano = parts[1]
	}

	return parts[0] + "," + nano
}
