$( document ).ready(function() { 

    $.ajax({
        url: proxyPath + "welcomeannounce",
        method: "GET",
    }).done(function(data, textStatus, jqXHR) {
        $('#welcomeannounce_text').trumbowyg({
            btns: [
                ['foreColor'], 
                ['backColor'],
                ['viewHTML'],
                ['formatting'],
                ['strong', 'em', 'del'],
                ['superscript', 'subscript'],
                ['link'],
                ['justifyLeft', 'justifyCenter', 'justifyRight', 'justifyFull'],
                ['unorderedList', 'orderedList'],
                ['horizontalRule'],
                ['removeformat']
            ]
        });
        $('#welcomeannounce_text').trumbowyg('html', data.welcomeannounce_text);

    }).fail(function(jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    }); 
})

//
// product_slFormatter formatter
//
function saveWelcomeAnnounce() {
    var welcomeannounce_text = $("#welcomeannounce_text").trumbowyg('html');

    $.ajax({
        url: proxyPath + "welcomeannounce",
        method: "PUT",
        dataType: 'json',
        data: {"welcomeannounce_text": welcomeannounce_text},
    }).done(function(jqXHR, textStatus, errorThrown) {
        msg = global.t("welcomeannounce_text_modificationsuccess", container.PersonLanguage) 
        global.displayMessage(msg, "success");
    }).fail(function(jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });  
}