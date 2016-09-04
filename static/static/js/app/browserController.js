(function () {

    angular.module('baseApp')
        .controller('browserController', ['$scope', 'dataService', '$location', '$routeParams', '$rootScope', '$http', browserController]);

    function browserController($scope, dataService, $location, $routeParams, $rootScope, $http) {
      var self = this;
      $scope.path = $routeParams.path;
      $scope.loading = true;
      $scope.treedata = [
        { "label" : "Loading ... Please Wait.", "id" : "main:build", "type": "file", "children" : []}
      ];

      $scope.convertUnixBitstoPermsString = function(mode) {
        if (mode == "" || mode == undefined){
          return "-";
        }

        var s = [];
        for (var i = 2; i >= 0; i--) {
          s.push((mode >> i * 3) & 4 ? 'r' : '-');
          s.push((mode >> i * 3) & 2 ? 'w' : '-');
          s.push((mode >> i * 3) & 1 ? 'x' : '-');
        }
        // optional
        if ((mode >> 9) & 4) // setuid
          s[2] = s[2] === 'x' ? 's' : 'S';
        if ((mode >> 9) & 2) // setgid
          s[5] = s[5] === 'x' ? 's' : 'S';
        if ((mode >> 9) & 1) // sticky
          s[8] = s[8] === 'x' ? 't' : 'T';
        return s.join('');
      }

      $scope.canEditSelection = function(){
        if ($scope.tree.currentNode && $scope.tree.currentNode.type == 'file'){
          return ($scope.tree.currentNode.id.match(/\/definitions\//g) || []).length || ($scope.tree.currentNode.id.match(/\/base\//g) || []).length;
        }else {
          return false;
        }
      }

      $scope.canNewFolderSelection = function(){
        if ($scope.tree.currentNode && $scope.tree.currentNode.type == 'folder'){
          return ($scope.tree.currentNode.id.match(/\/base/g) || []).length;
        }else {
          return false;
        }
      }

      $scope.canNewDefFileSelection = function(){
        if ($scope.tree.currentNode && $scope.tree.currentNode.type == 'folder'){
          return ($scope.tree.currentNode.id.match(/\/definitions/g) || []).length;
        }else {
          return false;
        }
      }

      $scope.canNewFileSelection = function(){
        if ($scope.tree.currentNode && $scope.tree.currentNode.type == 'folder'){
          return ($scope.tree.currentNode.id.match(/\/base/g) || []).length;
        }else {
          return false;
        }
      }

      $scope.canDownloadFileSelection = function(){
        if ($scope.tree.currentNode && $scope.tree.currentNode.type == 'file'){
          return ($scope.tree.currentNode.id.match(/\/build\//g) || []).length;
        }else {
          return false;
        }
      }

      $scope.canDeleteFileSelection = function(){
        if ($scope.tree.currentNode){
          return ($scope.tree.currentNode.id.match(/\/base\//g) || []).length;
        }else {
          return false;
        }
      }

      $scope.delete = function(){
        console.log($scope.tree.currentNode.id);
        if (confirm("Are you sure you want to delete: " + $scope.tree.currentNode.id + "?")){
          $http.get("/api/file/delete?path=" + $scope.tree.currentNode.id, {}).then(function (response) {
            console.log(response);
            $scope.treedata = [{ "label" : "Loading ... Please Wait.", "id" : "main:build", "type": "file", "children" : []}];
            self.update()
          }, function errorCallback(response) {
            console.log(response);
          });
        }
      }

      $scope.download = function(){
        var URL = "/api/file/download/workspace?path=" + $scope.tree.currentNode.id.replace(/^\/build\//, '');
        var win = window.open(URL, '_blank');
      }

      $scope.newFile = function(){
        var fileName = prompt("Please enter the name of the new file.", "");
        if (fileName == null || fileName == ""){
          return;
        }

        console.log($scope.tree.currentNode.id + "/" + fileName);
        $http.get("/api/file/new/file?path=" + $scope.tree.currentNode.id + "/" + fileName, {}).then(function (response) {
          console.log(response);
          if (response.data.success){
            var newNode = { "label" : fileName, "id" : $scope.tree.currentNode.id + "/" + fileName, "type": "file", "media": "-"};
            if ($scope.tree.currentNode.children){
              $scope.tree.currentNode.children[$scope.tree.currentNode.children.length] = newNode;
            } else {
              $scope.tree.currentNode.children = [newNode];
            }
          }
        }, function errorCallback(response) {
          console.log(response);
        });
      }

      ///api/file/new/definition
      $scope.newDefinition = function(){
        var fileName = prompt("Please enter the name of the new definition file.", "");
        if (fileName == null || fileName == ""){
          return;
        }

        console.log($scope.tree.currentNode.id + "/" + fileName);
        $http.get("/api/file/new/definition?path=" + $scope.tree.currentNode.id + "/" + fileName, {}).then(function (response) {
          console.log(response);
          if (response.data.success){
            var newNode = { "label" : fileName, "id" : $scope.tree.currentNode.id + "/" + fileName, "type": "file", "media": "JSON definition file"};
            if ($scope.tree.currentNode.children){
              $scope.tree.currentNode.children[$scope.tree.currentNode.children.length] = newNode;
            } else {
              $scope.tree.currentNode.children = [newNode];
            }
          }
        }, function errorCallback(response) {
          console.log(response);
        });
      }


      $scope.newFolder = function(){
        var folderName = prompt("Please enter the name of the new folder.", "");
        if (folderName == null || folderName == ""){
          return;
        }

        console.log($scope.tree.currentNode.id + "/" + folderName);
        $http.get("/api/file/new/folder?path=" + $scope.tree.currentNode.id + "/" + folderName, {}).then(function (response) {
          console.log(response);
          if (response.data.success){
            var newNode = { "label" : folderName, "id" : $scope.tree.currentNode.id + "/" + folderName, "type": "folder", "media": "-", "collapsed": true};
            if ($scope.tree.currentNode.children){
              $scope.tree.currentNode.children[$scope.tree.currentNode.children.length] = newNode;
            } else {
              $scope.tree.currentNode.children = [newNode];
            }
          }
        }, function errorCallback(response) {
          console.log(response);
        });
      }

      $scope.edit = function(){
        if ($scope.tree.currentNode && $scope.tree.currentNode.type == 'file'){
          if (($scope.tree.currentNode.id.match(/\/base\//g) || []).length){
            $location.path("/edit/file/-1/" + $scope.tree.currentNode.id.replace(/^\/base\//, ''));
          }
          if (($scope.tree.currentNode.id.match(/\/definitions\//g) || []).length){
            $location.path("/edit/definition/-1/" + $scope.tree.currentNode.id.replace(/^\/definitions\//, ''));
          }
        }
      }

      self.update = function(){
        dataService.getBrowserFileData(function(data){
          $scope.treedata = [
            { "label" : "Run Definitions", "id" : "/definitions", "type": "folder", "children" : data.definitions, "collapsed": true, "media": "-"},
            { "label" : "Build workspace", "id" : "/build", "type": "folder", "children" : data.build, "collapsed": true, "media": "-"},
            { "label" : "Base folders", "id" : "/base", "type": "folder", "children" : data.base, "collapsed": true, "media": "-"}
          ];
        })
      }
      self.update();


    }



    angular.module('baseApp').filter('bytes', function() {
    	return function(bytes, precision) {
    		if (isNaN(parseFloat(bytes)) || !isFinite(bytes)) return '-';
    		if (typeof precision === 'undefined') precision = 1;
    		var units = ['bytes', 'kB', 'MB', 'GB', 'TB', 'PB'],
    			number = Math.floor(Math.log(bytes) / Math.log(1024));
    		return (bytes / Math.pow(1024, Math.floor(number))).toFixed(precision) +  ' ' + units[number];
    	}
    });

})();
