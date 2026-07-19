-- +goose Up
-- +goose StatementBegin
-- update Cat properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Black & White / Tuxedo",
        "Blue Cream",
        "Blue Point",
        "Brown / Chocolate",
        "Buff & White",
        "Buff / Tan / Fawn",
        "Calico",
        "Chocolate Point",
        "Cream / Ivory",
        "Cream Point",
        "Dilute Calico",
        "Dilute Tortoiseshell",
        "Flame Point",
        "Gray & White",
        "Gray / Blue / Silver",
        "Lilac Point",
        "Orange & White",
        "Orange / Red",
        "Seal Point",
        "Smoke",
        "Tabby (Brown / Chocolate)",
        "Tabby (Buff / Tan / Fawn)",
        "Tabby (Gray / Blue / Silver)",
        "Tabby (Leopard / Spotted)",
        "Tabby (Orange / Red)",
        "Tabby (Tiger Striped)",
        "Torbie",
        "Tortoiseshell",
        "White"
      ],
      "title": "Color",
      "description": "The color of the cat",
      "x-group": "Traits"
    },
    "coat_length": {
      "type": "string",
      "enum": ["Hairless", "Short", "Medium", "Long"],
      "title": "Coat length",
      "description": "The length of the cat's coat",
      "x-group": "Traits"
    },
    "good_with_cats": {
      "type": "boolean",
      "title": "Good with cats",
      "description": "Whether the cat is good with other cats",
      "x-group": "Behaviour"
    },
    "good_with_dogs": {
      "type": "boolean",
      "title": "Good with dogs",
      "description": "Whether the cat is good with dogs",
      "x-group": "Behaviour"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the cat is good with other animals",
      "x-group": "Behaviour"
    },
    "good_with_children": {
      "type": "boolean",
      "title": "Good with children",
      "description": "Whether the cat is good with children",
      "x-group": "Behaviour"
    },
    "good_with_people": {
      "type": "boolean",
      "title": "Good with people",
      "description": "Whether the cat is good with people",
      "x-group": "Behaviour"
    },
    "house_trained": {
      "type": "boolean",
      "title": "House trained",
      "description": "Whether the cat is house trained",
      "x-group": "Behaviour"
    },
    "declawed": {
      "type": "boolean",
      "title": "Declawed",
      "description": "Whether the cat is declawed",
      "x-group": "Health"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the cat has special needs",
      "x-group": "Health"
    },
    "spayed_neutered": {
      "type": "boolean",
      "title": "Spayed/Neutered",
      "description": "Whether the cat is spayed/neutered",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the cat is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Cat';


-- update Dog properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Apricot / Beige",
        "Bicolor",
        "Black",
        "Brindle",
        "Brown / Chocolate",
        "Golden",
        "Gray / Blue / Silver",
        "Harlequin",
        "Merle (Blue)",
        "Merle (Red)",
        "Red / Chestnut / Orange",
        "Sable",
        "Tricolor (Brown, Black, & White)",
        "White / Cream",
        "Yellow / Tan / Blond / Fawn"
      ],
      "title": "Color",
      "description": "The color of the dog",
      "x-group": "Traits"
    },
    "coat_length": {
      "type": "string",
      "enum": ["Hairless", "Short", "Medium", "Long", "Wire", "Curly"],
      "title": "Coat length",
      "description": "The length of the dog's coat",
      "x-group": "Traits"
    },
    "good_with_cats": {
      "type": "boolean",
      "title": "Good with cats",
      "description": "Whether the dog is good with other cats",
      "x-group": "Behaviour"
    },
    "good_with_dogs": {
      "type": "boolean",
      "title": "Good with dogs",
      "description": "Whether the dog is good with dogs",
      "x-group": "Behaviour"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the dog is good with other animals",
      "x-group": "Behaviour"
    },
    "good_with_children": {
      "type": "boolean",
      "title": "Good with children",
      "description": "Whether the dog is good with children",
      "x-group": "Behaviour"
    },
    "good_with_people": {
      "type": "boolean",
      "title": "Good with people",
      "description": "Whether the dog is good with people",
      "x-group": "Behaviour"
    },
    "house_trained": {
      "type": "boolean",
      "title": "House trained",
      "description": "Whether the dog is house trained",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the dog has special needs",
      "x-group": "Health"
    },
    "spayed_neutered": {
      "type": "boolean",
      "title": "Spayed/Neutered",
      "description": "Whether the dog is spayed/neutered",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the dog is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Dog';


-- update Rabbit properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Black",
        "Blue / Gray",
        "Brown / Chocolate",
        "Cream",
        "Lilac",
        "Orange / Red",
        "Sable",
        "Silver Marten",
        "Tan",
        "Tortoiseshell",
        "White"
      ],
      "title": "Color",
      "description": "The color of the rabbit",
      "x-group": "Traits"
    },
    "coat_length": {
      "type": "string",
      "enum": ["Short", "Medium", "Long"],
      "title": "Coat length",
      "description": "The length of the rabbits's coat",
      "x-group": "Traits"
    },
    "good_with_cats": {
      "type": "boolean",
      "title": "Good with cats",
      "description": "Whether the rabbit is good with other cats",
      "x-group": "Behaviour"
    },
    "good_with_dogs": {
      "type": "boolean",
      "title": "Good with dogs",
      "description": "Whether the rabbit is good with dogs",
      "x-group": "Behaviour"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the rabbit is good with other animals",
      "x-group": "Behaviour"
    },
    "good_with_children": {
      "type": "boolean",
      "title": "Good with children",
      "description": "Whether the rabbit is good with children",
      "x-group": "Behaviour"
    },
    "good_with_people": {
      "type": "boolean",
      "title": "Good with people",
      "description": "Whether the rabbit is good with people",
      "x-group": "Behaviour"
    },
    "house_trained": {
      "type": "boolean",
      "title": "House trained",
      "description": "Whether the rabbit is house trained",
      "x-group": "Behaviour"
    },
    "handling_tolerant": {
      "type": "boolean",
      "title": "Handling tolerant",
      "description": "Whether the rabbit tolerates being picked up and handled",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the dog has special needs",
      "x-group": "Health"
    },
    "spayed_neutered": {
      "type": "boolean",
      "title": "Spayed/Neutered",
      "description": "Whether the rabbit is spayed/neutered",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the rabbit is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Rabbit';


-- update Donkey properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Appaloosa",
        "Bay",
        "Bay Roan",
        "Black",
        "Blue Roan",
        "Brown",
        "Buckskin",
        "Champagne",
        "Chestnut / Sorrel",
        "Cremello",
        "Dapple Gray",
        "Dun",
        "Gray",
        "Grullo",
        "Liver",
        "Paint",
        "Palomino",
        "Perlino",
        "Piebald",
        "Pinto",
        "Red Roan",
        "Silver Bay",
        "Silver Buckskin",
        "Silver Dapple",
        "White"
      ],
      "title": "Color",
      "description": "The color of the donkey",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the donkey is good with other animals",
      "x-group": "Behaviour"
    },
    "good_with_children": {
      "type": "boolean",
      "title": "Good with children",
      "description": "Whether the donkey is good with children",
      "x-group": "Behaviour"
    },
    "good_with_people": {
      "type": "boolean",
      "title": "Good with people",
      "description": "Whether the donkey is good with people",
      "x-group": "Behaviour"
    },
    "handled_regularly": {
      "type": "boolean",
      "title": "Handled regularly",
      "description": "Whether the donkey is accustomed to regular human handling",
      "x-group": "Behaviour"
    },
    "saddle_trained": {
      "type": "boolean",
      "title": "Saddle trained",
      "description": "Whether the donkey is trained to be ridden",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the donkey has special needs",
      "x-group": "Health"
    },
    "castrated": {
      "type": "boolean",
      "title": "Castrated",
      "description": "Whether the donkey is castrated",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the donkey is vaccinated",
      "x-group": "Health"
    },
    "shod": {
      "type": "boolean",
      "title": "Shod",
      "description": "Whether the donkey wears horseshoes",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Donkey';


-- update Horse properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Appaloosa",
        "Bay",
        "Bay Roan",
        "Black",
        "Blue Roan",
        "Brown",
        "Buckskin",
        "Champagne",
        "Chestnut / Sorrel",
        "Cremello",
        "Dapple Gray",
        "Dun",
        "Gray",
        "Grullo",
        "Liver",
        "Paint",
        "Palomino",
        "Perlino",
        "Piebald",
        "Pinto",
        "Red Roan",
        "Silver Bay",
        "Silver Buckskin",
        "Silver Dapple",
        "White"
      ],
      "title": "Color",
      "description": "The color of the horse",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the horse is good with other animals",
      "x-group": "Behaviour"
    },
    "good_with_children": {
      "type": "boolean",
      "title": "Good with children",
      "description": "Whether the horse is good with children",
      "x-group": "Behaviour"
    },
    "good_with_people": {
      "type": "boolean",
      "title": "Good with people",
      "description": "Whether the horse is good with people",
      "x-group": "Behaviour"
    },
    "handled_regularly": {
      "type": "boolean",
      "title": "Handled regularly",
      "description": "Whether the horse is accustomed to regular human handling",
      "x-group": "Behaviour"
    },
    "saddle_trained": {
      "type": "boolean",
      "title": "Saddle trained",
      "description": "Whether the horse is trained to be ridden",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the horse has special needs",
      "x-group": "Health"
    },
    "castrated": {
      "type": "boolean",
      "title": "Castrated",
      "description": "Whether the horse is castrated",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the horse is vaccinated",
      "x-group": "Health"
    },
    "shod": {
      "type": "boolean",
      "title": "Shod",
      "description": "Whether the horse wears horseshoes",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Horse';


-- update Mule properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Appaloosa",
        "Bay",
        "Bay Roan",
        "Black",
        "Blue Roan",
        "Brown",
        "Buckskin",
        "Champagne",
        "Chestnut / Sorrel",
        "Cremello",
        "Dapple Gray",
        "Dun",
        "Gray",
        "Grullo",
        "Liver",
        "Paint",
        "Palomino",
        "Perlino",
        "Piebald",
        "Pinto",
        "Red Roan",
        "Silver Bay",
        "Silver Buckskin",
        "Silver Dapple",
        "White"
      ],
      "title": "Color",
      "description": "The color of the mule",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the mule is good with other animals",
      "x-group": "Behaviour"
    },
    "good_with_children": {
      "type": "boolean",
      "title": "Good with children",
      "description": "Whether the mule is good with children",
      "x-group": "Behaviour"
    },
    "good_with_people": {
      "type": "boolean",
      "title": "Good with people",
      "description": "Whether the mule is good with people",
      "x-group": "Behaviour"
    },
    "handled_regularly": {
      "type": "boolean",
      "title": "Handled regularly",
      "description": "Whether the mule is accustomed to regular human handling",
      "x-group": "Behaviour"
    },
    "saddle_trained": {
      "type": "boolean",
      "title": "Saddle trained",
      "description": "Whether the mule is trained to be ridden",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the mule has special needs",
      "x-group": "Health"
    },
    "castrated": {
      "type": "boolean",
      "title": "Castrated",
      "description": "Whether the mule is castrated",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the mule is vaccinated",
      "x-group": "Health"
    },
    "shod": {
      "type": "boolean",
      "title": "Shod",
      "description": "Whether the mule wears horseshoes",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Mule';


-- update Pony properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Appaloosa",
        "Bay",
        "Bay Roan",
        "Black",
        "Blue Roan",
        "Brown",
        "Buckskin",
        "Champagne",
        "Chestnut / Sorrel",
        "Cremello",
        "Dapple Gray",
        "Dun",
        "Gray",
        "Grullo",
        "Liver",
        "Paint",
        "Palomino",
        "Perlino",
        "Piebald",
        "Pinto",
        "Red Roan",
        "Silver Bay",
        "Silver Buckskin",
        "Silver Dapple",
        "White"
      ],
      "title": "Color",
      "description": "The color of the pony",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the pony is good with other animals",
      "x-group": "Behaviour"
    },
    "good_with_children": {
      "type": "boolean",
      "title": "Good with children",
      "description": "Whether the pony is good with children",
      "x-group": "Behaviour"
    },
    "good_with_people": {
      "type": "boolean",
      "title": "Good with people",
      "description": "Whether the pony is good with people",
      "x-group": "Behaviour"
    },
    "handled_regularly": {
      "type": "boolean",
      "title": "Handled regularly",
      "description": "Whether the pony is accustomed to regular human handling",
      "x-group": "Behaviour"
    },
    "saddle_trained": {
      "type": "boolean",
      "title": "Saddle trained",
      "description": "Whether the pony is trained to be ridden",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the pony has special needs",
      "x-group": "Health"
    },
    "castrated": {
      "type": "boolean",
      "title": "Castrated",
      "description": "Whether the pony is castrated",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the pony is vaccinated",
      "x-group": "Health"
    },
    "shod": {
      "type": "boolean",
      "title": "Shod",
      "description": "Whether the pony wears horseshoes",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Pony';


-- update Miniature Horse properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Appaloosa",
        "Bay",
        "Bay Roan",
        "Black",
        "Blue Roan",
        "Brown",
        "Buckskin",
        "Champagne",
        "Chestnut / Sorrel",
        "Cremello",
        "Dapple Gray",
        "Dun",
        "Gray",
        "Grullo",
        "Liver",
        "Paint",
        "Palomino",
        "Perlino",
        "Piebald",
        "Pinto",
        "Red Roan",
        "Silver Bay",
        "Silver Buckskin",
        "Silver Dapple",
        "White"
      ],
      "title": "Color",
      "description": "The color of the miniature horse",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the miniature horse is good with other animals",
      "x-group": "Behaviour"
    },
    "good_with_children": {
      "type": "boolean",
      "title": "Good with children",
      "description": "Whether the miniature horse is good with children",
      "x-group": "Behaviour"
    },
    "good_with_people": {
      "type": "boolean",
      "title": "Good with people",
      "description": "Whether the miniature horse is good with people",
      "x-group": "Behaviour"
    },
    "handled_regularly": {
      "type": "boolean",
      "title": "Handled regularly",
      "description": "Whether the miniature horse is accustomed to regular human handling",
      "x-group": "Behaviour"
    },
    "saddle_trained": {
      "type": "boolean",
      "title": "Saddle trained",
      "description": "Whether the miniature horse is trained to be ridden",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the miniature horse has special needs",
      "x-group": "Health"
    },
    "castrated": {
      "type": "boolean",
      "title": "Castrated",
      "description": "Whether the miniature horse is castrated",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the miniature horse is vaccinated",
      "x-group": "Health"
    },
    "shod": {
      "type": "boolean",
      "title": "Shod",
      "description": "Whether the miniature horse wears horseshoes",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Miniature Horse';


-- update Button-Quail properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the button-quail",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the button-quail is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the button-quail has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the button-quail is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Button-Quail';


-- update Chicken properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the chicken",
      "x-group": "Traits"
    },
    "egg_production": {
      "type": "string",
      "enum": ["Low", "Medium", "High"],
      "title": "Egg production",
      "description": "Typical egg production level of the chicken",
      "x-group": "Traits"
    },
    "laying_status": {
      "type": "boolean",
      "title": "Currently laying",
      "description": "Whether the chicken is currently laying eggs",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the chicken is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the chicken has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the chicken is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Chicken';


-- update Dove properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the dove",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the dove is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the dove has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the dove is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Dove';


-- update Duck properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the duck",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the duck is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the duck has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the duck is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Duck';


-- update Emu properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the emu",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the emu is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the emu has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the emu is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Emu';


-- update Finch properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the finch",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the finch is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the finch has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the finch is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Finch';


-- update Goose properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the goose",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the goose is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the goose has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the goose is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Goose';


-- update Guinea Fowl properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the guinea fowl",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the guinea fowl is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the guinea fowl has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the guinea fowl is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Guinea Fowl';


-- update Ostritch properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the ostritch",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the ostritch is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the ostritch has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the ostritch is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Ostritch';


-- update Parakeet properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the parakeet",
      "x-group": "Traits"
    },
    "talking_ability": {
      "type": "boolean",
      "title": "Talking ability",
      "description": "Whether the parakeet can talk",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the parakeet is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the parakeet has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the parakeet is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Parakeet';


-- update Parrot properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the parrot",
      "x-group": "Traits"
    },
    "talking_ability": {
      "type": "boolean",
      "title": "Talking ability",
      "description": "Whether the parrot can talk",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the parrot is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the parrot has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the parrot is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Parrot';


-- update Peacock / Peafowl properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the peafowl",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the peafowl is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the peafowl has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the peafowl is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Peacock / Peafowl';


-- update Pheasant properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the pheasant",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the pheasant is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the pheasant has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the pheasant is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Pheasant';


-- update Quail properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the quail",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the quail is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the quail has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the quail is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Quail';


-- update Rhea properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the rhea",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the rhea is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the rhea has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the rhea is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Rhea';


-- update Swan properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the swan",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the swan is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the swan has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the swan is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Swan';


-- update Toucan properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the toucan",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the toucan is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the toucan has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the toucan is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Toucan';


-- update Turkey properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Buff",
        "Gray",
        "Green",
        "Olive",
        "Orange",
        "Pink",
        "Purple / Violet",
        "Red",
        "Rust / Rufous",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the turkey",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the turkey is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the turkey has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the turkey is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Turkey';


-- update Alpaca properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Black",
        "Black & White",
        "Brindle",
        "Brown",
        "Gray",
        "Pink",
        "Red",
        "Roan",
        "Spotted",
        "Tan",
        "White"
      ],
      "title": "Color",
      "description": "The color of the alpaca",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the alpaca is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the alpaca has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the alpaca is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Alpaca';


