package main

import (
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/valyala/fasthttp"
)

var redisdb *redis.Client

func fastHTTPRawHandler(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) == "GET" {
		//

		switch string(ctx.Path()) {

		case "/w/transaction":
			status, description, statusCode, transactionForDb, transactionTime := verifyTransaction(ctx)
			if status == false {
				//

				makeResponse(status, description, statusCode, ctx)
				return

			} else {
				//

				errRedis := redisdb.ZAdd("RAW TRANSACTIONS", redis.Z{
					Score:  float64(transactionTime),
					Member: transactionForDb,
				})

				if isRedisError(errRedis) {
					//

					makeResponse(false, "Internal Server Error", fasthttp.StatusInternalServerError, ctx)
					return
				}

				makeResponse(true, "OK", fasthttp.StatusOK, ctx)
			}

		case "/w/getBalance":
			status, description, statusCode, sender := verifyBalanceRequest(ctx)
			if status == false {
				//

				makeResponse(status, description, statusCode, ctx)
				return

			} else {
				//

				zScore := redisdb.ZScore("BALANCE", sender)
				if isRedisError(zScore) {
					//

					makeResponse(false, "Internal Server Error", fasthttp.StatusInternalServerError, ctx)
					return
				}

				makeBalanceResponse(true, strconv.FormatFloat(zScore.Val(), 'f', 8, 64), fasthttp.StatusOK, ctx)
			}

		case "/w/tranStatus":
			status, description, statusCode, key := verifyTranStatusRequest(ctx)
			if status == false {
				//

				makeResponse(status, description, statusCode, ctx)
				return

			} else {
				//

				zScore := redisdb.ZScore("COMPLETE TRANSACTIONS", key)
				if isRedisError(zScore) {
					//

					makeResponse(true, "TRANSACTION NOT FOUND", fasthttp.StatusOK, ctx)
					return
				}

				makeResponse(true, "TRANSACTION OK", fasthttp.StatusOK, ctx)
			}

		case "/w/getStat":
			var intCmd *redis.IntCmd
			var floatCmd *redis.FloatCmd

			var blockHeight int64
			var tranCount int64
			var ltpb int64
			var ltps int64
			var upd int64
			var tpd float64

			intCmd = redisdb.ZCard("VNCCHAIN")
			if isRedisError(intCmd) {
				//

				makeResponse(false, "Internal Server Error", fasthttp.StatusInternalServerError, ctx)
				return
			}

			blockHeight = intCmd.Val()

			intCmd = redisdb.ZCard("COMPLETE TRANSACTIONS")
			if isRedisError(intCmd) {
				//

				makeResponse(false, "Internal Server Error", fasthttp.StatusInternalServerError, ctx)
				return
			}

			tranCount = intCmd.Val()
			if blockHeight == 0 {
				//

				ltpb = 0

			} else {
				//

				ltpb = int64(math.Round(float64(tranCount) / float64(blockHeight)))
			}

			intCmd = redisdb.ZCount("COMPLETE TRANSACTIONS", strconv.FormatInt(blockHeight-11, 10), strconv.FormatInt(blockHeight-1, 10))
			if isRedisError(intCmd) {
				//

				makeResponse(false, "Internal Server Error", fasthttp.StatusInternalServerError, ctx)
				return
			}

			ltps = intCmd.Val()
			ltps = int64(math.Round(float64(ltps) / 30.0))
			floatCmd = redisdb.ZScore("MONEY MOVE", time.Now().Format("2006-01-02"))
			if isRedisError(floatCmd) {
				//

				makeResponse(false, "Internal Server Error", fasthttp.StatusInternalServerError, ctx)
				return
			}

			tpd = floatCmd.Val()
			upd = int64(math.Round(math.Sqrt(float64(tranCount))))

			makeStatisticsResponse(true, tranCount, ltps, blockHeight, ltpb, tpd, upd, fasthttp.StatusOK, ctx)

			return

		default:
			//

			ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		}

		return
	}

	ctx.Error("Unsupported method", fasthttp.StatusMethodNotAllowed)
}

func main() {
	//

	var redisDBNumInt int = 1
	redisHost, ok := os.LookupEnv("REDIS_PORT_6379_TCP_ADDR")
	if !ok {
		//

		redisHost = "0.0.0.0"
	}

	redisPort, ok := os.LookupEnv("REDIS_PORT_6379_TCP_PORT")
	if !ok {
		//

		redisPort = "6379"
	}

	redisDBNum, ok := os.LookupEnv("REDIS_PORT_6379_DB_NUM")
	if !ok {
		//

		redisDBNumInt = 1

	} else {
		//

		if redisDBNumInt64, err := strconv.ParseInt(redisDBNum, 10, 64); err != nil {
			//

			redisDBNumInt = int(redisDBNumInt64)

		} else {
			//

			redisDBNumInt = 1
		}
	}

	redisdb = redis.NewClient(&redis.Options{
		Addr:         net.JoinHostPort(redisHost, redisPort),
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		DB:           redisDBNumInt,
	})

	statusCmd := redisdb.Ping()
	if isRedisError(statusCmd) {
		//

		log.Fatalln("No connection to REDIS. ", statusCmd.Err())
		return
	}

	defer redisdb.Close()

	server := &fasthttp.Server{
		Handler:          fastHTTPRawHandler,
		DisableKeepalive: true,
		GetOnly:          true,
	}

	log.Fatalln(server.ListenAndServe(":5000"))
}
