(function() {

    var app = angular.module('baseApp').factory('dataService', ['$rootScope', '$http', '$timeout', dataService]);

    GET_DEF_URL = "/api/definitions"


    function dataService($rootScope, $http, $timeout){
      var self = this;
      self.loading = false;//if true dimmer is shown
      self.loadPending = 0;//how many AJAX requests are pending

      self.connectionEstablished = false;
      self.connectionLost = false;//also set if API call = error.
      self.error = false;//set if there is an error

      self.buildDefinitions = [];


      //EXPOSED METHODS
      self.isLoading = function(){
        return self.loading;
      }
      self.getDefinitions = function(){
          return self.buildDefinitions;
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
      return self;
    };

})();
