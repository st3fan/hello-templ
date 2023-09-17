// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"bytes"
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func renderComponent(c echo.Context, status int, component templ.Component) error {
	var b bytes.Buffer
	if err := component.Render(context.Background(), &b); err != nil {
		return err
	}
	return c.HTMLBlob(status, b.Bytes())
}
