//
// product and storage view common code
//

//
// startup actions
//
$( document ).ready(function() {  
    $('body').keypress(function(event){
        var keycode = (event.keyCode ? event.keyCode : event.which);
        if(keycode == '13'){
            search();	
        }
    });
    $('#s_casnumber,#s_empiricalformula,#s_storelocation,#s_name,#s_signalword').on("select2:close", function(e) {
        setTimeout(function() {
            $('.select2-container-active').removeClass('select2-container-active');
            $(':focus').blur();
        }, 1);
    });
});
$('#table').on('load-success.bs.table refresh.bs.table', function () {
      
    //console.log("tableload")

    // getting request parameters
    var urlParams = new URLSearchParams(window.location.search);

    // highlight row if needed
    if (urlParams.has("hl")) {
        $("tr[storage_id=" + urlParams.get("hl") + "]").addClass("animated bounce slow");
        $("tr[product_id=" + urlParams.get("hl") + "]").addClass("animated bounce slow");
    }
    
    var storelocationpromise = $.Deferred();
    var entitypromise = $.Deferred();
    var storagepromise = $.Deferred();
    var productpromise = $.Deferred();
    var bookmarkpromise = $.Deferred();
    var historypromise = $.Deferred();
    var productData, storageData, storelocationData, entityData, historyData, bookmarkData;

    // display titles, switch products<>storages selector
    if (urlParams.has("entity")) {
        e = urlParams.get("entity");
        entitypromise = $.ajax({
            url: proxyPath + "entities/" + e,
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            entityData = data;
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status)
        });
    } else {
        entitypromise.resolve();
    }
    if (urlParams.has("storelocation")) {
        s = urlParams.get("storelocation");
        storelocationpromise = $.ajax({
            url: proxyPath + "storelocations/" + s,
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            storelocationData = data;
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status)
        });
    } else {
        storelocationpromise.resolve();
    }
    if (urlParams.has("storage")) {
        s = urlParams.get("storage");
        storagepromise = $.ajax({
            url: proxyPath + "storages/" + s,
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            storageData = data;
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status)
        });
        // hiding the storage<>product view switch
        $(".toggleable").hide();
    } else {
        storagepromise.resolve();
    }
    if (urlParams.has("product")) {
        p = urlParams.get("product");
        // setting the title
        productpromise = $.ajax({
            url: proxyPath + "products/" + p,
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            productData = data;
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status)
        });
        // hiding the storage<>product view switch
        $(".toggleable").hide();
        
        // we should put "POST" here but this lead to a 405 error
        hasPermission("storages", "-2", "PUT").done(function(){
            // adding a create storage button for the product
            b = $("<button>").addClass("store btn btn-link").attr("type", "button");
            b.attr("onclick", "window.location.href = '" + proxyPath + "vc/storages?product=" + p + "'");
            i = $("<span>").addClass("mdi mdi-24px mdi-forklift iconlabel").html(global.t("storeagain_text", container.PersonLanguage));
            b.append(i);
            $("#button-store").html(b);
        })
        
        hasPermission("storages", "-2", "GET").done(function(){
            // adding a show stock button for the product
            b = $("<button>").addClass("store btn btn-link").attr("type", "button").attr("data-toggle", "modal").attr("data-target", "#stock");
            b.attr("onclick", "showStock(" + p + ")");
            i = $("<span>").addClass("mdi mdi-24px mdi-sigma iconlabel").html(global.t("totalstock_text", container.PersonLanguage));;
            b.append(i);
            $("#button-stock").html(b);
           
        })
        
    } else {
        productpromise.resolve();
    }
    if (urlParams.has("history")) {
        historyData = true;
    }
    historypromise.resolve();
    if (urlParams.has("bookmark")) {
        bookmarkData = true;
    }
    bookmarkpromise.resolve();
    
    $.when(entitypromise, storelocationpromise, storagepromise, productpromise, bookmarkpromise).done(function() {
        if (bookmarkData !== undefined || historyData !== undefined || entityData !== undefined || storelocationData !== undefined || storageData !== undefined || productData !== undefined) {
            $("#filter-item").html("");
        }
        if (bookmarkData !== undefined) {
            $("#filter-item").append(global.createTitle(global.t("menu_bookmark", container.PersonLanguage), "bookmark"));
        };
        if (historyData !== undefined) {
            $("#filter-item").append(global.createTitle(global.t("storage_history", container.PersonLanguage), "history"));
        };
        if (entityData !== undefined) {
            $("#filter-item").append(global.createTitle(entityData.entity_name, "entity"));
        };
        if (storelocationData !== undefined) {
            $("#filter-item").append(global.createTitle(storelocationData.storelocation_name.String, "storelocation"));
        };
        if (storageData !== undefined) {
            $("#filter-item").append(global.createTitle(storageData.product.name.name_label + " " + storageData.product.product_specificity.String + " (" + storageData.product.casnumber.casnumber_label + ") - " + storageData.storelocation.storelocation_fullpath, "storage"));
        };
        if (productData !== undefined) {
            $("#filter-item").append(global.createTitle(productData.name.name_label + " (" + productData.casnumber.casnumber_label + ") " + productData.product_specificity.String, "product"));
        };
    }); 

    // need to clean the request from its parameters to avoid selecting
    // former search parameters
    //cleanQueryParams();
});
//
// close buttons actions
//
function closeEdit() {
    
    // getting request parameters
    var urlParams = new URLSearchParams(window.location.search);
    
    $("#list-collapse").collapse("show");
    $("#edit-collapse").collapse("hide"); 
    if (!urlParams.get("product") && !urlParams.get("entity")) {
        $(".toggleable").show()
    }
    $('div#search').collapse('show');
}

