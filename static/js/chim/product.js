//
// request performed at table data loading
//
$( document ).ready(function() {  
    //
    // update form validation
    //
    $( "#product" ).validate({
        // ignore required to validate select2
        ignore: "",
        errorClass: "alert alert-danger",
        rules: {
            name: {
                required: true,
            },
            empiricalformula: {
                required: true,
                remote: {
                    url: "",
                    type: "post",
                    beforeSend: function(jqXhr, settings) {
                        id = -1
                        if ($("form#product input#product_id").length) {
                            id = $("form#product input#product_id").val()
                        }
                        settings.url = proxyPath + "validate/product/" + id + "/empiricalformula/";
                    },
                    data: {
                        empiricalformula: function() {
                            return $('select#empiricalformula').select2('data')[0].text;
                        },
                    },
                },
            },
            casnumber: {
                required: true,
                remote: {
                    url: "",
                    type: "post",
                    beforeSend: function(jqXhr, settings) {
                        id = -1
                        if ($("form#product input#product_id").length) {
                            id = $("form#product input#product_id").val()
                        }
                        settings.url = proxyPath + "validate/product/" + id + "/casnumber/";
                    },
                    data: {
                        casnumber: function() {
                            return $('select#casnumber').select2('data')[0].text;
                        },
                        product_specificity:  function() {
                            return $('#product_specificity').val();
                        },
                    },
                },
            },
            cenumber: {
                remote: {
                    url: "",
                    type: "post",
                    beforeSend: function(jqXhr, settings) {
                        id = -1
                        if ($("form#product input#product_id").length) {
                            id = $("form#product input#product_id").val()
                        }
                        settings.url = proxyPath + "validate/product/" + id + "/cenumber/";
                    },
                    data: {
                        cenumber: function() {
                            return $('select#cenumber').select2('data')[0].text;
                        },
                    },
                },
            },
        },
        messages: {
            name: {
                required: global.t("required_input", container.PersonLanguage)
            },
            empiricalformula: {
                required: global.t("required_input", container.PersonLanguage)
            },
            casnumber: {
                required: global.t("required_input", container.PersonLanguage)
            }
        }, 
    });

    //
    // search form
    //
    $('select#s_storelocation').select2({
        templateResult: formatStorelocation,
        //placeholder: "store location",
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
        //placeholder: "select a cas number",
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
       //placeholder: "select a name",
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
        //placeholder: "select a formula",
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
        //placeholder: "select signal word",
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
    //     placeholder: "direct store location access",
    //     ajax: {
    //         url: proxyPath + 'storelocations',
    //         delay: 400,
    //             data: function (params) {
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
    // casnumber select2
    //
    $('select#casnumber').select2({
        tags: true,
        placeholder: global.t("product_cas_placeholder", container.PersonLanguage),
        createTag: function (params) {
            if ($("input#exactMatchCasNumber").val() == "true") {
                return null
            }
            return {
                id: params.term,
                text: params.term,
            }
        },
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

                isExactMatch=false;
                
                // looking for an exact match
                $.each(data.rows, function( index, value ) {
                    if(value.c == 1) {
                        isExactMatch=true;
                    }
                });
                
                // there is a match: setting the input field
                if (isExactMatch) {
                    $("input#exactMatchCasNumber").val("true");
                } else {
                    $("input#exactMatchCasNumber").val("false");
                }

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.casnumber_label;
                    obj.id = obj.id || obj.casnumber_id;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-casnumber-results li").length + 10;

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
    // cenumber select2
    //
    $('select#cenumber').select2({
        tags: true,
        placeholder: global.t("product_ce_placeholder", container.PersonLanguage),
        allowClear: true,
        createTag: function (params) {
            if ($("input#exactMatchCeNumber").val() == "true") {
                return null
            }
            return {
                id: params.term,
                text: params.term,
            }
        },
        ajax: {
            url: proxyPath + 'products/cenumbers/',
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
                    $("input#exactMatchCeNumber").val("true");
                } else {
                    $("input#exactMatchCeNumber").val("false");
                }

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.cenumber_label.String;
                    obj.id = obj.id || obj.cenumber_id.Int64;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-cenumber-results li").length + 10;

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
    // physicalstate select2
    //
    $('select#physicalstate').select2({
        allowClear: true,
        tags: true,
        placeholder: global.t("product_physicalstate_placeholder", container.PersonLanguage),
        createTag: function (params) {
            if ($("input#exactMatchPhysicalstate").val() == "true") {
                return null
            }
            return {
                id: params.term,
                text: params.term,
            }
        },
        ajax: {
            url: proxyPath + 'products/physicalstates/',
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
                    $("input#exactMatchPhysicalstate").val("true");
                } else {
                    $("input#exactMatchPhysicalstate").val("false");
                }

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.physicalstate_label.String;
                    obj.id = obj.id || obj.physicalstate_id.Int64;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-physicalstate-results li").length + 10;

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
    // signalword select2
    //
    $('select#signalword').select2({
        templateResult: formatSignalWord,
        allowClear: true,
        placeholder: global.t("product_signalword_placeholder", container.PersonLanguage),
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
    }).on("change", function (e) {
        $(this).valid(); // FIXME: see https://github.com/select2/select2/issues/3901
    });

    //
    // classofcompound select2
    //
    $('select#classofcompound').select2({
        allowClear: true,
        tags: true,
        placeholder: global.t("product_classofcompound_placeholder", container.PersonLanguage),
        createTag: function (params) {
            if ($("input#exactMatchClassofcompounds").val() == "true") {
                return null
            }
            return {
                id: params.term,
                text: params.term,
            }
        },
        ajax: {
            url: proxyPath + 'products/classofcompounds/',
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
                    $("input#exactMatchClassofcompounds").val("true");
                } else {
                    $("input#exactMatchClassofcompounds").val("false");
                }

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.classofcompound_label;
                    obj.id = obj.id || obj.classofcompound_id;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-classofcompounds-results li").length + 10;

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
    // name select2
    //
    $('select#name').select2({
        tags: true,
        placeholder: global.t("product_name_placeholder", container.PersonLanguage),
        createTag: function (params) {
            if ($("input#exactMatchName").val() == "true") {
                return null
            }
            return {
                id: params.term,
                text: params.term,
            }
        },
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

                isExactMatch=false;
                
                // looking for an exact match
                $.each(data.rows, function( index, value ) {
                    if(value.c == 1) {
                        isExactMatch=true;
                    }
                });
                
                // there is a match: setting the input field
                if (isExactMatch) {
                    $("input#exactMatchName").val("true");
                } else {
                    $("input#exactMatchName").val("false");
                }

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.name_label;
                    obj.id = obj.id || obj.name_id;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-name-results li").length + 10;

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
    // empirical formula select2
    //
    $('select#empiricalformula').select2({
        tags: true,
        placeholder: global.t("product_empiricalformula_placeholder", container.PersonLanguage),
        createTag: function (params) {
            if ($("input#exactMatchEmpiricalFormula").val() == "true") {
                return null
            }
            return {
                id: params.term,
                text: params.term,
            }
        },
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

                isExactMatch=false;
                
                // looking for an exact match
                $.each(data.rows, function( index, value ) {
                    if(value.c == 1) {
                        isExactMatch=true;
                    }
                });
                
                // there is a match: setting the input field
                if (isExactMatch) {
                    $("input#exactMatchEmpiricalFormula").val("true");
                } else {
                    $("input#exactMatchEmpiricalFormula").val("false");
                }

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.empiricalformula_label;
                    obj.id = obj.id || obj.empiricalformula_id;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-empiricalformula-results li").length + 10;

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
    // linear formula select2
    //
    $('select#linearformula').select2({
        tags: true,
        allowClear: true,
        placeholder: global.t("product_linearformula_placeholder", container.PersonLanguage),
        createTag: function (params) {
            if ($("input#exactMatchLinearFormula").val() == "true") {
                return null
            }
            return {
                id: params.term,
                text: params.term,
            }
        },
        ajax: {
            url: proxyPath + 'products/linearformulas/',
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
                    $("input#exactMatchLinearFormula").val("true");
                } else {
                    $("input#exactMatchLinearFormula").val("false");
                }

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.linearformula_label.String;
                    obj.id = obj.id || obj.linearformula_id.Int64;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-linearformula-results li").length + 10;

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
    // synonyms select2
    //
    $('select#synonyms').select2({
        tags: true,
        placeholder: global.t("product_synonyms_placeholder", container.PersonLanguage),
        createTag: function (params) {
            if ($("input#exactMatchSynonyms").val() == "true") {
                return null
            }
            return {
                id: params.term,
                text: params.term,
            }
        },
        ajax: {
            url: proxyPath + 'products/synonyms/',
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
                    $("input#exactMatchSynonyms").val("true");
                } else {
                    $("input#exactMatchSynonyms").val("false");
                }

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.name_label;
                    obj.id = obj.id || obj.name_id;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-synonyms-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {more: selectnbitems<data.total}
                };
            }
        }
    });

    //
    // symbols select2
    //
    $('select#symbols').select2({
        templateResult: formatSymbol,
        closeOnSelect: false,
        placeholder: global.t("product_symbols_placeholder", container.PersonLanguage),
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

    //
    // hazardstatements select2
    //
    $('select#hazardstatements').select2({
        templateResult: formatHazardStatement,
        closeOnSelect: false,
        placeholder: global.t("product_hazardstatements_placeholder", container.PersonLanguage),
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

    //
    // precautionarystatements select2
    //
    $('select#precautionarystatements').select2({
        templateResult: formatPrecautionaryStatement,
        closeOnSelect: false,
        placeholder: global.t("product_precautionarystatements_placeholder", container.PersonLanguage),
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
});

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
            var a = $("<a>").attr("href", proxyPath + "download/" + data.exportfn).html("<span class='mdi mdi-48px mdi-file-download'></span>");
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
    
    hasPermission("storages", "", "POST").done(function(){
        $(".store").fadeIn();
        localStorage.setItem("storages::POST", true);
    }).fail(function(){
        localStorage.setItem("storages::POST", false);
    })  
    hasPermission("storages", "-2", "GET").done(function(){
        $("#switchview").removeClass("d-none");

        $(".storages").fadeIn();
        $(".ostorages").fadeIn();
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

    $("table#table").find("tr").each(function( index, b ) {
        hasPermission("products", $(b).attr("product_id"), "DELETE").done(function(){
            $("#delete"+$(b).attr("product_id")).fadeIn();
            localStorage.setItem("products:" + $(b).attr("product_id") + ":DELETE", true);
        }).fail(function(){
            localStorage.setItem("products:" + $(b).attr("product_id") + ":DELETE", false);
        }) 
    });
});
$('#table').on('expand-row.bs.table', function (e, index, row, $detail) {
    var mol = row['product_molformula'];
    var Info = {
        color: "#FFFFFF",
        height: 300,
        width: 300,
        use: "HTML5",
        disableInitialConsole: true
    };

    if (mol.Valid) {
        Jmol.getTMApplet("myJmol", Info);
        $("#jsmol" + row["product_id"]).html(myJmol._code);
        myJmol.__loadModel(mol.String);
    }
})

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

    html.push("<div id='jsmol" + row["product_id"] + "'>")
    html.push("</div>")

    html.push("<div class='row mt-sm-3'>")

        html.push("<div class='col-sm-10'>")
            $.each(row["synonyms"], function (key, value) {
                html.push("<span>" + value["name_label"] + "</span> ");
            });
        html.push("</div>")

        html.push("<div class='col-sm-2'>")
            html.push("<span class='iconlabel'>id</span> " + row["product_id"])
        html.push("</div>")

    html.push("</div>")


    html.push("<div class='row mt-sm-3'>")

        html.push("<div class='col-sm-4'>")
            html.push("<span class='iconlabel'>" + global.t("casnumber_label_title", container.PersonLanguage) + "</span> " + row["casnumber"]["casnumber_label"])
            if (row["casnumber"]["casnumber_cmr"]["Valid"]) {
                html.push("<span class='iconlabel'>" + global.t("casnumber_cmr_title", container.PersonLanguage) + "</span> " + row["casnumber"]["casnumber_cmr"]["String"])
            }
        html.push("</div>")

        if (row["cenumber"]["cenumber_id"]["Valid"]) {
            html.push("<div class='col-sm-4'><span class='iconlabel'>" + global.t("cenumber_label_title", container.PersonLanguage) + "</span> " + row["cenumber"]["cenumber_label"]["String"] + "</div>")
        }

        if (row["product_msds"]["Valid"]) {
            html.push("<div class='col-sm-4'><span class='iconlabel'>" + global.t("product_msds_title", container.PersonLanguage) + "</span> <a href='" + row["product_msds"]["String"] + "'><span class='mdi mdi-link-variant mdi-24px'></span></a></div>")
        }

    html.push("</div>")


    html.push("<div class='row mt-sm-3'>")

        html.push("<div class='col-sm-4'><span class='iconlabel'>" + global.t("empiricalformula_label_title", container.PersonLanguage) + "</span> " + row["empiricalformula"]["empiricalformula_label"] + "</div>")
        
        if (row["linearformula"]["linearformula_id"]["Valid"]) {
            html.push("<div class='col-sm-4'><span class='iconlabel'>" + global.t("linearformula_label_title", container.PersonLanguage) + "</span> " + row["linearformula"]["linearformula_label"]["String"] + "</div>")
        }
        
        if (row["product_threedformula"]["Valid"] && row["product_threedformula"]["String"] != "") {
            html.push("<div class='col-sm-4'><span class='iconlabel'>" + global.t("product_threedformula_title", container.PersonLanguage) + "</span> <a href='" + row["product_threedformula"]["String"] + "'><span class='mdi mdi-link-variant mdi-24px'></span></a></div>")
        }

    html.push("</div>")


    html.push("<div class='row mt-sm-3'>")

        html.push("<div class='col-sm-4'>")
        $.each(row["symbols"], function (key, value) {
            html.push("<img src='data:" + value["symbol_image"] + "' alt='" + value["symbol_label"] + "' title='" + value["symbol_label"] + "'/>");
        });
        html.push("</div>")

        html.push("<div class='col-sm-4'>")
        if (row["signalword"]["signalword_label"]["Valid"]) {
            html.push("<span class='iconlabel'>" + global.t("signalword_label_title", container.PersonLanguage) + "</span> " + row["signalword"]["signalword_label"]["String"])
        }
        html.push("</div>")

        if (row["physicalstate"]["physicalstate_id"]["Valid"]) {
            html.push("<div class='col-sm-4'><span class='iconlabel'>" + global.t("physicalstate_label_title", container.PersonLanguage) + "</span> " + row["physicalstate"]["physicalstate_label"]["String"] + "</div>")
        }

    html.push("</div>")


    html.push("<div class='row mt-sm-3'>")

        html.push("<div class='col-sm-4'>")
        if (row["hazardstatements"] != null && row["hazardstatements"].length != 0) {
            html.push("<div><span class='iconlabel'>" + global.t("hazardstatement_label_title", container.PersonLanguage) + "</span></div>")
            html.push("<ul>")
            $.each(row["hazardstatements"], function (key, value) {
                html.push("<li>" + value["hazardstatement_reference"] + ": <i>" + value["hazardstatement_label"] + "</i></li>");
            });
            html.push("</ul>")
        }
        html.push("</div>")

        html.push("<div class='col-sm-4'>")
        if (row["precautionarystatements"] != null && row["precautionarystatements"].length != 0) {
            html.push("<div><span class='iconlabel'>" + global.t("precautionarystatement_label_title", container.PersonLanguage) + "</span></div>")
            html.push("<ul>")
            $.each(row["precautionarystatements"], function (key, value) {
                html.push("<li>" + value["precautionarystatement_reference"] + ": <i>" + value["precautionarystatement_label"] + "</i></li>");
            });
        html.push("</ul>")
        }
        html.push("</div>")

        html.push("<div class='col-sm-4'>")
        if (row["classofcompound"] != null && row["classofcompound"].length != 0) {
            html.push("<div><span class='iconlabel'>" + global.t("classofcompound_label_title", container.PersonLanguage) + "</span></div>")
            html.push("<ul>")
            $.each(row["classofcompound"], function (key, value) {
                html.push("<li>" + value["classofcompound_label"] + "</li>");
            });
            html.push("</ul>")
        }
        html.push("</div>")

    html.push("</div>")


    html.push("<div class='row mt-sm-3'>")

        if (row["product_disposalcomment"]["Valid"] && row["product_disposalcomment"]["String"] != "") {
            html.push("<div class='col-sm-12'><span class='iconlabel'>" + global.t("product_disposalcomment_title", container.PersonLanguage) + "</span> " + row["product_disposalcomment"]["String"] + "</div>")
        }

    html.push("</div>")

    html.push("<div class='row mt-sm-3'>")

        if (row["product_remark"]["Valid"] && row["product_remark"]["String"] != "") {
            html.push("<div class='col-sm-12'><span class='iconlabel'>" + global.t("product_remark_title", container.PersonLanguage) + "</span> " + row["product_remark"]["String"] + "</div>")
        }

    html.push("</div>")


    html.push("<div class='row mt-sm-3'>")

        html.push("<div class='col-sm-12'>")

        if (row["product_radioactive"]["Bool"]) {
            html.push("<span title='" + global.t("product_radioactive_title", container.PersonLanguage) + "' class='mdi mdi-36px mdi-radioactive'></span>")
        }
        if (row["product_restricted"]["Bool"]) {
            html.push("<span title='" + global.t("product_restricted_title", container.PersonLanguage) + "' class='mdi mdi-36px mdi-hand'></span>")
        }

        html.push("</div>")

    html.push("</div>")


    html.push("<div class='row mt-sm-4'>")
        html.push("<div class='col-sm-12'><p class='blockquote-footer'>" + row["person"]["person_email"] + "</p></div>")
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
    var actions = [];

    actions.push('<div class="float-right">');

    $.each(row.symbols, function( index, value ) {
        if (value.symbol_label == "SGH02") {
            actions.push('<img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAIvSURBVFiFzdgxbI5BGMDx36uNJsJApFtFIh2QEIkBtYiYJFKsrDaJqYOofhKRMFhsFhNNMBgkTAaD0KUiBomN1EpYBHWGnvi8vu9r7/2eti65vHfv3T3P/57nuXvfuyqlJCRVVQuk1AqRl1LqP9NKpJxbETKjocLgYi0VaLl49wXBxUIFwsVDBcGFQWEAg1FwYZbCGMajLBfmPkzgUZRbw2IKFzGPrRFw/bpvD/bn8jUkXM719f3A9eu+k3iXA/92Bnub2yYx1NgDfbrvXIYZx8dcThjBExxvPOmGltqLIzmuEt63QSVczc+z/2whSw2ThpbajS+4UgOq59O4gYFSuGaByWb8zKvwN8RXXKiBPc7PLaWx3ARqY37O1CBe5/cvO1huVy+ZnfSX7y9MYxRTNeX32lZj+/sXWNfVnV3g1tT/aJeQ5vAGp3L9eXbjTFv7NzzM9VncSSnNF2lp4MqjNYvcxwEcy+0HcQg32/q8Kndl+YrcgM9Z4YdsrZ21PtvxHT9yv1vNgr8cbiIrnMUmbKu177PwVZjLgKPNt4sCOKzF0ww32aF9CA+yxSZKoTqDlVnucI6lMxhpg76OuxhrKr8oIENyXx/xxQKTE/hUkIdLJ1tlRd3TwtF/KtcuSalVVdUwdvQe+Fd6ljhfl9NzRKT5I8cvq/B+xi3vzFfk+FaqbEUPvEtVuipXBIspX9VLlW4Q/8U1VGe4EKgYsED3tefBgt271y7dUlV/ygHpF8bRglXiwx7BAAAAAElFTkSuQmCC" alt="flammable" title="flammable">');
        }
    });

    if (row.product_sc != 0) {
        actions.push('<button id="storages' + pid + '" class="storages btn btn-link btn-sm" style="display: none;" title="storages" type="button">',
        '<span class="mdi mdi-24px mdi-cube-unfolded"><i>' + row.product_sc + '</i></span>',
        '</button>');
    } else {
        actions.push('<button class="btn btn-link btn-sm"><span class="mdi mdi-24px mdi-blank">&nbsp;</span></button>');
    }

    if (row.product_tsc != 0) {
        actions.push('<button id="ostorages' + pid + '" class="ostorages btn btn-link btn-sm" style="display: none;" title="global availability" type="button">',
        '<span class="mdi mdi-24px mdi-cube-scan">',
        '</button>');
    } else {
        actions.push('<button class="btn btn-link btn-sm"><span class="mdi mdi-24px mdi-blank">&nbsp;</span></button>');
    }

    actions.push(
    '<button id="store' + pid + '" class="store btn btn-link btn-sm" style="display: none;" title="store" type="button">',
        '<span class="mdi mdi-24px mdi-forklift">',
    '</button>',
    '<button id="edit' + pid + '" class="edit btn btn-link btn-sm" style="display: none;" title="edit" type="button">',
        '<span class="mdi mdi-24px mdi-border-color">',
    '</button>',
    '<button id="delete' + pid + '" class="delete btn btn-link btn-sm" style="display: none;" title="delete" type="button">',
        '<span class="mdi mdi-24px mdi-delete">',
    '</button>',
    '<button id="bookmark' + pid + '" class="bookmark btn btn-link btn-sm" title="(un)bookmark" type="button">',
        '<span id="bookmark' + pid + '" class="mdi mdi-24px mdi-' + bookmarkicon + '">',
    '</button>',
    '<div class="collapse" id="ostorages-collapse-' + pid + '"></div>'
    );

    if (row.casnumber.casnumber_cmr.Valid) {
        actions.push('<span title="CMR" class="mdi mdi-16px mdi-alert-outline text-danger"></span><span class="text-danger">' + row.casnumber.casnumber_cmr.String + '</span>');
    }    
    if (row.product_restricted.Valid && row.product_restricted.Bool) {
        actions.push('<span title="restricted access" class="mdi mdi-16px mdi-hand"></span>');
    }  

    actions.push('</div>');

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
    'click .ostorages': function (e, value, row, index) {
        operateOStorages(e, value, row, index)
    },
    'click .edit': function (e, value, row, index) {
        operateEdit(e, value, row, index)
    },
    'click .delete': function (e, value, row, index) {
        // hiding possible previous confirmation button
        $("button#delete" + row.product_id).confirmation("show").off( "confirmed.bs.confirmation");
        $("button#delete" + row.product_id).confirmation("show").off( "canceled.bs.confirmation");
        
        // ask for confirmation and then delete
        $("button#delete" + row.product_id).confirmation("show").on( "confirmed.bs.confirmation", function() {
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
function operateOStorages(e, value, row, index) {
   // getting the product
   $.ajax({
        url: proxyPath + "storages/others?product=" + row['product_id'],
        method: "GET",
    }).done(function(data, textStatus, jqXHR) {
        var html = [];
        $.each(data["rows"], function (key, value) {
            html.push("<p><span class='iconlabel'>" + value.entity_name + "</span><span class='blockquote-footer'>" + value.entity_description+ "</span></p>");
        });
        
        $("#ostorages-collapse-" + row['product_id']).html(html.join('&nbsp;'));
    }).fail(function(jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });

    // finally collapsing the view
    $('#ostorages-collapse-' + row['product_id']).collapse('show');
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

        if (data.signalword.signalword_id.Valid) {
            var newOption = new Option(data.signalword.signalword_label.String, data.signalword.signalword_id.Int64, true, true);
            $('select#signalword').append(newOption).trigger('change');
        }
       
        for(var i in data.symbols) {
           var newOption = new Option(data.symbols[i].symbol_label, data.symbols[i].symbol_id, true, true);
           $('select#symbols').append(newOption).trigger('change');
        }
        
        for(var i in data.classofcompound) {
            var newOption = new Option(data.classofcompound[i].classofcompound_label, data.classofcompound[i].classofcompound_id, true, true);
            $('select#classofcompound').append(newOption).trigger('change');
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

    //
    // mol file selector
    //
    $('#product_molformula').change( function () {
		if ( ! window.FileReader ) {
			return alert( 'FileReader API is not supported by your browser.' );
		}
        molfile = $('#product_molformula')[0];
        if ( molfile.files && molfile.files[0] ) {
			file = molfile.files[0]; // The file
			fr = new FileReader(); // FileReader instance
            fr.onload = function () {
                $('#hidden_product_molformula_content').append(fr.result);
            };
			fr.readAsText(file);
		} else {
			// Handle errors here
			alert( "file not selected or browser incompatible." )
		}
    });

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
    // save store location callback
    //
    var createCallBack = function createCallback(data, textStatus, jqXHR) {
        global.displayMessage("product " + data.name.name_label + " created", "success");
        setTimeout(function(){ window.location = proxyPath + "v/products?product="+data.product_id+"&hl="+data.product_id; }, 1000);
    }
    var updateCallBack = function updateCallback(data, textStatus, jqXHR) {
        global.displayMessage("product " + data.name.name_label + " updated", "success");
        setTimeout(function(){ window.location = proxyPath + "v/products?product="+data.product_id+"&hl="+data.product_id; }, 1000);
    }
    function saveProduct() {
        var form = $("#product");
        if (! form.valid()) {
            return;
        };

        var product_id = $("input#product_id").val(),
            product_specificity = $("input#product_specificity").val(),
            product_threedformula = $("input#product_threedformula").val(),

            product_molformula = $("#hidden_product_molformula_content").html(),

            product_msds = $("input#product_msds").val(),
            product_disposalcomment = $("textarea#product_disposalcomment").val(),
            product_remark = $("textarea#product_remark").val(),
            product_restricted = $("input#product_restricted:CHECKED").val(),
            product_radioactive = $("input#product_radioactive:CHECKED").val(),
            casnumber = $('select#casnumber').select2('data')[0],
            cenumber = $('select#cenumber').select2('data')[0],
            empiricalformula = $('select#empiricalformula').select2('data')[0],
            linearformula = $('select#linearformula').select2('data')[0],
            name_ = $('select#name').select2('data')[0],
            physicalstate = $('select#physicalstate').select2('data')[0],
            signalword = $('select#signalword').select2('data')[0],
            classofcompound = $('select#classofcompound').select2('data'),
            synonyms = $('select#synonyms').select2('data'),
            symbols = $('select#symbols').select2('data'),
            hazardstatements = $('select#hazardstatements').select2('data'),
            precautionarystatements = $('select#precautionarystatements').select2('data'),
            ajax_url = proxyPath + "products",
            ajax_method = "POST",
            ajax_callback = createCallBack,
            data = {};

            if ($("form#product input#product_id").length) {
                ajax_url = proxyPath + "products/" + product_id
                ajax_method = "PUT"
                ajax_callback = updateCallBack
            }

            $.each(symbols, function( index, s ) {
                data["symbols." + index +".symbol_id"] = s.id;
                data["symbols." + index +".symbol_label"] = s.text;
            });
            $.each(synonyms, function( index, s ) {
                data["synonyms." + index +".name_id"] = s.id == s.text ? -1 : s.id;
                data["synonyms." + index +".name_label"] = s.text;
            });            
            $.each(classofcompound, function( index, s ) {
                data["classofcompound." + index +".classofcompound_id"] = s.id == s.text ? -1 : s.id;
                data["classofcompound." + index +".classofcompound_label"] = s.text;
            });
            $.each(hazardstatements, function( index, s ) {
                data["hazardstatements." + index +".hazardstatement_id"] = s.id;
                data["hazardstatements." + index +".hazardstatement_reference"] = s.text;
            });
            $.each(precautionarystatements, function( index, s ) {
                data["precautionarystatements." + index +".precautionarystatement_id"] = s.id;
                data["precautionarystatements." + index +".precautionarystatement_reference"] = s.text;
            });
            $.extend(data, {
                "product_id": product_id,
                "product_disposalcomment": product_disposalcomment,
                "product_remark": product_remark,
                "product_restricted": product_restricted == "on" ? true : false,
                "product_radioactive": product_radioactive == "on" ? true : false,
                "casnumber.casnumber_id": casnumber.id == casnumber.text ? -1 : casnumber.id,
                "casnumber.casnumber_label": casnumber.text,
                "empiricalformula.empiricalformula_id": empiricalformula.id == empiricalformula.text ? -1 : empiricalformula.id,
                "empiricalformula.empiricalformula_label": empiricalformula.text,
                "name.name_id": name_.id == name_.text ? -1 : name_.id,
                "name.name_label": name_.text,
            });
            if (product_molformula !== "") {
                $.extend(data, {
                    "product_molformula": product_molformula,
                });  
            }
            if (product_specificity !== "") {
                $.extend(data, {
                    "product_specificity": product_specificity,
                });                    
            }
            if (product_msds !== "") {
                $.extend(data, {
                    "product_msds": product_msds,
                });                    
            } 
            if (product_threedformula !== "") {
                $.extend(data, {
                    "product_threedformula": product_threedformula,
                });                    
            } 
            if (cenumber !== undefined) {
                $.extend(data, {
                    "cenumber.cenumber_id": cenumber.id == cenumber.text ? -1 : cenumber.id,
                    "cenumber.cenumber_label": cenumber === undefined ? "" : cenumber.text,
                });                    
            }
            if (physicalstate !== undefined) {
                $.extend(data, {
                    "physicalstate.physicalstate_id": physicalstate.id == physicalstate.text ? -1 : physicalstate.id,
                    "physicalstate.physicalstate_label": physicalstate.text,
                });                    
            }
            if (signalword !== undefined) {
                $.extend(data, {
                    "signalword.signalword_id": signalword.id == signalword.text ? -1 : signalword.id,
                    "signalword.signalword_label": signalword.text,
                });                    
            }
            if (linearformula !== undefined) {
                $.extend(data, {
                    "linearformula.linearformula_id": linearformula.id == linearformula.text ? -1 : linearformula.id,
                    "linearformula.linearformula_label": linearformula.text,
                });                    
            } 
            $.ajax({
                url: ajax_url,
                method: ajax_method,
                dataType: 'json',
                data: data,
            }).done(ajax_callback).fail(function(jqXHR, textStatus, errorThrown) {
                handleHTTPError(jqXHR.statusText, jqXHR.status)
            });  
        }

    //
    // magical selector
    //
    function magic() {
        var m = $('textarea#magical').val();
        if (m != "") {
            $.ajax({
                url: proxyPath + "products/magic",
                method: "POST",
                data: {"msds": m}
            }).done(function(data, textStatus, jqXHR) {
                $('select#hazardstatements').val(null).trigger('change');
                $('select#hazardstatements').find('option').remove();
                $('select#precautionarystatements').val(null).trigger('change');
                $('select#precautionarystatements').find('option').remove();
                
                var hs = data.hs,
                    ps = data.ps;
                for(var i= 0; i < hs.length; i++)
                {
                   var newOption = new Option(data.hs[i].hazardstatement_reference, data.hs[i].hazardstatement_id, true, true);
                   $('select#hazardstatements').append(newOption).trigger('change');
                }
                for(var i= 0; i < ps.length; i++)
                {
                   var newOption = new Option(data.ps[i].precautionarystatement_reference, data.ps[i].precautionarystatement_id, true, true);
                   $('select#precautionarystatements').append(newOption).trigger('change');
                }
            }).fail(function(jqXHR, textStatus, errorThrown) {
                handleHTTPError(jqXHR.statusText, jqXHR.status)
            });
        }
    }

    //
    // linear to empirical formula converter
    //
    function linearToEmpirical() {
        var f = $('select#linearformula').select2('data')[0].text;
        if (f != "") {
            $.ajax({
                url: proxyPath + "products/l2eformula/" + f,
                method: "GET",
            }).done(function(data, textStatus, jqXHR) {
                $("#fconverter").attr("data-content", data);
                $("#fconverter").popover('show');
            }).fail(function(jqXHR, textStatus, errorThrown) {
                handleHTTPError(jqXHR.statusText, jqXHR.status)
            });
        }
    }

    //
    // zero empirical formula selector
    //
    function noEmpiricalFormula() {
        $.ajax({
            url: proxyPath + 'products/empiricalformulas/',
            type: "GET",
            data: {
                search: "XXXX",
                page: 1,
                offset: 0,
                limit: 1,
            },
        }).done(function(data, textStatus, jqXHR) {
            console.log(data)
            var newOption = new Option(data.rows[0].empiricalformula_label, data.rows[0].empiricalformula_id, true, true);
            $('select#empiricalformula').append(newOption).trigger('change');
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status)
        });
    }

    //
    // zero CAS number selector
    //
    function noCASNumber() {
        $.ajax({
            url: proxyPath + 'products/casnumbers/',
            type: "GET",
            data: {
                search: "0000-00-0",
                page: 1,
                offset: 0,
                limit: 1,
            },
        }).done(function(data, textStatus, jqXHR) {
            console.log(data)
            var newOption = new Option(data.rows[0].casnumber_label, data.rows[0].casnumber_id, true, true);
            $('select#casnumber').append(newOption).trigger('change');
        }).fail(function(jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status)
        });
    }

    //
    // magical selector how to
    //
    function howToMagicalSelector() {
        window.open(proxyPath + "img/magicalselector.webm", '_blank');
    }