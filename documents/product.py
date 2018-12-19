# -*- coding: utf-8; -*-
#
# (c) 2011-2015 Thomas Bellembois thomas.bellembois@ens-lyon.fr
#
# This file is part of Chimithèque.
#
# Chimithèque is free software; you can redistribute it and/or modify
# it under the terms of the Cecill as published by the CEA, CNRS and INRIA
# either version 2 of the License, or (at your option) any later version.
#
# Chimithèque is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the Cecill along with Chimithèque.
# If not, see <http://www.cecill.info/licences/>.
#
# SVN Revision:
# -------------
# $Id: product.py 222 2015-07-21 16:06:35Z tbellemb2 $
#
from c_entity_mapper import ENTITY_MAPPER
from c_product_mapper import PRODUCT_MAPPER
from c_storage_mapper import STORAGE_MAPPER
from c_stock_mapper import STOCK_MAPPER
from c_store_location_mapper import STORE_LOCATION_MAPPER
from c_person_mapper import PERSON_MAPPER
from c_exposure_card_mapper import EXPOSURE_CARD_MAPPER
from c_exposure_item_mapper import EXPOSURE_ITEM_MAPPER
from chimitheque_logger import chimitheque_logger
from chimitheque_validators import IS_IN_DB_AND_USER_STORE_LOCATION, IS_IN_DB_AND_USER_ENTITY
from plugin_paginator import Paginator, PaginateSelector, PaginateInfo
from tempfile import NamedTemporaryFile, mkdtemp
from types import StringType
from collections import OrderedDict
from datetime import datetime
from time import strftime
from gluon import current
import base64
import Levenshtein
import chimitheque_commons as cc
import csv
import tarfile
import os
import json

mylogger = chimitheque_logger()

crud.messages.record_created = cc.get_string("PRODUCT_CREATED")
crud.messages.record_updated = cc.get_string("PRODUCT_UPDATED")
crud.messages.record_deleted = cc.get_string("PRODUCT_DELETED")


def ajax_ef_calculator():
    '''
    return the empirical formula from the linear formula
    '''
    mylogger.debug(message='request.vars:%s' %request.vars)
    f = request.vars['f']
    return cc.linear_to_empirical_formula(f)

def ajax_add_coc():
    '''
    add a new class of compounds
    '''
    mylogger.debug(message='request.vars:%s' %request.vars)
    label = request.vars['text']
    if len(label) == 0:
        return '1;%s' %cc.get_string("ENTER_A_LABEL")

    # check that the coc does NOT already exists
    count = db(db.class_of_compounds.label==label).count()
    if count != 0:
        return '1;%s' %cc.get_string("COC_ALREADY_EXIST")
    else:
        id = db.class_of_compounds.insert(label=label)
        return '0;' + str(id)


def ajax_sp():
    '''
    magical safety phrases selector
    '''
    mylogger.debug(message='request.vars:%s' %request.vars)
    text = request.vars['text']
    if len(text) == 0:
        return ''

    r = re.compile(r"(((\d){1,2}/(\d){1,2})|((\d){1,2}))" ,re.MULTILINE|re.DOTALL)
    m = r.findall(text)
    mylogger.debug(message='m:%s' %m)
    matches = []
    if m:
        for _m in m:
            matches.append(string.replace(_m[0], ' ', ''))
    mylogger.debug(message='matches:%s' %matches)
    ret = ''
    if (len(matches) != 0):
        for match in matches:
            hs = db(db.safety_phrase.reference==match).select(db.safety_phrase.id, db.safety_phrase.label).first()
            if hs: # should always return something
                l = str(current.T(hs.label))
                ret = '%s;%s|%s' %(ret, hs.id, base64.b64encode(l).encode('utf-8'))

    ret = ret[1:]
    mylogger.debug(message='ret:%s' %ret)
    return ret


def ajax_rp():
    '''
    magical risk phrases selector
    '''
    mylogger.debug(message='request.vars:%s' %request.vars)
    text = request.vars['text']
    if len(text) == 0:
        return ''

    r = re.compile(r"((\d){1,2}/(\d){1,2}/(\d){1,2}|(\d){1,2}/(\d){1,2}|(\d){1,2})" ,re.MULTILINE|re.DOTALL)
    m = r.findall(text)
    mylogger.debug(message='m:%s' %m)
    matches = []
    if m:
        for _m in m:
            matches.append(string.replace(_m[0], ' ', ''))
    mylogger.debug(message='matches:%s' %matches)
    ret = ''
    if (len(matches) != 0):
        for match in matches:
            hs = db(db.risk_phrase.reference==match).select(db.risk_phrase.id,db.risk_phrase.label).first()
            if hs: # should always return something
                l = str(current.T(hs.label))
                ret = '%s;%s|%s' %(ret, hs.id, base64.b64encode(l).encode('utf-8'))
    ret = ret[1:]
    mylogger.debug(message='ret:%s' %ret)
    return ret


def ajax_hs():
    '''
    magical hazard statements selector
    '''
    mylogger.debug(message='request.vars:%s' %request.vars)
    text = request.vars['text']
    if len(text) == 0:
        return ''

    r = re.compile(r"((EU)?H(\s)?(\d){3}(i|D|Df|F|Fd|FD|d|f|fd|A)?)" ,re.MULTILINE|re.DOTALL)
    m = r.findall(text)
    mylogger.debug(message='m:%s' %m)
    matches = []
    if m:
        for _m in m:
            matches.append(string.replace(_m[0], ' ', ''))
    mylogger.debug(message='matches:%s' %matches)
    ret = ''
    if (len(matches) != 0):
        for match in matches:
            hs = db(db.hazard_statement.reference==match).select(db.hazard_statement.id, db.hazard_statement.label).first()
            if hs: # should always return something
                l = str(current.T(hs.label))
                ret = '%s;%s|%s' %(ret, hs.id, base64.b64encode(l).encode('utf-8'))
    ret = ret[1:]
    mylogger.debug(message='ret:%s' %ret)
    return ret


def ajax_ps():
    '''
    magical precautionary statements selector
    '''
    mylogger.debug(message='request.vars:%s' %request.vars)
    text = request.vars['text']
    if len(text) == 0:
        return ''

    r = re.compile(r"((P(\s)?(\d){3}(\s)?(\+)?(\s)?){1,3})" ,re.MULTILINE|re.DOTALL)
    m = r.findall(text)
    mylogger.debug(message='m:%s' %m)
    matches = []
    if m:
        for _m in m:
            matches.append(string.replace(_m[0], ' ', ''))
    mylogger.debug(message='matches:%s' %matches)
    ret = ''
    if (len(matches) != 0):
        for match in matches:
            ps = db(db.precautionary_statement.reference==match).select(db.precautionary_statement.id, db.precautionary_statement.label).first()
            if ps: # should always return something
                l = str(current.T(ps.label))
                ret = '%s;%s|%s' %(ret, ps.id, base64.b64encode(l).encode('utf-8'))
    ret = ret[1:]
    mylogger.debug(message='ret:%s' %ret)
    return ret


def ajax_check_empirical_formula():
    """Return sorted empirical formula or an error."""
    mylogger.debug(message='request.vars:%s' % request.vars)
    empirical_formula = request.vars['empirical_formula_test']
    if len(empirical_formula) == 0:
        return ''

    sorted_formula, error = cc.sort_empirical_formula(empirical_formula)

    if error:
        return DIV(error)
    else:
        return DIV(sorted_formula)


def ajax_check_cas():
    """Return the products with the same CAS number."""
    mylogger.debug(message='request.vars:%s' % request.vars)
    cas_number = request.vars['cas_number']

    if len(cas_number) == 0:
        return ''

    _cas, error = cc.is_cas_number(cas_number)

    products = db((db.product.cas_number == cas_number) &
                  (db.product.archive == False)).select(db.product.name,
                                                        db.product.specificity,
                                                        db.product.id)
    mylogger.debug(message='products:%s' % products)

    _ul = UL()
    for product in products:
        _ul.append(LI(A('%(product_name)s %(product_specificity)s' % {'product_name': db.product.name.represent(product.name),
                                                                      'product_specificity': product.specificity},
                        _href=URL(request.application,
                                  'product',
                                  'search',
                                  vars={'product_id': product.id}))))

    if len(products) > 0:
        return DIV(cc.get_string("PRODUCT_WITH_SAME_CAS"), _ul)
    elif error:
        return DIV('%s - %s' % (cc.get_string("CAS_ERROR"), error), _ul)
    else:
        return DIV(cc.get_string("CAS_AVAILABLE"), _ul)


@auth.requires_login()
def bookmark():
    """Bookmark a product."""
    product_id = request.args[0]

    mylogger.debug(message='product_id:%s' % product_id)

    product_mapper = PRODUCT_MAPPER()
    _product = product_mapper.find(product_id=product_id)[0]

    _product.bookmark(user_id=auth.user.id)

    # in case of error, return json.dumps({'error': 'Error message'})
    return json.dumps({'success': True})


@auth.requires_login()
def unbookmark():
    """Unbookmark a product."""
    product_id = request.args[0]

    mylogger.debug(message='product_id:%s' % product_id)

    product_mapper = PRODUCT_MAPPER()
    _product = product_mapper.find(product_id=product_id)[0]

    _product.unbookmark(user_id=auth.user.id)

    # in case of error, return json.dumps({'error': 'Error message'})
    return json.dumps({'success': True})


@auth.requires_login()
def quick_add_exposure_card():
    """Add a product to the auth user active exposure card."""
    product_id = request.args[0]

    mylogger.debug(message='product_id:%s' % product_id)

    product_mapper = PRODUCT_MAPPER()
    person_mapper = PERSON_MAPPER()
    exposure_card_mapper = EXPOSURE_CARD_MAPPER()
    exposure_item_mapper = EXPOSURE_ITEM_MAPPER()
    _product = product_mapper.find(product_id=product_id)[0]
    _person = person_mapper.find(person_id=auth.user.id)[0]

    _active_ec = _person.get_active_exposure_card()

    if _active_ec is None:
        _title = strftime("%c")
        _active_ec = _person.create_exposure_card(title=_title)

    _exposure_item = _active_ec.add_exposure_item_for_product(_product)
    exposure_item_mapper.update(_exposure_item)
    exposure_card_mapper.update(_active_ec)

    person_mapper.update_exposure_card(_person)

    # in case of error, return json.dumps({'error': 'Error message'})
    return json.dumps({'success': True})

