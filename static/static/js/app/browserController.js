(function () {

    angular.module('baseApp')
        .controller('browserController', ['$scope', 'dataService', '$location', '$routeParams', '$rootScope', browserController]);

    function browserController($scope, dataService, $location, $routeParams, $rootScope) {
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

      $scope.canNewFileSelection = function(){
        if ($scope.tree.currentNode && $scope.tree.currentNode.type == 'folder'){
          return !($scope.tree.currentNode.id.match(/\/build/g) || []).length;
        }else {
          return false;
        }
      }

      dataService.getBrowserFileData(function(data){
        $scope.treedata = [
          { "label" : "Run Definitions", "id" : "/definitions", "type": "folder", "children" : data.definitions, "collapsed": true, "media": "-"},
          { "label" : "Build workspace", "id" : "/build", "type": "folder", "children" : data.build, "collapsed": true, "media": "-"},
          { "label" : "Base folders", "id" : "/base", "type": "folder", "children" : data.base, "collapsed": true, "media": "-"}
        ];
      })


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
