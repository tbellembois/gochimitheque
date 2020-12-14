// getting welcome announce
$(document).ready(function () {
    $.ajax({
        url: proxyPath + "welcomeannounce",
        method: "GET",
    }).done(function (data, textStatus, jqXHR) {
        if (data.welcomeannounce_html != "") {
            var html = [];
            html.push("<div class='card'>")
            html.push("<div class='card-header text-center'>")
            html.push("<span class='mdi mdi-information-outline mdi-36px'/>")
            html.push("</div>")
            html.push("<div class='card-body'>")
            html.push(data.welcomeannounce_html)
            html.push("</div>")
            html.push("</div>")

            $("#wannounce").html(html.join(''));
        }
    }).fail(function (jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });
})

// send a password initialization link
function resetPassword() {
    var email = $("#person_email").val(),
        captcha = $("input#captcha_text").val(),
        uid = $("input#captcha_uid").val();
    if (email == "") {
        msg = gjsUtils.translate("resetpassword_warning_enteremail", container.PersonLanguage)
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
        }).done(function (token) {
            msg = gjsUtils.translate("resetpassword_message_mailsentto", container.PersonLanguage)
            gjsUtils.message(msg + " " + email, "success");
        }).fail(function (jqXHR, textStatus, errorThrown) {
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
        gjsUtils.message(gjsUtils.translate("resetpassword_warning_enteremail", container.PersonLanguage), "warning");
    } else {
        $.ajax({
            url: proxyPath + "get-token",
            method: 'POST',
            data: {
                person_email: email,
                person_password: password
            }
        }).done(function (token) {
            //console.log(token);
            // store in web storage
            //window.localStorage.setItem('token', token);
            // cleaning local storage
            localStorage.clear();
            window.location.replace(proxyPath + "v/products");
        }).fail(function (jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.responseText, jqXHR.status)
        });
    }
}

// get a captcha to resolve
function getCaptcha() {
    var email = $("#person_email").val();
    if (email == "") {
        msg = gjsUtils.translate("resetpassword_warning_enteremail", container.PersonLanguage)
        gjsUtils.message(msg, "warning");
    } else {
        msg = gjsUtils.translate("resetpassword_areyourobot", container.PersonLanguage)
        gjsUtils.message(msg, "success");
        $("#getcaptcha").fadeOut("slow");
        $.ajax({
            url: proxyPath + "captcha",
            method: 'GET',
        }).done(function (data) {
            $("input#captcha_uid").val(data.uid);
            $("img#captcha-img").attr("src", "data:image/png;base64," + data.image);
            $("div#captcha-row").removeClass("invisible");
        }).fail(function (jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.responseText, jqXHR.status)
        });
    }
}