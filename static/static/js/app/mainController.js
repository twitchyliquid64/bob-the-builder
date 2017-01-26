(function () {

    angular.module('baseApp')
        .controller('mainController', ['$scope', 'dataService', '$location', '$routeParams', mainController]);

    function mainController($scope, dataService, $location, $routeParams) {
      var self = this;
      $scope.dataService = dataService;

      $scope.$on('$routeChangeSuccess', function() { //apparently routeParams isnt always immediately populated
        $scope.currentIndex = $routeParams.defID;
      });

      $scope.navBuild = function(index){
        console.log("Should be navigating to:", dataService.getDefinitions()[index]);
        $location.path("/definition/" + index);
        $scope.currentIndex = index;
      }

      $scope.browser = function(){
        $location.path("/browser/");
      }
      $scope.cron = function(){
        $location.path("/cron");
      }

      $scope.navDashboard = function(){
        $location.path("/");
      }

      $scope.documentation = function(){
        $location.path("/documentation");
      }

      $scope.updateActiveStatus = function(){
        $scope.currentlyDash = $location.path() == "/";
        $scope.currentlyDocumentation = $location.path() == "/documentation";
        $scope.currentIndex = $routeParams.defID;
        $scope.currentlyBrowser = $location.path().startsWith("/browser/");
        $scope.currentlyCron = $location.path().startsWith("/cron");
      }
      $scope.updateActiveStatus();

      $scope.$on('$routeChangeSuccess', function(event, toState, toParams, fromState, fromParams){
        $scope.updateActiveStatus();
      })
    }


})();
