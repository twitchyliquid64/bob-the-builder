

<div class="ui modal" id="miniDocumentationModal">
  <div class="header">
    Documentation
  </div>

  <div class="content" style="min-height: 450px;">
    <div class="ui fuild accordion" id="documentation-defedit" style="background-color: white;">
      <div class="title">
        <i class="dropdown icon"></i>
        What are the important fields?
      </div>
      <div class="content" style="">
        <p>At the root level of the definition JSON, you have:</p>
        <table class="ui celled table">
          <thead>
            <tr>
              <th>Field</th>
              <th>Details</th>
              <th></th>
            </tr>
        </thead>
          <tbody>
            <tr>
              <td>name</td>
              <td>Name of the run definition. Shows in the UI, and must be unique.</td>
              <td></td>
            </tr>
            <tr>
              <td>icon</td>
              <td>Icon to show in the UI for that build definition. See a list of icons <a href="http://semantic-ui.com/elements/icon.html">here</a>.</td>
              <td><div class="ui teal horizontal label">Optional</div></td>
            </tr>
            <tr>
              <td>base-folder</td>
              <td>If specified, files are copied from this location into the workspace prior to executing. This stage occurs after the clean phase and after (if any) the git clone phase.
                <br>
                The folder you wish to copy from must be located in build/. The string path you put in for this field must be relative to this directory.
              </td>
              <td><div class="ui teal horizontal label">Optional</div></td>
            </tr>
            <tr>
              <td>git-src</td>
              <td>This option should be set if you want to clone a git repository into your build area before running the steps. The passed value is passed to 'git clone .' which downloads the contents of the repo into the build folder.</td>
              <td><div class="ui teal horizontal label">Optional</div></td>
            </tr>
            <tr>
              <td>apt-packages-required</td>
              <td>This option should be set if you want to ensure a set of apt-get packages are installed on your system. If you system does not support apt-get, do not set this field in the JSON file.<br>This field should be a JSON list of package names.</td>
              <td><div class="ui teal horizontal label">Optional</div></td>
            </tr>
            <tr>
              <td>steps</td>
              <td>A JSON list of step structures. Each step describes an action the build to take, and is run in sequence.</td>
              <td></td>
            </tr>
            <tr>
              <td>params</td>
              <td>A JSON list of build parameters, which allow you to specific input data to the build.</td>
              <td><div class="ui teal horizontal label">Optional</div></td>
            </tr>
            <tr>
              <td>hide-from-log</td>
              <td>	If this boolean is set, the result of the execution will not be shown on the dashboard.</td>
              <td><div class="ui teal horizontal label">Optional</div></td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="title">
        <i class="dropdown icon"></i>
        What kinds of steps can I use?
      </div>
      <div class="content" style="">

        <table class="ui celled table">
          <thead>
            <tr>
              <th>Type</th>
              <th>Description</th>
              <th>Details</th>
            </tr>
        </thead>
          <tbody>
            <tr>
              <td>CMD</td>
              <td>Runs the given command with the specified arguments.</td>
              <td>
                <ul>
                  <li>'command' - name of the command to run. Do not put a path.</li>
                  <li>'args' <div class="ui olive horizontal label">template</div> - JSON List of arguments to pass to the command. No bash-style escaping permitted.</li>
                  <li>'can-fail' - if true, the exit code of the command can be zero without failing the run or stopping it from progressing.</li>
                </ul>
              </td>
            </tr>
            <tr>
              <td>EXEC</td>
              <td>Runs a bash script present in the build workspace.</td>
              <td>
                <ul>
                  <li>'command' - Path to the bash script.</li>
                </ul>
              </td>
            </tr>
            <tr>
              <td>S3_UPLOAD</td>
              <td>Uploads and overwrites the specified file to AWS. <br><br>
                <div class="ui orange horizontal label">NOTE</div> AWS information must be populated in the configuration file.</td>
              <td>
                <ul>
                  <li>'filename' <div class="ui olive horizontal label">template</div> - Path to the file relative to the build directory.</li>
                  <li>'region' - Name of the AWS region the bucket is in.</li>
                  <li>'bucket' - Name of the AWS bucket.</li>
                  <li>'filename-destination' <div class="ui olive horizontal label">template</div> - Path where the file is to be stored on the S3 bucket. If this parameter is empty or not provided, the file path of the source file will be used.</li>
                  <li>'ACL' - either 'public' or 'private'. This refers to the ACL applied on the object in S3.</li>
                </ul>
              </td>
            </tr>
            <tr>
              <td>ENV_SET</td>
              <td>
                Allows you to set environment variables for the build system, and any subsequent tasks.
              </td>
              <td>
                <ul>
                  <li>'key' <div class="ui olive horizontal label">template</div> - Name of the environment variable.</li>
                  <li>'value' <div class="ui olive horizontal label">template</div> - Value to set the environment variable to.</li>
                </ul>
              </td>
            </tr>
            <tr>
              <td>TAR_TO_S3</td>
              <td>
                Adds the given directories contents and the given files to a tar file, which is then compressed with gzip and streamed to S3. AWS information must be populated in the configuration file. This operation is suitable in low memory environments as the archive and compression routines are streamed on the given data.
              </td>
              <td>
                <ul>
                  <li>'region' - Name of the AWS region the bucket is in.</li>
                  <li>'bucket' - Name of the AWS bucket.</li>
                  <li>'filename-destination' <div class="ui olive horizontal label">template</div> - Path where the file is to be stored on the S3 bucket.</li>
                  <li>'directories' - JSON list of directories whoes files will be recursively added to the archive.</li>
                  <li>'files' - JSON list of files which will be added to the archive.</li>
                </ul>
              </td>
            </tr>
          </tbody>
        </table>

      </div>
      <div class="title">
        <i class="dropdown icon"></i>
        Can you give me some examples of steps?
      </div>
      <div class="content">
        <pre ng-non-bindable>
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
          },
          {
            "type": "TAR_TO_S3",
            "skip-condition": "{{getParameter `backupGit` | eq `true` | not}}",
            "filename-destination": "codebackup_{{.Day}}.{{.Month}}.{{.Year}}.tar.gz",
            "files": [
              "bob-the-builder.tar.gz",
              "misc-scripts.tar.gz"
            ],
            "bucket": "example.bucket",
            "region": "ap-southeast-2"
          },
        </pre>
      </div>
      <div class="title">
        <i class="dropdown icon"></i>
        What kind of parameters can I use?
      </div>
      <div class="content">
        <pre ng-non-bindable>
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
            "label": "Configuration",
            "varname": "filein",
            "filename": "main.out"
          }
        </pre>      </div>
      <div class="title">
        <i class="dropdown icon"></i>
        Tempate Snippets
      </div>
      <div class="content">
        <pre ng-non-bindable>
          // If backupGit parameter is not equal to 'true'
          {{getParameter `backupGit` | eq `true` | not}}

          // construct filename with the current date contained
          codebackup_{{.Day}}.{{.Month}}.{{.Year}}.tar.gz

          // returns true if the run has the given tag set
          {{ hasTag `tagName`  }}

          // returns all tags, space delimited
          {{ allTags }}
        </pre>
      </div>
    </div>
  </div>

  <div class="actions" style="text-align: initial; display: flex;">
    <div style="flex: 1; text-align: right;">
      <div class="ui button cancel">Close</div>
    </div>
  </div>
