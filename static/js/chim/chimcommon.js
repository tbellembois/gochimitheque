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

// handle HTTP errors
function handleHTTPError(msgText, msgStatus) {
    //console.log(msgStatus);
    switch(msgStatus) {
    case 401:
        global.displayMessage(msgText, "danger");
        break;
    case 403:
        global.displayMessage(msgText, "danger");
        break;
    case 500:
        global.displayMessage(msgText, "danger");
        break;
    default:
        global.displayMessage(msgText, "light");
        break;
    }
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

function cleanQueryParams() {

    //console.log("cleanQueryParams")

    // root url
    var root = window.location.protocol + '//' + window.location.hostname + ":" + window.location.port + window.location.pathname;
    // var urlParams = new URLSearchParams(window.location.search);
    // urlParams.delete("search");
    // urlParams.delete("sort");
    // urlParams.delete("order");
    // urlParams.delete("offset");
    // urlParams.delete("limit");
    // urlParams.delete("export");

    // window.history.replaceState("", "", root + "?" + urlParams.toString());
    window.history.replaceState("", "", root);
    $("#hidden_s_entity").val("");
    $("#hidden_s_history").val("");
    $("#hidden_s_storage").val("");
    $("#hidden_s_storage_archive").val("")
    $("#hidden_s_product").val("");
    $("#hidden_s_bookmark").val("");
    $("#filter-item").html("");
}

function exportAll() {
    // root path
    path = window.location.pathname
    // root url
    root = window.location.protocol + '//' + window.location.hostname + ":" + window.location.port + path;
    // building url parameters from last AJAX query parameters
    p = lastQueryParams.data;
    // filtering empty values
    newp = {};
    $.each(p, function(k, v) {
        if (v !== "") {
            newp[k] = v;
        }
    });

    // adding export param
    newp["export"] = true;

    // redirecting
    window.location.href = root + "?" + $.param(newp);
}

function switchProductStorageView() {
    // root path
    path = window.location.pathname
    // replacing path element
    if (path.indexOf("storages") >= 0) {
        path = path.replace("storages", "products");
    } else {
        path = path.replace("products", "storages");
    }
    // root url
    root = window.location.protocol + '//' + window.location.hostname + ":" + window.location.port + path;
    // building url parameters from last AJAX query parameters
    p = lastQueryParams.data;
    // filtering empty values
    newp = {};
    $.each(p, function(k, v) {
        if (v !== "") {
            newp[k] = v;
        }
    });

    // redirecting
    window.location.href = root + "?" + $.param(newp);
}

function clearsearch() {
    // root url
    root = window.location.protocol + '//' + window.location.hostname + ":" + window.location.port + proxyPath;
    window.location.href = root + "v/products";
}

function search() {

    //console.log("search")

    cleanQueryParams();

    var s_custom_name_part_of;
    var s_storage_barecode;
    var s_symbols;
    var s_hazardstatements;
    var s_precautionarystatements;

    query = {};

    if ($('select#s_storelocation').hasClass("select2-hidden-accessible")) {
        // Select2 has been initialized
        // storelocation_id
        i = $('select#s_storelocation').select2('data')[0];
        if (i != undefined) {
            //s_storelocation = i.id;
            $.extend(query, {
                "storelocation": i.id
            });
            $('#s_storelocation').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
        }
    }
    if ($('select#s_name').hasClass("select2-hidden-accessible")) {
        // Select2 has been initialized
        // name_id
        i = $('select#s_name').select2('data')[0];
        if (i != undefined) {
            s_name = i.id;
            $.extend(query, {
                "name": i.id
            });
            $('#s_name').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
        }   
    }
    if ($('select#s_empiricalformula').hasClass("select2-hidden-accessible")) {
        // Select2 has been initialized
        // empiricalformula_id
        i = $('select#s_empiricalformula').select2('data')[0];
        if (i != undefined) {
            //s_empiricalformula = $('select#s_empiricalformula').select2('data')[0].id;
            $.extend(query, {
                "empiricalformula": i.id
            });
            $('#s_empiricalformula').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
        }
    }
    if ($('select#s_casnumber').hasClass("select2-hidden-accessible")) {
        // Select2 has been initialized
        // casnumber_id
        i = $('select#s_casnumber').select2('data')[0];
        if (i != undefined) {
            //s_casnumber = $('select#s_casnumber').select2('data')[0].id;
            $.extend(query, {
                "casnumber": i.id
            });
            $('#s_casnumber').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
        }
    }
    if ($('select#s_signalword').hasClass("select2-hidden-accessible")) {
        // Select2 has been initialized
        // signalword_id
        i = $('select#s_signalword').select2('data')[0];
        if (i != undefined) {
            //s_signalword = $('select#s_signalword').select2('data')[0].id;
            $.extend(query, {
                "signalword": i.id
            });
            $('#s_signalword').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
        }
    }
    if ($('select#s_symbols').hasClass("select2-hidden-accessible")) {
        // Select2 has been initialized
        // symbols_id
        i = $('select#s_symbols').select2('data');
        if (i.length != 0) {
            s_symbols = [];
            i.forEach(function(e) {
                s_symbols.push(e.id);
            });
            $.extend(query, {
                "symbols": s_symbols
            });
            $('#s_symbols').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
        }
    }
    if ($('select#s_hazardstatements').hasClass("select2-hidden-accessible")) {
        // Select2 has been initialized
        // hazardstatements_id
        i = $('select#s_hazardstatements').select2('data');
        if (i.length != 0) {
            s_hazardstatements = [];
            i.forEach(function(e) {
                s_hazardstatements.push(e.id);
            });
            $.extend(query, {
                "hazardstatements": s_hazardstatements
            });
            $('#s_hazardstatements').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
        }
    }
    if ($('select#s_precautionarystatements').hasClass("select2-hidden-accessible")) {
        // Select2 has been initialized
        // precautionarystatement_id
        i = $('select#s_precautionarystatements').select2('data');
        if (i.length != 0) {
            s_precautionarystatements = [];
            i.forEach(function(e) {
                s_precautionarystatements.push(e.id);
            });
            $.extend(query, {
                "precautionarystatements": s_precautionarystatements
            });
            $('#s_precautionarystatements').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
        }
    }
    if ($('#s_casnumber_cmr:checked').length > 0) {
        //s_casnumber_cmr = true;
        $.extend(query, {
            "casnumber_cmr": true
        });
        $('#s_casnumber_cmr').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
    }

    s_storage_barecode = $('#s_storage_barecode').val() ;
    if (s_storage_barecode != "") {
        //params["storage_barecode"] = s_storage_barecode;
        $.extend(query, {
            "storage_barecode": s_storage_barecode
        });
        $('#s_storage_barecode').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
    }
    s_custom_name_part_of = $('#s_custom_name_part_of').val() ;
    if (s_custom_name_part_of != "") {
        //params["custom_name_part_of"] = s_custom_name_part_of;
        $.extend(query, {
            "custom_name_part_of": s_custom_name_part_of
        });
        $('#s_custom_name_part_of').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
    }

    var $table = $('#table');
    $table.bootstrapTable('refresh', {query: query});
}
