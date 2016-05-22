#TODOS#
 - Everything - All I have currently is boilerplate.


# README #

This README would normally document whatever steps are necessary to get your application up and running.

### What is this repository for? ###

* Stores code for the simple build automator 'bob the builder'
* Version 0.0.1

### How do I get set up? ###

* Install Go
* Get a copy of this repository
* 'go build'
* './bobthebuilder' or 'bobthebuilder.exe' (make sure testconfig.json is in your working directory)
* Follow the steps below to setup your first build definition


### How do I setup something to be built automatically? ###

There are two components to every build definition. First, there is a configuration JSON file in the definitions/ folder. Second, there is as folder in base/, which should contain any initial files to be copied into the build folder before the build process commences.

All folders are relative to the working directory of bobthebuilder when invoked.

The JSON config file is simple a json file you whack in /definitions. It should be structured like this:


```json
{
  "name": "Build libc",
  "base-folder": "libc-base",
  "apt-packages-required": [
    "build-essentials"
  ],
  "steps": [
    {
      "command": "./build.sh",
      "can-fail": false
    }
  ]
}
```
