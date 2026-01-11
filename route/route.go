package route

import (
	"net/http"

	"tatria/controller"
)

func Routes(router *http.ServeMux) {
	router.HandleFunc("/process", controller.Process)
}
