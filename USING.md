====== Getting a development environment ======
===== Install programs =====
==== Windows (10) ====
<del>
Just install it on Windows right?
</del>
  * Make sure latest updates are applied
  * Type "Turn windows features on or off" into the search bar and hit enter
  * Tick the "Enable windows subsystem for linux" box
  * Restart your computer
  * Go to the windows store and install Ubuntu
  * Follow the prompts to create a user etc.
  * Run ''sudo apt-get install update-manager-core''
  * Run ''sudo nano /etc/update-manager/release-upgrades''
  * Change the line that starts with ''Prompt'' to ''Prompt=normal''
  * Run ''sudo do-release-upgrade''
  * Follow on-screen prompts

  * Follow Linux instructions

==== Ubuntu ====
//Probably applies to most linux distros, with subtle changes//
  * Sassc: ''sudo apt-get install sassc''
  * Git: ''sudo apt-get install git''
  * Make: ''sudo apt-get install make''
  * Go: ''sudo apt-get install golang-go''
    * As long as you have a vaguely up to date Ubuntu (16.04+) the version of Go that you can get by default will be fine, but just in case follow instructions on https://github.com/golang/go/wiki/Ubuntu, add the ppa, then apt-get install required version (Check by running ''go version'')
=== Golang manual install ===
    * Go to the [[https://golang.org/dl/|Go downloads page]] and scroll to version 1.6
    * Download ''go1.6.linux-amd64.tar.gz''
      * To do this in terminal, visit the webpage in another browser, 
      * copy the address of the link
      * Run ''wget [address]''
    * Open a terminal at its location and unzip the file to /usr/local: ''sudo tar -C /usr/local -xzf go1.6.linux-amd64.tar.gz''
    * Now add the following lines to the end of ''~/.bashrc'' to add Go to your PATH:\\ ''export PATH=$PATH:/usr/local/go/bin''\\ ''export GOPATH=~/go''
    * Run ''source ~/.bashrc'' to update PATH values

===== Pulling from git =====
  * Make sure GOPATH is set (''echo $GOPATH'')
  * ''go get github.com/UniversityRadioYork/2016-site''
    * This downloads 2016-site to ~/go/src/github.com/UniversityRadioYork/2016-site along with all necessary dependencies
  * To update your working copy of the repo, running ''go get'' in the 2016-site directory, but you can also checkout to a specific branch if required:
    * ''git checkout this-is-a-branch-name'' 

===== Running the server =====
  * Copy ''config.toml.example'' to ''config.toml''
  * Create a new file called ''.myradio.key'' and paste your [[https://ury.org.uk/ceedox/computing:software:in-house:myury:api#getting_a_key|myradio api-key]] into it
  * Run ''make run''
  * Assuming there are no errors, go to http://localhost:3000 in your browser and there should be a website

===== Editing files =====
=== Ubuntu ===
Trivial. Open a file editor and browse to ~/go/src/github.com/UniversityRadioYork/2016-site

=== Windows ===
  * Windows likes to hide the actual files for the Ubuntu subsystem.
  * Depending on the version of the subsystem installed, they'll be located somewhere inside Users/<User>/AppData/ . See this AskUbuntu post for more details: https://askubuntu.com/q/759880
  * Probably a good idea to make a shortcut for it somewhere in your Documents folder or similar
  * After making an edit, re-run the website by pressing ctrl^c in the linux window, then type ''make run'' again
