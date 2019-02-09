//
// request performed at table data loading
//
function getData(params) {
    // getting request parameters
    var urlParams = new URLSearchParams(window.location.search);
    
    // saving the query parameters
    lastQueryParams = params;
    $.ajax({
        url: proxyPath + "storages",
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
        handleHTTPError(jqXHR.statusText, jqXHR.status);
    });
}

//
// view storages button
//
$('#s_storage_archive_button').on('click', function () {
    var $table = $('#table'),
    btntitle = "",
    btnicon = "";
    if ($('#s_storage_archive_button').hasClass("active")) {
        updateQueryStringParam("storage_archive", "false");
        btntitle = "show deleted";
        btnicon = "delete";
    } else {
        updateQueryStringParam("storage_archive", "true");
        btntitle = "do not show deleted";       
        btnicon = "delete-forever";
    }
    $table.bootstrapTable('refresh');
    $('#s_storage_archive_button').attr("title", btntitle);
    $('#s_storage_archive_button > span').attr("class", "mdi mdi-24px mdi-"+btnicon);
});

//
// stock
//
function showStockRecursive(sl, depth) {
    // pre checking if there is a stock or not for sl
    var hasStock = false;
    for (var i in sl.stock) { 
        var stock = sl.stock[i];
        if (stock.total !== 0 || stock.current !== 0) {
            hasStock = true;
            break;
        }
    }
    
    if (hasStock) {
        var html = [("<div class='row mt-sm-3'>")];
        for (i=1; i<=depth; i++) {
            html.push("<div class='col-sm-1'>&nbsp;</div>");
        }
        html.push("<div class='col' style='color: " + sl.storelocation_color.String + "'>" + sl.storelocation_name.String + "</div>");
        
        for (var i in sl.stock) {
            var stock = sl.stock[i];
            
            if (stock.total === 0 && stock.current === 0) {
            } else {
                html.push("<div class='col'><i title='total stock including sub store locations' class='material-icons'>functions</i> " + stock.total + " <b>" + stock.unit.unit_label.String + "</b></div>");
                html.push("<div class='col'><i title='total in this store location' class='material-icons'>extensions</i> " + stock.current + " <b>" + stock.unit.unit_label.String + "</b></div>");
                
            }
        }
        
        $("#stock-body").append(html.join(""));
        $("#stock-body").append("</div>");
    }
    
    if (sl.children !== null) {
        depth++;
        for  (var key in sl.children) {
            showStockRecursive(sl.children[key], depth);
        }
    }
}
function showStock(pid) {
    $.ajax({
        url: proxyPath + "stocks/" + pid,
        method: "GET",
    }).done(function(data, textStatus, jqXHR) {
        $("#stock-body").html("");
        for (var key in data) {
            showStockRecursive(data[key], 0)
        }
    }).fail(function(jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status);
    });
}

