package domain

import (
	"errors"
	"time"
)

func (m *OrderMessage) ToDomain() (*Order, error) {
	if m.OrderUID == "" {
		return nil, errors.New("order_uid is required")

	}

	if m.Payment.Transaction == "" {
		return nil, errors.New("payment.transaction is required")

	}

	if len(m.Items) == 0 {
		return nil, errors.New("items must be non empty")

	}

	created, err := time.Parse(time.RFC3339, m.DateCreated)
	if err != nil {
		return nil, errors.New("invalid date_created format")
	}
	payDT := time.Unix(m.Payment.PaymentDT, 0).UTC()

	return &Order{
		OrderUID:    m.OrderUID,
		TrackNumber: m.TrackNumber,
		Entry:       m.Entry,
		Delivery:    m.Delivery,
		Payment: Payment{
			Transaction:  m.Payment.Transaction,
			RequestID:    m.Payment.RequestID,
			Currency:     m.Payment.Currency,
			Provider:     m.Payment.Provider,
			Amount:       m.Payment.Amount,
			PaymentDT:    payDT,
			Bank:         m.Payment.Bank,
			DeliveryCost: m.Payment.DeliveryCost,
			GoodsTotal:   m.Payment.GoodsTotal,
			CustomFee:    m.Payment.CustomFee,
		},
		Items:             m.Items,
		Locale:            m.Locale,
		InternalSignature: m.InternalSignature,
		CustomerID:        m.CustomerID,
		DeliveryService:   m.DeliveryService,
		ShardKey:          m.ShardKey,
		SmID:              m.SmID,
		DateCreated:       created.UTC(),
		OofShard:          m.OofShard,
	}, nil

}
