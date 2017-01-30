[![Build Status](https://travis-ci.org/twitchyliquid64/bob-the-builder.svg?branch=master)](https://travis-ci.org/twitchyliquid64/bob-the-builder)

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
* Make sure the 'definitions' and 'bases' folder exist in your working directory.
* './bobthebuilder' or 'bobthebuilder.exe' (make sure config.json is in your working directory)
* Follow the steps below to setup your first build definition


### How do I setup something to be run automatically? ###

There are two components to every run definition. First, there is a configuration JSON file in the definitions/ folder. Second, there is as folder in base/, which should contain any initial files to be copied into the build folder before the build process commences.

All folders are relative to the working directory of bobthebuilder when invoked.

The JSON config file is simply a json file you whack in /definitions. It should be structured like this:


```json
{
  "name": "Build libc",
  "icon": "rocket",
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
      "skip-condition": "{{getParameter `include` | eq `true` | not}}",
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
  ],
  "params": [
    {
      "type": "check",
      "label": "Update browserify",
      "varname": "update",
      "default": false
    }
  ]
}
```

#### Special JSON definition fields

| Field Name    | Description                                   |
| ------------- |:----------------------------------------------|
| *name*          | Name of the run definition. Shows in the UI.  |
| *icon*          | Icon to show in the UI for that build definition. See a list of icons [here](http://semantic-ui.com/elements/icon.html)  |
| *base-folder*   | Folder from which files are copied into the build folder during initialisation. This stage occurs after the clean phase and after (if any) the git clone phase.  |
| *git-src*       | This option should be set if you want to clone a git repository into your build area before running the steps. The passed value is passed to 'git clone <passed value> .' which downloads the contents of the repo into the build folder.  |
| *apt-packages-required*  | This option should be set if you want to ensure a set of apt-get packages are installed on your system. If you system does not support apt-get, do not set this field in the JSON file.  |
| *params*  | List of configurable parameters which can be set to customize the build.  |



#### Available step types

| Type          | Description           | Parameters  |
| ------------- |:---------|       ----------------------------------------|
| *CMD*         | Runs the command with the specified arguments | <ul><li>'command' - name of the command to run. Do not put a path.</li><li>'args' <sup>template</sup> - List of arguments to pass to the command. No escaping permitted.</li><li>'can-fail' - if true, the exit code of the command can be zero without failing the run or stopping it from progressing.</li> </ul>|
| *EXEC*        | Runs the script specified in 'command' using bash | <ul><li>'command' - Path to the script relative to the build directory.</li><li>'can-fail' - if true, the exit code of the command can be zero without failing the run or stopping it from progressing.</li></ul>|
| *S3_UPLOAD*   | Uploads and overwrites the specified file to AWS. AWS information must be populated in the configuration file. | <ul><li>'filename' <sup>template</sup> - Path to the file relative to the build directory.</li><li>'region' - Name of the AWS region the bucket is in..</li><li>'bucket' - Name of the AWS bucket.</li><li>'filename-destination' <sup>template</sup> - Path where the file is to be stored on the S3 bucket. If this parameter is empty or not provided, the file path of the source file will be used.</li><li>'ACL' - either 'public' or 'private'. This refers to the ACL applied on the object in S3.</li></ul>|
| *ENV_SET*     | Allows you to set environment variables for the build system, and any subsequent tasks. | <ul><li>'key' <sup>template</sup> - Name of the environment variable</li><li>'value' <sup>template</sup> - Value to set the environment variable to.</li></ul>|
| *TAR_TO_S3*   | Adds the given directories contents and the given files to a tar file, which is then compressed with gzip and streamed to S3. AWS information must be populated in the configuration file. This operation is suitable in low memory environments as the archive and compression routines are streamed on the given data. | <ul><li>'region' - Name of the AWS region the bucket is in..</li><li>'bucket' - Name of the AWS bucket.</li><li>'filename-destination' <sup>template</sup> - Path where the file is to be stored on the S3 bucket.</li><li>'directories' - List of directories whoes files will be recursively added to the archive.</li><li>'files' - List of files which will be added to the archive.</li></ul>|

##### Additional step attributes

These fields can be set on any step.

| Name  | Description |
| ----- |:------------|
| HideFromSteps | This is a boolean field. If set to true, the step will not appear in the UI (top-right) when viewing the definition. |
| Conditional | This is a template field which allows you to define an expression, where it will skip evaluation of the step if the expression evaluates len(output) > 0 and output != 'false'. This field can be omitted - in which case the step will always evaluate. |

##### 'Template' fields

You may have noticed that a couple of the parameters in certain step types are 'templates'. This means they support Go's powerful templating system, which you can use to provide values dynamically. For instance, you can make the S3 uploader types prefix your files with todays date, or substitute in a tag name, or anything else that the go text/template engine supports. See [here](https://golang.org/pkg/text/template/) for details.

Additionally, the following functions are available for your templates:

 * hasTag(tagName) - returns True if a tag is set for the current run.
 * getParameter(variablename) - returns the value of the parameter, or '' if none exist. In the case of a checkbox, it returns a type boolean instead of a string.
 * allTags() - returns all tags, space delimited.

### Need more configurability? Use Parameters!

Parameters allow you to define a form so a user can populate information - just before your definitions run.

parameters are setup by adding new structures to the params list (in the JSON file).

All structures must have at minimum the following attributes:

 * label - Name of the field.
 * type - one of the field types specified below.
 * varname - name of the variable exported to your templates.

| Type           | Description                                                                                                          |
| -------------  |:---------------------------------------------------------------------------------------------------------------------|
| *check*        | adds a checkbox to the workflow. You may also specify a default which must be the string value true or false.        |
| *text*         | adds a text input to the workflow. You may also specify a placeholder.                                               |
| *select*       | adds a dropdown with configurable items.                                                                             |
| *branchselect* | adds a dropdown which is automatically populated with a list of branches in a remote git repository.                 |
| *file*         | adds a file upload field which will place the contents of the given file in a file on the workspace.                 |

#### Example structures

```json
{
  "type": "branchselect",
  "label": "Branch",
  "varname": "branch",
  "default": "master",
  "options": {
    "branchNamesOnly": true,
    "git-url": "https://github.com/twitchyliquid64/bob-the-builder"
  }
}

{
  "type": "select",
  "label": "Master Control",
  "varname": "mc",
  "items": {
    "Hello": "Hi"
  },
  "default": "Hi",
}

{
  "type": "text",
  "label": "Name",
  "varname": "stuff",
  "default": "robert",
}

},
{
  "type": "check",
  "label": "Backup build artifacts",
  "varname": "backup",
  "default": true
},
{
  "type": "file",
  "label": "Upload configuration",
  "varname": "file_input_filename",
  "filename": "input.json"
}
```

### Server config

You can see all the available settings defined in this struct [here](https://github.com/twitchyliquid64/bob-the-builder/blob/master/src/bobthebuilder/config/structure.go#L5).


### Todo


- [x] Add more valid tags
- [ ] Support parameter customizations for CRON jobs
- [ ] Support sending an email on success / failure
- [x] cronController to ingest crons-loaded event and refresh datasource.
