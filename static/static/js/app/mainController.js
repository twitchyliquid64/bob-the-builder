(function () {

    angular.module('baseApp')
        .controller('mainController', ['$scope', 'dataService', '$location', mainController]);

    function mainController($scope, dataService, $location) {
      var self = this;
      $scope.dataService = dataService;
    }


})();
