module github.com/jerrydevs/pokedex

go 1.21.5

replace pokeapi v0.0.0 => ./pokeapi

replace pokecache v0.0.0 => ./pokecache

require pokeapi v0.0.0

require pokecache v0.0.0 // indirect
