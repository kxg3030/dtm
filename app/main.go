/*
 * Copyright (c) 2021 yedf. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/yedf/dtm/bench"
	"github.com/yedf/dtm/common"
	"github.com/yedf/dtm/dtmcli"
	"github.com/yedf/dtm/dtmcli/dtmimp"
	"github.com/yedf/dtm/dtmsvr"
	"github.com/yedf/dtm/examples"

	_ "go.uber.org/automaxprocs"
)

var Version, Commit, Date string

var usage = `dtm is a lightweight distributed transaction manager.

usage:
    dtm [command]

Available commands:
    version           print dtm version
    dtmsvr            run dtm as a server
    dev               create all needed table and run dtm as a server
    bench             start bench server

    quick_start       run quick start example (dtm will create needed table)
    qs                same as quick_start
`

func main() {
	if len(os.Args) == 1 {
		fmt.Println(usage)
		for name := range examples.Samples {
			fmt.Printf("%4s%-18srun a sample includes %s\n", "", name, strings.ReplaceAll(name, "_", " "))
		}
		return
	}
	if os.Args[1] == "version" {
		fmt.Printf("version: %s commit: %s built at: %s\n", Version, Commit, Date)
		return
	}
	dtmimp.Logf("starting dtm....")
	common.MustLoadConfig()
	dtmcli.SetCurrentDBType(common.DtmConfig.DB["driver"])
	if os.Args[1] != "dtmsvr" {
		common.WaitDBUp()
		dtmsvr.PopulateDB(true)
		examples.PopulateDB(true)
	}
	dtmsvr.StartSvr()
	go dtmsvr.CronExpiredTrans(-1)

	switch os.Args[1] {
	case "bench":
		bench.StartSvr()
	case "dev", "dtmsvr":
	default:

	}
	select {}
}
