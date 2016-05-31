(function () {

    angular.module('baseApp')
        .controller('defViewController', ['$scope', 'dataService', '$location', '$routeParams', defViewController]);

    function defViewController($scope, dataService, $location, $routeParams) {
      var self = this;
      $scope.dataService = dataService;
      $scope.defObject = dataService.getDefinitions()[$routeParams.defID];

      $scope.$on('$routeChangeSuccess', function() { //apparently routeParams isnt always immediately populated
        $scope.defObject = dataService.getDefinitions()[$routeParams.defID];
        console.log($scope.defObject);
      });

      $scope.getStepTitle = function(type){
        if (type == 'CMD')return "Run command";
        if (type == 'EXEC')return "Run script";
      }
    }


})();