-- update Cow properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Black",
        "Black & White",
        "Brindle",
        "Brown",
        "Gray",
        "Pink",
        "Red",
        "Roan",
        "Spotted",
        "Tan",
        "White"
      ],
      "title": "Color",
      "description": "The color of the cow",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the cow is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the cow has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the cow is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Cow';


-- update Goat properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Black",
        "Black & White",
        "Brindle",
        "Brown",
        "Gray",
        "Pink",
        "Red",
        "Roan",
        "Spotted",
        "Tan",
        "White"
      ],
      "title": "Color",
      "description": "The color of the goat",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the goat is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the goat has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the goat is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Goat';


-- update Llama properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Black",
        "Black & White",
        "Brindle",
        "Brown",
        "Gray",
        "Pink",
        "Red",
        "Roan",
        "Spotted",
        "Tan",
        "White"
      ],
      "title": "Color",
      "description": "The color of the llama",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the llama is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the llama has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the llama is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Llama';


-- update Pig properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Black",
        "Black & White",
        "Brindle",
        "Brown",
        "Gray",
        "Pink",
        "Red",
        "Roan",
        "Spotted",
        "Tan",
        "White"
      ],
      "title": "Color",
      "description": "The color of the pig",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the pig is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the pig has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the pig is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Pig';


