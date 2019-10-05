//
// request performed at table data loading
//
$( document ).ready(function() {
    //
    // form validation
    //
    $( "#storage" ).validate({
        errorClass: "alert alert-danger",
        rules: {
            storelocation: {
                required: true,
            },
        },
        messages: {
            storelocation: {
                required: global.t("required_input", container.PersonLanguage)
            },
        },
    });
    $( "#borrowing" ).validate({
        errorClass: "alert alert-danger",
        rules: {
            borrower: {
                required: true,
            },
        },
        messages: {
            borrower: {
                required: global.t("required_input", container.PersonLanguage)
            },
        },
    });

    //
    // search form
    //
    $('select#s_storelocation').select2({
        templateResult: formatStorelocation,
        allowClear: true,
        placeholder: "store location",
        ajax: {
            url: proxyPath + 'storelocations',
            delay: 400,
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
            dataType: 'json',
            processResults: function (data) {
                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.storelocation_fullpath;
                    obj.id = obj.id || obj.storelocation_id.Int64;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-s_storelocation-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {more: selectnbitems<data.total}
                };
            }
        }
    });

    $('select#s_casnumber').select2({
        tags: false,
        allowClear: true,
        placeholder: "select a cas number",
        ajax: {
            url: proxyPath + 'products/casnumbers/',
            delay: 400,
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
            dataType: 'json',
            processResults: function (data) {
                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.casnumber_label;
                    obj.id = obj.id || obj.casnumber_id;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-s_casnumber-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {more: selectnbitems<data.total}
                };
            }
        }
    });

    $('select#s_name').select2({
        tags: false,
        allowClear: true,
        placeholder: "select a name",
        ajax: {
            url: proxyPath + 'products/names/',
            delay: 400,
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
            dataType: 'json',
            processResults: function (data) {

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.name_label;
                    obj.id = obj.id || obj.name_id;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-s_name-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {more: selectnbitems<data.total}
                };
            }
        }
    });

    $('select#s_empiricalformula').select2({
        tags: false,
        allowClear: true,
        placeholder: "select a formula",
        ajax: {
            url: proxyPath + 'products/empiricalformulas/',
            delay: 400,
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
            dataType: 'json',
            processResults: function (data) {

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.empiricalformula_label;
                    obj.id = obj.id || obj.empiricalformula_id;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-s_empiricalformula-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {more: selectnbitems<data.total}
                };
            }
        }
    });

    $('select#s_signalword').select2({
        templateResult: formatSignalWord,
        allowClear: true,
        placeholder: "select signal word",
        ajax: {
            url: proxyPath + 'products/signalwords/',
            delay: 400,
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
            dataType: 'json',
            processResults: function (data) {
                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.signalword_label.String;
                    obj.id = obj.id || obj.signalword_id.Int64;
                    return obj;
                });
                    return {
                    results: newdata,
                };
            }
        }
    });

    $('select#s_symbols').select2({
        templateResult: formatSymbol,
        closeOnSelect: false,
        ajax: {
            url: proxyPath + 'products/symbols/',
            delay: 400,
            dataType: 'json',
            processResults: function (data) {
                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.symbol_label;
                    obj.id = obj.id || obj.symbol_id;
                    return obj;
                });
                    return {
                    results: newdata,
                };
            }
        }
    });

    $('select#s_hazardstatements').select2({
        templateResult: formatHazardStatement,
        closeOnSelect: false,
        ajax: {
            url: proxyPath + 'products/hazardstatements/',
            delay: 400,
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
            dataType: 'json',
            processResults: function (data) {
                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.hazardstatement_reference;
                    obj.id = obj.id || obj.hazardstatement_id;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-hazardstatements-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {more: selectnbitems<data.total}
                };
            }
        }
    });

    $('select#s_precautionarystatements').select2({
        templateResult: formatPrecautionaryStatement,
        closeOnSelect: false,
        ajax: {
            url: proxyPath + 'products/precautionarystatements/',
            delay: 400,
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
            dataType: 'json',
            processResults: function (data) {
                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.precautionarystatement_reference;
                    obj.id = obj.id || obj.precautionarystatement_id;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-precautionarystatements-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {more: selectnbitems<data.total}
                };
            }
        }
    });

    //
    // store locations selector select2
    //
    // $('select#storelocationselector').select2({
    //     templateResult: formatStorelocation,
    //     //placeholder: "direct store location access",
    //     ajax: {
    //         url: proxyPath + 'storelocations',
    //         data: function (params) {
    //             var query = {
    //                 search: params.term,
    //                 page: params.page || 1,
    //                 offset: (params.page-1)*10 || 0,
    //                 limit: 10
    //             }

    //             // Query parameters will be ?search=[term]&page=[page]
    //             return query;
    //         },
    //         dataType: 'json',
    //         processResults: function (data) {
    //             // replacing name by text expected by select2
    //             var newdata = $.map(data.rows, function (obj) {
    //                 obj.text = obj.text || obj.storelocation_fullpath;
    //                 obj.id = obj.id || obj.storelocation_id.Int64;
    //                 return obj;
    //             });
    //             // getting the number of loaded select elements
    //             selectnbitems = $("ul#select2-storelocationselector-results li").length + 10;
                
    //             return {
    //                 results: newdata,
    //                 pagination: {more: selectnbitems<data.total}
    //             };
    //         }
    //     }
    // }).on("select2:select", function (e) {
    //     var data = e.params.data;
    //     var slid = data.storelocation_id.Int64;
    //     window.location.href = proxyPath + "v/products?storelocation=" + slid;
    // });

    //
    // store location select2
    //
    $('select#storelocation').select2({
        templateResult: formatStorelocation,
        placeholder: "store location",
        ajax: {
            url: proxyPath + 'storelocations',
            delay: 400,
            data: function (params) {
                var query = {
                    storelocation_canstore: true,
                    permission: "w",
                    search: params.term,
                    page: params.page || 1,
                    offset: (params.page-1)*10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data) {
                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.storelocation_name.String;
                    obj.id = obj.id || obj.storelocation_id.Int64;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-storelocation-results li").length + 10;
                
                return {
                    results: newdata,
                    pagination: {more: selectnbitems<data.total}
                };
            }
        }
    }).on("change", function (e) {
        $(this).valid(); // FIXME: see https://github.com/select2/select2/issues/3901
    });

    //
    // unit select2
    //
    $('select#unit').select2({
        allowClear: true,
        placeholder: "select a unit",
        ajax: {
            url: proxyPath + 'storages/units',
            delay: 400,
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
            dataType: 'json',
            processResults: function (data) {
                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.unit_label.String;
                    obj.id = obj.id || obj.unit_id.Int64;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-unit-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {more: selectnbitems<data.total}
                };
            }
        }
    }).on("change", function (e) {
        $(this).valid(); // FIXME: see https://github.com/select2/select2/issues/3901
    });

    //
    // suppliers select2
    //
    $('select#supplier').select2({
        allowClear: true,
        tags: true,
        placeholder: "select or enter a supplier",
        createTag: function (params) {
            if ($("input#exactMatchSupplier").val() == "true") {
                return null
            }
            return {
                id: params.term,
                text: params.term,
            }
        },
        ajax: {
            url: proxyPath + 'storages/suppliers',
            delay: 400,
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
            dataType: 'json',
            processResults: function (data) {
                isExactMatch=false;
                
                // looking for an exact match
                $.each(data.rows, function( index, value ) {
                    if(value.c == 1) {
                        isExactMatch=true;
                    }
                });
                
                // there is a match: setting the input field
                if (isExactMatch) {
                    $("input#exactMatchSupplier").val("true");
                } else {
                    $("input#exactMatchSupplier").val("false");
                }

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.supplier_label.String;
                    obj.id = obj.id || obj.supplier_id.Int64;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-supplier-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {more: selectnbitems<data.total}
                };
            }
        }
    }).on("change", function (e) {
        $(this).valid(); // FIXME: see https://github.com/select2/select2/issues/3901
    });

    //
    // borrowers select2
    //
    $('select#borrower').select2({
        placeholder: "select a borrower",
        ajax: {
            url: proxyPath + 'people',
            dataType: 'json',
            data: function (params) {
                var query = {
                order: "asc",
                sort: "person_email",
                }
                return query;
            },
            processResults: function (data) {
            // replacing email by text expected by select2
            var newdata = $.map(data.rows, function (obj) {
                obj.text = obj.text || obj.person_email;
                obj.id = obj.id || obj.person_id;
                return obj;
            });
            return {
                results: newdata
            };
            }
        }
    });
});

