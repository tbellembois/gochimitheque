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

// displays and fadeout the given message
function displayMessage(msgText, type) {
    var d = $("<div>");
    d.attr("role", "alert");
    d.addClass("alert alert-" + type);
    d.text(msgText);
    $("body").prepend(d.delay(800).fadeOut("slow"));
}

// // buidPermissionWidget returns a permission widget for the given person
// function buildPermissionWidget(persId) {
//     var widget = $("div").addClass("row");
//     var items = [
//                 'product', 
//                 'rproduct', 
//                 'storage', 
//                 'astorage', 
//                 'storelocation',
//                 'classofcompounds',
//                 'supplier'];

//     items.forEach(function(item, index, array) {
//         widget.append($("div").addClass("row").append($("div").addClass("col-sm").html(item)));
//         console.log(item, index);
//     });

//     return widget
// }