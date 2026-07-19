-- +goose Up
-- +goose StatementBegin
-- insert cat breeds
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Abyssinian'),
    ('American Bobtail'),
    ('American Curl'),
    ('American Shorthair'),
    ('American Wirehair'),
    ('Applehead Siamese'),
    ('Balinese'),
    ('Bengal'),
    ('Birman'),
    ('Bombay'),
    ('British Shorthair'),
    ('Burmese'),
    ('Burmilla'),
    ('Calico'),
    ('Canadian Hairless'),
    ('Chartreux'),
    ('Chausie'),
    ('Chinchilla'),
    ('Cornish Rex'),
    ('Cymric'),
    ('Devon Rex'),
    ('Dilute Calico'),
    ('Dilute Tortoiseshell'),
    ('Domestic Long Hair'),
    ('Domestic Medium Hair'),
    ('Domestic Short Hair'),
    ('Egyptian Mau'),
    ('Exotic Shorthair'),
    ('Extra-Toes Cat / Hemingway Polydactyl'),
    ('Havana'),
    ('Himalayan'),
    ('Japanese Bobtail'),
    ('Javanese'),
    ('Korat'),
    ('LaPerm'),
    ('Maine Coon'),
    ('Manx'),
    ('Munchkin'),
    ('Nebelung'),
    ('Norwegian Forest Cat'),
    ('Ocicat'),
    ('Oriental Long Hair'),
    ('Oriental Short Hair'),
    ('Oriental Tabby'),
    ('Persian'),
    ('Pixiebob'),
    ('Ragamuffin'),
    ('Ragdoll'),
    ('Russian Blue'),
    ('Scottish Fold'),
    ('Selkirk Rex'),
    ('Siamese'),
    ('Siberian'),
    ('Silver'),
    ('Singapura'),
    ('Snowshoe'),
    ('Somali'),
    ('Sphynx / Hairless Cat'),
    ('Tabby'),
    ('Tiger'),
    ('Tonkinese'),
    ('Torbie'),
    ('Tortoiseshell'),
    ('Toyger'),
    ('Turkish Angora'),
    ('Turkish Van'),
    ('Tuxedo'),
    ('York Chocolat')
) as b("name")
where asp."name" = 'Cat';


