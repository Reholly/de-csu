CREATE TABLE IF NOT EXISTS candle_type
(
    type INT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS stock
(
    isin VARCHAR(255) PRIMARY KEY,
    ticker VARCHAR(255) UNIQUE,
    sector VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS candle
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stock_ticker VARCHAR(255) REFERENCES stock(ticker) NOT NULL,
    open VARCHAR(255) NOT NULL,
    close VARCHAR(255) NOT NULL,
    time TIMESTAMP NOT NULL,
    type_id INT REFERENCES candle_type(type)
);