//
// table data loading
//
function queryParams(params) {
    
    //console.log("queryParams")

    // getting request parameters
    // window.location.search only populated on product/storage view
    var urlParams = new URLSearchParams(window.location.search);
    
    // search form parameters
    var storelocation = urlParams.get("storelocation");
    var name = urlParams.get("name");
    var casnumber = urlParams.get("casnumber");
    var empiricalformula = urlParams.get("empiricalformula");
    var storage_barecode = urlParams.get("storage_barecode");
    var storage_archive = urlParams.get("storage_archive");
    var custom_name_part_of = urlParams.get("custom_name_part_of");
    var signalword = urlParams.get("signalword");
    var symbols = urlParams.getAll("symbols[]");
    var hazardstatements = urlParams.getAll("hazardstatements[]");
    var precautionarystatements = urlParams.getAll("precautionarystatements[]");
    var casnumber_cmr = urlParams.get("casnumber_cmr");
    // parameters passed by url
    var entity = urlParams.get("entity");
    var history = urlParams.get("history");
    var storage = urlParams.get("storage");
    var storage_archive = urlParams.get("storage_archive");
    var product = urlParams.get("product");
    var bookmark = urlParams.get("bookmark");

    //
    // populating the form in case of product/storage view switch
    // parameters are passed by URL
    //
    // hidden form parameters
    if (entity != null) {
        $('#hidden_s_entity').val(entity);
    }
    if (product != null) {
        $('#hidden_s_product').val(product);
    }
    if (history != null) {
        $('#hidden_s_history').val(history);
    }
    if (storage != null) {
        $('#hidden_s_storage').val(storage);
    }
    if (storage_archive != null) {
        $('#hidden_s_storage_archive').val(storage_archive);
    }
    if (bookmark != null) {
        $('#hidden_s_bookmark').val(bookmark);
    }
    // search form parameters
    if (storelocation != null) {
        $("#advancedsearch").collapse('show');
        $.ajax({
            url: proxyPath + "storelocations/" + storelocation,
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            var newOption = new Option(data.storelocation_name.String, data.storelocation_id.Int64, true, true);
            $('#s_storelocation').append(newOption).trigger('change');
            $('#s_storelocation').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status);
        })
    }
    if (name != null) {
        $("#advancedsearch").collapse('show');
        $.ajax({
            url: proxyPath + "products/names/" + name,
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            var newOption = new Option(data.name_label, data.name_id, true, true);
            $('#s_name').append(newOption).trigger('change');
            $('#s_name').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status);
        })
    }
    if (casnumber != null) {
        $.ajax({
            url: proxyPath + "products/casnumbers/" + casnumber,
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            var newOption = new Option(data.casnumber_label, data.casnumber_id, true, true);
            $('#s_casnumber').append(newOption).trigger('change');
            $('#s_casnumber').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status);
        })
    }
    if (empiricalformula != null) {
        $.ajax({
            url: proxyPath + "products/empiricalformulas/" + empiricalformula,
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            var newOption = new Option(data.empiricalformula_label, data.empiricalformula_id, true, true);
            $('#s_empiricalformula').append(newOption).trigger('change');
            $('#s_empiricalformula').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status);
        })
    }           
    if (storage_barecode != null) {
        $("input#s_storage_barecode").val(storage_barecode)
        $('input#s_storage_barecode').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
    }
    if (custom_name_part_of != null) {
        $("input#s_custom_name_part_of").val(custom_name_part_of)
        $('input#s_custom_name_part_of').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
    }
    if (signalword != null) {
        $("#advancedsearch").collapse('show');
        $.ajax({
            url: proxyPath + "products/signalwords/" + signalword,
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            var newOption = new Option(data.signalword_label.String, data.signalword_id.Int, true, true);
            $('#s_signalword').append(newOption).trigger('change');
            $('#s_signalword').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status);
        })
    }
    if (symbols.length != 0) {
        $("#advancedsearch").collapse('show');
        $.each(symbols, function (index, symbol) { 
            $.ajax({
                url: proxyPath + "products/symbols/" + symbol,
                method: "GET",
            }).done(function(data, textStatus, jqXHR) {
                var newOption = new Option(data.symbol_label, data.symbol_id, true, true);
                $('#s_symbols').append(newOption).trigger('change');
                $('#s_symbols').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
            }).fail(function(jqXHR, textStatus, errorThrown) {
                handleHTTPError(jqXHR.statusText, jqXHR.status);
            }) 
        });
    }
    if (hazardstatements.length != 0) {
        $("#advancedsearch").collapse('show');
        $.each(hazardstatements, function (index, hs) { 
            $.ajax({
                url: proxyPath + "products/hazardstatements/" + hs,
                method: "GET",
            }).done(function(data, textStatus, jqXHR) {
                var newOption = new Option(data.hazardstatement_reference, data.hazardstatement_id, true, true);
                $('#s_hazardstatements').append(newOption).trigger('change');
                $('#s_hazardstatements').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
            }).fail(function(jqXHR, textStatus, errorThrown) {
                handleHTTPError(jqXHR.statusText, jqXHR.status);
            }) 
        });
    }
    if (precautionarystatements.length != 0) {
        $("#advancedsearch").collapse('show');
        $.each(precautionarystatements, function (index, ps) { 
            $.ajax({
                url: proxyPath + "products/precautionarystatements/" + ps,
                method: "GET",
            }).done(function(data, textStatus, jqXHR) {
                var newOption = new Option(data.precautionarystatement_reference, data.precautionarystatement_id, true, true);
                $('#s_precautionarystatements').append(newOption).trigger('change');
                $('#s_precautionarystatements').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
            }).fail(function(jqXHR, textStatus, errorThrown) {
                handleHTTPError(jqXHR.statusText, jqXHR.status);
            }) 
        });
    }
    if (casnumber_cmr == "true") {
        $("#advancedsearch").collapse('show');
        $('#s_casnumber_cmr').prop('checked', true);
        $('#s_casnumber_cmr').after("<span style='position: absolute; left: -5px; color: orange;' class='mdi mdi-checkbox-blank blink'></span>");
    }
    // storage view specific
    // if (storage_archive != null && storage_archive == "true") {
    //     $("#s_storage_archive").prop( "checked", true )
    // }

    //
    // populating bootstrap table ajax query parameters
    //
    if ($("#hidden_s_entity").val() != "") {
        params["entity"] = $("#hidden_s_entity").val();
    }
    if ($("#hidden_s_history").val() != "") {
        params["history"] = $("#hidden_s_history").val();
    }
    if ($("#hidden_s_storage").val() != "") {
        params["storage"] = $("#hidden_s_storage").val();
    }
    if ($("#hidden_s_storage_archive").val() != "") {
        params["storage_archive"] = $("#hidden_s_storage_archive").val();
    }
    if ($("#hidden_s_product").val() != "") {
        params["product"] = $("#hidden_s_product").val();
    }
    if ($("#hidden_s_bookmark").val() != "") {
        params["bookmark"] = $("#hidden_s_bookmark").val();
    }
    if (urlParams.has("export")) {
        params["export"] = urlParams.get("export")
    }

    // for select2 items we gather the value:
    // - from the url in case of storage/product view switch
    // - html input in case of table sorting use
    // -> we can NOT rely only on html input selection because
    //    they are initialize asynchronously by ajax calls 
    if (storelocation != null) {
        params["storelocation"] = storelocation;
    } else if ($('select#s_storelocation').hasClass("select2-hidden-accessible")) {
        sl = $('select#s_storelocation').select2('data')[0];
        if (sl != undefined) {
            params["storelocation"] = sl.id;
        }
    }
    if (name != null) {
        params["name"] = name;
    } else if ($('select#s_name').hasClass("select2-hidden-accessible")) {
        na = $('select#s_name').select2('data')[0];
        if (na != undefined) {
            params["name"] = na.id;
        }
    }
    if (casnumber != null) {
        params["casnumber"] = casnumber;
    } else if ($('select#s_casnumber').hasClass("select2-hidden-accessible")) {
        cas = $('select#s_casnumber').select2('data')[0];
        if (cas != undefined) {
            params["casnumber"] = cas.id;
        }
    }
    if (empiricalformula != null) {
        params["empiricalformula"] = empiricalformula;
    } else if ($('select#s_empiricalformula').hasClass("select2-hidden-accessible")) {
        ef = $('select#s_empiricalformula').select2('data')[0];
        if (ef != undefined) {
            params["empiricalformula"] = ef.id;
        }
    }
    if (signalword != null) {
        params["signalword"] = signalword;
    } else if ($('select#s_signalword').hasClass("select2-hidden-accessible")) {
        sw = $('select#s_signalword').select2('data')[0];
        if (sw != undefined) {
            params["signalword"] = sw.id;
        }
    }
    if (hazardstatements.length != 0) {
        params["hazardstatements"] = hazardstatements;
    } else if ($('select#s_hazardstatements').hasClass("select2-hidden-accessible")) {
        hs = $('select#s_hazardstatements').select2('data');
        if (hs.length != 0) {
            s_hazardstatements = [];
            hs.forEach(function(e) {
                s_hazardstatements.push(e.id);
            });
            params["hazardstatements"] = s_hazardstatements;
        }
    }
    if (precautionarystatements.length != 0) {
        params["precautionarystatements"] = precautionarystatements;
    } else if ($('select#s_precautionarystatements').hasClass("select2-hidden-accessible")) {
        ps = $('select#s_precautionarystatements').select2('data');
        if (ps.length != 0) {
            s_precautionarystatements = [];
            ps.forEach(function(e) {
                s_precautionarystatements.push(e.id);
            });
            params["precautionarystatements"] = s_precautionarystatements;
        }
    }
    if (symbols.length != 0) {
        params["symbols"] = symbols;
    } else if ($('select#s_symbols').hasClass("select2-hidden-accessible")) {
        sy = $('select#s_symbols').select2('data');
        if (sy.length != 0) {
            s_symbols = [];
            sy.forEach(function(e) {
                s_symbols.push(e.id);
            });
            params["symbols"] = s_symbols;
        }
    }

    // non select2 fields
    if ($('#s_storage_barecode').val() != "") {
        params["storage_barecode"] = storage_barecode;
    }
    if ($('#s_custom_name_part_of').val() != "") {
        params["custom_name_part_of"] = custom_name_part_of;
    }
    if ($('#s_casnumber_cmr:checked').length > 0) {
        params["casnumber_cmr"] = true;
    }

    return params;
}