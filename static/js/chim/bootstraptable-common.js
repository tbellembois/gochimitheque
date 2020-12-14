function bt_formatLoadingMessage() {
    return gjsUtils.translate("bt_loadingMessage", container.PersonLanguage);
}
function bt_formatRecordsPerPage(pageNumber) {
    return pageNumber + " " + gjsUtils.translate("bt_recordsPerPage", container.PersonLanguage);
}
function bt_formatShowingRows(pageFrom, pageTo, totalRows) {
    return pageFrom + "/" + pageTo + "(" + totalRows + " " + gjsUtils.translate("bt_showingRowsTotal", container.PersonLanguage) + ")";
}
function bt_formatSearch() {
    return gjsUtils.translate("bt_search", container.PersonLanguage);
}
function bt_formatNoMatches() {
    return gjsUtils.translate("bt_noMatches", container.PersonLanguage);
}