-- insert dog breeds
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Affenpinscher'),
    ('Afghan Hound'),
    ('Airedale Terrier'),
    ('Akbash'),
    ('Akita'),
    ('Alaskan Malamute'),
    ('American Bulldog'),
    ('American Bully'),
    ('American Eskimo Dog'),
    ('American Foxhound'),
    ('American Hairless Terrier'),
    ('American Staffordshire Terrier'),
    ('American Water Spaniel'),
    ('Anatolian Shepherd'),
    ('Appenzell Mountain Dog'),
    ('Aussiedoodle'),
    ('Australian Cattle Dog / Blue Heeler'),
    ('Australian Kelpie'),
    ('Australian Shepherd'),
    ('Australian Terrier'),
    ('Basenji'),
    ('Basset Hound'),
    ('Beagle'),
    ('Bearded Collie'),
    ('Beauceron'),
    ('Bedlington Terrier'),
    ('Belgian Shepherd / Laekenois'),
    ('Belgian Shepherd / Malinois'),
    ('Belgian Shepherd / Sheepdog'),
    ('Belgian Shepherd / Tervuren'),
    ('Bernedoodle'),
    ('Bernese Mountain Dog'),
    ('Bichon Frise'),
    ('Black and Tan Coonhound'),
    ('Black Labrador Retriever'),
    ('Black Mouth Cur'),
    ('Black Russian Terrier'),
    ('Bloodhound'),
    ('Blue Lacy'),
    ('Bluetick Coonhound'),
    ('Boerboel'),
    ('Bolognese'),
    ('Border Collie'),
    ('Border Terrier'),
    ('Borzoi'),
    ('Boston Terrier'),
    ('Bouvier des Flandres'),
    ('Boxer'),
    ('Boykin Spaniel'),
    ('Briard'),
    ('Brittany Spaniel'),
    ('Brussels Griffon'),
    ('Bull Terrier'),
    ('Bullmastiff'),
    ('Cairn Terrier'),
    ('Canaan Dog'),
    ('Cane Corso'),
    ('Cardigan Welsh Corgi'),
    ('Carolina Dog'),
    ('Catahoula Leopard Dog'),
    ('Cattle Dog'),
    ('Caucasian Sheepdog / Caucasian Ovtcharka'),
    ('Cavachon'),
    ('Cavalier King Charles Spaniel'),
    ('Cavapoo'),
    ('Chesapeake Bay Retriever'),
    ('Chihuahua'),
    ('Chinese Crested Dog'),
    ('Chinese Foo Dog'),
    ('Chinook'),
    ('Chiweenie'),
    ('Chocolate Labrador Retriever'),
    ('Chow Chow'),
    ('Cirneco dell''Etna'),
    ('Clumber Spaniel'),
    ('Cockapoo'),
    ('Cocker Spaniel'),
    ('Collie'),
    ('Coonhound'),
    ('Corgi'),
    ('Coton de Tulear'),
    ('Curly-Coated Retriever'),
    ('Dachshund'),
    ('Dalmatian'),
    ('Dandie Dinmont Terrier'),
    ('Doberman Pinscher'),
    ('Dogo Argentino'),
    ('Dogue de Bordeaux'),
    ('Dutch Shepherd'),
    ('English Bulldog'),
    ('English Cocker Spaniel'),
    ('English Coonhound'),
    ('English Foxhound'),
    ('English Pointer'),
    ('English Setter'),
    ('English Shepherd'),
    ('English Springer Spaniel'),
    ('English Toy Spaniel'),
    ('Entlebucher'),
    ('Eskimo Dog'),
    ('Feist'),
    ('Field Spaniel'),
    ('Fila Brasileiro'),
    ('Finnish Lapphund'),
    ('Finnish Spitz'),
    ('Flat-Coated Retriever'),
    ('Fox Terrier'),
    ('Foxhound'),
    ('French Bulldog'),
    ('Galgo Spanish Greyhound'),
    ('German Pinscher'),
    ('German Shepherd Dog'),
    ('German Shorthaired Pointer'),
    ('German Spitz'),
    ('German Wirehaired Pointer'),
    ('Giant Schnauzer'),
    ('Glen of Imaal Terrier'),
    ('Golden Retriever'),
    ('Goldendoodle'),
    ('Gordon Setter'),
    ('Great Dane'),
    ('Great Pyrenees'),
    ('Greater Swiss Mountain Dog'),
    ('Greyhound'),
    ('Hamiltonstovare'),
    ('Harrier'),
    ('Havanese'),
    ('Hound'),
    ('Hovawart'),
    ('Husky'),
    ('Ibizan Hound'),
    ('Icelandic Sheepdog'),
    ('Illyrian Sheepdog'),
    ('Irish Setter'),
    ('Irish Terrier'),
    ('Irish Water Spaniel'),
    ('Irish Wolfhound'),
    ('Italian Greyhound'),
    ('Jack Russell Terrier'),
    ('Japanese Chin'),
    ('Jindo'),
    ('Kai Dog'),
    ('Karelian Bear Dog'),
    ('Keeshond'),
    ('Kerry Blue Terrier'),
    ('Kishu'),
    ('Klee Kai'),
    ('Komondor'),
    ('Kuvasz'),
    ('Kyi Leo'),
    ('Labradoodle'),
    ('Labrador Retriever'),
    ('Lakeland Terrier'),
    ('Lancashire Heeler'),
    ('Leonberger'),
    ('Lhasa Apso'),
    ('Lowchen'),
    ('Lurcher'),
    ('Maltese'),
    ('Maltipoo'),
    ('Manchester Terrier'),
    ('Maremma Sheepdog'),
    ('Mastiff'),
    ('McNab'),
    ('Miniature Bull Terrier'),
    ('Miniature Dachshund'),
    ('Miniature Pinscher'),
    ('Miniature Poodle'),
    ('Miniature Schnauzer'),
    ('Mixed Breed'),
    ('Morkie'),
    ('Mountain Cur'),
    ('Mountain Dog'),
    ('Munsterlander'),
    ('Neapolitan Mastiff'),
    ('New Guinea Singing Dog'),
    ('Newfoundland Dog'),
    ('Norfolk Terrier'),
    ('Norwegian Buhund'),
    ('Norwegian Elkhound'),
    ('Norwegian Lundehund'),
    ('Norwich Terrier'),
    ('Nova Scotia Duck Tolling Retriever'),
    ('Old English Sheepdog'),
    ('Otterhound'),
    ('Papillon'),
    ('Parson Russell Terrier'),
    ('Patterdale Terrier / Fell Terrier'),
    ('Pekingese'),
    ('Pembroke Welsh Corgi'),
    ('Peruvian Inca Orchid'),
    ('Petit Basset Griffon Vendeen'),
    ('Pharaoh Hound'),
    ('Pit Bull Terrier'),
    ('Plott Hound'),
    ('Pointer'),
    ('Polish Lowland Sheepdog'),
    ('Pomeranian'),
    ('Pomsky'),
    ('Poodle'),
    ('Portuguese Podengo'),
    ('Portuguese Water Dog'),
    ('Presa Canario'),
    ('Pug'),
    ('Puggle'),
    ('Puli'),
    ('Pumi'),
    ('Pyrenean Shepherd'),
    ('Rat Terrier'),
    ('Redbone Coonhound'),
    ('Retriever'),
    ('Rhodesian Ridgeback'),
    ('Rottweiler'),
    ('Rough Collie'),
    ('Saint Bernard'),
    ('Saluki'),
    ('Samoyed'),
    ('Sarplaninac'),
    ('Schipperke'),
    ('Schnauzer'),
    ('Schnoodle'),
    ('Scottish Deerhound'),
    ('Scottish Terrier'),
    ('Sealyham Terrier'),
    ('Setter'),
    ('Shar-Pei'),
    ('Sheep Dog'),
    ('Sheepadoodle'),
    ('Shepherd'),
    ('Shetland Sheepdog / Sheltie'),
    ('Shiba Inu'),
    ('Shih poo'),
    ('Shih Tzu'),
    ('Shollie'),
    ('Siberian Husky'),
    ('Silky Terrier'),
    ('Skye Terrier'),
    ('Sloughi'),
    ('Smooth Collie'),
    ('Smooth Fox Terrier'),
    ('South Russian Ovtcharka'),
    ('Spaniel'),
    ('Spanish Water Dog'),
    ('Spinone Italiano'),
    ('Spitz'),
    ('Staffordshire Bull Terrier'),
    ('Standard Poodle'),
    ('Standard Schnauzer'),
    ('Sussex Spaniel'),
    ('Swedish Vallhund'),
    ('Tennessee Treeing Brindle'),
    ('Terrier'),
    ('Thai Ridgeback'),
    ('Tibetan Mastiff'),
    ('Tibetan Spaniel'),
    ('Tibetan Terrier'),
    ('Tosa Inu'),
    ('Toy Fox Terrier'),
    ('Toy Manchester Terrier'),
    ('Treeing Walker Coonhound'),
    ('Vizsla'),
    ('Weimaraner'),
    ('Welsh Springer Spaniel'),
    ('Welsh Terrier'),
    ('West Highland White Terrier / Westie'),
    ('Wheaten Terrier'),
    ('Whippet'),
    ('White German Shepherd'),
    ('Wire Fox Terrier'),
    ('Wirehaired Dachshund'),
    ('Wirehaired Pointing Griffon'),
    ('Wirehaired Terrier'),
    ('Xoloitzcuintli / Mexican Hairless'),
    ('Yellow Labrador Retriever'),
    ('Yorkshire Terrie')
) as b("name")
where asp."name" = 'Dog';