//
// table row attributes
//
function rowAttributes(row, index) {
    return {"storage_id":row["storage_id"].Int64}
}
//
// table detail formatter
//
function detailFormatter(index, row) {
    var html = [];
    
    html.push("<div class='row'>")
    
    html.push("<div class='col-sm-3'>")
    if ( row["storage_qrcode"] != null ) {
        html.push("<img src='data:image/png;base64," + row["storage_qrcode"] + "'>")
    }
    html.push("</div>")
    
    html.push("<div class='col-sm-9'>")
    html.push("<div class='row mt-sm-3'>")
    html.push("<div class='col-sm-6'><i class='material-icons'>local_offer</i> " + row["product"]["name"]["name_label"] + "</p></div>")
    html.push("<div class='col-sm-6'><i class='material-icons'>extension</i> " + row["storelocation"]["storelocation_name"]["String"] + "</p></div>")
    html.push("</div>")            
    
    html.push("<div class='row mt-sm-3'>")
    html.push("<div class='col-sm-3'><p><b>quantity</b> " + row["storage_quantity"]["Float64"] + " " + row["unit"]["unit_label"]["String"] + "</p></div>")
    html.push("<div class='col-sm-3'><b>barecode</b> " + row["storage_barecode"]["String"] + "</p></div>")
    html.push("</div>")   
    
    html.push("<div class='row mt-sm-2'>")
    html.push("<div class='col-sm-3'><b>batch number</b> " + row["storage_batchnumber"]["String"] + "</p></div>")
    html.push("<div class='col-sm-3'><p><b>supplier</b> " + row["supplier"]["supplier_label"]["String"] + "</p></div>")
    html.push("</div>")    
    
    html.push("<div class='row mt-sm-2'>")
    html.push("<div class='col-sm-12'><b>entry date</b> " + dateFormatter(row["storage_entrydate"], null, null, null) + "</p></div>")
    html.push("<div class='col-sm-12'><b>exit date</b> " + dateFormatter(row["storage_exitdate"], null, null, null) + "</p></div>")
    html.push("<div class='col-sm-12'><b>opening date</b> " + dateFormatter(row["storage_openingdate"], null, null, null) + "</p></div>")
    html.push("<div class='col-sm-12'><b>expiration date</b> " + dateFormatter(row["storage_expirationdate"], null, null, null) + "</p></div>")
    html.push("</div>")   
    
    html.push("<div class='row mt-sm-4'>")
    html.push("<div class='col-sm-8'>" + dateFormatter(row["storage_creationdate"], null, null, null) + " <i>(" + dateFormatter(row["storage_modificationdate"], null, null, null) + ")</i></div>")
    html.push("<div class='col-sm-4'>" + row["person"]["person_email"] + "</div>")
    html.push("</div>")   
    
    html.push("<div class='row'>")
    if (row["storage_todestroy"]["Bool"]) {
        html.push("<div class='col-sm-12'><i title='to destroy' class='material-icons'>delete_outline</i></div>")
    }
    html.push("</div>")             
    html.push("</div>")
    
    html.push("</div>")  
    
    return html.join('');
}
//
// date formatter
//
function dateFormatter(value, row, index, field) {
    if (value.Valid != undefined && value.Valid == false) {
        return ""
    } else {
        date = new Date(value);
        return date.toLocaleString();
    }
}
//
// storage_idFormatter
//
function storage_idFormatter(value, row, index, field) {
    if (row.storage_id.Valid) {
        return row.storage_id.Int64;
    } else {
        return value;
    }            
}
//
// storelocation_name formatter
//
function storelocation_fullpathFormatter(value, row, index, field) {
    if (row.storelocation.storelocation_color.Valid) {
        return "<span style='color:" + row.storelocation.storelocation_color.String + ";'>" + value + "</span>";
    } else {
        return "<span>" + value + "</span>";
    }
}
//
// storage_barecode formatter
//
function storage_barecodeFormatter(value, row, index, field) {
    if (row.storage_barecode.Valid) {
        return row.storage_barecode.String;
    } else {
        return "";
    }
}

