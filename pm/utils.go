package pm

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"borgor/print"
)

// die will print out a given error message and then die();
func die(text string, params ...interface{}) {
	msg := fmt.Sprintf(text, params...)

	if strings.Contains(msg, "*!") {
		parts := strings.Split(msg, "*!")
		print.WriteC(print.Red, parts[0])
		print.PrintC(print.ThinRed, parts[1])
	} else {
		print.PrintC(print.Red, msg)
	}

	os.Exit(-1)
}

// dieErr will print out a given error message + a golang error and then die();
func dieErr(text string, err error) {
	print.PrintC(print.Red, text)
	fmt.Println(err.Error())

	os.Exit(-1)
}

// DownloadFile will download a url to a local file.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func SetUpTemp() {
	os.Mkdir(dbDir+"/.tmp", os.ModePerm)
}

func CleanUpTemp() {
	os.RemoveAll(dbDir + "/.tmp")
}

func CopyFile(sourceFile string, destinationFile string) error {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)
		return err
	}

	return nil
}

func CopyDirectoryToDirectory(src string, dst string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	os.MkdirAll(dst, os.ModePerm)

	for _, file := range files {
		if file.IsDir() {
			err = CopyDirectoryToDirectory(src+"/"+file.Name(), dst+"/"+file.Name()+"/")
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(src+"/"+file.Name(), dst+"/"+file.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DirIsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
