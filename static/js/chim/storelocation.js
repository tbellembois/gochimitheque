$( document ).ready(function() { 

    //
    // INITIAL SETUP
    //

    // populating search input if needed
    if (URLValues.has("search")) {
        $('#table').bootstrapTable('resetSearch', URLValues.get("search")[0]);
    }

    //
    // FORM VALIDATION
    //

    $( "#storelocation" ).validate({
        // ignore required to validate select2
        ignore: "",
        errorClass: "alert alert-danger",
        rules: {
            storelocation_name: {
                required: true,
            },
            entity: {
                required: true,
            },
        },
        messages: {
            storelocation_name: {
                required: global.t("required_input", container.PersonLanguage)
            },
            entity: {
                required: global.t("required_input", container.PersonLanguage)
            },
        }, 
    });

    //
    // SELECT2 SETUP
    //
    
    // entity select2
    $('select#entity').select2({
        ajax: {
            url: proxyPath + 'entities',
            dataType: 'json',
            data: function (params) {
                var query = {
                    search: params.term,
                    page: params.page || 1,
                    offset: (params.page-1)*10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            processResults: function (data) {
            // replacing email by text expected by select2
            var newdata = $.map(data.rows, function (obj) {
                obj.text = obj.text || obj.entity_name;
                obj.id = obj.id || obj.entity_id;
                return obj;
            });
            // getting the number of loaded select elements
            selectnbitems = $("ul#select2-entity li").length + 10;

            return {
                results: newdata,
                pagination: {more: selectnbitems<data.total}
            };
            }
        }
    });
    $('select#entity').on('select2:select', function (e) {
        $('select#storelocation').val(null).trigger('change');
        $('select#storelocation').find('option').remove();
    });

    // storelocation select2
    $('select#storelocation').select2({
        placeholder: "select an entity first",
        ajax: {
            url: proxyPath + 'storelocations',
            dataType: 'json',
            data: function (params) {
                eid = $('select#entity').select2('data')[0].id;
                var query = {
                    entity: eid,
                    search: params.term,
                    page: params.page || 1,
                    offset: (params.page-1)*10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            processResults: function (data) {
            // replacing email by text expected by select2
            var newdata = $.map(data.rows, function (obj) {
                obj.text = obj.text || obj.storelocation_name.String;
                obj.id = obj.id || obj.storelocation_id.Int64;
                return obj;
            });
            // getting the number of loaded select elements
            selectnbitems = $("ul#select2-storelocation li").length + 10;

            return {
                results: newdata,
                pagination: {more: selectnbitems<data.total}
            };
            }
        }
    });

    // color picker
    $("#storelocation_color").colorpicker();

});

//
// TABLE SETUP
//

// table data loading
function getData(params) {
    $.ajax({
        url: proxyPath + "storelocations",
        method: "GET",
        dataType: "JSON",
        data: params.data,
    }).done(function(data, textStatus, jqXHR) {
        params.success({
            rows: data.rows,
            total: data.total,
        });
    }).fail(function(jqXHR, textStatus, errorThrown) {
        params.error(jqXHR.statusText);                
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });
}
function queryParams(params) {
    // getting request parameters
    var urlParams = new URLSearchParams(window.location.search);

    if (urlParams.has("entity")) {
        params["entity"] = urlParams.get("entity");
    }
    return params;
}

// when table is loaded
$('#table').on('load-success.bs.table refresh.bs.table', function () {
    $("button.edit").each(function( index, b ) {
        hasPermission("storelocations", $(b).attr("slid"), "PUT").done(function(){
            $("#edit"+$(b).attr("slid")).fadeIn();
            localStorage.setItem("storelocations:" + $(b).attr("slid") + ":PUT", true);
        }).fail(function(){
            localStorage.setItem("storelocations:" + $(b).attr("slid") + ":PUT", false);
        })
    });
    $("button.edit").each(function( index, b ) {
        hasPermission("storelocations", $(b).attr("slid"), "DELETE").done(function(){
            $("#delete"+$(b).attr("slid")).fadeIn();
            localStorage.setItem("storelocations:" + $(b).attr("slid") + ":DELETE", true);
        }).fail(function(){
            localStorage.setItem("storelocations:" + $(b).attr("slid") + ":DELETE", false);
        })
    });
});

//
// TABLE FORMATTERS
//

// storelocation_idFormatter formatter
function storelocation_idFormatter(value, row, index, field) {
    if (value.Valid) {
        return value.Int64;
    } else {
        return "";
    }
}
// storelocation_colorFormatter formatter
function storelocation_colorFormatter(value, row, index, field) {
    if (value.Valid) {
        return '<div style="background-color: ' + value.String + '">&nbsp;</div>';
    } else {
        return "";
    }
}
// storelocation_canstoreFormatter formatter
function storelocation_canstoreFormatter(value, row, index, field) {
    if (value.Valid && value.Bool) {
        return '<span class="mdi mdi-check mdi-36px"></span>';
    } else {
        return '<span class="mdi mdi-close mdi-36px"></span>';
    }
}
// storelocationFormatter formatter
function storelocationFormatter(value, row, index, field) {
    if (value.storelocation_name.Valid) {
        return value.storelocation_name.String;
    } else {
        return "";
    }
}

// actions formatter
function operateFormatter(value, row, index) {
    // show action buttons if permitted
    slid = row.storelocation_id.Int64

    // buttons are hidden by default
    var actions = [
    '<button id="edit' + slid + '" slid="' + slid + '" class="edit btn btn-link btn-sm" style="display: none;" title="edit" type="button">',
        '<span class="mdi mdi-24px mdi-border-color"',
    '</button>',
    '<button id="delete' + slid + '" slid="' + slid + '" class="delete btn btn-link btn-sm" style="display: none;" title="delete" type="button">',
        '<span class="mdi mdi-24px mdi-delete"',
    '</button>'];

    return actions.join('&nbsp;')    
}

//
// TABLE ACTIONS DEFINITION
//

window.operateEvents = {
    'click .edit': function (e, value, row, index) {
        operateEdit(e, value, row, index)
    },
    'click .delete': function (e, value, row, index) {
        // hiding possible previous confirmation button
        $("button#delete" + row.storelocation_id.Int64).confirmation("show").off( "confirmed.bs.confirmation");
        $("button#delete" + row.storelocation_id.Int64).confirmation("show").off( "canceled.bs.confirmation");
        
        // ask for confirmation and then delete
        $("button#delete" + row.storelocation_id.Int64).confirmation("show").on( "confirmed.bs.confirmation", function() {
            $.ajax({
                url: proxyPath + "storelocations/" + row['storelocation_id'],
                method: "DELETE",
            }).done(function(data, textStatus, jqXHR) {
                global.displayMessage("store location deleted", "success");
                var $table = $('#table');
                $table.bootstrapTable('refresh');
            }).fail(function(jqXHR, textStatus, errorThrown) {
                handleHTTPError(jqXHR.statusText, jqXHR.status)
            });
        }).on( "canceled.bs.confirmation", function() {
        });
    }
};
function operateEdit(e, value, row, index) {
    // clearing selections
    $('select#entity').val(null).trigger('change');
    $('select#entity').find('option').remove();
    $('select#storelocation').val(null).trigger('change');
    $('select#storelocation').find('option').remove();

    // getting the store location
    $.ajax({
        url: proxyPath + "storelocations/" + row['storelocation_id'].Int64,
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
        var newOption = new Option(data.entity.entity_name, data.entity.entity_id, true, true);
        $('select#entity').append(newOption).trigger('change');
        if (data.storelocation.storelocation_id.Valid) {
            var newOption = new Option(data.storelocation.storelocation_name.String, data.storelocation.storelocation_id.Int64, true, true);
            $('select#storelocation').append(newOption).trigger('change');
        }
    }).fail(function(jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });

    // finally collapsing the view
    $('#edit-collapse').collapse('show');
    $('#list-collapse').collapse('hide');
};

//
// close buttons actions
//
function closeEdit() { $("#list-collapse").collapse("show"); $("#edit-collapse").collapse("hide"); }
        
//
// save store location callback
//
var createCallBack = function createCallback(data, textStatus, jqXHR) {
    global.displayMessage(global.t("storelocation_created_message", container.PersonLanguage) + ": " + data.storelocation_name.String, "success");
    setTimeout(function(){ window.location = proxyPath + "v/storelocations"; }, 1000);
}
var updateCallBack = function updateCallback(data, textStatus, jqXHR) {
    $('#list-collapse').collapse('show');
    $('#edit-collapse').collapse('hide');
    var $table = $('#table');
    var index = $('input#index').val();
    $table.bootstrapTable('updateRow', {
        index: index,
        row: {
            "storelocation_name": data.storelocation_name,
            "storelocation_canstore": data.storelocation_canstore,
            "storelocation_color": data.storelocation_color,
            "entity.entity_name": data.entity.entity_name,
        }
    });
    $table.bootstrapTable('refresh');
    global.displayMessage(global.t("storelocation_updated_message", container.PersonLanguage) + ": " + data.storelocation_name.String, "success");
}
function saveStoreLocation() {
    var form = $("#storelocation");
    if (! form.valid()) {
        return;
    };

    var storelocation_id = $("input#storelocation_id").val(),
        storelocation_name = $("input#storelocation_name").val(),
        storelocation_color = $("input#storelocation_color").val(),
        storelocation_canstore = $("input#storelocation_canstore:CHECKED").val(),
        entity = $('select#entity').select2('data')[0],
        storelocation = $('select#storelocation').select2('data')[0],
        ajax_url = proxyPath + "storelocations",
        ajax_method = "POST",
        ajax_callback = createCallBack,
        data = {};

        if (storelocation !== undefined) {
            $.extend(data, {
                "storelocation.storelocation.storelocation_id": storelocation.id,
                "storelocation.storelocation.storelocation_name": storelocation.text,
            });                    
        }

        if ($("form#storelocation input#storelocation_id").length) {
            ajax_url = proxyPath + "storelocations/" + storelocation_id
            ajax_method = "PUT"
            ajax_callback = updateCallBack
        }

        $.extend(data, {
            "storelocation_id": storelocation_id,
            "storelocation_name": storelocation_name,
            "storelocation_color": storelocation_color,
            "storelocation_canstore": storelocation_canstore == "on" ? true : false,
            "entity.entity_id": entity.id,
            "entity.entity_name": entity.text,
        });

        // lazily clearing all the cache storage
        localStorage.clear();
        $.ajax({
            url: ajax_url,
            method: ajax_method,
            dataType: 'json',
            data: data,
        }).done(ajax_callback).fail(function(jqXHR, textStatus, errorThrown) {           
            handleHTTPError(jqXHR.statusText, jqXHR.status)
        });  
    };
        



 


 