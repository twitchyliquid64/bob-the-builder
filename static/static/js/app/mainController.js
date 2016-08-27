(function () {

    angular.module('baseApp')
        .controller('mainController', ['$scope', 'dataService', '$location', '$routeParams', mainController]);

    function mainController($scope, dataService, $location, $routeParams) {
      var self = this;
      $scope.dataService = dataService;
      $scope.currentlyDash = $location.path() == "/";
      $scope.currentlyDocumentation = $location.path() == "/documentation";
      $scope.currentIndex = $routeParams.defID;
      $scope.currentlyBrowser = $location.path().startsWith("/browser/");

      $scope.$on('$routeChangeSuccess', function() { //apparently routeParams isnt always immediately populated
        $scope.currentIndex = $routeParams.defID;
      });

      $scope.navBuild = function(index){
        console.log("Should be navigating to:", dataService.getDefinitions()[index]);
        $location.path("/definition/" + index);
        $scope.currentlyDash = false;
        $scope.currentlyBrowser = false;
        $scope.currentlyDocumentation = false;
        $scope.currentIndex = index;
      }

      $scope.browser = function(){
        $location.path("/browser/");
        $scope.currentlyBrowser = true;
        $scope.currentlyDash = false;
        $scope.currentlyDocumentation = false;
      }

      $scope.navDashboard = function(){
        $location.path("/");
        $scope.currentlyDash = true;
        $scope.currentlyBrowser = false;
        $scope.currentlyDocumentation = false;
      }

      $scope.documentation = function(){
        $scope.currentlyDocumentation = true;
        $scope.currentlyBrowser = false;
        $scope.currentlyDash = false;
        $location.path("/documentation");
      }
    }


})();
