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
    
    // FIXME: this should be done before table rendering to avoid 2 ajax calls
    // populating search input if needed
    //if (URLValues.has("search")) {
    if (urlParams.has("search")) {
        //$('#table').bootstrapTable('resetSearch', URLValues.get("search")[0]);
        $('#table').bootstrapTable('resetSearch', urlParams.get("search"));
    }
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
    //if (URLValues.has("entity")) {
    if (urlParams.has("entity")) {
        //e = URLValues.get("entity")[0];
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
        //e = URLValues.get("entity")[0];
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
        //e = URLValues.get("entity")[0];
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
    //if (URLValues.has("product")) {
    if (urlParams.has("product")) {
        //p = URLValues.get("product")[0];
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
        
        hasPermission(proxyPath + "f/storages/-2", "POST").done(function(){
            // adding a create storage button for the product
            b = $("<button>").addClass("store btn btn-primary btn-sm").attr("title", "store this product").attr("type", "button");
            b.attr("onclick", "window.location.href = '" + proxyPath + "vc/storages?product=" + p + "'");
            i = $("<span>").addClass("mdi mdi-24px mdi-forklift");
            b.append(i);
            $("#button-store").html(b);
        })
        
        hasPermission(proxyPath + "f/storages/-2", "GET").done(function(){
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
    if (name != null) {
        $.ajax({
            url: proxyPath + "products/names/" + name,
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            var newOption = new Option(data.name_label, data.name_id, false, false);
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
            var newOption = new Option(data.casnumber_label, data.casnumber_id, false, false);
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
            var newOption = new Option(data.empiricalformula_label, data.empiricalformula_id, false, false);
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
        $.ajax({
            url: proxyPath + "products/signalwords/" + signalword,
            method: "GET",
        }).done(function(data, textStatus, jqXHR) {
            var newOption = new Option(data.signalword_label.String, data.signalword_id.Int, false, false);
            $('#s_signalword').append(newOption).trigger('change');
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status);
        })
    }
    if (symbols != null) {
        console.log(symbols)
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
    
    return params;
}