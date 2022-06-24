package miscellaneous

import (
	"encoding/json"
	CONSTANT "hecruit-backend/constant"
	DB "hecruit-backend/database"
	UTIL "hecruit-backend/util"
	"net/http"
	"strings"
)

func PaddleCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response = make(map[string]interface{})

	customDetail := map[string]string{}
	err := json.Unmarshal([]byte(strings.ReplaceAll(r.FormValue("passthrough"), "+", "")), &customDetail)
	if err != nil {
		UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
		return
	}

	switch r.FormValue("alert_name") {
	case "subscription_created":
		_, err = DB.UpdateSQL(CONSTANT.CompaniesTable, map[string]string{
			"id": customDetail["company_id"],
		}, map[string]string{
			"subscription_id": r.FormValue("subscription_id"),
			"update_url":      r.FormValue("update_url"),
			"cancel_url":      r.FormValue("cancel_url"),
			"plan":            "paid",
			"plan_amount":     "99",
			"payment_status":  r.FormValue("status"),
		})
		if err != nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	case "subscription_payment_succeeded":
		_, _, err = DB.InsertWithUniqueID(CONSTANT.InvoicesTable, map[string]string{
			"order_id":       r.FormValue("order_id"),
			"sale_gross":     r.FormValue("sale_gross"),
			"title":          "Hecruit PRO",
			"description":    "Monthly Billing",
			"receipt_url":    r.FormValue("receipt_url"),
			"company_id":     customDetail["company_id"],
			"payment_status": r.FormValue("status"),
		}, "id")
		if err != nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
		_, err = DB.UpdateSQL(CONSTANT.CompaniesTable, map[string]string{
			"id": customDetail["company_id"],
		}, map[string]string{
			"next_bill_date": r.FormValue("next_bill_date"),
		})
		if err != nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}

		if strings.EqualFold(r.FormValue("initial_payment"), "1") {
			// send first subscription email
		} else {
			// send next month mail
		}
	case "subscription_updated":
		_, err = DB.UpdateSQL(CONSTANT.CompaniesTable, map[string]string{
			"id": customDetail["company_id"],
		}, map[string]string{
			"payment_status": r.FormValue("status"),
		})
		if err != nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	case "subscription_payment_failed":
		_, err = DB.UpdateSQL(CONSTANT.CompaniesTable, map[string]string{
			"id": customDetail["company_id"],
		}, map[string]string{
			"update_url": r.FormValue("update_url"),
			"cancel_url": r.FormValue("cancel_url"),
		})
		if err != nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	case "subscription_cancelled":
		_, err = DB.UpdateSQL(CONSTANT.CompaniesTable, map[string]string{
			"id": customDetail["company_id"],
		}, map[string]string{
			"plan":           "free",
			"plan_amount":    "0",
			"next_bill_date": r.FormValue("cancellation_effective_date"),
		})
		if err != nil {
			UTIL.SetReponse(w, CONSTANT.StatusCodeBadRequest, err.Error(), CONSTANT.ShowDialog, response)
			return
		}
	}

	UTIL.SetReponse(w, CONSTANT.StatusCodeOk, "", CONSTANT.ShowDialog, response)
}
