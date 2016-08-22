(function () {

    angular.module('baseApp')
        .controller('editDefController', ['$scope', 'dataService', '$location', '$routeParams', '$rootScope', editDefController]);

    function editDefController($scope, dataService, $location, $routeParams, $rootScope) {
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

      dataService.getDefinitionFile($routeParams.defID, function(data){
        console.log(data);
        self.editor.setValue(data, -1);
        $scope.loading = false;
      });
    }

})();
