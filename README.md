### What is bob the builder? ###

Bob-the-builder adds a frontend and framework around your scripts. Whilst intended for automating builds, it can be used to assist the running of any UNIX script or program.

[Click here for a (GIF) Demo](https://s3-ap-southeast-2.amazonaws.com/ciphersink.net.current.workingfolder/Bob%20the%20builder.gif )




### What is this repository for? ###

* Stores code for the simple run automator 'bob the builder'
* Version 0.0.2

### How do I get set up? ###

* Install Go
* Get a copy of this repository
* 'go build'
* './bobthebuilder' or 'bobthebuilder.exe' (make sure config.json is in your working directory)
* Follow the steps below to setup your first build definition


### How do I setup something to be run automatically? ###

There are two components to every run definition. First, there is a configuration JSON file in the definitions/ folder. Second, there is as folder in base/, which should contain any initial files to be copied into the build folder before the build process commences.

All folders are relative to the working directory of bobthebuilder when invoked.

The JSON config file is simply a json file you whack in /definitions. It should be structured like this:


```json
{
  "name": "Build libc",
  "base-folder": "arm base",
  "git-src": "https://github.com/twitchyliquid64/bob-the-builder",
  "apt-packages-required": [
    "build-essentials",
    "screen",
    "htop"
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
    },
    {
      "type": "S3_UPLOAD",
      "bucket": "example.com.my.bucket",
      "region": "ap-northwest-2",
      "filename": "src/main.c",
      "filename-destination": "main.c",
      "ACL": "public"
    }
  ]
}
```

#### Special JSON definition fields

| Field Name    | Description                                   |
| ------------- |:----------------------------------------------|
| *name*          | Name of the run definition. Shows in the UI.  |
| *base-folder*   | Folder from which files are copied into the build folder during initialisation. This stage occurs after the clean phase and after (if any) the git clone phase.  |
| *git-src*       | This option should be set if you want to clone a git repository into your build area before running the steps. The passed value is passed to 'git clone <passed value> .' which downloads the contents of the repo into the build folder.  |
| *apt-packages-required*  | This option should be set if you want to ensure a set of apt-get packages are installed on your system. If you system does not support apt-get, do not set this field in the JSON file.  |



#### Available step types

| Type          | Description           | Parameters  |
| ------------- |:----------------------|       -----|
| *CMD*           | Runs the command with the specified arguments | <ul><li>'command' - name of the command to run. Do not put a path.</li><li>'args' - List of arguments to pass to the command. No escaping permitted.</li><li>'can-fail' - if true, the exit code of the command can be zero without failing the run or stopping it from progressing.</li> </ul>|
| *EXEC*           | Runs the script specified in 'command' using bash | <ul><li>'command' - Path to the script relative to the build directory.</li><li>'can-fail' - if true, the exit code of the command can be zero without failing the run or stopping it from progressing.</li></ul>|
| *S3_UPLOAD*           | Uploads and overwrites the specified file to AWS. AWS information must be populated in the configuration file. | <ul><li>'filename' - Path to the file relative to the build directory.</li><li>'region' - Name of the AWS region the bucket is in..</li><li>'bucket' - Name of the AWS bucket.</li><li>'filename-destination' - Path where the file is to be stored on the S3 bucket. If this parameter is empty or not provided, the file path of the source file will be used.</li><li>'ACL' - either 'public' or 'private'. This refers to the ACL applied on the object in S3.</li></ul>|



### Development Todos ###
 - Implement Artifactory
 - Implement version number system
 - Implement other step/phase types - upload to S3, set environment variable, etc
 - Definition reload button / refresh UI
