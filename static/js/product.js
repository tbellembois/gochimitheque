//
// request performed at table data loading
//
function getData(params) {
    // saving the query parameters
    lastQueryParams = params;
    $.ajax({
        url: proxyPath + "products",
        method: "GET",
        dataType: "JSON",
        data: params.data,
    }).done(function(data, textStatus, jqXHR) {
        params.success({
            rows: data.rows,
            total: data.total,
        });
        if (data.total == 0) {
            var $table = $('#table');
            $table.bootstrapTable('removeAll');
        }
        if (data.exportfn != "") {
            var a = $("<a>").attr("href", proxyPath + "download/" + data.exportfn).text("download");
            $("#exportlink-body").html("");
            $("#exportlink-body").append(a);
            $("#exportlink").modal("show");
        }
    }).fail(function(jqXHR, textStatus, errorThrown) {
        params.error(jqXHR.statusText);                
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });
}

//
// when table is loaded
//
$('#table').on('load-success.bs.table refresh.bs.table', function () {
    //FIXME: perms
    hasPermission("storages", "", "POST").done(function(){
        $(".store").fadeIn();
        localStorage.setItem("storages::POST", true);
    }).fail(function(){
        localStorage.setItem("storages::POST", false);
    })  
    hasPermission("storages", "-2", "GET").done(function(){
        $(".storages").fadeIn();
        localStorage.setItem("storages:-2:GET", true);
    }).fail(function(){
        localStorage.setItem("storages:-2:GET", false);
    }) 
    hasPermission("products", "-1", "PUT").done(function(){
        $(".edit").fadeIn();
        localStorage.setItem("products:-1:PUT", true);
    }).fail(function(){
        localStorage.setItem("products:-1:PUT", false);
    }) 
    hasPermission("products", "-1", "DELETE").done(function(){
        $(".delete").fadeIn();
        localStorage.setItem("products:-1:DELETE", true);
    }).fail(function(){
        localStorage.setItem("products:-1:DELETE", false);
    }) 
});

//
// table row attributes
//
function rowAttributes(row, index) {
    return {"product_id":row["product_id"]}
}

