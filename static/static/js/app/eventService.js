(function () {

  angular.module('baseApp')
    .service('eventService', ['$rootScope', function($rootScope) {
      var self = this;
      var location = window.location;
      self.connected = false;
      var ws = new WebSocket("ws://" + location.hostname+(location.port ? ':'+location.port: '') + "/ws/events");

      ws.onopen = function()
      {
        console.log("Events ws opened.");
        $rootScope.$apply(function(){
          self.connected = true;
          $rootScope.$broadcast('ws-events-connected');
        });
      };

      ws.onmessage = function (evt)
      {
        var received_msg = evt.data;
        $rootScope.$apply(function(){
          var msg = JSON.parse(evt.data);
          console.log(msg);
          if (msg.Type == "RUN-STARTED"){
            $rootScope.$broadcast('ws-events-run-started', {run: msg.Data, index: msg.Index});
            $rootScope.$broadcast('ws-event-run-start-' + msg.Index);
          } else if (msg.Type == "RUN-FINISHED"){
            $rootScope.$broadcast('ws-events-run-finished', {run: msg.Data, index: msg.Index});
            $rootScope.$broadcast('ws-event-run-finish-' + msg.Index);
          }

        });
      };

      ws.onclose = function()
      {
        console.log("event ws closed.");
        $rootScope.$apply(function(){
          self.connected = false;
          $rootScope.$broadcast('ws-events-closed');
        });
      };

      ws.onerror = function(ev){
        console.log("event ws error:", ev);
      }
    }]);
})();
