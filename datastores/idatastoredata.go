package datastores

// TAG is the list of chemical products tags
// They are inserted into the database during its creation.
const TAG = `
3D Cell Culture
Acid
Antibody
Anticancer Drug
Antigen
Apoptose Inhibitor
Base
Buffer
Cell Counting
Cell Culture
Cell Culture Medium
Cell Staining
Cell Seeding
Cloning
Coating
Cryopreservation
Cryostorage
Dissociation Enzyme
Electrophoresis
ELISA
Enzyme
Flow Cytometry
Fluorescence
Fluorescent Probe
Growth Factor
IgG
IgG1
IgG2
IgG2a
IgM
Immucytochemistry
Immunofluorescence
iPS (Induced Pluripotent Stem Cell)
Isotypic Control
Messenchymal Stem Cell
Microscopy
Molecular Biology
Monoclonal
Nuclear Staining
PARP Inhibitor
PCR
Polyclonal
Primary Antibody
Protein
Recombinant Antibody
Recombinant Protein
Secondary Antibody
Sequencing
Solvent
Small Molecule
Spheroid
Stem Cell
Viability Test
Western Blot
`

const CLASSOFCOMPOUND = `
amine
haloalkane
carboxylic acid
indicateur coloré
solvant_organique
sel
reactif_divers
matiere_active_poudre
matiere_active_solution
aldehyde
alkane
alcohol
amino acid
amide
oxime
ester
alcyne
diazo
keton
lanthanide
chiral
halogenoformiate
alcene
phosphorous coumpound
silane
silyl
carbodiimide
peptidic coupling reagent
ether
palladium complex
nitrile
acyl halide
lactone
Colorant
Aralkyl Amine
hydrazine
carbamate
imine
phenol
Nitro
N-hydroxylamine
iodo
thiol
isocyanate
pyridine
quinone
chloro
phosphine
N-oxyde
siloxane
imidazole
aromatic
boronic acid
nitroso
triazolyl
sucre
bromo
iron complex
thiazol
nitrite
pyridazine
fluoro
enone
aniline
boronic ester
thioether
oxazoline
morpholine
sulfonic acid
thiocyanate
anhydride
aluminium
nitrate
sulfate
acetal
mercapto
acyl bromide
epoxyde
chloroformate
cyanide
benzopyran
furan
halogenated aromatic
ammonium salt
phosphonium salt
catalyst
platinium complex
piperazine
azo
metal
cetone
sulfone
acide
radical
indole
cerium
phosphate
borate
carboxylate salt
oxyde métallique
porphyrin
organomagnesien
lanthanum
amidine
urée
hydrure
metallocene
bore
protective group
titanium complex
polymer
hydrate
sulfonate
lactam
nucleotide
enzyme
hydrazide
manganese complex
peracid
perchlorate
peroxide
thiocarbamoyl
organotin
carbonate
sulfide
alloy
azide
pyrocarbonate
sodium
tungsten
imide
molybdenum complex
organolithium
Rhenium
thiophene
acridone
fluorene
protein
tetrazol
ruthenium
triazine
phosphie
phosphite
silver salt
cyanuric acid
diazomethane precursor
acrylate
benzophenone
phosphoramidite
phosphonate
diacid
iridium complex
iridium
oxirane
tosylate
strontium
hydroxide
piperidine
zirconium
alkoxyde
alkoxide
vanadium
acetophenone
phosphine oxide
phosphoester
sulfonimide
sulfinate
buffer
anthracene
germanium
lipase
tetrazole
quinine
quinidine
complex
allyl
tin
carbohydrate
platinum complex
chlorydrate
phtalimide
deuterated solvent
gadolinium complex
peptides
carbazole
arene
thiocarbonyl
organometallic
diol
copper salt
silazane
phenanthrolin
Bromide
ionic liquid
magnesium complex
benzothiazole
antibiotic
tensio-actifs
Zinc
pyrimidine
resin
sel métallique
gold salt
nitride
benzaldehyde
pyrazole
benzoic acid
pyridinim salt
`

// CATEGORY is the list of chemical products categories
// They are inserted into the database during its creation.
const CATEGORY = `
Antibody
Cell Culture Medium & supplement
Drug
Cellular Growth Factor
Matrix / Coating
Comercial Kit
Fluorescent Probe
Cellular Viability Reagent
Cleaning Product
Maintenance, Calibration Reagent
`

// SUPPLIER is the list of chemical products suppliers
// They are inserted into the database during its creation.
const SUPPLIER = `
Abcam
Acros Organics
Ajinomoto
AMSBio
Biological Industries
Bio-Rad
Biosolve
Calbiochem
Carbosynth
CliniSciences
Corning
Dutcher
Enzo Lifesciences
Eurobio
GE Healthcare
Gibco
Graeger
Interchim
Interscience
Invitrogen
Life Technologies
MC2
Merck / Sigma
Mettler Toledo
Millipore 
Molecular Probes
MP Biomedicals
Novo Nordisk
OriGene
OzBiosciences
Panpharma
PanReach AppliChem
Peprotech
Promega
R&D Systems
UGAP
Santa Cruz
Sarstedt
Selleckchem
StemCell
StemPro
ThermoFisher
TheWell Bioscience
Trevigen
VWR
`

// PRODUCER is the list of chemical products producers
// They are inserted into the database during its creation.
const PRODUCER = `
Abcam
Acros Organics
Ajinomoto
Biological Industries
Bio-Rad
Biosolve
Calbiochem
Carbosynth
Corning
Dutcher
Enzo Lifesciences
Eurobio
GE Healthcare
Gibco
Graeger
Interscience
Invitrogen
Life Technologies
Mettler Toledo
Millipore
Molecular Probes
MP Biomedicals
Novo Nordisk
OriGene
OzBiosciences
Panpharma
PanReach AppliChem
Peprotech
Promega
R&D Systems
Santa Cruz
Selleckchem
Merck / Sigma
StemCell
StemPro
ThermoFisher
TheWell Bioscience
Trevigen
`

// PRECAUTIONARYSTATEMENT is the list of chemical products precautionary statements
// They are inserted into the database during its creation.
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

