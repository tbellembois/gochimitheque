package models

const PRECAUTIONARYSTATEMENT = `
.Absorb spillage to prevent material damage.	P390
.Avoid breathing dust/fume/gas/mist/vapours/spray.	P261
.Avoid contact during pregnancy/while nursing.	P263
.Avoid release to the environment.	P273
.Brush off loose particles from skin.	P335
.Brush off loose particles from skin. Immese in cool water/wrap in wet bandages.	P335+P334
.Call a POISON CENTER or doctor/physician if you feel unwell.	P312
.Call a POISON CENTER or doctor/physician.	P311
.Collect spillage.	P391
.Contaminated work clothing should not be allowed out of the workplace.	P272
.Dispose of contents/container to…	P501
.Do no eat	P270
.Do not allow contact with air.	P222
.Do not breathe dust/fume/gas/mist/vapours/spray.	P260
.Do not expose ot temperatures exceeding 50°C/ 122°F.	P412
.DO NOT fight fire when fire reaches explosives.	P373
.Do not get in eyes	P262
.Do not handle until all safety precautions have been read and understood.	P202
.Do NOT induce vomiting.	P331
.Do not spray on an open flame or other ignition source.	P211
.Do not subject to grinding/shock/…/friction.	P250
.Eliminate all ignition sources if safe to do so.	P381
.Evacuate area.	P380
.Explosion risk in case of fire.	P372
.Fight fire remotely due to the risk of explosion.	P375
.Fight fire with normal precautions from a reasonable distance	P374
.Gently wash with plenty of soap and water.	P350
.Get immediate medical advice/attention.	P315
.Get medical advice/attention if you feel unwell.	P314
.Get medical advice/attention.	P313
.Ground/bond container and receiving equipment.	P240
.Handle under inert gas.	P231
.Handle under inert gas. Protect from moisture.	P231+P232
.If breathing is difficult	P341
.If experiencing respiratory symptoms:	P342
.If experiencing respiratory symptoms: Call a POISON CENTER or doctor/physician.	P342+P311
.IF exposed or concerned:	P308
.IF exposed or concerned: Get medical advice/attention.	P308+P313
.IF exposed or if you feel unwell:	P309
.IF exposed or if you feel unwell: Call a POISON CENTER or doctor/physician.	P309+P311
.IF exposed:	P307
.IF exposed: Call a POISON CENTER or doctor/physician.	P307+P311
.If eye irritation persists:	P337
.If eye irritation persists: Get medical advice/attention.	P337+P313
.IF IN EYES:	P305
.IF IN EYES: Rinse cautiously with water for several minuts. Remove contact lenses	P305+P351+P338
.IF INHALED:	P304
.IF INHALED: Call a POISON CENTER or doctor/physician if you feel unwell.	P304+P312
.IF INHALED: If breathing is difficult	P304+P341
.IF INHALED: Remove to fresh air and keep at rest in a position comfortable for breathing.	P304+P340
.If medical advice is needed	P101
.IF ON CLOTHING:	P306
.IF ON CLOTHING: rinse immediately contaminated clothing and skin with plenty of water before removing clothes.	P306+P360
.IF ON SKIN (or hair): Remove/Take off immediately all contaminated clothing. Rinse skin with water/shower.	P303+P361+P353
.IF ON SKIN:	P302
.IF ON SKIN:	P303
.IF ON SKIN: Gently wash with plenty of soap and water.	P302+P350
.IF ON SKIN: Immerse in cool water/wrap in wet bandages.	P302+P334
.IF ON SKIN: Wash with plenty of soap and water.	P302+P352
.If skin irritation occurs:	P332
.If skin irritation occurs: Get medical advice/attention.	P332+P313
.If skin irritation or rash occurs:	P333
.If skin irritation or rash occurs: Get medical advice/attention.	P333+P313
.IF SWALLOWED:	P301
.IF SWALLOWED: Call a POISON CENTER or doctor/physician if you feel unwell.	P301+P312
.IF SWALLOWED: Immediately call a POISON CENTER or doctor/physician.	P301+P310
.IF SWALLOWED: rinse mouth. Do NOT induce vomiting.	P301+P330+P331
.Immediately call a POISON CENTER or doctor/physician.	P310
.Immerse in cool water/wrap in wet bandages.	P334
.In case of fire:	P370
.In case of fire: Evacuate area.	P370+P380
.In case of fire: Evacuate area. Fight fire remotely due to the risk of explosion.	P370+P380+P375
.In case of fire: Stop leak if safe to do so.	P370+P376
.In case of fire: Use… for extinction.	P370+P378
.In case of inadequate ventilation wear respiratory protection.	P285
.In case of major fire and large quantities:	P371
.In case of major fire and large quantities: Evacuate area. Fight fire remotely due to the risk of explosion.	P371+P380+P375
.Keep away from any possible contact with water	P223
.Keep away from heat/sparks/open flames/hot surfaces. – No smoking.	P210
.Keep container tightly closed.	P233
.Keep cool.	P235
.Keep cool. Protect from sunlight.	P235+P410
.Keep only in original container.	P234
.Keep out of reach of children.	P102
.Keep reduction valves free from grease and oil.	P244
.Keep wetted with…	P230
.Keep/Store away from clothing/…/combustible materials.	P220
.Leaking gas fire - Do not extinguish	P377
.Maintain air gap between stacks/pallets.	P407
.Obtain special instructions before use.	P201
.Pressurized container: Do not pierce or burn	P251
.Protect from moisture.	P232
.Protect from sunlight.	P410
.Protect from sunlight. Do no expose to temperatures exceeding 50°C/ 122°F.	P410+P412
.Protect from sunlight. Store in a well-ventilated place.	P410+P403
.Read label before use.	P103
.Remove contact lenses	P338
.Remove to fresh air and keep at rest in a position comfortable for breathing.	P340
.Remove/Take off immediately all contaminated clothing.	P361
.Rinse cautiously with water for several minutes.	P351
.Rinse immediately contaminated clothing and skin with plenty of water before removing clothes.	P360
.Rinse mouth.	P330
.Rinse skin with water/shower.	P353
.Specific measures (see … on this label).	P322
.Specific treatment (see … on this label).	P321
.Specific treatment is urgent (see… on this label).	P320
.Stop leak if safe to do so.	P376
.Store at temperatures not exceeding…°C/…°F.	P411
.Store at temperatures not exceeding…°C/…°F. Keep cool.	P411+P235
.Store aways from other materials.	P420
.Store bulk masses greater than … kg/… lbs at temperatures not exceeding …°C/…°F.	P413
.Store contents under …	P422
.Store in a closed container.	P404
.Store in a dry place.	P402
.Store in a dry place. Store in a closed container.	P402+P404
.Store in a well-ventilated place.	P403
.Store in a well-ventilated place. Keep container tightly closed.	P403+P233
.Store in a well-ventilated place. Keep cool.	P403+P235
.Store in corrosive resistant/… container with a resistant inner liner.	P406
.Store locked up.	P405
.Store…	P401
.Take any precaution to avoid mixing with combustibles…	P221
.Take off contaminated clothing and wash before reuse.	P362
.Take precautionary measures against static discharge.	P243
.Thaw frosted parts with lukewarm water. Do no rub affected area.	P336
.Use explosion-proof electrical/ventilating/lighting/…/ equipment.	P241
.Use only non-sparking tools.	P242
.Use only outdoors or in a well-ventilated area.	P271
.Use personal protective equipment as required.	P281
.Use… for extinction.	P378
.Wash … thoroughly after handling.	P264
.Wash contaminated clothing before reuse.	P363
.Wash with plenty of soap and water.	P352
.Wear cold insulating gloves/face shield/eye protection.	P282
.Wear fire/flame resistant/retardant clothing.	P283
.Wear protective gloves/protective clothing/eye protection/face protection.	P280
.Wear respiratory protection.	P284
`

