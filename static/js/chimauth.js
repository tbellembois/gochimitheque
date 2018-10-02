// handle HTTP errors
function handleHTTPError(msgText, msgStatus) {
    //console.log(msgStatus);
    switch(msgStatus) {
    case 401:
        displayMessage(msgText, "danger");
        // redirect on 401 errors
        window.location.replace("/login");
        break;
    case 403:
        displayMessage(msgText, "danger");
        break;
    case 500:
        displayMessage(msgText, "danger");
        break;
    default:
        displayMessage(msgText, "light");
        break;
    }
}

// get a JWT token for the user
function getToken() {
    var email = $("#person_email").val(),
        password = $("#person_password").val();

    $.ajax({
        url: "/get-token",
        method: 'POST',
        data: {
            person_email: email,
            person_password: password
        }
    }).done(function(token) {
        //console.log(token);
        // store in web storage
        //window.localStorage.setItem('token', token);
        window.location.replace("/v/entities");
    });
}

// returns the cookie value
function readCookie(name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(';');
    for(var i=0;i < ca.length;i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1,c.length);
        if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
    }
    return null;
}

function hasPermission(url, method, itemId) {
    return $.ajax({
        url: url,
        itemId: itemId,
        method: method,
    });
}

// jwt in web storage
//$.ajaxPrefilter(function( options ) {
//    options.beforeSend = function (xhr) {
//        xhr.setRequestHeader('Authorization', 'Bearer '+localStorage.getItem('token'));
//    }
//});
