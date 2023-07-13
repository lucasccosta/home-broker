package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Order               []*Order
	Transactions        []*Transaction
	OrdersChannel       chan *Order
	OrdersChannelOutput chan *Order
	WaitGroup           *sync.WaitGroup
}

func NewBook(orderChannel chan *Order, orderChannelOutput chan *Order, waitGroup *sync.WaitGroup) *Book {
	return &Book{
		Order:               []*Order{},
		Transactions:        []*Transaction{},
		OrdersChannel:       orderChannel,
		OrdersChannelOutput: orderChannelOutput,
		WaitGroup:           waitGroup,
	}
}

func (book *Book) Trade() {
	buyOrders := NewOrderQueue()
	sellOrders := NewOrderQueue()

	heap.Init(buyOrders)
	heap.Init(sellOrders)

	for order := range book.OrdersChannel {
		if order.OrderType == "BUY" {
			buyOrders.Push(order)
			if sellOrders.Len() > 0 && sellOrders.Orders[0].Price <= order.Price {
				sellOrder := sellOrders.Pop().(*Order)
				if sellOrder.PendingShares > 0 {
					transaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price)
					book.AddTransaction(transaction, book.WaitGroup)
					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)
					book.OrdersChannel <- sellOrder
					book.OrdersChannelOutput <- order

					if sellOrder.PendingShares > 0 {
						sellOrders.Push(sellOrder)
					}
				}
			}
		} else if order.OrderType == "SELL" {
			sellOrders.Push(order)
			if buyOrders.Len() > 0 && buyOrders.Orders[0].Price >= order.Price {
				buyOrder := buyOrders.Pop().(*Order)
				if buyOrder.PendingShares > 0 {
					transaction := NewTransaction(order, buyOrder, order.Shares, buyOrder.Price)
					book.AddTransaction(transaction, book.WaitGroup)
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					book.OrdersChannel <- buyOrder
					book.OrdersChannelOutput <- order
					if buyOrder.PendingShares > 0 {
						buyOrders.Push(buyOrder)
					}
				}
			}
		}
	}
}

func (book *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	minShares := sellingShares
	if buyingShares < minShares {
		minShares = buyingShares
	}

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.ID, -minShares)
	transaction.AddBuyOrderPendingShares(-minShares)
	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.ID, minShares)
	transaction.AddSellOrderPendingShares(-minShares)

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)

	transaction.CloseBuyOrder()
	transaction.CloseSellOrder()

	book.Transactions = append(book.Transactions, transaction)
}
