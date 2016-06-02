### TODOS ###
 - Run from user interface
 - Implement the different phase types
 - Collect results
 - WS for UI - producers/consumers architecture
 - Display results to user


### What is this repository for? ###

* Stores code for the simple run automator 'bob the builder'
* Version 0.0.1

### How do I get set up? ###

* Install Go
* Get a copy of this repository
* 'go build'
* './bobthebuilder' or 'bobthebuilder.exe' (make sure config.json is in your working directory)
* Follow the steps below to setup your first build definition


### How do I setup something to be built automatically? ###

There are two components to every build definition. First, there is a configuration JSON file in the definitions/ folder. Second, there is as folder in base/, which should contain any initial files to be copied into the build folder before the build process commences.

All folders are relative to the working directory of bobthebuilder when invoked.

The JSON config file is simply a json file you whack in /definitions. It should be structured like this:


```json
{
  "name": "Build libc",
  "base-folder": "arm base",
  "git-src": "https://github.com/twitchyliquid64/bob-the-builder",
  "apt-packages-required": [
    "build-essentials"
  ],
  "steps": [
    {
      "type": "CMD",
      "command": "ping",
      "args": [
        "-c",
        "2",
        "google.com"
      ],
      "can-fail": true
    },
    {
      "type": "CMD",
      "command": "git",
      "args": [
        "log"
      ]
    },
    {
      "type": "EXEC",
      "command": "build.sh",
      "can-fail": false
    }
  ]
}
```
