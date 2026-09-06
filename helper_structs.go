package main

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Name represents a localized name
type Name struct {
	Name     string           `json:"name"`
	Language NamedAPIResource `json:"language"`
}

// EncounterMethodRate represents encounter method rates in a location area
type EncounterMethodRate struct {
	EncounterMethod NamedAPIResource `json:"encounter_method"`
	VersionDetails  []struct {
		Rate    int              `json:"rate"`
		Version NamedAPIResource `json:"version"`
	} `json:"version_details"`
}

// PokemonEncounter represents Pokémon encounter details
type PokemonEncounter struct {
	Pokemon        NamedAPIResource `json:"pokemon"`
	VersionDetails []struct {
		Version          NamedAPIResource `json:"version"`
		MaxChance        int              `json:"max_chance"`
		EncounterDetails []struct {
			MinLevel        int                `json:"min_level"`
			MaxLevel        int                `json:"max_level"`
			ConditionValues []NamedAPIResource `json:"condition_values"`
			Chance          int                `json:"chance"`
			Method          NamedAPIResource   `json:"method"`
		} `json:"encounter_details"`
	} `json:"version_details"`
}

// LocationArea represents a sub-area within a location
type LocationArea struct {
	ID                   int64                 `json:"id"`
	Name                 string                `json:"name"`
	GameIndex            int64                 `json:"game_index"`
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	Location             NamedAPIResource      `json:"location"`
	Names                []Name                `json:"names"`
	PokemonEncounters    []PokemonEncounter    `json:"pokemon_encounters"`
}

// LocationAreaList represents the paginated location-area collection returned by the API.
type LocationAreaList struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous *string            `json:"previous"`
	Results  []NamedAPIResource `json:"results"`
}
