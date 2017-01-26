(function() {

    var app = angular.module('baseApp', ['ngRoute', 'angularTreeview']);

    //routing
    app.config(['$routeProvider',
      function($routeProvider) {
        $routeProvider.when('/', {templateUrl: '/static/views/dash.html'});
        $routeProvider.when('/definition/:defID', {templateUrl: '/static/views/def.html', controller: 'defViewController'});
        $routeProvider.when('/edit/definition/:defID/:name', {templateUrl: '/static/views/editdef.html', controller: 'editDefController'});
        $routeProvider.when('/edit/file/:defID/:path*', {templateUrl: '/static/views/editbasefile.html', controller: 'editFileController'});
        $routeProvider.when("/documentation", {templateUrl: "/documentation/readme"});
        $routeProvider.when("/cron", {templateUrl: "/static/views/cron.html", controller: 'cronController'});
        $routeProvider.when("/browser/:path*?", {templateUrl: "/static/views/browser.html", controller: 'browserController'});
    }]);
})();
