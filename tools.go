package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/valyala/fasthttp"
)

type Response struct {
	STATUS bool
	DESCR  string
}

type BalanceResponse struct {
	STATUS  bool
	BALANCE string
}

type SimpleTransaction struct {
	TT        string
	SENDER    string
	RECEIVER  string
	TTOKEN    string
	CTOKEN    string
	TST       string
	SIGNATURE string
}

type SimpleTransactionForVerify struct {
	TT       string
	SENDER   string
	RECEIVER string
	TTOKEN   string
	CTOKEN   string
	TST      string
}

type HelloTransaction struct {
	TT        string
	SENDER    string
	SIGNATURE string
}

type HelloTransactionForVerify struct {
	TT     string
	SENDER string
}

type StatisticsResponse struct {
	TCOUNT  string
	LTPS    string
	BHEIGHT string
	LTPB    string
	TPD     string
	VMAP    []string
	UPD     string
}

type redisError interface {
	Err() error
}

func isRedisError(err redisError) bool {
	//

	return err.Err() != nil
}

func makeResponse(status bool,
	description string,
	statusCode int,
	ctx *fasthttp.RequestCtx) {
	//

	response := Response{status, description}
	jsResponse, _ := json.Marshal(response)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)
	ctx.SetBody(jsResponse)

	log.Println(statusCode, description, string(jsResponse))
}

func makeBalanceResponse(status bool,
	balance string,
	statusCode int,
	ctx *fasthttp.RequestCtx) {
	//

	response := BalanceResponse{status, balance}
	jsResponse, _ := json.Marshal(response)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)
	ctx.SetBody(jsResponse)
}

func makeStatisticsResponse(status bool,
	tranCount int64,
	ltps int64,
	blockHeight int64,
	ltpb int64,
	tpd float64,
	upd int64,
	statusCode int,
	ctx *fasthttp.RequestCtx) {
	//

	response := StatisticsResponse{strconv.FormatInt(tranCount, 10),
		strconv.FormatInt(ltps, 10),
		strconv.FormatInt(blockHeight, 10),
		strconv.FormatInt(ltpb, 10),
		strconv.FormatFloat(tpd, 'f', 8, 64),
		[]string{"50.11", "8.68"},
		strconv.FormatInt(upd, 10)}

	jsResponse, _ := json.Marshal(response)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)
	ctx.SetBody(jsResponse)
}
