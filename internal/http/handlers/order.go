package handlers

import "net/http"

type HTTPHandlers struct {
	orderInformer *orderInformer
}

func NewHTTPHandlers(orderInformer *orderInformer) *HTTPHandlers {
	return &HTTPHandlers{
		orderInformer: orderInformer,
	}
}

/*
endpoint: /orders
method: POST
description: Create a new order
info: JSON body with order details
*/
func (h *HTTPHandlers) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {

}

/*
endpoint: /orders/{id}
method: GET
description: Retrieve an order by ID

*/
func (h *HTTPHandlers) HandleGetOrder(w http.ResponseWriter, r *http.Request) {

}

/*
endpoint: /orders
method: GET
description: Retrieve all orders
*/
func (h *HTTPHandlers) HandleGetOrders(w http.ResponseWriter, r *http.Request) {
}

/*
endpoint: /orders/{id}
method: DELETE
description: Delete an order by ID

*/

func (h *HTTPHandlers) HandleDeleteOrder(w http.ResponseWriter, r *http.Request) {

}
