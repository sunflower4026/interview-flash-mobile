package migrations

import (
	"database/sql"
)

type V0004 struct{}

func init() {
	RegisterMigration(&V0004{})
}

func (v *V0004) Version() string {
	return "0004"
}

func (v *V0004) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
	DROP VIEW view_transactions;
	CREATE VIEW view_transactions AS
	SELECT 
		t.id,
		t.user_id,
		t.transaction_type,
		CASE
			WHEN t.amount >= 0 THEN 'CREDIT'
			WHEN t.amount < 0 THEN 'DEBIT'
		END AS account_type,
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

func (v *V0004) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`
	DROP VIEW view_transactions;
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
