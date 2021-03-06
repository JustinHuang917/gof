// Copyright 2012 The Justin Huang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofcore

func init() {
	initApplication()
	initRouters()
	initHandlers()
	initInvoker()
	initModel()
	initSession()
}