-- insert rabbit breeds
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('American'),
    ('American Fuzzy Lop'),
    ('American Sable'),
    ('Angora Rabbit'),
    ('Belgian Hare'),
    ('Beveren'),
    ('Britannia Petite'),
    ('Bunny Rabbit'),
    ('Californian'),
    ('Champagne D''Argent'),
    ('Checkered Giant'),
    ('Chinchilla'),
    ('Cinnamon'),
    ('Creme D''Argent'),
    ('Dutch'),
    ('Dwarf'),
    ('Dwarf Eared'),
    ('English Lop'),
    ('English Spot'),
    ('Flemish Giant'),
    ('Florida White'),
    ('French Lop'),
    ('Harlequin'),
    ('Havana'),
    ('Himalayan'),
    ('Holland Lop'),
    ('Hotot'),
    ('Jersey Wooly'),
    ('Lilac'),
    ('Lionhead'),
    ('Lop Eared'),
    ('Mini Lop'),
    ('Mini Rex'),
    ('Netherland Dwarf'),
    ('New Zealand'),
    ('Palomino'),
    ('Polish'),
    ('Rex'),
    ('Rhinelander'),
    ('Satin'),
    ('Silver'),
    ('Silver Fox'),
    ('Silver Marten'),
    ('Ta')
) as b("name")
where asp."name" = 'Rabbit';


-- insert horse breeds (Donkey species)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Donkey')
) as b("name")
where asp."name" = 'Donkey';

-- insert horse breeds (Miniature Horse)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Miniature Horse')
) as b("name")
where asp."name" = 'Miniature Horse';

-- insert horse breeds (Mule)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Mule')
) as b("name")
where asp."name" = 'Mule';

-- insert horse breeds (Pony)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Connemara'),
    ('Pony'),
    ('Pony of the Americas'),
    ('Shetland Pony')
) as b("name")
where asp."name" = 'Pony';

