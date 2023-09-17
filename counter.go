// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"net/http"
	"sync"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var counters = map[string]*Counter{}

type Counter struct {
	mu    sync.Mutex
	id    string
	value int
}

func NewCounter(id string) *Counter {
	return &Counter{id: id, value: 0}
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += 1
}

func (c *Counter) GetValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func (c *Counter) GetID() string {
	return c.id
}

//

func listCountersHandler(c echo.Context) error {
	var counterViews []templ.Component
	keys := maps.Keys(counters)
	slices.Sort(keys)
	for _, counterID := range keys {
		counterViews = append(counterViews, counterView(counters[counterID]))
	}
	component := counterPage(counterViews)
	return renderComponent(c, http.StatusOK, component)
}

func incrementCounterHandler(c echo.Context) error {
	counter, ok := counters[c.Param("counterID")]
	if !ok {
		return c.HTML(http.StatusNotFound, "") // TODO How to deal with errors like this with HTMX
	}

	counter.Increment()
	component := counterValue(counter)
	return renderComponent(c, http.StatusOK, component)
}
