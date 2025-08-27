package pg

import (
	"Ex-L0/internal/domain"
	"Ex-L0/internal/repository"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) UpsertOrderTx(ctx context.Context, o *domain.Order) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// orders
	_, err = tx.Exec(ctx, `
		INSERT INTO public.orders(
			order_uid, track_number, entry, locale, internal_signature, customer_id,
			delivery_service, shardkey, sm_id, date_created, oof_shard
		)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (order_uid) DO UPDATE SET
			track_number       = EXCLUDED.track_number,
			entry              = EXCLUDED.entry,
			locale             = EXCLUDED.locale,
			internal_signature = EXCLUDED.internal_signature,
			customer_id        = EXCLUDED.customer_id,
			delivery_service   = EXCLUDED.delivery_service,
			shardkey           = EXCLUDED.shardkey,
			sm_id              = EXCLUDED.sm_id,
			date_created       = EXCLUDED.date_created,
			oof_shard          = EXCLUDED.oof_shard
	`,
		o.OrderUID, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature, o.CustomerID,
		o.DeliveryService, o.ShardKey, o.SmID, o.DateCreated, o.OofShard,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO public.delivery(
			order_uid, name, phone, zip, city, address, region, email
		)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (order_uid) DO UPDATE SET
			name    = EXCLUDED.name,
			phone   = EXCLUDED.phone,
			zip     = EXCLUDED.zip,
			city    = EXCLUDED.city,
			address = EXCLUDED.address,
			region  = EXCLUDED.region,
			email   = EXCLUDED.email
	`,
		o.OrderUID, o.Delivery.Name, o.Delivery.Phone, o.Delivery.Zip, o.Delivery.City,
		o.Delivery.Address, o.Delivery.Region, o.Delivery.Email,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO public.payment(
			order_uid, "transaction", request_id, currency, provider, amount, payment_dt,
			bank, delivery_cost, goods_total, custom_fee
		)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (order_uid) DO UPDATE SET
			"transaction" = EXCLUDED."transaction",
			request_id    = EXCLUDED.request_id,
			currency      = EXCLUDED.currency,
			provider      = EXCLUDED.provider,
			amount        = EXCLUDED.amount,
			payment_dt    = EXCLUDED.payment_dt,
			bank          = EXCLUDED.bank,
			delivery_cost = EXCLUDED.delivery_cost,
			goods_total   = EXCLUDED.goods_total,
			custom_fee    = EXCLUDED.custom_fee
	`,
		o.OrderUID, o.Payment.Transaction, o.Payment.RequestID, o.Payment.Currency, o.Payment.Provider,
		o.Payment.Amount, o.Payment.PaymentDT, o.Payment.Bank, o.Payment.DeliveryCost, o.Payment.GoodsTotal, o.Payment.CustomFee,
	)
	if err != nil {
		return err
	}

	if _, err = tx.Exec(ctx, `DELETE FROM public.items WHERE order_uid=$1`, o.OrderUID); err != nil {
		return err
	}
	for _, it := range o.Items {
		_, err = tx.Exec(ctx, `
			INSERT INTO public.items(
				order_uid, chrt_id, track_number, price, rid, name, sale, size,
				total_price, nm_id, brand, status
			)
			VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		`,
			o.OrderUID, it.ChrtID, it.TrackNumber, it.Price, it.RID, it.Name, it.Sale, it.Size,
			it.TotalPrice, it.NmID, it.Brand, it.Status,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *Repo) GetByID(ctx context.Context, uid string) (*domain.Order, error) {
	var o domain.Order

	// order
	row := r.pool.QueryRow(ctx, `
		SELECT
			order_uid, track_number, entry, locale, internal_signature, customer_id,
			delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM public.orders
		WHERE order_uid=$1
	`, uid)
	if err := row.Scan(
		&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature, &o.CustomerID,
		&o.DeliveryService, &o.ShardKey, &o.SmID, &o.DateCreated, &o.OofShard,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	// delivery
	row = r.pool.QueryRow(ctx, `
		SELECT name, phone, zip, city, address, region, email
		FROM public.delivery
		WHERE order_uid=$1
	`, uid)
	if err := row.Scan(
		&o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip, &o.Delivery.City,
		&o.Delivery.Address, &o.Delivery.Region, &o.Delivery.Email,
	); err != nil {
		return nil, err
	}

	// payment
	row = r.pool.QueryRow(ctx, `
		SELECT "transaction", request_id, currency, provider, amount, payment_dt,
		       bank, delivery_cost, goods_total, custom_fee
		FROM public.payment
		WHERE order_uid=$1
	`, uid)
	if err := row.Scan(
		&o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency, &o.Payment.Provider,
		&o.Payment.Amount, &o.Payment.PaymentDT, &o.Payment.Bank, &o.Payment.DeliveryCost,
		&o.Payment.GoodsTotal, &o.Payment.CustomFee,
	); err != nil {
		return nil, err
	}

	// items
	rows, err := r.pool.Query(ctx, `
		SELECT chrt_id, track_number, price, rid, name, sale, size,
		       total_price, nm_id, brand, status
		FROM public.items
		WHERE order_uid=$1
		ORDER BY id
	`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var it domain.Item
		if err := rows.Scan(
			&it.ChrtID, &it.TrackNumber, &it.Price, &it.RID, &it.Name, &it.Sale, &it.Size,
			&it.TotalPrice, &it.NmID, &it.Brand, &it.Status,
		); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &o, nil
}

func (r *Repo) SelectRecentIDS(ctx context.Context, n int) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT order_uid
		FROM public.orders
		ORDER BY date_created DESC
		LIMIT $1
	`, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	uids := make([]string, 0, n)
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			return nil, err
		}
		uids = append(uids, u)
	}
	return uids, rows.Err()
}
