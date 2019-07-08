package main

import (
	"encoding/hex"
	"encoding/json"
	"math"
	secp "secp256k1-go"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

func verifyTransaction(ctx *fasthttp.RequestCtx) (status bool,
	description string,
	statusCode int,
	transactionForDB []byte,
	transactionTime int64) {
	//

	ttype := string(ctx.FormValue("TT"))
	sender := string(ctx.FormValue("SENDER"))
	receiver := string(ctx.FormValue("RECEIVER"))
	tst := string(ctx.FormValue("TST"))

	if len(ttype) == 0 {
		//

		return false, "CAN NOT FIND ATTRIBUTE - TT", fasthttp.StatusBadRequest, transactionForDB, transactionTime
	}

	if len(sender) != 66 {
		//

		return false, "WRONG ATTRIBUTE - SENDER", fasthttp.StatusBadRequest, transactionForDB, transactionTime
	}

	if len(receiver) != 66 {
		//

		return false, "WRONG ATTRIBUTE - RECEIVER", fasthttp.StatusBadRequest, transactionForDB, transactionTime
	}

	if len(tst) != 10 {
		//

		return false, "WRONG ATTRIBUTE - TST", fasthttp.StatusBadRequest, transactionForDB, transactionTime
	}

	transactionTime, err := strconv.ParseInt(tst, 10, 64)
	if err != nil {
		//

		return false, "WRONG ATTRIBUTE - TST", fasthttp.StatusBadRequest, transactionForDB, transactionTime
	}

	timestamp := time.Unix(transactionTime, 0)
	if int64(math.Abs(float64(time.Since(timestamp)/time.Second))) > 10 {
		//

		return false, "WRONG ATTRIBUTE - TST", fasthttp.StatusBadRequest, transactionForDB, transactionTime
	}

	switch ttype {

	case "ST":
		ttoken := string(ctx.FormValue("TTOKEN"))
		ctoken := string(ctx.FormValue("CTOKEN"))
		sign := string(ctx.FormValue("SIGNATURE"))

		if len(ttoken) == 0 {
			//

			return false, "CAN NOT FIND ATTRIBUTE - TTOKEN", fasthttp.StatusBadRequest, transactionForDB, transactionTime
		}

		if len(ctoken) == 0 {
			//

			return false, "CAN NOT FIND ATTRIBUTE - CTOKEN", fasthttp.StatusBadRequest, transactionForDB, transactionTime

		} else {
			//

			_, err = strconv.ParseFloat(ctoken, 64)
			if err != nil {
				//

				return false, "WRONG ATTRIBUTE - CTOKEN", fasthttp.StatusBadRequest, transactionForDB, transactionTime
			}
		}

		if len(sign) != 130 {
			//

			return false, "WRONG SIGNATURE", fasthttp.StatusBadRequest, transactionForDB, transactionTime
		}

		transcationForVerify := SimpleTransactionForVerify{ttype, sender, receiver, ttoken, ctoken, tst}
		js, err := json.Marshal(transcationForVerify)
		if err != nil {
			//

			return false, "Internal Server Error", fasthttp.StatusInternalServerError, transactionForDB, transactionTime
		}

		decodedSignature, err := hex.DecodeString(sign)
		if err != nil {
			//

			return false, "WRONG SIGNATURE", fasthttp.StatusBadRequest, transactionForDB, transactionTime
		}

		publicKey, err := hex.DecodeString(sender)
		if err != nil {
			//

			return false, "WRONG ATTRIBUTE - SENDER", fasthttp.StatusBadRequest, transactionForDB, transactionTime
		}

		if secp.VerifySignature(js, decodedSignature, publicKey) != 1 {
			//

			return false, "CAN'T VERIFY SIGNATURE", fasthttp.StatusBadRequest, transactionForDB, transactionTime
		}

		transcation := SimpleTransaction{ttype, sender, receiver, ttoken, ctoken, tst, sign}
		transactionForDB, err = json.Marshal(transcation)
		if err != nil {
			//

			return false, "Internal Server Error", fasthttp.StatusInternalServerError, transactionForDB, transactionTime
		}

		return true, "OK", fasthttp.StatusOK, transactionForDB, transactionTime
	}

	return false, "UNKNOWN TRANSACTION TYPE", fasthttp.StatusBadRequest, transactionForDB, transactionTime
}

func verifyBalanceRequest(ctx *fasthttp.RequestCtx) (status bool, description string, statusCode int, sender string) {
	//

	ttoken := string(ctx.FormValue("TTOKEN"))
	sender = string(ctx.FormValue("SENDER"))

	if len(ttoken) == 0 {
		//

		return false, "CAN NOT FIND ATTRIBUTE - TTOKEN", fasthttp.StatusBadRequest, sender
	}

	if len(sender) != 66 {
		//

		return false, "WRONG ATTRIBUTE - SENDER", fasthttp.StatusBadRequest, sender
	}

	return true, "", fasthttp.StatusOK, sender
}

func verifyTranStatusRequest(ctx *fasthttp.RequestCtx) (status bool, description string, statusCode int, key string) {
	//

	sender := string(ctx.FormValue("SENDER"))
	receiver := string(ctx.FormValue("RECEIVER"))
	sign := string(ctx.FormValue("SIGNATURE"))
	key = sender + receiver + sign

	if len(receiver) != 66 {
		//

		return false, "WRONG ATTRIBUTE - RECEIVER", fasthttp.StatusBadRequest, key
	}

	if len(sender) != 66 {
		//

		return false, "WRONG ATTRIBUTE - SENDER", fasthttp.StatusBadRequest, key
	}

	if len(sign) != 130 {
		//

		return false, "WRONG SIGNATURE", fasthttp.StatusBadRequest, key
	}

	return true, "", fasthttp.StatusOK, key
}