//
// table items actions
//
function operateFormatter(value, row, index) {
    // show action buttons if permitted
    sid = row.storage_id.Int64
    eid = row.storelocation.entity.entity_id
    
    var borrowingicon = "hand-okay",
    borrowingtitle = "borrow";
    if (row.borrowing.borrowing_id.Valid) {
        borrowingicon = "hand-pointing-left";
        borrowingtitle = "unborrow"
    }
    
    if (row.storage_archive.Valid && row.storage_archive.Bool) {
        // this is an archive
        var actions = [
            '<button id="clone' + sid + '" class="clone btn btn-primary btn-sm" style="display: none;" title="clone" type="button">',
            '<span class="mdi mdi-24px mdi-content-copy"></span>',
            '</button>',
            '<button id="restore' + sid + '" class="restore btn btn-primary btn-sm" style="display: none;" title="restore" type="button">',
            '<span class="mdi mdi-24px mdi-undo"></span>',
            '</button>',
            '<button id="delete' + sid + '" class="delete btn btn-primary btn-sm" style="display: none;" title="delete" type="button">',
            '<span class="mdi mdi-24px mdi-delete"></span>',
            '</button>']; 
    } else if (row.storage.storage_id.Valid) {
        // this is an history
        var actions = [
            '<button id="clone' + sid + '" class="clone btn btn-primary btn-sm" style="display: none;" title="clone" type="button">',
            '<span class="mdi mdi-24px mdi-content-copy"></span>',
            '</button>'];
    } else {
        // buttons are hidden by default
        var actions = [
            '<button id="edit' + sid + '" class="edit btn btn-primary btn-sm" style="display: none;" title="edit" type="button">',
            '<span class="mdi mdi-24px mdi-border-color"></span>',
            '</button>',
            '<button id="clone' + sid + '" class="clone btn btn-primary btn-sm" style="display: none;" title="clone" type="button">',
            '<span class="mdi mdi-24px mdi-content-copy"></span>',
            '</button>',
            '<button id="archive' + sid + '" class="archive btn btn-primary btn-sm" style="display: none;" title="delete" type="button">',
            '<span class="mdi mdi-24px mdi-delete"></span>',
            '</button>',
            '<button id="borrow' + sid + '" data-target="#borrow" class="borrow btn btn-primary btn-sm" style="display: none;" title="' + borrowingtitle + '" type="button">',
            '<span class="mdi mdi-24px mdi-' + borrowingicon + '"></span>',
            '</button>'];
    }
                
    if (row.storage.storage_id.Valid) {
        // this is an history
        actions.push('<span class="mdi mdi-24px mdi-history"></span>');
    }
    if (row.storage_archive.Valid && row.storage_archive.Bool) {
        // this is an archive
        actions.push('<span class="mdi mdi-24px mdi-delete"></span>');
    }
    
    if (row.storage_creationdate != row.storage_modificationdate && !row.storage.storage_id.Valid) {
        actions.push('<button id="history' + sid + '" class="history btn btn-primary btn-sm" title="show history" type="button">');
        actions.push('<span class="mdi mdi-24px mdi-history"></span>');
        actions.push('</button>');
    }
    
    hasPermission(proxyPath + "f/storages/" + eid, "PUT", sid).done(function(){
        $("#edit"+this.itemId).fadeIn();
        $("#clone"+this.itemId).fadeIn();
        $("#borrow"+this.itemId).fadeIn();
    })
    hasPermission(proxyPath + "f/storages/" + eid, "DELETE", sid).done(function(){
        $("#delete"+this.itemId).fadeIn();
        $("#archive"+this.itemId).fadeIn();
        $("#restore"+this.itemId).fadeIn();
    })
    
    return actions.join('&nbsp;')    
}
            
