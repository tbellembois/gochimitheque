package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"github.com/tbellembois/gochimitheque/helpers"
)

// DownloadExportHandler serves the temporary export files
func (env *Env) DownloadExportHandler(w http.ResponseWriter, r *http.Request) *helpers.AppError {

	vars := mux.Vars(r)
	var (
		id  string // temporary file id
		err error
		tb  []byte
		ok  bool
	)

	if id, ok = vars["id"]; !ok {
		return &helpers.AppError{
			Error:   err,
			Message: "no temporary file id",
			Code:    http.StatusBadRequest}
	}

	// full temporary file path
	ftp := path.Join(os.TempDir(), "chimitheque-"+id)

	// reading the file
	if tb, err = ioutil.ReadFile(ftp); err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error reading the temporary file",
		}
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=chimitheque.csv")

	// stream file
	b := bytes.NewBuffer(tb)
	if _, err := b.WriteTo(w); err != nil {
		return &helpers.AppError{
			Error:   err,
			Code:    http.StatusInternalServerError,
			Message: "error streaming the temporary file",
		}
	}

	w.Write([]byte("export finished"))

	return nil
}