// HAZARDSTATEMENT is the list of chemical products hazard statements
// They are inserted into the database during its creation.
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
.May cause cancer <state route of exposure if it is conclusively proven that no other routes of exposure cause the hazard>.	H350	C1
.May cause cancer by inhalation.	H350i	C1
.May cause damage to organs <or state all organs affected	H371	
.May cause damage to organs <or state all organs affected	H373	
.May cause drowsiness or dizziness.	H336	
.May cause fire or explosion; strong oxidizer.	H271	
.May cause genetic defects <state route of exposure if it is conclusively proven that no other routes of exposure cause the hazard>.	H340	M1
.May cause harm to breast-fed children.	H362	L
.May cause long lasting harmful effects to aquatic life.	H413	
.May cause or intensify fire; oxidizer.	H270	
.May cause respiratory irritation.	H335	
.May damage fertility or the unborn child <state specific effect if known > <state route of exposure if it is conclusively proven that no other routes of exposure cause the hazard>.	H360	R1
.May damage fertility.	H360F	R1
.May damage fertility. May damage the unborn child.	H360FD	R1
.May damage fertility. Suspected of damaging the unborn child.	H360Fd	R1
.May damage the unborn child.	H360D	R1
.May damage the unborn child. Suspected of damaging fertility.	H360Df	R1
.May form explosive peroxides.	EUH019	
.May intensify fire; oxidizer.	H272	
.May mass explode in fire.	H205	
.Reacts violently with water.	EUH014	
.Repeated exposure may cause skin dryness or cracking.	EUH066	
.Risk of explosion if heated under confinement.	EUH044	
.Safety data sheet available on request	EUH210	
.Self-heating in large quantities; may catch fire.	H252	
.Self-heating: may catch fire.	H251	
.Suspected of causing cancer <state route of exposure if it is conclusively proven that no other routs of exposure cause the hazard>.	H351	C2
.Suspected of causing genetic defects <state route of exposure if it is conclusively proven that no other routes of exposure cause the hazard>.	H341	M2
.Suspected of damaging fertility or the unborn child <state specific effect if known> <state route of exposure if it is conclusively proven that no other routes of exposure cause the hazard>.	H361	R2
.Suspected of damaging fertility.	H361f	R2
.Suspected of damaging fertility. Suspected of damaging the unborn child.	H361fd	R2
.Suspected of damaging the unborn child.	H361d	R2
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

