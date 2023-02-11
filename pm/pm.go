package pm

import (
	"io/ioutil"
	"net/http"
	"os"

	"borgor/print"
	"github.com/artdarek/go-unzip"
)

// Package PM
// ==========
// This package holds all package managing related functionality

func Get(op []string) {
	// Initialize the local package database
	InitializeDB()

	// check if a package name was given
	if len(op) != 1 {
		die("Get expected 1 argument but got %d!", len(op))
	}

	// remember the package name
	pack := op[0]

	// download this package
	DownloadPackage(pack)

	// done!
	print.PrintCF(print.Green, "Package '%s' has been installed!", pack)
	db.Close()
}

func DownloadPackage(pack string) {
	// look up the package
	print.PrintC(print.ThinYellow, "Looking up package on rps...")
	resp, err := http.Get("http://rps.rect.ml/api/lookup/" + pack)
	if err != nil {
		dieErr("Could not reach rps api endpoint!", err)
	}

	// package doesnt exist
	if resp.StatusCode == 404 {
		die("Could not find package '%s' on rps!*! Did you spell it correctly?", pack)
	}

	// something else went wrong
	if resp.StatusCode != 200 {
		die("Ouch!*! RPS gave back a status code of '%d', please tell an admin about this.", resp.StatusCode)
	}

	// if everything's alright
	// load the json response into a package object
	body, _ := ioutil.ReadAll(resp.Body)
	pdata := JsonToPackage(body)

	// close the response stream
	resp.Body.Close()

	print.PrintCF(print.ThinGreen, "Found package '%s', version: %s!", pdata.Name, pdata.Version)

	// set up a temp dir for the package to go
	SetUpTemp()

	// download this package from rps
	DownloadFile("./.tmp/pack.zip", "http://rps.rect.ml/api/download/"+pack)

	// create a directory for the package's contents
	err = os.Mkdir("./.tmp/pack", os.ModePerm)
	if err != nil {
		CleanUpTemp()
		die("Unable to create temporary directory for extraction!")
	}

	// unzip the package file
	uz := unzip.New("./.tmp/pack.zip", "./.tmp/pack")
	err = uz.Extract()

	if err != nil {
		CleanUpTemp()
		die("The package file could not be unzipped.")
	}

	// check if the package contains the .ll
	// if it doenst -> this package is broken
	if _, err := os.Stat("./.tmp/pack/" + pack + ".ll"); err != nil {
		CleanUpTemp()
		die("The package's module file (.ll) could not be located!*! This indicates a broken package.")
	}

	print.PrintC(print.ThinYellow, "Copying module file...")

	// if the module exists, copy it to the ./packages directory
	err = CopyFile("./.tmp/pack/"+pack+".ll", "./packages/"+pack+".ll")
	if err != nil {
		CleanUpTemp()
		die("The package could not be copied to ./packages!")
	}

	// remember the module file of this package
	pdata.File = "./packages/" + pack + ".ll"
	pdata.Dependencies = ""

	// remove the .ll from the temp dir
	os.Remove("./.tmp/pack/" + pack + ".ll")

	// check if this package has any dependencies
	empty, err := DirIsEmpty("./.tmp/pack/")

	// if so, copy them
	if !empty && err == nil {
		print.PrintC(print.ThinYellow, "Copying dependencies...")

		// if the dependency dir already exists, overwrite it
		err = os.RemoveAll("./packages/" + pack)
		if err != nil {
			CleanUpTemp()
			die("Could not create package dependency folder '%s'!", "./packages/"+pack)
		}

		// create a dependency dir
		err = os.Mkdir("./packages/"+pack, os.ModePerm)
		if err != nil {
			CleanUpTemp()
			die("Could not create package dependency folder '%s'!", "./packages/"+pack)
		}

		// copy all deps
		err = CopyDirectoryToDirectory("./.tmp/pack/", "./packages/"+pack+"/")
		if err != nil {
			CleanUpTemp()
			die("Could not copy package dependencies!")
		}

		// store the dep location in the package data object
		pdata.Dependencies = "./packages/" + pack

		print.PrintC(print.ThinGreen, "Dependencies have been copied!")
	} else {
		print.PrintC(print.ThinGreen, "Package does not need any dependencies!")
	}

	// make sure the package is registered in the db
	pdata.EnsureInDB()

	// clean up our temp directory
	CleanUpTemp()
}