@auth.requires_login()
def quick_delete_exposure_card():
    """Delete a product of the auth user active exposure card."""
    product_id = request.args[0]

    mylogger.debug(message='product_id:%s' % product_id)

    product_mapper = PRODUCT_MAPPER()
    person_mapper = PERSON_MAPPER()
    exposure_card_mapper = EXPOSURE_CARD_MAPPER()
    exposure_item_mapper = EXPOSURE_ITEM_MAPPER()
    _product = product_mapper.find(product_id=product_id)[0]
    _person = person_mapper.find(person_id=auth.user.id)[0]

    _active_ec = _person.get_active_exposure_card()
    mylogger.debug(message='_active_ec.exposure_items:%s' % _active_ec.exposure_items)
    #_active_ec.delete_exposure_item_for_product(_product)
    #mylogger.debug(message='_active_ec.exposure_items:%s' % _active_ec.exposure_items)
    _exposure_item = _active_ec.get_exposure_item_for_product(_product)
    exposure_item_mapper.delete(_exposure_item)

    #exposure_card_mapper.update(_active_ec)

    # in case of error, return json.dumps({'error': 'Error message'})
    return json.dumps({'success': True})

@auth.requires_login()
@auth.requires(auth.has_permission('create_pc') or
               auth.has_permission('admin'))
def create():
    """Create or clone a product card."""
    if request.args(0):
        #
        # cloning the product id passed in args(0)
        #
        row = db(db.product.id == request.args(0)).select().first()
        db.product.name.default = row.name
        db.product.synonym.default = row.synonym
        db.product.cas_number.default = row.cas_number
        db.product.physical_state.default = row.physical_state
        db.product.class_of_compounds.default = row.class_of_compounds
        db.product.signal_word.default = row.signal_word
        db.product.risk_phrase.default = row.risk_phrase
        db.product.safety_phrase.default = row.safety_phrase
        db.product.hazard_statement.default = row.hazard_statement
        db.product.precautionary_statement.default = row.precautionary_statement
        db.product.empirical_formula.default = row.empirical_formula
        db.product.disposal_comment.default = row.disposal_comment
        db.product.signal_word.default = row.signal_word
        db.product.remark.default = row.remark
        db.product.hazard_code.default = row.hazard_code
        db.product.symbol.default = row.symbol

    # if the connected user can NOT create restricted product cards, disabling
    # the checkbox
    if (not auth.has_permission('admin')) and (not auth.has_permission('create_rpc')):
        db.product.restricted_access.writable = False

    form = SQLFORM(db.product, submit_button=cc.get_string("SUBMIT"))
    if form.accepts(request.vars, session):
        redirect(URL(a=request.application,
                     c=request.controller,
                     f='search',
                     vars={'product_id': form.vars['id']}))
    elif form.errors:
        session.flash = DIV(cc.get_string("MISSING_FIELDS"), _class="flasherror")

    return dict(form=form)


def _delete():
    """Delete a product card. Only orphan product cards can be deleted.

    This function should not be called directly but by the delete and delete_restricted
    functions either.
    """
    # getting the orphan product card ids
    not_orphan_products = db(db.storage).select(db.storage.product, groupby=db.storage.product)
    not_orphan_product_ids = [str(r.product) for r in not_orphan_products]

    # checking that the product card to be deleted is orphan and deleting it
    if request.args(0) and (str(request.args(0)) not in not_orphan_product_ids):
        db(db.product.id == request.args(0)).delete()

        session.flash = cc.get_string("PRODUCT_DELETED")
    else:
        session.flash = cc.get_string("PRODUCT_CAN_NOT_BE_DELETED")

    return redirect(URL(request.application, 'product', 'search'))


@auth.requires_login()
@auth.requires(auth.has_permission('delete_rpc') or
               auth.has_permission('admin'))
def delete_restricted():
    """Delete a restricted product card."""
    return _delete()


@auth.requires_login()
@auth.requires(auth.has_permission('delete_pc') or
               auth.has_permission('admin'))
def delete():
    """Delete a not restricted product card."""
    return _delete()


@auth.requires_login()
@auth.requires(auth.has_permission('update_pc') or
               auth.has_permission('admin'))
def update():
    """Update a product card."""
    mylogger.debug(message='request.vars:%s' % request.vars)

    if (not auth.has_permission('admin')) and (not auth.has_permission('update_rpc')):
        db.product.restricted_access.writable = False

    form = crud.update(db.product,
                       request.args(0),
                       next=URL(a=request.application,
                                c=request.controller,
                                f='search',
                                vars={'product_id': request.args(0)}),
                       onaccept=lambda form: auth.archive(form, archive_table=db.product_history, archive_current=False),
                       ondelete=lambda form: auth.archive(form, archive_table=db.product_history, archive_current=False))

    if form.errors:
        session.flash = DIV(cc.get_string("MISSING_FIELDS"), _class="flasherror")

    cache.ram.clear(regex='.*/product/card')

    return dict(form=form)


@auth.requires(auth.has_permission('read_pc') or
               auth.has_permission('admin'))
@auth.requires_login()
def list_history():
    """List the product changes."""
    mylogger.debug(message='list_history')
    product_id = request.args(0)
    rows_product_history = db((db.product_history.current_record == product_id) &
                              (db.person.id == db.product_history.person)).select(db.person.email,
                                                                                  db.product_history.id,
                                                                                  db.product_history.modification_datetime,
                                                                                  orderby=db.product_history.modification_datetime)

    return dict(product_id=product_id,
                rows_product_history=rows_product_history)


@cache(request.env.path_info, time_expire=3600, cache_model=cache.ram)
@auth.requires_login()
@auth.requires(auth.has_permission('read_pc') or
               auth.has_permission('admin'))
def card():
    """Return the product card"""
    is_history = 'is_history' in request.vars
    _product = PRODUCT_MAPPER().find(product_id=request.args[0], history=is_history)[0]

    d = dict(product=_product, is_history=is_history)

    return response.render(d)


def details_reload():

    return details()


# No cache because the output depend the user
#@cache(request.env.path_info, time_expire=3600, cache_model=cache.ram)
@auth.requires_login()
@auth.requires(auth.has_permission('read_pc') or
               auth.has_permission('admin'))
def details():

    mylogger.debug(message='details')

    _product = PRODUCT_MAPPER().find(product_id=request.args[0])[0]

    _user_entities = ENTITY_MAPPER().find(person_id=auth.user.id)

    _stocks = STOCK_MAPPER().find(product_id=_product.id, entity_id=[_entity.id for _entity in _user_entities])

    d = dict(product=_product, stocks=_stocks)

    return response.render(d)


