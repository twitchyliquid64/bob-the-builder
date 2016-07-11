(function() {

    var app = angular.module('baseApp').factory('dataService', ['$rootScope', '$http', '$timeout', 'eventService', '$window', dataService]);

    GET_DEF_URL = "/api/definitions"
    GET_STATUS_URL = "/api/status"
    GET_HISTORY_URL = "/api/history"
    GET_STATUS_URL = "/api/status"


    function dataService($rootScope, $http, $timeout, eventService, $window){
      var self = this;
      self.loading = true;//if true dimmer is shown
      self.loadPending = 0;//how many AJAX requests are pending

      self.connectionEstablished = false;
      self.connectionLost = false;//also set if API call = error.
      self.error = false;//set if there is an error

      self.buildDefinitions = [];
      self.history = [];
      self.status = {index: -1, run: null};
      self.loadingMessage = "Loading build data";
      self.reloadQueued = false;
      self.serverStats = {prettyMemUsage: "-"};




      //EXPOSED METHODS
      self.isLoading = function(){
        return self.loading;
      }
      self.getDefinitions = function(){
        return self.buildDefinitions;
      }
      self.getHistory = function(){
        return self.history;
      }
      self.getStatus = function(){
        return self.status;
      }
      self.queueRun = function(index, version){
        if (  self.reloadQueued)return;

        if (version == null || version == "")version = "0.0.1";

        var defName = self.buildDefinitions[index].name;
        $http.get("/api/queue/new?version=" + version + "&name=" + defName, {}).then(function (response) {
        }, function errorCallback(response) {
          console.log(response);
          self._error();
        });
      }
      self.queueRunWithOptions = function(index, options){
        console.log("queueRunWithOptions: ", options);
        if (  self.reloadQueued)return;

        options.name = self.buildDefinitions[index].name;
        $http.post("/api/queue/newWithOptions", options).then(function (response) {
        }, function errorCallback(response) {
          console.log(response);
          self._error();
        });
      }
      //END EXPOSED METHODS





      self._loadDefinitions = function(){//called at end of factory (init)
        self._incrementLoadPendingCounter();
        $http.get(GET_DEF_URL, {}).then(function (response) {
          self.buildDefinitions = response.data;
          $rootScope.$broadcast('definitions-loaded');
          self._decrementLoadPendingCounter();
        }, function errorCallback(response) {
          console.log(response);
          self._decrementLoadPendingCounter();
          self._error();
        });
      }

      self._loadHistory = function(){//called at end of factory (init)
        self._incrementLoadPendingCounter();
        $http.get(GET_HISTORY_URL, {}).then(function (response) {
          self.history = response.data;
          self._preprocessHistory();
          self._decrementLoadPendingCounter();
        }, function errorCallback(response) {
          console.log(response);
          self._decrementLoadPendingCounter();
          self._error();
        });
      }

      self._loadStatus = function(){//called at end of factory (init)
        self._incrementLoadPendingCounter();
        $http.get(GET_STATUS_URL, {}).then(function (response) {
          self.status = response.data;
          self._decrementLoadPendingCounter();
        }, function errorCallback(response) {
          console.log(response);
          self._decrementLoadPendingCounter();
          self._error();
        });
      }

      self._preprocessHistory = function(){
        for (var i = 0; i < self.history.length; i++) {
          self.history[i].startMom = moment(self.history[i].startTime);
        };
        self.history.reverse();
      };



      //handle events from the eventService
      $rootScope.$on('ws-events-connected', function(event, args) {
        self.connectionEstablished = true;
      });
      $rootScope.$on('ws-events-closed', function(event, args) {
        self.connectionLost = true;
      });
      $rootScope.$on('ws-events-run-started', function(event, statusObj){
        self.status = statusObj;
      });
      $rootScope.$on('ws-events-run-finished', function(event, statusObj){
        self.status = statusObj;
        self.status.index = -1;
        self._loadHistory();
      });
      $rootScope.$on('ws-events-reload-started', function(event, statusObj){
        self.loadingMessage = "Definitions reload in progress"
        self._incrementLoadPendingCounter();
        self.loading = true;
      });
      $rootScope.$on('ws-events-reload-finished', function(event, statusObj){
        self.loadingMessage = "Definitions reload completed"
        $window.location.reload(true);
      });
      $rootScope.$on('ws-events-reload-queued', function(event, statusObj){
        self.loadingMessage = "Definitions reload queued, waiting to run"
        self.reloadQueued = true;
      });
      $rootScope.$on('ws-server-stats', function(event, sStats){
        self.serverStats = sStats;
        self.serverStats.prettyMemUsage = Math.round(sStats.mem.ActualFree / 1024 / 1024);
        if (self.serverStats.prettyMemUsage >= 1000){
          self.serverStats.prettyMemUsage = (self.serverStats.prettyMemUsage / 1024).toFixed(1) + " GB";
        } else {
          self.serverStats.prettyMemUsage = self.serverStats.prettyMemUsage =  + " MB";
        }
      });


      self.requestDefinitionsReload = function(){
        self.loadingMessage = "Queuing reload operation"
        self.loading = true;
        self._incrementLoadPendingCounter();
        $http.get("/api/definitions/reload", {}).then(function (response) {
        }, function errorCallback(response) {
          console.log(response);
          self.loadingMessage = "Queuing failed - please refresh your browser."
          self._error();
        });
      }


      self._incrementLoadPendingCounter = function(){
        //self.loading = true;
        self.loadPending += 1;
      }
      self._decrementLoadPendingCounter = function(){
        self.loadPending -= 1;
        if (self.loadPending == 0){
          $timeout(function(){self.loading = false;}, 400);
        }
      }
      self._error = function(){
        self.error = true;
        self.connectionLost = true;
      }




      self._loadDefinitions();
      self._loadHistory();
      self._loadStatus();
      return self;
    };

})();
