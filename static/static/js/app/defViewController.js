(function () {

    angular.module('baseApp')
        .controller('defViewController', ['$scope', 'dataService', '$location', '$routeParams', '$rootScope', defViewController]);

    function defViewController($scope, dataService, $location, $routeParams, $rootScope) {
      var self = this;
      $scope.dataService = dataService;
      $scope.defObject = dataService.getDefinitions()[$routeParams.defID];
      $scope.buildQueued = false;
      $scope.running = false;
      $scope.phases = [];
      $scope.content = [];
      $scope.contentCursor = -1;//-1 = end of string as a special case

      $scope.isRunning = function(){
        return $scope.running;
      }
      $scope.isOtherDefinitionRunning = function(){
        var i = dataService.getStatus().index;
        return i >= 0 && i != $routeParams.defID;
      }
      $scope.parseDuration = function(milliseconds){
        if (milliseconds < (1000 * 50)){
          return "" + (Math.round(milliseconds / 100) / 10) + " seconds";
        }
        return moment.duration(milliseconds, "milliseconds").humanize();
      }


      $scope.edit = function(){
        console.log("edit()");
        $location.path("/edit/definition/" + $routeParams.defID + "/" + $scope.defObject.name);
      }
      $scope.editStep = function(step){
        if (step.type == "EXEC"){
          $location.path("/edit/file/" + $routeParams.defID + "/" + $scope.defObject['base-folder'] + "/" + step.command);
        } else {
          $location.path("/edit/definition/" + $routeParams.defID + "/" + $scope.defObject.name);
        }
      }

      $scope.run = function(){
        if($scope.buildQueued || $scope.running)return;//cant queue another one when one is queued or already running
        dataService.queueRun($routeParams.defID, "");
        $scope.buildQueued = true;
        self.lastRunWasWithOptions = false;
      }

      $scope.runOptions = function(){
        $rootScope.$broadcast('runOptionsModal-start', $scope.defObject);
      }
      self.runOptionsFormSubmitted = function(args){
        if($scope.buildQueued || $scope.running)return;
        if(args.defObj.name != $scope.defObject.name)return;

        //$scope.defObject = args.defObj; //is updating the defObject with new stuff such a good idea?
        dataService.queueRunWithOptions($routeParams.defID, {
          tags: args.tags,
          isPhysDisabled: args.isPhysDisabled,
          version: args.version,
          params: args.params
        });
        $scope.buildQueued = true;
        self.lastRunWasWithOptions = true;
      }



      //used for the 'steps' section - type corresponds to the set type in the JSON definition file.
      $scope.getStepTitle = function(type){
        if (type == 'CMD')return "Run command";
        if (type == 'EXEC')return "Run script";
        if (type == 'S3_UPLOAD')return "S3 file upload";
        if (type == 'S3_UPLOAD_FOLDER')return "S3 folder upload";
        if (type == 'ENV_SET')return "Set environment variable";
        if (type == 'TAR_TO_S3')return "Archive to S3";
        if (type == 'SEND_EMAIL')return "Send Email";
      }

      $scope.getCodeOutput = function(phase){
        if(phase.errorCode == 954321)return "Phase Skipped";
        return "Error Code: " + phase.errorCode;
      }

      self.stepTypeToIcons = {
        "CMD": {"terminal": true},
        "EXEC": {"rocket": true},
        "S3_UPLOAD": {"cloud upload": true},
        "S3_UPLOAD_FOLDER": {"cloud upload": true},
        "ENV_SET": {"level down": true},
        "TAR_TO_S3": {"archive": true},
        "SEND_EMAIL": {"send": true}
      }

      $scope.getStepIcons = function(stepType){
        return self.stepTypeToIcons[stepType];
      }

      //used for 'run output' section - phase.type corresponds to the value in .Type for the phase struct
      $scope.getPhaseTitle = function(phase){
        if (phase.type == "BASE-INSTALL"){
          return "Install base";
        } else if (phase.type == "CLEAN"){
          return "Clean workspace";
        } else if (phase.type == "COMMAND"){
          return "Run command: " + phase.Command;
        } else if (phase.type == "GIT-CLONE"){
          return "Git Clone"
        } else if (phase.type == "SCRIPT-EXEC"){
          return "Run script: " + phase.ScriptPath;
        } else if (phase.type == "APT-CHECK"){
          return "Check & install dependencies";
        } else if (phase.type == "S3UP_BASIC"){
          return "S3 file upload";
        } else if (phase.type == "S3_UPLOAD_FOLDER"){
          return "S3 folder upload";
        } else if (phase.type == "SET_ENV"){
          return "Set environment variable";
        } else if (phase.type == "TAR_TO_S3"){
          return "Archive to S3";
        } else if (phase.type == "SEND_EMAIL"){
          return "Send Email";
        }
      }
      $scope.getStepDetail = function(step){
        if (step.type == 'S3_UPLOAD'){
          return step.filename;
        }else if (step.type == 'ENV_SET'){
          return step.key;
        }else if (step.type == 'S3_UPLOAD_FOLDER'){
          return step.filename + (step['filename-destination'] ? (' --> ' + step['filename-destination']) : '');
        }else if (step.type == 'TAR_TO_S3'){
          step.files = step.files || [];
          step.directories = step.directories || [];
          return "" + step.files.length + " candidate file(s), " + step.directories.length + " candidate dirs";
        }else if (step.type == 'SEND_EMAIL'){
          return step.to ? step.to.join(',') : 'Sends an email to the default address.';
        }else {
          return step.command;
        }
      }









      //run when build definition run finishes.
      self.buildFinishedEvent = function(evt, args){
        //console.log("defViewController.buildFinishedEvent()");
        $scope.buildQueued = false;
        $scope.running = false;
        if (!self.lastRunWasWithOptions)dataService.requestUpdateDefinitionValuesFromServer();
      }
      self.buildStartedEvent = function(evt, args){
        //console.log("defViewController.buildStartedEvent()");
        $scope.running = true;
        $scope.phases = [];
        $scope.content = [];
      }
      self.phaseStartedEvent = function(args){
        console.log("defViewController.phaseStartedEvent(): ", args);
        $scope.phases[$scope.phases.length] = args.phase;
        $scope.content[$scope.phases.length-1] = "";
      }
      self.phaseFinishedEvent = function(args){
        //console.log("defViewController.phaseFinishedEvent(): ", args);
        $scope.phases[args.phase.index] = args.phase;
        //console.log($scope.phases, $scope.content);
      }
      self.phaseDataEvent = function(args){
        console.log("defViewController.phaseDataEvent(): ", args);
        //$scope.phases[args.phase.index] = args.phase;
        $scope.content[args.phase.index] += args.content

        controlIndex = $scope.content[args.phase.index].indexOf("CONTROL<CHAR-RETURN>");
        count = 0
        while (controlIndex > -1 && count < 8){ //has a /r - transform the content
          count += 1;
          //console.log("Control Index: " + controlIndex);
          relevantContent = $scope.content[args.phase.index].slice(0, controlIndex);
          //console.log("Relevant content: " + relevantContent);
          lastNewLine = relevantContent.lastIndexOf("\n");
          //console.log("Last New Line: " + lastNewLine);
          $scope.content[args.phase.index] = relevantContent.slice(0, lastNewLine+1) + $scope.content[args.phase.index].slice(controlIndex + "CONTROL<CHAR-RETURN>".length);
          //console.log("Content: " + $scope.content[args.phase.index]);
          controlIndex = $scope.content[args.phase.index].indexOf("CONTROL<CHAR-RETURN>");
          //console.log("Control Index: " + controlIndex);
        }
        //console.log($scope.phases, $scope.content);
      }



      //construct necessary event listeners.
      self.buildFinishListenerFactory = function(){
        return $rootScope.$on('ws-event-run-finish-'+$routeParams.defID, function(event, args) {
          self.buildFinishedEvent(args);
        });
      }
      self.buildStartListenerFactory = function(){
        return $rootScope.$on('ws-event-run-start-'+$routeParams.defID, function(event, args) {
          self.buildStartedEvent(args);
        });
      }
      self.buildPhaseStartListenerFactory = function(){
        return $rootScope.$on('ws-event-phase-started-'+$routeParams.defID, function(event, args) {
          self.phaseStartedEvent(args);
        });
      }
      self.buildPhaseFinishListenerFactory = function(){
        return $rootScope.$on('ws-event-phase-finished-'+$routeParams.defID, function(event, args) {
          self.phaseFinishedEvent(args);
        });
      }
      self.buildPhaseDataListenerFactory = function(){
        return $rootScope.$on('ws-event-phase-data-'+$routeParams.defID, function(event, args) {
          self.phaseDataEvent(args);
        });
      }
      self.runOptionsFormSubmittedListenerFactory = function(){
        return $rootScope.$on('runOptionsModal-finished', function(event, args) {
          self.runOptionsFormSubmitted(args);
        });
      }

      self.buildFinishListener = self.buildFinishListenerFactory();
      self.buildStartListener = self.buildStartListenerFactory();
      self.buildPhaseStartListener = self.buildPhaseStartListenerFactory();
      self.buildPhaseFinishListener = self.buildPhaseFinishListenerFactory();
      self.buildPhaseDataListener = self.buildPhaseDataListenerFactory();
      self.runOptionsFormSubmittedListener = self.runOptionsFormSubmittedListenerFactory();

      // event handlers
      $rootScope.$on('definitions-loaded', function(event, args) {
        $scope.defObject = dataService.getDefinitions()[$routeParams.defID];
      });

      $scope.$on('$routeChangeSuccess', function() { //apparently routeParams isnt always immediately populated
        $scope.defObject = dataService.getDefinitions()[$routeParams.defID];

        //destroy and remake event listeners based on our defID
        self.buildFinishListener();
        self.buildStartListener();
        self.buildPhaseStartListener();
        self.buildPhaseFinishListener();
        self.buildPhaseDataListener();
        self.runOptionsFormSubmittedListener();
        self.buildFinishListener = self.buildFinishListenerFactory();
        self.buildStartListener = self.buildStartListenerFactory();
        self.buildPhaseStartListener = self.buildPhaseStartListenerFactory();
        self.buildPhaseFinishListener = self.buildPhaseFinishListenerFactory();
        self.buildPhaseDataListener = self.buildPhaseDataListenerFactory();
        self.runOptionsFormSubmittedListener = self.runOptionsFormSubmittedListenerFactory();
      });

      $scope.$on('$destroy', function() {//destroy listeners
        self.buildFinishListener();
        self.buildStartListener();
        self.buildPhaseStartListener();
        self.buildPhaseFinishListener();
        self.buildPhaseDataListener();
        self.runOptionsFormSubmittedListener();
      });
    }


})();
