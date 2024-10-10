package migrations

import (
	"database/sql"
)

type V0003 struct{}

func init() {
	RegisterMigration(&V0003{})
}

func (v *V0003) Version() string {
	return "0003"
}

func (v *V0003) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TYPE transaction_type AS ENUM ('TOPUP', 'PAYMENT', 'TRANSFER');

	CREATE TABLE transactions (
		id UUID PRIMARY KEY, -- Unique identifier for the transaction
		user_id UUID REFERENCES users(id) ON DELETE CASCADE, -- User who made the transaction
		transaction_type transaction_type NOT NULL,
		
		-- Single amount field: positive for credit, negative for debit
		amount BIGINT NOT NULL, -- Use positive for credit (add) and negative for debit (subtract)

		remarks VARCHAR(255), -- Additional information for payments and transfers
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE INDEX transactions_user_id_idx ON transactions(user_id);

	CREATE VIEW view_transactions AS
	SELECT 
		t.id,
		t.user_id,
		t.transaction_type,
		t.amount,
		t.remarks,
		t.created_at,

		-- Calculate the balance_before: cumulative sum of amounts before the current transaction
		(
			SELECT COALESCE(SUM(t_prev.amount), 0)
			FROM transactions t_prev
			WHERE t_prev.user_id = t.user_id
			AND t_prev.created_at < t.created_at
		) AS balance_before,

		-- Calculate the balance_after: cumulative sum of amounts including the current transaction
		(
			SELECT COALESCE(SUM(t_prev.amount), 0)
			FROM transactions t_prev
			WHERE t_prev.user_id = t.user_id
			AND t_prev.created_at <= t.created_at
		) AS balance_after

	FROM transactions t
	ORDER BY t.created_at;
	`)
	if err != nil {
		return err
	}
	return nil
}

func (v *V0003) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`
	DROP VIEW IF EXISTS view_transactions;
	DROP TABLE IF EXISTS transactions;
	DROP TYPE IF EXISTS transaction_type;
	`)
	if err != nil {
		return err
	}
	return nil
}