-- update Pot Bellied properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Black",
        "Black & White",
        "Brindle",
        "Brown",
        "Gray",
        "Pink",
        "Red",
        "Roan",
        "Spotted",
        "Tan",
        "White"
      ],
      "title": "Color",
      "description": "The color of the pot bellied",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the pot bellied is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the pot bellied has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the pot bellied is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Pot Bellied';


-- update Sheep properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Black",
        "Black & White",
        "Brindle",
        "Brown",
        "Gray",
        "Pink",
        "Red",
        "Roan",
        "Spotted",
        "Tan",
        "White"
      ],
      "title": "Color",
      "description": "The color of the sheep",
      "x-group": "Traits"
    },
    "coat_length": {
      "type": "string",
      "enum": ["Hairless", "Short", "Medium", "Long"],
      "title": "Coat type",
      "description": "The length of the sheep's coat",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the sheep is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the sheep has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the sheep is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Sheep';


-- update Guinea Pig properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "coat_length": {
      "type": "string",
      "enum": ["Hairless", "Short", "Medium", "Long"],
      "title": "Coat length",
      "description": "The length of the guinea pig's coat",
      "x-group": "Traits"
    },
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Albino",
        "Black",
        "Black Sable",
        "Blue / Gray",
        "Brown / Chocolate",
        "Calico",
        "Champagne",
        "Cinnamon",
        "Cream",
        "Orange / Red",
        "Sable",
        "Tan",
        "Tortoiseshell",
        "White",
        "White (Dark-Eyed)"
      ],
      "title": "Color",
      "description": "The color of the guinea pig",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the guinea pig is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the guinea pig has special needs",
      "x-group": "Health"
    },
    "spayed_neutered": {
      "type": "boolean",
      "title": "Spayed/Neutered",
      "description": "Whether the guinea pig is spayed/neutered",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the guinea pig is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Guinea Pig';


-- update Chinchilla properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "coat_length": {
      "type": "string",
      "enum": ["Short", "Medium", "Long"],
      "title": "Coat length",
      "description": "The length of the chinchilla's coat",
      "x-group": "Traits"
    },
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Albino",
        "Black",
        "Black Sable",
        "Blue / Gray",
        "Brown / Chocolate",
        "Calico",
        "Champagne",
        "Cinnamon",
        "Cream",
        "Orange / Red",
        "Sable",
        "Tan",
        "Tortoiseshell",
        "White",
        "White (Dark-Eyed)"
      ],
      "title": "Color",
      "description": "The color of the chincilla",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the chinchilla is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the chinchilla has special needs",
      "x-group": "Health"
    },
    "spayed_neutered": {
      "type": "boolean",
      "title": "Spayed/Neutered",
      "description": "Whether the chinchilla is spayed/neutered",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the chinchilla is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Chinchilla';


-- update Rat properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Albino",
        "Black",
        "Black Sable",
        "Blue / Gray",
        "Brown / Chocolate",
        "Calico",
        "Champagne",
        "Cinnamon",
        "Cream",
        "Orange / Red",
        "Sable",
        "Tan",
        "Tortoiseshell",
        "White",
        "White (Dark-Eyed)"
      ],
      "title": "Color",
      "description": "The color of the rat",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the rat is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the rat has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the rat is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Rat';


