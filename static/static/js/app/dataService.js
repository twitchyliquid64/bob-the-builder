(function() {

    var app = angular.module('baseApp').factory('dataService', ['$rootScope', '$http', '$timeout', 'eventService', dataService]);

    GET_DEF_URL = "/api/definitions"
    GET_STATUS_URL = "/api/status"
    GET_HISTORY_URL = "/api/history"
    GET_STATUS_URL = "/api/status"


    function dataService($rootScope, $http, $timeout, eventService){
      var self = this;
      self.loading = true;//if true dimmer is shown
      self.loadPending = 0;//how many AJAX requests are pending

      self.connectionEstablished = false;
      self.connectionLost = false;//also set if API call = error.
      self.error = false;//set if there is an error

      self.buildDefinitions = [];
      self.history = [];
      self.status = {index: -1, run: null};





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
      self.queueRun = function(index){
        var defName = self.buildDefinitions[index].name;
        $http.get("/api/queue/new?name=" + defName, {}).then(function (response) {
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