@auth.requires_login()
@auth.requires(auth.has_permission('admin'))
def import_from_csv():
    """Import products from a Chimithèque export."""
    temp_dir = mkdtemp()
    form = SQLFORM.factory(Field('upload',
                                 'upload',
                                 uploadfolder=temp_dir,
                                 label=cc.get_string("UPLOAD"),
                                 requires=IS_NOT_EMPTY()))

    if form.accepts(request.vars):

        mylogger.ram(cc.get_string("STARTING_TASK"), index=0)
        filename = form.vars.upload
        mylogger.debug(message='filename:%s' % filename)

        mylogger.ram('opening archive file', index=-1)
        tar_file = tarfile.open(os.path.join(temp_dir, filename), 'r:gz')
        mylogger.ram('extracting archive file', index=-1)
        mylogger.debug(message='temp_dir:%s' % temp_dir)

        try:
            mylogger.debug(message='extraction start')
            tar_file.extractall(path=temp_dir)
            mylogger.debug(message='extraction end')
            mylogger.ram('extraction done', index=-1)
        except Exception as ex:
            mylogger.debug(message='exception type:%s' % type(ex))
            mylogger.debug(message='exception args:%s' % ex.args)
            return dict(form=form)

        tar_file.close()

        mylogger.ram('loading local product database', index=-1)
        mylogger.debug(message='select product')
        local_database_products = db(db.product).select()
        local_database_product_cas_numbers = [p.cas_number for p in local_database_products]

        source_name = {}
        mylogger.debug(message='opening name.csv')
        mylogger.ram('opening name.csv', index=-1)
        csv_name_reader = csv.reader(open(os.path.join(temp_dir, 'name.csv'), 'rb'), delimiter=',')
        for row in csv_name_reader:
            source_name[row[0]] = row[1]  # id -> label

        destination_name = {}
        for row in db(db.name).select():
            destination_name[row.label] = row.id  # label -> id

        source_empirical_formula = {}
        mylogger.debug(message='opening empirical_formula.csv')
        mylogger.ram('opening empirical_formula.csv', index=-1)
        csv_empirical_formula_reader = csv.reader(open(os.path.join(temp_dir, 'empirical_formula.csv'), 'rb'), delimiter=',')
        for row in csv_empirical_formula_reader:
            source_empirical_formula[row[0]] = row[1]  # id -> label

        destination_empirical_formula = {}
        for row in db(db.empirical_formula).select():
            destination_empirical_formula[row.label] = row.id  # label -> id

        source_linear_formula = {}
        mylogger.debug(message='opening linear_formula.csv')
        mylogger.ram('opening linear_formula.csv', index=-1)
        csv_linear_formula_reader = csv.reader(open(os.path.join(temp_dir, 'linear_formula.csv'), 'rb'), delimiter=',')
        for row in csv_linear_formula_reader:
            source_linear_formula[row[0]] = row[1]  # id -> label

        destination_linear_formula = {}
        for row in db(db.linear_formula).select():
            destination_linear_formula[row.label] = row.id  # label -> id

        mylogger.debug(message='opening product.csv')
        mylogger.ram('opening product.csv', index=-1)
        mylogger.ram(' ', index=-1)
        csv_product_reader = csv.DictReader(open(os.path.join(temp_dir, 'product.csv'), 'rb'), delimiter=',')
        _count_imported = 0
        _count_not_imported = 0

        # compatibility trick
        if 'product.cas_number' in csv_product_reader.fieldnames:
            _table_product = 'product'
        else:
            _table_product = 'PRODUCT'

        #
        # data remote<->local id mapping
        #

        # physical_state mapping
        mylogger.debug(message='local database physical_state mapping')
        source_physical_state = {}
        csv_physical_state_reader = csv.reader(open(os.path.join(temp_dir, 'physical_state.csv'), 'rb'), delimiter=',')

        for row in csv_physical_state_reader:
            source_physical_state[row[1]] = row[0]  # label -> id

        mapping_physical_state = {}
        try:
            for _row in db(db.physical_state).select():
                mylogger.debug(message='physical_state local_id:%s remote_id:%s' % (_row.id, source_physical_state[_row.label]))
                mapping_physical_state[source_physical_state[_row.label]] = _row.id
        except KeyError:
            mylogger.debug(message='key:%s does not exist in the source database' % _row.label)

        # hazard_statement
        mylogger.debug(message='local database hazard_statement mapping')
        source_hazard_statement = {}
        csv_hazard_statement_reader = csv.reader(open(os.path.join(temp_dir, 'hazard_statement.csv'), 'rb'), delimiter=',')

        for row in csv_hazard_statement_reader:
            source_hazard_statement[row[1]] = row[0]  # label -> id

        mapping_hazard_statement = {}
        try:
            for _row in db(db.hazard_statement).select():
                mylogger.debug(message='hazard_statement local_id:%s remote_id:%s' % (_row.id, source_hazard_statement[_row.reference]))
                mapping_hazard_statement[source_hazard_statement[_row.reference]] = _row.id
        except KeyError:
            mylogger.debug(message='key:%s does not exist in the source database' % _row.reference)

        # precautionary_statement
        mylogger.debug(message='local database precautionary_statement mapping')
        source_precautionary_statement = {}
        csv_precautionary_statement_reader = csv.reader(open(os.path.join(temp_dir, 'precautionary_statement.csv'), 'rb'), delimiter=',')

        for row in csv_precautionary_statement_reader:
            source_precautionary_statement[row[1]] = row[0]  # label -> id

        mapping_precautionary_statement = {}
        try:
            for _row in db(db.precautionary_statement).select():
                mylogger.debug(message='precautionary_statement local_id:%s remote_id:%s' % (_row.id, source_precautionary_statement[_row.reference]))
                mapping_precautionary_statement[source_precautionary_statement[_row.reference]] = _row.id
        except KeyError:
            mylogger.debug(message='key:%s does not exist in the source database' % _row.reference)

        # risk_phrase
        mylogger.debug(message='local database risk_phrase mapping')
        source_risk_phrase = {}
        csv_risk_phrase_reader = csv.reader(open(os.path.join(temp_dir, 'risk_phrase.csv'), 'rb'), delimiter=',')

        for row in csv_risk_phrase_reader:
            source_risk_phrase[row[1]] = row[0]  # label -> id

        mapping_risk_phrase = {}
        try:
            for _row in db(db.risk_phrase).select():
                mylogger.debug(message='risk_phrase local_id:%s remote_id:%s' % (_row.id, source_risk_phrase[_row.reference]))
                mapping_risk_phrase[source_risk_phrase[_row.reference]] = _row.id
        except KeyError:
            mylogger.debug(message='key:%s does not exist in the source database' % _row.reference)

        # safety_phrase
        mylogger.debug(message='local database safety_phrase mapping')
        source_safety_phrase = {}
        csv_safety_phrase_reader = csv.reader(open(os.path.join(temp_dir, 'safety_phrase.csv'), 'rb'), delimiter=',')

        for row in csv_safety_phrase_reader:
            source_safety_phrase[row[1]] = row[0]  # label -> id

        mapping_safety_phrase = {}
        try:
            for _row in db(db.safety_phrase).select():
                mylogger.debug(message='safety_phrase local_id:%s remote_id:%s' % (_row.id, source_safety_phrase[_row.reference]))
                mapping_safety_phrase[source_safety_phrase[_row.reference]] = _row.id
        except KeyError:
            mylogger.debug(message='key:%s does not exist in the source database' % _row.reference)

        # symbol mapping
        mylogger.debug(message='local database symbol mapping')
        source_symbol = {}
        csv_symbol_reader = csv.reader(open(os.path.join(temp_dir, 'symbol.csv'), 'rb'), delimiter=',')

        for row in csv_symbol_reader:
            source_symbol[row[1]] = row[0]  # label -> id

        mapping_symbol = {}
        try:
            for _row in db(db.symbol).select():
                mylogger.debug(message='symbol local_id:%s remote_id:%s' % (_row.id, source_symbol[_row.label]))
                mapping_symbol[source_symbol[_row.label]] = _row.id
        except KeyError:
            mylogger.debug(message='key:%s does not exist in the source database' % _row.label)

        # hazard_code mapping
        mylogger.debug(message='local database hazard_code mapping')
        source_hazard_code = {}
        csv_hazard_code_reader = csv.reader(open(os.path.join(temp_dir, 'hazard_code.csv'), 'rb'), delimiter=',')

        for row in csv_hazard_code_reader:
            source_hazard_code[row[1]] = row[0]  # label -> id

        mapping_hazard_code = {}
        try:
            for _row in db(db.hazard_code).select():
                mylogger.debug(message='hazard_code local_id:%s remote_id:%s' % (_row.id, source_hazard_code[_row.label]))
                mapping_hazard_code[source_hazard_code[_row.label]] = _row.id
        except KeyError:
            mylogger.debug(message='key:%s does not exist in the source database' % _row.label)

        # signal_word mapping
        mylogger.debug(message='local database signal_word mapping')
        source_signal_word = {}
        csv_signal_word_reader = csv.reader(open(os.path.join(temp_dir, 'signal_word.csv'), 'rb'), delimiter=',')

        for row in csv_signal_word_reader:
            source_signal_word[row[1]] = row[0]  # label -> id

        mapping_signal_word = {}
        try:
            for _row in db(db.signal_word).select():
                mylogger.debug(message='signal_word local_id:%s remote_id:%s' % (_row.id, source_signal_word[_row.label]))
                mapping_signal_word[source_signal_word[_row.label]] = _row.id
        except KeyError:
            mylogger.debug(message='key:%s does not exist in the source database' % _row.label)

        #
        # product insert
        #
        for row in csv_product_reader:

            mylogger.ram('imported: %s - not imported: %s' % (_count_imported, _count_not_imported), index=-2)
            mylogger.debug(message='imported: %s - not imported: %s' % (_count_imported, _count_not_imported))

            mylogger.debug(message='row:%s' % row)
            _cas_number = row['%s.cas_number' % _table_product]
            _ce_number = row['%s.ce_number' % _table_product] if row['%s.ce_number' % _table_product] != '<NULL>' else None
            _name_id = row['%s.name' % _table_product]
            _synonym_id = row['%s.synonym' % _table_product]
            _specificity = row['%s.specificity' % _table_product]
            _empirical_formula_id = row['%s.empirical_formula' % _table_product] if row['%s.empirical_formula' % _table_product] != '<NULL>' else None
            _linear_formula_id = row['%s.linear_formula' % _table_product] if row['%s.linear_formula' % _table_product] != '<NULL>' else None
            _msds = row['%s.msds' % _table_product]
            _class_of_compound_ids = row['%s.class_of_compounds' % _table_product]

            if row['%s.physical_state' % _table_product] in mapping_physical_state.keys():
                _physical_state_id = mapping_physical_state[row['%s.physical_state' % _table_product]] if row['%s.physical_state' % _table_product] != '<NULL>' else None
            else:
                _physical_state_id = None

            if row['%s.signal_word' % _table_product] in mapping_signal_word.keys():
                _signal_word_id = mapping_signal_word[row['%s.signal_word' % _table_product]] if row['%s.signal_word' % _table_product] != '<NULL>' else None
            else:
                _signal_word_id = None

            _hazard_code_id = '|'
            for i in row['%s.hazard_code' % _table_product].split('|'):
                if i in mapping_hazard_code.keys():
                    _hazard_code_id = _hazard_code_id + str(mapping_hazard_code[i]) + '|'

            _symbol_id = '|'
            for i in row['%s.symbol' % _table_product].split('|'):
                if i in mapping_symbol.keys():
                    _symbol_id = _symbol_id + str(mapping_symbol[i]) + '|'

            _risk_phrase_id = '|'
            for i in row['%s.risk_phrase' % _table_product].split('|'):
                if i in mapping_risk_phrase.keys():
                    _risk_phrase_id = _risk_phrase_id + str(mapping_risk_phrase[i]) + '|'

            _safety_phrase_id = '|'
            for i in row['%s.safety_phrase' % _table_product].split('|'):
                if i in mapping_safety_phrase.keys():
                    _safety_phrase_id = _safety_phrase_id + str(mapping_safety_phrase[i]) + '|'

            _hazard_statement_id = '|'
            for i in row['%s.hazard_statement' % _table_product].split('|'):
                if i in mapping_hazard_statement.keys():
                    _hazard_statement_id = _hazard_statement_id + str(mapping_hazard_statement[i]) + '|'

            _precautionary_statement_id = '|'
            for i in row['%s.precautionary_statement' % _table_product].split('|'):
                if i in mapping_precautionary_statement.keys():
                    _precautionary_statement_id = _precautionary_statement_id + str(mapping_precautionary_statement[i]) + '|'

            _disposal_comment = row['%s.disposal_comment' % _table_product]
            _remark = row['%s.remark' % _table_product]
            if _cas_number not in local_database_product_cas_numbers:
                mylogger.debug(message='%s to be imported' % _cas_number)

                _name_label = source_name[_name_id]
                mylogger.debug(message='_name_label:%s' % _name_label)
                if _name_label in destination_name.keys():
                    _new_name_id = destination_name[_name_label]
                else:
                    # inserting the new name
                    _new_name_id = db.name.insert(label=_name_label)
                    destination_name[_name_label] = _new_name_id

                if _empirical_formula_id:
                    mylogger.debug(message='_empirical_formula_id:%s' % _empirical_formula_id)
                    if _empirical_formula_id == '0':
                        _new_empirical_formula_id = '1'
                    else:
                        _empirical_formula_label = source_empirical_formula[_empirical_formula_id]
                        if _empirical_formula_label in destination_empirical_formula.keys():
                            _new_empirical_formula_id = destination_empirical_formula[_empirical_formula_label]
                        else:
                            # inserting the new empirical formula
                            _new_empirical_formula_id = db.empirical_formula.insert(label=_empirical_formula_label)
                            destination_empirical_formula[_empirical_formula_label] = _new_empirical_formula_id
                else:
                    _new_empirical_formula_id = None
                mylogger.debug(message='_new_empirical_formula_id:%s' % _new_empirical_formula_id)

                if _linear_formula_id:
                    mylogger.debug(message='_linear_formula_id:%s' % _linear_formula_id)
                    if _linear_formula_id == '0':
                        _new_linear_formula_id = '1'
                    else:
                        _linear_formula_label = source_linear_formula[_linear_formula_id]
                        if _linear_formula_label in destination_linear_formula.keys():
                            _new_linear_formula_id = destination_linear_formula[_linear_formula_label]
                        else:
                            # inserting the new linear formula
                            _new_linear_formula_id = db.linear_formula.insert(label=_linear_formula_label)
                            destination_linear_formula[_linear_formula_label] = _new_linear_formula_id
                else:
                    _new_linear_formula_id = None
                mylogger.debug(message='_new_linear_formula_id:%s' % _new_linear_formula_id)

                mylogger.debug(message='_synonym_id:%s' % _synonym_id)
                if _synonym_id == '|0|':
                    _new_synonym_id = [0]
                else:
                    _new_synonym_id = []
                    for _id in _synonym_id.split('|'):
                        mylogger.debug(message='_id:%s' % _id)
                        if _id != '':
                            # _id is not in source_name.keys() in case of accidental name entry deletion
                            # should never happen
                            _synonym_label = source_name[_id] if _id in source_name.keys() else None
                            if _synonym_label is None:
                                break
                            mylogger.debug(message='_synonym_label:%s' % _synonym_label)
                            if _synonym_label in destination_name.keys():
                                mylogger.debug(message='_synonym_label exist:%s' % _synonym_label)
                                _new_synonym_id.append(int(destination_name[_synonym_label]))
                            else:
                                mylogger.debug(message='_synonym_label not exist:%s' % _synonym_label)
                                # inserting the new synonym
                                _new_id = db.name.insert(label=_synonym_label)
                                _new_synonym_id.append(_new_id)
                                destination_name[_synonym_label] = _new_id

                mylogger.debug(message='_new_synonym_id' % _new_synonym_id)

                mylogger.debug(message='''
                cas_number:%(cas_number)s
                risk_phrase:%(risk_phrase)s
                hazard_statement:%(hazard_statement)s
                precautionary_statement:%(precautionary_statement)s
                symbol:%(symbol)s
                signal_word:%(signal_word)s
                physical_state:%(physical_state)s
                hazard_code:%(hazard_code)s''' % {'cas_number': _cas_number,
                                                  'physical_state': _physical_state_id,
                                                  'hazard_code': _hazard_code_id,
                                                  'signal_word': _signal_word_id,
                                                  'risk_phrase': _risk_phrase_id,
                                                  'hazard_statement': _hazard_statement_id,
                                                  'precautionary_statement': _precautionary_statement_id,
                                                  'symbol': _symbol_id})

                db.product.insert(cas_number=_cas_number,
                                  ce_number=_ce_number,
                                  name=_new_name_id,
                                  synonym=_new_synonym_id,
                                  person=auth.user.id,
                                  specificity=_specificity,
                                  empirical_formula=_new_empirical_formula_id,
                                  linear_formula=_new_linear_formula_id,
                                  msds=_msds,
                                  physical_state=_physical_state_id,
                                  # class_of_compounds=[ i for i in _class_of_compound_ids.split('|') if i != '' ],
                                  hazard_code=[i for i in _hazard_code_id.split('|')] if _hazard_code_id is not None else None,
                                  symbol=[i for i in _symbol_id.split('|') if i != ''],
                                  signal_word=_signal_word_id,
                                  risk_phrase=[i for i in _risk_phrase_id.split('|') if i != ''],
                                  safety_phrase=[i for i in _safety_phrase_id.split('|') if i != ''],
                                  hazard_statement=[i for i in _hazard_statement_id.split('|') if i != ''],
                                  precautionary_statement=[i for i in _precautionary_statement_id.split('|') if i != ''],
                                  disposal_comment=_disposal_comment,
                                  remark=_remark)

                local_database_product_cas_numbers.append(_cas_number)
                _count_imported = _count_imported + 1
            else:
                _count_not_imported = _count_not_imported + 1
                mylogger.debug(message='%s already exist' % _cas_number)

        db.commit()

        mylogger.debug(message='_count_imported:%i' % _count_imported)
        mylogger.ram('imported %i products' % _count_imported, index=-1)

    return dict(form=form)


