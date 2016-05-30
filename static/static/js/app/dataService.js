(function() {

    var app = angular.module('baseApp').factory('dataService', ['$rootScope', '$http', '$timeout', dataService]);

    GET_DEF_URL = "/api/definitions"
    GET_STATUS_URL = "/api/status"
    GET_HISTORY_URL = "/api/history"


    function dataService($rootScope, $http, $timeout){
      var self = this;
      self.loading = false;//if true dimmer is shown
      self.loadPending = 0;//how many AJAX requests are pending

      self.connectionEstablished = false;
      self.connectionLost = false;//also set if API call = error.
      self.error = false;//set if there is an error

      self.buildDefinitions = [];
      self.history = [];





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

      self._preprocessHistory = function(){
        for (var i = 0; i < self.history.length; i++) {
          self.history[i].startMom = moment(self.history[i].startTime);
        };
      };






      self._incrementLoadPendingCounter = function(){
        self.loading = true;
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
      return self;
    };

})();
