// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Create some test counters ...
	for i := 1; i <= 3; i++ {
		counter := NewCounter(fmt.Sprintf("Counter-%d", i))
		counters[counter.GetID()] = counter
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/hello", helloHandler)

	e.GET("/counters", listCountersHandler)
	e.POST("/counters/:counterID/increment", incrementCounterHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
