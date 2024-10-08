package storage

import (
	"context"
	"de-reasearch-project-csu/candles/models"
	_ "embed"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed init.sql
var initSql string

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: pool,
	}
}

func (s *Storage) Init(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, initSql)
	return err
}

func (s *Storage) InsertCandles(ctx context.Context, candles []models.Candle) error {
	batch := &pgx.Batch{}
	for i := range candles {
		batch.Queue("INSERT INTO candle(stock_ticker, open, close, time, type_id) VALUES($1, $2, $3, $4, $5)",
			candles[i].StockISIN,
			candles[i].Open,
			candles[i].Close,
			candles[i].Time,
			candles[i].Type,
		)
	}
	result := s.pool.SendBatch(ctx, batch)
	if err := result.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) InsertStocks(ctx context.Context, stocks []models.Stock) error {
	batch := &pgx.Batch{}
	for i := range stocks {
		batch.Queue("INSERT INTO stock(isin, ticker, sector) VALUES($1, $2, $3)",
			stocks[i].ISIN,
			stocks[i].Ticker,
			stocks[i].MarketSector,
		)
	}
	result := s.pool.SendBatch(ctx, batch)
	if err := result.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) InsertCandleTypes(ctx context.Context, types []models.CandleType) error {
	batch := &pgx.Batch{}
	for i := range types {
		batch.Queue("INSERT INTO candle_type(type) VALUES($1)", types[i])
	}

	result := s.pool.SendBatch(ctx, batch)
	if err := result.Close(); err != nil {
		return err
	}

	return nil
}
