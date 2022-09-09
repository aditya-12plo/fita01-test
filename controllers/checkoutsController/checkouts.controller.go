package checkoutsController

import (
	"encoding/json"
	"fita-test-01/enums"
	models "fita-test-01/models"
	checkoutsRepository "fita-test-01/repositories/checkoutsRepository"
	responseRepository "fita-test-01/repositories/responseRepository"
	uniquerandRepository "fita-test-01/repositories/uniquerandRepository"
	"fita-test-01/repositories/validatorRepository"
	basketsService "fita-test-01/services/basketsService"
	checkoutsService "fita-test-01/services/checkoutsService"
	"fita-test-01/services/productsService"
	"fita-test-01/services/promotionsService"
	"fita-test-01/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func Index(Context *gin.Context) {

	data, err := checkoutsService.ReturnAllCheckouts()
	if err != nil {
		arr := make(map[string]interface{})
		arr["message"] = err
		res := responseRepository.Result{Status: http.StatusInternalServerError, Datas: nil, Errors: arr}
		Context.JSON(http.StatusInternalServerError, res)
		return
	}
	res := responseRepository.Result{Status: 200, Datas: data, Errors: nil}
	Context.JSON(http.StatusOK, res)
	return
}

func CheckoutData(Context *gin.Context) {

	const N = 10000
	rng := uniquerandRepository.NewUniqueRand(2 * N)
	buyerId := rng.Int()
	var objDatas []*checkoutsRepository.DataCheckout
	var objCheckouts []*checkoutsRepository.DataCheckoutRepos
	var datasCheckouts []*models.Checkouts
	var objValidate validatorRepository.ValidatorResult

	if err := Context.ShouldBindJSON(&objDatas); err != nil {
		arr := make(map[string]interface{})
		arr["message"] = "must be use raw with json data"
		res := responseRepository.Result{Status: http.StatusInternalServerError, Datas: nil, Errors: arr}
		Context.JSON(http.StatusInternalServerError, res)
		return
	}

	v := validator.New()
	count := 1
	for _, objData := range objDatas {

		err := v.Struct(objData)

		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				objValidate = validatorRepository.ValidatorResult{
					FieldName:   utils.KebabStyle(err.StructField(), "lower") + "." + strconv.Itoa(count),
					Validated:   err.ActualTag(),
					TypeData:    err.Kind().String(),
					Value:       err.Value(),
					Param:       err.Param(),
					Description: enums.FIELD_REQUIRED + " (" + utils.KebabStyle(err.StructField(), "lower") + ")",
				}

				res := responseRepository.Result{Status: http.StatusUnprocessableEntity, Datas: nil, Errors: objValidate}
				Context.JSON(http.StatusUnprocessableEntity, res)
				return
			}
		}
		count++
	}

	count2 := 1
	for _, objData := range objDatas {

		var promobjCheckout checkoutsRepository.DataCheckoutRepos
		promobjCheckout.Sku = objData.Sku
		promobjCheckout.Qty = objData.Qty

		checkSku, errSku := productsService.GetProductBySku(objData.Sku)
		if errSku != nil {
			basketsService.DeleteBasketsByBuyyerId(buyerId)
			arr := make(map[string]interface{})
			arr["message"] = "Sku Not Found on line " + strconv.Itoa(count2)
			res := responseRepository.Result{Status: http.StatusUnprocessableEntity, Datas: nil, Errors: arr}
			Context.JSON(http.StatusUnprocessableEntity, res)
			return
		}

		basketsService.InsertBaskets(buyerId, objData.Sku, objData.Qty)

		promobjCheckout.CheckSku = checkSku

		checkPromo, errPromo := promotionsService.GetPromoBySku(objData.Sku)
		if errPromo != nil {
			promobjCheckout.CheckPromo = nil

		} else {

			promobjCheckout.CheckPromo = checkPromo

		}

		objCheckouts = append(objCheckouts, &promobjCheckout)

		count2++
	}

	var sku_array []string

	for _, objCheckout := range objCheckouts {

		checkSkuOnArray := checkSkuOnArray(sku_array, objCheckout.Sku)
		if checkSkuOnArray {
			continue
		}

		sku_array = append(sku_array, objCheckout.Sku)

		getBasketDatas, errBasketData := basketsService.ReturnAllBasketsBySku(buyerId, objCheckout.Sku)

		if errBasketData != nil {
			basketsService.DeleteBasketsByBuyyerId(buyerId)
			arr := make(map[string]interface{})
			arr["message"] = errBasketData.Error()
			res := responseRepository.Result{Status: http.StatusUnprocessableEntity, Datas: nil, Errors: arr}
			Context.JSON(http.StatusUnprocessableEntity, res)
			return
		}

		totalBasketQty := 0
		for _, getBasketData := range getBasketDatas {
			totalBasketQty += getBasketData.Qty
		}

		var inrecSku, _ = json.Marshal(objCheckout.CheckSku)
		var FinalSku map[string]interface{}
		json.Unmarshal(inrecSku, &FinalSku)

		var inrecPromo, _ = json.Marshal(objCheckout.CheckPromo)
		var FinalPromo map[string]interface{}
		json.Unmarshal(inrecPromo, &FinalPromo)

		var float64QtySku float64 = FinalSku["qty"].(float64)
		var intQtySku int = int(float64QtySku)

		if intQtySku < totalBasketQty {
			basketsService.DeleteBasketsByBuyyerId(buyerId)
			arr := make(map[string]interface{})
			arr["message"] = "Inventory Not Available For Sku " + objCheckout.Sku
			res := responseRepository.Result{Status: http.StatusUnprocessableEntity, Datas: nil, Errors: arr}
			Context.JSON(http.StatusUnprocessableEntity, res)
			return
		}

		var var_PromoType string
		var_PromoType_readonly, ok := FinalPromo["promo_type"].(string)
		if ok {
			var_PromoType = var_PromoType_readonly
		} else {
			var_PromoType = ""
		}

		var var_PromoCode string
		var_PromoCode_readonly, ok := FinalPromo["promo_code"].(string)
		if ok {
			var_PromoCode = var_PromoCode_readonly
		} else {
			var_PromoCode = ""
		}

		if var_PromoType == "free_sku" {

			var var_MinimumQty int
			var_MinimumQty = 0
			var_MinimumQty_readonly, ok := FinalPromo["minimum_qty"].(float64)
			if ok {
				var_MinimumQty = int(var_MinimumQty_readonly)
			}

			if totalBasketQty >= var_MinimumQty {

				numberFree := totalBasketQty / var_MinimumQty

				var intnumberFree int = int(numberFree)

				var inrecPromoDetails, _ = json.Marshal(FinalPromo["details_promo"])
				var FinalPromoDetails []map[string]interface{}
				json.Unmarshal(inrecPromoDetails, &FinalPromoDetails)

				var var_PriceSku float64
				var_PriceSku = 0
				var_PriceSku_readonly, ok := FinalSku["price"].(float64)
				if ok {
					var_PriceSku = var_PriceSku_readonly
				}

				floatQty := float64(totalBasketQty)

				var var_PriceTotal float64
				var_PriceTotal = 0
				var_PriceTotal_readonly, ok := FinalSku["price"].(float64)
				if ok {
					var_PriceTotal = var_PriceTotal_readonly * floatQty
				}

				var datasCheckout models.Checkouts
				datasCheckout.IdBuyer = buyerId
				datasCheckout.Sku = objCheckout.Sku
				datasCheckout.Qty = totalBasketQty
				datasCheckout.Price = var_PriceSku
				datasCheckout.TotalPrice = var_PriceTotal
				datasCheckout.PromoCode = ""
				datasCheckout.PromoType = ""
				datasCheckout.Discount = 0
				datasCheckouts = append(datasCheckouts, &datasCheckout)

				for _, DetailPromo := range FinalPromoDetails {

					var float64Qty float64 = DetailPromo["qty"].(float64)
					var intQty int = int(float64Qty)
					var_PromoQty := intQty * intnumberFree

					var datasCheckoutPromo models.Checkouts
					datasCheckoutPromo.IdBuyer = buyerId
					datasCheckoutPromo.Sku = DetailPromo["sku"].(string)
					datasCheckoutPromo.Qty = var_PromoQty
					datasCheckoutPromo.Price = 0
					datasCheckoutPromo.TotalPrice = 0
					datasCheckoutPromo.PromoCode = var_PromoCode
					datasCheckoutPromo.PromoType = var_PromoType
					datasCheckoutPromo.Discount = 0
					datasCheckouts = append(datasCheckouts, &datasCheckoutPromo)

				}

			} else {

				var var_PriceSku float64
				var_PriceSku = 0
				var_PriceSku_readonly, ok := FinalSku["price"].(float64)
				if ok {
					var_PriceSku = var_PriceSku_readonly
				}

				floatQty := float64(totalBasketQty)

				var var_PriceTotal float64
				var_PriceTotal = 0
				var_PriceTotal_readonly, ok := FinalSku["price"].(float64)
				if ok {
					var_PriceTotal = var_PriceTotal_readonly * floatQty
				}

				var datasCheckout models.Checkouts
				datasCheckout.IdBuyer = buyerId
				datasCheckout.Sku = objCheckout.Sku
				datasCheckout.Qty = totalBasketQty
				datasCheckout.Price = var_PriceSku
				datasCheckout.TotalPrice = var_PriceTotal
				datasCheckout.PromoCode = ""
				datasCheckout.PromoType = ""
				datasCheckout.Discount = 0
				datasCheckouts = append(datasCheckouts, &datasCheckout)
			}

		} else if var_PromoType == "free_one_sku" {

			var var_MinimumQty int
			var_MinimumQty = 0
			var_MinimumQty_readonly, ok := FinalPromo["minimum_qty"].(float64)
			if ok {
				var_MinimumQty = int(var_MinimumQty_readonly)
			}

			if totalBasketQty >= var_MinimumQty {

				numberFree := totalBasketQty / var_MinimumQty

				var intnumberFree int = int(numberFree)

				var var_PriceSku float64
				var_PriceSku = 0
				var_PriceSku_readonly, ok := FinalSku["price"].(float64)
				if ok {
					var_PriceSku = var_PriceSku_readonly
				}

				floatQty := float64(totalBasketQty)
				floatnumberFree := float64(intnumberFree)

				var var_PriceTotal float64
				var_PriceTotal = 0

				var_PriceTotal_readonly, ok := FinalSku["price"].(float64)
				if ok {
					var_PriceTotal = var_PriceTotal_readonly * floatQty
				}

				var var_PriceTotalFree float64
				var_PriceTotalFree = var_PriceTotal_readonly * floatnumberFree

				var_PriceTotal = var_PriceTotal - var_PriceTotalFree

				var datasCheckout models.Checkouts
				datasCheckout.IdBuyer = buyerId
				datasCheckout.Sku = objCheckout.Sku
				datasCheckout.Qty = totalBasketQty
				datasCheckout.Price = var_PriceSku
				datasCheckout.TotalPrice = var_PriceTotal
				datasCheckout.PromoCode = var_PromoCode
				datasCheckout.PromoType = var_PromoType
				datasCheckout.Discount = 0
				datasCheckouts = append(datasCheckouts, &datasCheckout)

			} else {

				var var_PriceSku float64
				var_PriceSku = 0
				var_PriceSku_readonly, ok := FinalSku["price"].(float64)
				if ok {
					var_PriceSku = var_PriceSku_readonly
				}

				floatQty := float64(totalBasketQty)

				var var_PriceTotal float64
				var_PriceTotal = 0
				var_PriceTotal_readonly, ok := FinalSku["price"].(float64)
				if ok {
					var_PriceTotal = var_PriceTotal_readonly * floatQty
				}

				var datasCheckout models.Checkouts
				datasCheckout.IdBuyer = buyerId
				datasCheckout.Sku = objCheckout.Sku
				datasCheckout.Qty = totalBasketQty
				datasCheckout.Price = var_PriceSku
				datasCheckout.TotalPrice = var_PriceTotal
				datasCheckout.PromoCode = ""
				datasCheckout.PromoType = ""
				datasCheckout.Discount = 0
				datasCheckouts = append(datasCheckouts, &datasCheckout)

			}

		} else if var_PromoType == "discountAlexa" {

			var var_MinimumQty int
			var_MinimumQty = 0
			var_MinimumQty_readonly, ok := FinalPromo["minimum_qty"].(float64)
			if ok {
				var_MinimumQty = int(var_MinimumQty_readonly)
			}

			if totalBasketQty >= var_MinimumQty {

				numberFree := totalBasketQty / var_MinimumQty

				var intnumberFree int = int(numberFree)

				var var_PriceSku float64
				var_PriceSku = 0
				var_PriceSku_readonly, ok := FinalSku["price"].(float64)
				if ok {
					var_PriceSku = var_PriceSku_readonly
				}

				var var_DiscountSku float64
				var_DiscountSku = 0
				var_DiscountSku_readonly, ok := FinalPromo["discount"].(float64)
				if ok {
					var_DiscountSku = var_DiscountSku_readonly
				}

				floatQty := float64(totalBasketQty)
				totalDiscount := var_DiscountSku * float64(intnumberFree)
				dataDiscount := totalDiscount / 100
				var_PriceTotal := var_PriceSku * floatQty
				var_PriceTotalAfterDiscount := var_PriceTotal * dataDiscount

				var datasCheckout models.Checkouts
				datasCheckout.IdBuyer = buyerId
				datasCheckout.Sku = objCheckout.Sku
				datasCheckout.Qty = totalBasketQty
				datasCheckout.Price = var_PriceSku
				datasCheckout.TotalPrice = var_PriceTotalAfterDiscount
				datasCheckout.PromoCode = var_PromoCode
				datasCheckout.PromoType = var_PromoType
				datasCheckout.Discount = float64(totalDiscount)
				datasCheckouts = append(datasCheckouts, &datasCheckout)

			} else {

				var var_PriceSku float64
				var_PriceSku = 0
				var_PriceSku_readonly, ok := FinalSku["price"].(float64)
				if ok {
					var_PriceSku = var_PriceSku_readonly
				}

				floatQty := float64(totalBasketQty)

				var var_PriceTotal float64
				var_PriceTotal = 0
				var_PriceTotal_readonly, ok := FinalSku["price"].(float64)
				if ok {
					var_PriceTotal = var_PriceTotal_readonly * floatQty
				}

				var datasCheckout models.Checkouts
				datasCheckout.IdBuyer = buyerId
				datasCheckout.Sku = objCheckout.Sku
				datasCheckout.Qty = totalBasketQty
				datasCheckout.Price = var_PriceSku
				datasCheckout.TotalPrice = var_PriceTotal
				datasCheckout.PromoCode = ""
				datasCheckout.PromoType = ""
				datasCheckout.Discount = 0
				datasCheckouts = append(datasCheckouts, &datasCheckout)

			}

		} else {

			var var_PriceSku float64
			var_PriceSku = 0
			var_PriceSku_readonly, ok := FinalSku["price"].(float64)
			if ok {
				var_PriceSku = var_PriceSku_readonly
			}

			floatQty := float64(totalBasketQty)

			var var_PriceTotal float64
			var_PriceTotal = 0
			var_PriceTotal_readonly, ok := FinalSku["price"].(float64)
			if ok {
				var_PriceTotal = var_PriceTotal_readonly * floatQty
			}

			var datasCheckout models.Checkouts
			datasCheckout.IdBuyer = buyerId
			datasCheckout.Sku = objCheckout.Sku
			datasCheckout.Qty = totalBasketQty
			datasCheckout.Price = var_PriceSku
			datasCheckout.TotalPrice = var_PriceTotal
			datasCheckout.PromoCode = ""
			datasCheckout.PromoType = ""
			datasCheckout.Discount = 0
			datasCheckouts = append(datasCheckouts, &datasCheckout)
		}

	}

	var sku_input_array []string
	for _, datasCheckout := range datasCheckouts {

		checkSkuOnArray := checkSkuOnArray(sku_input_array, datasCheckout.Sku)
		if checkSkuOnArray {
			continue
		}

		sku_input_array = append(sku_input_array, datasCheckout.Sku)

		getDataCheckout, errCheck := getDataCheckout(datasCheckouts, datasCheckout.Sku)
		if errCheck != nil {
			basketsService.DeleteBasketsByBuyyerId(buyerId)
			arr := make(map[string]interface{})
			arr["message"] = "Sku Not Found on line "
			res := responseRepository.Result{Status: http.StatusUnprocessableEntity, Datas: nil, Errors: arr}
			Context.JSON(http.StatusUnprocessableEntity, res)
			return
		}

		var inrecCheckoutDetails, _ = json.Marshal(getDataCheckout)
		var FinalCheckoutDetails []map[string]interface{}
		json.Unmarshal(inrecCheckoutDetails, &FinalCheckoutDetails)

		var var_totalSkuCheckoutQty float64
		var_totalSkuCheckoutQty = 0
		for _, FinalCheckoutDetail := range FinalCheckoutDetails {

			var var_TotalCheckoutQty float64
			var_TotalCheckoutQty = 0
			var_TotalCheckoutQty_readonly, ok := FinalCheckoutDetail["qty"].(float64)
			if ok {
				var_TotalCheckoutQty = var_TotalCheckoutQty_readonly
			}

			var_totalSkuCheckoutQty += var_TotalCheckoutQty
		}

		checkSku, errSku := productsService.GetProductBySku(datasCheckout.Sku)
		if errSku != nil {
			basketsService.DeleteBasketsByBuyyerId(buyerId)
			arr := make(map[string]interface{})
			arr["message"] = errSku.Error()
			res := responseRepository.Result{Status: http.StatusUnprocessableEntity, Datas: nil, Errors: arr}
			Context.JSON(http.StatusUnprocessableEntity, res)
			return
		}

		if checkSku.Qty < int(var_totalSkuCheckoutQty) {
			basketsService.DeleteBasketsByBuyyerId(buyerId)
			arr := make(map[string]interface{})
			arr["message"] = "Inventroy not available for Sku " + datasCheckout.Sku
			res := responseRepository.Result{Status: http.StatusUnprocessableEntity, Datas: nil, Errors: arr}
			Context.JSON(http.StatusUnprocessableEntity, res)
			return
		}

	}

	for _, dataCheckoutInsert := range datasCheckouts {

		_, errQty := productsService.UpdateProductBySku(dataCheckoutInsert.Sku, dataCheckoutInsert.Qty)
		if errQty != nil {
			basketsService.DeleteBasketsByBuyyerId(buyerId)
			arr := make(map[string]interface{})
			arr["message"] = errQty.Error()
			res := responseRepository.Result{Status: http.StatusUnprocessableEntity, Datas: nil, Errors: arr}
			Context.JSON(http.StatusUnprocessableEntity, res)
			return
		}

		basketsService.DeleteBasketsByBuyyerId(buyerId)
		checkoutsService.InsertCheckouts(buyerId, dataCheckoutInsert.Sku, dataCheckoutInsert.Qty, dataCheckoutInsert.Price, dataCheckoutInsert.TotalPrice, dataCheckoutInsert.PromoCode, dataCheckoutInsert.PromoType, dataCheckoutInsert.Discount)
	}

	arr := make(map[string]interface{})
	arr["datasCheckouts"] = datasCheckouts
	// arr["objDatas"] = objCheckouts

	res := responseRepository.Result{Status: 200, Datas: arr, Errors: buyerId}
	Context.JSON(http.StatusOK, res)
	return

}

func getDataCheckout(s []*models.Checkouts, str string) ([]*models.Checkouts, error) {
	var datas []*models.Checkouts
	for _, datasCheckout := range s {
		if datasCheckout.Sku == str {
			datas = append(datas, datasCheckout)
		}
	}

	return datas, nil
}

func checkSkuOnArray(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