function getData(params) {
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
            var a = $("<a>").attr("href", proxyPath + "download/" + data.exportfn).html("<span class='mdi mdi-48px mdi-file-download'></span>");
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
// when table is loaded
//
$('#table').on('load-success.bs.table refresh.bs.table', function () {
    $("table#table").find("tr").each(function( index, b ) {
        hasPermission("storages", $(b).attr("storage_id"), "PUT").done(function(){
            $("#edit"+$(b).attr("storage_id")).fadeIn();
            $("#clone"+$(b).attr("storage_id")).fadeIn();
            $("#borrow"+$(b).attr("storage_id")).fadeIn();
            localStorage.setItem("storages:" + $(b).attr("storage_id") + ":PUT", true);
        }).fail(function(){
            localStorage.setItem("storages:" + $(b).attr("storage_id") + ":PUT", false);
        })
        hasPermission("storages", $(b).attr("storage_id"), "DELETE").done(function(){
            $("#delete"+$(b).attr("storage_id")).fadeIn();
            $("#archive"+$(b).attr("storage_id")).fadeIn();
            $("#restore"+$(b).attr("storage_id")).fadeIn();
            localStorage.setItem("storages:" + $(b).attr("storage_id") + ":DELETE", true);
        }).fail(function(){
            localStorage.setItem("storages:" + $(b).attr("storage_id") + ":DELETE", false);
        })
    });
});

//
// view storages button
//
$('#s_storage_archive_button').on('click', function () {
    var $table = $('#table'),
    btntitle = "",
    btnicon = "";
    if ($('#s_storage_archive_button').hasClass("active")) {
        updateQueryStringParam("storage_archive", "false");
        btntitle = global.t("showdeleted_text", container.PersonLanguage);
        btnicon = "delete";
    } else {
        updateQueryStringParam("storage_archive", "true");
        btntitle = global.t("hidedeleted_text", container.PersonLanguage);     
        btnicon = "delete-forever";
    }
    $table.bootstrapTable('refresh');
    $('#s_storage_archive_button').attr("title", btntitle);
    $('#s_storage_archive_button > span').attr("class", "iconlabel mdi mdi-24px mdi-"+btnicon);
    $('#s_storage_archive_button > span').text(btntitle);
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
        for (i=1; i<=depth-1; i++) {
            html.push("<div class='col-sm-1'><span class='mdi mdi-square-medium mdi-24px'></span></div>");
        }
        html.push("<div class='col' style='color: " + sl.storelocation_color.String + "'>" + sl.storelocation_name.String + "</div>");
        
        for (var i in sl.stock) {
            var stock = sl.stock[i];
            
            if (stock.total === 0 && stock.current === 0) {
            } else {
                html.push("<div class='col'><span class='iconlabel'>" + global.t("stock_storelocation_title", container.PersonLanguage) + "</span> " + stock.current + " <b>" + stock.unit.unit_label.String + "</b></div>");
                html.push("<div class='col'><span class='iconlabel'>" + global.t("stock_storelocation_sub_title", container.PersonLanguage) + "</span> " + stock.total + " <b>" + stock.unit.unit_label.String + "</b></div>");
                
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

        html.push("<div class='row mb-sm-3'>")
            html.push("<div class='col-sm-6'>")
                html.push("<span class='iconlabel'>id</span> " + row["storage_id"]["Int64"])
            html.push("</div>")
        html.push("</div>")

        html.push("<div class='row mb-sm-3'>")
            html.push("<div class='col-sm-6'><span class='mdi mdi-24px mdi-tag'></span> " + row["product"]["name"]["name_label"] + "</div>")
            html.push("<div class='col-sm-6'><span class='mdi mdi-24px mdi-docker'></span> " + row["storelocation"]["storelocation_name"]["String"] + "</div>")
        html.push("</div>")
    
        html.push("<div class='row mb-sm-3'>")
            html.push("<div class='col-sm-6'><span class='iconlabel'>" + global.t("storage_quantity_title", container.PersonLanguage) + "</span> " + row["storage_quantity"]["Float64"] + " " + row["unit"]["unit_label"]["String"] + "</div>")
            html.push("<div class='col-sm-6'><span class='iconlabel'>" + global.t("storage_barecode_title", container.PersonLanguage) + "</span> " + row["storage_barecode"]["String"] + "</div>")
        html.push("</div>")
    
        html.push("<div class='row mb-sm-3'>")
            if (row["storage_batchnumber"]["Valid"] && row["storage_batchnumber"]["String"] != "") {
                html.push("<div class='col-sm-6'><span class='iconlabel'> " + global.t("storage_batchnumber_title", container.PersonLanguage) + "</span> " + row["storage_batchnumber"]["String"] + "</div>")
            }
            if (row["supplier"]["supplier_label"]["Valid"]) {
                html.push("<div class='col-sm-6'><span class='iconlabel'> " + global.t("supplier_label_title", container.PersonLanguage) + "</span> " + row["supplier"]["supplier_label"]["String"] + "</div>")
            }
        html.push("</div>")
    
        html.push("<div class='row mb-sm-3'>")
            if (row["storage_entrydate"]["Valid"]) {
                html.push("<div class='col-sm-12'><span class='iconlabel'> " + global.t("storage_entrydate_title", container.PersonLanguage) + "</span> " + dateFormatter(row["storage_entrydate"]["Time"], null, null, null) + "</div>")
            }
            if (row["storage_exitdate"]["Valid"]) {
                html.push("<div class='col-sm-12'><span class='iconlabel'> " + global.t("storage_exitdate_title", container.PersonLanguage) + "</span> " + dateFormatter(row["storage_exitdate"]["Time"], null, null, null) + "</div>")
            }
            if (row["storage_openingdate"]["Valid"]) {
                html.push("<div class='col-sm-12'><span class='iconlabel'> " + global.t("storage_openingdate_title", container.PersonLanguage) + "</span> " + dateFormatter(row["storage_openingdate"]["Time"], null, null, null) + "</div>")
            }
            if (row["storage_expirationdate"]["Valid"]) {
                html.push("<div class='col-sm-12'><span class='iconlabel'> " + global.t("storage_expirationdate_title", container.PersonLanguage) + "</span> " + dateFormatter(row["storage_expirationdate"]["Time"], null, null, null) + "</div>")
            }
        html.push("</div>")   
    
        html.push("<div class='row mb-sm-3'>")
        if (row["storage_comment"]["Valid"] && row["storage_comment"]["String"] != "") {
            html.push("<div class='col-sm-12'><span class='iconlabel'> " + global.t("storage_comment_title", container.PersonLanguage) + "</span> " + row["storage_comment"]["String"] + "</div>")
        }
        html.push("</div>")  

        html.push("<div class='row mb-sm-3'>")
            html.push("<div class='col-sm-8'><span class='iconlabel'> " + global.t("created", container.PersonLanguage) + "</span> " + dateFormatter(row["storage_creationdate"], null, null, null) + " <span class='iconlabel'> " + global.t("modified", container.PersonLanguage) + "</span> " + dateFormatter(row["storage_modificationdate"], null, null, null) + "</div>")
            html.push("<div class='col-sm-4'><p class='blockquote-footer'>" + row["person"]["person_email"] + "</p></div>")
        html.push("</div>")   
      
        html.push("<div class='row mb-sm-3'>")
            if (row["storage_todestroy"]["Bool"]) {
                html.push("<div class='col-sm-12'><span title='to destroy' class='mdi mdi-24px mdi-delete-sweep'></span></div>")
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
    if (value == "") {
        return ""
    } else {
        date = new Date(value);
        return date.toLocaleDateString();
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
// storage_quantityFormatter
//
function storage_quantityFormatter(value, row, index, field) {
    ret = "";
    if (row.storage_quantity.Valid) {
        ret += row.storage_quantity.Float64
    }
    if (row.unit.unit_label.Valid) {
        ret += " " + row.unit.unit_label.String
    } 
    return ret;
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
// storage_product formatter
//
function storage_productFormatter(value, row, index, field) {
    if (value == "") {
        return ""
    } else {
        return "<a href='" + proxyPath + "v/products?product=" + row["product"]["product_id"] + "'>" + row["product"]["name"]["name_label"] + "</a>";
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
    borrowingtitle = global.t("storage_borrow", container.PersonLanguage);
    if (row.borrowing.borrowing_id.Valid) {
        borrowingicon = "hand-pointing-left";
        borrowingtitle = global.t("storage_unborrow", container.PersonLanguage);
    }
    
    if (row.storage_archive.Valid && row.storage_archive.Bool) {
        // this is an archive
        var actions = [
            '<button id="clone' + sid + '" sid="' + sid + '" class="clone btn btn-link btn-sm" style="display: none;" title="' + global.t("storage_clone", container.PersonLanguage) + '" type="button">',
            '<span class="mdi mdi-24px mdi-content-copy"></span>',
            '</button>',
            '<button id="restore' + sid + '" sid="' + sid + '" class="restore btn btn-link btn-sm" style="display: none;" title="' + global.t("storage_restore", container.PersonLanguage) + '" type="button">',
            '<span class="mdi mdi-24px mdi-undo"></span>',
            '</button>',
            '<button id="delete' + sid + '" sid="' + sid + '" class="delete btn btn-link btn-sm" style="display: none;" title="' + global.t("delete", container.PersonLanguage) + '" type="button">',
            '<span class="mdi mdi-24px mdi-delete"></span>',
            '</button>']; 
    } else if (row.storage.storage_id.Valid) {
        // this is an history
        var actions = [
            '<button id="clone' + sid + '" class="clone btn btn-link btn-sm" style="display: none;" title="' + global.t("storage_clone", container.PersonLanguage) + '" type="button">',
            '<span class="mdi mdi-24px mdi-content-copy"></span>',
            '</button>'];
    } else {
        // buttons are hidden by default
        var actions = [
            '<button id="edit' + sid + '" sid="' + sid + '" class="edit btn btn-link btn-sm" style="display: none;" title="' + global.t("edit", container.PersonLanguage) + '" type="button">',
            '<span class="mdi mdi-24px mdi-border-color"></span>',
            '</button>',
            '<button id="clone' + sid + '" sid="' + sid + '" class="clone btn btn-link btn-sm" style="display: none;" title="' + global.t("storage_clone", container.PersonLanguage) + '" type="button">',
            '<span class="mdi mdi-24px mdi-content-copy"></span>',
            '</button>',
            '<button id="archive' + sid + '" sid="' + sid + '" class="archive btn btn-link btn-sm" style="display: none;" title="' + global.t("delete", container.PersonLanguage) + '" type="button">',
            '<span class="mdi mdi-24px mdi-delete"></span>',
            '</button>',
            '<button id="borrow' + sid + '" sid="' + sid + '" data-target="#borrow" class="borrow btn btn-link btn-sm" style="display: none;" title="' + borrowingtitle + '" type="button">',
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
    
    //if ((row.storage_creationdate != row.storage_modificationdate) && !row.storage.storage_id.Valid && !(row.storage_archive.Valid && row.storage_archive.Bool)) {
    if (row.storage_hc != 0) {
        actions.push('<button id="history' + sid + '" class="history btn btn-link btn-sm" title="' + global.t("storage_showhistory", container.PersonLanguage) + '" type="button">');
        actions.push('<span class="mdi mdi-24px mdi-history"></span>');
        actions.push('</button>');
    }
        
    return actions.join('&nbsp;')    
}

//
// items actions callback
//
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
        $("button#restore" + row.storage_id.Int64).confirmation("show").off( "confirmed.bs.confirmation");
        $("button#restore" + row.storage_id.Int64).confirmation("show").off( "canceled.bs.confirmation");
        
        // ask for confirmation and then restore
        $("button#restore" + row.storage_id.Int64).confirmation("show").on( "confirmed.bs.confirmation", function() {
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
        $("button#archive" + row.storage_id.Int64).confirmation("show").off( "confirmed.bs.confirmation");
        $("button#archive" + row.storage_id.Int64).confirmation("show").off( "canceled.bs.confirmation");
        
        // ask for confirmation and then delete
        $("button#archive" + row.storage_id.Int64).confirmation("show").on( "confirmed.bs.confirmation", function() {
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
        $("button#delete" + row.storage_id.Int64).confirmation("show").off( "confirmed.bs.confirmation");
        $("button#delete" + row.storage_id.Int64).confirmation("show").off( "canceled.bs.confirmation");
        
        // ask for confirmation and then delete
        $("button#delete" + row.storage_id.Int64).confirmation("show").on( "confirmed.bs.confirmation", function() {
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
function operateBorrow(e, value, row, index) {
    // clean select2 input selection
    $('select#borrower').val(null).trigger('change');
    $('select#borrower').find('option').remove();
    
    // get borrowed storage id
    $("input#bstorage_id").val(row['storage_id'].Int64);
    
    if (row['borrowing']['borrowing_id'].Valid) {
        // unborrow storage
        saveBorrowing();
    } else {
        // borrow storage, showing the modal
        $("#borrow").modal("show");
        var $table = $('#table');
        $table.bootstrapTable('refresh');

        // pre selecting the connected user
        var email = readCookie("email")
        var id = readCookie("id")
        var newOption = new Option(email, id, true, true);
        $('select#borrower').append(newOption).trigger('change');
    }   

}
function saveBorrowing() {
    var form = $("#borrowing");
    if (! form.valid()) {
        return;
    };

    var borrowing_comment = $("textarea#borrowing_comment").val(),
        borrower = $('select#borrower').select2('data')[0],
        storage_id = $("input#bstorage_id").val(),
        data = {};

    if (borrower !== undefined) {
        $.extend(data, {
            "borrowing_comment": borrowing_comment,
            "borrower.person_id": borrower.id,
        });
    }

    $.ajax({
        url: proxyPath + "borrowings/" + storage_id,
        method: "PUT",
        dataType: 'json',
        data: data,
    }).done(function(jqXHR, textStatus, errorThrown) {
       $("#borrow").modal("hide");
       global.displayMessage("borrowing modified", "success");
       var $table = $('#table');
       $table.bootstrapTable('refresh');
    }).fail(function(jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status);
    });  
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

    //
    // store location select2 formatter
    //
    function formatStorelocation (sl) {
        if (!sl.storelocation_id) {
            return sl.storelocation_fullpath;
        }
        var canstore = '<span class="mdi mdi-close"></span>';
        var icon = '<span class="mdi mdi-gesture" style="color: ' + sl.storelocation_color.String + ';"></span>';
        if (sl.storelocation_canstore.Valid && sl.storelocation_canstore.Bool) {
            canstore = '<span class="mdi mdi-check"></span>'
        }
        var s = $(
            '<div>' + icon + '<span>' + sl.storelocation_fullpath + '</span>' + canstore + '</div>'
        );
        return s;
    };
    //
    // signalwords select2 formatter
    //
    function formatSignalWord (signalword) {
        if (!signalword.signalword_id) {
            return signalword.signalword_label;
        }
        if (signalword.signalword_id.Valid) {
            return signalword.signalword_label.String;
        }
    };
    //
    // symbols select2 formatter
    //
    function formatSymbol (symbol) {
        if (!symbol.symbol_id) {
            return symbol.symbol_label;
        }
        var s = $(
            '<span><img src="data:' + symbol.symbol_image + '" title="' + symbol.symbol_label + '" /> ' + symbol.symbol_label + '</span>'
        );
        return s;
    };
    //
    // precautionary statements select2 formatter
    //
    function formatPrecautionaryStatement (ps) {
        if (!ps.precautionarystatement_id) {
            return ps.precautionarystatement_label;
        }
        var s = $(
            '<span><b>' + ps.precautionarystatement_reference + '</b> ' + ps.precautionarystatement_label + '</span>'
        );
        return s;
    };
    //
    // hazard statements select2 formatter
    //
    function formatHazardStatement (hs) {
        if (!hs.hazardstatement_id) {
            return hs.hazardstatement_label;
        }
        var s = $(
            '<span><b>' + hs.hazardstatement_reference + '</b> ' + hs.hazardstatement_label + '</span>'
        );
        return s;
    };
    
    //
    // save storage callback
    //
    var createCallBack = function createCallback(data, textStatus, jqXHR) {
        global.displayMessage("storage created", "success");
        setTimeout(function(){ window.location = proxyPath + "v/storages?product="+data.product.product_id+"&hl="+data.product.product_id; }, 1000);
    }
    var updateCallBack = function updateCallback(data, textStatus, jqXHR) {
        //- $('#list-collapse').collapse('show');
        //- $('#edit-collapse').collapse('hide');
        //- var $table = $('#table');
        //- var index = $('input#index').val();
        //- // appending hl parameter to highlight updated row
        //- updateQueryStringParam("hl", data.storage_id.Int64);
        //- $table.bootstrapTable('refresh');
        console.log(data)
        global.displayMessage("storage " + data.storage_id.Int64 + " updated", "success");
        setTimeout(function(){ window.location = proxyPath + "v/storages?product="+data.product.product_id+"&hl="+data.product.product_id; }, 1000);
    }

    function saveStorage() {
        var form = $("#storage");
        if (! form.valid()) {
            return;
        };

        var product_id = $("input#product_id").val(),
            storage_id = $("input#storage_id").val(),
            storage_nbitem = $("input#storage_nbitem").val(),
            storage_comment = $("input#storage_comment").val(),
            storage_quantity = $("input#storage_quantity").val(),
            storage_barecode = $("input#storage_barecode").val(),
            storage_entrydate = $("input#storage_entrydate").val(),
            storage_exitdate = $("input#storage_exitdate").val(),
            storage_openingdate = $("input#storage_openingdate").val(),
            storage_expirationdate = $("input#storage_expirationdate").val(),
            storage_reference = $("input#storage_reference").val(),
            storage_batchnumber = $("input#storage_batchnumber").val(),
            storage_todestroy = $("input#storage_todestroy:CHECKED").val(),
            unit = $('select#unit').select2('data')[0],
            supplier = $('select#supplier').select2('data')[0],
            storelocation = $('select#storelocation').select2('data')[0],
            ajax_url = proxyPath + "storages",
            ajax_method = "POST",
            ajax_callback = createCallBack,
            data = {};

            if ($("form#storage input#storage_id").length) {
                ajax_url = proxyPath + "storages/" + storage_id
                ajax_method = "PUT"
                ajax_callback = updateCallBack
            }

            $.extend(data, {
                "product.product_id": product_id,
                "storage_id": storage_id,
                "storage_nbitem": storage_nbitem == 0 ? 1 : storage_nbitem,
                "storage_comment": storage_comment,
                "storage_quantity": storage_quantity == "" ? 1 : storage_quantity,
                "storage_barecode": storage_barecode,
                "storage_reference": storage_reference,
                "storage_batchnumber": storage_batchnumber,
                "storage_todestroy": storage_todestroy == "on" ? true : false,
                "storelocation.storelocation_id": storelocation.id,
            });
            if (supplier !== undefined) {
                $.extend(data, {
                    "supplier.supplier_id": supplier.id == supplier.text ? -1 : supplier.id,
                    "supplier.supplier_label": supplier === undefined ? "" : supplier.text,
                });                    
            }
            if (unit !== undefined) {
                $.extend(data, {
                    "unit.unit_id": unit.id,
                    "unit.unit_label": unit.text,
                });                    
            }
            if (storage_entrydate !== "") {
                $.extend(data, {
                    "storage_entrydate": storage_entrydate,
                });
            }
            if (storage_exitdate !== "") {
                $.extend(data, {
                    "storage_exitdate": storage_exitdate,
                });
            }
            if (storage_openingdate !== "") {
                $.extend(data, {
                    "storage_openingdate": storage_openingdate,
                });
            }
            if (storage_expirationdate !== "") {
                $.extend(data, {
                    "storage_expirationdate": storage_expirationdate,
                });
            }

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
