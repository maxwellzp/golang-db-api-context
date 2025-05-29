use golang_application;

CREATE TABLE IF NOT EXISTS exchange_rates
(
    id                 INT AUTO_INCREMENT PRIMARY KEY,
    currency_code      VARCHAR(3)     NOT NULL,
    base_currency_code VARCHAR(3)     NOT NULL,
    rate               DECIMAL(15, 6) NOT NULL,
    date_updated       DATETIME       NOT NULL
);

