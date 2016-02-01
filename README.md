#1Password linux edition

![1Password_1](https://s3.amazonaws.com/nitecon/1Password/screenShots/1Password_1.png)
![1Password_2](https://s3.amazonaws.com/nitecon/1Password/screenShots/1Password_1.png)

## About 1Password
This application was started due to the need for a linux 1Password client.  There currently is no other packages that even compare closely to the functionality.  As seen in the screenshots above this does work however currently it requires compilation / manipulation of the main.go file to indicate where the 1password files live.
AND easy of use that 1Password provides.  I know there are many others like dashlane / lastpass / keypass etc and I've used them all, but none compares to 1password.
I Welcome contributions to this repository as with life sometimes I have time and sometimes I will have no time to work on this but I'm always appreciative of testers,
and merge requests.


## Getting started with the project

### Requirements for building this application:
- Golang (1.5+)
- NPM / Node

`go get -u https://github.com/Nitecon/1Password`

or if you would like to contribute, fork and clone it, (but make sure you check it out under the correct path in golang scope) fix / add functionality and PR back!

Once you have the repo cloned you should be able to immediately build it (if you run on a linux 64 bit host with golang installed + configured)
That will upgrade you to the revision that we currently use for the program.

## Why electron / some webkit garbage?
Answer to this is simple, it's easy to support multiple OS's which is obviously missing from the official clients.
Second and probably most importantly it will allow us an easy way to do non intrusuve updates as the frontend is
disconnected from the backend application that actually does the work etc.  The frontend in this case electron is only responsible for
spawning the application and eventually to check and do updates, downloads of the updates and respawn of the newly downloaded
go binary.  This greatly simplifies things and simplicity is very dear to me, so live with it.  Finally this allows us the ability
to build a cross OS app that can be deployed on windows / mac and perhaps in the future support more than basic functionality.

## How the application is structured / works and how I can contribute.
There is a couple of different ways you can help.
- Electron frontend to check a location for updates and downlaod / respawn
  - See the `app/main.js` file if you are interested in binary names / where it will be stored.
- If you are a go developer, think of this as a web app
- Which means we need both GUI / Backend help to interface with the 1Password api's etc.

More questions beyond reading the entry points (`app/main.js`) for electron and (`main.go`) for golang feel free to ask / create an issue.

## How does 1Password essentially work?
1Password design is based on a frontend / backend setup.  The frontend is handled by electron (chromium) as it already has window management etc solved and reduces our need to deal with native gui implementations across multiple OS's.  Node.js is responsible for spawning the OS specific golang binary application, and go into listening mode.  The go application upon startup responds to the electron application with a random port location that the application is listening to locally, and then redirects to that location for the actual "main" or homepage whatever you want to call it.  Beyond the redirection electron will have a couple of menu items over time for getting general application help / checking for updates, downloading updates / respawning the app and exit etc and thats just about it.
