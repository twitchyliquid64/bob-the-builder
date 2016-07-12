(function () {

    angular.module('baseApp')
        .controller('runOptionsController', ['$scope', 'dataService', '$location', '$routeParams', '$rootScope', runOptionsController]);

    function runOptionsController($scope, dataService, $location, $routeParams, $rootScope) {
      var self = this;
      $scope.buildParams = [];

      self.setParamsEvent = function(args){
        console.log("runOptionsController.setParamsEvent(): ", args);
        $scope.buildParams = args;
      }


      $rootScope.$on('runOptionsModal-setParamsEvent', function(event, args) {
        self.setParamsEvent(args);
      });
    }

})();
