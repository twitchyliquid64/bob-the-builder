(function () {

    angular.module('baseApp')
        .controller('mainController', ['$scope', 'dataService', '$location', mainController]);

    function mainController($scope, dataService, $location) {
      var self = this;
      $scope.dataService = dataService;

      $scope.navBuild = function(index){
        console.log("Should be navigating to:", dataService.getDefinitions()[index]);
      }
    }


})();
