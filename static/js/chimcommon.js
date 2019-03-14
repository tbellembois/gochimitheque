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

// var normalizeSqlNull = function(obj) {
//   newfdata = new Map()
//   $.each(obj, function(k, v) {
//     matchs = k.match(/(.+)\.String/);
//     matchi = k.match(/(.+)\.Int64/);
//     matchf = k.match(/(.+)\.Float64/);
//     matchb = k.match(/(.+)\.Bool/);
//     if (matchs !== null) {
//         fieldname = matchs[1];
//         valid = fdata[fieldname+".Valid"] == true;
//         if (valid) {
//             newfdata[fieldname] = v;
//         }
//     }  else if (matchi !== null) {
//         fieldname = matchi[1];
//         valid = fdata[fieldname+".Valid"] == true;
//         if (valid) {
//             newfdata[fieldname] = v;
//         }
//     }  else if (matchb !== null) {
//         fieldname = matchb[1];
//         valid = fdata[fieldname+".Valid"] == true;
//         if (valid) {
//             newfdata[fieldname] = v;
//         }
//     } else if (matchf !== null) {
//         fieldname = matchf[1];
//         valid = fdata[fieldname+".Valid"] == true;
//         if (valid) {
//             newfdata[fieldname] = v;
//         }
//     }
//     else {
//         newfdata[k] = v;
//     }
//   });
//   return newfdata;
// };

// function createTitle(msgText, type) {
//     var i=$("<i>").addClass("material-icons");
//     switch (type) {
//         case 'history':
//             i.text("alarm");
//             break;
//         case 'bookmark':
//             i.text("star");
//             break;
//         case 'entity':
//             i.text("store");
//             break;
//         case 'storelocation':
//             i.text("extension");
//             break;
//         case 'product':
//             i.text("local_offer");
//             break;
//         case 'storage':
//             i.text("inbox");
//             break;
//         default:
//             i.text("keyboard_tab");
//     }

//     var d = $("<div>");
//     var h = $("<span>");
//     d.addClass("mt-md-3 mb-md-3 row");
//     h.addClass("col-sm-11 align-bottom").text(msgText);
//     d.append(i);
//     d.append(h);
//     return d;
// }

// displays and fadeout the given message
// function displayMessage(msgText, type) {
//     var d = $("<div>");
//     d.attr("role", "alert");
//     d.addClass("alert alert-" + type);
//     d.text(msgText);
//     $("body").prepend(d.delay(800).fadeOut("slow"));
// }

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
        }
    }
    if ($('#s_casnumber_cmr:checked').length > 0) {
        //s_casnumber_cmr = true;
        $.extend(query, {
            "casnumber_cmr": true
        });
    }

    s_storage_barecode = $('#s_storage_barecode').val() ;
    if (s_storage_barecode != "") {
        //params["storage_barecode"] = s_storage_barecode;
        $.extend(query, {
            "storage_barecode": s_storage_barecode
        });
    }
    s_custom_name_part_of = $('#s_custom_name_part_of').val() ;
    if (s_custom_name_part_of != "") {
        //params["custom_name_part_of"] = s_custom_name_part_of;
        $.extend(query, {
            "custom_name_part_of": s_custom_name_part_of
        });
    }


    var $table = $('#table');
    $table.bootstrapTable('refresh', {query: query}); 
}
