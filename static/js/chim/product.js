//
// request performed at table data loading
//
$(document).ready(function () {
    //
    // type chooser radio
    //
    $("input[type=radio][name=typechooser]").change(function() {
        if (this.value == 'chem') {
            chemfy()
        }
        else if (this.value == 'bio') {
            biofy()
        }
    });

    //
    // update form validation
    //
    $("#product").validate({
        // ignore required to validate select2
        ignore: "",
        errorClass: "alert alert-danger",
        rules: {
            name: {
                required: true,
            },
            // product_batchnumber: {
            //     required: true,
            // },
            producerref: {
                required: true,
            },
            // supplierref: {
            //     required: true,
            // },
            // unit_concentration: {
            //     required: function(element) {
            //         return $("#product_concentration").val() != "";
            //       }            
            // },
            unit_temperature: {
                required: function(element) {
                    return $("#product_temperature").val() != "";
                  }            
            },
            empiricalformula: {
                required: true,
                remote: {
                    url: "",
                    type: "post",
                    beforeSend: function (jqXhr, settings) {
                        id = -1
                        if ($("form#product input#product_id").length) {
                            id = $("form#product input#product_id").val()
                        }
                        settings.url = proxyPath + "validate/product/" + id + "/empiricalformula/";
                    },
                    data: {
                        empiricalformula: function () {
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
                    beforeSend: function (jqXhr, settings) {
                        id = -1
                        if ($("form#product input#product_id").length) {
                            id = $("form#product input#product_id").val()
                        }
                        settings.url = proxyPath + "validate/product/" + id + "/casnumber/";
                    },
                    data: {
                        casnumber: function () {
                            return $('select#casnumber').select2('data')[0].text;
                        },
                        product_specificity: function () {
                            return $('#product_specificity').val();
                        },
                    },
                },
            },
            cenumber: {
                remote: {
                    url: "",
                    type: "post",
                    beforeSend: function (jqXhr, settings) {
                        id = -1
                        if ($("form#product input#product_id").length) {
                            id = $("form#product input#product_id").val()
                        }
                        settings.url = proxyPath + "validate/product/" + id + "/cenumber/";
                    },
                    data: {
                        cenumber: function () {
                            return $('select#cenumber').select2('data')[0].text;
                        },
                    },
                },
            },
        },
        messages: {
            name: {
                required: gjsUtils.translate("required_input", container.PersonLanguage)
            },
            // product_batchnumber: {
            //     required: gjsUtils.translate("required_input", container.PersonLanguage)
            // },
            empiricalformula: {
                required: gjsUtils.translate("required_input", container.PersonLanguage)
            },
            casnumber: {
                required: gjsUtils.translate("required_input", container.PersonLanguage)
            },
            producerref: {
                required: gjsUtils.translate("required_input", container.PersonLanguage)
            }
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
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;
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
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
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
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;
                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.casnumber_label.String;
                    obj.id = obj.id || obj.casnumber_id.Int64;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-s_casnumber-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
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
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;

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
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
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
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.empiricalformula_label.String;
                    obj.id = obj.id || obj.empiricalformula_id.Int64;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-s_empiricalformula-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
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
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;
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
            processResults: function (data, params) {
                params.page = params.page || 1;
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
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;
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
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
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
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;
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
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
                };
            }
        }
    });

    //
    // producerref select2
    //
    $('select#producerref').on('select2:select', function (e) {
        if ($('select#producerref').select2('data').length != 0 && $('select#producerref').select2('data')[0].producer) {
            id = $('select#producerref').select2('data')[0].producerref_id.Int64;
            label = $('select#producerref').select2('data')[0].producerref_label.String;
        
            var newOption = new Option(label, id, true, true);
            $('select#producer').append(newOption).trigger('change');
        }
    });
    $('select#producerref').select2({
        tags: true,
        templateResult: formatProducerRef,
        templateSelection: formatProducerRef2,
        placeholder: "select a reference",
        createTag: function (params) {
            if ($("input#exactMatchProducerRef").val() == "true") {
                return null
            }
            if ($('select#producer').select2('data').length == 0) {
                gjsUtils.message("to create a new reference select a producer first", "warning");
                return null
            }
            return {
                id: params.term,
                text: params.term,
                producerlabel: $('select#producer').select2('data')[0].text,
                producerid: $('select#producer').select2('data')[0].id,
            }
        },
        ajax: {
            url: proxyPath + 'products/producerrefs/',
            delay: 400,
            data: function (params) {
                var query = {
                    search: params.term,
                    page: params.page || 1,
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                if ($('select#producer').select2('data').length != 0) {
                    query.producer = $('select#producer').select2('data')[0].id;
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;

                isExactMatch = false;

                // looking for an exact match
                $.each(data.rows, function (index, value) {
                    if (value.c == 1) {
                        isExactMatch = true;
                    }
                });

                // there is a match: setting the input field
                if (isExactMatch) {
                    $("input#exactMatchProducerRef").val("true");
                } else {
                    $("input#exactMatchProducerRef").val("false");
                }

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.producerref_label.String;
                    obj.id = obj.id || obj.producerref_id.Int64;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-producerref-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
                };
            }
        }
    }).on("change", function (e) {
        $(this).valid(); // FIXME: see https://github.com/select2/select2/issues/3901
    });

    //
    // supplierref select2
    //
    $('select#supplierrefs').select2({
        tags: true,
        templateResult: formatSupplierRef,
        templateSelection: formatSupplierRef2,
        placeholder: "select a reference",
        createTag: function (params) {
            if ($("input#exactMatchSupplierRef").val() == "true") {
                return null
            }
            if ($('select#supplier').select2('data').length == 0) {
                gjsUtils.message("to create a new reference select a supplier first", "warning");
                return null
            }
            return {
                id: params.term,
                text: params.term,
                supplierlabel: $('select#supplier').select2('data')[0].supplier_label.String,
                supplierid: $('select#supplier').select2('data')[0].supplier_id.Int64,
            }
        },
        ajax: {
            url: proxyPath + 'products/supplierrefs/',
            delay: 400,
            data: function (params) {
                sid = $('select#supplier').select2('data').length != 0 && ($('select#supplier').select2('data')[0].id != $('select#supplier').select2('data')[0].text) ? $('select#supplier').select2('data')[0].id : -1;
                var query = {
                    supplier: sid,
                    search: params.term,
                    page: params.page || 1,
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;

                isExactMatch = false;

                // looking for an exact match
                $.each(data.rows, function (index, value) {
                    if (value.c == 1) {
                        isExactMatch = true;
                    }
                });

                // there is a match: setting the input field
                if (isExactMatch) {
                    $("input#exactMatchSupplierRef").val("true");
                } else {
                    $("input#exactMatchSupplierRef").val("false");
                }

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.supplierref_label;
                    obj.id = obj.id || obj.supplierref_id;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-supplierrefs-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
                };
            }
        }
    }).on("change", function (e) {
        $(this).valid(); // FIXME: see https://github.com/select2/select2/issues/3901
    });

    //
    // supplier select2
    //
    $('select#supplier').select2({
        //tags: true,
        placeholder: "select a supplier",
        allowClear: true,
        // createTag: function (params) {
        //     if ($("input#exactMatchSupplier").val() == "true") {
        //         return null
        //     }
        //     return {
        //         id: params.term,
        //         text: params.term,
        //     }
        // },
        ajax: {
            url: proxyPath + 'products/suppliers/',
            delay: 400,
            data: function (params) {
                var query = {
                    search: params.term,
                    page: params.page || 1,
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;

                isExactMatch = false;

                // looking for an exact match
                $.each(data.rows, function (index, value) {
                    if (value.c == 1) {
                        isExactMatch = true;
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
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
                };
            }
        }
    }).on("change", function (e) {
        $(this).valid(); // FIXME: see https://github.com/select2/select2/issues/3901
    });

    //
    // producer select2
    //
    $('select#producer').select2({
        //tags: true,
        placeholder: "select a producer",
        allowClear: true,
        // createTag: function (params) {
        //     if ($("input#exactMatchProducer").val() == "true") {
        //         return null
        //     }
        //     return {
        //         id: params.term,
        //         text: params.term,
        //     }
        // },
        ajax: {
            url: proxyPath + 'products/producers/',
            delay: 400,
            data: function (params) {
                var query = {
                    search: params.term,
                    page: params.page || 1,
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;

                isExactMatch = false;

                // looking for an exact match
                $.each(data.rows, function (index, value) {
                    if (value.c == 1) {
                        isExactMatch = true;
                    }
                });

                // there is a match: setting the input field
                if (isExactMatch) {
                    $("input#exactMatchProducer").val("true");
                } else {
                    $("input#exactMatchProducer").val("false");
                }

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.producer_label.String;
                    obj.id = obj.id || obj.producer_id.Int64;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-producer-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
                };
            }
        }
    }).on("change", function (e) {
        $(this).valid(); // FIXME: see https://github.com/select2/select2/issues/3901
    });
    
    //
    // tags select2
    //
    $('select#tags').select2({
        tags: true,
        placeholder: "select or enter tag(s)",
        createTag: function (params) {
            if ($("input#exactMatchTags").val() == "true") {
                return null
            }
            return {
                id: params.term,
                text: params.term,
            }
        },
        ajax: {
            url: proxyPath + 'products/tags/',
            delay: 400,
            data: function (params) {
                var query = {
                    search: params.term,
                    page: params.page || 1,
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;

                isExactMatch = false;

                // looking for an exact match
                $.each(data.rows, function (index, value) {
                    if (value.c == 1) {
                        isExactMatch = true;
                    }
                });

                // there is a match: setting$('select#producerref').select2('data')[0].producer the input field
                if (isExactMatch) {
                    $("input#exactMatchTags").val("true");
                } else {
                    $("input#exactMatchTags").val("false");
                }

                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.tag_label;
                    obj.id = obj.id || obj.tag_id;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-tags-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
                };
            }
        }
    });

    //
    // category select2
    //
    $('select#category').select2({
        tags: true,
        allowClear: true,
        placeholder: "select a category",
        ajax: {
            url: proxyPath + 'products/categories/',
            delay: 400,
            data: function (params) {
                var query = {
                    search: params.term,
                    page: params.page || 1,
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;
                // replacing name by text expected by select2
                var newdata = $.map(data.rows, function (obj) {
                    obj.text = obj.text || obj.category_label.String;
                    obj.id = obj.id || obj.category_id.Int64;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-category-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
                };
            }
        }
    }).on("change", function (e) {
        $(this).valid(); // FIXME: see https://github.com/select2/select2/issues/3901
    });

    //
    // unit_temperature select2
    //
    $('select#unit_temperature').select2({
        placeholder: "select a unit",
        ajax: {
            url: proxyPath + 'storages/units',
            delay: 400,
            data: function (params) {
                var query = {
                    search: params.term,
                    unit_type: "temperature",
                    page: params.page || 1,
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;
                // replacing name by text expected by select2
                // var newdata = $.map(data.rows, function (obj) {
                //     obj.text = obj.text || obj.unit_label.String;
                //     obj.id = obj.id || obj.unit_id.Int64;
                //     return obj;
                // });

                const concentration = new Array(); 
                const temperature = new Array(); 
                const volume = new Array();
                const weight = new Array();
                const length = new Array();
                for (const i in data.rows) {
                    t = data.rows[i].unit_type.String;
                    if (t == "concentration") {
                        concentration.push({
                            "id": data.rows[i].unit_id.Int64,
                            "text": data.rows[i].unit_label.String,
                        })
                    }
                    if (t == "volume") {
                        volume.push({
                            "id": data.rows[i].unit_id.Int64,
                            "text": data.rows[i].unit_label.String,
                        })
                    }
                    if (t == "weight") {
                        weight.push({
                            "id": data.rows[i].unit_id.Int64,
                            "text": data.rows[i].unit_label.String,
                        })
                    }
                    if (t == "length") {
                        length.push({
                            "id": data.rows[i].unit_id.Int64,
                            "text": data.rows[i].unit_label.String,
                        })
                    }
                    if (t == "temperature") {
                        temperature.push({
                            "id": data.rows[i].unit_id.Int64,
                            "text": data.rows[i].unit_label.String,
                        })
                    }
                }

                var newdata = new Array()
                if (volume.length != 0) {
                    newdata.push({
                        "text": "volume",
                        "children": volume,
                    })
                }
                if (length.length != 0) {
                    newdata.push({
                        "text": "length",
                        "children": length,
                    })
                }
                if (weight.length != 0) {
                    newdata.push({
                        "text": "weight",
                        "children": weight,
                    })
                }
                if (concentration.length != 0) {
                    newdata.push({
                        "text": "concentration",
                        "children": concentration,
                    })
                }
                if (temperature.length != 0) {
                    newdata.push({
                        "text": "temperature",
                        "children": temperature,
                    })
                }

                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-unit_temperature-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
                };
            }
        }
    }).on("change", function (e) {
        $(this).valid(); // FIXME: see https://github.com/select2/select2/issues/3901
    });

    //
    // casnumber select2
    //
    $('select#casnumber').select2({
        tags: true,
        placeholder: gjsUtils.translate("product_cas_placeholder", container.PersonLanguage),
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
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;

                isExactMatch = false;

                // looking for an exact match
                $.each(data.rows, function (index, value) {
                    if (value.c == 1) {
                        isExactMatch = true;
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
                    obj.text = obj.text || obj.casnumber_label.String;
                    obj.id = obj.id || obj.casnumber_id.Int64;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-casnumber-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
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
        placeholder: gjsUtils.translate("product_ce_placeholder", container.PersonLanguage),
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
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;
                isExactMatch = false;

                // looking for an exact match
                $.each(data.rows, function (index, value) {
                    if (value.c == 1) {
                        isExactMatch = true;
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
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
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
        placeholder: gjsUtils.translate("product_physicalstate_placeholder", container.PersonLanguage),
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
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;
                isExactMatch = false;

                // looking for an exact match
                $.each(data.rows, function (index, value) {
                    if (value.c == 1) {
                        isExactMatch = true;
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
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
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
        placeholder: gjsUtils.translate("product_signalword_placeholder", container.PersonLanguage),
        ajax: {
            url: proxyPath + 'products/signalwords/',
            delay: 400,
            data: function (params) {
                var query = {
                    search: params.term,
                    page: params.page || 1,
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;
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
        placeholder: gjsUtils.translate("product_classofcompound_placeholder", container.PersonLanguage),
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
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;
                isExactMatch = false;

                // looking for an exact match
                $.each(data.rows, function (index, value) {
                    if (value.c == 1) {
                        isExactMatch = true;
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
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
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
        placeholder: gjsUtils.translate("product_name_placeholder", container.PersonLanguage),
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
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;

                isExactMatch = false;

                // looking for an exact match
                $.each(data.rows, function (index, value) {
                    if (value.c == 1) {
                        isExactMatch = true;
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
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
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
        placeholder: gjsUtils.translate("product_empiricalformula_placeholder", container.PersonLanguage),
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
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;

                isExactMatch = false;

                // looking for an exact match
                $.each(data.rows, function (index, value) {
                    if (value.c == 1) {
                        isExactMatch = true;
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
                    obj.text = obj.text || obj.empiricalformula_label.String;
                    obj.id = obj.id || obj.empiricalformula_id.Int64;
                    return obj;
                });
                // getting the number of loaded select elements
                selectnbitems = $("ul#select2-empiricalformula-results li").length + 10;

                return {
                    results: newdata,
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
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
        placeholder: gjsUtils.translate("product_linearformula_placeholder", container.PersonLanguage),
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
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;

                isExactMatch = false;

                // looking for an exact match
                $.each(data.rows, function (index, value) {
                    if (value.c == 1) {
                        isExactMatch = true;
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
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
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
        placeholder: gjsUtils.translate("product_synonyms_placeholder", container.PersonLanguage),
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
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;

                isExactMatch = false;

                // looking for an exact match
                $.each(data.rows, function (index, value) {
                    if (value.c == 1) {
                        isExactMatch = true;
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
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
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
        placeholder: gjsUtils.translate("product_symbols_placeholder", container.PersonLanguage),
        ajax: {
            url: proxyPath + 'products/symbols/',
            delay: 400,
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;
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
        placeholder: gjsUtils.translate("product_hazardstatements_placeholder", container.PersonLanguage),
        ajax: {
            url: proxyPath + 'products/hazardstatements/',
            delay: 400,
            data: function (params) {
                var query = {
                    search: params.term,
                    page: params.page || 1,
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;
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
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
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
        placeholder: gjsUtils.translate("product_precautionarystatements_placeholder", container.PersonLanguage),
        ajax: {
            url: proxyPath + 'products/precautionarystatements/',
            delay: 400,
            data: function (params) {
                var query = {
                    search: params.term,
                    page: params.page || 1,
                    offset: (params.page - 1) * 10 || 0,
                    limit: 10
                }

                // Query parameters will be ?search=[term]&page=[page]
                return query;
            },
            dataType: 'json',
            processResults: function (data, params) {
                params.page = params.page || 1;
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
                    pagination: {
                        more: (params.page * 10) < data.total
                    }
                };
            }
        }
    });


    if ($("form#product input#product_id").length) {

    } else {
        // create
        $("input#showchem").attr("checked", true)
        chemfy()
    }
   
});

function getData(params) {
    // saving the query parameters
    lastQueryParams = params;
    $.ajax({
        url: proxyPath + "products",
        method: "GET",
        dataType: "JSON",
        data: params.data,
    }).done(function (data, textStatus, jqXHR) {
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
    }).fail(function (jqXHR, textStatus, errorThrown) {
        params.error(jqXHR.statusText);
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });
}

//
// when table is loaded
//
$('#table').on('load-success.bs.table refresh.bs.table', function () {

    hasPermission("storages", "", "POST").done(function () {
        $(".store").fadeIn();
        localStorage.setItem("storages::POST", true);
    }).fail(function () {
        localStorage.setItem("storages::POST", false);
    })
    hasPermission("storages", "-2", "GET").done(function () {
        $("#switchview").removeClass("d-none");

        $(".storages").fadeIn();
        $(".ostorages").fadeIn();
        localStorage.setItem("storages:-2:GET", true);
    }).fail(function () {
        localStorage.setItem("storages:-2:GET", false);
    })
    hasPermission("products", "-2", "PUT").done(function () {
        $(".edit").fadeIn();
        localStorage.setItem("products:-2:PUT", true);
    }).fail(function () {
        localStorage.setItem("products:-2:PUT", false);
    })

    $("table#table").find("tr[product_id]").each(function (index, b) {
        hasPermission("products", $(b).attr("product_id"), "DELETE").done(function () {
            $("#delete" + $(b).attr("product_id")).fadeIn();
            localStorage.setItem("products:" + $(b).attr("product_id") + ":DELETE", true);
        }).fail(function () {
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
    return { "product_id": row["product_id"] }
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
        html.push("<div class='col-sm-6'>")
            if (row["category"]["category_id"]["Valid"]) {
                html.push("<span class='iconlabel'>" + gjsUtils.translate("category_label_title", container.PersonLanguage) + "</span> " + row["category"]["category_label"]["String"])
            }
        html.push("</div>")
        html.push("<div class='col-sm-6'>")
            $.each(row["tags"], function (key, value) {
                html.push("<span class='mdi mdi-code-tags'>" + value["tag_label"] + "</span> ");
                html.push("<br/>");
            });
        html.push("</div>")
    html.push("</div>")

    html.push("<div class='row mt-sm-3'>")
        html.push("<div class='col-sm-6'>")
            if (row["producerref"]["producerref_id"]["Valid"]) {
                html.push("<span class='iconlabel'>" + gjsUtils.translate("producerref_label_title", container.PersonLanguage) + "</span> " + row["producerref"]["producerref_label"]["String"] + " <i>(" + row["producerref"]["producer"]["producer_label"]["String"] + "</i>)")
            }
        html.push("</div>")
        html.push("<div class='col-sm-6'>")
            if (row["supplierrefs"] != null) {
                html.push("<span class='iconlabel'>" + gjsUtils.translate("supplierref_label_title", container.PersonLanguage) + "</span><br/>")
                $.each(row["supplierrefs"], function (key, value) {
                    html.push("<span>" + value["supplierref_label"] + "</span> " + " <i>(" + value["supplier"]["supplier_label"]["String"] + "</i>)");
                    html.push("<br/>");
                });
            }
        html.push("</div>")
    html.push("</div>")
  
    html.push("<div class='row mt-sm-3'>")
        html.push("<div class='col-sm-4'>")
            if (row["casnumber"]["casnumber_id"]["Valid"]) {
                html.push("<span class='iconlabel'>" + gjsUtils.translate("casnumber_label_title", container.PersonLanguage) + "</span> " + row["casnumber"]["casnumber_label"]["String"])
            }
            if (row["casnumber"]["casnumber_cmr"]["Valid"]) {
                html.push("&nbsp;<span class='iconlabel'>" + gjsUtils.translate("casnumber_cmr_title", container.PersonLanguage) + "</span> " + row["casnumber"]["casnumber_cmr"]["String"])
            }
            $.each(row.hazardstatements, function (index, value) {
                if (value.hazardstatement_cmr.Valid) {
                    html.push("&nbsp;" + value.hazardstatement_cmr.String);
                }
            });
        html.push("</div>")
        html.push("<div class='col-sm-4'>")
            if (row["cenumber"]["cenumber_id"]["Valid"]) {
                html.push("<span class='iconlabel'>" + gjsUtils.translate("cenumber_label_title", container.PersonLanguage) + "</span> " + row["cenumber"]["cenumber_label"]["String"])
            }
        html.push("</div>")
        html.push("<div class='col-sm-4'>")
            if (row["product_msds"]["Valid"]) {
                html.push("<span class='iconlabel'>" + gjsUtils.translate("product_msds_title", container.PersonLanguage) + "</span> <a href='" + row["product_msds"]["String"] + "'><span class='mdi mdi-link-variant mdi-24px'></span></a>")
            }
        html.push("</div>")
    html.push("</div>")


    html.push("<div class='row mt-sm-3'>")
        if (row["empiricalformula"]["empiricalformula_id"]["Valid"]) {
            html.push("<div class='col-sm-4'><span class='iconlabel'>" + gjsUtils.translate("empiricalformula_label_title", container.PersonLanguage) + "</span> " + row["empiricalformula"]["empiricalformula_label"]["String"] + "</div>")
        }
        if (row["linearformula"]["linearformula_id"]["Valid"]) {
            html.push("<div class='col-sm-4'><span class='iconlabel'>" + gjsUtils.translate("linearformula_label_title", container.PersonLanguage) + "</span> " + row["linearformula"]["linearformula_label"]["String"] + "</div>")
        }
        if (row["product_threedformula"]["Valid"] && row["product_threedformula"]["String"] != "") {
            html.push("<div class='col-sm-4'><span class='iconlabel'>" + gjsUtils.translate("product_threedformula_title", container.PersonLanguage) + "</span> <a href='" + row["product_threedformula"]["String"] + "'><span class='mdi mdi-link-variant mdi-24px'></span></a></div>")
        }
    html.push("</div>")

    // html.push("<div class='row mt-sm-3'>")
    //     html.push("<div class='col-sm-6'>")
    //         if (row["product_batchnumber"]["Valid"]) {
    //             html.push("<span class='iconlabel'>" + gjsUtils.translate("product_batchnumber_title", container.PersonLanguage) + "</span> " + row["product_batchnumber"]["String"])
    //         }
    //     html.push("</div>")
    //     html.push("<div class='col-sm-6'>")
    //         if (row["product_expirationdate"]["Valid"]) {

    //             date = new Date(row["product_expirationdate"]["Time"]);

    //             html.push("<span class='iconlabel'>" + gjsUtils.translate("product_expirationdate_title", container.PersonLanguage) + "</span> " + date.toLocaleDateString())
    //         }
    //     html.push("</div>")
    // html.push("</div>")

    html.push("<div class='row mt-sm-3'>")
        html.push("<div class='col-sm-6'>")
            // if (row["product_concentration"]["Valid"]) {
            //     html.push("<span class='iconlabel'>" + gjsUtils.translate("product_concentration_title", container.PersonLanguage) + "</span> " + row["product_concentration"]["Int64"])
            // }
            // if (row["unit_concentration"]["unit_id"]["Valid"]) {
            //     html.push(row["unit_concentration"]["unit_label"]["String"])
            // }
        html.push("</div>")
        html.push("<div class='col-sm-6'>")
            if (row["product_temperature"]["Valid"]) {
                html.push("<span class='iconlabel'>" + gjsUtils.translate("product_temperature_title", container.PersonLanguage) + "</span> " + row["product_temperature"]["Int64"])
            }
            if (row["unit_temperature"]["unit_id"]["Valid"]) {
                html.push(row["unit_temperature"]["unit_label"]["String"])
            }
        html.push("</div>")
    html.push("</div>")


    html.push("<div class='row mt-sm-3'>")

    html.push("<div class='col-sm-4'>")
    $.each(row["symbols"], function (key, value) {
        html.push("<img src='data:" + value["symbol_image"] + "' alt='" + value["symbol_label"] + "' title='" + value["symbol_label"] + "'/>");
    });
    html.push("</div>")
    html.push("<div class='col-sm-4'>")
    if (row["signalword"]["signalword_label"]["Valid"]) {
        html.push("<span class='iconlabel'>" + gjsUtils.translate("signalword_label_title", container.PersonLanguage) + "</span> " + row["signalword"]["signalword_label"]["String"])
    }
    html.push("</div>")
    html.push("<div class='col-sm-4'>")
    if (row["physicalstate"]["physicalstate_id"]["Valid"]) {
        html.push("<div class='col-sm-4'><span class='iconlabel'>" + gjsUtils.translate("physicalstate_label_title", container.PersonLanguage) + "</span> " + row["physicalstate"]["physicalstate_label"]["String"] + "</div>")
    }
    html.push("</div>")
    html.push("</div>")


    html.push("<div class='row mt-sm-3'>")

    html.push("<div class='col-sm-4'>")
    if (row["hazardstatements"] != null && row["hazardstatements"].length != 0) {
        html.push("<div><span class='iconlabel'>" + gjsUtils.translate("hazardstatement_label_title", container.PersonLanguage) + "</span></div>")
        html.push("<ul>")
        $.each(row["hazardstatements"], function (key, value) {
            html.push("<li>" + value["hazardstatement_reference"] + ": <i>" + value["hazardstatement_label"] + "</i></li>");
        });
        html.push("</ul>")
    }
    html.push("</div>")

    html.push("<div class='col-sm-4'>")
    if (row["precautionarystatements"] != null && row["precautionarystatements"].length != 0) {
        html.push("<div><span class='iconlabel'>" + gjsUtils.translate("precautionarystatement_label_title", container.PersonLanguage) + "</span></div>")
        html.push("<ul>")
        $.each(row["precautionarystatements"], function (key, value) {
            html.push("<li>" + value["precautionarystatement_reference"] + ": <i>" + value["precautionarystatement_label"] + "</i></li>");
        });
        html.push("</ul>")
    }
    html.push("</div>")

    html.push("<div class='col-sm-4'>")
    if (row["classofcompound"] != null && row["classofcompound"].length != 0) {
        html.push("<div><span class='iconlabel'>" + gjsUtils.translate("classofcompound_label_title", container.PersonLanguage) + "</span></div>")
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
        html.push("<div class='col-sm-12'><span class='iconlabel'>" + gjsUtils.translate("product_disposalcomment_title", container.PersonLanguage) + "</span> " + row["product_disposalcomment"]["String"] + "</div>")
    }

    html.push("</div>")

    html.push("<div class='row mt-sm-3'>")

    if (row["product_remark"]["Valid"] && row["product_remark"]["String"] != "") {
        html.push("<div class='col-sm-12'><span class='iconlabel'>" + gjsUtils.translate("product_remark_title", container.PersonLanguage) + "</span> " + row["product_remark"]["String"] + "</div>")
    }

    html.push("</div>")


    html.push("<div class='row mt-sm-3'>")

    html.push("<div class='col-sm-12'>")

    if (row["product_radioactive"]["Bool"]) {
        html.push("<span title='" + gjsUtils.translate("product_radioactive_title", container.PersonLanguage) + "' class='mdi mdi-36px mdi-radioactive'></span>")
    }
    if (row["product_restricted"]["Bool"]) {
        html.push("<span title='" + gjsUtils.translate("product_restricted_title", container.PersonLanguage) + "' class='mdi mdi-36px mdi-hand'></span>")
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
// casnumberFormatter formatter
//
function casnumberFormatter(value, row, index, field) {
    if (value.Valid) {
        return value.String;
    } else {
        return "";
    }
}
//
// empiricalformulaFormatter formatter
//
function empiricalformulaFormatter(value, row, index, field) {
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

    actions.push('<div class="float-right" style="position: relative">');

    $.each(row.symbols, function (index, value) {
        if (value.symbol_label == "SGH02") {
            actions.push('<img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmCAYAAACoPemuAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAN1wAADdcBQiibeAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAIvSURBVFiFzdgxbI5BGMDx36uNJsJApFtFIh2QEIkBtYiYJFKsrDaJqYOofhKRMFhsFhNNMBgkTAaD0KUiBomN1EpYBHWGnvi8vu9r7/2eti65vHfv3T3P/57nuXvfuyqlJCRVVQuk1AqRl1LqP9NKpJxbETKjocLgYi0VaLl49wXBxUIFwsVDBcGFQWEAg1FwYZbCGMajLBfmPkzgUZRbw2IKFzGPrRFw/bpvD/bn8jUkXM719f3A9eu+k3iXA/92Bnub2yYx1NgDfbrvXIYZx8dcThjBExxvPOmGltqLIzmuEt63QSVczc+z/2whSw2ThpbajS+4UgOq59O4gYFSuGaByWb8zKvwN8RXXKiBPc7PLaWx3ARqY37O1CBe5/cvO1huVy+ZnfSX7y9MYxRTNeX32lZj+/sXWNfVnV3g1tT/aJeQ5vAGp3L9eXbjTFv7NzzM9VncSSnNF2lp4MqjNYvcxwEcy+0HcQg32/q8Kndl+YrcgM9Z4YdsrZ21PtvxHT9yv1vNgr8cbiIrnMUmbKu177PwVZjLgKPNt4sCOKzF0ww32aF9CA+yxSZKoTqDlVnucI6lMxhpg76OuxhrKr8oIENyXx/xxQKTE/hUkIdLJ1tlRd3TwtF/KtcuSalVVdUwdvQe+Fd6ljhfl9NzRKT5I8cvq/B+xi3vzFfk+FaqbEUPvEtVuipXBIspX9VLlW4Q/8U1VGe4EKgYsED3tefBgt271y7dUlV/ygHpF8bRglXiwx7BAAAAAElFTkSuQmCC" alt="flammable" title="flammable">');
        }
    });

    if (row.product_sc != 0 || row.product_asc != 0) {
        actions.push('<button id="storages' + pid + '" class="storages btn btn-link btn-sm" style="display: none;" title="storages" type="button">',
            '<span class="mdi mdi-24px mdi-cube-unfolded"><i>' + row.product_sc + ' (' + row.product_asc + ')</i></span>',
            '</button>');
    } else {
        actions.push('<button disabled id="storages' + pid + '" class="storages btn btn-link btn-sm" style="display: none;" title="storages" type="button">',
            '<span class="mdi mdi-24px mdi-cube-unfolded"><i>' + row.product_sc + '</i></span>',
            '</button>');
    }

    if (row.product_tsc != 0) {
        actions.push('<button id="ostorages' + pid + '" class="ostorages btn btn-link btn-sm" style="display: none;" title="global availability" type="button">',
            '<span class="mdi mdi-24px mdi-cube-scan">',
            '</button>');
    } else {
        actions.push('<button disabled id="ostorages' + pid + '" class="ostorages btn btn-link btn-sm" style="display: none;" title="global availability" type="button">',
            '<span class="mdi mdi-24px mdi-cube-scan">',
            '</button>');
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

    actions.push('<div class="position-absolute" style="right: 0px; bottom: -8px;">');
    if (row.casnumber.casnumber_cmr.Valid) {
        actions.push('<span title="CMR" class="text-danger font-italic">' + row.casnumber.casnumber_cmr.String + '</span>');
    }
    $.each(row.hazardstatements, function (index, value) {
        if (value.hazardstatement_cmr.Valid) {
            actions.push('<span title="CMR" class="text-danger font-italic">' + value.hazardstatement_cmr.String + '</span>');
        }
    });
    actions.push('</div>');
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

        // ask for confirmation and then delete
        $("button#delete" + row.product_id).on("click", function () {
            $.ajax({
                url: proxyPath + "products/" + row['product_id'],
                method: "DELETE",
            }).done(function (data, textStatus, jqXHR) {
                gjsUtils.message(gjsUtils.translate("product_deleted_message", container.PersonLanguage), "success");
                var $table = $('#table');
                $table.bootstrapTable('refresh');
            }).fail(function (jqXHR, textStatus, errorThrown) {
                handleHTTPError(jqXHR.statusText, jqXHR.status)
            });
        }).css("color", "red").find("span.mdi-delete").removeClass("mdi-delete").addClass("mdi-delete-sweep").prop("title", gjsUtils.translate("confirm", container.PersonLanguage));
    }
};
function operateBookmark(e, value, row, index) {
    // toggling the bookmark
    $.ajax({
        url: proxyPath + "bookmarks/" + row['product_id'],
        method: "PUT",
    }).done(function (data, textStatus, jqXHR) {
        if ($("span#bookmark" + data.product_id).hasClass("mdi-bookmark")) {
            $("span#bookmark" + data.product_id).removeClass("mdi-bookmark");
            $("span#bookmark" + data.product_id).addClass("mdi-bookmark-outline");
        } else {
            $("span#bookmark" + data.product_id).removeClass("mdi-bookmark-outline");
            $("span#bookmark" + data.product_id).addClass("mdi-bookmark");
        }
    }).fail(function (jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });
}
function operateOStorages(e, value, row, index) {
    // getting the product
    $.ajax({
        url: proxyPath + "storages/others?product=" + row['product_id'],
        method: "GET",
    }).done(function (data, textStatus, jqXHR) {
        var html = [];
        $.each(data["rows"], function (key, value) {
            html.push("<p><span class='iconlabel'>" + value.entity_name + "</span><span class='blockquote-footer'>" + value.entity_description + "</span></p>");
        });

        $("#ostorages-collapse-" + row['product_id']).html(html.join('&nbsp;'));
    }).fail(function (jqXHR, textStatus, errorThrown) {
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

    // bio
    $('input#addproducer').val(null);
    $('input#addsupplier').val(null);
    // $('input#product_batchnumber').val(null);
    // $('input#product_concentration').val(null);
    // $('input#product_expirationdate').val(null);
    $('input#product_temperature').val(null);
    $('input#product_sheet').val(null);

    $('select#category').val(null).trigger('change');
    $('select#category').find('option').remove();
    $('select#tags').val(null).trigger('change');
    $('select#tags').find('option').remove();
    $('select#producerref').val(null).trigger('change');
    $('select#producerref').find('option').remove();
    $('select#producer').val(null).trigger('change');
    $('select#producer').find('option').remove();
    $('select#supplierrefs').val(null).trigger('change');
    $('select#supplierrefs').find('option').remove();
    $('select#supplier').val(null).trigger('change');
    $('select#supplier').find('option').remove();
    // $('select#unit_concentration').val(null).trigger('change');
    // $('select#unit_concentration').find('option').remove();
    $('select#unit_temperature').val(null).trigger('change');
    $('select#unit_temperature').find('option').remove();

    // getting the product
    $.ajax({
        url: proxyPath + "products/" + row['product_id'],
        method: "GET",
    }).done(function (data, textStatus, jqXHR) {
        // flattening response data
        fdata = flatten(data);

        // processing sqlNull values
        //newfdata = normalizeSqlNull(fdata)
        newfdata = gjsUtils.normalizeSqlNull(fdata);

        // autofilling form
        $("#edit-collapse").autofill(newfdata, { "findbyname": false });
        // setting index hidden input
        $("input#index").val(index);

        // select2 is not autofilled - we need a special operation
        var newOption = new Option(data.name.name_label, data.name.name_id, true, true);
        $('select#name').append(newOption).trigger('change');

        if (data.casnumber.casnumber_id.Valid) {
            var newOption = new Option(data.casnumber.casnumber_label.String, data.casnumber.casnumber_id.Int64, true, true);
            $('select#casnumber').append(newOption).trigger('change');
        }

        if (data.empiricalformula.empiricalformula_id.Valid) {
            var newOption = new Option(data.empiricalformula.empiricalformula_label.String, data.empiricalformula.empiricalformula_id.Int64, true, true);
            $('select#empiricalformula').append(newOption).trigger('change');
        }

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

        if (data.category.category_id.Valid) {
            var newOption = new Option(data.category.category_label.String, data.category.category_id.Int64, true, true);
            $('select#category').append(newOption).trigger('change');
        }

        // if (data.unit_concentration.unit_id.Valid) {
        //     var newOption = new Option(data.unit_concentration.unit_label.String, data.unit_concentration.unit_id.Int64, true, true);
        //     $('select#unit_concentration').append(newOption).trigger('change');
        // }

        if (data.unit_temperature.unit_id.Valid) {
            var newOption = new Option(data.unit_temperature.unit_label.String, data.unit_temperature.unit_id.Int64, true, true);
            $('select#unit_temperature').append(newOption).trigger('change');
        }

        if (data.producerref.producerref_id.Valid) {
            var newOption = new Option(data.producerref.producerref_label.String, data.producerref.producerref_id.Int64, true, true);
            $('select#producerref').append(newOption).trigger('change');

            var newOption = new Option(data.producerref.producer.producer_label.String, data.producerref.producer.producer_id.Int64, true, true);
            $('select#producer').append(newOption).trigger('change');
        }

        for (var i in data.supplierrefs) {
            var newOption = new Option(data.supplierrefs[i].supplier.supplier_label.String + ": " + data.supplierrefs[i].supplierref_label, data.supplierrefs[i].supplierref_id, true, true);
            $('select#supplierrefs').append(newOption).trigger('change');
        }

        for (var i in data.tags) {
            var newOption = new Option(data.tags[i].tag_label, data.tags[i].tag_id, true, true);
            $('select#tags').append(newOption).trigger('change');
        }

        for (var i in data.symbols) {
            var newOption = new Option(data.symbols[i].symbol_label, data.symbols[i].symbol_id, true, true);
            $('select#symbols').append(newOption).trigger('change');
        }

        for (var i in data.classofcompound) {
            var newOption = new Option(data.classofcompound[i].classofcompound_label, data.classofcompound[i].classofcompound_id, true, true);
            $('select#classofcompound').append(newOption).trigger('change');
        }

        for (var i in data.synonyms) {
            var newOption = new Option(data.synonyms[i].name_label, data.synonyms[i].name_id, true, true);
            $('select#synonyms').append(newOption).trigger('change');
        }

        for (var i in data.hazardstatements) {
            var newOption = new Option(data.hazardstatements[i].hazardstatement_reference, data.hazardstatements[i].hazardstatement_id, true, true);
            $('select#hazardstatements').append(newOption).trigger('change');
        }

        for (var i in data.precautionarystatements) {
            var newOption = new Option(data.precautionarystatements[i].precautionarystatement_reference, data.precautionarystatements[i].precautionarystatement_id, true, true);
            $('select#precautionarystatements').append(newOption).trigger('change');
        }

        // chem/bio detection
        if ($("input#product_batchnumber").val()) {
            biofy()
        } else {
            chemfy()
        }
    }).fail(function (jqXHR, textStatus, errorThrown) {
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
$('#product_molformula').change(function () {
    if (!window.FileReader) {
        return alert('FileReader API is not supported by your browser.');
    }
    molfile = $('#product_molformula')[0];
    if (molfile.files && molfile.files[0]) {
        file = molfile.files[0]; // The file
        fr = new FileReader(); // FileReader instance
        fr.onload = function () {
            $('#hidden_product_molformula_content').append(fr.result);
        };
        fr.readAsText(file);
    } else {
        // Handle errors here
        alert("file not selected or browser incompatible.")
    }
});

//
// store location select2 formatter
//
function formatStorelocation(sl) {
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
// producerref formatter
//
function formatProducerRef(pr) {
    if (!pr.producerref_id) {
        return pr.producerref_label;
    }
    if (pr.producerref_id.Valid) {
        return $("<span><b>" + pr.producer.producer_label.String + "</b>: " + pr.producerref_label.String + "</span>");
    }
};
function formatProducerRef2(pr) {
    if (!pr.producerref_id) {
        return pr.text;
    }
    if (pr.producerref_id.Valid) {
        return $("<span>" + pr.producer.producer_label.String + ": " + pr.producerref_label.String + "</span>");
    }
};

//
// supplierref formatter
//
function formatSupplierRef(sr) {
    if (!sr.supplierref_id) {
        return sr.supplierref_label;
    } else {
        return $("<span><b>" + sr.supplier.supplier_label.String + "</b>: " + sr.supplierref_label + "</span>");
    }
};
function formatSupplierRef2(sr) {
    if (!sr.supplier && (sr.text != sr.id)) {
        return $("<span>" + sr.text + "</span>");
    }
    if (!sr.supplierref_id) {
        return $("<span>" + sr.supplierlabel + ": " + sr.text + "</span>");
    } else {
        return $("<span>" + sr.supplier.supplier_label.String + ": " + sr.supplierref_label + "</span>");
    }
};

//
// signalwords select2 formatter
//
function formatSignalWord(signalword) {
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
function formatSymbol(symbol) {
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
function formatPrecautionaryStatement(ps) {
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
function formatHazardStatement(hs) {
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
    gjsUtils.message(gjsUtils.translate("product_created_message", container.PersonLanguage) + ": " + data.name.name_label, "success");
    setTimeout(function () { window.location = proxyPath + "v/products?product=" + data.product_id + "&hl=" + data.product_id; }, 1000);
}
var updateCallBack = function updateCallback(data, textStatus, jqXHR) {
    gjsUtils.message(gjsUtils.translate("product_updated_message", container.PersonLanguage) + ": " + data.name.name_label, "success");
    setTimeout(function () { window.location = proxyPath + "v/products?product=" + data.product_id + "&hl=" + data.product_id; }, 1000);
}
function saveProduct() {
    var form = $("#product");
    if (!form.valid()) {
        return;
    };

    var product_id = $("input#product_id").val(),
        // product_batchnumber = $("input#product_batchnumber").val(),
        // product_concentration = $("input#product_concentration").val(),
        // product_expirationdate = $("input#product_expirationdate").val(),
        product_temperature = $("input#product_temperature").val(),
        product_specificity = $("input#product_specificity").val(),
        product_threedformula = $("input#product_threedformula").val(),

        product_molformula = $("#hidden_product_molformula_content").html(),

        product_sheet = $("input#product_sheet").val(),
        product_msds = $("input#product_msds").val(),
        product_disposalcomment = $("textarea#product_disposalcomment").val(),
        product_remark = $("textarea#product_remark").val(),
        product_restricted = $("input#product_restricted:CHECKED").val(),
        product_radioactive = $("input#product_radioactive:CHECKED").val(),

        //unit_concentration = $('select#unit_concentration').select2('data')[0],
        unit_temperature = $('select#unit_temperature').select2('data')[0],
        casnumber = $('select#casnumber').select2('data')[0],
        cenumber = $('select#cenumber').select2('data')[0],
        empiricalformula = $('select#empiricalformula').select2('data')[0],
        linearformula = $('select#linearformula').select2('data')[0],
        name_ = $('select#name').select2('data')[0],
        physicalstate = $('select#physicalstate').select2('data')[0],
        signalword = $('select#signalword').select2('data')[0],
        category = $('select#category').select2('data')[0],
        producerref = $('select#producerref').select2('data')[0],
        classofcompound = $('select#classofcompound').select2('data'),
        synonyms = $('select#synonyms').select2('data'),
        symbols = $('select#symbols').select2('data'),
        hazardstatements = $('select#hazardstatements').select2('data'),
        precautionarystatements = $('select#precautionarystatements').select2('data'),
        tags = $('select#tags').select2('data'),
        supplierrefs = $('select#supplierrefs').select2('data'),
        ajax_url = proxyPath + "products",
        ajax_method = "POST",
        ajax_callback = createCallBack,
        data = {};

    if ($("form#product input#product_id").length != 0) {
        ajax_url = proxyPath + "products/" + product_id
        ajax_method = "PUT"
        ajax_callback = updateCallBack
    }

    $.each(symbols, function (index, s) {
        data["symbols." + index + ".symbol_id"] = s.id;
        data["symbols." + index + ".symbol_label"] = s.text;
    });
    $.each(synonyms, function (index, s) {
        data["synonyms." + index + ".name_id"] = s.id == s.text ? -1 : s.id;
        data["synonyms." + index + ".name_label"] = s.text;
    });
    $.each(classofcompound, function (index, s) {
        data["classofcompound." + index + ".classofcompound_id"] = s.id == s.text ? -1 : s.id;
        data["classofcompound." + index + ".classofcompound_label"] = s.text;
    });
    $.each(hazardstatements, function (index, s) {
        data["hazardstatements." + index + ".hazardstatement_id"] = s.id;
        data["hazardstatements." + index + ".hazardstatement_reference"] = s.text;
    });
    $.each(precautionarystatements, function (index, s) {
        data["precautionarystatements." + index + ".precautionarystatement_id"] = s.id;
        data["precautionarystatements." + index + ".precautionarystatement_reference"] = s.text;
    });
    $.each(tags, function (index, s) {
        data["tags." + index + ".tag_id"] = s.id == s.text ? -1 : s.id;
        data["tags." + index + ".tag_label"] = s.text;
    });
    $.each(supplierrefs, function (index, s) {
        data["supplierrefs." + index + ".supplierref_id"] = s.id == s.text ? -1 : s.id;
        data["supplierrefs." + index + ".supplierref_label"] = s.text;
        data["supplierrefs." + index + ".supplier.supplier_label"] = s.supplierlabel;
        data["supplierrefs." + index + ".supplier.supplier_id"] = s.supplierid;
    });
    $.extend(data, {
        "product_id": product_id,
        "product_disposalcomment": product_disposalcomment,
        "product_remark": product_remark,
        "product_restricted": product_restricted == "on" ? true : false,
        "product_radioactive": product_radioactive == "on" ? true : false,
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
    // if (product_batchnumber !== "") {
    //     $.extend(data, {
    //         "product_batchnumber": product_batchnumber,
    //     });
    // }
    // if (product_concentration !== "") {
    //     $.extend(data, {
    //         "product_concentration": product_concentration,
    //     });
    // }
    // if (product_expirationdate !== "") {
    //     $.extend(data, {
    //         "product_expirationdate": product_expirationdate,
    //     });
    // }
    if (product_temperature !== "") {
        $.extend(data, {
            "product_temperature": product_temperature,
        });
    }
    if (product_sheet !== "") {
        $.extend(data, {
            "product_sheet": product_sheet,
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
    // if (unit_concentration !== undefined) {
    //     $.extend(data, {
    //         "unit_concentration.unit_id": unit_concentration.id == unit_concentration.text ? -1 : unit_concentration.id,
    //         "unit_concentration.unit_label": unit_concentration === undefined ? "" : unit_concentration.text,
    //     });
    // }
    if (unit_temperature !== undefined) {
        $.extend(data, {
            "unit_temperature.unit_id": unit_temperature.id == unit_temperature.text ? -1 : unit_temperature.id,
            "unit_temperature.unit_label": unit_temperature === undefined ? "" : unit_temperature.text,
        });
    }
    if (category !== undefined) {
        $.extend(data, {
            "category.category_id": category.id == category.text ? -1 : category.id,
            "category.category_label": category === undefined ? "" : category.text,
        });
    }
    if (producerref !== undefined) {
        $.extend(data, {
            "producerref.producerref_id": producerref.id == producerref.text ? -1 : producerref.id,
            "producerref.producerref_label": producerref === undefined ? "" : producerref.text,
            "producerref.producer.producer_id": producerref.producerid,
            "producerref.producer.producer_label": producerref.producerlabel,
        });
    }
    if (casnumber !== undefined) {
        $.extend(data, {
            "casnumber.casnumber_id": casnumber.id == casnumber.text ? -1 : casnumber.id,
            "casnumber.casnumber_label": casnumber === undefined ? "" : casnumber.text,
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
    if (empiricalformula !== undefined) {
        $.extend(data, {
            "empiricalformula.empiricalformula_id": empiricalformula.id == empiricalformula.text ? -1 : empiricalformula.id,
            "empiricalformula.empiricalformula_label": empiricalformula.text,
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
    }).done(ajax_callback).fail(function (jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });
}

//
// magical selector
//
function magic() {
    //console.log("magic")
    var m = $('textarea#magical').val();
    if (m != "") {
        $.ajax({
            url: proxyPath + "products/magic",
            method: "POST",
            data: { "msds": m }
        }).done(function (data, textStatus, jqXHR) {
            $('select#hazardstatements').val(null).trigger('change');
            $('select#hazardstatements').find('option').remove();
            $('select#precautionarystatements').val(null).trigger('change');
            $('select#precautionarystatements').find('option').remove();

            var hs = data.hs,
                ps = data.ps;
            for (var i = 0; i < hs.length; i++) {
                var newOption = new Option(data.hs[i].hazardstatement_reference, data.hs[i].hazardstatement_id, true, true);
                $('select#hazardstatements').append(newOption).trigger('change');
            }
            for (var i = 0; i < ps.length; i++) {
                var newOption = new Option(data.ps[i].precautionarystatement_reference, data.ps[i].precautionarystatement_id, true, true);
                $('select#precautionarystatements').append(newOption).trigger('change');
            }
        }).fail(function (jqXHR, textStatus, errorThrown) {
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
        }).done(function (data, textStatus, jqXHR) {
            $("#fconverter").attr("data-content", data);
            $("#fconverter").popover('show');
        }).fail(function (jqXHR, textStatus, errorThrown) {
            handleHTTPError(jqXHR.statusText, jqXHR.status)
        });
    }
}

//
// zero empirical formula selector
//
function noEmpiricalFormula() {
    $("select#empiricalformula").rules("remove", "required");
    $("span#empiricalformula.badge").remove();
}

//
// zero CAS number selector
//
function noCASNumber() {
    $("select#casnumber").rules("remove", "required");
    $("span#casnumber.badge").remove();
}

//
// magical selector how to
//
function howToMagicalSelector() {
    window.open(proxyPath + "img/magicalselector.webm", '_blank');
}

//
// hide/show chem/bio
//
function chemfy() {
    $(".chem").show()
    $(".bio").hide()

    $("select#producerref").rules("remove", "required");
    $("input#product_batchnumber").rules("remove", "required");

    $("select#empiricalformula").rules("add", "required");
    $("select#casnumber").rules("add", "required");
}
function biofy() {
    $(".chem").hide()
    $(".bio").show()

    $("select#producerref").rules("add", "required");
    $("input#product_batchnumber").rules("add", "required");

    $("select#empiricalformula").rules("remove", "required");
    $("select#casnumber").rules("remove", "required");
}

//
// add supplier
//
function addSupplier() {
    data = {};
    data["supplier_label"] = $("input#addsupplier").val();

    if (data["supplier_label"] == "") {
        return
    }

    $.ajax({
        url: proxyPath + "products/suppliers",
        method: "POST",
        dataType: 'json',
        data: data,
    }).done(function(){
        gjsUtils.message("supplier added", "success")
        $("input#addsupplier").val("");
    }).fail(function (jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });
}

//
// add producer
//
function addProducer() {
    data = {};
    data["producer_label"] = $("input#addproducer").val();

    if (data["producer_label"] == "") {
        return
    }

    $.ajax({
        url: proxyPath + "products/producers",
        method: "POST",
        dataType: 'json',
        data: data,
    }).done(function(){
        gjsUtils.message("producer added", "success")
        $("input#addproducer").val("");
    }).fail(function (jqXHR, textStatus, errorThrown) {
        handleHTTPError(jqXHR.statusText, jqXHR.status)
    });
}