// CMR_CAS is a list of product CAS numbers that are CMRs
// They are inserted into the database during its creation.
const CMR_CAS = `
100-00-5,C2 M2
100-40-3,C2
100-42-5,R2
100-44-7,C1B
100-63-0,C1B M2
10028-18-9,C1A M2 R1B
10039-54-0,C2
10043-35-3,R1B
10046-00-1,C2
100683-97-4,C1B
100683-98-5,C1B
100683-99-6,C1B
100684-02-4,C1B
100684-03-5,C1B
100684-04-6,C1B
100684-05-7,C1B
100684-33-1,C1B
100684-37-5,C1B
100684-38-6,C1B
100684-49-9,C1B
100684-51-3,C1B
100801-63-6,C1B M1B
100801-65-8,C1B M1B
100801-66-9,C1B M1B
100988-63-4,M2
101-14-4,C1B
101-21-3,C2
101-61-1,C1B
101-68-8,C2
101-77-9,C1B M2
101-80-4,C1B M1B R2
101-90-6,C2 M2
10101-96-9,C1A
10108-64-2,C1B M1B R1B
101205-02-1,R2
10124-36-4,C1B M1B R1B
10124-43-3,C1B M2 R1B
101316-45-4,C1B
101316-49-8,C1B
101316-56-7,C1B M1B
101316-57-8,C1B
101316-59-0,C1B
101316-62-5,C1B M1B
101316-63-6,C1B M1B
101316-66-9,C1B M1B
101316-67-0,C1B M1B
101316-69-2,C1B
101316-70-5,C1B
101316-71-6,C1B
101316-72-7,C1B
101316-76-1,C1B M1B
101316-83-0,C1A
101316-84-1,C1A
101316-85-2,C1B
101316-86-3,C1B M1B
101316-87-4,C1B M1B
10141-05-6,C1B M2 R1B
101463-69-8,LACT
101631-14-5,C1B
101631-20-3,C1B M1B
101794-74-5,C1B
101794-75-6,C1B
101794-76-7,C1B
101794-90-5,C1B M1B
101794-91-6,C1B M1B
101794-97-2,C1B M1B
101795-01-1,C1B M1B
101896-26-8,C1B M1B
101896-27-9,C1B M1B
101896-28-0,C1B M1B
102-06-7,R2
102110-14-5,C1B M1B
102110-15-6,C1B M1B
102110-55-4,C1B M1B
1024-57-3,C2
103-33-3,C1B M2
103112-35-2,C1B
103122-66-3,C1B M1B
10325-94-7,C1B M1B
10332-33-9,R1B
103361-09-7,R1B
10381-36-9,C1A
104-91-6,M2
104653-34-1,R1B
10486-00-7,R1B
105024-66-6,R1B
10588-01-9,C1B M1B R1B
106-46-7,C2
106-47-8,C1B
106-49-0,C2
106-87-6,C2
106-88-7,C2
106-89-8,C1B
106-91-2,C1B M2 R1B
106-92-3,C2 M2 R2
106-93-4,C1B
106-94-5,R1B
106-97-8,C1A M1B
106-99-0,C1A M1B
10605-21-7,M1B R1B
107-05-1,C2 M2
107-06-2,C1B
107-13-1,C1B
107-20-0,C2
107-22-2,M2
107-30-2,C1A
107534-96-3,R2
108-05-4,C2
108-45-2,M2
108-88-3,R2
108-91-8,R2
108-95-2,M2
108225-03-2,C1B
109-86-4,R1B
109-99-9,C2
110-00-9,C1B M2
110-05-4,M2
110-49-6,R1B
110-54-3,R2
110-71-4,R1B
110-80-5,R1B
110-85-0,R2
110-88-3,R2
110235-47-7,C2
11099-02-8,C1A
111-15-9,R1B
111-41-1,R1B
111-44-4,C2
111-77-3,R2
111-96-6,R1B
11113-50-1,R1B
11113-74-9,C1A M2 R1B
11113-75-0,C1A M2
11132-10-8,C1A M2 R1B
11138-47-9,R1B
1116-54-7,C1B
111988-49-9,C2 R1B
112-49-2,R1B
1120-71-4,C1B
114565-66-1,C2
115-96-8,C2 R1B
115662-06-1,R2
117-81-7,R1B
117-82-8,R1B
118-74-1,C1B
118134-30-8,R2
118612-00-3,C2
118658-99-4,C1B
119-90-4,C1B
119-93-7,C1B
119738-06-6,M2 R1B
120-32-1,C2 R2
120-71-8,C1B
12001-28-4,C1A
12001-29-5,C1A
12004-35-2,C1A
12007-00-0,C1A
12007-01-1,C1A
12007-02-2,C1A
12008-41-2,R1B
120187-29-3,M2
12031-65-1,C1A
12035-36-8,C1A
12035-38-0,C1A
12035-39-1,C1A
12035-64-2,C1A
12035-71-1,C1A M2
12035-72-2,C1A M2
12040-72-1,R1B
12054-48-7,C1A M2 R1B
12056-51-8,C2
12059-14-2,C1A
12068-61-0,C1A
121-14-2,C1B M2 R2
121-69-7,C2
121158-58-5,R1B
12137-12-1,C1A
12142-88-0,C1A
121575-60-8,C1B
121620-46-0,C1B M1B
121620-47-1,C1B M1B
121620-48-2,C1B M1B
12172-73-5,C1A
12179-04-3,R1B
122-34-9,C2
122-60-1,C1B M2
122-66-7,C1B
12201-89-7,C1A
122070-78-4,C1B
122070-79-5,C1B M1B
122070-80-8,C1B M1B
122384-77-4,C1B
122384-78-5,C1B M1B
12267-73-1,R1B
12280-03-4,R1B
123-30-8,M2
123-31-9,C2 M2
123-39-7,R1B
123-73-9,M2
123-91-1,C2
123312-89-0,C2
1239-45-8,M2
12427-38-2,R2
125051-32-3,R2
12510-42-8,C1A
125116-23-6,R2
12519-85-6,C1A
126-73-8,C2
126-99-8,C1B
12607-70-4,C1A M2 R1B
12619-90-8,C1A
12653-76-8,C1A
126535-15-7,C2
12656-85-8,C1B R1A
12673-58-4,C1A
127-18-4,C2
127-19-5,R1B
12737-30-3,C1A
129-73-7,C2 M2
1303-00-0,C1B R1B
1303-28-2,C1A
1303-86-2,R1B
1303-96-4,R1B
130328-20-0,R2
1304-56-9,C1B
1306-19-0,C1B M2 R2
1306-23-6,C1B M2 R2
130728-76-6,M2
1309-64-4,C2
131-18-0,R1B
131-52-2,C2
1313-27-5,C2
1313-99-1,C1A
13138-45-9,C1A M2 R1B
1314-04-1,C1A M2
1314-05-2,C1A
1314-06-3,C1A
1314-62-1,M2 R2
13171-21-6,M2
132-32-1,C1B
132207-32-0,C1A
1327-53-3,C1A
133-06-2,C2
133-07-3,C2
1330-43-4,R1B
1333-82-0,C1A M1B R2
1335-32-6,C2 R1A
13360-57-1,C1B
133855-98-8,C2 R1B
13424-46-9,R1A
1344-37-2,C1B R1A
13462-88-9,C1A M2 R1B
13462-90-3,C1A M2 R1B
13463-39-3,C2 R1B
13465-08-2,C2
13477-70-8,C1A
135-88-6,C2
13517-20-9,R1B
13595-25-0,R2
13637-71-3,C1A M2 R1B
13654-40-5,C1A M2 R1B
13674-87-8,C2
13689-92-4,C1A M2 R1B
137-17-7,C1B
13765-19-0,C1B
13770-89-3,C1A M2 R1B
13775-54-7,C1A
138164-12-2,R2
13840-56-7,R1B
13842-46-1,C1A M2 R1B
138526-69-9,C2
139-40-2,C2
139-65-1,C1B
139001-49-3,C2 R2
139528-85-1,C2
140-41-0,C2
140698-96-0,C2
141112-29-0,R2
14177-51-6,C1A
14177-55-0,C1A
142-64-3,R2
1420-07-1,R1B
14216-75-2,C1A M2 R1B
142891-20-1,C2
143-50-0,C2
14332-34-4,C1A
143322-57-0,R2
143390-89-0,C2
143860-04-2,R1B
144177-62-8,R2
14448-18-1,C1A
14507-36-9,C1A
14550-87-9,C1A M2 R1B
1464-53-5,C1B M1B
14708-14-6,C1A M2 R1B
14721-18-7,C1A
148-24-3,R1B
14816-18-3,R2
1484-13-5,M2
14874-78-3,C1A
149-57-5,R2
149591-38-8,R2
14977-61-8,C1B M1B
149961-52-4,C2 R2
149979-41-9,C2 R2
14998-37-9,C1A M2 R1B
150-68-5,C2
15060-62-5,C1A M2 R1B
151-56-4,C1B M1B
15120-21-5,R1B
15159-40-7,C2
151798-26-4,R1B
151882-81-4,C2
15245-44-0,R1A
15375-21-0,R2
15545-48-9,C2 R2
15571-58-1,R1B
15586-38-6,C1A M2 R1B
156-43-4,M2
15606-95-8,C1A
156145-66-3,R2
15699-18-0,C1A M2 R1B
15780-33-3,C1A
1582-09-8,C2
15843-02-4,C1A M2 R1B
15851-52-2,C1A
15852-21-8,C1A
158894-67-8,C2
1589-47-5,R1B
1593-77-7,R2
15972-60-8,C2
159939-85-2,R2
16039-61-5,C1A M2 R1B
16071-86-6,C1B
16083-14-0,C1A M2 R1B
16118-49-3,C2 R1B
16337-84-1,C1A M2 R1B
163879-69-4,R2
164058-22-4,C1B
166242-53-1,C2
1671-49-4,R2
16812-54-7,C1A M2
1689-83-4,R2
1689-84-5,R2
1689-99-2,R2
1694-09-3,C2
17010-21-8,C2
17169-61-8,C1A
17570-76-2,R1A
1763-23-1,C2 R1B LACT
17630-75-0,R2
17804-35-2,M1B R1B
18283-82-4,C1A M2 R1B
183196-57-8,R1B
1836-75-5,C1B R1B
18718-11-1,C1A
18721-51-2,C1A M2 R1B
189278-12-4,C2
1897-45-6,C2
19098-16-9,C2
192-97-2,C1B
1937-37-7,C1B R2
19372-20-4,C1A
19398-06-2,C2
194-55-6,R2 LACT
1951-97-9,R2
19750-95-9,C2
19900-65-3,C2
199327-61-2,R1B
20108-78-5,R2
202197-26-0,M2
203313-25-1,R2
2040-90-6,M1B R2
205-82-3,C1B
205-99-2,C1B
20543-06-0,C1A
207-08-9,C1B
20845-01-6,C2
21041-95-2,C1B M1B
21049-39-8,C2 R1B LACT
210555-94-5,R1B
21136-70-9,C1A
2122-19-2,R2
214353-17-0,C1B
21436-97-5,C1B
2164-08-1,C2
21784-78-1,C1A
218-01-9,C1B M2
2186-24-5,M2
2186-25-6,M2
220444-73-5,C2
2210-79-9,M2
2212-67-1,C2 R2
221354-37-6,R2
2223-95-2,C1A M2 R1B
22398-80-7,C1B R2
2243-62-1,C2
22605-92-1,C1A M2 R1B
2303-16-4,C2
23085-60-1,R2
23103-98-2,C2
2312-35-8,C2
2314-97-8,M2
23564-05-8,M2
2385-85-5,C2 R2 LACT
23950-58-5,C2
2425-06-1,C1B
2426-08-6,C2 M2
2431-50-7,C2
2437-29-8,R2
2439-01-2,R2
2451-62-9,M1B
24602-86-6,R1B
24613-89-6,C1B
2475-45-8,C1B
25154-52-3,R2
25155-23-1,R1B
25321-14-6,C1B M2 R2
2536-05-2,C2
25383-07-7,R2
25637-99-4,R2 LACT
25808-74-6,R1A
2593-15-9,C2
2602-46-2,C1B R2
26043-11-8,C1A M2 R1B
26157-73-3,M2
26447-14-3,M2
26447-40-5,C2
26471-62-5,C2
2687-91-4,R1B
27016-75-7,C1A
27083-27-8,C2
27140-08-5,C1B M2
27366-72-9,R1B
27610-48-6,M2
27637-46-3,C1A M2 R1B
2795-39-3,C2 R1B LACT
2832-40-8,C2
28772-56-7,R1B
288-32-4,R1B
288-88-0,R2
29081-56-9,C2 R1B LACT
29317-63-3,C1A M2 R1B
29457-72-5,C2 R1B LACT
301-04-2,R1A
302-01-2,C1B
302-97-6,R2
3033-77-0,C1B M2 R2
309-00-2,C2
3108-42-7,C2 R1B LACT
3165-93-3,C1B M2
31717-87-0,R2
31748-25-1,C1A
32289-58-0,C2
32534-81-9,LACT
32536-52-0,R1B
330-54-1,C2
330-55-2,C2 R1B
3327-22-8,C2
3333-67-3,C1A M2 R1B
334-88-3,C1B
3349-06-2,C1A M2 R1B
3349-08-4,C1A M2 R1B
335-67-1,C2 R1B LACT
335-76-2,C2 R1B LACT
335104-84-2,R2
34123-59-6,C2
34256-82-1,C2 R2
34492-97-2,C1A
35554-44-0,C2
36026-88-7,C1A
36341-27-2,C1A
36734-19-7,C2
3691-35-8,R1B
3724-43-4,R1B
37244-98-7,R1B
373-02-4,C1A M2 R1B
37321-15-6,C1A
375-95-1,C2 R1B LACT
37894-46-5,R1B
3825-26-1,C2 R1B LACT
3830-45-3,C2 R1B LACT
3861-47-0,R2
3878-19-1,C2
3906-55-6,C1A M2 R1B
39156-41-7,C1B M2
39300-45-3,R1B
39807-15-3,R2
39819-65-3,C1A M2 R1B
399-95-1,C1B
40722-80-3,C1B M1B
41107-56-6,M2
41483-43-6,C2
4149-60-4,C2 R1B LACT
4170-30-3,M2
420-04-2,C2 R2
4454-16-4,C1A M2 R1B
4464-23-7,C2
485-31-4,R1B
492-80-8,C2
495-54-5,M2
4995-91-9,C1A M2 R1B
50-00-0,C1B M2
50-29-3,C2
50-32-8,C1B M1B R1B
50471-44-8,C2 R1B
5064-31-3,C2
51-79-6,C1B
51085-52-0,C2
51229-78-8,R2
513-78-0,C1B M1B
513-79-1,C1B M2 R1B
5146-66-7,M1B
51594-55-9,C1B
51818-56-5,C1A M2 R1B
52033-74-6,C1B M2
5216-25-1,C1B R2
52234-82-9,M2
52502-12-2,C1A
52625-25-9,C1A M2 R1B
53-70-3,C1B
531-85-1,C1A
531-86-2,C1A
532-82-1,M2
534-52-1,M2
53933-48-5,C2
540-23-8,C2
540-25-0,C2
540-73-8,C1B
5406-86-0,R2
541-69-5,M2
542-56-3,C1B M2
542-83-6,C2
542-88-1,C1A
547-67-1,C1A
5470-11-1,C2
548-62-9,C1B
55-38-9,M2
55219-65-3,R1B LACT
553-00-4,C1A
553-71-9,C1A M2 R1B
5543-57-7,R1A
5543-58-8,R1A
556-52-5,C1B M2 R1B
556-67-2,R2
557-19-7,C1A
5571-36-8,R1B
56-23-5,C2
56-55-3,C1B
56073-07-5,R1B
56073-10-0,R1A
5625-90-1,C1B M2
56634-95-8,R2
569-61-9,C1B
569-64-2,R2
57-14-7,C1B
57-57-8,C1B
57-74-9,C2
57044-25-4,C1B M2 R1B
573-58-0,C1B R2
57583-34-3,R2
57583-35-4,R2
57966-95-7,R2
58-89-9,LACT
581-89-5,C1B
5836-29-3,R1B
584-84-9,C2
58591-45-0,C1A
5873-54-1,C2
59-88-1,C1B M2
591-78-6,R2
592-62-1,C1B R1B
593-60-2,C1B
59653-74-6,M1B
60-09-3,C1B
60-34-4,C1B
60-35-5,C2
60-57-1,C2
60168-88-9,R2 LACT
602-01-7,C1B M2 R2
602-87-9,C1B
605-50-5,R1B
60568-05-0,C2
606-20-2,C1B M2 R2
6094-40-2,R2
61-82-5,R2
610-39-9,C1B M2 R2
612-52-2,C1A
612-82-8,C1B
613-35-4,C1B M2
615-05-4,C1B M2
615-28-1,C2 M2
61571-06-0,R1B
6164-98-3,C2
61789-28-4,C1B
61789-60-4,C1B
618-85-9,C1B M2 R2
619-15-8,C1B M2 R2
62-53-3,C2 M2
62-55-5,C1B
62-56-6,C2 R2
62-75-9,C1B
621-64-7,C1B
624-83-9,R2
625-45-6,R1B
629-14-1,R1A
63-25-2,C2
630-08-0,R1A
63681-54-9,M2
64-67-5,C1B M1B
64-86-8,M1B
64485-90-1,C2
64741-41-9,C1B M1B
64741-42-0,C1B M1B
64741-45-3,C1B
64741-46-4,C1B M1B
64741-47-5,C1B M1B
64741-48-6,C1B M1B
64741-50-0,C1A
64741-51-1,C1A
64741-52-2,C1A
64741-53-3,C1A
64741-54-4,C1B M1B
64741-55-5,C1B M1B
64741-57-7,C1B
64741-59-9,C1B
64741-60-2,C1B
64741-61-3,C1B
64741-62-4,C1B
64741-63-5,C1B M1B
64741-64-6,C1B M1B
64741-65-7,C1B M1B
64741-66-8,C1B M1B
64741-67-9,C1B
64741-68-0,C1B M1B
64741-69-1,C1B M1B
64741-70-4,C1B M1B
64741-74-8,C1B M1B
64741-75-9,C1B
64741-76-0,C1B
64741-77-1,C2
64741-78-2,C1B M1B
64741-80-6,C1B
64741-81-7,C1B
64741-82-8,C1B
64741-83-9,C1B M1B
64741-84-0,C1B M1B
64741-86-2,C1B
64741-87-3,C1B M1B
64741-88-4,C1B
64741-89-5,C1B
64741-90-8,C1B
64741-91-9,C1B
64741-92-0,C1B M1B
64741-95-3,C1B
64741-96-4,C1B
64741-97-5,C1B
64742-01-4,C1B
64742-03-6,C1B
64742-04-7,C1B
64742-05-8,C1B
64742-11-6,C1B
64742-12-7,C1B
64742-13-8,C1B
64742-14-9,C1B
64742-15-0,C1B M1B
64742-18-3,C1A
64742-19-4,C1A
64742-20-7,C1A
64742-21-8,C1A
64742-22-9,C1B M1B
64742-23-0,C1B M1B
64742-27-4,C1A
64742-28-5,C1A
64742-29-6,C1B
64742-30-9,C1B
64742-34-3,C1A
64742-35-4,C1A
64742-36-5,C1B
64742-37-6,C1B
64742-38-7,C1B
64742-41-2,C1B
64742-44-5,C1B
64742-45-6,C1B
64742-46-7,C1B
64742-48-9,C1B M1B
64742-49-0,C1B M1B
64742-52-5,C1B
64742-53-6,C1B
64742-54-7,C1B
64742-55-8,C1B
64742-56-9,C1B
64742-57-0,C1B
64742-59-2,C1B
64742-61-6,C1B
64742-62-7,C1B
64742-63-8,C1B
64742-64-9,C1B
64742-65-0,C1B
64742-66-1,C1B M1B
64742-67-2,C1B
64742-68-3,C1B
64742-69-4,C1B
64742-70-7,C1B
64742-71-8,C1B
64742-73-0,C1B M1B
64742-75-2,C1B
64742-76-3,C1B
64742-78-5,C1B
64742-79-6,C1B
64742-80-9,C1B
64742-82-1,C1B M1B
64742-83-2,C1B M1B
64742-86-5,C1B
64742-89-8,C1B M1B
64742-90-1,C1B
64742-95-6,C1B M1B
64743-01-7,C1B
64969-36-4,C1B
65195-55-3,R2
65229-23-4,C1A
65277-42-1,R1B
65321-67-7,C1B
65322-65-8,C2 M2
65405-96-1,C1A M2 R1B
65756-41-4,M2
65996-78-3,C1B M1B
65996-79-4,C1B M1B
65996-82-9,C1B M1B
65996-83-0,C1B M1B
65996-84-1,C1B M1B
65996-85-2,C1B M1B
65996-86-3,C1B M1B
65996-87-4,C1B M1B
65996-88-5,C1B M1B
65996-89-6,C1A
65996-90-9,C1A
65996-91-0,C1B
65996-92-1,C1B
65996-93-2,C1A M1B R1B
66-81-9,M2 R1B
66246-88-6,R2
66938-41-8,M2
67-66-3,C2 R2
67129-08-2,C2
67564-91-4,R2
67891-79-6,C1B M1B
67891-80-9,C1B M1B
67952-43-6,C1A M2 R1B
68-12-2,R1B
680-31-9,C1B M1B
68016-03-5,C1A
6804-07-5,C1B
68049-83-2,R1B
6807-17-6,R1B
68130-19-8,C1A R1A
68130-36-9,C1A
68131-49-7,C1B M1B
68131-75-9,C1A M1B
68134-59-8,C1A M2 R1B
68157-60-8,C2
68186-89-0,C1A
68187-57-5,C1B
68188-48-7,C1B
683-18-1,M2 R1B
68307-98-2,C1A M1B
68307-99-3,C1A M1B
68308-00-9,C1A M1B
68308-01-0,C1A M1B
68308-03-2,C1B M1B
68308-04-3,C1A M1B
68308-05-4,C1A M1B
68308-06-5,C1A M1B
68308-07-6,C1A M1B
68308-08-7,C1A M1B
68308-09-8,C1A M1B
68308-10-1,C1A M1B
68308-11-2,C1A M1B
68308-12-3,C1B M1B
68333-22-2,C1B
68333-25-5,C1B
68333-26-6,C1B
68333-27-7,C1B
68333-28-8,C1B
68334-30-5,C2
68391-11-7,C1B M1B
68409-99-4,C1A M1B
68410-05-9,C1B M1B
68410-71-9,C1B M1B
68410-96-8,C1B M1B
68410-97-9,C1B M1B
68410-98-0,C1B M1B
68425-29-6,C1B M1B
68425-35-4,C1B M1B
68475-57-0,C1A M1B
68475-58-1,C1A M1B
68475-59-2,C1A M1B
68475-60-5,C1A M1B
68475-70-7,C1B M1B
68475-79-6,C1B M1B
68475-80-9,C1B
68476-26-6,C1A M1B
68476-29-9,C1A M1B
68476-30-2,C2
68476-31-3,C2
68476-32-4,C1B
68476-33-5,C1B
68476-34-6,C2
68476-40-4,C1A M1B
68476-42-6,C1A M1B
68476-46-0,C1B M1B
68476-47-1,C1B M1B
68476-49-3,C1A M1B
68476-50-6,C1B M1B
68476-55-1,C1B M1B
68476-85-7,C1A M1B
68476-86-8,C1A M1B
68477-23-6,C1B M1B
68477-29-2,C1B
68477-30-5,C1B
68477-31-6,C1B
68477-33-8,C1A M1B
68477-34-9,C1B M1B
68477-35-0,C1A M1B
68477-38-3,C1B
68477-50-9,C1B M1B
68477-53-2,C1B M1B
68477-55-4,C1B M1B
68477-61-2,C1B M1B
68477-65-6,C1A M1B
68477-66-7,C1A M1B
68477-67-8,C1A M1B
68477-68-9,C1A M1B
68477-69-0,C1A M1B
68477-70-3,C1A M1B
68477-71-4,C1A M1B
68477-72-5,C1A M1B
68477-73-6,C1A M1B
68477-74-7,C1A M1B
68477-75-8,C1A M1B
68477-76-9,C1A M1B
68477-77-0,C1A M1B
68477-79-2,C1A M1B
68477-80-5,C1A M1B
68477-81-6,C1A M1B
68477-82-7,C1A M1B
68477-83-8,C1A M1B
68477-84-9,C1A M1B
68477-85-0,C1A M1B
68477-86-1,C1A M1B
68477-87-2,C1A M1B
68477-89-4,C1B M1B
68477-90-7,C1A M1B
68477-91-8,C1A M1B
68477-92-9,C1A M1B
68477-93-0,C1A M1B
68477-94-1,C1A M1B
68477-95-2,C1A M1B
68477-96-3,C1A M1B
68477-97-4,C1A M1B
68477-98-5,C1A M1B
68477-99-6,C1A M1B
68478-00-2,C1A M1B
68478-01-3,C1A M1B
68478-02-4,C1A M1B
68478-03-5,C1A M1B
68478-04-6,C1A M1B
68478-05-7,C1A M1B
68478-12-6,C1B M1B
68478-13-7,C1B
68478-15-9,C1B M1B
68478-16-0,C1B M1B
68478-17-1,C1B
68478-21-7,C1A M1B
68478-22-8,C1A M1B
68478-24-0,C1A M1B
68478-25-1,C1A M1B
68478-26-2,C1A M1B
68478-27-3,C1A M1B
68478-28-4,C1A M1B
68478-29-5,C1A M1B
68478-30-8,C1A M1B
68478-32-0,C1A M1B
68478-33-1,C1A M1B
68478-34-2,C1A M1B
68512-61-8,C1B
68512-62-9,C1B
68512-78-7,C1B M1B
68512-91-4,C1A M1B
68513-02-0,C1B M1B
68513-03-1,C1B M1B
68513-14-4,C1A M1B
68513-15-5,C1A M1B
68513-16-6,C1A M1B
68513-17-7,C1A M1B
68513-18-8,C1A M1B
68513-19-9,C1A M1B
68513-63-3,C1B M1B
68513-66-6,C1A M1B
68513-69-9,C1B
68513-87-1,C1B M1B
68514-15-8,C1B M1B
68514-31-8,C1A M1B
68514-36-3,C1A M1B
68514-79-4,C1B M1B
68515-42-4,R1B
68515-50-4,R1B
68515-84-4,C1A
68516-20-1,C1B M1B
68527-15-1,C1A M1B
68527-16-2,C1A M1B
68527-18-4,C1B
68527-19-5,C1A M1B
68527-21-9,C1B M1B
68527-22-0,C1B M1B
68527-23-1,C1B M1B
68527-26-4,C1B M1B
68527-27-5,C1B M1B
68553-00-4,C1B
68555-24-8,C1B M1B
68602-82-4,C1A M1B
68602-83-5,C1A M1B
68602-84-6,C1A M1B
68603-00-9,C1B M1B
68603-01-0,C1B M1B
68603-03-2,C1B M1B
68603-08-7,C1B M1B
68606-10-0,C1B M1B
68606-11-1,C1B M1B
68606-25-7,C1A M1B
68606-26-8,C1A M1B
68606-27-9,C1A M1B
68606-34-8,C1A M1B
68607-11-4,C1A M1B
68607-30-7,C1B
68610-24-2,C1A
68694-11-1,R1B
68783-00-6,C1B
68783-04-0,C1B
68783-06-2,C1A M1B
68783-07-3,C1A M1B
68783-08-4,C1B
68783-09-5,C1B M1B
68783-12-0,C1B M1B
68783-13-1,C1B
68783-64-2,C1A M1B
68783-65-3,C1A M1B
68783-66-4,C1B M1B
68814-67-5,C1A M1B
68814-89-1,C1B
68814-90-4,C1A M1B
68815-21-4,C1B M1B
68911-58-0,C1A M1B
68911-59-1,C1A M1B
68918-99-0,C1A M1B
68919-00-6,C1A M1B
68919-01-7,C1A M1B
68919-02-8,C1A M1B
68919-03-9,C1A M1B
68919-04-0,C1A M1B
68919-05-1,C1A M1B
68919-06-2,C1A M1B
68919-07-3,C1A M1B
68919-08-4,C1a M1B
68919-09-5,C1A M1B
68919-10-8,C1A M1B
68919-11-9,C1A M1B
68919-12-0,C1A M1B
68919-20-0,C1A M1B
68919-37-9,C1B M1B
68919-39-1,C1B M1B
68921-08-4,C1B M1B
68921-09-5,C1B M1B
68937-63-3,C1B M1B
68952-76-1,C1A M1B
68952-77-2,C1A M1B
68952-79-4,C1A M1B
68952-80-7,C1A M1B
68952-81-8,C1A M1B
68952-82-9,C1A M1B
68955-27-1,C1B
68955-28-2,C1A M1B
68955-29-3,C1B M1B
68955-33-9,C1A M1B
68955-34-0,C1A M1B
68955-35-1,C1B M1B
68955-36-2,C1B
68989-88-8,C1A M1B
68990-61-4,C1B
69012-50-6,C1A
69094-18-4,C2
69227-51-6,M2
6923-22-4,M2
69806-50-4,R1B
70-25-7,C1B
70225-14-8,C2 R1B LACT
70321-67-4,C1B M1B
70321-79-8,C1B
70321-80-1,C1B
70592-76-6,C1B
70592-77-7,C1B
70592-78-8,C1B
70657-70-4,R1B
70692-93-2,C1A
70987-78-9,C1B M2
71-43-2,C1A M1B
71-48-7,C1B M2 R1B
71720-48-4,C1A M2 R1B
71751-41-2,R2
71868-10-5,R1B
71888-89-6,R1B
71957-07-8,C1A M2 R1B
7226-23-5,R2
72319-19-8,C1A M2 R1B
72490-01-8,C2
72623-85-9,C1B
72623-86-0,C1B
72623-87-1,C1B
73665-18-6,C1B M1B
74-83-9,M2
74-87-3,C2
74-88-4,C2
74-96-4,C2
74070-46-5,C2
74195-78-1,C1A
7425-14-1,R2
7439-92-1,R1A LACT
7439-97-6,R1B
7440-02-0,C2
7440-41-7,C1B
7440-43-9,C1B M2 R2
7446-27-7,R1A
74499-35-7,R1B
74646-29-0,C1A
74753-18-7,C1B
74869-21-9,C1B
74869-22-0,C1B
7487-94-7,M2 R2
75-00-3,C2
75-01-4,C1A
75-07-0,C2
75-09-2,C2
75-12-7,R1B
75-15-0,R2
75-21-8,C1B M1B
75-26-3,R1A
75-28-5,C1A M1B
75-35-4,C2
75-55-8,C1B
75-56-9,C1B M1B
75-91-2,M2
75113-37-0,M2 R1B
753-73-1,R2
75660-25-2,M2
7572-29-4,C2
7580-31-6,C1A M2 R1B
75980-60-8,R2
76-01-7,C2
76-44-8,C2
76-87-9,C2 R2
7632-04-4,R1B
764-41-0,C1B
7646-79-9,C1B M2 R1B
77-09-8,C1B M2 R2
77-58-7,M2 R1B
77-78-1,C1B M2
7718-54-9,C1A M2 R1B
77182-82-2,R1B
77402-03-0,C1B M1B
77402-05-2,C1B M1B
77536-66-4,C1A
77536-67-5,C1A
77536-68-6,C1A
7757-95-1,C1A
7758-01-2,C1B
7758-97-6,C1B R1A
7775-11-3,C1B M1B R1B
7778-50-9,C1B M1B R1B
7778-73-6,C2
777891-21-1,R2
7784-40-9,C1A R1A
7786-81-4,C1A M2 R1B
7789-00-6,C1B M1B
7789-06-2,C1B
7789-09-5,C1B M1B R1B
7790-79-6,C1B M1B R1B
7790-80-9,C2
78-59-1,C2
78-79-5,C1B M2
78-87-5,C1B
78-88-6,M2
7803-49-8,C2
79-00-5,C2
79-01-6,C1B M2
79-06-1,C1B M1B R2
79-07-2,R2
79-16-3,R1B
79-44-7,C1B
79-46-9,C1B
79234-33-6,M2
79241-46-6,R2
79622-59-6,R2
79815-20-6,R2
80-05-7,R1B
8001-35-2,C2
8001-58-9,C1B
8002-05-9,C1B
8006-61-9,C1B M1B
8007-45-2,C1A
8009-03-8,C1B
8018-01-7,R2
8030-30-6,C1B M1B
8032-32-4,C1B M1B
80387-97-9,R1B
8052-41-3,C1B M1B
80844-07-1,LACT
81-14-1,C2
81-15-2,C2
81-81-2,R1A
81880-96-8,M2
823-40-5,M2
82413-20-5,C2 R1B
82560-54-1,R2
82657-04-3,C2
83056-32-0,M2
838-88-0,C1B
83968-67-6,M2
84-61-7,R1B
84-65-1,C1B
84-69-5,R1B
84-74-2,R1B
84-75-3,R1B
84196-22-5,M2
842-07-9,C2 M2
84245-12-5,C1B M1B R1B
84332-86-5,C2
84650-02-2,C1B
84650-03-3,C1B M1B
84650-04-4,C1B M1B
84776-45-4,C1A M2 R1B
84777-06-0,R1B
84852-15-3,R2
84852-35-7,C1A M2 R1B
84852-36-8,C1A M2 R1B
84852-37-9,C1A M2 R1B
84852-39-1,C1A M2 R1B
84988-93-2,C1B M1B
84989-03-7,C1B M1B
84989-04-8,C1B M1B
84989-05-9,C1B M1B
84989-06-0,C1B M1B
84989-07-1,C1B M1B
84989-09-3,C1B M1B
84989-10-6,C1B
84989-11-7,C1B
84989-12-8,C1B M1B
85-68-7,R1B
85029-51-2,C1B M1B
85029-74-9,C1B
85116-53-6,C1B
85116-58-1,C1B M1B
85116-59-2,C1B M1B
85116-60-5,C1B M1B
85116-61-6,C1B M1B
85117-03-9,C1B
85135-77-9,C1A M2 R1B
85136-74-9,C1B
85166-19-4,C1A M2 R1B
85407-90-5,M2
85508-43-6,C1A M2 R1B
85508-44-7,C1A M2 R1B
85508-45-8,C1A M2 R1B
85508-46-9,C1A M2 R1B
85509-19-9,C2 R1B
85535-84-8,C2
85535-85-9,LACT
85536-17-0,C1B M1B
85536-19-2,C1B M1B
85536-20-5,C1B M1B
85551-28-6,C1A M2 R1B
85954-11-6,C2
86-88-4,C2
86290-81-5,C1B M1B
86552-32-1,C2
87-62-7,C2
87-66-1,M2
87-86-5,C2
872-50-4,R1B
87691-88-1,R2
87741-01-3,C1A M1B
87820-88-0,C2
88-06-2,C2
88-10-8,C2
88-12-0,C2
88-72-2,C1B M1B R2
88-85-7,R1B
88671-89-0,R2
90-04-0,C1B M2
90-41-5,C2
90-94-8,C1B M2
900-95-8,C2 R2
90035-08-8,R1B
90622-53-0,C1B
90622-55-2,C1A M1B
90640-80-5,C1B
90640-81-6,C1B M1B
90640-82-7,C1B M1B
90640-84-9,C1B
90640-85-0,C1B
90640-86-1,C1B
90640-87-2,C1B M1B
90640-88-3,C1B M1B
90640-89-4,C1B M1B
90640-90-7,C1B M1B
90640-91-8,C1B
90640-92-9,C1B
90640-93-0,C1B
90640-94-1,C1B
90640-95-2,C1B
90640-96-3,C1B
90640-97-4,C1B
90640-99-6,C1B M1B
90641-00-2,C1B M1B
90641-01-3,C1B M1B
90641-02-4,C1B M1B
90641-03-5,C1B M1B
90641-04-6,C1B M1B
90641-05-7,C1B M1B
90641-06-8,C1B M1B
90641-07-9,C1B
90641-08-0,C1B
90641-09-1,C1B
90641-11-5,C1B M1B
90641-12-6,C1B M1B
90657-55-9,R2
90669-57-1,C1B
90669-58-2,C1B
90669-59-3,C1B
90669-74-2,C1B
90669-75-3,C1B
90669-76-4,C1B
90669-77-5,C1B
90669-78-6,C1B
90989-38-1,C1B M1B
90989-39-2,C1B M1B
90989-41-6,C1B M1B
90989-42-7,C1B M1B
91-08-7,C2
91-20-3,C2
91-22-5,C1B M2
91-23-6,C1B
91-59-8,C1A
91-94-1,C1B
91-95-2,C1B M2
91079-47-9,C1B M1B
91082-50-7,C1B
91082-52-9,C1B M1B
91082-53-0,C1B M1B
91697-23-3,C1B
91697-41-5,C1A M2 R1B
91770-57-9,C1B
91995-14-1,C1B
91995-15-2,C1B M1B
91995-16-3,C1B M1B
91995-17-4,C1B M1B
91995-18-5,C1B M1B
91995-20-9,C1B M1B
91995-31-2,C1B M1B
91995-34-5,C1B
91995-35-6,C1B M1B
91995-38-9,C1B M1B
91995-39-0,C1B
91995-40-3,C1B
91995-41-4,C1B M1B
91995-42-5,C1B
91995-45-8,C1B
91995-48-1,C1B M1B
91995-49-2,C1B M1B
91995-50-5,C1B M1B
91995-51-6,C1B
91995-52-7,C1B
91995-53-8,C1B M1B
91995-54-9,C1B
91995-61-8,C1B M1B
91995-66-3,C1B M1B
91995-68-5,C1B M1B
91995-73-2,C1B
91995-75-4,C1B
91995-76-5,C1B
91995-77-6,C1B
91995-78-7,C1B
91995-79-8,C1B
92-67-1,C1A
92-87-5,C1A
92-93-3,C1B
92045-12-0,C1B
92045-14-2,C1B
92045-15-3,C1A M1B
92045-16-4,C1A M1B
92045-17-5,C1A M1B
92045-18-6,C1A M1B
92045-19-7,C1A M1B
92045-20-0,C1A M1B
92045-22-2,C1A M1B
92045-23-3,C1A M1B
92045-29-9,C1B
92045-42-6,C1B
92045-43-7,C1B
92045-49-3,C1B M1B
92045-50-6,C1B M1B
92045-51-7,C1B M1B
92045-52-8,C1B M1B
92045-53-9,C1B M1B
92045-55-1,C1B M1B
92045-57-3,C1B M1B
92045-58-4,C1B M1B
92045-59-5,C1B M1B
92045-60-8,C1B M1B
92045-61-9,C1B M1B
92045-62-0,C1B M1B
92045-63-1,C1B M1B
92045-64-2,C1B M1B
92045-65-3,C1B M1B
92045-71-1,C1B
92045-72-2,C1B
92045-77-7,C1B
92045-80-2,C1A M1B
92061-86-4,C1B
92061-92-2,C1B M1B
92061-93-3,C1B
92061-94-4,C1B
92061-97-7,C1B
92062-00-5,C1B
92062-04-9,C1B
92062-09-4,C1B
92062-10-7,C1B
92062-11-8,C1B
92062-15-2,C1B M1B
92062-20-9,C1B
92062-22-1,C1B M1B
92062-26-5,C1B M1B
92062-27-6,C1B M1B
92062-28-7,C1B M1B
92062-29-8,C1B M1B
92062-33-4,C1B M1B
92062-34-5,C1B
92062-36-7,C1B M1B
92128-94-4,C1B M1B
92129-09-4,C1B
92129-57-2,C1A M2 R1B
92201-59-7,C1B
92201-60-0,C1B
92201-97-3,C1B M1B
92704-08-0,C1B
93107-30-3,R2
93165-19-6,C1B M1B
93165-55-0,C1B M1B
93571-75-6,C1B M1B
93572-29-3,C1B M1B
93572-35-1,C1B M1B
93572-36-2,C1B M1B
93572-43-1,C1B
93629-90-4,M2
93763-10-1,C1B
93763-11-2,C1B
93763-33-8,C1B M1B
93763-34-9,C1B M1B
93763-38-3,C1B
93763-85-0,C1B
93821-38-6,C1B M1B
93821-66-0,C1B
93920-09-3,C1A M2 R1B
93920-10-6,C1A M2 R1B
93924-31-3,C1B
93924-32-4,C1B
93924-33-5,C1B
93924-61-9,C1B
93983-68-7,C1A M2 R1B
94-59-7,C1B M2
94114-03-1,C1B M1B
94114-13-3,C1B
94114-29-1,C1B M1B
94114-40-6,C1B M1B
94114-46-2,C1B
94114-47-3,C1B
94114-48-4,C1B
94114-52-0,C1B M1B
94114-53-1,C1B M1B
94114-54-2,C1B M1B
94114-55-3,C1B
94114-56-4,C1B M1B
94114-57-5,C1B M1B
94114-58-6,C2
94114-59-7,C2
94247-67-3,M2
94361-06-5,R1B
94551-87-8,C1A M2 R1A
94723-86-1,R1B
94733-08-1,C1B
94733-09-2,C1B
94733-15-0,C1B
94733-16-1,C1B
95-06-7,C1B
95-53-4,C1B
95-54-5,C2 M2
95-55-6,M2
95-69-2,C1B M2
95-80-7,C1B M2 R2
95009-23-7,C1B M1B
95371-04-3,C1B
95371-05-4,C1B
95371-07-6,C1B
95371-08-7,C1B
95465-89-7,C1A M1B
96-09-3,C1B
96-12-8,C1B M1B R1A
96-13-9,C1B R2
96-18-4,C1B R1B
96-23-1,C1B
96-29-7,C2
96-45-7,R1B
96314-26-0,R2
96690-55-0,C1B M1B
97-56-3,C1B
97-99-4,R1B
97488-73-8,C1B
97488-74-9,C1B
97488-95-4,C1B
97488-96-5,C1B
97675-85-9,C1B
97675-86-0,C1B
97675-87-1,C1B
97675-88-2,C2
97722-04-8,C1B
97722-06-0,C1B
97722-08-2,C1B
97722-09-3,C1B
97722-10-6,C1B
97722-19-5,C1A M1B
97862-76-5,C1B
97862-77-6,C1B
97862-78-7,C1B
97862-81-2,C1B
97862-82-3,C1B
97862-83-4,C1B
97862-97-0,C1B
97862-98-1,C1B
97863-04-2,C1B
97863-05-3,C1B
97863-06-4,C1B
97926-43-7,C1B M1B
97926-59-5,C1B
97926-68-6,C1B
97926-70-0,C1B
97926-71-1,C1B
97926-76-6,C1B
97926-77-7,C1B
97926-78-8,C1B
98-00-0,C2
98-01-1,C2
98-07-7,C1B
98-54-4,R2
98-73-7,R1B
98-87-3,C2
98-95-3,C2 R1B
98219-46-6,C1B M1B
98219-47-7,C1B M1B
98219-64-8,C1B
99-55-8,C2
99105-77-8,R2
993-16-8,R2
99464-83-2,M2
99610-72-7,R2
`

// CMR_H is a list of H phrases that are CMRs
// They are NOT inserted into the database.
var CMR_H = map[string]string{
	"H340":   "M1",
	"H341":   "M2",
	"H350":   "C1",
	"H350i":  "C1",
	"H351":   "C2",
	"H360":   "R1",
	"H360F":  "R1",
	"H360D":  "R1",
	"H360Fd": "R1",
	"H360Df": "R1",
	"H360FD": "R1",
	"H361":   "R2",
	"H361f":  "R2",
	"H361d":  "R2",
	"H361fd": "R2",
	"H362":   "L",
}
