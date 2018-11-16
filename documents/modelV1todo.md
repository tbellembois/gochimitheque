```bash

db.define_table('person',
                Field('creator',
                      'reference person',
                      ondelete='NO ACTION',
                      label=cc.get_string("DB_PERSON_CREATOR_LABEL"),
                      comment=cc.get_string("DB_PERSON_CREATOR_COMMENT"),
                      compute=lambda r: db.person[auth.user.id] if auth.user else None,
                      represent=lambda r: str(db(db.person.id == r).select(db.person.email).first().email) if r else None),
              
                Field('contact',
                      'text',
                      label=cc.get_string("DB_PERSON_CONTACT_LABEL"),
                      comment=cc.get_string("DB_PERSON_CONTACT_COMMENT")),
                Field('password',
                      'password',
                      label=cc.get_string("DB_PERSON_PASSWORD_LABEL"),
                      comment=cc.get_string("DB_PERSON_PASSWORD_COMMENT"),
                      writable=False,
                      readable=False,
                      required=True,
                      notnull=True),
                Field('creation_date',
                      'date',
                      label=cc.get_string("DB_PERSON_CREATION_DATE_LABEL"),
                      comment=cc.get_string("DB_PERSON_CREATION_DATE_COMMENT"),
                      default=datetime.now(),
                      writable=False,
                      readable=True),
                Field('virtual',
                      'boolean',
                      writable=False,
                      readable=False,
                      default=False),
                Field('archive',
                      'boolean',
                      writable=False,
                      readable=False,
                      default=False),
                Field('exposure_card',
                      'list:reference exposure_card',
                      writable=False,
                      readable=False),
 
                # web2py auth
                format=lambda r: r.email)

db.entity.manager.requires = [IS_EMPTY_OR(IS_LIST_OF(CLEANUP()))]

db.person.first_name.requires = IS_NOT_EMPTY()
db.person.last_name.requires = IS_NOT_EMPTY()
db.person.email.requires = [IS_NOT_EMPTY(), IS_EMAIL(), IS_NOT_IN_DB(db, db.person.email)]
db.person.password.requires = [IS_NOT_EMPTY(), CRYPT(key=settings['hmac_key'])]
db.person.registration_key.requires = [IS_ONE_SELECTED(tuple_list=[('unactive', 'unactive'), ('active', 'active')])]

db.define_table('product',
                Field('is_cmr',
                      'boolean',
                      label=cc.get_string("DB_PRODUCT_IS_CMR_LABEL"),
                      comment=cc.get_string("DB_PRODUCT_IS_CMR_COMMENT"),
                      compute=lambda r: compute_product_is_cmr(r),
                      default=False),
                Field('cmr_cat',
                      'string',
                      label=cc.get_string("DB_PRODUCT_CMR_CATEGORY_LABEL"),
                      compute=lambda r: compute_product_cmr_cat(r),
                      writable=False,
                      default=None,
                      represent=lambda r: r.replace('|', ' ') if r else None),
                Field('archive',
                      'boolean',
                      writable=False,
                      default=False),
                Field('creation_datetime',
                      'datetime',
                      label=cc.get_string("DB_PRODUCT_CREATION_DATETIME_LABEL"),
                      comment=cc.get_string("DB_PRODUCT_CREATION_DATETIME_COMMENT"),
                      default=datetime.now(),
                      writable=False,
                      readable=True),
                format=lambda r: represent_product(r))

db.product.hazard_code.requires = IS_EMPTY_OR(IS_IN_DB(db, db.hazard_code.id, db.symbol._format, multiple=True))
db.product.symbol.requires = IS_EMPTY_OR(IS_IN_DB(db, db.symbol.id, db.symbol._format, multiple=True))
db.product.cas_number.requires = IS_CONFIRM_EMPTY_OR('cas_number',
                                                     [IS_VALID_CAS(),
                                                      IS_UNIQUE_WITH_SPECIFICITY(request.vars.specificity, request.function)])
db.product.ce_number.requires = IS_EMPTY_OR(CLEANUP())
db.product.physical_state.requires = IS_EMPTY_OR(IS_IN_DB(db, db.physical_state.id, '%(label)s'))
db.product.class_of_compounds.requires = [IS_EMPTY_OR(IS_IN_DB(db, db.class_of_compounds.id,
                                                               label=db.class_of_compounds._format,
                                                               multiple=True,
                                                               sort=db.class_of_compounds.label))]
db.product.signal_word.requires = IS_EMPTY_OR(IS_IN_DB(db, db.signal_word.id, '%(label)s'))
db.product.name.requires = IS_NOT_EMPTY()
db.product.synonym.requires = IS_EMPTY_OR(IS_LIST_OF(CLEANUP()))
db.product.risk_phrase.requires = IS_EMPTY_OR(IS_IN_DB(db, db.risk_phrase.id, '(%(reference)s) %(label)s', multiple=True))
db.product.safety_phrase.requires = IS_EMPTY_OR(IS_IN_DB(db, db.safety_phrase.id, '(%(reference)s) %(label)s', multiple=True))
db.product.hazard_statement.requires = IS_EMPTY_OR(IS_IN_DB(db, db.hazard_statement.id, '(%(reference)s) %(label)s', multiple=True))
db.product.precautionary_statement.requires = IS_EMPTY_OR(IS_IN_DB(db, db.precautionary_statement.id, '(%(reference)s) %(label)s', multiple=True))
db.product.msds.requires = IS_CONFIRM_EMPTY_OR('msds', IS_NOT_EMPTY())
db.product.empirical_formula.requires = IS_CONFIRM_EMPTY_OR('empirical_formula', IS_IN_DB(db, db.empirical_formula.id, '%(label)s'))
db.product.linear_formula.requires = IS_EMPTY_OR(CLEANUP())
db.product.remark.widget = lambda field, value: SQLFORM.widgets.text.widget(field, value, _rows=5)
db.define_table('product_history',
                Field('current_record', db.product),
                Field('modification_datetime', 'datetime', writable=False, default=datetime.now()),
                db.product)

db.define_table('bookmark',
                Field('person',
                      db.person,
                      compute=lambda r: db.person[auth.user.id] if auth.user else None),
                Field('product',
                      db.product))

db.define_table('storage',
                Field('archive',
                      'boolean',
                      label=cc.get_string("DB_STORAGE_ARCHIVE_LABEL"),
                      comment=cc.get_string("DB_STORAGE_ARCHIVE_COMMENT"),
                      writable=False,
                      default=False),
                # computed field to make coding easier
                Field('computed_entity',
                      'integer',
                      compute=lambda r: db(db.store_location.id == (r['STORE_LOCATION'])).select(db.store_location.entity).first().entity if r else None,
                      writable=False,
readable=False))

db.storage.product.requires = IS_NOT_EMPTY()
db.storage.store_location.requires = IS_IN_DB_AND_USER_STORE_LOCATION(db(db.store_location.can_store==True), db.store_location.id, db.store_location._format, orderby=db.store_location.label_full_path)
db.storage.volume_weight.requires = IS_EMPTY_OR(IS_FLOAT_IN_RANGE(cc.MIN_FLOAT, cc.MAX_FLOAT))
# prevent users from giving a volume_weight without a unit
db.storage.unit.requires = IS_IN_DB(db, db.unit.id, db.unit._format) if request.vars.volume_weight != '' else IS_EMPTY_OR(IS_IN_DB(db, db.unit.id, db.unit._format))
db.storage.nb_items.requires = IS_EMPTY_OR(IS_INT_IN_RANGE(1, 31))
db.storage.supplier.requires = IS_EMPTY_OR(IS_IN_DB(db, db.supplier.id, '%(label)s'))

db.define_table('storage_history',
                Field('current_record', db.storage),
                Field('modification_datetime', 'datetime', writable=False, default=datetime.now()),
db.storage)

db.define_table('stock',
                Field('maximum',
                      'double',
                      label=cc.get_string("DB_STOCK_MAXIMUM_LABEL"),
                      comment=cc.get_string("DB_STOCK_MAXIMUM_COMMENT")),
                Field('maximum_unit',
                      db.unit,
                      label=cc.get_string("DB_STOCK_MAXIMUM_UNIT_LABEL"),
                      comment=cc.get_string("DB_STOCK_MAXIMUM_UNIT_COMMENT"),
                      represent=lambda r: str(db(db.unit.id == r).select(db.unit.label).first().label) if r else None),
                Field('minimum',
                      'double',
                      label=cc.get_string("DB_STOCK_MINIMUM_LABEL"),
                      comment=cc.get_string("DB_STOCK_MINIMUM_COMMENT")),
                Field('minimum_unit',
                      db.unit,
                      label=cc.get_string("DB_STOCK_MINIMUM_UNIT_LABEL"),
                      comment=cc.get_string("DB_STOCK_MINIMUM_UNIT_COMMENT"),
                      represent=lambda r: str(db(db.unit.id == r).select(db.unit.label).first().label) if r else None),
                Field('product',
                      db.product,
                      label=cc.get_string("DB_STOCK_PRODUCT_LABEL"),
                      comment=cc.get_string("DB_STOCK_PRODUCT_COMMENT"),
                      writable=False),
                Field('entity',
                      db.entity,
                      label=cc.get_string("DB_STOCK_ENTITY_LABEL"),
                      comment=cc.get_string("DB_STOCK_ENTITY_COMMENT"),
                      writable=False))

db.stock.minimum.requires = [IS_NOT_EMPTY(), IS_FLOAT_IN_RANGE(cc.MIN_FLOAT, cc.MAX_FLOAT)]
db.stock.maximum.requires = [IS_NOT_EMPTY(), IS_FLOAT_IN_RANGE(cc.MIN_FLOAT, cc.MAX_FLOAT)]

db.define_table('cpe',
                Field('label',
                      'string',
                      length=255,
                      label=cc.get_string("DB_CPE_LABEL"),
                      comment=cc.get_string("DB_CPE_COMMENT"),
                      required=True,
                      notnull=True,
                      unique=True,
                      represent=lambda r: r.label),
                format=lambda r: r.label)

db.define_table('ppe',
                Field('label',
                      'string',
                      length=255,
                      label=cc.get_string("DB_PPE_LABEL"),
                      comment=cc.get_string("DB_PPE_COMMENT"),
                      required=True,
                      notnull=True,
                      unique=True,
                      represent=lambda r: r.label),
format=lambda r: r.label)

db.define_table('exposure_item',
                Field('creation_datetime',
                      'datetime',
                      label=cc.get_string("DB_EXPOSURE_ITEM_CREATION_DATETIME_LABEL"),
                      comment=cc.get_string("DB_EXPOSURE_ITEM_CREATION_DATETIME_COMMENT"),
                      default=datetime.now(),
                      writable=False,
                      readable=True),
                Field('product',
                      db.product,
                      label=cc.get_string("DB_EXPOSURE_ITEM_PRODUCT_LABEL"),
                      comment=cc.get_string("DB_EXPOSURE_ITEM_PRODUCT_COMMENT"),
                      required=True,
                      notnull=True,
                      represent=lambda r: represent_product(r)),
                Field('kind_of_work',
                      'text',
                      label=cc.get_string("DB_EXPOSURE_ITEM_KIND_OF_WORK_LABEL"),
                      comment=cc.get_string("DB_EXPOSURE_ITEM_KIND_OF_WORK_COMMENT")),
                Field('cpe',
                      'list:reference cpe',
                      label=cc.get_string("DB_EXPOSURE_ITEM_CPE_LABEL"),
                      comment=cc.get_string("DB_EXPOSURE_ITEM_CPE_COMMENT"),
                      represent=lambda r: XML(' <br/>'.join(['%s' %(T(row.label)) \
                                              for row in db(db.cpe.id.belongs(r)).select()])) \
                                              if r else None),
                Field('ppe',
                      'list:reference ppe',
                      label=cc.get_string("DB_EXPOSURE_ITEM_PPE_LABEL"),
                      comment=cc.get_string("DB_EXPOSURE_ITEM_PPE_COMMENT"),
                      represent=lambda r: XML(' <br/>'.join(['%s' %(T(row.label)) \
                                              for row in db(db.ppe.id.belongs(r)).select()])) \
                                              if r else None),
                Field('nb_exposure',
                      'integer',
                      label=cc.get_string("DB_EXPOSURE_ITEM_NB_EXPOSURE_LABEL"),
                      comment=cc.get_string("DB_EXPOSURE_ITEM_NB_EXPOSURE_COMMENT"),
                      default=1),
                Field('exposure_time',
                      'time',
                      label=cc.get_string("DB_EXPOSURE_ITEM_EXPOSURE_TIME_LABEL"),
                      comment=cc.get_string("DB_EXPOSURE_ITEM_EXPOSURE_TIME_COMMENT")),
                Field('simultaneous_risk',
                      'text',
                      label=cc.get_string("DB_EXPOSURE_ITEM_SIMULTANEAOUS_RISK_LABEL"),
                      comment=cc.get_string("DB_EXPOSURE_ITEM_SIMULTANEAOUS_RISK_COMMENT")))

db.define_table('exposure_card',
                Field('title',
                      'string',
                      default='card: %s' % datetime.now()),
                Field('accidental_exposure_type',
                      'text',
                       label=cc.get_string("DB_EXPOSURE_CARD_ACCIDENTAL_EXPOSURE_TYPE_LABEL"),
                       comment=cc.get_string("DB_EXPOSURE_CARD_ACCIDENTAL_EXPOSURE_TYPE_COMMENT")),
                Field('accidental_exposure_datetime',
                      'datetime',
                       label=cc.get_string("DB_EXPOSURE_CARD_ACCIDENTAL_EXPOSURE_DATETIME_LABEL"),
                       comment=cc.get_string("DB_EXPOSURE_CARD_ACCIDENTAL_EXPOSURE_DATETIME_COMMENT")),
                Field('accidental_exposure_duration_and_extent',
                      'text',
                       label=cc.get_string("DB_EXPOSURE_CARD_ACCIDENTAL_EXPOSURE_DURATION_AND_EXTENT_LABEL"),
                       comment=cc.get_string("DB_EXPOSURE_CARD_ACCIDENTAL_EXPOSURE_DURATION_AND_EXTENT_COMMENT")),
                Field('creation_datetime',
                      'datetime',
                      default=datetime.now(),
                      writable=False,
                      readable=True),
                Field('modification_datetime',
                      'datetime',
                      default=datetime.now(),
                      compute=lambda r: datetime.now(),
                      writable=False,
                      readable=True),
                Field('archive',
                      'boolean',
                      writable=False,
                      default=False),
                Field('exposure_item',
'list:reference exposure_item'))

db.exposure_item.product.requires = []
db.exposure_item.cpe.requires = IS_EMPTY_OR(IS_IN_DB(db, db.cpe.id, '%(label)s', multiple=True))
db.exposure_item.ppe.requires = IS_EMPTY_OR(IS_IN_DB(db, db.ppe.id, '%(label)s', multiple=True))
db.exposure_item.nb_exposure.requires = IS_INT_IN_RANGE(1, 3650)
db.exposure_item.exposure_time.requires = IS_TIME()

db.define_table('borrow',
                Field('creation_datetime',
                      'datetime',
                      label=cc.get_string("DB_USE_CREATION_DATETIME_LABEL"),
                      comment=cc.get_string("DB_USE_CREATION_DATETIME_COMMENT"),
                      default=datetime.now(),
                      writable=False,
                      readable=True),
                Field('person',
                      db.person,
                      label=cc.get_string("DB_USE_PERSON_LABEL"),
                      comment=cc.get_string("DB_USE_PERSON_COMMENT"),
                      compute=lambda r: db.person[auth.user.id] if auth.user else None,
                      writable=False,
                      readable=True,
                      represent=lambda r: str(db(db.person.id == r).select(db.person.email).first().email) if r else None),
                Field('borrower',
                      db.person,
                      label=cc.get_string("DB_USE_BORROWER_LABEL"),
                      comment=cc.get_string("DB_USE_BORROWER_COMMENT"),
                      represent=lambda r: str(db(db.person.id == r).select(db.person.email).first().email) if r else None,
                      required=True,
                      notnull=True),
                Field('storage',
                      db.storage,
                      label=cc.get_string("DB_USE_STORAGE_LABEL"),
                      comment=cc.get_string("DB_USE_STORAGE_COMMENT"),
                      writable=False,
                      readable=False),
                Field('comment',
                      'text',
                      label=cc.get_string("DB_USE_COMMENT_LABEL"),
                      comment=cc.get_string("DB_USE_COMMENT_COMMENT")))

db.borrow.borrower.requires = [IS_IN_DB(db, db.person.id, label=db.person._format, sort=db.person.email)]
```