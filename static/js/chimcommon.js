// https://coderwall.com/p/w22s0w/recursive-merge-flatten-objects-in-plain-old-vanilla-javascript
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

var normalizeSqlNull = function(obj) {
  newfdata = new Map()
  $.each(obj, function(k, v) {
    matchs = k.match(/(.+)\.String/);
    matchi = k.match(/(.+)\.Int64/);
    matchf = k.match(/(.+)\.Float64/);
    matchb = k.match(/(.+)\.Bool/);
    if (matchs !== null) {
        fieldname = matchs[1];
        valid = fdata[fieldname+".Valid"] == true;
        if (valid) {
            newfdata[fieldname] = v;
        }
    }  else if (matchi !== null) {
        fieldname = matchi[1];
        valid = fdata[fieldname+".Valid"] == true;
        if (valid) {
            newfdata[fieldname] = v;
        }
    }  else if (matchb !== null) {
        fieldname = matchb[1];
        valid = fdata[fieldname+".Valid"] == true;
        if (valid) {
            newfdata[fieldname] = v;
        }
    } else if (matchf !== null) {
        fieldname = matchf[1];
        valid = fdata[fieldname+".Valid"] == true;
        if (valid) {
            newfdata[fieldname] = v;
        }
    }
    else {
        newfdata[k] = v;
    }
  });
  return newfdata;
};

function createTitle(msgText) {
    var d = $("<div>");
    var i = $("<i>");
    var h = $("<h6>");
    d.addClass("text-center mx-auto mt-md-5 mb-md-5");
    i.addClass("material-icons").text("keyboard_tab");
    h.text(msgText);
    d.append(i);
    d.append(h);
    return d;
}

// displays and fadeout the given message
function displayMessage(msgText, type) {
    var d = $("<div>");
    d.attr("role", "alert");
    d.addClass("alert alert-" + type);
    d.text(msgText);
    $("body").prepend(d.delay(800).fadeOut("slow"));
}

function highlightRow(attr, id) {
    $("tr[" + attr + "=" + id + "]").fadeTo("slow", 0.30);
    $("tr[" + attr + "=" + id + "]").fadeTo("slow", 1);
    $("tr[" + attr + "=" + id + "]").fadeTo("slow", 0.30);
    $("tr[" + attr + "=" + id + "]").fadeTo("slow", 1);
    $("tr[" + attr + "=" + id + "]").fadeTo("slow", 0.30);
    $("tr[" + attr + "=" + id + "]").fadeTo("slow", 1);
}

// https://gist.github.com/excalq/2961415
var updateQueryStringParam = function (key, value) {

    var baseUrl = [location.protocol, '//', location.host, location.pathname].join(''),
        urlQueryString = document.location.search,
        newParam = key + '=' + value,
        params = '?' + newParam;

    // If the "search" string exists, then build params from it
    if (urlQueryString) {
        var updateRegex = new RegExp('([\?&])' + key + '[^&]*');
        var removeRegex = new RegExp('([\?&])' + key + '=[^&;]+[&;]?');

        if( typeof value == 'undefined' || value == null || value == '' ) { // Remove param if value is empty
            params = urlQueryString.replace(removeRegex, "$1");
            params = params.replace( /[&;]$/, "" );

        } else if (urlQueryString.match(updateRegex) !== null) { // If param exists already, update it
            params = urlQueryString.replace(updateRegex, "$1" + newParam);

        } else { // Otherwise, add it to end of query string
            params = urlQueryString + '&' + newParam;
        }
    }

    // no parameter was set so we don't need the question mark
    params = params == '?' ? '' : params;

    window.history.replaceState({}, "", baseUrl + params);
};