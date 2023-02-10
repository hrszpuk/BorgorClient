# ReCT Package System
The ultimate command line tool for managing ReCT Package System packages for the ReCT programming language.

## How to build/install 
In order to build and install `rps` you must have the Go programming language installed.
A guide to the installation of Go for your platform is available [here](https://go.dev/doc/install).

### Cloning the repository
To begin installing `rps`, we first must clone the project using git.
For this, you can use either git via command line, a git client, or even download the most recent source from the releases tab.

To clone using the git command line client you can use the commands below:
``` 
git clone git@github.com:RedCubeDev-ByteSpace/Rect-Package-System-Cli.git
```

Please open the project folder in the command line:
``` 
cd Rect-Package-System-Cli
```

### Building the project
Next, we need to build the project executable. 
This will allow you to run the project code on your computer.
The `main` branch of the repository may be unstable, so it is recommended you build from the latest tagged commit.

You can switch to the latest tag commit with the command below:
``` 
git checkout $(git describe --tags $(git rev-list --tags --max-count=1))
```

Then, assuming you have installed the Go programming language as mentioned in the requirements, please 
``` 
go build -v -o "rps"
```
You should now have an executable file in the project directory called borgor.

### System-wide installation
On Window, you can add this directly using software (please see guide [here](https://learn.microsoft.com/en-us/previous-versions/office/developer/sharepoint-2010/ee537574(v=office.14))).

On MacOS and Linux, you can either move the executable into a folder that's already on the path (such as `sudo mv rps /usr/bin`),
or you can add the file to path directly by adding `export PATH="$PATH:path/to/executable"` to your `~/.bash_profile` file (make sure to restart the shell).

## Contributors 
<a href="https://github.com/RedCubeDev-ByteSpace/Rect-Package-System-Cli/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=RedCubeDev-ByteSpace/Rect-Package-System-Cli" />
</a>