</div>









<div class="ui modal" id="runOptionsModal" ng-controller="runOptionsController">
  <div class="header">
    Run Options
  </div>

  <div class="content">

    <h4 class="ui dividing header">
      <i class="options icon"></i>
      <div class="content">
        Run parameters <i class="notched circle loading icon" ng-if="loadingParams"></i>
      </div>
    </h4>
    <div class="ui basic segment form"  style="min-height: 50px;">
      <div ng-repeat="param in defObj.params">

        <div ng-if="param.type == 'text'" class="field" style="margin-top: 7px;">
          <label>{{param.label}}</label>
          <input type="text" name="{{param.varname}}" ng-model="param.default" placeholder="{{param.placeholder}}" id="{{'runopt-field-' + $index}}">
        </div>


        <div ng-if="param.type == 'check'" class="field" style="margin-top: 7px;">
          <div class="ui checkbox">
            <input type="checkbox" tabindex="0" ng-model="param.default" id="{{'runopt-field-' + $index}}">
            <label>{{param.label}}</label>
          </div>
        </div>

        <div ng-if="param.type == 'branchselect'" style="margin-top: 7px;">
          <div class="ui fluid search selection dropdown" id="{{'runopt-field-' + $index}}">
            <input type="hidden" name="{{'nm-runopt-field-' + $index}}">
            <i class="dropdown icon"></i>
            <div class="default text">{{param.label}}</div>
            <div class="menu">
            </div>
          </div>
        </div>

        <div ng-if="param.type == 'select'" style="margin-top: 7px;">
          <div class="ui fluid search selection dropdown" id="{{'runopt-field-' + $index}}">
            <input type="hidden" name="{{'nm-runopt-field-' + $index}}">
            <i class="dropdown icon"></i>
            <div class="default text">{{param.label}}</div>
            <div class="menu">
              <div class="item" data-value="{{value}}" ng-repeat="(display, value) in param.items">
                {{display}}
              </div>
            </div>
          </div>
        </div>

        <div ng-if="param.type == 'file'" class="field" style="margin-top: 7px;">
          <label>{{param.label}}</label>
          <input type="file" name="{{param.varname}}" id="{{'runopt-field-' + $index}}" file-change="fileChange">
        </div>

      </div>
      <p ng-if="defObj.params == null || defObj.params == undefined || defObj.params.length == 0">No run parameters exist in the build definition.</p>
    </div>



    <h4 class="ui dividing header">
      <i class="tag icon"></i>
      <div class="content">
        Tags
      </div>
    </h4>

    <div class="ui basic segment" style="min-height: 50px;">
      <div class="ui fluid multiple search selection dropdown" style="border: 0px;" id="runOptionsModal-tagsDropdown">
        <input name="tags" type="hidden">
        <div class="default text"><i class="tags icon"></i> Add tags...</div>
        <div class="menu">
          <div class="item" data-value="production">
            <div class="ui red empty circular label"></div>
            Production
          </div>
          <div class="item" data-value="staging">
            <div class="ui yellow empty circular label"></div>
            Staging
          </div>
          <div class="item" data-value="test">
            <div class="ui green empty circular label"></div>
            Test
          </div>
          <div class="item" data-value="backup">
            <div class="ui purple empty circular label"></div>
            Backup
          </div>
          <div class="item" data-value="build">
            <div class="ui orange empty circular label"></div>
            Build
          </div>
          <div class="item" data-value="deploy">
            <div class="ui blue empty circular label"></div>
            Deploy
          </div>
          <div class="item" data-value="batch">
            <div class="ui teal empty circular label"></div>
            Batch
          </div>
          <div class="item" data-value="verify">
            <div class="ui pink empty circular label"></div>
            Verify
          </div>
          <div class="item" data-value="manual">
            <div class="ui cyan empty circular label"></div>
            Manual
          </div>
        </div>
      </div>
    </div>

    <h4 class="ui dividing header">
      <i class="adjust icon"></i>
      <div class="content">
        Auxillary
      </div>
    </h4>

    <div class="ui basic segment" style="min-height: 50px;">
      <div class="ui checkbox">
        <input type="checkbox" name="disphys" id="runOptionsModal-disablephys">
        <label>Disable physical indicators</label>
      </div>
      <div class="ui checkbox" style="margin-left: 3em;">
        <input type="checkbox" name="preventpostbuild">
        <label>Prevent post-build triggers</label>
      </div>
    </div>
  </div>

  <div class="actions" style="text-align: initial; display: flex;">

    <div style="flex: 1; text-align: left;">
      <div class="ui labeled input">
        <div class="ui label">
          version
        </div>
        <input type="text" name="version" placeholder="0.0.1" ng-model="defObj['last-version']" style="width: 78px;">
      </div>
    </div>


    <div style="flex: 1; text-align: right;">
      <div class="ui button cancel">Cancel</div>
      <div class="ui button primary icon ok">
        <i class="icon rocket"></i>
        Run
      </div>
    </div>
  </div>
</div>
