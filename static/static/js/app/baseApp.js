(function() {

    var app = angular.module('baseApp', ['ngRoute']);

    //routing
    app.config(['$routeProvider',
      function($routeProvider) {
        $routeProvider.when('/', {templateUrl: '/static/views/dash.html'});
        //$routeProvider.when('/admin/entities', {templateUrl: '/view/entities', controller: 'entityViewerAdminController'});
    }]);
})();
