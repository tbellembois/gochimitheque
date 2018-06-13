//https://coderwall.com/p/w22s0w/recursive-merge-flatten-objects-in-plain-old-vanilla-javascript
var merge = function(objects) {
    var out = {};
  
    for (var i = 0; i < objects.length; i++) {
      for (var p in objects[i]) {
        out[p] = objects[i][p];
      }
    }
  
    return out;
}

var flatten = function(obj, name, stem) {
    var out = {};
    var newStem = (typeof stem !== 'undefined' && stem !== '') ? stem + '.' + name : name;
    
    if (typeof obj !== 'object') {
      out[newStem] = obj;
      return out;
    }
    
    for (var p in obj) {
      var prop = flatten(obj[p], p, newStem);
      out = merge([out, prop]);
    }

    return out;
};

function displayMessage(msgText, type) {
    var d = $("<div>");
    d.attr("role", "alert");
    d.addClass("alert alert-" + type);
    d.text(msgText);
    $("body").prepend(d.delay(800).fadeOut("slow"));
}

function handleHTTPError(msgText, msgStatus) {
    console.log(msgStatus);
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
        console.log(token);
        // store in web storage
        //window.localStorage.setItem('token', token);
        window.location.replace("/v/entities");
    });
}

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

function hasPermission(personId, perm, item, itemId) {
    return $.ajax({
        url: "/haspermission/" + personId + "/" + perm + "/" + item + "/" + itemId,
        itemId: itemId,
        method: 'GET',
    });
}

// jwt in web storage
//$.ajaxPrefilter(function( options ) {
//    options.beforeSend = function (xhr) {
//        xhr.setRequestHeader('Authorization', 'Bearer '+localStorage.getItem('token'));
//    }
//});