@auth.requires_login()
@auth.requires(auth.has_permission('admin'))
def export_to_csv():
    """Export products into CSV."""
    mylogger.ram(cc.get_string("STARTING_TASK"), index=0)
    tmp_file_return_fd = NamedTemporaryFile()

    mylogger.ram('table product', index=1)
    tmp_file_product_fd = NamedTemporaryFile()
    tmp_file_product_fd.write(str(db((db.product)).select(db.product.cas_number,
                                                          db.product.ce_number,
                                                          db.product.name,
                                                          db.product.synonym,
                                                          db.product.specificity,
                                                          db.product.empirical_formula,
                                                          db.product.linear_formula,
                                                          db.product.msds,
                                                          db.product.physical_state,
                                                          db.product.class_of_compounds,
                                                          db.product.hazard_code,
                                                          db.product.symbol,
                                                          db.product.signal_word,
                                                          db.product.risk_phrase,
                                                          db.product.safety_phrase,
                                                          db.product.hazard_statement,
                                                          db.product.precautionary_statement,
                                                          db.product.disposal_comment,
                                                          db.product.remark)))

    tmp_file_product_fd.flush()
    os.fsync(tmp_file_product_fd)

    mylogger.ram('table name', index=2)
    tmp_file_name_fd = NamedTemporaryFile()
    tmp_file_name_fd.write(str(db((db.name)).select(db.name.id, db.name.label)))
    tmp_file_name_fd.flush()
    os.fsync(tmp_file_name_fd)

    mylogger.ram('table empirical_formula', index=3)
    tmp_file_empirical_formula_fd = NamedTemporaryFile()
    tmp_file_empirical_formula_fd.write(str(db((db.empirical_formula)).select(db.empirical_formula.id, db.empirical_formula.label)))
    tmp_file_empirical_formula_fd.flush()
    os.fsync(tmp_file_empirical_formula_fd)

    mylogger.ram('table linear_formula', index=4)
    tmp_file_linear_formula_fd = NamedTemporaryFile()
    tmp_file_linear_formula_fd.write(str(db((db.linear_formula)).select(db.linear_formula.id, db.linear_formula.label)))
    tmp_file_linear_formula_fd.flush()
    os.fsync(tmp_file_linear_formula_fd)

    mylogger.ram('table physical_state', index=5)
    tmp_file_physical_state_fd = NamedTemporaryFile()
    tmp_file_physical_state_fd.write(str(db((db.physical_state)).select(db.physical_state.id, db.physical_state.label)))
    tmp_file_physical_state_fd.flush()
    os.fsync(tmp_file_physical_state_fd)

    mylogger.ram('table class_of_compounds', index=6)
    tmp_file_class_of_compounds_fd = NamedTemporaryFile()
    tmp_file_class_of_compounds_fd.write(str(db((db.class_of_compounds)).select(db.class_of_compounds.id, db.class_of_compounds.label)))
    tmp_file_class_of_compounds_fd.flush()
    os.fsync(tmp_file_class_of_compounds_fd)

    mylogger.ram('table hazard_code', index=7)
    tmp_file_hazard_code_fd = NamedTemporaryFile()
    tmp_file_hazard_code_fd.write(str(db((db.hazard_code)).select(db.hazard_code.id, db.hazard_code.label)))
    tmp_file_hazard_code_fd.flush()
    os.fsync(tmp_file_hazard_code_fd)

    mylogger.ram('table symbol', index=8)
    tmp_file_symbol_fd = NamedTemporaryFile()
    tmp_file_symbol_fd.write(str(db((db.symbol)).select(db.symbol.id, db.symbol.label)))
    tmp_file_symbol_fd.flush()
    os.fsync(tmp_file_hazard_code_fd)

    mylogger.ram('table signal_word', index=9)
    tmp_file_signal_word_fd = NamedTemporaryFile()
    tmp_file_signal_word_fd.write(str(db((db.signal_word)).select(db.signal_word.id, db.signal_word.label)))
    tmp_file_signal_word_fd.flush()
    os.fsync(tmp_file_signal_word_fd)

    mylogger.ram('table risk_phrase', index=10)
    tmp_file_risk_phrase_fd = NamedTemporaryFile()
    tmp_file_risk_phrase_fd.write(str(db((db.risk_phrase)).select(db.risk_phrase.id, db.risk_phrase.reference, db.risk_phrase.label)))
    tmp_file_risk_phrase_fd.flush()
    os.fsync(tmp_file_risk_phrase_fd)

    mylogger.ram('table safety_phrase', index=11)
    tmp_file_safety_phrase_fd = NamedTemporaryFile()
    tmp_file_safety_phrase_fd.write(str(db((db.safety_phrase)).select(db.safety_phrase.id, db.safety_phrase.reference, db.safety_phrase.label)))
    tmp_file_safety_phrase_fd.flush()
    os.fsync(tmp_file_safety_phrase_fd)

    mylogger.ram('table hazard_statement', index=12)
    tmp_file_hazard_statement_fd = NamedTemporaryFile()
    tmp_file_hazard_statement_fd.write(str(db((db.hazard_statement)).select(db.hazard_statement.id, db.hazard_statement.reference, db.hazard_statement.label)))
    tmp_file_hazard_statement_fd.flush()
    os.fsync(tmp_file_hazard_statement_fd)

    mylogger.ram('table precautionary_statement', index=13)
    tmp_file_precautionary_statement_fd = NamedTemporaryFile()
    tmp_file_precautionary_statement_fd.write(str(db((db.precautionary_statement)).select(db.precautionary_statement.id, db.precautionary_statement.reference, db.precautionary_statement.label)))
    tmp_file_precautionary_statement_fd.flush()
    os.fsync(tmp_file_precautionary_statement_fd)

    mylogger.ram('building tar file', index=14)
    return_tar = tarfile.open(tmp_file_return_fd.name, "w:gz")
    return_tar.add(tmp_file_product_fd.name, arcname='product.csv')
    return_tar.add(tmp_file_name_fd.name, arcname='name.csv')
    return_tar.add(tmp_file_empirical_formula_fd.name, arcname='empirical_formula.csv')
    return_tar.add(tmp_file_linear_formula_fd.name, arcname='linear_formula.csv')
    return_tar.add(tmp_file_physical_state_fd.name, arcname='physical_state.csv')
    return_tar.add(tmp_file_class_of_compounds_fd.name, arcname='class_of_compounds.csv')
    return_tar.add(tmp_file_hazard_code_fd.name, arcname='hazard_code.csv')
    return_tar.add(tmp_file_symbol_fd.name, arcname='symbol.csv')
    return_tar.add(tmp_file_signal_word_fd.name, arcname='signal_word.csv')
    return_tar.add(tmp_file_risk_phrase_fd.name, arcname='risk_phrase.csv')
    return_tar.add(tmp_file_safety_phrase_fd.name, arcname='safety_phrase.csv')
    return_tar.add(tmp_file_hazard_statement_fd.name, arcname='hazard_statement.csv')
    return_tar.add(tmp_file_precautionary_statement_fd.name, arcname='precautionary_statement.csv')

    tmp_file_return_fd.flush()
    os.fsync(tmp_file_return_fd)

    return_tar.close()

    tmp_file_product_fd.close()
    tmp_file_name_fd.close()
    tmp_file_empirical_formula_fd.close()
    tmp_file_physical_state_fd.close()
    tmp_file_class_of_compounds_fd.close()
    tmp_file_hazard_code_fd.close()
    tmp_file_symbol_fd.close()
    tmp_file_signal_word_fd.close()
    tmp_file_risk_phrase_fd.close()
    tmp_file_safety_phrase_fd.close()
    tmp_file_hazard_statement_fd.close()
    tmp_file_precautionary_statement_fd.close()

    response.headers = {'Content-disposition:': 'attachment; filename=chimitheque_db.tar.gz',
                        'Content-type:': 'application/gzip'}

    return response.stream(open(tmp_file_return_fd.name))


