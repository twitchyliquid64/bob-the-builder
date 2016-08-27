(function () {

    angular.module('baseApp')
        .controller('editFileController', ['$scope', 'dataService', '$location', '$routeParams', '$rootScope', editFileController]);

    function editFileController($scope, dataService, $location, $routeParams, $rootScope) {
      var self = this;
      $scope.defID = $routeParams.defID;
      $scope.path = $routeParams.path;
      $scope.loading = true;

      var editor = ace.edit("editor");
      self.editor = editor;
      //editor.setTheme("ace/theme/monokai");

      if ($scope.path.substr($scope.path.length-5) === ".json")
        editor.session.setMode("ace/mode/json");
      if ($scope.path.substr($scope.path.length-3) === ".sh")
        editor.session.setMode("ace/mode/nix");

      editor.renderer.setScrollMargin(10, 10);
      editor.setOptions({
          autoScrollEditorIntoView: true,
          maxLines: Infinity,
          minLines: 38
      });
      editor.resize();

      $scope.back = function(){
        if ($scope.defID >= 0)
          $location.path("/definition/" + $routeParams.defID);
        else
          $location.path("/browser/");
      }

      $scope.isFromDefinition = function(){
        return $scope.defID >= 0;
      }

      $scope.save = function(){
        $scope.back()
        dataService.saveFile($scope.path, self.editor.getValue());
      }

      dataService.getFile($routeParams.path, function(data){
        console.log(data);
        self.editor.setValue(data, -1);
        $scope.loading = false;
      });
    }

})();
