package entity

type OrderQueue struct {
	Orders []*Order
}

func (orderQueue *OrderQueue) Less(firstOrder, secondOrder int) bool {
	return orderQueue.Orders[firstOrder].Price < orderQueue.Orders[secondOrder].Price
}

func (orderQueue *OrderQueue) Swap(firstOrder, secondOrder int) {
	orderQueue.Orders[firstOrder], orderQueue.Orders[secondOrder] = orderQueue.Orders[secondOrder], orderQueue.Orders[firstOrder]
}

func (orderQueue *OrderQueue) Len() int {
	return len(orderQueue.Orders)
}

// 												"value any"
func (orderQueue *OrderQueue) Push(value interface{}) {
	// 																	Type casting - força o valor a ser do tipo Order
	orderQueue.Orders = append(orderQueue.Orders, value.(*Order))
}

func (orderQueue *OrderQueue) Pop() interface{} {
	oldValue := orderQueue.Orders
	ordersQuantity := len(oldValue)
	item := oldValue[ordersQuantity-1]
	orderQueue.Orders = oldValue[0 : ordersQuantity-1]

	return item
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{}
}
