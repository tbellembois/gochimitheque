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

// DownloadExportHandler serve the export files.
func (env *Env) DownloadExportHandler(w http.ResponseWriter, r *http.Request) *models.AppError {
	var (
		fileID string
		ok     bool
		err    error
	)

	vars := mux.Vars(r)

	if fileID, ok = vars["id"]; !ok {
		return &models.AppError{
			OriginalError: err,
			Message:       "no query file id",
			Code:          http.StatusBadRequest,
		}
	}

	fileFullPath := path.Join(os.TempDir(), "chimitheque-"+fileID)

	var fileData []byte

	if fileData, err = ioutil.ReadFile(fileFullPath); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error reading the file",
		}
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=chimitheque.csv")

	// Stream file.
	b := bytes.NewBuffer(fileData)
	if _, err = b.WriteTo(w); err != nil {
		return &models.AppError{
			OriginalError: err,
			Code:          http.StatusInternalServerError,
			Message:       "error streaming the file",
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
