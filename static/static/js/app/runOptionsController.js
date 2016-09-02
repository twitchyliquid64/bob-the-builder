(function () {

    angular.module('baseApp')
        .controller('runOptionsController', ['$scope', 'dataService', '$location', '$routeParams', '$rootScope', runOptionsController]);

    function runOptionsController($scope, dataService, $location, $routeParams, $rootScope) {
      var self = this;
      $scope.defObj = {};
      $scope.loadingParams = true;

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
        $scope.loadingParams = true;
        setTimeout(function(){self.initFields();}, 200);//to be sure digest cycles have completed
        $('#runOptionsModal').modal('show');
      }
      var startListener = $rootScope.$on('runOptionsModal-start', function(event, defObj) {
        self.setupAndShow(defObj);
      });

      $scope.$on('$destroy', function() {
        startListener();
      });


      //called just before modal show to initialize any javascript on field views.
      self.initFields = function(){
        if ($scope.defObj.params){
          for (var i = 0; i < $scope.defObj.params.length; i++)
          {
            if ($scope.defObj.params[i].type == "branchselect"){
              $('#runopt-field-' + i).dropdown({
                apiSettings: {
                  url: '/api/lookup/buildparam?name=' + $scope.defObj.name + '&param=' + i,
                  onResponse: function(remoteResponse) {
                    console.log(remoteResponse);
                    return remoteResponse;
                  }
                }
              });
            } else if ($scope.defObj.params[i].type == "select"){
              $('#runopt-field-' + i).dropdown();
            }
          }
        }
        $scope.loadingParams = false;
      }

      $scope.fileChange = function(event){
        var files = event.target.files;
        var reader = new FileReader();
        reader.readAsBinaryString(files[0]);
        reader.onload = function(e){
            $scope.defObj.params[parseInt(event.target.id.split("-")[2])].data = e.target.result;
        }
      }

      //returns the namespace of user-selected build parameters. Defaults are set for untouched fields.
      self.getBuildParamsValues = function(){
        var parameters = {};
        if ($scope.defObj.params){
          for (var i = 0; i < $scope.defObj.params.length; i++)
          {
            if ($scope.defObj.params[i].type == "text"){
              parameters[$scope.defObj.params[i].varname] = $scope.defObj.params[i].default;
            }
            if ($scope.defObj.params[i].type == "file"){
              parameters[$scope.defObj.params[i].varname] = $scope.defObj.params[i].data;
            }
            if ($scope.defObj.params[i].type == "check"){
              parameters[$scope.defObj.params[i].varname] = $scope.defObj.params[i].default ? "true" : "false";
            }
            if ($scope.defObj.params[i].type == "branchselect"){
              parameters[$scope.defObj.params[i].varname] = $('#runopt-field-' + i).dropdown('get value');
              if (parameters[$scope.defObj.params[i].varname] == '' && $scope.defObj.params[i].default){
                parameters[$scope.defObj.params[i].varname] = $scope.defObj.params[i].default
              }
            }
            if ($scope.defObj.params[i].type == "select"){
              parameters[$scope.defObj.params[i].varname] = $('#runopt-field-' + i).dropdown('get value');
              if (parameters[$scope.defObj.params[i].varname] == '' && $scope.defObj.params[i].default){
                parameters[$scope.defObj.params[i].varname] = $scope.defObj.params[i].default
              }
            }
          }
        }
        console.log("buildParameters: ", parameters);
        return parameters;
      }


      //callbacks called when the user presses a button
      self.submitPressed = function(){
        var tagArray = $('#runOptionsModal-tagsDropdown').dropdown('get value').split(",")
        var isPhysDisabled = $("#runOptionsModal-disablephys").prop('checked');
        var version = $scope.defObj['last-version'];
        $('#runOptionsModal').modal('hide');
        $rootScope.$broadcast('runOptionsModal-finished', {version: version, isPhysDisabled: isPhysDisabled, tags: tagArray, defObj: $scope.defObj, params: self.getBuildParamsValues()});
      }
      self.cancelPressed = function(){
        $('#runOptionsModal').modal('hide');
      }
    }



    angular.module('baseApp').directive('fileChange', function() {
      return {
        restrict: 'A',
        link: function (scope, element, attrs) {
          var onChangeHandler = scope.$eval(attrs.fileChange);
          element.bind('change', onChangeHandler);
        }
      };
    });

})();
