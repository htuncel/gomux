package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"main/utils"
)

// ExampleController entails all the methods concerning examples a boring man
type ExampleController struct{}

// ExampleService interface of the example controller
type ExampleService interface {
	tokenHandler(w http.ResponseWriter, r *http.Request)
}

// NewExampleController to test file upload
func NewExampleController() *ExampleController {
	return new(ExampleController)
}

// TokenHandler godoc
// @Summary Get token
// @Description Get token
// @Tags token
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /example/token [get]
func (e *ExampleController) TokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, errToken := utils.GenerateToken()
	claims, _ := utils.VerifyToken(token)
	host := utils.GetClientIP(w, r)
	token2 := utils.GetToken(w, r)

	if errToken != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": errToken.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"data":   token,
			"claims": claims,
			"host":   host,
			"token2": token2,
		},
	)
}

// MultipleFileUpload godoc
// @Summary	Upload multiple file
// @Description	Upload multiple file
// @Tags		file
// @Accept		multipart/form-data
// @Produce		json
// @Param		files	formData[]file true  "files"
// @Success		200		{object}map[string]string
// @Failure		400		{object}map[string]string
// @Router		/example/multifile [post]
func (e *ExampleController) MultipleFileUpload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// 32MB is the default used by FormFile
	r.ParseMultipartForm(32 << 20)
	files := r.MultipartForm.File["files"]

	success := 0
	failure := 0
	for _, file := range files {
		f, err := file.Open()
		defer f.Close()
		filename := file.Filename
		extension := path.Ext(filename)

		if err != nil {
			failure++
			continue
		}

		if extension != ".xml" {
			failure++
			continue
		}
		dst, err := os.Create("upload/xml/" + filename)
		defer dst.Close()
		if err != nil {
			failure++
			continue
		}
		if _, err := io.Copy(dst, f); err != nil {
			failure++
			continue
		}
		success++
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"success": fmt.Sprintf("%d files uploaded!", success),
		"failure": fmt.Sprintf("%d files could not be uploaded!", failure),
	})
}