-- insert horse breeds (Horse)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Appaloosa'),
    ('Arabian'),
    ('Belgian'),
    ('Clydesdale'),
    ('Curly Horse'),
    ('Draft'),
    ('Friesian'),
    ('Gaited'),
    ('Grade'),
    ('Haflinger'),
    ('Icelandic Horse'),
    ('Lipizzan'),
    ('Missouri Foxtrotter'),
    ('Morgan'),
    ('Mule'),
    ('Mustang'),
    ('Paint / Pinto'),
    ('Palomino'),
    ('Paso Fino'),
    ('Percheron'),
    ('Peruvian Paso'),
    ('Quarterhorse'),
    ('Rocky Mountain Horse'),
    ('Saddlebred'),
    ('Standardbred'),
    ('Tennessee Walker'),
    ('Thoroughbred'),
    ('Warmblood')
) as b("name")
where asp."name" = 'Horse';


-- insert bird breeds (Button-Quail)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Button-Quail')
) as b("name")
where asp."name" = 'Button-Quail';

-- insert bird breeds (Chicken)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Chicken')
) as b("name")
where asp."name" = 'Chicken';

-- insert bird breeds (Dove)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Dove'),
    ('Pigeon')
) as b("name")
where asp."name" = 'Dove';

-- insert bird breeds (Duck)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Duck')
) as b("name")
where asp."name" = 'Duck';

-- insert bird breeds (Emu)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Emu')
) as b("name")
where asp."name" = 'Emu';

-- insert bird breeds (Finch)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Finch')
) as b("name")
where asp."name" = 'Finch';

-- insert bird breeds (Goose)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Goose')
) as b("name")
where asp."name" = 'Goose';

-- insert bird breeds (Guinea Fowl)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Guinea Fowl')
) as b("name")
where asp."name" = 'Guinea Fowl';

-- insert bird breeds (Ostritch)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Ostritch')
) as b("name")
where asp."name" = 'Ostritch';

-- insert bird breeds (Parakeet)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Parakeet (Other)')
) as b("name")
where asp."name" = 'Parakeet';

-- insert bird breeds (Peacock / Peafowl)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Peacock / Peafowl')
) as b("name")
where asp."name" = 'Peacock';

-- insert bird breeds (Pheasant)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Pheasant')
) as b("name")
where asp."name" = 'Pheasant';

-- insert bird breeds (Quail)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Quail')
) as b("name")
where asp."name" = 'Quail';

-- insert bird breeds (Rhea)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Rhea')
) as b("name")
where asp."name" = 'Rhea';

-- insert bird breeds (Swan)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Swan')
) as b("name")
where asp."name" = 'Swan';

-- insert bird breeds (Toucan)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Toucan')
) as b("name")
where asp."name" = 'Toucan';

-- insert bird breeds (Turkey)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Turkey')
) as b("name")
where asp."name" = 'Turkey';

-- insert bird breeds (Parrot)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('African Grey'),
    ('Amazon'),
    ('Brotogeris'),
    ('Budgie / Budgerigar'),
    ('Caique'),
    ('Canary'),
    ('Chicken'),
    ('Cockatiel'),
    ('Cockatoo'),
    ('Conure'),
    ('Eclectus'),
    ('Kakariki'),
    ('Lory / Lorikeet'),
    ('Lovebird'),
    ('Macaw'),
    ('Ostrich'),
    ('Parrot (Other)'),
    ('Parrotlet'),
    ('Pionus'),
    ('Poicephalus / Senegal'),
    ('Quaker Parakeet'),
    ('Ringneck / Psittacula'),
    ('Rosella')
) as b("name")
where asp."name" = 'Parrot';


-- insert barnyard breeds (Alpaca)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Alpaca')
) as b("name")
where asp."name" = 'Alpaca';

-- insert barnyard breeds (Cow)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Angus'),
    ('Cow'),
    ('Holstein'),
    ('Jersey')
) as b("name")
where asp."name" = 'Cow';

-- insert barnyard breeds (Goat)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Alpine'),
    ('Angora'),
    ('Boer'),
    ('Goat'),
    ('LaMancha'),
    ('Myotonic / Fainting'),
    ('Nigerian Dwarf'),
    ('Nubian'),
    ('Oberhasli'),
    ('Pygmy'),
    ('Saanen'),
    ('Toggenburg')
) as b("name")
where asp."name" = 'Goat';

-- insert barnyard breeds (Pig)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Duroc'),
    ('Hampshire'),
    ('Landrace'),
    ('Pig'),
    ('Yorkshire')
) as b("name")
where asp."name" = 'Pig';