@auth.requires_login()
def search():
    mylogger.debug(message='request.vars:%s' % str(request.vars))
    mylogger.debug(message='request.args:%s' % str(request.args))

    if len(request.vars) == 0 and len(request.args) == 0:
        request.vars['borrower'] = auth.user.id

    # request arguments cleanup for autocomplete widgets
    if 'risk_phrase' in request.vars and request.vars['risk_phrase'] == '':
        del request.vars['risk_phrase']
    if 'safety_phrase' in request.vars and request.vars['safety_phrase'] == '':
        del request.vars['safety_phrase']
    if 'hazard_statement' in request.vars and request.vars['hazard_statement'] == '':
        del request.vars['hazard_statement']
    if 'precautionary_statement' in request.vars and request.vars['precautionary_statement'] == '':
        del request.vars['precautionary_statement']

    # some init
    query_list = []
    did_you_mean = []  # suggestions list
    rows = None
    join_bookmark = False  # join the bookmark table ?
    join_storage = False  # join the storage table ?
    join_borrow = False  # join the borrow table ?
    nb_entries = 1  # number of results
    label = ''  # request title, ie. "products in the Chemical Lab.
    user_entity = [_entity.id for _entity in ENTITY_MAPPER().find(person_id=auth.user.id)]
    page = int(request.vars['page']) if 'page' in request.vars else 0
    result_per_page = int(request.vars['result_per_page']) if 'result_per_page' in request.vars else 10
    export_csv = 'export_csv' in request.vars
    export_html = 'export_html' in request.vars
    is_did_you_mean = 'is_did_you_mean' in request.vars

    # no way to pass the "keep_last_search" variable while clicking on a "x results per page" link
    if 'paginate' in request.vars:
        request.vars['keep_last_search'] = True

    if 'request' in request.vars and request.vars['request'] == 'entity':
        session.search_request = request.vars['request']

        if (not auth.has_permission('select_sc')) and (not auth.has_permission('admin')):
            raise HTTP(403, "Not authorized")

        request.vars['entity'] = [request.vars['is_in_entity']]

    elif 'request' in request.vars and request.vars['request'] == 'store_location':
        session.search_request = request.vars['request']

        if (not auth.has_permission('read_sc')) and (not auth.has_permission('admin')):
            raise HTTP(403, "Not authorized")

        request.vars['store_location'] = [request.vars['is_in_store_location']]

    elif 'request' in request.vars and request.vars['request'] == 'all':
        session.search_request = request.vars['request']
        query_list.append(db.product.id > 0)
        label = cc.get_string("ALL_PRODUCT")

    elif 'request' in request.vars and request.vars['request'] == 'organization':
        session.search_request = request.vars['request']

        if (not auth.has_permission('select_sc')) and (not auth.has_permission('admin')):
            raise HTTP(403, "Not authorized")

        join_storage = True
        query_list.append(db.storage.id > 0)
        query_list.append(db.storage.archive == False)
        label = '%s %s' % (cc.get_string("PRODUCT_STORED_AT"), settings['organization'])

    #
    # restoring session vars if keep_last_search
    #
    if 'keep_last_search' in request.vars:
        if session.search_bookmark:
            request.vars['bookmark'] = session.search_bookmark
        if session.search_display_by:
            request.vars['display_by'] = session.search_display_by
        if session.search_order_by:
            request.vars['order_by'] = session.search_order_by
        if session.search_product_id:
            request.vars['product_id'] = session.search_product_id
        if session.search_result_per_page:
            request.vars['result_per_page'] = result_per_page = session.search_result_per_page
        if session.search_page:
            request.vars['page'] = page = session.search_page
        if session.search_request:
            request.vars['request'] = session.search_request
        if session.search_borrow_entity:
            request.vars['borrow_entity'] = session.search_borrow_entity
        if session.search_entity:
            request.vars['entity'] = session.search_entity
        if session.search_store_location:
            request.vars['store_location'] = session.search_store_location
        if session.search_include_children_store_location:
            request.vars['include_children_store_location'] = session.search_include_children_store_location
        if session.search_archive:
            request.vars['archive'] = session.search_archive
        if session.search_is_cmr:
            request.vars['is_cmr'] = session.search_is_cmr
        if session.search_is_radio:
            request.vars['is_radio'] = session.search_is_radio
        if session.search_to_destroy:
            request.vars['to_destroy'] = session.search_to_destroy
        if session.search_cas_number:
            request.vars['cas_number'] = session.search_cas_number
        if session.search_ce_number:
            request.vars['ce_number'] = session.search_ce_number
        if session.search_entry_datetime:
            request.vars['product_datetime'] = session.search_entry_datetime
        if session.search_entry_datetime:
            request.vars['entry_datetime'] = session.search_entry_datetime
        if session.search_exit_datetime:
            request.vars['exit_datetime'] = session.search_exit_datetime
        if session.search_barecode:
            request.vars['barecode'] = session.search_barecode
        if session.search_comment:
            request.vars['comment'] = session.search_comment
        if session.search_name:
            request.vars['name'] = session.search_name
        if session.search_empirical_formula:
            request.vars['empirical_formula'] = session.search_empirical_formula
        if session.search_linear_formula:
            request.vars['linear_formula'] = session.search_linear_formula
        if session.search_physical_state:
            request.vars['physical_state'] = session.search_physical_state
        if session.search_exact_coc:
            request.vars['exact_coc'] = session.search_exact_coc
        if session.search_class_of_compounds:
            request.vars['class_of_compounds'] = session.search_class_of_compounds
        if session.search_risk_phrase:
            request.vars['risk_phrase'] = session.search_risk_phrase
        if session.search_safety_phrase:
            request.vars['safety_phrase'] = session.search_safety_phrase
        if session.search_hazard_statement:
            request.vars['hazard_statement'] = session.search_hazard_statement
        if session.search_precautionary_statement:
            request.vars['precautionary_statement'] = session.search_precautionary_statement
        if session.search_hazard_code:
            request.vars['hazard_code'] = session.search_hazard_code
        if session.search_symbol:
            request.vars['symbol'] = session.search_symbol
        if session.search_person_pc:
            request.vars['person_pc'] = session.search_person_pc
        if session.search_person_sc:
            request.vars['person_sc'] = session.search_person_sc
        if session.search_person_asc:
            request.vars['person_asc'] = session.search_person_asc
        if session.search_borrower:
            request.vars['borrower'] = session.search_borrower
        del(request.vars['keep_last_search'])

    #
    # and then cleaning up session vars
    #
    for key in ['search_bookmark',
                'search_display_by',
                'search_order_by',
                'search_product_id',
                'search_result_per_page',
                'search_page',
                'search_request',
                'search_borrow_entity',
                'search_entity',
                'search_store_location',
                'search_include_children_store_location',
                'search_product_id',
                'search_archive',
                'search_is_cmr',
                'search_is_radio',
                'search_to_destroy',
                'search_cas_number',
                'search_ce_number',
                'search_product_datetime',
                'search_entry_datetime',
                'search_exit_datetime',
                'search_barecode',
                'search_comment',
                'search_name',
                'search_empirical_formula',
                'search_linear_formula',
                'search_physical_state',
                'search_exact_coc',
                'search_class_of_compounds',
                'search_risk_phrase',
                'search_safety_phrase',
                'search_hazard_statement',
                'search_precautionary_statement',
                'search_hazard_code',
                'search_symbol',
                'search_person_pc',
                'search_person_sc',
                'search_person_asc',
                'search_borrower',
                ]:
        if key in session:
            mylogger.debug(message='key:%s' % str(key))
            mylogger.debug(message='session[key]:%s' % str(session[key]))
            del session[key]

    session.search_result_per_page = result_per_page
    session.search_page = page

    #
    # display by product or storage
    #
    if 'display_by' in request.vars and request.vars['display_by'] == 'storage':
        session.search_display_by = 'storage'
        display_by_storage = True
        join_storage = True
    else:
        session.search_display_by = 'product'
        display_by_storage = False
    #
    # order by borrower or storage
    #
    if 'order_by' in request.vars and request.vars['order_by'] == 'storage':
        session.search_order_by = 'storage'
        order_by_storage = True
    else:
        session.search_order_by = 'borrower'
        order_by_storage = False

    mylogger.debug(message='request.vars:%s' % str(request.vars))

    #
    # building the request
    #
    if (not auth.has_permission('admin')) and (not auth.has_permission('read_rpc')):
        query_list.append(db.product.restricted_access == False)

    if request.vars:
        # product_id search just used after a new product card creation
        if 'product_id' in request.vars and request.vars['product_id'] is not None:
            session.search_product_id = request.vars['product_id']
            if type(request.vars['product_id']) is StringType:
                request.vars['product_id'] = [request.vars['product_id']]
            query_list.append(db.product.id.belongs(request.vars['product_id']))
            if request.vars['product_id'][0] != '-1':
                label += '%s %s<br/>' % ('product_id', request.vars['product_id'][0])
        if 'not_archive' in request.vars:
            session.search_archive = False
            query_list.append(db.storage.archive == False)
            join_storage = True
        if ('archive' in request.vars) and (auth.has_permission('read_archive') or auth.has_permission('admin')):
            session.search_archive = True
            query_list.append(db.storage.archive == True)
            join_storage = True
            label += ' %s<br/>' % cc.get_string("SEARCH_ARCHIVE")
        if 'bookmark' in request.vars:
            session.search_bookmark = True
            query_list.append(db.bookmark.person == auth.user.id)
            join_bookmark = True
            label += ' %s<br/>' % cc.get_string("BOOKMARK")
        if ('entity' in request.vars) and (len(request.vars['entity']) != 0):
            if type(request.vars['entity']) is StringType:
                request.vars['entity'] = [request.vars['entity']]
            if (not auth.has_permission('select_sc')) and (not auth.has_permission('admin')):
                raise HTTP(403, "Not authorized")
            session.search_entity = request.vars['entity']
            query_list.append(db.store_location.entity.belongs(tuple(request.vars['entity'])))
            query_list.append(db.storage.store_location == db.store_location.id)

            join_storage = True
            entity_label = " & ".join([e.role for e in db(db.entity.id.belongs(tuple(request.vars['entity']))).select(db.entity.role, cacheable=True)])
            label += ' %s<br/>' % entity_label
        if ('store_location' in request.vars) and (len(request.vars['store_location']) != 0):
            if type(request.vars['store_location']) is StringType:
                request.vars['store_location'] = [request.vars['store_location']]
            if (not auth.has_permission('read_sc')) and (not auth.has_permission('admin')):
                raise HTTP(403, "Not authorized")

            session.search_store_location = request.vars['store_location']

            _children_sl = []
            if 'include_children_store_location' in request.vars:
                session.search_include_children_store_location = request.vars['include_children_store_location']

                store_location_mapper = STORE_LOCATION_MAPPER()
                for _sl in request.vars['store_location']:
                    _childrens = store_location_mapper.find(store_location_id=_sl)[0].retrieve_children()
                    mylogger.debug(message='_childrens:%s' % str(_childrens))
                    _children_sl.extend([_child.id for _child in _childrens])
                    mylogger.debug(message='_children_sl:%s' % str(_children_sl))

                mylogger.debug(message='_children_sl:%s' % str(_children_sl))

            query_list.append(db.storage.archive == False)
            query_list.append(db.storage.store_location.belongs(tuple(request.vars['store_location'] + _children_sl)))
            join_storage = True
            sl_label = " | ".join([sl.label for sl in db(db.store_location.id.belongs(tuple(request.vars['store_location']))).select(db.store_location.label,
                                                                                                                                     cacheable=True)])
            label += ' %s' % sl_label

            if 'include_children_store_location' in request.vars:
                label += ' (%s)<br/>' % cc.get_string('SEARCH_INCLUDE_CHILDREN_STORE_LOCATION')
            else:
                label += '<br/>'

        if 'borrow_entity' in request.vars and request.vars['borrow_entity'] != '0':
            session.search_borrow_entity = request.vars['borrow_entity']
            query_list.append(db.store_location.entity.belongs(tuple(user_entity)))
            query_list.append(db.storage.store_location == db.store_location.id)
            join_borrow = True
            join_storage = True
            label += '%s<br/>' % cc.get_string("SEARCH_BORROW_ENTITY")
        if 'borrower' in request.vars and request.vars['borrower'] != '0':
            session.search_borrower = request.vars['borrower']
            query_list.append(db.borrow.borrower == request.vars['borrower'])
            join_borrow = True
            join_storage = True
            label += '%s<br/>' % cc.get_string("SEARCH_BORROW")
        if 'physical_state' in request.vars and request.vars['physical_state'] != '0':
            session.search_physical_state = request.vars['physical_state']
            query_list.append(db.product.physical_state == request.vars['physical_state'])
            ps_label = " & ".join([ps.label for ps in db(db.physical_state.id == request.vars['physical_state']).select(db.physical_state.label,
                                                                                                                        cacheable=True)])
            label += ' %s<br/>' % ps_label
        if 'is_cmr' in request.vars:
            session.search_is_cmr = True
            query_list.append(db.product.is_cmr == True)
            label += ' %s<br/>' % cc.get_string("DB_PRODUCT_IS_CMR_LABEL")
        if 'is_radio' in request.vars:
            session.search_is_radio = True
            query_list.append(db.product.is_radio == True)
            label += ' %s<br/>' % cc.get_string("DB_PRODUCT_IS_RADIO_LABEL")
        if 'to_destroy' in request.vars:
            session.search_to_destroy = True
            query_list.append(db.storage.to_destroy == True)
            join_storage = True
            label += ' %s<br/>' % cc.get_string("DB_STORAGE_TO_DESTROY_LABEL")
        if 'cas_number' in request.vars and request.vars['cas_number'] != '':
            session.search_cas_number = request.vars['cas_number']
            query_list.append(db.product.cas_number == request.vars['cas_number'])
            label += ' %s %s<br/>' % (cc.get_string("DB_PRODUCT_CAS_NUMBER_LABEL"), request.vars['cas_number'])
        if 'ce_number' in request.vars and request.vars['ce_number'] != '':
            session.search_ce_number = request.vars['ce_number']
            query_list.append(db.product.ce_number == request.vars['ce_number'])
            label += ' %s %s<br/>' % (cc.get_string("DB_PRODUCT_CE_NUMBER_LABEL"), request.vars['ce_number'])
        if 'product_datetime' in request.vars and request.vars['product_datetime'] != '':
            session.search_product_datetime = request.vars['product_datetime']
            query_list.append(db.product.creation_datetime >= datetime.strptime(request.vars['product_datetime'], '%s' % T('%Y-%m-%d %H:%M:%S')))
            label += ' %s %s<br/>' % (cc.get_string("DB_PRODUCT_CREATION_DATETIME_LABEL"), request.vars['product_datetime'])
        if 'entry_datetime' in request.vars and request.vars['entry_datetime'] != '':
            session.search_entry_datetime = request.vars['entry_datetime']
            query_list.append(db.storage.entry_datetime >= datetime.strptime(request.vars['entry_datetime'], '%s' % T('%Y-%m-%d %H:%M:%S')))
            join_storage = True
            label += ' %s %s<br/>' % (cc.get_string("DB_STORAGE_ENTRY_DATETIME_LABEL"), request.vars['entry_datetime'])
        if 'exit_datetime' in request.vars and request.vars['exit_datetime'] != '':
            session.search_exit_datetime = request.vars['exit_datetime']
            query_list.append(db.storage.exit_datetime <= datetime.strptime(request.vars['exit_datetime'], '%s' % T('%Y-%m-%d %H:%M:%S')))
            join_storage = True
            label += ' %s %s<br/>' % (cc.get_string("DB_STORAGE_EXIT_DATETIME_LABEL"), request.vars['exit_datetime'])
        if 'barecode' in request.vars and request.vars['barecode'] != '':
            session.search_barecode = request.vars['barecode']
            query_list.append(db.storage.barecode.like('%%%s%%' % request.vars['barecode'].strip()))
            join_storage = True
            label += ' %s %s<br/>' % (cc.get_string("DB_STORAGE_BARECODE_LABEL"), request.vars['barecode'])
        if 'comment' in request.vars and request.vars['comment'] != '':
            session.search_comment = request.vars['comment']
            query_list.append(db.storage.comment.like('%%%s%%' % request.vars['comment'].strip()))
            join_storage = True
            label += ' %s %s<br/>' % (cc.get_string("DB_STORAGE_COMMENT_LABEL"), request.vars['comment'])
        if 'name' in request.vars and request.vars['name'] != '':
            # lazy search - we replace "-" by "_" - we could do it better...
            session.search_name = request.vars['name']

            _search = request.vars['name'].upper().strip().replace('-', '_').replace("'", "%%").replace("`", "%%")
            _req1 = db.name.label.like('%%%s%%' % _search)
            _req2 = []
            for _name in db(db.name.label.like('%%%s%%' % _search)).select(db.name.id,
                                                                           cacheable=True):
                mylogger.debug(message='_name.id:%s' % _name.id)
                _req2.append(db.product.synonym.contains(_name.id))
            _req2.append(_req1)
            query_list.append(cc.or_ify(_req2))

            label += ' %s %s<br/>' % (cc.get_string("DB_NAME_LABEL"), request.vars['name'])

            # did you mean ?
            _did_you_mean = []
            _names = db(db.name).select(db.name.label_nost, cacheable=True)
            for _name in _names:
                _dist = Levenshtein.distance(_name.label_nost, request.vars['name'].upper().strip())
                if _dist <= 4:
                    _did_you_mean.append((_name.label_nost, _dist))
            if len(_did_you_mean) > 0:
                did_you_mean = sorted(_did_you_mean, key=lambda k: k[1])

        if 'empirical_formula' in request.vars and request.vars['empirical_formula']:
            session.search_empirical_formula = request.vars['empirical_formula']
            query_list.append(db.product.empirical_formula == request.vars['empirical_formula'])
            ef_label = db(db.empirical_formula.id == request.vars['empirical_formula']).select(db.empirical_formula.label, cacheable=True).first().label
            label += ' %s %s<br/>' % (cc.get_string("DB_PRODUCT_EMPIRICAL_FORMULA_LABEL"), ef_label)
        if 'linear_formula' in request.vars and request.vars['linear_formula']:
            session.search_linear_formula = request.vars['linear_formula']
            query_list.append(db.product.linear_formula == request.vars['linear_formula'])
            lf_label = db(db.linear_formula.id == request.vars['linear_formula']).select(db.linear_formula.label, cacheable=True).first().label
            label += ' %s %s<br/>' % (cc.get_string("DB_PRODUCT_LINEAR_FORMULA_LABEL"), lf_label)
        if 'class_of_compounds' in request.vars and request.vars['class_of_compounds'] != '':
            if type(request.vars['class_of_compounds']) is StringType:
                request.vars['class_of_compounds'] = [request.vars['class_of_compounds']]
            subQuery = (db.product.class_of_compounds == -1)
            if 'exact_coc' in request.vars:
                session.search_exact_coc = True
                subQuery = subQuery.__or__(db.product.class_of_compounds.contains(request.vars['class_of_compounds'], all=True))
                for coc in db(~db.class_of_compounds.id.belongs(request.vars['class_of_compounds'])).select(cacheable=True):
                    subQuery = subQuery.__and__(~db.product.class_of_compounds.contains(coc.id))
                coc_label = 'exact '
            else:
                for class_of_compounds in request.vars['class_of_compounds']:
                    subQuery = subQuery.__or__(db.product.class_of_compounds.contains(class_of_compounds))
                coc_label = ''
            query_list.append(subQuery)
            session.search_class_of_compounds = request.vars['class_of_compounds']
            coc_label += " & ".join([coc.label for coc in db(db.class_of_compounds.id.belongs(tuple(request.vars['class_of_compounds']))).select(db.class_of_compounds.label,
                                                                                                                                                 cacheable=True)])
            label += ' %s<br/>' % coc_label
        if 'risk_phrase' in request.vars:
            if type(request.vars['risk_phrase']) is StringType:
                request.vars['risk_phrase'] = [request.vars['risk_phrase']]
            subQuery = (db.product.risk_phras == -1)
            for risk_phrase in request.vars['risk_phrase']:
                subQuery = subQuery.__or__(db.product.risk_phrase.contains(risk_phrase))
            query_list.append(subQuery)
            session.search_risk_phrase = request.vars['risk_phrase']
            rp_label = " & ".join([rp.reference for rp in db(db.risk_phrase.id.belongs(tuple(request.vars['risk_phrase']))).select(db.risk_phrase.reference,
                                                                                                                                   cacheable=True)])
            label += ' %s %s<br/>' % (cc.get_string("DB_PRODUCT_RISK_PHRASE_LABEL"), rp_label)
        if 'safety_phrase' in request.vars:
            if type(request.vars['safety_phrase']) is StringType:
                request.vars['safety_phrase'] = [request.vars['safety_phrase']]
            subQuery = (db.product.safety_phrase == -1)
            for safety_phrase in request.vars['safety_phrase']:
                subQuery = subQuery.__or__(db.product.safety_phrase.contains(safety_phrase))
            query_list.append(subQuery)
            session.search_safety_phrase = request.vars['safety_phrase']
            sp_label = " & ".join([sp.reference for sp in db(db.safety_phrase.id.belongs(tuple(request.vars['safety_phrase']))).select(db.safety_phrase.reference,
                                                                                                                                       cacheable=True)])
            label += ' %s %s<br/>' % (cc.get_string("DB_PRODUCT_SAFETY_PHRASE_LABEL"), sp_label)
        if 'hazard_statement' in request.vars:
            if type(request.vars['hazard_statement']) is StringType:
                request.vars['hazard_statement'] = [request.vars['hazard_statement']]
            subQuery = (db.product.hazard_statement == -1)
            for hazard_statement in request.vars['hazard_statement']:
                subQuery = subQuery.__or__(db.product.hazard_statement.contains(hazard_statement))
            query_list.append(subQuery)
            session.search_hazard_statement = request.vars['hazard_statement']
            hs_label = " & ".join([hs.reference for hs in db(db.hazard_statement.id.belongs(tuple(request.vars['hazard_statement']))).select(db.hazard_statement.reference,
                                                                                                                                             cacheable=True)])
            label += ' %s %s<br/>' % (cc.get_string("DB_PRODUCT_HAZARD_STATEMENT_LABEL"), hs_label)
        if 'precautionary_statement' in request.vars:
            if type(request.vars['precautionary_statement']) is StringType:
                request.vars['precautionary_statement'] = [request.vars['precautionary_statement']]
            subQuery = (db.product.precautionary_statement == -1)
            for precautionary_statement in request.vars['precautionary_statement']:
                subQuery = subQuery.__or__(db.product.precautionary_statement.contains(precautionary_statement))
            query_list.append(subQuery)
            session.search_precautionary_statement = request.vars['precautionary_statement']
            ps_label = " & ".join([ps.reference for ps in db(db.precautionary_statement.id.belongs(tuple(request.vars['precautionary_statement']))).select(db.precautionary_statement.reference,
                                                                                                                                                           cacheable=True)])
            label += ' %s %s<br/>' % (cc.get_string("DB_PRODUCT_PRECAUTIONARY_STATEMENT_LABEL"), ps_label)
        if 'hazard_code' in request.vars:
            if type(request.vars['hazard_code']) is StringType:

                request.vars['hazard_code'] = [request.vars['hazard_code']]
            subQuery = (db.product.hazard_code == -1)
            for hazard_code in request.vars['hazard_code']:
                subQuery = subQuery.__or__(db.product.hazard_code.contains(hazard_code))
            query_list.append(subQuery)
            session.search_hazard_code = request.vars['hazard_code']
            hc_label = " & ".join([hc.label for hc in db(db.hazard_code.id.belongs(tuple(request.vars['hazard_code']))).select(db.hazard_code.label,
                                                                                                                               cacheable=True)])
            label += ' %s %s<br/>' % (cc.get_string("DB_PRODUCT_HAZARD_CODE_LABEL"), hc_label)
        if 'symbol' in request.vars:
            if type(request.vars['symbol']) is StringType:
                request.vars['symbol'] = [request.vars['symbol']]
            subQuery = (db.product.symbol == -1)
            for symbol in request.vars['symbol']:
                subQuery = subQuery.__or__(db.product.symbol.contains(symbol))
            query_list.append(subQuery)
            session.search_symbol = request.vars['symbol']
            symbol_label = " & ".join([symbol.label for symbol in db(db.symbol.id.belongs(tuple(request.vars['symbol']))).select(db.symbol.label,
                                                                                                                                 cacheable=True)])
            label += ' %s %s<br/>' % (cc.get_string("DB_PRODUCT_SYMBOL_LABEL"), symbol_label)
        if ('person_pc' in request.vars and request.vars['person_pc'] != ''):
            session.search_person_pc = request.vars['person_pc']
            query_list.append(db.product.person == request.vars['person_pc'])
        if ('person_sc' in request.vars and request.vars['person_sc'] != ''):
            session.search_person_sc = request.vars['person_sc']
            query_list.append(db.storage.person == request.vars['person_sc'])
            join_storage = True
        if ('person_asc' in request.vars and request.vars['person_asc'] != ''):
            session.search_person_asc = request.vars['person_asc']
            query_list.append(db.storage.person == request.vars['person_asc'])
            join_storage = True

    #
    # join with the name table
    #
    query_list.append(db.product.name == db.name.id)

    #
    # performing additional join if needed
    #
    if join_storage:
        query_list.append(db.product.id == db.storage.product)
        if display_by_storage:
            query_list.append(db.storage.store_location == db.store_location.id)
            query_list.append(db.storage.archive == False)
            query_list.append(db.store_location.entity.belongs(user_entity))
    if join_borrow:
        query_list.append(db.storage.id == db.borrow.storage)
    if join_bookmark:
        query_list.append(db.bookmark.product == db.product.id)

    mylogger.debug(message='query_list:%s' % query_list)

    if len(query_list) != 0:
        #
        # building the final query
        #
        finalQuery = (db.product.id > 0)
        for query in query_list:
            mylogger.debug(message='query:%s' % str(query))
            finalQuery = finalQuery.__and__(query)
        mylogger.debug(message='finalQuery:%s' % str(finalQuery))

        #
        # pagination
        #
        range_min = page * result_per_page
        range_max = range_min + result_per_page
        mylogger.debug(message='page:%s' % page)
        mylogger.debug(message='result_per_page:%s' % result_per_page)
        mylogger.debug(message='range_min:%s' % range_min)
        mylogger.debug(message='range_min:%s' % range_max)

        theset = db(finalQuery)
        if not display_by_storage:
            nb_entries = theset.count(distinct=db.product.id)
        else:
            nb_entries = theset.count()
        mylogger.debug(message='nb_entries:%i' % nb_entries)

        paginate_selector = PaginateSelector(anchor='main')
        paginator = Paginator(paginate=paginate_selector.paginate,
                              extra_vars={'keep_last_search': True},
                              anchor='main',
                              renderstyle=False)
        paginator.records = nb_entries
        paginate_info = PaginateInfo(paginator.page, paginator.paginate, paginator.records)

        #
        # executing the query
        #
        left = None
        if (not export_csv) and (not export_html):
            _limitby = paginator.limitby()
        else:
            _limitby = None

        select_fields = [db.product.ALL, db.name.ALL]

        if export_csv:
            select_fields.append(db.storage.ALL)
            select_fields.append(db.borrow.ALL)
            left = db.borrow.on(db.storage.id == db.borrow.storage)

        if (('entity' in request.vars) or ('store_location' in request.vars)) and (export_csv or export_html):
            select_fields.append(db.storage.to_destroy)
            select_fields.append(db.borrow.borrower)

        if display_by_storage:
            if order_by_storage:
                mylogger.debug(message='order by storage')
                select_fields.append(db.storage.ALL)
                _order_by = db.storage.store_location
            else:
                mylogger.debug(message='order by borrower')
                select_fields.append(db.borrow.ALL)
                select_fields.append(db.storage.ALL)
                left = db.borrow.on(db.storage.id == db.borrow.storage)
                _order_by = (db.borrow.borrower|db.storage.store_location)

        if is_did_you_mean:
            _order_by = (~(db.name.label == request.vars['name'])).case() | db.name.label_nost
            _distinct = False
        elif display_by_storage:
            #_order_by = db.storage.store_location
            _distinct = False
        else:
            _order_by = db.name.label_nost
            _distinct = True

        mylogger.debug(message='_order_by:%s' % _order_by)
        allrows = theset.select(*select_fields,
                                orderby=_order_by,
                                left=left,
                                # distinct generates error with the CASE clause in POSTGRES
                                distinct=_distinct,
                                limitby=_limitby,
                                cacheable=True)
        rows = allrows

        storages = None
        products = None
        mylogger.debug(message='len(rows):%s' % len(rows))
        if len(rows) > 0:
            if not display_by_storage:
                products = PRODUCT_MAPPER().find(product_id=[row.product.id for row in rows], orderby=_order_by)
                mylogger.debug(message='len(products):%s' % len(products))
            else:
                storages = STORAGE_MAPPER().find(storage_id=[row.storage.id for row in rows], archive=session.search_archive, orderby=_order_by)
                mylogger.debug(message='len(storages):%s' % len(storages))

        #
        # export
        #
        if export_csv:

            export_row = []
            fields_order = ['product.name',
                            'product.synonym',
                            'product.specificity',
                            'product.restricted_access',
                            'product.creation_datetime',
                            'product.archive',
                            'product.person',
                            'product.cas_number',
                            'product.ce_number',
                            'product.empirical_formula',
                            'product.linear_formula',
                            'product.td_formula',
                            'product.msds',
                            'product.is_cmr',
                            'product.is_radio',
                            'product.cmr_cat',
                            'product.class_of_compounds',
                            'product.physical_state',
                            'product.risk_phrase',
                            'product.safety_phrase',
                            'product.hazard_statement',
                            'product.precautionary_statement',
                            'product.hazard_code',
                            'product.signal_word',
                            'product.symbol',
                            'product.remark',
                            'product.disposal_comment',
                            'storage.store_location',
                            'storage.volume_weight',
                            'storage.unit',
                            'storage.barecode',
                            'storage.comment',
                            'storage.batch_number',
                            'storage.supplier',
                            'storage.creation_datetime',
                            'storage.entry_datetime',
                            'storage.exit_datetime',
                            'storage.opening_datetime',
                            'storage.person',
                            'storage.archive',
                            'borrow.borrower']

            fields_banned = ['name.id',
                             'product.id',
                             'name.label',
                             'name.label_nost',
                             'storage.id',
                             'storage.reference',
                             'storage.computed_entity',
                             'storage.nb_items',
                             'storage.product',
                             'borrow.storage',
                             'borrow.id',
                             'borrow.person']

            _total = len(rows)
            _current = 1
            mylogger.ram('start building export...', index=0)
            for row in rows:
                mylogger.debug(message='row:%s' % str(row))
                mylogger.ram('built %s rows from %s' % (_current, _total), index=1)
                _current = _current + 1

                _leafs = cc.flatten_row(row)
                mylogger.debug(message='_leafs:%s' % str(_leafs))
                for _f in fields_banned:
                    if _f in _leafs.keys():
                        del _leafs[_f]
                _leafs = OrderedDict(sorted(_leafs.items(), key=lambda k: fields_order.index(k[0]) if k[0] in fields_order else 99))
                mylogger.debug(message='_leafs:%s' % str(_leafs))

                export_row.append(_leafs.values())

            mylogger.ram('done!', index=3)
            mylogger.ram('[END]', index=4)
            field_names = _leafs.keys()

            mylogger.debug(message='field_names:%s' % str(field_names))
            mylogger.debug(message='export_row:%s' % str(export_row))
            response.view = 'product/export_chimitheque.csv'

            return dict(filename='export_chimitheque.csv',
                        csvdata=export_row,
                        field_names=field_names)

    #
    # building the search form
    #
    db.product.is_cmr.writable = True
    db.product.is_radio.writable = True
    db.product.name.widget = SQLFORM.widgets.string.widget
    db.product.empirical_formula.widget = CHIMITHEQUE_MULTIPLE_widget(db.empirical_formula.label, minchar=1, configuration={'search': {'func_lambda': 'lambdaempiricalformula'}, 'create': {'add_in_db': True}, 'update': {'add_in_db': True}})
    db.product.linear_formula.widget = CHIMITHEQUE_MULTIPLE_widget(db.linear_formula.label)
    db.product.cas_number.widget = lambda field, value: SQLFORM.widgets.string.widget(field, value)  # removing the "required" id
    db.product.risk_phrase.widget = CHIMITHEQUE_MULTIPLE_widget(db.risk_phrase.reference, minchar=1, configuration={'*': {'multiple': True}})
    db.product.risk_phrase.comment = cc.get_string("SEARCH_RISK_PHRASE")
    db.product.safety_phrase.widget = CHIMITHEQUE_MULTIPLE_widget(db.safety_phrase.reference, minchar=1, configuration={'*': {'multiple': True}})
    db.product.safety_phrase.comment = cc.get_string("SEARCH_SAFETY_PHRASE")
    db.product.hazard_statement.widget = CHIMITHEQUE_MULTIPLE_widget(db.hazard_statement.reference, minchar=1, configuration={'*': {'multiple': True}})
    db.product.hazard_statement.comment = cc.get_string("SEARCH_HAZARD_STATEMENT")
    db.product.precautionary_statement.widget = CHIMITHEQUE_MULTIPLE_widget(db.precautionary_statement.reference, minchar=1, configuration={'*': {'multiple': True}})
    db.product.precautionary_statement.comment = cc.get_string("SEARCH_PRECAUTIONARY_STATEMENT")

    # prepopulating form values + default values in the following form declaration
    db.product.name.default = request.vars['name']
    db.product.empirical_formula.default = request.vars['empirical_formula']
    db.product.linear_formula.default = request.vars['linear_formula']
    db.product.cas_number.default = request.vars['cas_number']
    db.product.ce_number.default = request.vars['ce_number']
    db.product.class_of_compounds.default = request.vars['class_of_compounds']
    db.product.physical_state.default = request.vars['physical_state']
    db.product.hazard_code.default = request.vars['hazard_code']
    db.product.symbol.default = request.vars['symbol']
    db.product.safety_phrase.default = request.vars['safety_phrase']
    db.product.risk_phrase.default = request.vars['risk_phrase']
    db.product.hazard_statement.default = request.vars['hazard_statement']
    db.product.precautionary_statement.default = request.vars['precautionary_statement']

    form = SQLFORM.factory(Field('store_location',
                                 'list:reference store_location',
                                 default=request.vars['store_location'],
                                 widget=lambda field, value: SQLFORM.widgets.options.widget(field, value, _multiple="multiple", _size='12'),
                                 requires=IS_IN_DB_AND_USER_STORE_LOCATION(db((db.store_location.id > 0) &
                                                                              (db.store_location.can_store == True)),
                                                                           db.store_location.id,
                                                                           db.store_location._format,
                                                                           multiple=True,
                                                                           orderby=db.store_location.label_full_path)),
                           Field('include_children_store_location',
                                 'boolean',
                                 default=session.search_include_children_store_location),
                           Field('entity',
                                 'list:reference entity',
                                 default=request.vars['entity'],
                                 requires=IS_IN_DB_AND_USER_ENTITY(db(db.entity.id > 0),
                                                                   db.entity.id,
                                                                   db.entity._format,
                                                                   multiple=True,
                                                                   orderby=db.entity.role)),
                           Field('is_cmr',
                                 'boolean',
                                 default=session.search_is_cmr),
                           Field('is_radio',
                                 'boolean',
                                 default=session.search_is_radio),
                           Field('archive',
                                 'boolean',
                                 default=session.search_archive),
                           Field('bookmark',
                                 'boolean',
                                 default=session.search_bookmark),
                           db.product.name,
                           db.product.cas_number,
                           db.product.ce_number,
                           Field('product_datetime',
                                 'datetime',
                                 default=datetime.strptime(request.vars['product_datetime'],
                                                           '%Y-%m-%d %H:%M:%S')
                                 if ('product_datetime' in request.vars and request.vars['product_datetime'] != '')
                                 else None),
                           Field('entry_datetime',
                                 'datetime',
                                 default=datetime.strptime(request.vars['entry_datetime'],
                                                           '%Y-%m-%d %H:%M:%S')
                                 if ('entry_datetime' in request.vars and request.vars['entry_datetime'] != '')
                                 else None),
                           Field('exit_datetime',
                                 'datetime',
                                 default=datetime.strptime(request.vars['exit_datetime'],
                                                           '%Y-%m-%d %H:%M:%S')
                                 if ('exit_datetime' in request.vars and request.vars['exit_datetime'] != '')
                                 else None),
                           Field('barecode',
                                 'string',
                                 default=request.vars['barecode']),
                           Field('comment',
                                 'string',
                                 default=request.vars['comment']),
                           db.product.empirical_formula,
                           db.product.linear_formula,
                           Field('to_destroy',
                                 'boolean',
                                 default=session.search_to_destroy),
                           db.product.physical_state,
                           Field('exact_coc',
                                 'boolean',
                                 default=request.vars['exact_coc']),
                           db.product.class_of_compounds,
                           db.product.hazard_code,
                           db.product.symbol,
                           db.product.risk_phrase,
                           db.product.safety_phrase,
                           db.product.hazard_statement,
                           db.product.precautionary_statement,
                           _action='/%s/%s/search' % (request.application, request.controller),
                           submit_button=cc.get_string("SEARCH"))

    mylogger.debug(message='request.vars:%s' % request.vars)

    return dict(form=form,
                products=products,
                storages=storages,
                auth_person=PERSON_MAPPER().find(person_id=auth.user.id)[0],
                nb_entries=nb_entries,
                label=label,
                did_you_mean=did_you_mean,
                is_did_you_mean=is_did_you_mean,
                paginator=paginator,
                paginate_selector=paginate_selector,
                paginate_info=paginate_info)