-- update Mouse properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Albino",
        "Black",
        "Black Sable",
        "Blue / Gray",
        "Brown / Chocolate",
        "Calico",
        "Champagne",
        "Cinnamon",
        "Cream",
        "Orange / Red",
        "Sable",
        "Tan",
        "Tortoiseshell",
        "White",
        "White (Dark-Eyed)"
      ],
      "title": "Color",
      "description": "The color of the mouse",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the mouse is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the mouse has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the mouse is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Mouse';


-- update Ferret properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Albino",
        "Black",
        "Black Sable",
        "Blue / Gray",
        "Brown / Chocolate",
        "Calico",
        "Champagne",
        "Cinnamon",
        "Cream",
        "Orange / Red",
        "Sable",
        "Tan",
        "Tortoiseshell",
        "White",
        "White (Dark-Eyed)"
      ],
      "title": "Color",
      "description": "The color of the ferret",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the ferret is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the ferret has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the ferret is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Ferret';


-- update Hamster properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Albino",
        "Black",
        "Black Sable",
        "Blue / Gray",
        "Brown / Chocolate",
        "Calico",
        "Champagne",
        "Cinnamon",
        "Cream",
        "Orange / Red",
        "Sable",
        "Tan",
        "Tortoiseshell",
        "White",
        "White (Dark-Eyed)"
      ],
      "title": "Color",
      "description": "The color of the hamster",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the hamster is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the hamster has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the hamster is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Hamster';


