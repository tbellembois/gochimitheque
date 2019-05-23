function bt_formatLoadingMessage() {
    return global.t("bt_loadingMessage", container.PersonLanguage);
}
function bt_formatRecordsPerPage(pageNumber) {
    return pageNumber + " " + global.t("bt_recordsPerPage", container.PersonLanguage);
}
function bt_formatShowingRows(pageFrom, pageTo, totalRows) {
    return global.t("bt_showingRowsFrom", container.PersonLanguage) + pageFrom + global.t("bt_showingRowsTo", container.PersonLanguage) + pageTo + global.t("bt_showingRowsTotal", container.PersonLanguage);
}