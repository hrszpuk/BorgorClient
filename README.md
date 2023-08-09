# :hamburger: Borgor - RPS Client (fork)
The ultimate command line tool for managing ReCT Package System packages for the ReCT programming language.

## :warning: This project has been archived :warning:
Borgor not longer works as the development of the rect programming language has ended and the maintanence of rps.rect-lang.org (or rps-dev.rect-lang.org) has come to an end.
For this reason, I have decided the project will be publicly archived to indicate the fact that borgor not longer works and will not longer be maintained.

## How to build/install 
In order to build and install `borgor` you must have the Go programming language installed.
A guide to the installation of Go for your platform is available [here](https://go.dev/doc/install).

### Cloning the repository
To begin installing `borgor`, we first must clone the project using git.
For this, you can use either git via command line, a git client, or even download the most recent source from the releases tab.

To clone using the git command line client you can use the commands below:
``` 
git clone git@github.com:hrszpuk/BorgorClient.git
```

Please open the project folder in the command line:
``` 
cd BorgorClient
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
go build -v -o borgor
```
You should now have an executable file in the project directory called borgor.

### System-wide installation
On Windows, you can add this directly using software (please see guide [here](https://learn.microsoft.com/en-us/previous-versions/office/developer/sharepoint-2010/ee537574(v=office.14))).

On MacOS and Linux, you can either move the executable into a folder that's already on the path (such as `sudo mv borgor /usr/bin`),
or you can add the file to path directly by adding `export PATH="$PATH:path/to/executable"` to your `~/.bash_profile` file (make sure to restart the shell).

## Contributors 
<a href="https://github.com/hrszpuk/BorgorClient/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=hrszpuk/BorgorClient" />
</a>
