package core

import "time"

type listStyle int

var (
	searchVal = ""

	//done               = make(chan struct{})
	tickerTimeInterval = time.Millisecond * 10

	//cursorPos       = 0
	resultsListMaxH = 0

	listStyleShort    listStyle = 0
	listStyleLong     listStyle = 1
	selectedListStyle           = listStyleShort
)
