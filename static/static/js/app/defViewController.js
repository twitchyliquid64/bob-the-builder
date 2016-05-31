(function () {

    angular.module('baseApp')
        .controller('defViewController', ['$scope', 'dataService', '$location', '$routeParams', '$rootScope', defViewController]);

    function defViewController($scope, dataService, $location, $routeParams, $rootScope) {
      var self = this;
      $scope.dataService = dataService;
      $scope.defObject = dataService.getDefinitions()[$routeParams.defID];
      $scope.buildQueued = false;
      $scope.running = false;

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

      $scope.getStepTitle = function(type){
        if (type == 'CMD')return "Run command";
        if (type == 'EXEC')return "Run script";
      }



      //run when build definition run finishes.
      self.buildFinishedEvent = function(evt, args){
        console.log("defViewController.buildFinishedEvent()");
        $scope.buildQueued = false;
        $scope.running = false;
      }
      self.buildStartedEvent = function(evt, args){
        console.log("defViewController.buildStartedEvent()");
        $scope.running = true;
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
      self.buildFinishListener = self.buildFinishListenerFactory();
      self.buildStartListener = self.buildStartListenerFactory();

      // event handlers
      $scope.$on('$routeChangeSuccess', function() { //apparently routeParams isnt always immediately populated
        $scope.defObject = dataService.getDefinitions()[$routeParams.defID];

        //destroy and remake event listeners based on our defID
        self.buildFinishListener();
        self.buildStartListener();
        self.buildFinishListener = self.buildFinishListenerFactory();
        self.buildStartListener = self.buildStartListenerFactory();
      });
    }


})();
