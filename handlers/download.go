package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"github.com/tbellembois/gochimitheque/models"
)

// DownloadExportHandler serves the temporary export files
func (env *Env) DownloadExportHandler(w http.ResponseWriter, r *http.Request) *models.AppError {

	vars := mux.Vars(r)
	var (
		id  string // temporary file id
		err error
		tb  []byte
		ok  bool
	)

	if id, ok = vars["id"]; !ok {
		return &models.AppError{
			Error:   err,
			Message: "no temporary file id",
			Code:    http.StatusBadRequest}
	}

	// full temporary file path
	ftp := path.Join(os.TempDir(), "chimitheque-"+id)

	// reading the file
	if tb, err = ioutil.ReadFile(ftp); err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error reading the temporary file",
		}
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=chimitheque.csv")

	// stream file
	b := bytes.NewBuffer(tb)
	if _, err = b.WriteTo(w); err != nil {
		return &models.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error streaming the temporary file",
		}
	}

	if _, err = w.Write([]byte("export finished")); err != nil {
		return &models.AppError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}
