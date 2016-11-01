// Copyright 2015 opentsdb-goclient authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// Package main shows the sample of how to use github.com/bluebreezecf/opentsdbclient/client
// to communicate with the OpenTSDB with the pre-define rest apis.
// (http://opentsdb.net/docs/build/html/api_http/index.html#api-endpoints)
//
package main

import (
	"fmt"
	"time"

	"github.com/bluebreezecf/opentsdb-goclient/client"
	"github.com/bluebreezecf/opentsdb-goclient/config"
)

func main() {
	opentsdbCfg := config.OpenTSDBConfig{
		OpentsdbHost: "127.0.0.1:4242",
	}
	tsdbClient, err := client.NewClient(opentsdbCfg)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	//0. Ping
	if err = tsdbClient.Ping(); err != nil {
		fmt.Println(err.Error())
		return
	}

	//2.1 POST /api/query to query
	fmt.Println("Begin to test POST /api/query.")
	st2 := time.Now().Unix()
	st1 := st2 - 60
	queryParam := client.QueryParam{
		Start: st1,
		End:   st2,
	}
	name := []string{"container__memory__rss", "container__network__receive__errors__total", "container__network__receive__bytes__total"}
	tags := map[string]string{
		"kubernetes_io_hostname": "kube-node-1",
	}
	subqueries := make([]client.SubQuery, 0)
	for _, metric := range name {
		subQuery := client.SubQuery{
			Aggregator: "sum",
			Metric:     metric,
			Tags:       tags,
		}
		subqueries = append(subqueries, subQuery)
	}
	queryParam.Queries = subqueries
	fmt.Println(queryParam)
	if queryResp, err := tsdbClient.Query(queryParam); err != nil {
		fmt.Printf("Error occurs when querying: %v", err)
	} else {
		for _, item := range queryResp.QueryRespCnts {
			fmt.Printf("%s:\n%v\n", item.Metric, item)
		}
	}
	fmt.Println("Finish testing POST /api/query.")

	//2.2 POST /api/query/last
	fmt.Println("Begin to test POST /api/query/last.")
	time.Sleep(1 * time.Second)
	subqueriesLast := make([]client.SubQueryLast, 0)
	for _, metric := range name {
		subQueryLast := client.SubQueryLast{
			Metric: metric,
			Tags:   tags,
		}
		subqueriesLast = append(subqueriesLast, subQueryLast)
	}
	queryLastParam := client.QueryLastParam{
		Queries:      subqueriesLast,
		ResolveNames: true,
		BackScan:     24,
	}
	if queryLastResp, err := tsdbClient.QueryLast(queryLastParam); err != nil {
		fmt.Printf("Error occurs when querying last: %v", err)
	} else {
		fmt.Printf("%s", queryLastResp.String())
	}
	fmt.Println("Finish testing POST /api/query/last.")

	
}