const HAZARDSTATEMENT = `
.Can become flammable in use.	EUH209A
.Can become highly flammable in use or can become flammable in use.	EUH209
.Can become highly flammable in use.	EUH30
.Catches fire spontaneously if exposed to air.	H250
.Causes damage to organs <or state all organs affected	H370
.Causes damage to organs <or state all organs affected	H372
.Causes serious eye damage.	H318
.Causes serious eye irritation.	H319
.Causes severe skin burns and eye damage.	H314
.Causes skin irritation.	H315
.Contact with acids liberates toxic gas.	EUH031
.Contact with acids liberates very toxic gas.	EUH032
.Contact with water liberates toxic gas.	EUH029
.Contains (name of sensitising substance). May produce an allergic reaction.	EUH208
.Contains chromium (VI). May produce an allergic reaction.	EUH203
.Contains epoxy constituents. See information supplied by the manufacturer.	EUH205
.Contains gas under pressure; may explode if heated.	H280
.Contains isocyanates. See information supplied by the manufacturer.	EUH204
.Contains lead. Should not be used on surfaces liable to be chewed or sucked by children.Warning! Contains lead.	EUH201
.Contains refrigerated gas; may cause cryogenic burns or injury.	H281
.Corrosive to the respiratory tract.	EUH071
.Cyanoacrylate. Danger. Bonds skin and eyes in seconds. Keep out of the reach of children.	EUH202
.Explosive	H202
.Explosive when dry	EUH001
.Explosive with or without contact with air.	EUH006
.Explosive; fire	H203
.Explosive; mass explosion hazard.	H201
.Extremely flammable aerosol.	H222
.Extremely flammable gas.	H220
.Extremely flammable liquid and vapour.	H224
.Fatal if inhaled.	H330
.Fatal if swallowed.	H300
.Fatal in contact with skin.	H310
.Fire or projection hazard.	H204
.Flammable aerosol.	H223
.Flammable gas.	H221
.Flammable liquid and vapour.	H226
.Flammable solid.	H228
.Harmful if inhaled.	H332
.Harmful if swallowed.	H302
.Harmful in contact with skin.	H312
.Harmful to aquatic life with long lasting effects.	H412
.Hazardous to the ozone layer.	EUH059
.Heating may cause a fire or explosion.	H241
.Heating may cause a fire.	H242
.Heating may cause an explosion.	H240
.Highly flammable liquid and vapour.	H225
.In contact with water releases flammable gas.	H261
.In contact with water releases flammable gases which may ignite spontaneously.	H260
.In use may form flammable/explosive vapour-air mixture.	EUH018
.May be corrosive to metals.	H290
.May be fatal if swallowed and enters airways.	H304
.May cause allergy or asthma symptoms or breathing difficulties if inhaled.	H334
.May cause an allergic skin reaction.	H317
.May cause cancer <state route of exposure if it is conclusively proven that no other routes of exposure cause the hazard>.	H350
.May cause cancer by inhalation.	H350i
.May cause damage to organs <or state all organs affected	H371
.May cause damage to organs <or state all organs affected	H373
.May cause drowsiness or dizziness.	H336
.May cause fire or explosion; strong oxidizer.	H271
.May cause genetic defects <state route of exposure if it is conclusively proven that no other routes of exposure cause the hazard>.	H340
.May cause harm to breast-fed children.	H362
.May cause long lasting harmful effects to aquatic life.	H413
.May cause or intensify fire; oxidizer.	H270
.May cause respiratory irritation.	H335
.May damage fertility or the unborn child <state specific effect if known > <state route of exposure if it is conclusively proven that no other routes of exposure cause the hazard>.	H360
.May damage fertility.	H360F
.May damage fertility. May damage the unborn child.	H360FD
.May damage fertility. Suspected of damaging the unborn child.	H360Fd
.May damage the unborn child.	H360D
.May damage the unborn child. Suspected of damaging fertility.	H360Df
.May form explosive peroxides.	EUH019
.May intensify fire; oxidizer.	H272
.May mass explode in fire.	H205
.Reacts violently with water.	EUH014
.Repeated exposure may cause skin dryness or cracking.	EUH066
.Risk of explosion if heated under confinement.	EUH044
.Safety data sheet available on request	EUH210
.Self-heating in large quantities; may catch fire.	H252
.Self-heating: may catch fire.	H251
.Suspected of causing cancer <state route of exposure if it is conclusively proven that no other routs of exposure cause the hazard>.	H351
.Suspected of causing genetic defects <state route of exposure if it is conclusively proven that no other routes of exposure cause the hazard>.	H341
.Suspected of damaging fertility or the unborn child <state specific effect if known> <state route of exposure if it is conclusively proven that no other routes of exposure cause the hazard>.	H361
.Suspected of damaging fertility.	H361f
.Suspected of damaging fertility. Suspected of damaging the unborn child.	H361fd
.Suspected of damaging the unborn child.	H361d
.To avoid risks to human health and the environment	EUH401
.Toxic by eye contact	EUH070
.Toxic if inhaled.	H331
.Toxic if swallowed.	H301
.Toxic in contact with skin.	H311
.Toxic to aquatic life with long lasting effects.	H411
.Unstable explosives.	H200
.Very toxic to aquatic life with long lasting effects.	H410
.Very toxic to aquatic life.	H400
.Warning! Contains cadmium. Dangerous fumes are formed during use. See informationsupplied by the manufacturer. Comply with the safety instructions. Contains (name of sensitising substance). May produce an allergic reaction	EUH207
.Warning! Contains lead.	EUH201A
.Warning! Do not use together with other products. May release dangerous gases (chlorine).	EUH206
`