//
// table detail formatter
//
function detailFormatter(index, row) {
    var html = [];

    html.push("<div class='row mt-sm-3'>")
    html.push("<div class='col-sm-12'>")
    $.each(row["synonyms"], function (key, value) {
        html.push("<span>" + value["name_label"] + "</span> ");
    });
    html.push("</div>")
    html.push("</div>")

    html.push("<div class='row mt-sm-3'>")
    html.push("<div class='col-sm-4'><span title='empirical formula' class='mdi mdi-36px mdi-alpha-e-box-outline'></span> " + row["empiricalformula"]["empiricalformula_label"] + "</div>")
    html.push("<div class='col-sm-4'><span title='linear formula' class='mdi mdi-36px mdi-alpha-l-box-outline'></span> " + row["linearformula"]["linearformula_label"]["String"]  + "</div>")
    html.push("<div class='col-sm-4'><span title='3d formula' class='mdi mdi-36px mdi-cube-outline'></span> " + row["product_threedformula"]["String"] + "</div>")
    html.push("</div>")

    html.push("<div class='row mt-sm-4'>")
    html.push("<div class='col-sm-4'><b>ce</b> " + row["cenumber"]["cenumber_label"]["String"] + "</div>")
    html.push("<div class='col-sm-8'><b>msds</b> <a href='" + row["product_msds"]["String"] + "'>" + row["product_msds"]["String"] + "</a></div>")
    html.push("</div>")

    html.push("<div class='row mt-md-3'>")
    html.push("<div class='col-sm-6'>")
    $.each(row["symbols"], function (key, value) {
        html.push("<img src='data:" + value["symbol_image"] + "' alt='" + value["symbol_label"] + "' title='" + value["symbol_label"] + "'/>");
    });
    html.push("</div>")
    html.push("<div class='col-sm-6'>")
    html.push("<p><b>signal word</b>: " + row["signalword"]["signalword_label"]["String"] + "</p>")
    html.push("</div>")
    html.push("</div>")

    html.push("<div class='row mt-md-3'>")
    html.push("<div class='col-sm-6'>")
    html.push("<b>hazardstatements</b>:")
    html.push("<ul>")
    $.each(row["hazardstatements"], function (key, value) {
        html.push("<li>" + value["hazardstatement_reference"] + ": <i>" + value["hazardstatement_label"] + "</i></li>");
    });
    html.push("</ul>")
    html.push("</div>")
    html.push("<div class='col-sm-4'>")
    html.push("<b>precautionarystatements</b>:")
    html.push("<ul>")
    $.each(row["precautionarystatements"], function (key, value) {
        html.push("<li>" + value["precautionarystatement_reference"] + ": " + value["precautionarystatement_label"] + "</li>");
    });
    html.push("</ul>")
    html.push("</div>")
    html.push("</div>")

    html.push("<div class='row mt-md-3'>")
    html.push("<div class='col-sm-6'><b>classofcompound</b> " + row["classofcompound"]["classofcompound_label"]["String"] + "</div>")
    html.push("<div class='col-sm-6'><b>physical state</b>: " + row["physicalstate"]["physicalstate_label"]["String"] + "</div>")
    html.push("</div>")

    html.push("<div class='row mt-md-3'>")
    html.push("<div class='col-sm-6'><b>disposal comment</b>: " + row["product_disposalcomment"]["String"] + "</div>")
    html.push("<div class='col-sm-6'><b>remark</b>: " + row["product_remark"]["String"] + "</div>")
    html.push("</div>")

    html.push("<div class='row mt-sm-3'>")
    html.push("<div class='col-sm-12'>")
    if (row["product_radioactive"]["Bool"]) {
        html.push("<span title='radioactive' class='mdi mdi-36px mdi-radioactive'></span>")
    }
    if (row["product_restricted"]["Bool"]) {
        html.push("<span title='restricted access' class='mdi mdi-36px mdi-hand'></span>")
    }
    html.push("</div>")
    html.push("</div>")

    html.push("<div class='row mt-sm-4'>")
    html.push("<div class='col-sm-12'>" + row["person"]["person_email"] + "</div>")
    html.push("</div>") 
    
    return html.join('');
}

//
// product_specificityFormatter formatter
//
function product_specificityFormatter(value, row, index, field) {
    if (value.Valid) {
        return value.String;
    } else {
        return "";
    }
}
//
// product_slFormatter formatter
//
function product_slFormatter(value, row, index, field) {
    if (value.Valid) {
        return value.String;
    } else {
        return "";
    }
}

//
// table items actions
//
function operateFormatter(value, row, index) {
    // show action buttons if permitted
    pid = row.product_id

    var bookmarkicon = "bookmark-outline"
    if (row.bookmark.bookmark_id.Valid) {
        bookmarkicon = "bookmark"
    }

    // buttons are hidden by default
    var actions = [
    '<button id="storages' + pid + '" class="storages btn btn-primary btn-sm" style="display: none;" title="storages" type="button">',
        '<span class="mdi mdi-24px mdi-cube-unfolded"',
    '</button>',
    '<button id="store' + pid + '" class="store btn btn-primary btn-sm" style="display: none;" title="store" type="button">',
        '<span class="mdi mdi-24px mdi-forklift"',
    '</button>',
    '<button id="edit' + pid + '" class="edit btn btn-primary btn-sm" style="display: none;" title="edit" type="button">',
        '<span class="mdi mdi-24px mdi-border-color"',
    '</button>',
    '<button id="delete' + pid + '" class="delete btn btn-primary btn-sm" style="display: none;" title="delete" type="button">',
        '<span class="mdi mdi-24px mdi-delete"',
    '</button>',
    '<button id="bookmark' + pid + '" class="bookmark btn btn-primary btn-sm" title="(un)bookmark" type="button">',
        '<span id="bookmark' + pid + '" class="mdi mdi-24px mdi-' + bookmarkicon + '">',
    '</button>',
    ];

    if (row.casnumber.casnumber_cmr.Valid) {
        actions.push('<span title="CMR" class="mdi mdi-16px mdi-alert-outline"></span>' + row.casnumber.casnumber_cmr.String);
    }    
    if (row.product_restricted.Valid && row.product_restricted.Bool) {
        actions.push('<span title="restricted access" class="mdi mdi-16px mdi-hand"></span>');
    }  

    return actions.join('&nbsp;');
}