-- update Hedgehog properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Albino",
        "Black",
        "Black Sable",
        "Blue / Gray",
        "Brown / Chocolate",
        "Calico",
        "Champagne",
        "Cinnamon",
        "Cream",
        "Orange / Red",
        "Sable",
        "Tan",
        "Tortoiseshell",
        "White",
        "White (Dark-Eyed)"
      ],
      "title": "Color",
      "description": "The color of the hedgehog",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the hedgehog is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the hedgehog has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the hedgehog is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Hedgehog';


-- update Gerbil properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Albino",
        "Black",
        "Black Sable",
        "Blue / Gray",
        "Brown / Chocolate",
        "Calico",
        "Champagne",
        "Cinnamon",
        "Cream",
        "Orange / Red",
        "Sable",
        "Tan",
        "Tortoiseshell",
        "White",
        "White (Dark-Eyed)"
      ],
      "title": "Color",
      "description": "The color of the gerbil",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the gerbil is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the gerbil has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the gerbil is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Gerbil';


-- update Degu properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Albino",
        "Black",
        "Black Sable",
        "Blue / Gray",
        "Brown / Chocolate",
        "Calico",
        "Champagne",
        "Cinnamon",
        "Cream",
        "Orange / Red",
        "Sable",
        "Tan",
        "Tortoiseshell",
        "White",
        "White (Dark-Eyed)"
      ],
      "title": "Color",
      "description": "The color of the degu",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the degu is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the degu has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the degu is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Degu';


