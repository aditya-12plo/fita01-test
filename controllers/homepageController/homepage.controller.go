package homepageController

import (
	responseRepository "fita-test-01/repositories/responseRepository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(Context *gin.Context) {

	arr := make(map[string]interface{})
	arr["message"] = "Welcome to testing 01"

	res := responseRepository.Result{Status: 200, Datas: arr, Errors: nil}
	Context.JSON(http.StatusOK, res)
}
