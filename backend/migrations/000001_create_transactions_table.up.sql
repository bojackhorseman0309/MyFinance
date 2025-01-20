CREATE TABLE IF NOT EXISTS WalletTransaction
(
    id               bigserial PRIMARY KEY,
    created_at       timestamp without time zone not null default now(),
    title            text                        NOT NULL,
    category         text                        NOT NULL,
    account          text                        NOT NULL,
    amount           decimal(10, 2)              NOT NULL,
    currency         text                        NOT NULL,
    transactionType  text                        NOT NULL,
    transferAmount   decimal(10, 2),
    transferCurrency text,
    toAccount        text,
    receiveAmount    decimal(10, 2),
    receiveCurrency  text,
    description      text,
    dueDate          timestamp without time zone not null default now(),
    walletId         text                        NOT NULL
);