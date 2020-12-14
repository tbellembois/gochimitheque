window.onload = function () {

    // cookie set by the GetTokenHandler() function after login
    var email = readCookie("email")
    var urlParams = new URLSearchParams(window.location.search);
    var message = urlParams.get("message");

    // displaying logged user
    if (email != null) {
        document.getElementById("logged").innerHTML = email;
    }
    if (message != null) {
        gjsUtils.message(message, "success");
    }

    // displaying version
    //document.getElementById("appversion").innerHTML = buildID;

    // showing menu items given the connected person permissions
    hasPermission("products", "-2", "GET").done(function () {
        $("#menu_scan_qrcode").show();
        $("#menu_list_products").show();
        $("#menu_list_bookmarks").show();
        localStorage.setItem("products:-2:GET", true);
    }).fail(function () {
        localStorage.setItem("products:-2:GET", false);
    })
    hasPermission("products", "", "POST").done(function () {
        $("#menu_create_product").show();
        localStorage.setItem("products::POST", true);
    }).fail(function () {
        localStorage.setItem("products::POST", false);
    })
    hasPermission("entities", "-2", "GET").done(function () {
        $("#menu_entities").show();
        localStorage.setItem("entities:-2:GET", true);
    }).fail(function () {
        localStorage.setItem("entities:-2:GET", false);
    })
    hasPermission("entities", "", "POST").done(function () {
        $("#menu_create_entity").show();
        localStorage.setItem("entities::POST", true);
    }).fail(function () {
        localStorage.setItem("entities::POST", false);
    })
    hasPermission("entities", "-2", "PUT").done(function () {
        $("#menu_update_welcomeannounce").show();
        localStorage.setItem("entities:-2:PUT", true);
    }).fail(function () {
        localStorage.setItem("entities:-2:PUT", false);
    })
    hasPermission("storages", "-2", "GET").done(function () {
        $("#menu_storelocations").show();
        localStorage.setItem("storages:-2:GET", true);
    }).fail(function () {
        localStorage.setItem("storages:-2:GET", false);
    })
    hasPermission("storelocations", "", "POST").done(function () {
        $("#menu_create_storelocation").show();
        localStorage.setItem("storelocations::POST", true);
    }).fail(function () {
        localStorage.setItem("storelocations::POST", false);
    })
    hasPermission("people", "-2", "GET").done(function () {
        $("#menu_people").show();
        localStorage.setItem("people:-2:GET", true);
    }).fail(function () {
        localStorage.setItem("people:-2:GET", false);
    })
    hasPermission("people", "", "POST").done(function () {
        $("#menu_create_person").show();
        localStorage.setItem("people::POST", true);
    }).fail(function () {
        localStorage.setItem("people::POST", false);
    })

};