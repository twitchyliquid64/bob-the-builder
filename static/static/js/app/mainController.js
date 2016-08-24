(function () {

    angular.module('baseApp')
        .controller('mainController', ['$scope', 'dataService', '$location', '$routeParams', mainController]);

    function mainController($scope, dataService, $location, $routeParams) {
      var self = this;
      $scope.dataService = dataService;
      $scope.currentlyDash = $location.path() == "/";
      $scope.currentlyDocumentation = $location.path() == ""
      $scope.currentIndex = $routeParams.defID;

      $scope.$on('$routeChangeSuccess', function() { //apparently routeParams isnt always immediately populated
        $scope.currentIndex = $routeParams.defID;
      });

      $scope.navBuild = function(index){
        console.log("Should be navigating to:", dataService.getDefinitions()[index]);
        $location.path("/definition/" + index);
        $scope.currentlyDash = false;
        $scope.currentIndex = index;
      }

      $scope.navDashboard = function(){
        $location.path("/");
        $scope.currentlyDash = true;
      }

      $scope.documentation = function(){
        $scope.currentlyDocumentation = true;
        $location.path("/documentation");
      }
    }


})();
