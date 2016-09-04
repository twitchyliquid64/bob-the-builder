(function () {

    angular.module('baseApp')
        .controller('editDefController', ['$scope', 'dataService', '$location', '$routeParams', '$rootScope', '$http', editDefController]);

    function editDefController($scope, dataService, $location, $routeParams, $rootScope, $http) {
      var self = this;
      $scope.defID = $routeParams.defID;
      $scope.name = $routeParams.name;
      $scope.loading = true;

      var editor = ace.edit("editor");
      self.editor = editor;
      //editor.setTheme("ace/theme/monokai");
      editor.session.setMode("ace/mode/json");
      editor.renderer.setScrollMargin(10, 10);
      editor.setOptions({
          autoScrollEditorIntoView: true,
          maxLines: Infinity,
          minLines: 38
      });
      editor.resize();

      $scope.back = function(){
        $location.path("/definition/" + $routeParams.defID);
      }

      $scope.save = function(){
        $scope.back()
        dataService.saveDefinitionFile($scope.defID, self.editor.getValue());
      }




      self.cancelPressed = function(){
        $('#miniDocumentationModal').modal('hide');
      }
      $scope.openDocs = function(){
        $('#miniDocumentationModal').modal({//setup button callbacks + general parameters
          closable: true,
          autofocus: true,
          dimmerSettings: {
            closable: false,
            opacity: 0,
          },
          onDeny: self.cancelPressed
        });
        $('#miniDocumentationModal').modal('show');
        $('#documentation-defedit').accordion({
          onOpen: function(){
            $('#miniDocumentationModal').modal('refresh');
          }
        });
      }


      if ($routeParams.defID != "-1"){ //defID known
        dataService.getDefinitionFile($routeParams.defID, function(data){
          console.log(data);
          self.editor.setValue(data, -1);
          $scope.loading = false;
        });
      } else {
        $http.get("http://localhost:8010/api/definition/getIdByName?fname="+$routeParams.name).then(function (response) {
          $location.path("/edit/definition/" + response.data + "/" + $routeParams.name);
        }, function errorCallback(response) {
          console.log(response);
          $scope.loading = false;
        });
      }

    }

})();
