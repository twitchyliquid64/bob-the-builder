(function() {

    var app = angular.module('baseApp', ['ngRoute']);

    //routing
    app.config(['$routeProvider',
      function($routeProvider) {
        $routeProvider.when('/', {templateUrl: '/static/views/dash.html'});
        $routeProvider.when('/definition/:defID', {templateUrl: '/static/views/def.html', controller: 'defViewController'});
        $routeProvider.when('/edit/definition/:defID/:name', {templateUrl: '/static/views/editdef.html', controller: 'editDefController'});
        $routeProvider.when('/edit/file/:defID/:path*', {templateUrl: '/static/views/editbasefile.html', controller: 'editFileController'});
        $routeProvider.when("/documentation", {templateUrl: "/documentation/readme"});
        //$routeProvider.when('/admin/entities', {templateUrl: '/view/entities', controller: 'entityViewerAdminController'});
    }]);
})();
