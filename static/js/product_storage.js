//
// product and storage view common code
//

//
// startup actions
//
$( document ).ready(function() {  
});
$('#table').on('load-success.bs.table refresh.bs.table', function () {
    
    // getting request parameters
    var urlParams = new URLSearchParams(window.location.search);
    
    // highlight row if needed
    if (urlParams.has("hl")) {
        $("tr[storage_id=" + urlParams.get("hl") + "]").addClass("animated bounce slow");
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
        
        hasPermission("storages", "-2", "POST").done(function(){
            // adding a create storage button for the product
            b = $("<button>").addClass("store btn btn-primary btn-sm").attr("title", "store this product").attr("type", "button");
            b.attr("onclick", "window.location.href = '" + proxyPath + "vc/storages?product=" + p + "'");
            i = $("<span>").addClass("mdi mdi-24px mdi-forklift");
            b.append(i);
            $("#button-store").html(b);
        })
        
        hasPermission("storages", "-2", "GET").done(function(){
            // adding a create show stock button for the product
            b = $("<button>").addClass("store btn btn-primary btn-sm").attr("title", "show stock of this product").attr("type", "button").attr("data-toggle", "modal").attr("data-target", "#stock");
            b.attr("onclick", "showStock(" + p + ")");
            i = $("<span>").addClass("mdi mdi-24px mdi-sigma");
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
        $("#filter-item").html("");
        if (bookmarkData !== undefined) {
            $("#filter-item").append(global.createTitle("my bookmarked products", "bookmark"));
        };
        if (historyData !== undefined) {
            $("#filter-item").append(global.createTitle("history", "history"));
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

    cleanQueryParams();
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
    
    // getting request parameters
    var urlParams = new URLSearchParams(window.location.search);
    
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
    // need to populate the form in case of product/storage view switch
    if (storelocation != null) {
        $.ajax({
            url: proxyPath + "storelocations/" + storelocation,
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            var newOption = new Option(data.storelocation_name.String, data.storelocation_id.Int64, true, true);
            $('#s_storelocation').append(newOption).trigger('change');
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status);
        })
    }
    if (name != null) {
        $.ajax({
            url: proxyPath + "products/names/" + name,
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            var newOption = new Option(data.name_label, data.name_id, true, true);
            $('#s_name').append(newOption).trigger('change');
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
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status);
        })
    }           
    if (storage_barecode != null) {
        $("input#storage_barecode").val(storage_barecode)
    }
    if (custom_name_part_of != null) {
        $("input#custom_name_part_of").val(custom_name_part_of)
    }
    if (signalword != null) {
        $("#advancedsearch").collapse('show');
        $.ajax({
            url: proxyPath + "products/signalwords/" + signalword,
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            var newOption = new Option(data.signalword_label.String, data.signalword_id.Int, true, true);
            $('#s_signalword').append(newOption).trigger('change');
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
            }).fail(function(jqXHR, textStatus, errorThrown) {
                handleHTTPError(jqXHR.statusText, jqXHR.status);
            }) 
        });
    }

    // storage view specific
    if (storage_archive != null && storage_archive == "true") {
        $("#s_storage_archive").prop( "checked", true )
    }

    if (urlParams.has("entity")) {
        params["entity"] = urlParams.get("entity")
    }
    if (urlParams.has("storelocation")) {
        params["storelocation"] = urlParams.get("storelocation")
    }
    if (urlParams.has("product")) {
        params["product"] = urlParams.get("product")
    }
    if (urlParams.has("bookmark")) {
        params["bookmark"] = urlParams.get("bookmark")
    }
    if (urlParams.has("export")) {
        params["export"] = urlParams.get("export")
    }
    
    if (storelocation != null) {
        params["storelocation"] = storelocation
    }
    if (name != null) {
        params["name"] = name
    }
    if (casnumber != null) {
        params["casnumber"] = casnumber
    }
    if (empiricalformula != null) {
        params["empiricalformula"] = empiricalformula
    }
    if (storage_barecode != null) {
        params["storage_barecode"] = storage_barecode
    }
    if (custom_name_part_of != null && custom_name_part_of != "") {
        params["custom_name_part_of"] = custom_name_part_of
    }
    if (signalword != null) {
        params["signalword"] = signalword
    }
    if (hazardstatements.length != 0) {
        params["hazardstatements"] = hazardstatements
    }
    if (precautionarystatements.length != 0) {
        params["precautionarystatements"] = precautionarystatements
    }
    
    return params;
}