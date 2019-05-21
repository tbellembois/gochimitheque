// getting welcome announce
$( document ).ready(function() {
    $.ajax({
        url: proxyPath + "welcomeannounce",
        method: "GET",
    }).done(function(data, textStatus, jqXHR) {
        $("#wannounce").html(data.welcomeannounce_text);
    }).fail(function(jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });  
})

// send a password initialization link
function resetPassword() {
    var email = $("#person_email").val(),
        captcha = $("input#captcha_text").val(),
        uid = $("input#captcha_uid").val();
    if (email == "") {
        msg = global.t("resetpassword_warning_enteremail", container.PersonLanguage) 
    } else {
        $("#resetpassword").fadeOut("slow");
        $.ajax({
            url: proxyPath + "reset-password",
            method: 'POST',
            data: {
                person_email: email,
                captcha_text: captcha,
                captcha_uid: uid,
            }
        }).done(function(token) {
            msg = global.t("resetpassword_message_mailsentto", container.PersonLanguage) 
            global.displayMessage(msg + " " + email, "success");
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.responseText, jqXHR.status)
        });
        $("#person_email").val("");
        $("input#captcha_text").val("");
        $("div#captcha-row").addClass("invisible");
    }
}

// get a JWT token for the user
function getToken() {
    var email = $("#person_email").val(),
        password = $("#person_password").val();

    if (email == "") {
        global.displayMessage(global.t("resetpassword_warning_enteremail", container.PersonLanguage), "warning");
    } else {
        $.ajax({
            url: proxyPath + "get-token",
            method: 'POST',
            data: {
                person_email: email,
                person_password: password
            }
        }).done(function(token) {
            //console.log(token);
            // store in web storage
            //window.localStorage.setItem('token', token);
            window.location.replace(proxyPath + "v/products");
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.responseText, jqXHR.status)
        });
    }
}

// get a captcha to resolve
function getCaptcha() {
    var email = $("#person_email").val();
    if (email == "") {
        msg = global.t("resetpassword_warning_enteremail", container.PersonLanguage) 
        global.displayMessage(msg, "warning");
    } else {
        msg = global.t("resetpassword_areyourobot", container.PersonLanguage) 
        global.displayMessage(msg, "success");
        $("#getcaptcha").fadeOut("slow");
        $.ajax({
            url: proxyPath + "captcha",
            method: 'GET',
        }).done(function(data) {
            $("input#captcha_uid").val(data.uid);
            $("img#captcha-img").attr("src", "data:image/png;base64," + data.image);
            $("div#captcha-row").removeClass("invisible");
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.responseText, jqXHR.status)
        });
    }
}