-- update Prairie Dog properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Albino",
        "Black",
        "Black Sable",
        "Blue / Gray",
        "Brown / Chocolate",
        "Calico",
        "Champagne",
        "Cinnamon",
        "Cream",
        "Orange / Red",
        "Sable",
        "Tan",
        "Tortoiseshell",
        "White",
        "White (Dark-Eyed)"
      ],
      "title": "Color",
      "description": "The color of the prairie dog",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the prairie dog is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the prairie dog has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the prairie dog is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Prairie Dog';


-- update Skunk properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Albino",
        "Black",
        "Black Sable",
        "Blue / Gray",
        "Brown / Chocolate",
        "Calico",
        "Champagne",
        "Cinnamon",
        "Cream",
        "Orange / Red",
        "Sable",
        "Tan",
        "Tortoiseshell",
        "White",
        "White (Dark-Eyed)"
      ],
      "title": "Color",
      "description": "The color of the skunk",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the skunk is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the skunk has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the skunk is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Skunk';


-- update Sugar Glider properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Agouti",
        "Albino",
        "Black",
        "Black Sable",
        "Blue / Gray",
        "Brown / Chocolate",
        "Calico",
        "Champagne",
        "Cinnamon",
        "Cream",
        "Orange / Red",
        "Sable",
        "Tan",
        "Tortoiseshell",
        "White",
        "White (Dark-Eyed)"
      ],
      "title": "Color",
      "description": "The color of the sugar glider",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the sugar glider is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the sugar glider has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the sugar glider is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Sugar Glider';