func Update(op []string) {
	// Initialize the local package database
	InitializeDB()

	if len(op) == 0 {
		UpdateAll()
		return
	}

	// check if a package name was given
	if len(op) != 1 {
		die("Update expected 0 or 1 arguments but got %d!", len(op))
	}

	pack := op[0]

	// try finding the package in question
	print.PrintCF(print.ThinYellow, "Looking up package '%s'...", pack)
	res, err := db.Query("SELECT * FROM Packages WHERE Name = ?", pack)
	ErrorDB(err)

	// check if the package is in the database
	if !res.Next() {
		die("Package '%s' could not be found!*! Was it not installed with rps, or has already been removed?", pack)
	}

	// if it is, check if its up to date
	pdata := GetPackage(res)
	res.Close()

	// look up the package
	print.PrintC(print.ThinYellow, "Looking up package on rps...")
	resp, err := http.Get("http://rps.rect.ml/api/lookup/" + pack)
	if err != nil {
		dieErr("Could not reach rps api endpoint!", err)
	}

	// package doesnt exist
	if resp.StatusCode == 404 {
		die("Could not find package '%s' on rps!*! Did you spell it correctly?", pack)
	}

	// something else went wrong
	if resp.StatusCode != 200 {
		die("Ouch!*! RPS gave back a status code of '%d', please tell an admin about this.", resp.StatusCode)
	}

	// if everything's alright
	// load the json response into a package object
	body, _ := ioutil.ReadAll(resp.Body)
	prdata := JsonToPackage(body)

	// close the response stream
	resp.Body.Close()

	// check if we need to update at all
	if pdata.Version == prdata.Version {
		print.PrintCF(print.Green, "Package '%s' is already up to date!", pdata.Name)
		os.Exit(0)
	}

	// if its out of date -> install the new version
	print.PrintCF(print.Green, "New version of package '%s' was found! (%s | %s)", pdata.Name, pdata.Version, prdata.Version)
	print.PrintC(print.Green, "Installing new version...")

	// download the new version
	DownloadPackage(pack)

	// done!
	print.PrintCF(print.Green, "Package '%s' has been updated!", pack)
	db.Close()
}

func UpdateAll() {
	print.PrintCF(print.ThinYellow, "Checking what packages need updates...")

	// select all packages
	res, err := db.Query("SELECT * FROM Packages")
	ErrorDB(err)

	// keep track of all packages
	packs := make([]Package, 0)

	// load all the package data
	for res.Next() {
		packs = append(packs, GetPackage(res))
	}

	// end the query
	res.Close()

	// keep track of out of data packages
	ood := make([]Package, 0)

	for _, pkg := range packs {
		rpkg := SilentLookup(pkg.Name)

		if pkg.Version != rpkg.Version {
			ood = append(ood, pkg)
		}
	}

	// if theres no update, dont update lol
	if len(ood) == 0 {
		print.PrintC(print.Green, "All packages are up to date!")
		os.Exit(0)
	} else {
		print.PrintCF(print.Green, "Installing new versions of %d packages!", len(ood))
	}

	// update  t h e m   a l l
	for _, pkg := range ood {
		// download the new version
		DownloadPackage(pkg.Name)
	}

	// don
	print.PrintC(print.Green, "All packages have been updated!")
	db.Close()
}

func SilentLookup(pack string) *Package {
	resp, err := http.Get("http://rps.rect.ml/api/lookup/" + pack)
	if err != nil {
		dieErr("Could not reach rps api endpoint!", err)
	}

	// package doesnt exist
	if resp.StatusCode == 404 {
		return nil
	}

	// something else went wrong
	if resp.StatusCode != 200 {
		die("Ouch!*! RPS gave back a status code of '%d', please tell an admin about this.", resp.StatusCode)
	}

	// if everything's alright
	// load the json response into a package object
	body, _ := ioutil.ReadAll(resp.Body)
	prdata := JsonToPackage(body)

	// close the response stream
	resp.Body.Close()

	return &prdata
}

func Remove(op []string) {
	// Initialize the local package database
	InitializeDB()

	// check if a package name was given
	if len(op) != 1 {
		die("Remove expected 1 argument but got %d!", len(op))
	}

	pack := op[0]

	// try finding the package in question
	print.PrintCF(print.ThinYellow, "Looking up package '%s'...", pack)
	res, err := db.Query("SELECT * FROM Packages WHERE Name = ?", pack)
	ErrorDB(err)

	// check if the package is in the database
	if !res.Next() {
		die("Package '%s' could not be found!*! Was it not installed with rps, or has already been removed?", pack)
	}

	// if it is, remove it
	pdata := GetPackage(res)
	res.Close()

	// remove its .ll module
	print.PrintCF(print.ThinGreen, "Removing module file '%s'...", pdata.File)
	os.Remove(pdata.File)

	// remove its dependencies (if they exist)
	if pdata.Dependencies != "" {
		print.PrintCF(print.ThinGreen, "Removing dependencies '%s'...", pdata.Dependencies)
		os.RemoveAll(pdata.Dependencies)
	}

	// remove its database entry
	print.PrintC(print.ThinGreen, "Removing database entry...")
	_, err = db.Exec("DELETE FROM Packages WHERE ID = ?", pdata.ID)
	ErrorDB(err)

	print.PrintC(print.Green, "Package has been removed!")
	db.Close()
}
