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

    $( "#entity" ).validate({
        errorClass: "alert alert-danger",
        rules: {
            entity_name: {
                required: true,
                remote: {
                    url: "",
                    type: "post",
                    beforeSend: function(jqXhr, settings) {
                        id = -1
                        if ($("form#entity input#entity_id").length) {
                            id = $("form#entity input#entity_id").val()
                        }
                        settings.url = proxyPath + "validate/entity/" + id + "/name/";
                    },
                },               
            },
        },
        messages: {
            entity_name: {
                required: global.t("required_input", container.PersonLanguage)
            }
        }, 
    });

    //
    // SELECT2 SETUP
    //

    // managers select2
    $('select#managers').select2({
        ajax: {
            url: proxyPath + 'people',
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
                obj.text = obj.text || obj.person_email;
                obj.id = obj.id || obj.person_id;
                return obj;
            });
            // getting the number of loaded select elements
            selectnbitems = $("ul#select2-managers li").length + 10;

            return {
                results: newdata,
                pagination: {more: selectnbitems<data.total}
            };
            }
        }
    });

});

//
// TABLE SETUP
//

// table data loading
function getData(params) {
    $.ajax({
        url: proxyPath + "entities",
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

// when table is loaded
$('#table').on('load-success.bs.table refresh.bs.table', function () {
    $("button.storelocations").each(function( index, b ) {
        hasPermission("entities", $(b).attr("eid"), "GET").done(function(){
            $("#storelocations"+$(b).attr("eid")).fadeIn();
            $("#members"+$(b).attr("eid")).fadeIn();
            localStorage.setItem("entities:" + $(b).attr("eid") + ":GET", true);
        }).fail(function(){
            localStorage.setItem("entities:" + $(b).attr("eid") + ":GET", false);
        })
    });
    $("button.storelocations").each(function( index, b ) {
        hasPermission("entities", $(b).attr("eid"), "PUT").done(function(){
            $("#edit"+$(b).attr("eid")).fadeIn();
            localStorage.setItem("entities:" + $(b).attr("eid") + ":PUT", true);
        }).fail(function(){
            localStorage.setItem("entities:" + $(b).attr("eid") + ":PUT", false);
        })
    });
    $("button.storelocations").each(function( index, b ) {
        hasPermission("entities", $(b).attr("eid"), "DELETE").done(function(){
            $("#delete"+$(b).attr("eid")).fadeIn();
            localStorage.setItem("entities:" + $(b).attr("eid") + ":DELETE", true);
        }).fail(function(){
            localStorage.setItem("entities:" + $(b).attr("eid") + ":DELETE", false);
        })
    });
});


//
// TABLE FORMATTERS
//

// managers formatter
function managersFormatter(value, row, index, field) {
    var html = [ "<ul>" ]; 
    $.each(value, function( index, m ) {
        html.push("<li>" + m.person_email + "</li>");
    });
    html.push("</ul>");
    return html.join("");
}

// actions formatter
function operateFormatter(value, row, index) {
    // show action buttons if permitted
    eid = row.entity_id

    // buttons are hidden by default
    var actions = [
    '<button id="storelocations' + eid + '" eid="' + eid + '" class="storelocations btn btn-link btn-sm" style="display: none;" title="' + global.t("storelocations", container.PersonLanguage) + '" type="button">',
        '<span class="mdi mdi-docker mdi-24px"><i>' + row.entity_slc + '</i></span>',
    '</button>',
    '<button id="members' + eid + '" eid="' + eid + '" class="members btn btn-link btn-sm" style="display: none;" title="' + global.t("members", container.PersonLanguage) + '" type="button">',
        '<span class="mdi mdi-account-group mdi-24px"><i>' + row.entity_pc + '</i></span>',
    '</button>',
    '<button id="edit' + eid + '" eid="' + eid + '" class="edit btn btn-link btn-sm" style="display: none;" title="' + global.t("edit", container.PersonLanguage) + '" type="button">',
        '<span class="mdi mdi-border-color mdi-24px"></span>',
    '</button>',
    '<button id="delete' + eid + '" eid="' + eid + '" class="delete btn btn-link btn-sm" style="display: none;" title="' + global.t("delete", container.PersonLanguage) + '" type="button">',
        '<span class="mdi mdi-delete mdi-24px"></span>',
    '</button>'];

    return actions.join('&nbsp;')    
}

//
// TABLE ACTIONS DEFINITION
//
window.operateEvents = {
    'click .storelocations': function (e, value, row, index) {
        window.location.href = proxyPath + "v/storelocations?entity=" + row['entity_id'];
    },
    'click .members': function (e, value, row, index) {
        window.location.href = proxyPath + "v/people?entity=" + row['entity_id'];
    },
    'click .edit': function (e, value, row, index) {
        // clearing selections
        $('select#managers').val(null).trigger('change');
        $('select#managers').find('option').remove();

        // getting the entity
        $.ajax({
            url: proxyPath + "entities/" + row['entity_id'],
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            // flattening response data
            fdata = flatten(data);
            // autofilling form
            $("#edit-collapse").autofill( fdata, {"findbyname": false } );
            // setting index hidden input
            $("input#index").val(index);
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status)
        });

        // getting the entity managers
        $.ajax({
            url: proxyPath + "entities/" + row['entity_id'] + "/people",
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            // select2 is not autofilled - we need a special operation
            for(var i in data) {
                var newOption = new Option(data[i].person_email, data[i].person_id, true, true);
                $('select#managers').append(newOption).trigger('change');
            }
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status)
        });

        // finally collapsing the view
        $('#edit-collapse').collapse('show');
        $('#list-collapse').collapse('hide');
    },
    'click .delete': function (e, value, row, index) {
        // hiding possible previous confirmation button
        $("button#delete" + row.entity_id).confirmation("show").off( "confirmed.bs.confirmation");
        $("button#delete" + row.entity_id).confirmation("show").off( "canceled.bs.confirmation");
        
        // ask for confirmation and then delete
        $("button#delete" + row.entity_id).confirmation("show").on( "confirmed.bs.confirmation", function() {
            $.ajax({
                url: proxyPath + "entities/" + row['entity_id'],
                method: "DELETE",
            }).done(function(data, textStatus, jqXHR) {
                global.displayMessage(global.t("entity_deleted_message", container.PersonLanguage), "success");
                var $table = $('#table');
                $table.bootstrapTable('refresh');
            }).fail(function(jqXHR, textStatus, errorThrown) {
                handleHTTPError(jqXHR.statusText, jqXHR.status)
            });
        }).on( "canceled.bs.confirmation", function() {
        });
    }
};

//
// close buttons actions
//
function closeViewEdit() { $("#list-collapse").collapse("show"); $("#edit-collapse").collapse("hide"); }

//
// save entity callback
//
var createCallBack = function createCallback(data, textStatus, jqXHR) {
    global.displayMessage(global.t("entity_created_message", container.PersonLanguage) + ": " + data.entity_name, "success");
    setTimeout(function(){ window.location = proxyPath + "v/entities"; }, 1000);
}
var updateCallBack = function updateCallback(data, textStatus, jqXHR) {
    $('#list-collapse').collapse('show');
    $('#edit-collapse').collapse('hide');
    var $table = $('#table');
    var index = $('input#index').val();
    $table.bootstrapTable('updateRow', {
        index: index,
        row: {
            "entity_name": data.entity_name,
            "entity_description": data.entity_description,
        }
    });
    $table.bootstrapTable('refresh');
    global.displayMessage(global.t("entity_updated_message", container.PersonLanguage) + ": " + data.entity_name, "success");
}
function saveEntity() {
    var form = $("#entity");
    if (!form.valid()) {
        return;d
    };

    var entity_id = $("input#entity_id").val(),
        entity_name = $("input#entity_name").val(),
        entity_description = $("input#entity_description").val(),
        managers = $('select#managers').select2('data'),
        ajax_url = proxyPath + "entities",
        ajax_method = "POST",
        ajax_callback = createCallBack,
        data = {};

        if ($("form#entity input#entity_id").length) {
            ajax_url = proxyPath + "entities/" + entity_id
            ajax_method = "PUT"
            ajax_callback = updateCallBack
        }

        $.each(managers, function( index, m ) {
            data["managers." + index +".person_id"] = m.id;
            data["managers." + index +".person_email"] = m.text;
        });
        $.extend(data, {
            "entity_id": entity_id,
            "entity_name": entity_name,
            "entity_description": entity_description,
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
}