-- update Amphibian properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Gray",
        "Green",
        "Iridescent",
        "Orange",
        "Purple",
        "Red",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the amphibian",
      "x-group": "Traits"
    },
    "diet_type": {
      "type": "string",
      "enum": ["herbivore", "carnivore", "omnivore", "insectivore", "special"],
      "title": "Diet type",
      "description": "Main dietary requirement of the amphibian",
      "x-group": "Traits"
    },
    "humidity_requirement": {
      "type": "string",
      "title": "Humidity requirement",
      "description": "Preferred habitat humidity or range",
      "x-group": "Traits"
    },
    "preferred_habitat": {
      "type": "string",
      "enum": [
        "aquatic",
        "terrestrial",
        "semi-aquatic",
        "arboreal",
        "desert",
        "tropical"
      ],
      "title": "Habitat type",
      "description": "Preferred habitat type for the amphibian",
      "x-group": "Traits"
    },
    "poisonous": {
      "type": "boolean",
      "title": "Poisonous",
      "description": "Whether amphibian is poisonous",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the amphibian is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the amphibian has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the amphibian is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Amphibian';


-- update Fish properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Gray",
        "Green",
        "Iridescent",
        "Orange",
        "Purple",
        "Red",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the fish",
      "x-group": "Traits"
    },
    "diet_type": {
      "type": "string",
      "enum": ["herbivore", "carnivore", "omnivore", "insectivore", "special"],
      "title": "Diet type",
      "description": "Main dietary requirement of the fish",
      "x-group": "Traits"
    },
    "humidity_requirement": {
      "type": "string",
      "title": "Humidity requirement",
      "description": "Preferred habitat humidity or range",
      "x-group": "Traits"
    },
    "preferred_habitat": {
      "type": "string",
      "enum": [
        "aquatic",
        "terrestrial",
        "semi-aquatic",
        "arboreal",
        "desert",
        "tropical"
      ],
      "title": "Habitat type",
      "description": "Preferred habitat type for fish",
      "x-group": "Traits"
    },
    "poisonous": {
      "type": "boolean",
      "title": "Poisonous",
      "description": "Whether amphibian is poisonous",
      "x-group": "Traits"
    },
    "freshwater": {
      "type": "boolean",
      "title": "Freshwater",
      "description": "Whether the fish lives in freshwater",
      "x-group": "Traits"
    },
    "saltawter": {
      "type": "boolean",
      "title": "Saltawter",
      "description": "Whether the fish lives in saltawter",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the fish is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the fish has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the fish is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Fish';


-- update Other Animal properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "description": {
      "type": "string",
      "title": "Information",
      "description": "Animal information",
      "x-group": "Traits"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the animal has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the animal is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Other Animal';


-- update Reptile properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Gray",
        "Green",
        "Iridescent",
        "Orange",
        "Purple",
        "Red",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the reptile",
      "x-group": "Traits"
    },
    "diet_type": {
      "type": "string",
      "enum": ["herbivore", "carnivore", "omnivore", "insectivore", "special"],
      "title": "Diet type",
      "description": "Main dietary requirement of the animal",
      "x-group": "Traits"
    },
    "humidity_requirement": {
      "type": "string",
      "title": "Humidity requirement",
      "description": "Preferred habitat humidity or range",
      "x-group": "Traits"
    },
    "preferred_habitat": {
      "type": "string",
      "enum": [
        "aquatic",
        "terrestrial",
        "semi-aquatic",
        "arboreal",
        "desert",
        "tropical"
      ],
      "title": "Habitat type",
      "description": "Preferred habitat type for reptile",
      "x-group": "Traits"
    },
    "poisonous": {
      "type": "boolean",
      "title": "Poisonous",
      "description": "Whether reptile is poisonous",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the reptile is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the reptile has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the reptile is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Reptile';


-- update Salamander / Newt properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Gray",
        "Green",
        "Iridescent",
        "Orange",
        "Purple",
        "Red",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the salamander/newt",
      "x-group": "Traits"
    },
    "diet_type": {
      "type": "string",
      "enum": ["herbivore", "carnivore", "omnivore", "insectivore", "special"],
      "title": "Diet type",
      "description": "Main dietary requirement of the animal",
      "x-group": "Traits"
    },
    "humidity_requirement": {
      "type": "string",
      "title": "Humidity requirement",
      "description": "Preferred habitat humidity or range",
      "x-group": "Traits"
    },
    "preferred_habitat": {
      "type": "string",
      "enum": [
        "aquatic",
        "terrestrial",
        "semi-aquatic",
        "arboreal",
        "desert",
        "tropical"
      ],
      "title": "Habitat type",
      "description": "Preferred habitat type for the salamander/newt",
      "x-group": "Traits"
    },
    "poisonous": {
      "type": "boolean",
      "title": "Poisonous",
      "description": "Whether salamander/newt is poisonous",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the salamander/newt is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the salamander/newt has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the salamander/newt is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Salamander / Newt';


-- update Snake properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Gray",
        "Green",
        "Iridescent",
        "Orange",
        "Purple",
        "Red",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the snake",
      "x-group": "Traits"
    },
    "humidity_requirement": {
      "type": "string",
      "title": "Humidity requirement",
      "description": "Preferred habitat humidity or range",
      "x-group": "Traits"
    },
    "preferred_habitat": {
      "type": "string",
      "enum": [
        "aquatic",
        "terrestrial",
        "semi-aquatic",
        "arboreal",
        "desert",
        "tropical"
      ],
      "title": "Habitat type",
      "description": "Preferred habitat type for the snake",
      "x-group": "Traits"
    },
    "venomous": {
      "type": "string",
      "title": "Venomous",
      "description": "Whether snake is venomous",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the snake is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the snake has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the snake is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Snake';


-- update Toad properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Gray",
        "Green",
        "Iridescent",
        "Orange",
        "Purple",
        "Red",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the animal",
      "x-group": "Traits"
    },
    "diet_type": {
      "type": "string",
      "enum": ["herbivore", "carnivore", "omnivore", "insectivore", "special"],
      "title": "Diet type",
      "description": "Main dietary requirement of the toad",
      "x-group": "Traits"
    },
    "humidity_requirement": {
      "type": "string",
      "title": "Humidity requirement",
      "description": "Preferred habitat humidity or range",
      "x-group": "Traits"
    },
    "preferred_habitat": {
      "type": "string",
      "enum": [
        "aquatic",
        "terrestrial",
        "semi-aquatic",
        "arboreal",
        "desert",
        "tropical"
      ],
      "title": "Habitat type",
      "description": "Preferred habitat type for the toad",
      "x-group": "Traits"
    },
    "poisonous": {
      "type": "boolean",
      "title": "Poisonous",
      "description": "Whether toad is poisonous",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the toad is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the toad has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the toad is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Toad';


-- update Tortoise properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Gray",
        "Green",
        "Iridescent",
        "Orange",
        "Purple",
        "Red",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the tortoise",
      "x-group": "Traits"
    },
    "diet_type": {
      "type": "string",
      "enum": ["herbivore", "carnivore", "omnivore", "insectivore", "special"],
      "title": "Diet type",
      "description": "Main dietary requirement of the animal",
      "x-group": "Traits"
    },
    "humidity_requirement": {
      "type": "string",
      "title": "Humidity requirement",
      "description": "Preferred habitat humidity or range",
      "x-group": "Traits"
    },
    "preferred_habitat": {
      "type": "string",
      "enum": [
        "aquatic",
        "terrestrial",
        "semi-aquatic",
        "arboreal",
        "desert",
        "tropical"
      ],
      "title": "Habitat type",
      "description": "Preferred habitat type for the tortoise",
      "x-group": "Traits"
    },
    "poisonous": {
      "type": "boolean",
      "title": "Poisonous",
      "description": "Whether tortoise is poisonous",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the tortoise is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the tortoise has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the tortoise is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Tortoise';


-- update Turtle properties schema
update "animal_specie"
set "properties_schema" = $$
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "color": {
      "type": "string",
      "enum": [
        "Black",
        "Blue",
        "Brown",
        "Gray",
        "Green",
        "Iridescent",
        "Orange",
        "Purple",
        "Red",
        "Tan",
        "White",
        "Yellow"
      ],
      "title": "Color",
      "description": "The color of the turtle",
      "x-group": "Traits"
    },
    "diet_type": {
      "type": "string",
      "enum": ["herbivore", "carnivore", "omnivore", "insectivore", "special"],
      "title": "Diet type",
      "description": "Main dietary requirement of the animal",
      "x-group": "Traits"
    },
    "humidity_requirement": {
      "type": "string",
      "title": "Humidity requirement",
      "description": "Preferred habitat humidity or range",
      "x-group": "Traits"
    },
    "preferred_habitat": {
      "type": "string",
      "enum": [
        "aquatic",
        "terrestrial",
        "semi-aquatic",
        "arboreal",
        "desert",
        "tropical"
      ],
      "title": "Habitat type",
      "description": "Preferred habitat type for the turtle",
      "x-group": "Traits"
    },
    "poisonous": {
      "type": "boolean",
      "title": "Poisonous",
      "description": "Whether turtle is poisonous",
      "x-group": "Traits"
    },
    "good_with_other_animals": {
      "type": "boolean",
      "title": "Good with other animals",
      "description": "Whether the turtle is good with other animals",
      "x-group": "Behaviour"
    },
    "special_needs": {
      "type": "boolean",
      "title": "Special needs",
      "description": "Whether the turtle has special needs",
      "x-group": "Health"
    },
    "vaccinated": {
      "type": "boolean",
      "title": "Vaccinated",
      "description": "Whether the turtle is vaccinated",
      "x-group": "Health"
    }
  }
}
$$::jsonb
where "name" = 'Turtle';
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
update "animal_specie" set "properties_schema" = null where "name" = 'Cat';
update "animal_specie" set "properties_schema" = null where "name" = 'Dog';
update "animal_specie" set "properties_schema" = null where "name" = 'Rabbit';
update "animal_specie" set "properties_schema" = null where "name" = 'Donkey';
update "animal_specie" set "properties_schema" = null where "name" = 'Horse';
update "animal_specie" set "properties_schema" = null where "name" = 'Mule';
update "animal_specie" set "properties_schema" = null where "name" = 'Pony';
update "animal_specie" set "properties_schema" = null where "name" = 'Miniature Horse';
update "animal_specie" set "properties_schema" = null where "name" = 'Button-Quail';
update "animal_specie" set "properties_schema" = null where "name" = 'Chicken';
update "animal_specie" set "properties_schema" = null where "name" = 'Dove';
update "animal_specie" set "properties_schema" = null where "name" = 'Duck';
update "animal_specie" set "properties_schema" = null where "name" = 'Emu';
update "animal_specie" set "properties_schema" = null where "name" = 'Finch';
update "animal_specie" set "properties_schema" = null where "name" = 'Goose';
update "animal_specie" set "properties_schema" = null where "name" = 'Guinea Fowl';
update "animal_specie" set "properties_schema" = null where "name" = 'Ostritch';
update "animal_specie" set "properties_schema" = null where "name" = 'Parakeet';
update "animal_specie" set "properties_schema" = null where "name" = 'Parrot';
update "animal_specie" set "properties_schema" = null where "name" = 'Peacock / Peafowl';
update "animal_specie" set "properties_schema" = null where "name" = 'Pheasant';
update "animal_specie" set "properties_schema" = null where "name" = 'Quail';
update "animal_specie" set "properties_schema" = null where "name" = 'Rhea';
update "animal_specie" set "properties_schema" = null where "name" = 'Swan';
update "animal_specie" set "properties_schema" = null where "name" = 'Toucan';
update "animal_specie" set "properties_schema" = null where "name" = 'Turkey';
update "animal_specie" set "properties_schema" = null where "name" = 'Alpaca';
update "animal_specie" set "properties_schema" = null where "name" = 'Cow';
update "animal_specie" set "properties_schema" = null where "name" = 'Goat';
update "animal_specie" set "properties_schema" = null where "name" = 'Llama';
update "animal_specie" set "properties_schema" = null where "name" = 'Pig';
update "animal_specie" set "properties_schema" = null where "name" = 'Pot Bellied';
update "animal_specie" set "properties_schema" = null where "name" = 'Sheep';
update "animal_specie" set "properties_schema" = null where "name" = 'Guinea Pig';
update "animal_specie" set "properties_schema" = null where "name" = 'Chinchilla';
update "animal_specie" set "properties_schema" = null where "name" = 'Rat';
update "animal_specie" set "properties_schema" = null where "name" = 'Mouse';
update "animal_specie" set "properties_schema" = null where "name" = 'Ferret';
update "animal_specie" set "properties_schema" = null where "name" = 'Hamster';
update "animal_specie" set "properties_schema" = null where "name" = 'Hedgehog';
update "animal_specie" set "properties_schema" = null where "name" = 'Gerbil';
update "animal_specie" set "properties_schema" = null where "name" = 'Degu';
update "animal_specie" set "properties_schema" = null where "name" = 'Prairie Dog';
update "animal_specie" set "properties_schema" = null where "name" = 'Skunk';
update "animal_specie" set "properties_schema" = null where "name" = 'Sugar Glider';
update "animal_specie" set "properties_schema" = null where "name" = 'Amphibian';
update "animal_specie" set "properties_schema" = null where "name" = 'Fish';
update "animal_specie" set "properties_schema" = null where "name" = 'Other Animal';
update "animal_specie" set "properties_schema" = null where "name" = 'Reptile';
update "animal_specie" set "properties_schema" = null where "name" = 'Salamander / Newt';
update "animal_specie" set "properties_schema" = null where "name" = 'Snake';
update "animal_specie" set "properties_schema" = null where "name" = 'Toad';
update "animal_specie" set "properties_schema" = null where "name" = 'Tortoise';
update "animal_specie" set "properties_schema" = null where "name" = 'Turtle';
-- +goose StatementEnd