// items actions callback
function operateBorrow(e, value, row, index) {
    $('select#borrower').val(null).trigger('change');
    $('select#borrower').find('option').remove();
    
    $("input#bstorage_id").val(row['storage_id'].Int64);
    
    if (row['borrowing']['borrowing_id'].Valid) {
        saveBorrowing();
    } else {
        $("#borrow").modal("show");
    }
    
    var $table = $('#table');
    $table.bootstrapTable('refresh');
}
function operateEdit(e, value, row, index) {
    // clearing selections
    $('input#storage_quantity').val(null);
    $('input#storage_entrydate').val(null);
    $('input#storage_exitdate').val(null);
    $('input#storage_openingdate').val(null);
    $('input#storage_expirationdate').val(null);
    $('input#storage_reference').val(null);
    $('input#storage_batchnumber').val(null);
    $('input#storage_barecode').val(null);
    $('input#storage_comment').val(null);
    
    $('select#storelocation').val(null).trigger('change');
    $('select#storelocation').find('option').remove();
    $('select#unit').val(null).trigger('change');
    $('select#unit').find('option').remove();
    $('select#supplier').val(null).trigger('change');
    $('select#supplier').find('option').remove();
    
    // getting the storage
    $.ajax({
        url: proxyPath + "storages/" + row['storage_id'].Int64,
        method: "GET",
    }).done(function(data, textStatus, jqXHR) {
        // flattening response data
        fdata = flatten(data);
        
        // processing sqlNull values
        newfdata = global.normalizeSqlNull(fdata);
        
        // autofilling form
        $("#edit-collapse").autofill( newfdata, {"findbyname": false } );
        
        // setting index hidden input
        $("input#index").val(index);
        
        // select2 is not autofilled - we need a special operation
        var newOption = new Option(data.storelocation.storelocation_name.String, data.storelocation.storelocation_id.Int64, true, true);
        $('select#storelocation').append(newOption).trigger('change');
        
        if (data.unit.unit_id.Valid) {
            var newOption = new Option(data.unit.unit_label.String, data.unit.unit_id.Int64, true, true);
            $('select#unit').append(newOption).trigger('change');
        }
        
        if (data.supplier.supplier_id.Valid) {
            var newOption = new Option(data.supplier.supplier_label.String, data.supplier.supplier_id.Int64, true, true);
            $('select#supplier').append(newOption).trigger('change');
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
window.operateEvents = {
    'click .edit': function (e, value, row, index) {
        operateEdit(e, value, row, index)
    },
    'click .borrow': function (e, value, row, index) {
        operateBorrow(e, value, row, index)
    },
    'click .history': function (e, value, row, index) {
        var urlParams = new URLSearchParams(window.location.search);
        window.location = proxyPath + "v/storages?storage="+row['storage_id'].Int64+"&history=true&" + urlParams;
    },
    'click .clone': function (e, value, row, index) {
        window.location = proxyPath + "vc/storages?storage="+row['storage_id'].Int64+"&product="+row['product']['product_id'];
    },
    'click .restore': function (e, value, row, index) {
        // hiding possible previous confirmation button
        $(this).confirmation("show").off( "confirmed.bs.confirmation");
        $(this).confirmation("show").off( "canceled.bs.confirmation");
        
        // ask for confirmation and then restore
        $(this).confirmation("show").on( "confirmed.bs.confirmation", function() {
            $.ajax({
                url: proxyPath + "storages/" + row['storage_id'].Int64 + "/r",
                method: "PUT",
            }).done(function(data, textStatus, jqXHR) {
                global.displayMessage("storage restored", "success");
                var $table = $('#table');
                $table.bootstrapTable('refresh');
            }).fail(function(jqXHR, textStatus, errorThrown) {
                handleHTTPError(jqXHR.statusText, jqXHR.status)
            });
        }).on( "canceled.bs.confirmation", function() {
        });
    },
    'click .archive': function (e, value, row, index) {
        // hiding possible previous confirmation button
        $(this).confirmation("show").off( "confirmed.bs.confirmation");
        $(this).confirmation("show").off( "canceled.bs.confirmation");
        
        // ask for confirmation and then delete
        $(this).confirmation("show").on( "confirmed.bs.confirmation", function() {
            $.ajax({
                url: proxyPath + "storages/" + row['storage_id'].Int64 + "/a",
                method: "DELETE",
            }).done(function(data, textStatus, jqXHR) {
                global.displayMessage("storage trashed", "success");
                var $table = $('#table');
                $table.bootstrapTable('refresh');
            }).fail(function(jqXHR, textStatus, errorThrown) {
                handleHTTPError(jqXHR.statusText, jqXHR.status)
            });
        }).on( "canceled.bs.confirmation", function() {
        });
    },
    'click .delete': function (e, value, row, index) {
        // hiding possible previous confirmation button
        $(this).confirmation("show").off( "confirmed.bs.confirmation");
        $(this).confirmation("show").off( "canceled.bs.confirmation");
        
        // ask for confirmation and then delete
        $(this).confirmation("show").on( "confirmed.bs.confirmation", function() {
            $.ajax({
                url: proxyPath + "storages/" + row['storage_id'].Int64,
                method: "DELETE",
            }).done(function(data, textStatus, jqXHR) {
                global.displayMessage("storage deleted", "success");
                var $table = $('#table');
                $table.bootstrapTable('refresh');
            }).fail(function(jqXHR, textStatus, errorThrown) {
                handleHTTPError(jqXHR.statusText, jqXHR.status)
            });
        }).on( "canceled.bs.confirmation", function() {
        });
    }
};