//
// items actions callback
//
window.operateEvents = {
    'click .bookmark': function (e, value, row, index) {
        operateBookmark(e, value, row, index)
    },
    'click .store': function (e, value, row, index) {
        window.location.href = proxyPath + "vc/storages?product=" + row['product_id'];
    },
    'click .storages': function (e, value, row, index) {
        var urlParams = new URLSearchParams(window.location.search);
        var url = proxyPath + "v/storages?product=" + row['product_id'];
        if (urlParams.has("storelocation")) {
            s = urlParams.get("storelocation");
            url = url + "&storelocation=" + s;
        }
        window.location.href = url;
    },
    'click .edit': function (e, value, row, index) {
        operateEdit(e, value, row, index)
    },
    'click .delete': function (e, value, row, index) {
        // hiding possible previous confirmation button
        $(this).confirmation("show").off( "confirmed.bs.confirmation");
        $(this).confirmation("show").off( "canceled.bs.confirmation");
        
        // ask for confirmation and then delete
        $(this).confirmation("show").on( "confirmed.bs.confirmation", function() {
            $.ajax({
                url: proxyPath + "products/" + row['product_id'],
                method: "DELETE",
            }).done(function(data, textStatus, jqXHR) {
                global.displayMessage("product deleted", "success");
                var $table = $('#table');
                $table.bootstrapTable('refresh');
            }).fail(function(jqXHR, textStatus, errorThrown) {
                handleHTTPError(jqXHR.statusText, jqXHR.status)
            });
        }).on( "canceled.bs.confirmation", function() {
        });
    }
};
function operateBookmark(e, value, row, index) {
    // toggling the bookmark
    $.ajax({
        url: proxyPath + "bookmarks/" + row['product_id'],
        method: "PUT",
    }).done(function(data, textStatus, jqXHR) {
        if ($("span#bookmark" + data.product_id).hasClass("mdi-bookmark")) {
            $("span#bookmark" + data.product_id).removeClass("mdi-bookmark");
            $("span#bookmark" + data.product_id).addClass("mdi-bookmark-outline");
        } else {
            $("span#bookmark" + data.product_id).removeClass("mdi-bookmark-outline");
            $("span#bookmark" + data.product_id).addClass("mdi-bookmark");
        }
    }).fail(function(jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });
}
function operateEdit(e, value, row, index) {
    // clearing selections
    $('textarea#product_remark').val(null);
    $('textarea#product_disposalcomment').val(null);
    $('input#product_specificity').val(null);
    $('input#product_msds').val(null);
    $('input#product_threedformula').val(null);

    $('select#casnumber').val(null).trigger('change');
    $('select#casnumber').find('option').remove();
    $('select#cenumber').val(null).trigger('change');
    $('select#cenumber').find('option').remove();
    $('select#name').val(null).trigger('change');
    $('select#name').find('option').remove();
    $('select#symbols').val(null).trigger('change');
    $('select#symbols').find('option').remove();
    $('select#synonyms').val(null).trigger('change');
    $('select#synonyms').find('option').remove();
    $('select#empiricalformula').val(null).trigger('change');
    $('select#empiricalformula').find('option').remove();
    $('select#linearformula').val(null).trigger('change');
    $('select#linearformula').find('option').remove();
    $('select#physicalstate').val(null).trigger('change');
    $('select#physicalstate').find('option').remove();
    $('select#classofcompound').val(null).trigger('change');
    $('select#classofcompound').find('option').remove();
    $('select#symbols').val(null).trigger('change');
    $('select#symbols').find('option').remove();
    $('select#signalword').val(null).trigger('change');
    $('select#signalword').find('option').remove();
    $('select#hazardstatements').val(null).trigger('change');
    $('select#hazardstatements').find('option').remove();
    $('select#precautionarystatements').val(null).trigger('change');
    $('select#precautionarystatements').find('option').remove();

    // getting the product
    $.ajax({
        url: proxyPath + "products/" + row['product_id'],
        method: "GET",
    }).done(function(data, textStatus, jqXHR) {
        // flattening response data
        fdata = flatten(data);

        // processing sqlNull values
        //newfdata = normalizeSqlNull(fdata)
        newfdata = global.normalizeSqlNull(fdata);

        // autofilling form
        $("#edit-collapse").autofill( newfdata, {"findbyname": false } );
        // setting index hidden input
        $("input#index").val(index);
        
        // select2 is not autofilled - we need a special operation
        var newOption = new Option(data.casnumber.casnumber_label, data.casnumber.casnumber_id, true, true);
        $('select#casnumber').append(newOption).trigger('change');
        
        var newOption = new Option(data.empiricalformula.empiricalformula_label, data.empiricalformula.empiricalformula_id, true, true);
        $('select#empiricalformula').append(newOption).trigger('change');
        
        var newOption = new Option(data.name.name_label, data.name.name_id, true, true);
        $('select#name').append(newOption).trigger('change');

        if (data.cenumber.cenumber_id.Valid) {
            var newOption = new Option(data.cenumber.cenumber_label.String, data.cenumber.cenumber_id.Int64, true, true);
            $('select#cenumber').append(newOption).trigger('change');
        }

        if (data.physicalstate.physicalstate_id.Valid) {
            var newOption = new Option(data.physicalstate.physicalstate_label.String, data.physicalstate.physicalstate_id.Int64, true, true);
            $('select#physicalstate').append(newOption).trigger('change');
        }

        if (data.classofcompound.classofcompound_id.Valid) {
            var newOption = new Option(data.classofcompound.classofcompound_label.String, data.classofcompound.classofcompound_id.Int64, true, true);
            $('select#classofcompound').append(newOption).trigger('change');
        }

        if (data.signalword.signalword_id.Valid) {
            var newOption = new Option(data.signalword.signalword_label.String, data.signalword.signalword_id.Int64, true, true);
            $('select#signalword').append(newOption).trigger('change');
        }
       
        for(var i in data.symbols) {
           var newOption = new Option(data.symbols[i].symbol_label, data.symbols[i].symbol_id, true, true);
           $('select#symbols').append(newOption).trigger('change');
        }
        
        for(var i in data.synonyms) {
           var newOption = new Option(data.synonyms[i].name_label, data.synonyms[i].name_id, true, true);
           $('select#synonyms').append(newOption).trigger('change');
        }

        for(var i in data.hazardstatements) {
           var newOption = new Option(data.hazardstatements[i].hazardstatement_reference, data.hazardstatements[i].hazardstatement_id, true, true);
           $('select#hazardstatements').append(newOption).trigger('change');
        }

        for(var i in data.precautionarystatements) {
           var newOption = new Option(data.precautionarystatements[i].precautionarystatement_reference, data.precautionarystatements[i].precautionarystatement_id, true, true);
           $('select#precautionarystatements').append(newOption).trigger('change');
        }
    }).fail(function(jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });

    // finally collapsing the view
    $('#edit-collapse').collapse('show');
    $('#list-collapse').collapse('hide');
    $('div#search').collapse('hide');
    $(".toggleable").hide();
}
