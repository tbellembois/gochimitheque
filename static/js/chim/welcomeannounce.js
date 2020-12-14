$(document).ready(function () {

    $.ajax({
        url: proxyPath + "welcomeannounce",
        method: "GET",
    }).done(function (data, textStatus, jqXHR) {
        $('#welcomeannounce_text').html(data.welcomeannounce_text);
    }).fail(function (jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });
})

//
// product_slFormatter formatter
//
function saveWelcomeAnnounce() {
    var welcomeannounce_text = $("#welcomeannounce_text").val();

    $.ajax({
        url: proxyPath + "welcomeannounce",
        method: "PUT",
        dataType: 'json',
        data: { "welcomeannounce_text": welcomeannounce_text },
    }).done(function (jqXHR, textStatus, errorThrown) {
        msg = gjsUtils.translate("welcomeannounce_text_modificationsuccess", container.PersonLanguage)
        gjsUtils.message(msg, "success");
    }).fail(function (jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });
}