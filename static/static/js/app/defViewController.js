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
      $scope.run = function(){
        if($scope.buildQueued || $scope.running)return;//cant queue another one when one is queued or already running
        dataService.queueRun($routeParams.defID);
        $scope.buildQueued = true;
      }

      //used for the 'steps' section - type corresponds to the set type in the JSON definition file.
      $scope.getStepTitle = function(type){
        if (type == 'CMD')return "Run command";
        if (type == 'EXEC')return "Run script";
        if (type == 'S3_UPLOAD')return "S3 file upload";
        if (type == 'ENV_SET')return "Set environment variable";
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
        } else if (phase.type == "SET_ENV"){
          return "Set environment variable";
        }
      }
      $scope.getStepDetail = function(step){
        if (step.type == 'S3_UPLOAD'){
          return step.filename;
        }else if (step.type == 'ENV_SET'){
          return step.key;
        }else {
          return step.command;
        }
      }



      //run when build definition run finishes.
      self.buildFinishedEvent = function(evt, args){
        //console.log("defViewController.buildFinishedEvent()");
        $scope.buildQueued = false;
        $scope.running = false;
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
        //console.log("defViewController.phaseDataEvent(): ", args);
        //$scope.phases[args.phase.index] = args.phase;
        if ($scope.contentCursor != -1){
          $scope.content[args.phase.index] = $scope.content[args.phase.index].slice(0, $scope.contentCursor);
          $scope.contentCursor = -1;
        }

        if (args.content.indexOf("CONTROL<CHAR-RETURN>") > -1){ //has a /r - transform the current content
          $scope.contentCursor = $scope.content[args.phase.index].lastIndexOf("\n");
          if ($scope.contentCursor == -1)$scope.contentCursor = 0;
        }
        $scope.content[args.phase.index] += args.content.replace("CONTROL<CHAR-RETURN>", "");
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

      self.buildFinishListener = self.buildFinishListenerFactory();
      self.buildStartListener = self.buildStartListenerFactory();
      self.buildPhaseStartListener = self.buildPhaseStartListenerFactory();
      self.buildPhaseFinishListener = self.buildPhaseFinishListenerFactory();
      self.buildPhaseDataListener = self.buildPhaseDataListenerFactory();




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
        self.buildFinishListener = self.buildFinishListenerFactory();
        self.buildStartListener = self.buildStartListenerFactory();
        self.buildPhaseStartListener = self.buildPhaseStartListenerFactory();
        self.buildPhaseFinishListener = self.buildPhaseFinishListenerFactory();
        self.buildPhaseDataListener = self.buildPhaseDataListenerFactory();
      });
    }


})();
