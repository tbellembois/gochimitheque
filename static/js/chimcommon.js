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
    } else {
        newfdata[k] = v;
    }
  });
  return newfdata;
};

// displays and fadeout the given message
function displayMessage(msgText, type) {
    var d = $("<div>");
    d.attr("role", "alert");
    d.addClass("alert alert-" + type);
    d.text(msgText);
    $("body").prepend(d.delay(800).fadeOut("slow"));
}
