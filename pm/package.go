package pm

import (
	"encoding/json"
)

type Package struct {
	ID           int    `json:"-"`
	Name         string `json:"PackageName"`
	Version      string `json:"Version"`
	File         string `json:"-"`
	Dependencies string `json:"-"`
}

func JsonToPackage(data []byte) Package {
	var pack Package
	err := json.Unmarshal(data, &pack)

	if err != nil {
		dieErr("Failed to to parse package json information!", err)
	}

	return pack
}

func (p Package) EnsureInDB() {
	// try deleting the package, in case it exists already
	db.Exec("DELETE FROM Packages WHERE Name = ?", p.Name)

	// add this new package
	_, err := db.Exec("INSERT INTO Packages (Name, Version, File, Dependencies) VALUES (?,?,?,?)", p.Name, p.Version, p.File, p.Dependencies)
	ErrorDB(err)
}
