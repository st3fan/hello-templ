// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//

type Counter struct {
	id    string
	value int
}

func NewCounter(id string) *Counter {
	return &Counter{id: id, value: 0}
}

func (c *Counter) Increment() {
	c.value += 1
}

func (c *Counter) GetValue() int {
	return c.value
}

func (c *Counter) GetID() string {
	return c.id
}

// Application state .. maybe Echo has a better way to handle this?

var counters = map[string]*Counter{}

//

func renderComponent(c echo.Context, status int, component templ.Component) error {
	var b bytes.Buffer
	if err := component.Render(context.Background(), &b); err != nil {
		return err
	}
	return c.HTMLBlob(status, b.Bytes())
}

func helloHandler(c echo.Context) error {
	component := hello("Stefan")
	return renderComponent(c, http.StatusOK, component)
}

// Turn into CounterPageHandler - how to do this in templ
func pageHandler(c echo.Context) error {
	var counterViews []templ.Component
	for counterID := range counters {
		counterViews = append(counterViews, counterView(counters[counterID]))
	}
	component := counterPage(counterViews) // TODO pass a list of counters here
	return renderComponent(c, http.StatusOK, component)
}

// TODO Can this be nicer?
func incrementCounterHandler(c echo.Context) error {
	counter, ok := counters[c.QueryParams().Get("counterID")]
	if !ok {
		return c.HTML(http.StatusNotFound, "") // How to deal errors like this with HTMX
	}

	counter.Increment()
	component := counterValue(counter)
	return renderComponent(c, http.StatusOK, component)
}

func main() {
	for i := 1; i <= 3; i++ {
		counter := NewCounter(fmt.Sprintf("Counter-%d", i))
		counters[counter.GetID()] = counter
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", pageHandler)
	e.GET("/hello", helloHandler)
	e.GET("/increment", incrementCounterHandler) // /counters and /counters/xxx/increment
	e.Logger.Fatal(e.Start(":8080"))
}
