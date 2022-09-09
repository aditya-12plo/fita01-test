package productsController

import (
	responseRepository "fita-test-01/repositories/responseRepository"
	productsService "fita-test-01/services/productsService"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(Context *gin.Context) {

	data, err := productsService.ReturnAllProducts()
	if err != nil {
		arr := make(map[string]interface{})
		arr["message"] = err
		res := responseRepository.Result{Status: 500, Datas: nil, Errors: arr}
		Context.JSON(500, res)
		return
	}
	res := responseRepository.Result{Status: 200, Datas: data, Errors: nil}
	Context.JSON(http.StatusOK, res)
	return
}
