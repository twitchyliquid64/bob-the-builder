(function () {

    angular.module('baseApp')
        .controller('runOptionsController', ['$scope', 'dataService', '$location', '$routeParams', '$rootScope', runOptionsController]);

    function runOptionsController($scope, dataService, $location, $routeParams, $rootScope) {
      var self = this;
      $scope.defObj = {};

      self.initModal = function(){
        $('#runOptionsModal-tagsDropdown').dropdown({ allowAdditions: true, });

        $('#runOptionsModal').modal({//setup button callbacks + general parameters
          closable: false,
          autofocus: false,
          dimmerSettings: {
            closable: false,
            opacity: 0,
          },
          onApprove: self.submitPressed,
          onDeny: self.cancelPressed
        });
      }

      self.incrementVersion = function(){ //set the version to the last version plus one.
        console.log($scope.defObj);
        if ($scope.defObj['last-version'] == null || $scope.defObj['last-version'] == undefined || $scope.defObj['last-version'] == ""){
          $scope.defObj['last-version'] = "0.0.1";
        } else {
          $scope.defObj['last-version'] = $scope.defObj['last-version'].replace(/\d+$/, function(n){ return ++n });
        }
        console.log($scope.defObj);
      }

      self.setupAndShow = function(definitionObj){
        $scope.defObj = jQuery.extend({}, definitionObj);//shallow copy
        self.incrementVersion();
        self.initModal();
        $('#runOptionsModal').modal('show');
      }
      var startListener = $rootScope.$on('runOptionsModal-start', function(event, defObj) {
        self.setupAndShow(defObj);
      });
      
      $scope.$on('$destroy', function() {
        startListener();
      });

      self.submitPressed = function(){
        var tagArray = $('#runOptionsModal-tagsDropdown').dropdown('get value').split(",")
        var isPhysDisabled = $("#runOptionsModal-disablephys").prop('checked');
        var version = $scope.defObj['last-version'];
        $('#runOptionsModal').modal('hide');
        $rootScope.$broadcast('runOptionsModal-finished', {version: version, isPhysDisabled: isPhysDisabled, tags: tagArray, defObj: $scope.defObj});
      }
      self.cancelPressed = function(){
        $('#runOptionsModal').modal('hide');
      }
    }

})();
