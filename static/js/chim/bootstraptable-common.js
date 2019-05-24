function bt_formatLoadingMessage() {
    return global.t("bt_loadingMessage", container.PersonLanguage);
}
function bt_formatRecordsPerPage(pageNumber) {
    return pageNumber + " " + global.t("bt_recordsPerPage", container.PersonLanguage);
}
function bt_formatShowingRows(pageFrom, pageTo, totalRows) {
    return pageFrom + "/" + pageTo + "(" + totalRows + " " + global.t("bt_showingRowsTotal", container.PersonLanguage) + ")";
}
function bt_formatSearch() {
    return global.t("bt_search", container.PersonLanguage);
}
function bt_formatNoMatches() {
    return global.t("bt_noMatches", container.PersonLanguage);
}