-- insert barnyard breeds (Pot Bellied)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Pot Bellied'),
    ('Vietnamese Pot Bellied')
) as b("name")
where asp."name" = 'Pot Bellied';

-- insert barnyard breeds (Sheep)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Barbados'),
    ('Merino'),
    ('Mouflon'),
    ('Sheep'),
    ('Shetland')
) as b("name")
where asp."name" = 'Sheep';


-- insert small & furry breeds (Chinchilla)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Chinchilla')
) as b("name")
where asp."name" = 'Chinchilla';

-- insert small & furry breeds (Degu)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Degu')
) as b("name")
where asp."name" = 'Degu';

-- insert small & furry breeds (Ferret)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Ferret')
) as b("name")
where asp."name" = 'Ferret';

-- insert small & furry breeds (Gerbil)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Gerbil')
) as b("name")
where asp."name" = 'Gerbil';

-- insert small & furry breeds (Guinea Pig)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Abyssinian'),
    ('Guinea Pig'),
    ('Peruvian'),
    ('Rex'),
    ('Short-Haired'),
    ('Silkie'),
    ('Sheltie'),
    ('Teddy')
) as b("name")
where asp."name" = 'Guinea Pig';

-- insert small & furry breeds (Hamster)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Hamster'),
    ('Dwarf Hamster')
) as b("name")
where asp."name" = 'Hamster';

-- insert small & furry breeds (Hedgehog)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Hedgehog')
) as b("name")
where asp."name" = 'Hedgehog';

-- insert small & furry breeds (Mouse)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Mouse')
) as b("name")
where asp."name" = 'Mouse';

-- insert small & furry breeds (Prairie Dog)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Prairie Dog')
) as b("name")
where asp."name" = 'Prairie Dog';

-- insert small & furry breeds (Rat)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Rat')
) as b("name")
where asp."name" = 'Rat';

-- insert small & furry breeds (Skunk)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Skunk')
) as b("name")
where asp."name" = 'Skunk';

-- insert small & furry breeds (Sugar Glider)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Sugar Glider')
) as b("name")
where asp."name" = 'Sugar Glider';


-- insert scales, fins & other breeds (Amphibian)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Bullfrog'),
    ('Frog'),
    ('Horned Frog'),
    ('Leopard Frog'),
    ('Tree Frog')
) as b("name")
where asp."name" = 'Amphibian';

-- insert scales, fins & other breeds (Fish)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Freshwater Fish'),
    ('Saltwater Fish'),
    ('Goldfish')
) as b("name")
where asp."name" = 'Fish';

-- insert scales, fins & other breeds (Reptile)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Bearded Dragon'),
    ('Chameleon'),
    ('Gecko'),
    ('Iguana'),
    ('Lizard'),
    ('Monitor'),
    ('Uromastyx'),
    ('Water Dragon')
) as b("name")
where asp."name" = 'Reptile';

-- insert scales, fins & other breeds (Salamander / Newt)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Fire Salamander'),
    ('Fire-Bellied Newt'),
    ('Oregon Newt'),
    ('Paddle Tailed Newt'),
    ('Tiger Salamander')
) as b("name")
where asp."name" = 'Salamander / Newt';

-- insert scales, fins & other breeds (Snake)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Ball Python'),
    ('Boa'),
    ('Boa Constrictor'),
    ('Bull'),
    ('Burmese Python'),
    ('Corn / Ribbon'),
    ('King / Milk'),
    ('Python'),
    ('Snake')
) as b("name")
where asp."name" = 'Snake';

-- insert scales, fins & other breeds (Toad)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Fire-Bellied'),
    ('Southern'),
    ('Toad')
) as b("name")
where asp."name" = 'Toad';

-- insert scales, fins & other breeds (Tortoise)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Leopard'),
    ('Red Foot'),
    ('Russian'),
    ('Sulcata'),
    ('Tortoise')
) as b("name")
where asp."name" = 'Tortoise';

-- insert scales, fins & other breeds (Turtle)
insert into "breed" ("animal_specie_id", "name")
select asp."id", b."name"
from "animal_specie" asp
cross join (
  values
    ('Asian Box'),
    ('Box'),
    ('Eastern Box'),
    ('Florida Box'),
    ('Mississippi Map Turtle'),
    ('Mud'),
    ('Musk'),
    ('Ornamental Box'),
    ('Painted'),
    ('Read-Earder Slider'),
    ('Snapping'),
    ('Soft Shell'),
    ('Three-Toed Box'),
    ('Turtle'),
    ('Yellow-Bellied Slider')
) as b("name")
where asp."name" = 'Turtle';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate table "breed" cascade;
-- +goose StatementEnd
