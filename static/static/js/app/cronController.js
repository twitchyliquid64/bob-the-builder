(function () {

    angular.module('baseApp')
        .controller('cronController', ['$scope', 'dataService', '$location', '$routeParams', '$rootScope', '$http', cronController]);

    function cronController($scope, dataService, $location, $routeParams, $rootScope, $http) {
      var self = this;
      $scope.loading = true;
      $scope.dataService = dataService;
      $scope.dirty = false;

      $scope.crons = [];

      self.init = function(){
        for (var i = 0; i < $scope.crons.length; i++)
        {
          (function(){
            var it = i;
            $('#target-definition-select-' + i).dropdown({
              onChange: function(value, text, $selectedItem) {
                console.log(it, $scope.crons, value);
                $scope.crons[it].TargetDefinition = value;
                $scope.setDirty();
                $scope.$apply();
              }
            });
            $('#cron-tagsDropdown-' + it).dropdown({
              onChange: function(value, text, $selectedItem) {
                $scope.setDirty();
              }
            });
            $('#cron-tagsDropdown-' + it).dropdown('set selected', $scope.crons[it].Tags);
          })();
        }
      }

      $scope.disabled = function(){
        return !$scope.dirty;
      }

      $scope.setDirty = function(){
        $scope.dirty = true;
      }

      $scope.doSave = function(){
        for (var i = 0; i < $scope.crons.length; i++)
        {
          $scope.crons[i].Tags = $('#cron-tagsDropdown-' + i).dropdown('get value').split(",");
        }
        dataService.updateCronEntries($scope.crons)
      }

      $scope.delete = function(index){
        $scope.crons.splice(index, 1);
        setTimeout(function(){self.init();}, 200);
      }
      $scope.clone = function(index){
        $scope.crons[$scope.crons.length] = JSON.parse(JSON.stringify($scope.crons[index]));
        setTimeout(function(){self.init();}, 200);
      }
      $scope.add = function(){
        console.log($scope.crons);
        $scope.crons[$scope.crons.length] = {"Spec":"@every 5h","TargetDefinition":dataService.getDefinitions()[0].name,"Tags":["cron"]};
        $scope.setDirty();
        setTimeout(function(){self.init();}, 200);
      }

      $scope.$on('$routeChangeSuccess', function() { //apparently routeParams isnt always immediately populated
        if ($location.path().startsWith("/cron")){
          $('#cron-accordion').accordion();
          setTimeout(function(){self.init();}, 200);
          $scope.crons = dataService.getCrons();
        }
      });
    }

})();
